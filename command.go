// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"bytes"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/gertd/go-pluralize"
)

const NULL_PRICE = "-1"

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

	// Handle add and remove.
	var responses []string
	if opt, ok := optMap["add"]; ok {
		msg, _ := b.modifyItem(opt.StringValue(), owner, opAdd)
		responses = append(responses, msg)
	}

	if opt, ok := optMap["remove"]; ok {
		msg, _ := b.modifyItem(opt.StringValue(), owner, opDel)
		responses = append(responses, msg)
	}

	if opt, ok := optMap["set"]; ok {
		msg, _ := b.modifyItem(opt.StringValue(), owner, opSet)
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

// modifyItem updates an item for a given owner based on a request.
//
// An appropriate message for the user along with a boolean indicating if the
// transaction was declined will be returned.
//
// Take note that the transaction could also fail due to database corruption or
// other such issues.
func (b backpack) modifyItem(request, owner string, op operation) (string, bool) {
	values := strings.Split(request, " ")

	// If we need to set instead of add to the record count.
	var absolute bool

	// Check if the request has a count.
	count, err := strconv.Atoi(values[0])
	if err == nil {
		// First argument is a number.
		if len(values) == 1 {
			// A single number means we're adding coins.
			values = append(values, "coin", NULL_PRICE)
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

	// We now have count and removed it from values if present.
	// Let's get the price.
	price := NULL_PRICE
	if len(values) != 1 {
		p, err := strconv.Atoi(values[len(values)-1])
		if err == nil {
			price = strconv.Itoa(p)
		}
		values = values[:len(values)-1]
	}

	// Make the name singular for storage.
	plur := pluralize.NewClient()
	name := plur.Singular(strings.Join(values, " "))

	rec := record{
		count: strconv.Itoa(count),
		name:  name,
		price: price,
	}

	var response bytes.Buffer
	var declined bool
	absNoPriceRec := record{
		count: strings.TrimPrefix(rec.count, "-"),
		name:  rec.name,
	}
	updated, old, err := updateRecord(rec, b.dir, owner, absolute)
	if _, ok := err.(*declinedError); ok {
		// Declined.
		declined = true
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
		return fmt.Sprintf(
			"Backpack update failed! Contant your local currator for help!",
		), true
	}

	// Add a little summary line.
	response.WriteString(fmt.Sprintf("\n%v has %v", owner, updated))
	log.Println(owner, strings.ToLower(op.String()), absNoPriceRec)
	return response.String(), declined
}
