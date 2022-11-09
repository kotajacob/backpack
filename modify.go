package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"github.com/gertd/go-pluralize"
)

// modifyItem updates an item for a given owner based on a request.
// An appropriate message for the user will be returned.
//
// Take note that the transaction can fail due to database corruption or other
// such issues.
func (b backpack) modifyItem(count, price int, item, owner, op string) string {
	log.Println(owner, op, count, item, price)

	// If we need to set instead of add to the record count.
	var absolute bool

	switch op {
	case "remove":
		count = -count
	case "set":
		absolute = true
	}

	// Make the name singular for storage.
	plur := pluralize.NewClient()
	name := strings.Trim(plur.Singular(item), " ")

	rec := record{
		count: count,
		name:  name,
		price: price,
	}

	var response bytes.Buffer
	absCount := count
	if absCount < 0 {
		absCount = -absCount
	}
	absNoPriceRec := record{
		count: absCount,
		name:  rec.name,
		price: UNCHANGED,
	}
	updated, old, err := updateRecord(rec, b.dir, owner, absolute)
	if _, ok := err.(*declinedError); ok {
		// Declined.
		response.WriteString(fmt.Sprintf(
			"%v does not have %v to remove",
			owner,
			absNoPriceRec,
		))
	} else if err == nil {
		switch op {
		case "set":
			response.WriteString(fmt.Sprintf(
				"Set the quantity from %v to %v",
				old,
				absNoPriceRec,
			))
		case "add":
			response.WriteString(fmt.Sprintf(
				"Added %v",
				absNoPriceRec,
			))
		case "remove":
			response.WriteString(fmt.Sprintf(
				"Removed %v",
				absNoPriceRec,
			))
		}
	} else {
		// Fatal error.
		log.Println(err)
		return FATAL_MSG
	}

	// Add a little summary line.
	response.WriteString(fmt.Sprintf("\n%v has %v", owner, updated))
	return response.String()
}
