// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"bytes"
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gertd/go-pluralize"
)

const FATAL_MSG = "Backpack update failed! Contact your local currator for help!"

type backpack struct {
	dir string
}

type operation uint

const (
	opAdd operation = iota
	opDel
	opSet
	opBuy
)

func (op operation) String() string {
	switch op {
	case opAdd:
		return "Added"
	case opDel:
		return "Removed"
	case opSet:
		return "Set"
	case opBuy:
		return "Bought"
	default:
		return "unknown operation"
	}
}

var invCommand = discordgo.ApplicationCommand{
	Name:        "inv",
	Description: "Manage Inventories",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "owner",
			Description: "Owner of the inventory",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "add",
			Description: "Add an item to an inventory",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "remove",
			Description: "Remove an item from an inventory",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "set",
			Description: "Set the number of an item in an inventory",
			Required:    false,
		},
		{
			Type:        discordgo.ApplicationCommandOptionString,
			Name:        "buy",
			Description: "Owner buys an item from channel's invenory",
			Required:    false,
		},
	},
}

// commandHandler is called (due to the AddHandler above) every time a new
// command is sent on any channel that the authenticated bot has access to.
func (b backpack) commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.ApplicationCommandData().Name != invCommand.Name {
		return
	}

	// Add the options to a map.
	opts := m.ApplicationCommandData().Options
	optMap := make(
		map[string]*discordgo.ApplicationCommandInteractionDataOption,
		len(opts),
	)
	for _, opt := range opts {
		optMap[opt.Name] = opt
	}

	// Figure out the owner of the inventory.
	var owner string
	if opt, ok := optMap["owner"]; ok {
		owner = opt.StringValue()
	} else {
		// Surround string with <# > to highlight it as a channel on discord.
		owner = fmt.Sprintf("<#%v>", m.ChannelID)
	}

	// Handle add, remove, and set.
	var operated bool
	var responses []string
	if opt, ok := optMap["add"]; ok {
		operated = true
		msg := b.modifyItem(opt.StringValue(), owner, opAdd)
		responses = append(responses, msg)
	}

	if opt, ok := optMap["remove"]; ok {
		operated = true
		msg := b.modifyItem(opt.StringValue(), owner, opDel)
		responses = append(responses, msg)
	}

	if opt, ok := optMap["set"]; ok {
		operated = true
		msg := b.modifyItem(opt.StringValue(), owner, opSet)
		responses = append(responses, msg)
	}

	// Print a table as fallback command if there was no operation given.
	if !operated {
		msg := b.displayInvetory(owner, false)
		responses = append(responses, msg)
	}

	// Send our response.
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: strings.Join(responses, "\n"),
		},
	})
}

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

// buyItem removes an item from the seller, removes the corresponding number of
// coins from the buyer, and then adds the item to the buyer.
//
// Unlike add, set, and remove, you do not specify a price in a buy request.
// Additionally, you must always specify a name as the shorthand add/set/remove
// coin functionality is not used in buy requests.
func (b backpack) buyItem(request, buyer, seller string) string {
	values := strings.Split(request, " ")

	// Check if the request has a count.
	count, err := strconv.Atoi(values[0])
	if err == nil {
		// First argument is a number.
		if len(values) == 1 {
			// Normally, a number means we're using coins, but you cannot buy
			// coins so reject the offer.
			return "You can't buy coins. Make sure you request an item."
		}
		values = values[1:]
	}

	// Request count should always be greater than 0!
	if count == 0 {
		if request == "" {
			return b.displayInvetory(seller, true)
		}
		return "You can't buy 0 of an item, silly!"
	} else if count < 0 {
		fixed := strconv.Itoa(-count) + " " + strings.Join(values, " ")
		return fmt.Sprintf("You've requested to give away your items?\n"+
			"Try again with: \"%v\"", fixed)
	}

	// Prepare the record requests.
	plur := pluralize.NewClient()
	name := plur.Singular(strings.Join(values, " "))
	itemFromSeller := record{
		count: -count, // Pass a negative count to seller.
		name:  name,
		price: UNCHANGED,
	}
	itemToBuyer := record{
		count: count, // Pass a positive count to buyer.
		name:  name,
		price: UNCHANGED,
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
		return FATAL_MSG
	}

	// Ensure that the item is actually for sale!
	if sellerOld.price == NOT_FOR_SALE {
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
			return FATAL_MSG
		}
		return response.String()
	}

	// Remove coins from buyer.
	sum := itemToBuyer.count * sellerOld.price
	coinsFromBuyer := record{
		count: -sum, // Negative to subtract.
		name:  COIN,
		price: NOT_FOR_SALE,
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
			return FATAL_MSG
		}
		return response.String()
	} else if err != nil {
		// Fatal error.
		log.Println(err)
		return FATAL_MSG
	}

	// Give item to buyer.
	_, _, err = updateRecord(itemToBuyer, b.dir, buyer, false)
	if err != nil {
		log.Printf("error in buy request %v: "+
			"item removed from seller, coins removed from buyer,"+
			"but failed to give %v to buyer", request, itemToBuyer)
		return FATAL_MSG
	}
	absCoins := coinsFromBuyer.count
	if absCoins < 0 {
		absCoins = -absCoins
	}
	response.WriteString(fmt.Sprintf(
		"%v bought %v for $%v\n",
		buyer,
		itemToBuyer,
		absCoins,
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

// modifyItem updates an item for a given owner based on a request.
// An appropriate message for the user will be returned.
//
// Take note that the transaction can fail due to database corruption or other
// such issues.
func (b backpack) modifyItem(request, owner string, op operation) string {
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
	name := plur.Singular(strings.Join(values, " "))

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
	log.Println(owner, strings.ToLower(op.String()), absNoPriceRec)
	return response.String()
}
