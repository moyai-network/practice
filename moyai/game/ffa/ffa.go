package ffa

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot/worlds"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/structure"
	"math/rand"
)

var ffas = map[game.Game]*world.World{}

func init() {
	pairs := [][]world.Block{
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

	for _, g := range game.Games() {
		size := 256
		w := world.Config{Entities: ent.Registry}.New()
		s := structure.GenerateBoxStructure([3]int{size, 20, size}, pairs[rand.Intn(len(pairs))]...)
		w.BuildStructure(cube.Pos{0, 0, 0}, s)
		w.SetSpawn(cube.Pos{size / 2, 2, size / 2})
		w.Handle(&worlds.Handler{})
		w.StopWeatherCycle()
		w.SetDefaultGameMode(world.GameModeAdventure)
		w.SetTime(6000)
		w.StopTime()
		w.SetTickRange(0)
		w.StopThundering()
		w.StopRaining()
		world.NewLoader(16, w, world.NopViewer{})
		ffas[g] = w
	}
}
