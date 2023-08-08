package lobby

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/form"
	"github.com/moyai-network/practice/moyai/user"
	"time"
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
	_ = data.SaveUser(u)

	return h
}

// HandleItemUse ...
func (h *Handler) HandleItemUse(_ *event.Context) {
	held, _ := h.p.HeldItems()

	val, ok := held.Value("lobby")
	if !ok {
		return
	}
	switch val {
	case 0:
		h.p.SendForm(form.NewFFA(AddPlayer))
	}
}

// HandleQuit ...
func (h *Handler) HandleQuit() {
	user.Remove(h.p)
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
