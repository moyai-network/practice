package duel

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/title"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot/worlds"
	"github.com/moyai-network/practice/moyai/data"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/structure"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"time"
)

type Match struct {
	id int64

	players [2]*player.Player

	g game.Game
	w *world.World

	beginning time.Time
	running   bool

	competitive bool
}

func NewMatch(p, op *player.Player, g game.Game, competitive bool) *Match {
	m := &Match{
		id:      rand.Int63(),
		players: [2]*player.Player{p, op},

		g: g,

		competitive: competitive,
	}
	return m
}

func (m *Match) Start() {
	dim := [3]int{50, 20, 80}
	m.prepare(dim)

	p1, p2 := m.players[0], m.players[1]

	game.DeQueue(p1)
	game.DeQueue(p2)

	p1.Inventory().Handle(inventory.NopHandler{})
	p1.Armour().Handle(inventory.NopHandler{})
	p1.SetNameTag(text.Colourf("<red>%s</red>", p1.Name()))
	p1.Handle(newHandler(p1, p2, m))
	m.w.AddEntity(p1)
	p1.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 10})
	kit.Apply(m.g.Kit(), p1)

	p2.Inventory().Handle(inventory.NopHandler{})
	p2.Armour().Handle(inventory.NopHandler{})
	p2.SetNameTag(text.Colourf("<red>%s</red>", p2.Name()))
	p2.Handle(newHandler(p2, p1, m))
	m.w.AddEntity(p2)
	p2.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 70})
	kit.Apply(m.g.Kit(), p2)

	m.beginning = time.Now()
	m.running = true
}

func (m *Match) prepare(dim [3]int) {
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

	m.w = w
}

func (m *Match) End(winner, loser *player.Player, forced bool) {
	if !m.running {
		return
	}
	m.running = false
	w := m.w

	u, _ := data.LoadUser(loser.Name())
	if u.Stats.KillStreak > u.Stats.BestKillStreak {
		u = u.WithBestKillStreak(u.Stats.KillStreak)
	}
	u = u.WithDeaths(u.Stats.Deaths + 1).WithIncreasedLoss(m.competitive)

	killer, _ := data.LoadUser(winner.Name())

	killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
	if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
		killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
	}
	killer = killer.WithKills(killer.Stats.Kills + 1).WithIncreasedWin(m.competitive)

	if m.competitive {
		earnings, losings := eloEarnings(u.GameElo(m.g), killer.GameElo(m.g)), eloLosings(killer.GameElo(m.g), u.GameElo(m.g))

		killer = killer.WithElo(m.g, killer.GameElo(m.g)+earnings)
		u = u.WithElo(m.g, u.GameElo(m.g)-losings)

		if !forced {
			loser.Message(text.Colourf("<red>New elo: %d (-%d)</red>", u.GameElo(m.g), losings))
		}

		winner.Message(text.Colourf("<green>New elo: %d (+%d)</green>", killer.GameElo(m.g), earnings))
	}

	_ = data.SaveUser(killer)
	_ = data.SaveUser(u.WithKillStreak(0))

	msg := text.Colourf("<green>Winner: </green><yellow>%s <grey>-</grey> <red>Loser: </red>%s</yellow>", winner.Name(), loser.Name())
	winner.Message(msg)
	loser.Message(msg)

	winner.SendTitle(title.New(text.Colourf("<b><green>VICTORY</green></b>\n <white>You won the match</white>")))
	loser.SendTitle(title.New(text.Colourf("<b><red>DEFEAT</red></b>\n <white>You lost the match</white>")))

	loser.SetGameMode(world.GameModeSpectator)
	loser.KnockBack(winner.Position(), 0.5, 0.5)

	loser.Inventory().Clear()
	winner.Inventory().Clear()
	loser.Armour().Clear()
	winner.Armour().Clear()

	c := player.New(loser.Name(), loser.Skin(), loser.Position())
	c.SetAttackImmunity(time.Millisecond * 1400)
	c.SetNameTag(loser.NameTag())
	c.SetScale(loser.Scale())
	w.AddEntity(c)

	for _, viewer := range w.Viewers(c.Position()) {
		viewer.ViewEntityAction(c, entity.DeathAction{})
	}

	c.KnockBack(winner.Position(), 0.5, 0.2)
	time.AfterFunc(time.Millisecond*1400, func() {
		_ = c.Close()
	})

	time.AfterFunc(time.Second*3, func() {
		if !forced {
			lobby(loser)
		}
		lobby(winner)
	})
}
