package data

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/role"
	"github.com/rcrowley/go-bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userCollection *mongo.Collection

	usersMu sync.Mutex
	users   = map[string]User{}
)

func init() {
	t := time.NewTicker(time.Second * 5)
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
		Ban, Mute carrot.Punishment
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
		Roles:       role.NewRoles([]carrot.Role{role.Default{}}, map[carrot.Role]time.Time{}),
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

// LoadUserOrCreate loads a user using the given name.
func LoadUserOrCreate(name string) (User, error) {
	usersMu.Lock()
	defer usersMu.Unlock()

	if u, ok := users[strings.ToLower(name)]; ok {
		return u, nil
	}
	filter := bson.M{"name": bson.M{"$eq": strings.ToLower(name)}}

	result := userCollection.FindOne(ctx(), filter)
	if err := result.Err(); err != nil {
		if err == mongo.ErrNoDocuments {
			return DefaultUser(name), nil
		}
		return User{}, err
	}
	var u User

	err := result.Decode(&u)
	if err != nil {
		return User{}, err
	}

	users[u.Name] = u

	return u, nil
}

// LoadUsersCond loads users using the given filter.
func LoadUsersCond(cond any) ([]User, error) {
	collection := db.Collection("users")
	count, err := collection.EstimatedDocumentCount(ctx())
	if err != nil {
		return nil, err
	}

	var data = make([]User, count)

	for d := 0; d < int(count); d++ {
		data[d] = User{}
	}

	cursor, err := collection.Find(ctx(), cond)
	if err != nil {
		return nil, err
	}

	if err = cursor.All(ctx(), &data); err != nil {
		return nil, err
	}

	usersMu.Lock()
	for i, u := range data {
		if d, ok := users[u.Name]; ok {
			data[i] = d
		}
	}
	usersMu.Unlock()

	return data, nil
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
