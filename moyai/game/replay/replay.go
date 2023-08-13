package replay

import (
	"github.com/df-mc/dragonfly/server/item"
	"github.com/go-gl/mathgl/mgl64"
)

// Input is an input
type Input interface{}

// Movement is an input sent when a player moves
type Movement struct {
	Input
	// pos is the new position.
	Pos mgl64.Vec3
	// pitch is the new pitch
	Pitch float64
	// yaw is the new yaw
	Yaw float64
}

// Inventory is an input when a player
// changes item/slots in its inventory
type Inventory struct {
	Input
	// slot is the new selected slot.
	Slot float64
}

// Use is an input when a player
// uses either a pearl or a pot
type Use struct {
	Input
	// item is the item that was used
	Item item.Stack
}

// Swing is an input when a player swings its arm
type Swing struct {
	Input
}

// Hurt is an input when a player swings is hurt
type Hurt struct {
	Input
}

// Death is an input when a player dies
type Death struct {
	Input
}
