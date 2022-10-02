// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const NULL_PRICE = "-1"

type backpack struct {
	dir string
}

type operation uint

const (
	opAdd operation = iota
	opDel
	opBuy
)

func (op operation) String() string {
	switch op {
	case opAdd:
		return "add"
	case opDel:
		return "remove"
	case opBuy:
		return "buy"
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
		// TODO: Is owner initialized?
	} else {
		// Surround string with <# > to make it highlight as a channel on
		// discord.
		owner = fmt.Sprintf("<#%v>", m.ChannelID)
	}

	// Handle options.
	shouldPrint := true // TODO: Perhaps we should just always print?
	var response string
	if opt, ok := optMap["add"]; ok {
		shouldPrint = false
		response += b.updateItem(opAdd, owner, opt.StringValue())
	}

	if opt, ok := optMap["remove"]; ok {
		shouldPrint = false
		response += b.updateItem(opDel, owner, opt.StringValue())
	}

	if shouldPrint {
		response += b.printInventory(owner)
	}

	// Send our response.
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func (b backpack) updateItem(op operation, owner, v string) string {
	values := strings.Split(v, " ")

	// If we need to set instead of add to the record count.
	var absolute bool

	count, err := strconv.Atoi(values[0])
	if err == nil {
		// First argument is a number.
		if len(values) == 1 {
			// A single number means we're adding coins.
			values = append(values, "coins", NULL_PRICE)
		}
		values = values[1:]
	} else {
		if op == opDel {
			absolute = true
			count = 0
		} else {
			count = 1
		}
	}

	// Invert count for delete.
	if op == opDel {
		count = -count
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

	record := record{
		count: strconv.Itoa(count),
		name:  strings.Join(values, " "),
		price: price,
	}
	if err := updateRecord(record, b.dir, owner, absolute); err != nil {
		log.Println(err)
		return fmt.Sprintf("%v %v to %v failed", op, record, owner)
	}
	return fmt.Sprintf("%v %v to %v succeeded", op, record, owner)
}

func (b backpack) printInventory(owner string) string {
	s := fmt.Sprintf("printing %v's inventory\n", owner)
	log.Println(s)
	return s
}
