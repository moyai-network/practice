package user

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/moose/sets"
)

var (
	users = sets.Set[*player.Player]{}
)

func Count() int {
	return len(users)
}

func Add(p *player.Player) {
	users.Add(p)
}

func Remove(p *player.Player) {
	users.Delete(p)
}

func LookupXUID(xuid string) (*player.Player, bool) {
	for p := range users {
		if p.XUID() == xuid {
			return p, true
		}
	}
	return nil, false
}

func Broadcast(key string, args ...any) {
	for _, p := range users.Values() {
		p.Message(lang.Translatef(p.Locale(), key, args...))
	}
}
