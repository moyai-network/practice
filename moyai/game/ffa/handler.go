package ffa

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/df-mc/atomic"
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
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Handler represents the player handler for FFA.
type Handler struct {
	*user.Handler
	p *player.Player
	g game.Game

	combat *carrot.Tag
	pearl  *carrot.CoolDown

	lastAttacker           atomic.Value[string]
	lastAttackerExpiration atomic.Value[time.Time]

	close chan struct{}
}

// newHandler returns a new FFA handler.
func newHandler(p *player.Player, g game.Game) *Handler {
	h := &Handler{
		Handler: p.Handler().(user.UserHandler).UserHandler(),
		p:       p,
		g:       g,

		close: make(chan struct{}, 0),
	}
	h.combat = carrot.NewTag(h.tag, h.unTag)
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
	}
}

func (h *Handler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	switch h.g {
	case game.Fist():
		*attackImmunity = 350 * time.Millisecond
	default:
		*attackImmunity = 490 * time.Millisecond
		*damage = *damage / 1.25
	}
	if src == (entity.FallDamageSource{}) {
		ctx.Cancel()
		return
	}
	w := h.p.World()

	if (h.p.Health()-h.p.FinalDamageFrom(*damage, src) <= 0) || src == (entity.VoidDamageSource{}) {
		ctx.Cancel()

		u, _ := data.LoadUser(h.p.Name())
		if u.Stats.KillStreak > u.Stats.BestKillStreak {
			u = u.WithBestKillStreak(u.Stats.KillStreak)
		}
		u = u.WithDeaths(u.Stats.Deaths + 1)
		_ = data.SaveUser(u.WithKillStreak(0))

		h.combat.Cancel()
		h.pearl.Cancel()

		killer, ok := data.LoadUser(h.lastAttacker.Load())
		if !ok || h.lastAttackerExpiration.Load().Before(time.Now()) {
			user.Broadcast("user.suicide", u.Roles.Highest().Colour(u.DisplayName))
			return
		}

		k, online := user.LookupXUID(killer.XUID)
		if online {
			killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
			if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
				killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
			}
			killer = killer.WithKills(killer.Stats.Kills + 1)

			_ = data.SaveUser(killer)
			if h.g == game.NoDebuff() {
				user.Broadcast("user.kill.pots", killer.DisplayName, potions(k), u.DisplayName, potions(h.p))
			} else {
				user.Broadcast("user.kill", killer.DisplayName, u.DisplayName)
			}

			c := player.New(h.p.Name(), h.p.Skin(), h.p.Position())
			c.SetAttackImmunity(time.Millisecond * 1400)
			c.SetNameTag(h.p.NameTag())
			c.SetScale(h.p.Scale())
			w.AddEntity(c)

			for _, viewer := range w.Viewers(c.Position()) {
				viewer.ViewEntityAction(c, entity.DeathAction{})
			}

			c.KnockBack(k.Position(), 0.5, 0.2)
			time.AfterFunc(time.Millisecond*1400, func() {
				_ = c.Close()
			})
		}
		kh, ok := k.Handler().(*Handler)
		if online && ok {
			kh.SendScoreBoard()
			kh.combat.Cancel()
			kh.pearl.Reset()
			kit.Apply(h.g.Kit(), kh.p)
		}
		lobby(h.p)
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	switch h.g {
	case game.Fist():
		*force, *height = 0.4, 0.375
	default:
		*force, *height = 0.38, 0.38
	}

	target, ok := e.(*player.Player)
	if !ok {
		return
	}

	if !target.OnGround() {
		max, min := maxMin(target.Position().Y(), h.p.Position().Y())
		if max-min >= 2.5 {
			*height = 0.38 / 1.25
		}
	}

	th, ok := target.Handler().(*Handler)
	if !ok {
		return
	}

	u, _ := data.LoadUser(h.p.Name())
	if *critical && u.Settings.Gameplay.CriticalEffect {
		ent := entity.NewText(text.Colourf("<red>Critical!</red>"), th.p.Position().Add(mgl64.Vec3{rand.Float64(), 1, rand.Float64()}))
		for _, viewer := range h.p.World().Viewers(ent.Position()) {
			viewer.HideEntity(ent)
		}
		h.p.World().AddEntity(ent)
		h.p.ShowEntity(ent)
		time.AfterFunc(time.Millisecond*500, func() {
			_ = ent.Close()
		})
	}

	h.combat.Set(time.Second * 15)
	th.combat.Set(time.Second * 15)

	th.lastAttacker.Store(h.p.Name())
	th.lastAttackerExpiration.Store(time.Now().Add(time.Second * 15))
	h.SendScoreBoard()
}

func maxMin(n, n2 float64) (max float64, min float64) {
	if n > n2 {
		return n, n2
	}
	return n2, n
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
	defer user.Remove(h.UserHandler())
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
		if online {
			killer = killer.WithKillStreak(killer.Stats.KillStreak + 1)
			if killer.Stats.KillStreak > killer.Stats.BestKillStreak {
				killer = killer.WithBestKillStreak(killer.Stats.KillStreak)
			}
			killer = killer.WithKills(killer.Stats.Kills + 1)

			_ = data.SaveUser(killer)
			if h.g == game.NoDebuff() {
				user.Broadcast("user.kill.pots", killer.Roles.Highest().Colour(killer.DisplayName), potions(k), u.Roles.Highest().Colour(u.DisplayName), potions(h.p))
			} else {
				user.Broadcast("user.kill", u.Roles.Highest().Colour(u.DisplayName), killer.Roles.Highest().Colour(killer.DisplayName))
			}
		}

		kh, ok := k.Handler().(*Handler)
		if online && ok {
			kh.SendScoreBoard()
			kh.combat.Cancel()
			kh.pearl.Reset()
			kit.Apply(h.g.Kit(), kh.p)
		}
	}
}

// UserHandler ...
func (h *Handler) UserHandler() *user.Handler {
	return h.Handler
}

// Close ...
func (h *Handler) Close() {
	h.combat.Cancel()
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

func (h *Handler) SendScoreBoard() {
	l := h.p.Locale()
	u, _ := data.LoadUser(h.p.Name())

	if !u.Settings.Display.Scoreboard {
		return
	}

	sb := scoreboard.New(carrot.GlyphFont("Moyai"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE002")

	_, _ = sb.WriteString(text.Colourf("<red>Kills:</red> <white>%d</white>", u.Stats.Kills))
	_, _ = sb.WriteString(text.Colourf("<red>Killstreak:</red> <white>%d</white>", u.Stats.KillStreak))
	_, _ = sb.WriteString(text.Colourf("<red>Deaths:</red> <white>%d</white>", u.Stats.Deaths))

	if h.combat.Active() {
		_, _ = sb.WriteString(text.Colourf("<red>Combat:</red> <white>%s</white>", parseDuration(h.combat.Remaining())))
	}

	_, _ = sb.WriteString("§a")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

	_, _ = sb.WriteString("\uE002")
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE002") {
			sb.Set(i, " "+li)
		}
	}

	_, _ = sb.WriteString("\uE002")
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

func (h *Handler) Game() game.Game {
	return h.g
}

func (h *Handler) tag(t *carrot.Tag) {
	if !t.Active() {
		h.p.SendPopup(lang.Translatef(h.p.Locale(), "combat.tag"))
	}
}

func (h *Handler) unTag(t *carrot.Tag) {
	h.p.SendPopup(lang.Translatef(h.p.Locale(), "combat.untag"))
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
