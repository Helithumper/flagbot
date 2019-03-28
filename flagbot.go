/*
 * flagbot.go
 * Author: github.com/helithumper
 * Flagbot is a bot for discord meant to remove SunshineCTF flags as they appear
 * Spring 2019
 */

package main

import (
	"bufio"
	"flag"
	"math/rand"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

var (
	token             string
	gifSlice          []string
	regexSlice        []regexp.Regexp
	responseSlice     []string
	configurationPath string
)

func init() {
	// Setup logging
	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
		DisableColors: true,
	})
	log.SetLevel(log.InfoLevel)

	// Parse out command line arguments
	flag.StringVar(&token, "t", "", "Bot Token")
	flag.StringVar(&configurationPath, "c", "configuration", "Configuration directory path")
	flag.Parse()

	// Check if discord token was properly parsed. If not exit
	if token == "" {
		log.Error("Error in parsing command. Proper Usage: flagbot -t <bot token> -c <configuration path>")
		os.Exit(1)
	}

	// Load files into their respective variables
	var err error
	gifSlice, err = readFileToSlice(configurationPath + "/gifs.txt")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading GIFs to Slice")
		os.Exit(1)
	}

	responseSlice, err = readFileToSlice(configurationPath + "/responses.txt")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading Responses to Slice")
		os.Exit(1)
	}

	var regexTextSlice []string
	regexTextSlice, err = readFileToSlice(configurationPath + "/patterns.txt")
	if err != nil {
		log.WithFields(log.Fields{
			"Error": err,
		}).Error("Error reading Regex Patterns to Slice")
		os.Exit(1)
	}
	for _, pattern := range regexTextSlice {
		regex := regexp.MustCompile(pattern)
		regexSlice = append(regexSlice, *regex)
	}
}

func main() {
	// Create a new discord bot using the token from the config file
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Error("Error creating Discord session: ", err)
		return
	}
	log.Info("Created Discord Bot")

	// Register ready handler
	dg.AddHandler(ready)

	// Register message creation handler
	dg.AddHandler(messageCreate)

	// Open websocket and listen
	err = dg.Open()
	if err != nil {
		log.Error("Error opening Discord session: ", err)
		return
	}
	log.Info("Created Websocket")

	// Wait here until CTRL-C or other term signal is received.
	log.Info("Flagbot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

// ready runs as soon as the bot is ready after deployment.
// TODO: swap out hard coded status message for something more dynamic.
//		 perhaps use an outside file to hold the string along with other
//		 general metadata for the application?
func ready(s *discordgo.Session, event *discordgo.Ready) {
	s.UpdateStatus(0, "Watching for flags (´･ω･`)")
}

// messageCreate is a hook for any new discord messages. Any time a message is created
// in a monitored channel, this function is run against it
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	go flagCheck(s, m)
}

// flagCheck checks a given message `m` during a session `s` to see if it matches
// the given flag patterns. If so, a fun message and GIF are printed to the screen
func flagCheck(s *discordgo.Session, m *discordgo.MessageCreate) {
	if isFlag(m.Content) {
		log.WithFields(log.Fields{
			"message": m.Content,
			"author":  m.Author.Username,
		}).Info("Removed Message")

		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.WithFields(log.Fields{
				"messageID": m.ID,
				"channelID": m.ChannelID,
				"Error":     err,
			}).Error("Could not delete message.")
		}

		randMessage := randItem(responseSlice)
		_, err = s.ChannelMessageSend(m.ChannelID, m.Author.Mention()+" "+randMessage)
		if err != nil {
			log.WithFields(log.Fields{
				"messageID": m.ID,
				"channelID": m.ChannelID,
				"Error":     err,
			}).Error("Could not send post-delete text.")
		}

		randGif := randItem(gifSlice)
		_, err = s.ChannelMessageSend(m.ChannelID, randGif)
		if err != nil {
			log.WithFields(log.Fields{
				"messageID": m.ID,
				"channelID": m.ChannelID,
				"Error":     err,
			}).Error("Could not send post-delete gif.")
		}
	}
}

// isFlag returns true if the given message matches a flag pattern.
// It does this by iterating through a global set of Regex patterns
// initialized earlier on in the program.
func isFlag(message string) bool {
	for _, pattern := range regexSlice {
		if pattern.MatchString(message) {
			return true
		}
	}
	return false
}

// readFileToSlice reads in the file at `filepath` and returns a slice
// containing each line of the file.
func readFileToSlice(filepath string) ([]string, error) {
	file, err := os.Open(filepath)

	if err != nil {
		return []string{}, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, scanner.Err()
}

// randLine takes in a `filepath` and returns a single line from said file
func randItem(Slice []string) string {
	rand.Seed(time.Now().Unix())
	return Slice[rand.Intn(len(Slice))]
}
