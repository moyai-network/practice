package data

import (
	"github.com/moyai-network/practice/moyai/game"
	"log"
	"strings"
	"sync"
	"time"

	"golang.org/x/exp/maps"

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

	DeviceID     string
	SelfSignedID string
	Address      string

	FirstLogin time.Time
	PlayTime   time.Duration

	Roles *role.Roles

	Punishments struct {
		Ban, Mute carrot.Punishment
	}

	Stats Stats

	Settings Settings
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

func (u User) WithIncreasedWin(competitive bool) User {
	stats := u.Stats
	if competitive {
		stats.CompetitiveWins++
	} else {
		stats.CasualWins++
	}
	u.Stats = stats
	return u
}

func (u User) WithIncreasedLoss(competitive bool) User {
	stats := u.Stats
	if competitive {
		stats.CompetitiveLosses++
	} else {
		stats.CasualLosses++
	}
	u.Stats = stats
	return u
}

func (u User) KDR() float64 {
	var kdr float64
	if u.Stats.Deaths > 0 {
		kdr = float64(u.Stats.Kills) / float64(u.Stats.Deaths)
	} else {
		kdr = float64(u.Stats.Kills)
	}
	return kdr
}

func (u User) WithElo(g game.Game, n int32) User {
	stats := u.Stats
	if stats.Elo == nil {
		stats.Elo = map[string]int32{}
	}
	if _, ok := stats.Elo[strings.ToLower(g.Name())]; !ok {
		stats.Elo[strings.ToLower(g.Name())] = 1000
	}
	stats.Elo[strings.ToLower(g.Name())] = n

	u.Stats = stats
	return u
}

func (u User) GameElo(g game.Game) int32 {
	stats := u.Stats
	if stats.Elo == nil {
		stats.Elo = map[string]int32{}
	}
	if _, ok := stats.Elo[strings.ToLower(g.Name())]; !ok {
		stats.Elo[strings.ToLower(g.Name())] = 1000
	}
	return stats.Elo[strings.ToLower(g.Name())]
}

func (u User) TotalElo() int32 {
	var tot int32 = 1000
	for _, g := range game.Games() {
		if !g.Duel() {
			continue
		}
		tot += u.GameElo(g) - 1000
	}
	return tot
}

func (u User) WithIncreasedPlayTime(inc time.Duration) User {
	u.PlayTime += inc
	return u
}

func (u User) WithSettings(s Settings) User {
	u.Settings = s
	return u
}

// DefaultUser creates a default user.
func DefaultUser(name string) User {
	return User{
		Name:        strings.ToLower(name),
		DisplayName: name,
		Roles:       role.NewRoles([]carrot.Role{role.Default{}}, map[carrot.Role]time.Time{}),
		Stats:       DefaultStats(),
		Settings:    DefaultSettings(),
		FirstLogin:  time.Now(),
	}
}

type Settings struct {
	Display struct {
		Scoreboard bool
		CPS        bool
	}
	Privacy struct {
		PrivateMessages bool
	}
}

func DefaultSettings() Settings {
	s := Settings{}
	s.Display.Scoreboard = true
	s.Display.CPS = true
	s.Privacy.PrivateMessages = true
	return s
}

type Stats struct {
	Kills, Deaths  int
	KillStreak     int
	BestKillStreak int

	Elo map[string]int32

	CasualWins        int
	CasualLosses      int
	CompetitiveWins   int
	CompetitiveLosses int
}

func DefaultStats() Stats {
	return Stats{
		Elo: map[string]int32{},
	}
}

// Users returns the user data for all users.
func Users() []User {
	usersMu.Lock()
	m := users
	usersMu.Unlock()

	cur, err := userCollection.Find(ctx(), bson.M{}, nil)
	if err != nil {
		return maps.Values(m)
	}
	var usrs []User

	_ = cur.All(ctx(), &usrs)

	for _, u := range usrs {
		if _, ok := m[u.Name]; !ok {
			m[u.Name] = u
		}
	}
	return maps.Values(m)
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
		if u.Roles == nil {
			u.Roles = role.NewRoles([]carrot.Role{role.Default{}}, map[carrot.Role]time.Time{})
		}
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
