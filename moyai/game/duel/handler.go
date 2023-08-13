package duel

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	*user.Handler
	p  *player.Player
	op *player.Player

	pearl     *carrot.CoolDown
	startTime time.Time

	replay []struct {
		Name   string
		Packet packet.Packet
	}
	replayMu sync.Mutex

	close chan struct{}
}

func newHandler(p *player.Player, op *player.Player) *Handler {
	var uHandler *user.Handler
	if uh, ok := p.Handler().(user.UserHandler); ok {
		uHandler = uh.UserHandler()
	} else {
		uHandler = user.NewHandler(p)
	}

	h := &Handler{
		Handler:   uHandler,
		p:         p,
		op:        op,
		startTime: time.Now(),

		close: make(chan struct{}, 0),
	}

	h.pearl = carrot.NewCoolDown(func(cd *carrot.CoolDown) {
		if !cd.Active() {
			h.p.SendPopup(lang.Translatef(h.p.Locale(), "pearl.cooldown"))
		}
	}, func(cd *carrot.CoolDown) {
		h.p.SendPopup(lang.Translatef(h.p.Locale(), "pearl.expired"))
	})

	t := time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-t.C:
				h.SendScoreBoard()
			case <-h.close:
				t.Stop()
				return
			}
		}
	}()
	return h
}
func (h *Handler) HandleItemUse(ctx *event.Context) {
	held, _ := h.p.HeldItems()
	switch held.Item().(type) {
	case item.EnderPearl:
		if h.pearl.Active() {
			ctx.Cancel()
			h.p.Message(text.Colourf("<red>You are on pearl cooldown</red>"))
			return
		}

		h.pearl.Set(time.Second * 15)
		h.SendScoreBoard()
	case item.SplashPotion:
	}
}

func (h *Handler) HandleMove(ctx *event.Context, pos mgl64.Vec3, pitch, yaw float64) {
}

func (h *Handler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	*damage = *damage / 1.25
	if src == (entity.FallDamageSource{}) {
		ctx.Cancel()
		return
	}

	if (h.p.Health()-h.p.FinalDamageFrom(*damage, src) <= 0) || src == (entity.VoidDamageSource{}) {
		ctx.Cancel()

		u, _ := data.LoadUser(h.p.Name())
		if u.Stats.KillStreak > u.Stats.BestKillStreak {
			u = u.WithBestKillStreak(u.Stats.KillStreak)
		}
		u = u.WithDeaths(u.Stats.Deaths + 1)
		_ = data.SaveUser(u.WithKillStreak(0))
		h.pearl.Cancel()

		killer, _ := data.LoadUser(h.op.Name())

		killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
		if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
			killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
		}
		killer = killer.WithKills(killer.Stats.Kills + 1)

		_ = data.SaveUser(killer)
		user.Broadcast("user.kill", u.Roles.Highest().Colour(u.DisplayName), potions(h.p), killer.Roles.Highest().Colour(killer.DisplayName), potions(h.op))

		h.UserHandler().SetRecentReplay(h.replay)
		h.op.Handler().(user.UserHandler).UserHandler().SetRecentReplay(h.replay)

		lobby(h.op)
		lobby(h.p)

	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	*force, *height = 0.394, 0.394
}

// bannedCommands is a list of commands disallowed in combat.
var bannedCommands = []string{"spawn", "rekit"}

func (h *Handler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	for _, bc := range bannedCommands {
		c, ok := cmd.ByAlias(bc)
		if !ok {
			continue
		}
		if c.Name() == command.Name() {
			h.p.Message(lang.Translatef(h.p.Locale(), "command.combat"))
			ctx.Cancel()
		}
	}
}

// HandleQuit ...
func (h *Handler) HandleQuit() {
	h.Close()
	u, _ := data.LoadUser(h.p.Name())
	if u.Stats.KillStreak > u.Stats.BestKillStreak {
		u = u.WithBestKillStreak(u.Stats.KillStreak)
	}
	u = u.WithDeaths(u.Stats.Deaths + 1)
	_ = data.SaveUser(u.WithKillStreak(0))

	h.pearl.Cancel()

	killer, _ := data.LoadUser(h.op.Name())

	killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
	if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
		killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
	}
	killer = killer.WithKills(killer.Stats.Kills + 1)

	_ = data.SaveUser(killer)
	user.Broadcast("user.kill", u.Roles.Highest().Colour(u.DisplayName), potions(h.p), killer.Roles.Highest().Colour(killer.DisplayName), potions(h.op))

	lobby(h.op)
}

// Close ...
func (h *Handler) Close() {
	h.pearl.Cancel()
	for _, e := range h.p.World().Entities() {
		if ent, ok := e.(*entity.Ent); ok {
			if be, ok := ent.Behaviour().(*entity.ProjectileBehaviour); ok {
				if be.Owner() == h.p {
					_ = e.Close()
				}
			}
		}
	}
	close(h.close)
}

func (h *Handler) AddReplayAction(p *player.Player, pk packet.Packet) {
	h.replayMu.Lock()
	h.replay = append(h.replay, ReplayAction{
		Name:   p.Name(),
		Packet: pk,
	})
	h.replayMu.Unlock()
}

// UserHandler ...
func (h *Handler) UserHandler() *user.Handler {
	return h.Handler
}

func (h *Handler) SendScoreBoard() {
	l := h.p.Locale()
	//u, _ := data.LoadUser(h.p.Name())

	sb := scoreboard.New(carrot.GlyphFont("PRACTICE"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE000")

	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>Opponent<grey>:</grey> <red>%s</red>", h.op.Name()))

	_, _ = sb.WriteString("\n\uE146\uE147\uE148\uE149\uE144\uE143")
	if h.pearl.Active() {
		_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>Ender Pearl<grey>:</grey> <black>%.0f</black>", h.pearl.Remaining().Seconds()))
	}

	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>Ping<grey>:</grey> <green>%dms</green> \uE145 <red>%dms</red>", h.p.Latency().Milliseconds()*2, h.op.Latency().Milliseconds()*2))
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>Time<grey>:</grey> <black>%s</black>", parseDuration(time.Since(h.startTime))))

	_, _ = sb.WriteString("\uE000")
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE000") {
			sb.Set(i, "  "+li)
		}
	}
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
	h.p.RemoveScoreboard()
	h.p.SendScoreboard(sb)
}

func parseDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	seconds := int(d.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}

	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

// potions returns the amount of potions the player has.
func potions(p *player.Player) (n int) {
	for _, i := range p.Inventory().Items() {
		if p, ok := i.Item().(item.SplashPotion); ok && p.Type == potion.StrongHealing() {
			n++
		}
	}
	return n
}
