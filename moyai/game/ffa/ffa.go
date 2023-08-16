package ffa

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/practice/moyai/structure"
	"math"
	"math/rand"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot/worlds"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game"
)

var lobby func(p *player.Player)

func InitializeLobby(f func(*player.Player)) {
	lobby = f
}

var ffas = map[game.Game]*world.World{}

func init() {
	pairs := [][]world.Block{
		{
			block.Concrete{Colour: item.ColourLime()},
			block.Concrete{Colour: item.ColourGreen()},
			block.GlazedTerracotta{Colour: item.ColourGreen()},
			block.Emerald{},
		},
		{
			block.Concrete{Colour: item.ColourLightBlue()},
			block.Concrete{Colour: item.ColourBlue()},
			block.Lapis{},
			block.Wool{Colour: item.ColourLightBlue()},
			block.Wool{Colour: item.ColourBlue()},
			block.StainedGlass{Colour: item.ColourBlue()},
		},
		{
			block.Amethyst{},
			block.Purpur{},
			block.Wool{Colour: item.ColourPurple()},
			block.Concrete{Colour: item.ColourPurple()},
			block.GlazedTerracotta{Colour: item.ColourPurple()},
			block.StainedGlass{Colour: item.ColourPurple()},
		},
		{
			block.NetherBricks{Type: block.CrackedNetherBricks()},
			block.NetherBricks{Type: block.ChiseledNetherBricks()},
			block.Wool{Colour: item.ColourRed()},
			block.Concrete{Colour: item.ColourRed()},
			block.Concrete{Colour: item.ColourBlack()},
			block.Coal{},
		},
	}

	for _, g := range game.Games() {
		size := 256
		w := world.Config{Entities: ent.Registry}.New()
		s := structure.GenerateBoxStructure([3]int{size, 50, size}, pairs[rand.Intn(len(pairs))]...)
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
		go world.NewLoader(16, w, world.NopViewer{}).Load(math.MaxInt)
		ffas[g] = w
	}
}
