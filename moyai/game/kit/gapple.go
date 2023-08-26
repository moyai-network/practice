package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

// Gapple represents the kit given when players join Gapple.
type Gapple struct {
}

// Items ...
func (n Gapple) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 10), item.NewEnchantment(enchantment.Sharpness{}, 2)).WithCustomName(text.Colourf("<red>Moyai</red>")),
		item.NewStack(item.GoldenApple{}, 8),
	}

	return items
}

// Armour ...
func (Gapple) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking{}, 10)
	protection := item.NewEnchantment(enchantment.Protection{}, 2)
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(protection, durability),
	}
}

// Effects ...
func (Gapple) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{effect.New(effect.Speed{}, 1, time.Hour*24).WithoutParticles()}
}
