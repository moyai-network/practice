package user

import (
	"fmt"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/oomph-ac/oomph/check"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"

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
	fmt.Println("delete")
	users.Delete(p)
	fmt.Println("deleted")
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

// DuelRequests ...
func (h *Handler) DuelRequests() (requests []string) {
	for xuid, t := range h.duelRequests {
		p, ok := LookupXUID(xuid)
		if !ok || t.Expired() {
			delete(h.duelRequests, xuid)
			continue
		}
		requests = append(requests, p.Name())
	}
	return
}

// Duel ...
func (h *Handler) Duel(p *player.Player, g game.Game) {
	h.duelRequests[p.XUID()] = newRequest(g)
}

func (h *Handler) AcceptDuel(t *player.Player) game.Game {
	defer delete(h.duelRequests, t.XUID())
	return h.duelRequests[t.XUID()].Game()
}

func (h *Handler) History() map[check.Check]float64 {
	h.historyMu.Lock()
	defer h.historyMu.Unlock()
	return h.history
}

func (h *Handler) SetRecentReplay(re []struct {
	Name   string
	Packet packet.Packet
}) {
	h.recentReplay = re
}

func (h *Handler) RecentReplay() ([]struct {
	Name   string
	Packet packet.Packet
}, bool) {
	return h.recentReplay, len(h.recentReplay) != 0
}

func (h *Handler) SetWatchingReplay(b bool) {
	h.watchingReplay.Store(b)
}

func (h *Handler) WatchingReplay() bool {
	return h.watchingReplay.Load()
}
