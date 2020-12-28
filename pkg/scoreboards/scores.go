package users

import (
	"log"
	"time"

	"github.com/adanalife/tripbot/pkg/database"
)

// CREATE TABLE scores (
//   id            SERIAL PRIMARY KEY,
//   user_id       INTEGER NOT NULL,
//   scoreboard_id INTEGER NOT NULL,
//   score         REAL DEFAULT 0.0, /* float32: https://github.com/go-pg/pg/wiki/Model-Definition */
//   date_created  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
// );

type Score struct {
	ID           uint16    `db:"id"`
	UserID       uint16    `db:"user_id"`
	ScoreboardID uint16    `db:"scoreboard_id"`
	Score        float32   `db:"score"`
	DateCreated  time.Time `db:"date_created"`
}

// User.save() will take the given user and store it in the DB
func (s Score) save() {
	if c.Conf.Verbose {
		log.Println("saving score", u)
	}
	query := `UPDATE scores SET score=:score, WHERE id = :id`
	_, err := database.Connection().NamedExec(query, u)
	if err != nil {
		terrors.Log(err, "error saving score")
	}
}

//// IsFollower returns true if the user is a follower
//func (u User) IsFollower() bool {
//	return twitch.UserIsFollower(u.Username)
//}

//// IsSubscriber returns true if the user is a subscriber
//func (u User) IsSubscriber() bool {
//	return twitch.UserIsSubscriber(u.Username)
//}

//// User.String prints a colored version of the user
//func (u User) String() string {
//	if u.IsBot {
//		return aurora.Gray(15, u.Username).String()
//	}
//	if c.UserIsAdmin(u.Username) {
//		return aurora.Gray(11, u.Username).String()
//	}
//	return aurora.Magenta(u.Username).String()
//}

//// FindOrCreate will try to find the user in the DB, otherwise it will create a new user
//func FindOrCreate(username string) User {
//	if c.Conf.Verbose {
//		log.Printf("FindOrCreate(%s)", username)
//	}
//	user := Find(username)
//	if user.ID != 0 {
//		return user
//	}
//	// create the user in the DB
//	return create(username)
//}

//// Find will look up the username in the DB, and return a User if possible
//func Find(username string) User {
//	var user User
//	query := `SELECT * FROM users WHERE username=$1`
//	err := database.Connection().Get(&user, query, username)
//	// spew.Config.ContinueOnMethod = true
//	// spew.Config.MaxDepth = 2
//	// spew.Dump(user)
//	if err != nil {
//		//TODO: is there a better way to do this?
//		return User{ID: 0}
//	}
//	return user
//}

//// HasCommandAvailable lets users run a command once a day,
//// unless they are a follower in which case they can run
//// as many as they like
//func (u *User) HasCommandAvailable() bool {
//	// followers get unlimited commands
//	if u.IsFollower() {
//		return true
//	}
//	// check if they ran a command in the last 24 hrs
//	now := time.Now()
//	if now.Sub(u.lastCmd) > 24*time.Hour {
//		log.Println("letting", u, "run a command")
//		// update their lastCmd time
//		u.lastCmd = now
//		return true
//	}
//	return false
//}

////TODO: maybe return an err here?
//// create() will actually create the DB record
//func create(username string) User {
//	log.Println("creating user", username)
//	tx := database.Connection().MustBegin()
//	// create a new row, using default vals and creating a single visit
//	tx.MustExec("INSERT INTO users (username, num_visits) VALUES ($1, $2)", username, 1)
//	tx.Commit()
//	return Find(username)
//}