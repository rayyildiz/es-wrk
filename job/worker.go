package job

import (
	"context"
	"fmt"
	"log"
	"time"

	//es "gopkg.in/olivere/elastic.v5"
	//"gopkg.in/olivere/elastic.v5/config"

	es "github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
)

type Worker struct {
	client *es.Client
}

func NewWorker(url, username, password string) (*Worker, error) {

	cfg := config.Config{
		URL: url,
	}

	if len(username) > 0 && len(password) > 0 {
		cfg.Username = username
		cfg.Password = password
	}

	b := new(bool)
	*b = false

	cfg.Sniff = b

	c, err := es.NewClientFromConfig(&cfg)

	if err != nil {
		return nil, fmt.Errorf("could not connect to elastic, %v", err)
	}
	ctx := context.Background()
	info, code, err := c.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("error getting info, %v", err)
	}
	log.Printf("Elasticsearch returned with code %d. Cluster name %s, version %s and node name %s", code, info.ClusterName, info.Version.Number, info.Name)

	w := Worker{c}

	return &w, nil
}

func (w *Worker) DoJob(numberOfPost, numberOfWorker int) {
	w.PrintInfo()

	start := time.Now()
	log.Print("Generating posts")
	posts := GetRandomPosts("1000", numberOfPost)
	diff := time.Now().Sub(start)
	log.Printf("%d posts generated in %v", len(posts), diff)
	ctx := context.Background()

	start = time.Now()
	/*
			bulkService := w.client.Bulk().Index("post").Type("post")
			for _, post := range posts {
				req := es.NewBulkIndexRequest().
					Index("post").
					Type("post").
					OpType("create").
					Id(post.ID).
					Doc(post)

				bulkService = bulkService.Add(req)
			}

			ctx := context.Background()
			resp, err := bulkService.Do(ctx)
			if err != nil {
				fmt.Println(err)
			}

			if resp != nil {
				log.Printf("[INFO] bulk took %v", resp.Took)
			}

			// Flush to make sure the documents got written.
		_, err = w.client.Flush().Index("post").Do(ctx)
		if err != nil {
			log.Printf("[ERROR] could not flush %v", err)
		}
	*/

	sem := make(chan bool, numberOfWorker)

	for _, post := range posts {
		sem <- true

		go func(post Post) {
			defer func() { <-sem }()

			_, err := w.client.Index().
				Index("post").
				Type("post").
				Id(post.ID).
				BodyJson(post).
				Do(ctx)

			if err != nil {
				log.Printf("[ERROR] coudl not insert, %v", err)
			}

		}(post)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	diff = time.Now().Sub(start)
	log.Printf("%d posts inserted into the %s in %v", len(posts), w.client.String(), diff)

	w.PrintInfo()
}

func (w *Worker) PrintInfo() {
	ctx := context.Background()
	result, err := w.client.Count("post").Do(ctx)
	if err != nil {
		log.Printf("[ERROR] coudl not search post, %v", err)
	} else {
		log.Printf("[INFO] Total posts %d", result)
	}

}
