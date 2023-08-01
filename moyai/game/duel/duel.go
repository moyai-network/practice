package duel

import (
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"sync"
)

var (
	duelsMu sync.Mutex
	duels   = map[*player.Player]*Match{}
)

func Lookup(p *player.Player) (*Match, bool) {
	m, ok := duels[p]
	return m, ok
}

func AddPlayer(p *player.Player, m *Match) {
	kit.Apply(m.g.Kit(), p)

	p.Inventory().Handle(inventory.NopHandler{})
	p.Armour().Handle(inventory.NopHandler{})

	duelsMu.Lock()
	duels[p] = m
	duelsMu.Unlock()

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))
}

func RemovePlayer(p *player.Player) {
	duelsMu.Lock()
	delete(duels, p)
	duelsMu.Unlock()
}
