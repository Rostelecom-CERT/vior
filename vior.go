package vior

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// InputPath is http handle path as you specified in the report-uri
	InputPath = "/csp-violation"
)

// DataStorage represents an interface for the actual reports storage
type DataStorage interface {
	Init() error               // initializes storage
	GetPipe() chan *ExtReport  // returns pipe consuming timestamped reports
	Save(csp *ExtReport) error // saves timestamped report to the storage
}

// Config represents application configuration
type Config struct {
	Storage DataStorage
}

// Request is a top struct of the CSP violation report request.
type Request struct {
	Report `json:"csp-report"`
}

// UnmarshalJSON is custom unmarshal function for the report.
// Any Report should contain at least:
//	document-uri
//	blocked-uri
//	violated-directive
//	original-policy
func (r *Request) UnmarshalJSON(data []byte) error {
	type alias Request
	req := &struct {
		*alias
	}{
		alias: (*alias)(r),
	}
	if err := json.Unmarshal(data, req); err != nil {
		return err
	}

	if req.Report.DocumentURI == "" || req.Report.BlockedURI == "" || req.Report.ViolatedDirective == "" || req.Report.OriginalPolicy == "" {
		return errors.New("wrong data supplied")
	}

	return nil
}

// Report represents Content Security Policy violation report
// Link:
// https://w3c.github.io/webappsec-csp/2/#directive-report-uri
// (8.2 Sample violation report)
type Report struct {
	DocumentURI       string `json:"document-uri"`
	Referrer          string `json:"referrer"`
	BlockedURI        string `json:"blocked-uri"`
	ViolatedDirective string `json:"violated-directive"`
	OriginalPolicy    string `json:"original-policy"`
}

// ExtReport is an extended Report with additional metadata
type ExtReport struct {
	Report
	Date     time.Time `json:"date"`
	RemoteIP net.IP    `json:"remote-ip"`
}

// ReportReceive handles fasthttp server requests
func (conf *Config) ReportReceive(ctx *fasthttp.RequestCtx) {
	var req Request

	if string(ctx.Path()) == InputPath && ctx.IsPost() {
		if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
			ctx.Error(err.Error(), fasthttp.StatusBadRequest)
		} else {
			rep := ExtReport{req.Report, time.Now().UTC(), ctx.RemoteIP()}
			conf.Storage.GetPipe() <- &rep

			ctx.SetStatusCode(fasthttp.StatusOK)
		}
	} else {
		ctx.Error("", fasthttp.StatusNotFound)
	}
}

// ListenAndServe wraps http server ListenAndServe call
func (conf *Config) ListenAndServe(endpoint string) {
	if err := fasthttp.ListenAndServe(endpoint, conf.ReportReceive); err != nil {
		conf.Shutdown()
		log.Fatal(err)
	}
}

// Shutdown closes all open resources
func (conf *Config) Shutdown() {
	close(conf.Storage.GetPipe())
}
