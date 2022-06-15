package main

import (
	_ "embed"
	"log"
	"net/http"
	"os"
	"strconv"

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
	_, shouldRenderStats := os.LookupEnv("MKS_SHOULD_RENDER_STATS")

	return env{url, uint(port), shouldRenderStats}
}

func main() {
	mux := http.NewServeMux()
	db := database.New()

	e := readEnv()

	h := handler.New(
		e.url,
		e.port,
		true,
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
