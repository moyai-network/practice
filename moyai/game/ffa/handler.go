package ffa

import (
	"github.com/df-mc/atomic"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose"
	"github.com/moyai-network/moose/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
	"time"
)

// Handler represents the player handler for FFA.
type Handler struct {
	*user.Handler
	p *player.Player
	g game.Game

	combat *moose.Tag
	pearl  *moose.CoolDown

	lastAttacker           atomic.Value[string]
	lastAttackerExpiration atomic.Value[time.Time]

	lobby func(*player.Player)
}

// newHandler returns a new FFA handler.
func newHandler(p *player.Player, g game.Game, lobby func(p *player.Player)) *Handler {
	h := &Handler{
		Handler: user.NewHandler(p),
		p:       p,
		g:       g,
		lobby:   lobby,
	}
	h.combat = moose.NewTag(h.tag, h.unTag)
	h.pearl = moose.NewCoolDown()
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
	}
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

		h.combat.Reset()
		h.pearl.Reset()

		h.lobby(h.p)

		killer, ok := data.LoadUser(h.lastAttacker.Load())
		if !ok || h.lastAttackerExpiration.Load().Before(time.Now()) {
			user.Broadcast("user.suicide", u.Roles.Highest().Colour(u.DisplayName))
			return
		}

		k, online := user.LookupXUID(killer.XUID)
		kh, ok := k.Handler().(*Handler)
		if online && ok {
			kh.sendScoreBoard()
			kh.combat.Cancel()
			kh.pearl.Reset()
			kit.Apply(h.g.Kit(), kh.p)
		}

		killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
		if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
			killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
		}
		killer = killer.WithKills(killer.Stats.Kills + 1)

		_ = data.SaveUser(killer)
		user.Broadcast("user.kill", u.Roles.Highest().Colour(u.DisplayName), killer.Roles.Highest().Colour(killer.DisplayName))
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	*force, *height = 0.394, 0.394
	target, ok := e.(*player.Player)
	if !ok {
		return
	}

	th, ok := target.Handler().(*Handler)
	if !ok {
		return
	}

	h.combat.Set(time.Second * 15)
	th.combat.Set(time.Second * 15)

	th.lastAttacker.Store(h.p.Name())
	th.lastAttackerExpiration.Store(time.Now().Add(time.Second * 15))
}

// bannedCommands is a list of commands disallowed in combat.
var bannedCommands = []string{"spawn", "rekit"}

func (h *Handler) HandleCommandExecution(ctx *event.Context, command cmd.Command, args []string) {
	if !h.combat.Active() {
		return
	}
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

func (h *Handler) HandleQuit() {
	if h.combat.Active() {
		u, _ := data.LoadUser(h.p.Name())
		if u.Stats.KillStreak > u.Stats.BestKillStreak {
			u = u.WithBestKillStreak(u.Stats.KillStreak)
		}
		u = u.WithDeaths(u.Stats.Deaths + 1)
		_ = data.SaveUser(u.WithKillStreak(0))

		killer, ok := data.LoadUser(h.lastAttacker.Load())
		if !ok || h.lastAttackerExpiration.Load().Before(time.Now()) {
			user.Broadcast("user.suicide", u.Roles.Highest().Colour(u.DisplayName))
			return
		}

		k, online := user.LookupXUID(killer.XUID)
		kh, ok := k.Handler().(*Handler)
		if online && ok {
			kh.sendScoreBoard()
			kh.combat.Cancel()
			kh.pearl.Reset()
			kit.Apply(h.g.Kit(), kh.p)
		}

		killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
		if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
			killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
		}
		killer = killer.WithKills(killer.Stats.Kills + 1)

		_ = data.SaveUser(killer)
		user.Broadcast("user.kill", u.Roles.Highest().Colour(u.DisplayName), killer.Roles.Highest().Colour(killer.DisplayName))
	}
	user.Remove(h.p)
}

func (h *Handler) tag(t *moose.Tag) {
	if !t.Active() {
		h.p.SendPopup(lang.Translatef(h.p.Locale(), "combat.tag"))
	}
}

func (h *Handler) unTag(t *moose.Tag) {
	h.p.SendPopup(lang.Translatef(h.p.Locale(), "combat.untag"))
}

func (h *Handler) sendScoreBoard() {
	l := h.p.Locale()
	u, _ := data.LoadUser(h.p.Name())

	sb := scoreboard.New(moose.GlyphFont("PRACTICE", item.ColourOrange()))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE000")

	_, _ = sb.WriteString(text.Colourf("<black>Kills</black><grey>: %d</grey>", u.Stats.Kills))
	_, _ = sb.WriteString(text.Colourf("<black>Deaths</black><grey>: %d</grey>", u.Stats.Deaths))
	_, _ = sb.WriteString(text.Colourf("<black>Kill Streak</black><grey>: %d</grey>", u.Stats.KillStreak))

	_, _ = sb.WriteString("\uE000")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE000") {
			sb.Set(i, " "+li)
		}
	}
	h.p.SendScoreboard(sb)
}
