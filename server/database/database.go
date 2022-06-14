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

func (d Database) Purge() int {
	purged := 0

	for k, v := range d {
		if v.createdAt.Add(v.expiresIn).Before(now()) {
			purged++
			delete(d, k)
		}
	}

	return purged
}
