package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Rostelecom-CERT/vior"
)

func main() {
	var (
		port       string
		elasticURL string
	)

	if port = os.Getenv("VIOR_PORT"); port == "" {
		port = "8080"
	}

	if elasticURL = os.Getenv("VIOR_ELASTIC_URL"); elasticURL == "" {
		elasticURL = "http://127.0.0.1:9200"
	}

	e, err := vior.NewElasticStorage(
		elasticURL,
		"csp-violations", // Index name
		"report",         // Doc type
	)
	if err != nil {
		log.Fatal(err)
	}

	v := &vior.Config{
		Storage: e,
	}

	v.ListenAndServe(fmt.Sprintf(":%s", port))
}
