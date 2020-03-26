package users

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/dmerrick/tripbot/pkg/config"
	"github.com/dmerrick/tripbot/pkg/database"
	terrors "github.com/dmerrick/tripbot/pkg/errors"
	"github.com/dmerrick/tripbot/pkg/helpers"
	"github.com/logrusorgru/aurora"
)

var Leaderboard [][]string
var initLeaderboardSize = 25
var maxLeaderboardSize = 50

// Leaderboard creates a leaderboard
func InitLeaderboard() {
	users := []User{}
	query := `SELECT * FROM users WHERE miles != 0 AND is_bot = false AND username!=$1 ORDER BY miles DESC LIMIT $2`
	database.Connection().Select(&users, query, strings.ToLower(config.ChannelName), initLeaderboardSize)
	for _, user := range users {
		miles := fmt.Sprintf("%.1f", user.Miles)
		pair := []string{user.Username, miles}
		Leaderboard = append(Leaderboard, pair)
	}
}

func UpdateLeaderboard() {
	for _, user := range LoggedIn {
		// skip adding this user if they're a bot (or me)
		if user.IsBot || helpers.UserIsAdmin(user.Username) {
			continue
		}
		insertIntoLeaderboard(*user)
	}
	// truncate Leaderboard if it gets too big
	if len(Leaderboard) > maxLeaderboardSize {
		Leaderboard = Leaderboard[:maxLeaderboardSize]
	}
}

// convert the string to a float32
func strToFloat32(str string) float32 {
	value, err := strconv.ParseFloat(str, 32)
	if err != nil {
		terrors.Log(err, "error parsing float")
		return 0.0
	}
	return float32(value)
}

func insertIntoLeaderboard(user User) {
	// first we remove this user from the board
	removeFromLeaderboard(user.Username)

	// get the current miles as a float
	miles := user.CurrentMiles()

	for i, pair := range Leaderboard {
		val := strToFloat32(pair[1])
		// see if our miles are higher
		if miles >= val {
			milesStr := fmt.Sprintf("%.1f", miles)
			newPair := []string{user.Username, milesStr}

			// insert into Leaderboard
			// https://github.com/golang/go/wiki/SliceTricks#insert
			Leaderboard = append(Leaderboard[:i], append([][]string{newPair}, Leaderboard[i:]...)...)
			return
		}
	}
}

// removeFromLeaderboard searches the Leaderboard for
// a username and removes it
func removeFromLeaderboard(username string) {
	for i, pair := range Leaderboard {
		if pair[0] == username {
			// delete from Leaderboard
			// https://github.com/golang/go/wiki/SliceTricks#delete
			Leaderboard = append(Leaderboard[:i], Leaderboard[i+1:]...)
			return
		}
	}
}

// this was used for development
func printLeaderboard() {
	for i, pair := range Leaderboard {
		fmt.Printf("%d: %s - %s\n", i+1, pair[1], aurora.Magenta(pair[0]))
	}
}
