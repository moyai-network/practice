package duel

import "github.com/sandertv/gophertunnel/minecraft/protocol/packet"

type ReplayAction struct {
	Name   string
	Packet packet.Packet
}
