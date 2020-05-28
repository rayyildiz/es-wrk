package main

import (
	"flag"
	"log"
	"os"
	"reflect"

	"go.rayyildiz.dev/eswrk/worker"
)

func main() {
	flags := flag.NewFlagSet("es-wrk", flag.ExitOnError)
	noOfPost := flags.Int("n", 20000, "Number of documents")
	esURL := flag.String("url", "http://localhost:9200", "Elasticsearch URL")
	esUsername := flag.String("username", "", "Elasticsearch username")
	esPassword := flag.String("password", "", "Elasticsearch Password")
	flags.Parse(os.Args[1:])

	wrk, err := worker.NewWorker(*esURL, *esUsername, *esPassword, reflect.TypeOf(Article{}))
	if err != nil {
		log.Printf("[ERROR] could not created worker, %v", err)
		os.Exit(1)
	}

	wrk.Start(*noOfPost)
}

// Article is a test object.
type Article struct {
	ID            string `json:"id"`
	Text          string `json:"text"`
	Language      string `json:"language"`
	PostDate      string `json:"postDate"`
	CurrentStatus int    `json:"currentStatus"`
	IsPublished   bool   `json:"isPublished"`
}
