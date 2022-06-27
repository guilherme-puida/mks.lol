package handler

import (
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/guilherme-puida/mks.lol/internal/database"
)

type renderData struct {
	Name string

	IsSuccess bool
	ShortUrl  template.URL

	IsError      bool
	ErrorMessage string

	ShouldRenderStats bool
	Uptime            string
	DatabaseLength    int
}

type Handler struct {
	url               string
	port              uint
	shouldRenderStats bool
	db                database.Database
	pageTemplate      *template.Template
	favicon           []byte
	startTime         time.Time
}

func New(
	url string,
	port uint,
	shouldRenderStats bool,
	db database.Database,
	pageTemplate string,
	favicon []byte,
) Handler {
	t := template.Must(template.New("page").Parse(pageTemplate))

	return Handler{
		url:               url,
		port:              port,
		shouldRenderStats: shouldRenderStats,
		db:                db,
		pageTemplate:      t,
		favicon:           favicon,
		startTime:         time.Now(),
	}
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		switch r.Method {
		case http.MethodGet:
			h.renderBase(w, renderData{})
		case http.MethodPost:
			err := r.ParseForm()
			if err != nil || !r.Form.Has("link") {
				h.renderError(w, "Error parsing form.", http.StatusBadRequest)
				return
			}

			slug := h.db.Insert(r.Form.Get("link"), getDurationFromStr(r.Form.Get("expiresIn")))
			w.WriteHeader(http.StatusCreated)
			h.renderBase(w, renderData{IsSuccess: true, ShortUrl: template.URL(h.url + "/" + slug)})
		default:
			h.renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	case "/favicon.ico":
		if r.Method != http.MethodGet {
			h.renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		h.serveBytes(w, r, "favicon.ico", h.favicon)
	default:
		if r.Method != http.MethodGet {
			h.renderError(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		link, ok := h.db.Get(strings.TrimPrefix(r.URL.Path, "/"))
		if !ok {
			h.renderError(w, "Link not found.", http.StatusNotFound)
			return
		}

		http.Redirect(w, r, link, http.StatusFound)
	}
}

func (h Handler) renderBase(w http.ResponseWriter, data renderData) {
	if data.Name == "" {
		data.Name = h.url
	}

	data.ShouldRenderStats = h.shouldRenderStats
	data.DatabaseLength = len(h.db)
	data.Uptime = time.Since(h.startTime).Round(time.Second).String()

	err := h.pageTemplate.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (h Handler) renderError(w http.ResponseWriter, message string, status int) {
	w.WriteHeader(status)
	h.renderBase(w, renderData{IsError: true, ErrorMessage: message})
}

func (h Handler) serveBytes(w http.ResponseWriter, r *http.Request, name string, content []byte) {
	http.ServeContent(w, r, name, h.startTime, bytes.NewReader(content))
}

func getDurationFromStr(str string) time.Duration {
	duration, ok := map[string]time.Duration{
		"5min":  5 * time.Minute,
		"15min": 15 * time.Minute,
		"30min": 30 * time.Minute,
		"1h":    1 * time.Hour,
		"6h":    6 * time.Hour,
		"12h":   12 * time.Hour,
		"1day":  24 * time.Hour,
		"1week": 7 * 24 * time.Hour,
	}[str]

	if !ok {
		return 5 * time.Minute
	}

	return duration
}
