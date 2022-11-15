package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gertd/go-pluralize"
)

// buyItem removes an item from the seller, removes the corresponding number of
// coins from the buyer, and then adds the item to the buyer.
//
// Unlike add, set, and remove, you do not specify a price in a buy request.
// Additionally, you must always specify a name as the shorthand add/set/remove
// coin functionality is not used in buy requests.
func (b backpack) buyItem(count int, item, buyer, seller string) string {
	log.Println(buyer, "bought", count, item, "from", seller)

	// Check if buyer and seller are the same person.
	if buyer == seller {
		return fmt.Sprintf(
			"bruh. alright...\n%v bought %v %v from %v",
			buyer,
			count, item,
			seller,
		)
	}

	// Invalid items.
	if item == "" {
		return "You forgot to request an item."
	}
	if item == "coins" || item == "" {
		return "You can't buy coins silly!"
	}

	// Request count should always be greater than 0!
	if count == 0 {
		return "You can't buy 0 of an item, silly!"
	} else if count < 0 {
		fixed := strconv.Itoa(-count) + " " + item
		return fmt.Sprintf("You've requested to give away your items?\n"+
			"Try again with: \"%v\"", fixed)
	}

	// Prepare the record requests.
	plur := pluralize.NewClient()
	name := strings.Trim(plur.Singular(item), " ")
	itemFromSeller := record{
		count: -count, // Pass a negative count to seller.
		name:  name,
		price: Unchanged,
	}
	itemToBuyer := record{
		count: count, // Pass a positive count to buyer.
		name:  name,
		price: Unchanged,
	}

	var response bytes.Buffer

	// Remove item from seller.
	sellerUpdated, sellerOld, err := updateRecord(
		itemFromSeller,
		b.dir,
		seller,
		false,
	)
	if _, ok := err.(*declinedError); ok {
		// Transaction declined. Seller's doesn't have enough in stock.
		response.WriteString(fmt.Sprintf(
			"%v does not have %v in stock\n",
			seller,
			itemToBuyer,
		))
		response.WriteString("Please choose one of the following items:\n")
		response.WriteString(b.displayInvetory(seller, true))
		return response.String()
	} else if err != nil {
		// Fatal error.
		log.Println(err)
		return FatalMessage
	}

	// Ensure that the item is actually for sale!
	if sellerOld.price == NotForSale {
		// Transaction declined. Item is not for sale!
		response.WriteString(
			fmt.Sprintf("%v does not have %v for sale\n", seller, itemToBuyer),
		)
		response.WriteString("Please choose one of the following items:\n")
		response.WriteString(b.displayInvetory(seller, true))

		// Revert seller inventory change!
		_, _, err := updateRecord(sellerOld, b.dir, seller, true)
		if err != nil {
			log.Printf(
				"error in buy request %v: "+
					"item was not for sale, but unrolling transaction failed: %v\n",
				itemToBuyer,
				err,
			)
			return FatalMessage
		}
		return response.String()
	}

	// Remove coins from buyer.
	sum := itemToBuyer.count * sellerOld.price
	coinsFromBuyer := record{
		count: -sum, // Negative to subtract.
		name:  Coin,
		price: NotForSale,
	}
	coinsToSeller := record{
		count: sum,
		name:  Coin,
		price: NotForSale,
	}
	_, buyerOld, err := updateRecord(coinsFromBuyer, b.dir, buyer, false)
	if _, ok := err.(*declinedError); ok {
		// Transaction declined. Buyer doesn't have enough coins.
		response.WriteString(
			fmt.Sprintf("%v has insufficient funds\n", buyer) +
				fmt.Sprintf("%v costs $%v\n", itemToBuyer, strconv.Itoa(sum)) +
				fmt.Sprintf("%v only has %v", buyer, buyerOld),
		)

		// Revert seller inventory change!
		_, _, err := updateRecord(sellerOld, b.dir, seller, true)
		if err != nil {
			log.Printf(
				"error in buy request %v: "+
					"buyer lacked coins, but unrolling transaction failed: %v\n",
				itemToBuyer,
				err,
			)
			return FatalMessage
		}
		return response.String()
	} else if err != nil {
		// Fatal error.
		log.Println(err)
		return FatalMessage
	}

	// Give coins to seller.
	_, _, err = updateRecord(coinsToSeller, b.dir, seller, false)
	if err != nil {
		log.Printf("error in buy request %v %v: "+
			"item removed from seller, coins removed from buyer,"+
			"but failed to give coins to seller\n", count, item)
		return FatalMessage
	}

	// Give item to buyer.
	_, _, err = updateRecord(itemToBuyer, b.dir, buyer, false)
	if err != nil {
		log.Printf("error in buy request %v %v: "+
			"item removed from seller, coins removed from buyer,"+
			"but failed to give %v to buyer\n", count, item, itemToBuyer)
		return FatalMessage
	}

	response.WriteString(fmt.Sprintf(
		"%v bought %v for $%v\n",
		buyer,
		itemToBuyer,
		sum,
	))
	response.WriteString(fmt.Sprintf(
		"%v has %v\n",
		buyer,
		itemToBuyer,
	))
	response.WriteString(fmt.Sprintf(
		"%v has %v",
		seller,
		sellerUpdated,
	))
	return response.String()
}
