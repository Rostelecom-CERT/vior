package vior

import (
	"context"
	"log"
	"strings"

	"github.com/olivere/elastic"
)

// ElasticStorage is an example of the storage for CSP reports,
// that implements DataStorage interface
type ElasticStorage struct {
	Pipe    chan *ExtReport
	URL     string
	Client  *elastic.Client
	IdxName string
	DocType string
	Ctx     context.Context
}

// NewElasticStorage bootstraps and initializes ElasticStorage
func NewElasticStorage(url string, idxname string, doctype string) (*ElasticStorage, error) {
	e := &ElasticStorage{
		Pipe:    make(chan *ExtReport),
		URL:     url,
		IdxName: idxname,
		DocType: doctype,
		Ctx:     context.Background(),
	}

	if err := e.Init(); err != nil {
		return nil, err
	}

	return e, nil
}

// Init initializes Elastic client, creates index and starts
// goroutine that pops reports from incomming channel
func (e *ElasticStorage) Init() error {
	ec, err := elastic.NewClient(
		elastic.SetURL(e.URL),
	)
	if err != nil {
		return err
	}
	e.Client = ec

	_, err = e.Client.CreateIndex(e.IdxName).Do(e.Ctx)
	if !strings.Contains(err.Error(), "already_exists_exception") != nil {
		return err
	}

	go func() {
		for rep := range e.Pipe {
			go e.Save(rep)
		}
	}()

	return nil
}

// Save saves the report in Elastic
func (e *ElasticStorage) Save(r *ExtReport) error {
	_, err := e.Client.Index().
		Index(e.IdxName).
		Type(e.DocType).
		BodyJson(*r).
		Do(e.Ctx)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetPipe returns a pipe to write reports to
func (e *ElasticStorage) GetPipe() chan *ExtReport {
	return e.Pipe
}
