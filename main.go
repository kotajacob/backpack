// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{
		Name: "basic-command",
		// All commands and options must have a description
		// Commands/options without description will fail the registration
		// of the command.
		Description: "Basic command",
	},
}

const ROOT_CMD = "/inv"

func main() {
	// Load bot token.
	token := os.Getenv("BACKPACK_TOKEN")
	if token == "" {
		log.Fatalf("token is missing, you must set BACKPACK_TOKEN")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Recieve messages.
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("backpack bot running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// If the message is "ping" reply with "pong"
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "pong")
	}

	// If the message is "pong" reply with "ping"
	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "ping")
	}
}
