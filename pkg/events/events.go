package events

import (
	"time"

	"github.com/dmerrick/danalol-stream/pkg/database"
)

type Event struct {
	Username    string    `db:"username"`
	Event       string    `db:"event"`
	DateCreated time.Time `db:"date_created"`
}

func Login(user string) {
	tx := database.DBCon.MustBegin()
	tx.NamedExec(
		"INSERT INTO events (username, event) VALUES (:username, :event)",
		&Event{user, "login"},
	)
	tx.Commit()
}

func Logout(user string) {
	tx := database.DBCon.MustBegin()
	tx.NamedExec(
		"INSERT INTO events (username, event) VALUES (:username, :event)",
		&Event{user, "logout"},
	)
	tx.Commit()
}
