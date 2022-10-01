// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// updateRecord updates a csv file in dir named owner.csv.
// A record always has 3 fields: count, name, price.
//
// absolute indicates that we should simply set the count instead of adding
// to the existing count.
func updateRecord(record []string, dir, owner string, absolute bool) string {
	path := filepath.Join(dir, owner+".csv")
	data, err := os.ReadFile(path)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("failed reading from %v: %v\n", path, err)
		return fmt.Sprintf("failed to add %v to %v", record, owner)
	}

	r := csv.NewReader(bytes.NewReader(data))

	// We're going to build up a new list of all the records to store. So we add
	// all the existing records whos names do not match the record we were
	// given. Then we either add our modified record.
	var records [][]string
	for {
		oldRecord, err := r.Read()
		if err == io.EOF {
			// Add the record.
			records = append(records, record)
			break
		}
		if err != nil {
			log.Printf("failed parsing from %v: %v\n", path, err)
			return fmt.Sprintf("failed to add %v to %v", record, owner)
		}

		if oldRecord[1] != record[1] {
			// Not current record.
			records = append(records, oldRecord)
			continue
		}

		if !absolute {
			// Add the count.
			count, err := strconv.Atoi(record[0])
			if err != nil {
				log.Printf("failed parsing count from %v: %v\n", path, err)
				return fmt.Sprintf("failed to add %v to %v", record, owner)
			}
			oldCount, err := strconv.Atoi(oldRecord[0])
			if err != nil {
				log.Printf("failed parsing count from %v: %v\n", path, err)
				return fmt.Sprintf("failed to add %v to %v", record, owner)
			}
			count += oldCount
			record[0] = strconv.Itoa(count)
		}

		// Use the old price if we were given NULL_PRICE.
		if record[2] == NULL_PRICE {
			record[2] = oldRecord[2]
		}
	}

	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.WriteAll(records)
	w.Flush()
	err = os.WriteFile(
		path,
		buf.Bytes(),
		0600,
	)
	if err != nil {
		log.Printf("failed writing %v to %v: %v\n", record, path, err)
		return fmt.Sprintf("failed to add %v to %v", record, owner)
	}

	s := fmt.Sprintf(
		"%v now has %v",
		owner,
		strings.Join(record[:len(record)-1], " "),
	)
	log.Println(s)
	return s
}
