package duel

import (
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot/worlds"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/game/structure"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

var pairs = [][]world.Block{
	{
		block.Concrete{Colour: item.ColourLime()},
		block.Concrete{Colour: item.ColourGreen()},
	},
	{
		block.Concrete{Colour: item.ColourLightBlue()},
		block.Concrete{Colour: item.ColourBlue()},
	},
	{
		block.Concrete{Colour: item.ColourLightGrey()},
		block.Concrete{Colour: item.ColourBlack()},
	},
}

func Start(p1, p2 *player.Player, g game.Game) {
	UnQueue(p1)
	UnQueue(p2)

	dim := [3]int{50, 20, 80}
	w := world.Config{Entities: ent.Registry, ReadOnly: true}.New()
	s := structure.GenerateBoxStructure(dim, pairs[rand.Intn(len(pairs))]...)
	id := rand.Int63()
	w.BuildStructure(cube.Pos{0, 0, 0}, s)
	w.Handle(&worlds.Handler{})
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeAdventure)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()

	p1.Inventory().Handle(inventory.NopHandler{})
	p1.Armour().Handle(inventory.NopHandler{})
	p1.SetNameTag(text.Colourf("<red>%s</red>", p1.Name()))
	p1.Handle(newHandler(p1, p2, id))
	w.AddEntity(p1)
	p1.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 10})
	kit.Apply(g.Kit(), p1)

	p2.Inventory().Handle(inventory.NopHandler{})
	p2.Armour().Handle(inventory.NopHandler{})
	p2.SetNameTag(text.Colourf("<red>%s</red>", p2.Name()))
	p2.Handle(newHandler(p2, p1, id))
	w.AddEntity(p2)
	p2.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 70})
	kit.Apply(g.Kit(), p2)
}
