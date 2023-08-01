package kit

import (
	"github.com/df-mc/dragonfly/server/entity/effect"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/session"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/moose"
	_ "unsafe"
)

// Apply ...
func Apply(kit moose.Kit, p *player.Player) {
	p.Inventory().Clear()
	p.Armour().Clear()

	p.SetHeldItems(item.Stack{}, item.Stack{})
	if s := player_session(p); s != session.Nop {
		_ = s.SetHeldSlot(0)
	}

	p.StopSneaking()
	p.StopSwimming()
	p.StopSprinting()
	p.StopFlying()
	p.ResetFallDistance()
	p.SetGameMode(world.GameModeSurvival)

	p.Heal(20, effect.InstantHealingSource{})
	p.SetFood(20)
	for _, eff := range p.Effects() {
		p.RemoveEffect(eff.Type())
	}

	inv := p.Inventory()
	armour := kit.Armour(p)
	for slot, it := range kit.Items(p) {
		_ = inv.SetItem(slot, it)
	}
	for _, eff := range kit.Effects(p) {
		p.AddEffect(eff)
	}
	p.Armour().Set(armour[0], armour[1], armour[2], armour[3])
}

// noinspection ALL
//
//go:linkname player_session github.com/df-mc/dragonfly/server/player.(*Player).session
func player_session(*player.Player) *session.Session
