package lobby

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
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

	kit.Apply(kit.Lobby{}, p)
	p.Handle(newHandler(p))
	p.Inventory().Handle(inventoryHandler{})
	p.Armour().Handle(inventoryHandler{})

	u, _ := data.LoadUser(p.Name())
	p.SetNameTag(u.Roles.Highest().Colour(p.Name()))

	updateScoreBoards()
}

func updateScoreBoards() {
	for _, e := range lobby.Entities() {
		p, ok := e.(*player.Player)
		if !ok {
			continue
		}
		l := p.Locale()
		u, _ := data.LoadUser(p.Name())

		sb := scoreboard.New(moose.GlyphFont("PRACTICE", item.ColourOrange()))
		sb.RemovePadding()
		_, _ = sb.WriteString("§r\uE000")
		_, _ = sb.WriteString(text.Colourf("<black>Role</black><grey>:</grey> %s", u.Roles.Highest().Colour(u.Roles.Highest().Name())))
		_, _ = sb.WriteString("\uE000")
		_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
		for i, li := range sb.Lines() {
			if !strings.Contains(li, "\uE000") {
				sb.Set(i, " "+li)
			}
		}
		p.SendScoreboard(sb)
	}
}
