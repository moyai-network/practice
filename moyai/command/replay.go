package command

import (
	"time"

	"github.com/df-mc/dragonfly/server/block"
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/go-gl/mathgl/mgl64"
	"github.com/moyai-network/carrot/sets"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game/lobby"
	"github.com/moyai-network/practice/moyai/game/structure"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/samber/lo"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

type ReplayRecent struct {
	Sub cmd.SubCommand `cmd:"recent"`
}

func (r ReplayRecent) Run(src cmd.Source, out *cmd.Output) {
	p, ok := src.(*player.Player)
	if !ok {
		return
	}
	u, ok := user.Lookup(p.Name())
	if h, ok := u.Handler().(user.UserHandler); ok {
		p.Message("debug: 1")
		h := h.UserHandler()
		if h.WatchingReplay() {
			p.Message("debug: already watching duel")
			return
		}

		p.Message("debug: 2")
		if r, ok := h.RecentReplay(); ok {
			p.Message("debug: 3")
			dim := [3]int{50, 20, 80}
			w := world.Config{Entities: ent.Registry, ReadOnly: true}.New()
			s := structure.GenerateBoxStructure(dim, block.Concrete{Colour: item.ColourBlack()})
			w.BuildStructure(cube.Pos{0, 0, 0}, s)
			w.AddEntity(p)
			p.Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 10})
			ss := sets.New(p.Name())
			names := lo.Map(r, func(item struct {
				Name   string
				Packet packet.Packet
			}, _ int) string {
				return item.Name
			})
			for _, n := range names {
				if !ss.Contains(n) {
					ss.Add(n)
				}
			}
			fakePlayerMap := map[string]*player.Player{}
			for n := range ss {
				fakePlayerMap[n] = player.New("Replay | "+n, p.Skin(), mgl64.Vec3{float64(dim[0] / 2), 2, 10})
				w.AddEntity(fakePlayerMap[n])
				fakePlayerMap[n].Teleport(mgl64.Vec3{float64(dim[0] / 2), 2, 10})
			}
			go func() {
				for _, action := range r {
					p, _ := fakePlayerMap[action.Name]
					pk, _ := action.Packet.(*packet.PlayerAuthInput)
					p.Move(mgl64.Vec3{float64(pk.Delta.X()), float64(pk.Delta.Y()), float64(pk.Delta.Z())}, float64(pk.Yaw), float64(pk.Pitch))
					time.Sleep(50 * time.Millisecond)
				}
				w.Close()
				lobby.AddPlayer(p)
			}()

		} else {
			p.Message("debug: no duel")
		}
	} else {
		p.Message("debug: no handler")
	}
}
