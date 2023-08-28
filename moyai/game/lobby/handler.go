package lobby

import (
	"encoding/hex"
	"net/netip"
	"strings"
	"time"

	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/form"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/moyai-network/practice/moyai/game/kit"
	"golang.org/x/crypto/sha3"

	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Handler represents the player handler for the lobby.
type Handler struct {
	*user.Handler
	p *player.Player
}

// newHandler returns a new lobby handler.
func newHandler(p *player.Player) *Handler {
	var uHandler *user.Handler
	if uh, ok := p.Handler().(user.UserHandler); ok {
		uHandler = uh.UserHandler()
	} else {
		uHandler = user.NewHandler(p)
	}

	h := &Handler{
		Handler: uHandler,
		p:       p,
	}

	u := data.LoadOrCreateUser(p.Name())
	u.DisplayName = p.Name()
	u.XUID = p.XUID()

	if s := player_session(p); s != session.Nop {
		u.DeviceID = s.ClientData().DeviceID
		u.SelfSignedID = s.ClientData().SelfSignedID
		sha := sha3.New256()
		addr, _ := netip.ParseAddrPort(p.Addr().String())
		sha.Write(addr.Addr().AsSlice())
		sha.Write([]byte(data.Salt))
		u.Address = hex.EncodeToString(sha.Sum(nil))
	}

	_ = data.SaveUser(u)

	return h
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(_ *event.Context) {
	held, _ := h.p.HeldItems()

	val, ok := held.Value("lobby")
	if ok {
		switch val {
		case 0:
			h.p.SendForm(form.NewFFA())
		case 1:
			h.p.SendForm(form.NewCasualQueue())
		case 2:
			h.p.SendForm(form.NewCompetitiveQueue())
		case 8:
			h.p.SendForm(form.NewSettings())
		}
	}

	val, ok = held.Value("queue")
	if ok {
		switch val {
		case 8:
			kit.Apply(kit.Lobby{}, h.p)
			game.DeQueue(h.p)
		}
	}
}

// HandleAttackEntity ...
func (h *Handler) HandleAttackEntity(ctx *event.Context, _ world.Entity, _, _ *float64, _ *bool) {
	ctx.Cancel()
}

// HandleHurt ...
func (h *Handler) HandleHurt(ctx *event.Context, _ *float64, _ *time.Duration, _ world.DamageSource) {
	ctx.Cancel()
}

// UserHandler ...
func (h *Handler) UserHandler() *user.Handler {
	return h.Handler
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

	var playing int
	for _, u := range user.All() {
		if _, ok := u.Handler().(*Handler); !ok {
			playing++
		}
	}

	_, _ = sb.WriteString(text.Colourf("<red>Online:</red> <white>%d</white>", user.Count()))
	_, _ = sb.WriteString(text.Colourf("<red>Playing:</red> <white>%d</white>", playing))

	_, _ = sb.WriteString("§a")
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))

	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE002") {
			sb.Set(i, " "+li)
		}
	}

	_, _ = sb.WriteString("\uE002")
	h.p.RemoveScoreboard()
	h.p.SendScoreboard(sb)
}
func (h *Handler) HandleQuit() {
	user.Remove(h.UserHandler())
	game.DeQueue(h.p)
	h.Close()
}

func (h *Handler) Close() {
	for _, u := range user.All() {
		if h, ok := u.Handler().(*Handler); ok {
			h.SendScoreBoard()
		}
	}
}
