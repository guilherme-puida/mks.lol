package database

import (
	"testing"
	"time"
)

func TestInsert(t *testing.T) {
	d := New()
	link := "https://example.com"
	slug := d.Insert(link, time.Second*10)

	if len(d) != 1 {
		t.Fatalf("Expected database to contain %d elements, got %d instead", 1, len(d))
	}

	if d[slug].Link != link {
		t.Fatalf("Expected entry to point to %q, got %q instead", link, d[slug].Link)
	}
}

func TestPurge(t *testing.T) {
	d := New()
	d.Insert("https://example.com", time.Second*10)
	slug := d.Insert("https://example2.com", time.Minute*10)

	now = func() time.Time {
		return time.Now().Add(time.Second * 11)
	}

	purged := d.Purge()
	if purged != 1 {
		t.Fatalf("Expected %d entries to be purged, got %d instead", 1, purged)
	}

	_, ok := d[slug]
	if !ok {
		t.Fatal("Expected one entry to still be in database, but it was purged")
	}
}
