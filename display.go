package main

import (
	"log"
	"path/filepath"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gertd/go-pluralize"
)

// displayInvetory returns a pretty table showing owner's inventory.
func (b backpack) displayInvetory(owner string, pricedOnly bool) string {
	path := filepath.Join(b.dir, owner+".csv")
	recs, err := loadRecords(path)
	if err != nil {
		log.Printf("error displaying inventory %v: %v\n", owner, err)
		return FatalMessage
	}
	if pricedOnly {
		return recs.forSale().String()
	}
	return recs.String()
}

// displayName capitalizes the first letter of the first word in an item's name.
func displayName(name string, count int) string {
	r, size := utf8.DecodeRuneInString(name)
	if r == utf8.RuneError || size == 0 {
		return name
	}
	name = string(unicode.ToTitle(r)) + name[1:]

	plur := pluralize.NewClient()
	return plur.Pluralize(name, count, false)
}

// normalizeName lowercases name and converts it to its singular representation.
func normalizeName(name string) string {
	plur := pluralize.NewClient()
	return strings.Trim(plur.Singular(name), " ")
}
