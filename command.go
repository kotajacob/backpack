// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
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
	fmt.Println(opts)

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

	if opt, ok := optMap["buy"]; ok {
		operated = true
		msg := b.buyItem(
			opt.StringValue(),
			owner,
			fmt.Sprintf("<#%v>", m.ChannelID),
		)
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
