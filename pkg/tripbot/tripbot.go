package tripbot

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/dmerrick/danalol-stream/pkg/config"
	"github.com/dmerrick/danalol-stream/pkg/database"
	"github.com/dmerrick/danalol-stream/pkg/events"
	"github.com/dmerrick/danalol-stream/pkg/store"
	mytwitch "github.com/dmerrick/danalol-stream/pkg/twitch"
	"github.com/gempir/go-twitch-irc/v2"
	"github.com/joho/godotenv"
	"github.com/kelvins/geocoder"
)

var lastVid, botUsername, clientAuthenticationToken, twitchClientID, googleMapsAPIKey string
var client *twitch.Client
var datastore *store.Store
var Uptime time.Time

// used to determine which help message to display
// randomized so it starts with a new one every restart
var helpIndex = rand.Intn(len(config.HelpMessages))

const followerMsg = "You must be a follower to run that command :)"

// all chat messages
func PrivateMessage(message twitch.PrivateMessage) {
	user := message.User.Name

	// log the user in if their login time isn't currently recorded
	events.LoginIfNecessary(user)

	if strings.HasPrefix(strings.ToLower(message.Message), "!help") {
		helpCmd(user)
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!uptime") {
		uptimeCmd(user)
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!optimized") {
		optimizedCmd(user)
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!oldmiles") {
		if isFollower(user) {
			oldMilesCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!miles") {
		if isFollower(user) {
			milesCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!sunset") {
		if isFollower(user) {
			sunsetCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!tripbot") {
		if isFollower(user) {
			tripbotCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!leaderboard") {
		if isFollower(user) {
			leaderboardCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!time") {
		if isFollower(user) {
			timeCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!date") {
		if isFollower(user) {
			dateCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}

	if strings.HasPrefix(strings.ToLower(message.Message), "!state") {
		if isFollower(user) {
			stateCmd(user)
		} else {
			client.Say(config.ChannelName, followerMsg)
		}
	}
}

func UserJoin(joinMessage twitch.UserJoinMessage) {
	events.LoginIfNecessary(joinMessage.User)
}

func UserPart(partMessage twitch.UserPartMessage) {
	events.LogoutIfNecessary(partMessage.User)
}

// send message to chat if someone subs
func UserNotice(message twitch.UserNoticeMessage) {
	msg := fmt.Sprintf("%s Your support powers me bleedPurple", message.Message)
	client.Say(config.ChannelName, msg)
}

// if the message comes from me, then post the message to chat
func Whisper(message twitch.WhisperMessage) {
	log.Println("whisper from", message.User.Name, ":", message.Message)
	if message.User.Name == config.ChannelName {
		client.Say(config.ChannelName, message.Message)
	}
}

func Initialize() *twitch.Client {
	var err error
	Uptime = time.Now()

	// load ENV vars from .env file
	err = godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// first we must check for required ENV vars
	if os.Getenv("DASHCAM_DIR") == "" {
		panic("You must set DASHCAM_DIR")
	}
	botUsername = os.Getenv("BOT_USERNAME")
	if botUsername == "" {
		panic("You must set BOT_USERNAME")
	}
	clientAuthenticationToken = os.Getenv("TWITCH_AUTH_TOKEN")
	if clientAuthenticationToken == "" {
		panic("You must set TWITCH_AUTH_TOKEN")
	}
	twitchClientID = os.Getenv("TWITCH_CLIENT_ID")
	if twitchClientID == "" {
		panic("You must set TWITCH_CLIENT_ID")
	}
	googleMapsAPIKey = os.Getenv("GOOGLE_MAPS_API_KEY")
	if googleMapsAPIKey == "" {
		panic("You must set GOOGLE_MAPS_API_KEY")
	}

	// set up geocoder (for translating coords to places)
	geocoder.ApiKey = googleMapsAPIKey

	// initialize the twitch API client
	//TODO: rename me to Initialize()
	_, err = mytwitch.FindOrCreateClient(twitchClientID)
	if err != nil {
		log.Fatal("unable to create twitch API client", err)
	}

	// initialize the SQL database
	database.DBCon, err = database.Initialize()
	if err != nil {
		log.Fatal("error initializing the DB", err)
	}

	// initialize the local datastore
	datastore = store.FindOrCreate(config.DBPath)

	client = twitch.NewClient(botUsername, clientAuthenticationToken)

	return client
}