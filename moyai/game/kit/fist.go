package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
)

// Fist represents the kit given when players join Fist.
type Fist struct{}

// Items ...
func (Fist) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Beef{Cooked: true}, 64),
	}
	return items
}

// Armour ...
func (Fist) Armour(*player.Player) [4]item.Stack {
	return [4]item.Stack{}
}

// Effects ...
func (Fist) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{}
}
