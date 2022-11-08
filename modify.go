package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gertd/go-pluralize"
)

// modifyItem updates an item for a given owner based on a request.
// An appropriate message for the user will be returned.
//
// Take note that the transaction can fail due to database corruption or other
// such issues.
func (b backpack) modifyItem(request, owner string, op operation) string {
	log.Println(owner, strings.ToLower(op.String()), request)
	values := strings.Split(request, " ")

	// If we need to set instead of add to the record count.
	var absolute bool

	// Check if the request has a count.
	count, err := strconv.Atoi(values[0])
	if err == nil {
		// First argument is a number.
		if len(values) == 1 {
			// A single number means we're adding coins.
			values = append(values, COIN, strconv.Itoa(NOT_FOR_SALE))
		}
		values = values[1:]
	} else {
		if op == opDel {
			count = 0
			op = opSet
		} else {
			count = 1
		}
	}

	switch op {
	case opDel:
		count = -count
	case opSet:
		absolute = true
	}

	// We now have count and have removed it from values if present.
	// Let's get the price.
	price := UNCHANGED
	if len(values) != 1 {
		p, err := strconv.Atoi(values[len(values)-1])
		if err == nil {
			price = p
			values = values[:len(values)-1]
		}
	}

	// Make the name singular for storage.
	plur := pluralize.NewClient()
	name := strings.Trim(plur.Singular(strings.Join(values, " ")), " ")

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
		if op == opSet {
			response.WriteString(fmt.Sprintf(
				"%v the quantity from %v to %v",
				op,
				old,
				absNoPriceRec,
			))
		} else {
			response.WriteString(fmt.Sprintf(
				"%v %v",
				op,
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
