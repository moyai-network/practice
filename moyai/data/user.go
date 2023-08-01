package data

import (
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/role"
	"github.com/rcrowley/go-bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"strings"
	"sync"
	"time"
)

var (
	userCollection *mongo.Collection

	usersMu sync.Mutex
	users   = map[string]User{}
)

func init() {
	t := time.NewTicker(time.Minute * 5)
	go func() {
		for range t.C {
			if err := Close(); err != nil {
				log.Println("error saving data:", err)
			}
		}
	}()
}

type User struct {
	XUID        string
	Name        string
	DisplayName string

	Roles *role.Roles

	Punishments struct {
		Ban, Mute moose.Punishment
	}
	Stats struct {
		Kills, Deaths  int
		KillStreak     int
		BestKillStreak int
	}
}

func (u User) WithKills(n int) User {
	stats := u.Stats
	stats.Kills = n
	u.Stats = stats
	return u
}

func (u User) WithDeaths(n int) User {
	stats := u.Stats
	stats.Deaths = n
	u.Stats = stats
	return u
}

func (u User) WithKillStreak(n int) User {
	stats := u.Stats
	stats.KillStreak = n
	u.Stats = stats
	return u
}

func (u User) WithBestKillStreak(n int) User {
	stats := u.Stats
	stats.BestKillStreak = n
	u.Stats = stats
	return u
}

// DefaultUser creates a default user.
func DefaultUser(name string) User {
	return User{
		Name:        strings.ToLower(name),
		DisplayName: name,
		Roles:       role.NewRoles([]moose.Role{role.Default{}}, map[moose.Role]time.Time{}),
	}
}

// LoadOrCreateUser loads a user or creates it, using the given name.
func LoadOrCreateUser(name string) User {
	u, ok := LoadUser(name)
	if !ok {
		return DefaultUser(name)
	}
	return u
}

// LoadUser loads a user using the given name.
func LoadUser(name string) (User, bool) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if u, ok := users[strings.ToLower(name)]; ok {
		return u, true
	}

	filter := bson.M{"name": bson.M{"$eq": strings.ToLower(name)}}

	result := userCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		return User{}, false
	}
	var u User

	err := result.Decode(&u)
	if err != nil {
		return User{}, false
	}
	users[u.Name] = u

	return u, true
}

// SaveUser saves the provided user into the database.
func SaveUser(u User) error {
	usersMu.Lock()
	users[u.Name] = u
	usersMu.Unlock()
	return nil
}

// Close closes and saves the data.
func Close() error {
	usersMu.Lock()
	defer usersMu.Unlock()

	for _, u := range users {
		filter := bson.M{"name": bson.M{"$eq": u.Name}}
		update := bson.M{"$set": u}

		res, err := userCollection.UpdateOne(ctx(), filter, update)
		if err != nil {
			return err
		}

		if res.MatchedCount == 0 {
			_, err = userCollection.InsertOne(ctx(), u)
			return err
		}
	}
	users = map[string]User{}

	return nil
}
