package ffa

import (
	"log"
	"time"

	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

func AddPlayer(p *player.Player, g game.Game, lobby func(player2 *player.Player)) {
	if c, closeable := p.Handler().(interface{ Close() }); closeable {
		c.Close()
	}

	w, ok := ffas[g]
	if !ok {
		log.Fatalln("no world found for game:", g.Name())
	}

	w.AddEntity(p)
	p.Teleport(w.Spawn().Vec3Middle())
	p.SetImmobile()
	kit.Apply(g.Kit(), p)

	h := newHandler(p, g, lobby)
	h.SendScoreBoard()
	p.Handle(h)
	p.Inventory().Handle(inventory.NopHandler{})
	p.Armour().Handle(inventory.NopHandler{})

	p.SetNameTag(text.Colourf("<red>%s</red>", p.Name()))
	h.SendScoreBoard()
	// TODO: remove this later
	time.AfterFunc(time.Millisecond*500, func() {
		p.SetMobile()
	})
}
