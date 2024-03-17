package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string // bad
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.Parse()
}

func main() {

	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatalln(err)
	}

	dg.AddHandler(messageCreate)

	dg.Identify.Intents = discordgo.IntentsGuildMessages

	if err := dg.Open(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("bot running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

type Gopher struct {
	Name string `json: "name"`
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.Identify.Intents = discordgo.IntentMessageContent
	fmt.Println(m.Content)
	if m.Content == "!gopher" {
		resp, err := http.Get("http://localhost:8080/gopher/dr-who")
		if err != nil {
			fmt.Print(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "dr-who.png", resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("Error: Can't get dr-who Gopher! :-(")
		}
	}

	if m.Content == "!random" {
		resp, err := http.Get("http://localhost:8080/gopher/random")
		if err != nil {
			fmt.Print(err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode == 200 {
			_, err = s.ChannelFileSend(m.ChannelID, "random-gopher.png", resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			fmt.Println("Error: Can't get random Gopher! :-(")
		}
	}
}
