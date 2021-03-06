// Package web serves the JavaScript player and WebRTC negotiation
package web

import (
	"gitlab.crans.org/nounous/ghostream/stream/ovenmediaengine"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/markbates/pkger"
	"gitlab.crans.org/nounous/ghostream/messaging"
)

// Options holds web package configuration
type Options struct {
	Enabled                     bool
	CustomCSS                   string
	Favicon                     string
	Hostname                    string
	ListenAddress               string
	Name                        string
	MapDomainToStream           map[string]string
	PlayerPoster                string
	SRTServerPort               string
	STUNServers                 []string
	ViewersCounterRefreshPeriod int
	WidgetURL                   string
	LegalMentionsEntity         string
	LegalMentionsAddress        string
	LegalMentionsFullAddress    []string
	LegalMentionsEmail          string
}

var (
	cfg *Options

	omeCfg *ovenmediaengine.Options

	// Preload templates
	templates *template.Template

	// Streams to get statistics
	streams *messaging.Streams
)

// Load templates with pkger
// templates will be packed in the compiled binary
func loadTemplates() error {
	templates = template.New("")
	return pkger.Walk("/web/template", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip non-templates
		if info.IsDir() || !strings.HasSuffix(path, ".html") {
			return nil
		}

		// Open file with pkger
		f, err := pkger.Open(path)
		if err != nil {
			return err
		}

		// Read and parse template
		temp, err := ioutil.ReadAll(f)
		if err != nil {
			return err
		}
		templates, err = templates.Parse(string(temp))
		return err
	})
}

// Serve HTTP server
func Serve(s *messaging.Streams, c *Options, ome *ovenmediaengine.Options) {
	streams = s
	cfg = c
	omeCfg = ome

	if !cfg.Enabled {
		// Web server is not enabled, ignore
		return
	}

	// Load templates
	if err := loadTemplates(); err != nil {
		log.Fatalln("Failed to load templates:", err)
	}

	// Set up HTTP router and server
	mux := http.NewServeMux()
	mux.HandleFunc("/", viewerHandler)
	mux.Handle("/static/", staticHandler())
	mux.HandleFunc("/_ws/", websocketHandler)
	mux.HandleFunc("/_stats/", statisticsHandler)
	log.Printf("HTTP server listening on %s", cfg.ListenAddress)
	log.Fatal(http.ListenAndServe(cfg.ListenAddress, mux))
}
