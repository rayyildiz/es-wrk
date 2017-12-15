package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	es "github.com/olivere/elastic"
	"github.com/olivere/elastic/config"
	"math/rand"
	"reflect"
	"strings"
)

// Worker represent the es-wrk.
type Worker struct {
	client    *es.Client
	generator *dataGenerator
	ttype     string
}

// NewWorker creates a worker to generate random data  and insert into the elasticsearch.
func NewWorker(url, username, password string, ttype reflect.Type) (*Worker, error) {
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

	client, err := es.NewClientFromConfig(&cfg)

	if err != nil {
		return nil, fmt.Errorf("could not connect to elastic, %v", err)
	}
	ctx := context.Background()
	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		// Handle error
		return nil, fmt.Errorf("error getting info, %v", err)
	}
	log.Printf("Elasticsearch returned with code %d. Cluster name %s, version %s and node name %s", code, info.ClusterName, info.Version.Number, info.Name)

	generator, err := NewGenerator(ttype)
	if err != nil {
		return nil, fmt.Errorf("generator, %v", err)
	}

	w := Worker{client, generator, strings.ToLower(ttype.Name())}

	return &w, nil
}

// Start is entrypoint for worker.
func (w *Worker) Start(numberOfElements int) {
	w.PrintInfo()

	start := time.Now()
	log.Print("[INFO] Generating random data")
	elems := w.generator.GetRandomElements(numberOfElements)
	diff := time.Now().Sub(start)
	log.Printf("[INFO] %d random data generated in %v", len(elems), diff)

	start = time.Now()
	//w.insert(numberOfWorker, posts)
	w.insertBulk(elems)
	diff = time.Now().Sub(start)
	log.Printf("[INFO] %d random data inserted into the %s in %v", len(elems), w.client.String(), diff)

	w.PrintInfo()
}

// PrintInfo prints the count of data in elasticsearch.
func (w *Worker) PrintInfo() {
	ctx := context.Background()
	result, err := w.client.Count(w.ttype).Do(ctx)
	if err != nil {
		log.Printf("[ERROR] could not search data for %s, %v", w.ttype, err)
	} else {
		log.Printf("[INFO] Total %s %d", w.ttype, result)
	}

}

func (w *Worker) insert(numberOfWorker int, elements []reflect.Value) error {
	ctx := context.Background()

	sem := make(chan bool, numberOfWorker)

	for _, elem := range elements {
		sem <- true

		go func(elem reflect.Value) {
			defer func() { <-sem }()

			_, err := w.client.Index().
				Index(w.ttype).
				Type(w.ttype).
				Id(randString(48)).
				BodyJson(elem).
				Do(ctx)

			if err != nil {
				log.Printf("[ERROR] could not insert, %v", err)
			}

		}(elem)
	}

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	return nil
}

func (w *Worker) insertBulk(elements []interface{}) error {
	ctx := context.Background()

	log.Printf("[INFO] inserting type %s", w.ttype)
	bulkService := w.client.Bulk()
	for i, elem := range elements {
		req := es.NewBulkIndexRequest().
			Index(w.ttype).
			Type(w.ttype).
			OpType("create").
			Id(randString(64)).
			Doc(elem)

		bulkService.Add(req)

		if i%5000 == 1 {
			log.Printf("[INFO] estimated size in bytes %d", bulkService.EstimatedSizeInBytes())
			resp, err := bulkService.Do(ctx)
			if err != nil {
				log.Printf("[ERROR] could not insert to elastic search, %v", err)
			}
			if resp != nil {
				log.Printf("[INFO] %d created random data ", len(resp.Created()))

				bulkService = w.client.Bulk()
			}
		}

	}

	if bulkService.NumberOfActions() > 0 {

		log.Printf("[INFO] estimated size in bytes %d", bulkService.EstimatedSizeInBytes())
		resp, err := bulkService.Do(ctx)
		if err != nil {
			log.Printf("[ERROR] could not insert to elastic search, %v", err)
		}

		if resp != nil {
			log.Printf("[INFO] %d created posts ", len(resp.Created()))
		}
	}

	_, err := w.client.Flush().Index(w.ttype).Do(ctx)
	if err != nil {
		log.Printf("[ERROR] could not flush %v", err)
	}

	return nil
}

const constLetters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ "

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = constLetters[rand.Intn(len(constLetters))]
	}
	return string(b)
}
