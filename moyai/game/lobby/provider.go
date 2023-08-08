package lobby

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/user"
)

func New(w *world.World) {
	lobby = w
}

var lobby *world.World

func Contains(p *player.Player) bool {
	return lobby == p.World()
}

func AddPlayer(p *player.Player) {
	user.Add(p)
	lobby.AddEntity(p)
	p.Teleport(lobby.Spawn().Vec3Middle())

	if c, closeable := p.Handler().(interface{ Close() }); closeable {
		c.Close()
	}

	kit.Apply(kit.Lobby{}, p)
	h := newHandler(p)
	h.SendScoreBoard()
	p.Handle(h)
	p.Inventory().Handle(inventoryHandler{})
	p.Armour().Handle(inventoryHandler{})

	u, _ := data.LoadUser(p.Name())
	p.SetNameTag(u.Roles.Highest().Colour(p.Name()))
}
