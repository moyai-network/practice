package duel

import (
	"fmt"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

type Handler struct {
	*user.Handler
	m *Match

	p  *player.Player
	op *player.Player

	hits  int
	pearl *carrot.CoolDown
}

func newHandler(p *player.Player, op *player.Player, m *Match) *Handler {
	h := &Handler{
		Handler: p.Handler().(user.UserHandler).UserHandler(),
		m:       m,
		p:       p,
		op:      op,
	}

	ongoingMu.Lock()
	defer ongoingMu.Unlock()
	ongoing[m.id] = []struct {
		Name   string
		Packet packet.Packet
	}{}

	h.pearl = carrot.NewCoolDown(func(cd *carrot.CoolDown) {
		if !cd.Active() {
			h.p.SendPopup(lang.Translatef(h.p.Locale(), "pearl.cooldown"))
		}
	}, func(cd *carrot.CoolDown) {
		h.p.SendPopup(lang.Translatef(h.p.Locale(), "pearl.expired"))
	})

	h.SendScoreBoard()
	t := time.NewTicker(time.Second)
	go func() {
		for range t.C {
			if !m.running {
				t.Stop()
				return
			}
			h.SendScoreBoard()
		}
	}()
	return h
}
func (h *Handler) HandleItemUse(ctx *event.Context) {
	if !h.m.running {
		return
	}
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

func (h *Handler) HandleHurt(ctx *event.Context, damage *float64, attackImmunity *time.Duration, src world.DamageSource) {
	if !h.m.running {
		return
	}
	switch h.m.g {
	default:
		*attackImmunity = 490 * time.Millisecond
	}
	*damage = *damage / 1.25
	if src == (entity.FallDamageSource{}) {
		ctx.Cancel()
		return
	}

	if s, ok := src.(entity.AttackDamageSource); ok {
		if s.Attacker != h.op {
			ctx.Cancel()
			return
		}
	}

	oph := h.op.Handler().(*Handler)
	oph.hits++

	h.SendScoreBoard()
	oph.SendScoreBoard()

	if h.m.g == game.Boxing() {
		*damage = 0
	}

	if (h.p.Health()-h.p.FinalDamageFrom(*damage, src) <= 0) || src == (entity.VoidDamageSource{}) || (h.m.g == game.Boxing() && h.op.Handler().(*Handler).hits >= 100) {
		ctx.Cancel()

		ongoingMu.Lock()
		defer ongoingMu.Unlock()
		h.UserHandler().SetRecentReplay(ongoing[h.m.id])
		h.op.Handler().(user.UserHandler).UserHandler().SetRecentReplay(ongoing[h.m.id])

		h.m.End(h.op, h.p, false)
	}
}

func (h *Handler) HandleAttackEntity(ctx *event.Context, e world.Entity, force, height *float64, critical *bool) {
	if !h.m.running {
		return
	}
	switch h.m.g {
	case game.Boxing():
		*critical = false
		*force, *height = 0.38, 0.38
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
	h.m.End(h.op, h.p, true)
	h.Close()
	user.Remove(h.UserHandler())
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
}

func (h *Handler) AddReplayAction(p *player.Player, pk packet.Packet) {
	ongoingMu.Lock()
	defer ongoingMu.Unlock()
	ongoing[h.m.id] = append(ongoing[h.m.id], struct {
		Name   string
		Packet packet.Packet
	}{
		Name:   p.Name(),
		Packet: pk,
	})
}

// UserHandler ...
func (h *Handler) UserHandler() *user.Handler {
	return h.Handler
}

func (h *Handler) SendScoreBoard() {
	if !h.m.running {
		return
	}

	l := h.p.Locale()
	u, _ := data.LoadUser(h.p.Name())

	if !u.Settings.Display.Scoreboard {
		return
	}

	sb := scoreboard.New(carrot.GlyphFont("Moyai"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE002")

	_, _ = sb.WriteString(text.Colourf("Duration: <red>%s</red>", parseDuration(time.Since(h.m.beginning))))

	if h.m.g == game.Boxing() {
		_, _ = sb.WriteString(text.Colourf("\n\nHits:\n  You: <red>%d</red>\n  Them: <red>%d</red>", h.hits, h.op.Handler().(*Handler).hits))
	}

	_, _ = sb.WriteString(text.Colourf("\nYour Ping: <red>%dms</red>\nTheir Ping: <red>%dms</red>", h.p.Latency().Milliseconds()*2, h.op.Latency().Milliseconds()*2))

	_, _ = sb.WriteString("§a")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

	_, _ = sb.WriteString("\uE002")
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE002") {
			sb.Set(i, " "+li)
		}
	}
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
