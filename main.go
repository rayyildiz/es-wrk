package main

import (
	"flag"
	"log"
	"os"

	"github.com/rayyildiz/es-wrk/job"
	"github.com/subosito/gotenv"
)

func init() {
	gotenv.Load()
}

func main() {
	cfg := loadConfig()

	flags := flag.NewFlagSet("es-wrk", flag.ExitOnError)
	noOfPost := flags.Int("n", 500, "Number of posts")
	noOfWorker := flags.Int("t", 8, "Number of thread")
	flags.Parse(os.Args[1:])

	wrk, err := job.NewWorker(cfg.url, cfg.username, cfg.password)
	if err != nil {
		log.Printf("[ERROR] could not created worker, %v", err)
		os.Exit(1)
	}

	wrk.DoJob(*noOfPost, *noOfWorker)
}

type config struct {
	url      string
	username string
	password string
}

func loadConfig() *config {
	cfg := config{}

	cfg.url = os.Getenv("ELASTICSEARCH_URL")
	cfg.username = os.Getenv("ELASTICSEARCH_USERNAME")
	cfg.password = os.Getenv("ELASTICSEARCH_PASSWORD")

	return &cfg
}
