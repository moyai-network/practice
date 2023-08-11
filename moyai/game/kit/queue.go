package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
)

// Queue represents the kit given when players join the queue.
type Queue struct{}

// Items ...
func (Queue) Items(p *player.Player) [36]item.Stack {
	return [36]item.Stack{
		8: item.NewStack(item.Dye{Colour: item.ColourRed()}, 1).WithCustomName(
			"§r"+text.Colourf("<red>Leave queue</red>"),
		).WithValue("queue", 8),
	}
}

// Armour ...
func (Queue) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}

// Effects ...
func (Queue) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{}
}
