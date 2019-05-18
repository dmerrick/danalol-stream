package store

import (
	"log"
	"os"
	"os/signal"
	"sort"
	"syscall"
	"time"

	"github.com/boltdb/bolt"
	"github.com/dmerrick/danalol-stream/pkg/helpers"
)

const (
	userJoinsBucket   = "user_joins"
	userWatchedBucket = "user_watched"
)

// this stores the current datastore
var currentDatastore *Store

func FindOrCreate() *Store {
	// use the pre-existing datastore if we have one
	if currentDatastore != nil {
		return currentDatastore
	}

	// initialize the database
	datastore := NewStore(helpers.DbPath)
	if err := datastore.Open(); err != nil {
		panic(err)
	}

	// catch CTRL-Cs and run datastore.Close()
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		datastore.Close()
		os.Exit(1)
	}()

	// save the current datastore for use later
	currentDatastore = datastore

	return datastore
}

type Store struct {
	path string
	db   *bolt.DB
}

func NewStore(path string) *Store {
	return &Store{
		path: path,
	}
}

// Open opens and initializes the database.
func (s *Store) Open() error {
	// Open underlying data store.
	db, err := bolt.Open(s.path, 0666, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return err
	}
	s.db = db

	// Initialize all the required buckets.
	if err := s.db.Update(func(tx *bolt.Tx) error {
		tx.CreateBucketIfNotExists([]byte(userJoinsBucket))
		tx.CreateBucketIfNotExists([]byte(userWatchedBucket))
		return nil
	}); err != nil {
		s.Close()
		return err
	}

	return nil
}

func (s *Store) Close() error {
	if s.db != nil {
		var onlineUsers = []string{}

		// first we make a list of all of the online users
		s.db.View(func(tx *bolt.Tx) error {
			joinedBucket := tx.Bucket([]byte(userJoinsBucket))
			err := joinedBucket.ForEach(func(k, _ []byte) error {
				user := string(k)
				onlineUsers = append(onlineUsers, user)
				return nil
			})
			return err
		})

		// then we loop over it and record the current watched duration
		for _, user := range onlineUsers {
			log.Println("logging out", user)
			s.RecordUserPart(string(user))
		}

		s.db.Close()
	}
	return nil
}

func (s *Store) MilesForUser(user string) int {
	prevDuration := s.DurationForUser(user)
	currDuration := s.CurrentViewDuration(user)
	return helpers.DurationToMiles(prevDuration + currDuration)
}

// fetch the current view duration
func (s *Store) CurrentViewDuration(user string) time.Duration {
	var joinTime time.Time
	err := s.db.View(func(tx *bolt.Tx) error {
		joinedBucket := tx.Bucket([]byte(userJoinsBucket))
		err := joinTime.UnmarshalText(joinedBucket.Get([]byte(user)))
		return err
	})
	if err != nil {
		log.Printf("encountered error getting current view duration: %s", err)
		return 0
	}
	return time.Since(joinTime)
}

func (s *Store) TopUsers(size int) []string {
	var reversedMap = make(map[int]string)
	var topUsers = []string{}

	sortedValues := []int{}

	s.db.View(func(tx *bolt.Tx) error {
		watchedBucket := tx.Bucket([]byte(userWatchedBucket))
		//TODO: do I need this?
		// joinedBucket := tx.Bucket([]byte(userJoinsBucket))

		err := watchedBucket.ForEach(func(k, v []byte) error {
			user := string(k)
			// don't include these in leaderboard
			//TODO: this should be a func
			for _, ignored := range helpers.IgnoredUsers {
				if user == ignored {
					return nil
				}
			}

			// fetch current view duration...
			// currentDuration := s.CurrentViewDuration(user)
			// ...and the previous view duration
			duration, err := time.ParseDuration(string(v))
			if err != nil {
				return err
			}
			// add them together
			intDuration := int(duration) // + int(currentDuration)
			sortedValues = append(sortedValues, intDuration)
			reversedMap[intDuration] = user
			return err
		})

		return err
	})

	// sort and reverse
	sort.Sort(sort.Reverse(sort.IntSlice(sortedValues)))

	// print the top 5
	for i := 0; i < size; i++ {
		topUsers = append(topUsers, reversedMap[sortedValues[i]])
	}

	return topUsers
}

func (s *Store) DurationForUser(user string) time.Duration {
	var previousDurationWatched *time.Duration

	s.db.View(func(tx *bolt.Tx) error {
		watchedBucket := tx.Bucket([]byte(userWatchedBucket))

		// fetch the previous duration watched from the DB
		duration, err := time.ParseDuration(string(watchedBucket.Get([]byte(user))))
		previousDurationWatched = &duration
		return err
	})
	return *previousDurationWatched
}

func (s *Store) PrintStats() {
	s.db.View(func(tx *bolt.Tx) error {
		watchedBucket := tx.Bucket([]byte(userWatchedBucket))
		err := watchedBucket.ForEach(func(k, v []byte) error {
			duration, err := time.ParseDuration(string(v))
			log.Printf("%s has watched %s.\n", k, duration)
			return err
		})
		return err
	})
}

func (s *Store) RecordUserJoin(user string) {
	log.Println(user, "joined the channel")

	s.db.Update(func(tx *bolt.Tx) error {
		joinedBucket := tx.Bucket([]byte(userJoinsBucket))
		currentTime, err := time.Now().MarshalText()
		if err != nil {
			return err
		}
		err = joinedBucket.Put([]byte(user), []byte(currentTime))
		return err
	})
}

func (s *Store) RecordUserPart(user string) {
	var joinTime time.Time
	var durationWatched time.Duration
	var previousDurationWatched time.Duration

	s.db.View(func(tx *bolt.Tx) error {
		joinedBucket := tx.Bucket([]byte(userJoinsBucket))
		watchedBucket := tx.Bucket([]byte(userWatchedBucket))

		// first find the time the user joined the channel
		err := joinTime.UnmarshalText(joinedBucket.Get([]byte(user)))
		if err != nil {
			return err
		}

		// seems like we did find a time, so calculate the duration watched
		durationWatched = time.Since(joinTime)

		// fetch the previous duration watched from the DB
		previousDurationWatched, err = time.ParseDuration(string(watchedBucket.Get([]byte(user))))
		return err
	})

	// calculate total duration watched
	totalDurationWatched := previousDurationWatched + durationWatched

	// update the DB with the total duration watched
	s.db.Update(func(tx *bolt.Tx) error {
		joinedBucket := tx.Bucket([]byte(userJoinsBucket))
		watchedBucket := tx.Bucket([]byte(userWatchedBucket))

		err := watchedBucket.Put([]byte(user), []byte(totalDurationWatched.String()))
		if err != nil {
			return err
		}
		// remove the user from the joined bucket
		err = joinedBucket.Delete([]byte(user))
		return err
	})

	log.Println(user, "left the channel, total watched:", totalDurationWatched)
}
