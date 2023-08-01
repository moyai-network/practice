package lobby

import (
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
)

type inventoryHandler struct {
	inventory.NopHandler
}

// HandleTake ...
func (inventoryHandler) HandleTake(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandlePlace ...
func (inventoryHandler) HandlePlace(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}

// HandleDrop ...
func (inventoryHandler) HandleDrop(ctx *event.Context, _ int, _ item.Stack) {
	ctx.Cancel()
}
