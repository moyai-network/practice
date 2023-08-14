package duel

import (
	"sync"

	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

var (
	ongoing map[int64][]struct {
		Name   string
		Packet packet.Packet
	}
	ongoingMu sync.Mutex
)

func init() {
	ongoing = make(map[int64][]struct {
		Name   string
		Packet packet.Packet
	})
}

type ReplayAction struct {
	Name   string
	Packet packet.Packet
}
