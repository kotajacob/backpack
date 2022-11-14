// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"fmt"
	"strconv"

	"github.com/bwmarrin/discordgo"
)

const FATAL_MSG = "Backpack failed! Contact your local currator for help!"

type backpack struct {
	dir string
}

var invCommand = discordgo.ApplicationCommand{
	Name:        "inv",
	Description: "Manage Inventories",
	Options: []*discordgo.ApplicationCommandOption{
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "view",
			Description: "View an inventory",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "owner",
					Description: "Whose inventory to view",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "add",
			Description: "Add an item",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "owner",
					Description: "Whose inventory to add to",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "quantity",
					Description: "The number of items to add",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The name of the item to add",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "price",
					Description: "The price of the item to add",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "remove",
			Description: "Remove an item",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "owner",
					Description: "Whose inventory to remove from",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "quantity",
					Description: "The number of items to remove",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The name of the item to remove",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "set",
			Description: "Set the quantity or price of an item",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "owner",
					Description: "Whose inventory to edit",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "quantity",
					Description: "The resultant number of items",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The name of the item",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "price",
					Description: "The price of the item",
					Required:    false,
				},
			},
		},
		{
			Type:        discordgo.ApplicationCommandOptionSubCommand,
			Name:        "buy",
			Description: "Buy an item",
			Required:    false,
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "buyer",
					Description: "Who's buying the item",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "seller",
					Description: "Who's selling the item",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "quantity",
					Description: "The number of items to buy",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "item",
					Description: "The name of the item to buy",
					Required:    false,
				},
			},
		},
	},
}

// commandHandler is called (due to the AddHandler above) every time a new
// command is sent on any channel that the authenticated bot has access to.
func (b backpack) commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.ApplicationCommandData().Name != invCommand.Name {
		return
	}

	if len(m.ApplicationCommandData().Options) != 1 {
		say("WTF ARE YOU DOING!?!?!", s, m)
		return
	}
	subcommand := m.ApplicationCommandData().Options[0]

	options := mapOptions(subcommand.Options)
	defaultOwner := fmt.Sprintf("<#%v>", m.ChannelID)

	if subcommand.Name == "view" {
		say(b.displayInvetory(
			getStringOrDefault(options, "owner", defaultOwner),
			false,
		), s, m)
		return
	}

	count, err := getIntOrDefault(options, "quantity", 1)
	if err != nil {
		say("Invalid quantity. Please use a whole number.", s, m)
		return
	}
	if subcommand.Name == "buy" {
		say(b.buyItem(
			count,
			getStringOrDefault(options, "item", ""),
			getStringOrDefault(options, "buyer", defaultOwner),
			getStringOrDefault(options, "seller", defaultOwner),
		), s, m)
		return
	}

	// Handle add, remove, and set.
	price, err := getIntOrDefault(options, "price", UNCHANGED)
	if err != nil {
		say("Invalid price. Please use a whole number.", s, m)
		return
	}
	say(b.modifyItem(
		count,
		price,
		getStringOrDefault(options, "item", COIN),
		getStringOrDefault(options, "owner", defaultOwner),
		subcommand.Name,
	), s, m)
}

// say something in the chat.
func say(msg string, s *discordgo.Session, m *discordgo.InteractionCreate) {
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: msg,
		},
	})
}

// getStringOrDefault will return the option or a default string.
func getStringOrDefault(
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
	key string,
	defaultValue string,
) string {
	if opt, ok := options[key]; ok {
		return opt.StringValue()
	}
	return defaultValue
}

// getIntOrDefault will return the option or a default int.
func getIntOrDefault(
	options map[string]*discordgo.ApplicationCommandInteractionDataOption,
	key string,
	defaultValue int,
) (int, error) {
	if opt, ok := options[key]; ok {
		i, err := strconv.Atoi(opt.StringValue())
		return i, err
	}
	return defaultValue, nil
}

// mapOptions takes a list of options and makes a map of them based on their name.
func mapOptions(
	options []*discordgo.ApplicationCommandInteractionDataOption,
) map[string]*discordgo.ApplicationCommandInteractionDataOption {
	optMap := make(
		map[string]*discordgo.ApplicationCommandInteractionDataOption,
		len(options),
	)
	for _, opt := range options {
		optMap[opt.Name] = opt
	}
	return optMap
}
