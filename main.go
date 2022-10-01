// License: AGPL-3.0-only
// (c) 2022 Dakota Walsh <kota@nilsu.org>
package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	// Load bot token.
	token := os.Getenv("BACKPACK_TOKEN")
	if token == "" {
		log.Fatalf("token is missing, you must set BACKPACK_TOKEN")
	}

	dir := os.Getenv("BACKPACK_DATA")
	if dir == "" {
		log.Fatalf("you must set BACKPACK_DATA")
	}
	info, err := os.Stat(dir)
	if errors.Is(err, os.ErrNotExist) {
		// Create the data directory.
		if err := os.MkdirAll(dir, 0777); err != nil {
			log.Fatalf("failed creating data directory: %v: %v\n", dir, err)
		}
	} else if !info.IsDir() {
		log.Fatalln("data path is a file instead of a directory")
	} else if err != nil {
		log.Fatalf("error reading data directory: %v: %v\n", dir, err)
	}

	b := backpack{
		dir: dir,
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("error creating Discord session: %v\n", err)
	}

	// Register the commandHandler func for InteractionCreate events.
	dg.AddHandler(b.commandHandler)

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
