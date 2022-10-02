package main

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/dustin/go-humanize"
	"github.com/gertd/go-pluralize"
)

// record represents a single entry in an inventory.
type record struct {
	count string
	name  string
	price string
}

// addCount adds to a record's count.
func (r *record) addCount(count string) error {
	current, err := strconv.Atoi(r.count)
	if err != nil {
		return fmt.Errorf(
			"failed parsing count of %v: %v\n",
			r,
			err,
		)
	}
	addend, err := strconv.Atoi(count)
	if err != nil {
		return fmt.Errorf(
			"failed parsing count %v: %v\n",
			count,
			err,
		)
	}
	sum := current + addend
	r.count = strconv.Itoa(sum)
	return nil
}

// String prints out a pretty message describing the record.
func (r record) String() string {
	var buf bytes.Buffer

	// Ignoring error to use 0 as fallback count.
	count, _ := strconv.Atoi(r.count)
	plur := pluralize.NewClient()
	buf.WriteString(plur.Pluralize(r.name, count, true))

	price, err := strconv.Atoi(r.price)
	if r.price != NULL_PRICE && err == nil {
		buf.WriteString(" for sale for $")
		buf.WriteString(humanize.Comma(int64(price)))
	}

	return buf.String()
}

type records []record
