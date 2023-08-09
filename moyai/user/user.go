package user

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/carrot/sets"
	"github.com/moyai-network/practice/moyai/data"
)

var (
	users = sets.Set[*player.Player]{}
)

// All returns a slice of all the users.
func All() []*player.Player {
	return users.Values()
}

func Count() int {
	return len(users)
}

func Add(p *player.Player) {
	users.Add(p)
}

func Remove(p *player.Player) {
	users.Delete(p)
}

func Lookup(name string) (*player.Player, bool) {
	for p := range users {
		if p.Name() == name {
			return p, true
		}
	}
	return nil, false
}

func LookupXUID(xuid string) (*player.Player, bool) {
	for p := range users {
		if p.XUID() == xuid {
			return p, true
		}
	}
	return nil, false
}

// Alert alerts all staff users with an action performed by a cmd.Source.
func Alert(s cmd.Source, key string, args ...any) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	for _, h := range All() {
		if u, _ := data.LoadUserOrCreate(h.Name()); role.Staff(u.Roles.Highest()) {
			h.Message(lang.Translatef(h.Locale(), "staff.alert", p.Name(), fmt.Sprintf(lang.Translate(h.Locale(), key), args...)))
		}
	}
}

func Broadcast(key string, args ...any) {
	for _, p := range users.Values() {
		p.Message(lang.Translatef(p.Locale(), key, args...))
	}
}
