// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

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

func main() {
	// Load bot token.
	token := os.Getenv("BACKPACK_TOKEN")
	if token == "" {
		log.Fatalf("token is missing, you must set BACKPACK_TOKEN")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v\n", err)
	}

	// Register the commandHandler func for InteractionCreate events.
	dg.AddHandler(commandHandler)

	// Recieve messages.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v\n", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("backpack bot running")

	cmd, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", &invCommand)
	if err != nil {
		log.Fatalf("cannot create '%v' command: %v\n", invCommand.Name, err)
	}
	log.Println("registered commands")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-stop

	err = dg.ApplicationCommandDelete(dg.State.User.ID, "", cmd.ID)
	if err != nil {
		log.Fatalf("cannot delete '%v' command: %v\n", cmd.Name, err)
	}

	// Cleanly close down the Discord session.
	dg.Close()
}

// commandHandler is called (due to the AddHandler above) every time a new
// command is sent on any channel that the authenticated bot has access to.
func commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if m.ApplicationCommandData().Name != invCommand.Name {
		return
	}

	// DEBUG
	log.Println("CHANNEL ID:", m.ChannelID)

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
	shouldPrint := true
	var response string
	if opt, ok := optMap["add"]; ok {
		shouldPrint = false
		response += addItem(owner, opt.StringValue())
	}

	if opt, ok := optMap["remove"]; ok {
		shouldPrint = false
		response += removeItem(owner, opt.StringValue())
	}

	if shouldPrint {
		response += printInventory(owner)
	}

	// Send our response.
	s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}

func addItem(owner, item string) string {
	s := fmt.Sprintf("added %v to %v\n", item, owner)
	log.Println(s)
	return s
}

func removeItem(owner, item string) string {
	s := fmt.Sprintf("removed %v from %v\n", item, owner)
	log.Println(s)
	return s
}

func printInventory(owner string) string {
	s := fmt.Sprintf("printing %v's inventory\n", owner)
	log.Println(s)
	return s
}
