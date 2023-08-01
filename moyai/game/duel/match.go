package duel

import (
	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/moose/worlds"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/structure"
	"math/rand"
	"time"
)

type Match struct {
	opponents [2]*player.Player
	g         game.Game
	startedAt time.Time

	duel bool
}

func NewDuel(p1, p2 *player.Player, g game.Game) *Match {
	m := &Match{opponents: [2]*player.Player{p1, p2}, g: g}
	return m
}

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

func (m *Match) Start() {
	dim := [3]int{50, 20, 80}
	w := world.Config{Entities: ent.Registry, ReadOnly: true}.New()
	s := structure.GenerateBoxStructure(dim, pairs[rand.Intn(len(pairs))]...)
	w.BuildStructure(cube.Pos{0, 0, 0}, s)
	w.Handle(&worlds.Handler{})
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeAdventure)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()

	p1 := m.opponents[0]
	p2 := m.opponents[1]

	w.AddEntity(p1)
	p1.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 10})
	AddPlayer(p1, m)

	w.AddEntity(p2)
	p2.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 70})
	AddPlayer(p2, m)

	m.startedAt = time.Now()
}

func (m *Match) End(winner *player.Player) {
}

func (m *Match) Beginning() time.Time {
	return m.startedAt
}
