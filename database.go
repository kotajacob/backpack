// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
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

type records []record

// updateRecord updates a record with v in a csv file located at dir/owner.csv.
//
// absolute indicates that we should set the count instead of adding to the
// existing count.
func updateRecord(v record, dir, owner string, absolute bool) error {
	path := filepath.Join(dir, owner+".csv")
	recs, err := loadRecords(path)
	if err != nil {
		return err
	}

	var found bool
	for i := range recs {
		if recs[i].name != v.name {
			continue
		}
		found = true

		if absolute {
			recs[i].count = v.count
		} else {
			recs[i].addCount(v.count)
		}

		// Use the new price only if we were given a price.
		if v.price != NULL_PRICE {
			recs[i].price = v.price
		}
	}
	if !found {
		recs = append(recs, v)
	}

	return storeRecords(path, recs)
}

// loadRecords reads a csv file located at path and parses the contents into a
// list of records.
func loadRecords(path string) (records, error) {
	var recs records

	d, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return recs, fmt.Errorf("failed reading %v: %v\n", path, err)
	}
	r := csv.NewReader(bytes.NewReader(d))

	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return recs, fmt.Errorf("failed parsing %v: %v\n", path, err)
		}
		rec := record{
			count: line[0],
			name:  line[1],
			price: line[2],
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// storeRecords writes a list of records to a csv file at path.
func storeRecords(path string, records records) error {
	var lines [][]string
	for _, r := range records {
		line := []string{r.count, r.name, r.price}
		lines = append(lines, line)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.WriteAll(lines)
	w.Flush()
	return os.WriteFile(path, buf.Bytes(), 0600)
}
