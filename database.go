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

// updateRecord updates a record with v in a csv file located at dir/owner.csv.
//
// absolute indicates that we should set the count instead of adding to the
// existing count.
//
// The updated record and the old record, or an error are returned.
func updateRecord(v record, dir, owner string, absolute bool) (record, record, error) {
	var updated record
	var old record
	path := filepath.Join(dir, owner+".csv")
	recs, err := loadRecords(path)
	if err != nil {
		return updated, old, err
	}

	var found bool
	for i := range recs {
		if recs[i].name != v.name {
			continue
		}
		found = true
		old = recs[i]

		if absolute {
			recs[i].count = v.count
		} else {
			if err := recs[i].addCount(v.count); err != nil {
				return updated, old, err
			}
		}

		// Use the new price only if we were given a price.
		if v.price != UNCHANGED {
			recs[i].price = v.price
		}
		updated = recs[i]
	}
	if !found {
		if v.count < 0 {
			return updated, old, &declinedError{v}
		}

		if v.price == UNCHANGED {
			v.price = NOT_FOR_SALE
		}

		updated = v
		recs = append(recs, v)
	}

	return updated, old, storeRecords(path, recs)
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
		count, err := strconv.Atoi(line[0])
		if err != nil {
			return recs, fmt.Errorf(
				"failed parsing count as int: %v\n",
				line[0],
			)
		}
		price, err := strconv.Atoi(line[2])
		if err != nil {
			return recs, fmt.Errorf(
				"failed parsing price as int: %v\n",
				line[0],
			)
		}
		rec := record{
			count: count,
			name:  line[1],
			price: price,
		}
		recs = append(recs, rec)
	}

	return recs, nil
}

// storeRecords writes a list of records to a csv file at path.
func storeRecords(path string, records records) error {
	var lines [][]string
	for _, r := range records {
		count := strconv.Itoa(r.count)
		price := strconv.Itoa(r.price)
		line := []string{
			count,
			r.name,
			price,
		}
		lines = append(lines, line)
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.WriteAll(lines)
	w.Flush()
	return os.WriteFile(path, buf.Bytes(), 0600)
}
