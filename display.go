package main

import "path/filepath"

// displayInvetory prints out a pretty table showing owner's inventory.
func (b backpack) displayInvetory(owner string, pricedOnly bool) string {
	path := filepath.Join(b.dir, owner+".csv")
	recs, err := loadRecords(path)
	if err != nil {
		return FATAL_MSG
	}
	if pricedOnly {
		return recs.forSale().String()
	}
	return recs.String()
}
