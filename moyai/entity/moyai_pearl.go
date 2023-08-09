package entity

import (
	_ "unsafe"

	"github.com/df-mc/dragonfly/server/block/cube/trace"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/df-mc/dragonfly/server/world/particle"
	"github.com/df-mc/dragonfly/server/world/sound"
	"github.com/go-gl/mathgl/mgl64"
)

// NewMoyaiPearl creates an EnderPearl entity. EnderPearl is a smooth, greenish-
// blue item used to teleport.
func NewMoyaiPearl(pos mgl64.Vec3, vel mgl64.Vec3, owner world.Entity) world.Entity {
	e := entity.Config{Behaviour: moyaiPearlConf.New(owner)}.New(entity.EnderPearlType{}, pos)
	e.SetVelocity(vel)
	return e
}

var moyaiPearlConf = entity.ProjectileBehaviourConfig{
	Gravity:               0.085,
	Drag:                  0.01,
	KnockBackHeightAddend: 0.388 - 0.45,
	KnockBackForceAddend:  0.39 - 0.3608,
	Particle:              particle.EndermanTeleport{},
	Sound:                 sound.Teleport{},
	Hit:                   teleport,
}

// teleporter represents a living entity that can teleport.
type teleporter interface {
	// Teleport teleports the entity to the position given.
	Teleport(pos mgl64.Vec3)
	entity.Living
}

// teleport teleports the owner of an Ent to a trace.Result's position.
func teleport(e *entity.Ent, target trace.Result) {
	if u, ok := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter); ok {
		p := e.Behaviour().(*entity.ProjectileBehaviour).Owner().(teleporter).(*player.Player)
		rot := p.Rotation()
		onGround := p.OnGround()
		for _, v := range p.World().Viewers(p.Position()) {
			v.ViewEntityMovement(p, e.Position(), rot, onGround)
		}

		e.World().PlaySound(u.Position(), sound.Teleport{})
		u.Teleport(target.Position())
		u.Hurt(5, entity.FallDamageSource{})
	}
}
