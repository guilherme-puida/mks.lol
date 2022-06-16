package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/guilherme-puida/mks.lol/internal/database"
	"github.com/guilherme-puida/mks.lol/internal/handler"
)

// pageTemplate should be able to render all handler.renderData struct members.
//go:embed web/page.tmpl
var pageTemplate string

//go:embed web/favicon.ico
var favicon []byte

type env struct {
	url               string
	port              uint
	shouldRenderStats bool
}

func readEnv() env {
	url := os.Getenv("MKS_URL")
	port, _ := strconv.ParseUint(os.Getenv("MKS_PORT"), 10, 64)
	shouldRenderStats := strings.ToLower(os.Getenv("MKS_SHOULD_RENDER_STATS"))

	var shouldRender bool
	switch shouldRenderStats {
	case "", "no", "false", "0":
		shouldRender = false
	default:
		shouldRender = true
	}

	return env{url, uint(port), shouldRender}
}

func main() {
	mux := http.NewServeMux()
	db := database.New()

	ticker := time.NewTicker(time.Second)

	go func() {
		for {
			select {
			case <-ticker.C:
				db.Purge()
			}
		}
	}()

	e := readEnv()

	h := handler.New(
		e.url,
		e.port,
		e.shouldRenderStats,
		db,
		pageTemplate,
		favicon,
	)

	mux.Handle("/", h)

	addr := ":" + strconv.FormatUint(uint64(e.port), 10)
	log.Printf("Starting server on port %s", addr)
	err := http.ListenAndServe(addr, mux)

	if err != nil {
		log.Fatalf("ListenAndServe fatal error: %s\n", err)
	}
}
