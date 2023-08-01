package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
)

type Handler struct {
	player.NopHandler
	p *player.Player
}

func NewHandler(p *player.Player) *Handler {
	h := &Handler{p: p}
	return h
}

// HandleFoodLoss ...
func (*Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

// HandleBlockBreak ...
func (*Handler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	ctx.Cancel()
}

// HandleBlockPlace ...
func (*Handler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	ctx.Cancel()
}
