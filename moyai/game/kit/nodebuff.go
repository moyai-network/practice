package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/enchantment"
	"github.com/df-mc/dragonfly/server/item/potion"
	"github.com/df-mc/dragonfly/server/player"
	"time"
)

// NoDebuff represents the kit given when players join NoDebuff.
type NoDebuff struct {
}

// Items ...
func (n NoDebuff) Items(*player.Player) [36]item.Stack {
	items := [36]item.Stack{
		item.NewStack(item.Sword{Tier: item.ToolTierDiamond}, 1).WithEnchantments(item.NewEnchantment(enchantment.Unbreaking{}, 10)),
		item.NewStack(item.EnderPearl{}, 16),
	}
	for i := 2; i < 36; i++ {
		items[i] = item.NewStack(item.SplashPotion{Type: potion.StrongHealing()}, 1)
	}
	return items
}

// Armour ...
func (NoDebuff) Armour(*player.Player) [4]item.Stack {
	durability := item.NewEnchantment(enchantment.Unbreaking{}, 10)
	return [4]item.Stack{
		item.NewStack(item.Helmet{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability),
		item.NewStack(item.Chestplate{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability),
		item.NewStack(item.Leggings{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability),
		item.NewStack(item.Boots{Tier: item.ArmourTierDiamond{}}, 1).WithEnchantments(durability),
	}
}

// Effects ...
func (n NoDebuff) Effects(*player.Player) []effect.Effect {
	return []effect.Effect{effect.New(effect.Speed{}, 1, time.Hour*24).WithoutParticles()}
}
