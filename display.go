package main

import (
	"log"
	"path/filepath"
)

// displayInvetory prints out a pretty table showing owner's inventory.
func (b backpack) displayInvetory(owner string, pricedOnly bool) string {
	path := filepath.Join(b.dir, owner+".csv")
	recs, err := loadRecords(path)
	if err != nil {
		log.Printf("error displaying inventory %v: %v", owner, err)
		return FATAL_MSG
	}
	if pricedOnly {
		return recs.forSale().String()
	}
	return recs.String()
}
