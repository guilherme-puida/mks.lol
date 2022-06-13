package main

import (
	"bytes"
	_ "embed"
	"flag"
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//go:embed content/index.tmpl
var tmpl string

//go:embed content/favicon.ico
var favicon []byte

var indexTemplate *template.Template

// dataEntry is the main structure stored. It contains information about the links currently stored in the database.
type dataEntry struct {
	link      string
	createdAt time.Time
	expiresIn time.Duration
}

// renderData is passed to the template when rendering a new page. It contains data about the previous operation.
type renderData struct {
	// Represents a successful operation.
	IsSuccess bool

	// Contains the short url generated.
	// Should only be set when IsSuccess is true.
	ShortUrl string

	// Represents an unsuccessful operation.
	IsError bool

	// Contains some sort of explanation about what went wrong.
	// Should only be set when IsError is true.
	ErrorMessage string

	BaseUrl string
}

type serverOptions struct {
	url       string
	port      uint
	https     bool
	startTime time.Time
}

// database is the data storage for the application. It maps a slug (a shortened link) to a dataEntry.
var database map[string]dataEntry

// durationMap maps the string representation of expiration times sent by the client in the form to an actual Go value.
var durationMap map[string]time.Duration

var options serverOptions

// insertEntry inserts an entry in the database.
func insertEntry(link string, expiresIn time.Duration) string {
	var slug string

	for {
		slug = strconv.FormatUint(rand.Uint64(), 36)
		_, ok := database[slug]
		if !ok {
			break
		}
	}

	database[slug] = dataEntry{link: link, expiresIn: expiresIn, createdAt: time.Now()}
	return slug
}

// generateShortUrl generates a complete url based on the service's base url, a slug,
// and whether the requests are being made over https or http.
func generateShortUrl(slug string) string {
	if options.https {
		return fmt.Sprintf("https://%s/%s", options.url, slug)
	}
	return fmt.Sprintf("http://%s/%s", options.url, slug)

}

// requestHandler switches on requests to determine which action should be taken.
func requestHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		switch r.URL.Path {
		case "/":
			indexTemplate.Execute(w, renderData{BaseUrl: options.url})
		case "/favicon.ico":
			http.ServeContent(w, r, "favicon.ico", options.startTime, bytes.NewReader(favicon))
		default:
			entry, ok := database[strings.TrimPrefix(r.URL.Path, "/")]
			if !ok {
				indexTemplate.Execute(w, renderData{IsError: true, ErrorMessage: "Link not found in database.",
					BaseUrl: options.url})
				return
			}
			http.Redirect(w, r, entry.link, http.StatusMovedPermanently)
		}
	} else if r.Method == http.MethodPost {
		switch r.URL.Path {
		case "/":
			r.ParseForm()
			form := r.Form
			link := form.Get("link")
			expiresIn := form.Get("expiresIn")
			slug := insertEntry(link, durationMap[expiresIn])
			indexTemplate.Execute(w, renderData{IsSuccess: true, ShortUrl: generateShortUrl(slug),
				BaseUrl: options.url})
		default:
			indexTemplate.Execute(w, renderData{IsError: true,
				ErrorMessage: fmt.Sprintf("Method %s not allowed in URL %s", r.Method, r.URL.Path),
				BaseUrl:      options.url})
		}
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// init is called before main, and used here to initialize all global variables.
func init() {
	// Initializing the random seed to the current time. This randomizes the values of the generated short URLs.
	rand.Seed(time.Now().UnixNano())

	// Initializing the database.
	database = make(map[string]dataEntry)

	// These are the expiration values that can be sent by the client.
	durationMap = map[string]time.Duration{
		"5min":  5 * time.Minute,
		"15min": 15 * time.Minute,
		"30min": 30 * time.Minute,
		"1h":    1 * time.Hour,
		"6h":    6 * time.Hour,
		"12h":   12 * time.Hour,
		"1day":  24 * time.Hour,
	}

	indexTemplate, _ = template.New("index").Parse(tmpl)
}

// deleteOldEntries handles deleting entries that have already expired.
func deleteOldEntries() {

	for k, v := range database {
		if v.createdAt.Add(v.expiresIn).Before(time.Now()) {
			delete(database, k)
		}
	}
}

func main() {
	urlFlag := flag.String("url", "mks.lol", "url used in rendered templates")
	portFlag := flag.Uint("port", 8080, "port that will listen for all requests")
	httpsFlag := flag.Bool("https", false, "use https instead of http in rendered templates")
	flag.Parse()

	options = serverOptions{url: *urlFlag, port: *portFlag, https: *httpsFlag, startTime: time.Now()}

	ticker := time.NewTicker(5 * time.Minute)

	go func() {
		for {
			select {
			case <-ticker.C:
				deleteOldEntries()
			}
		}
	}()

	http.HandleFunc("/", requestHandler)
	http.ListenAndServe(fmt.Sprintf(":%d", options.port), nil)
}
