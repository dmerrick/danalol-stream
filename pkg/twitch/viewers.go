package twitch

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/dmerrick/danalol-stream/pkg/config"
)

// chattersAPIURL is the URL to hit for current chatter list
var chattersAPIURL = "https://tmi.twitch.tv/group/user/" + strings.ToLower(config.ChannelName) + "/chatters"

// chattersResponse is the json returned by the Twitch chatters endpoint
type chattersResponse struct {
	Count    int                 `json:"chatter_count"`
	Chatters map[string][]string `json:"chatters"`
}

// currentChatters will contain the current viewers
var currentChatters chattersResponse

// ChatterCount returns the number of chatters (as reported by Twitch)
func ChatterCount() int {
	return currentChatters.Count
}

// Chatters returns a map where the keys are current chatters
// we use an empty struct for performance reasons
// c.p. https://stackoverflow.com/a/10486196
//TODO: consider using an int as the value and have that be the ID in the DB
func Chatters() map[string]struct{} {
	//TODO: maybe we don't want to make this every time?
	var chatters = make(map[string]struct{})
	for _, list := range currentChatters.Chatters {
		for _, chatter := range list {
			chatters[chatter] = struct{}{}
		}
	}
	return chatters
}

// UpdateChatters makes a request to the chatters API and updates currentChatters
func UpdateChatters() {
	var latestChatters chattersResponse

	client := http.Client{
		Timeout: 2 * time.Second,
	}

	req, err := http.NewRequest(http.MethodGet, chattersAPIURL, nil)
	if err != nil {
		log.Println("error creating request", err)
		return
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")

	res, err := client.Do(req)
	if err != nil {
		log.Println("error making request", err)
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading request body", err)
		return
	}

	err = json.Unmarshal(body, &latestChatters)
	if err != nil {
		log.Println("error unmarshalling json", err)
	}

	currentChatters = latestChatters
}

// PrintCurrentChatters prints the current chatters
//TODO: this was added for debugging purposes and can probably be removed
func PrintCurrentChatters() {
	usernames := make([]string, 0, len(Chatters()))
	for username := range Chatters() {
		usernames = append(usernames, username)
	}
	sort.Sort(sort.StringSlice(usernames))
	log.Printf("Current chatters: %s", strings.Join(usernames, ", "))
}
