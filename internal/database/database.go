package database

import (
	"math/rand"
	"strconv"
	"time"
)

type Entry struct {
	Link      string
	createdAt time.Time
	expiresIn time.Duration
}

type Database map[string]Entry

var now = time.Now

func New() Database {
	d := make(Database)
	return d
}

func (d Database) Insert(link string, expiresIn time.Duration) string {
	var slug string

	for {
		slug = strconv.FormatUint(rand.Uint64(), 36)
		_, ok := d[slug]
		if !ok {
			break
		}
	}

	d[slug] = Entry{Link: link, createdAt: now(), expiresIn: expiresIn}
	return slug
}

func (d Database) InsertFixed(slug, link string, expiresIn time.Duration) Database {
	d[slug] = Entry{Link: link, createdAt: now(), expiresIn: expiresIn}
	return d
}

func (d Database) Get(slug string) (string, bool) {
	entry, ok := d[slug]
	if !ok {
		return "", false
	}

	return entry.Link, true
}

func (d Database) Purge() int {
	purged := 0

	for k, v := range d {
		if now().After(v.createdAt.Add(v.expiresIn)) {
			purged++
			delete(d, k)
		}
	}

	return purged
}
