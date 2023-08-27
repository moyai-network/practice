package command

import (
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/samber/lo"
)

// PlayerInfo is a command that is used to info on an online player.
type PlayerInfo struct {
	Targets []cmd.Target `cmd:"target"`
}

// PlayerInfoOffline is a command that is used to info an offline player.
type PlayerInfoOffline struct {
	Target string `cmd:"target"`
}

// Run ...
func (b PlayerInfo) Run(src cmd.Source, o *cmd.Output) {
	l := locale(src)
	if len(b.Targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	t, ok := b.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	dT, ok := data.LoadUser(t.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	roles := strings.Join(lo.Map(dT.Roles.All(), func(r carrot.Role, _ int) string {
		return r.Colour(r.Name())
	}), ", ")
	addr := dT.Address[:9] + "..."
	// var ban, mute string
	// if b := dT.Punishments.Ban; !b.Expired() {
	// 	ban = text.Colourf("<gold>%s</gold", b.Occurrence)
	// }
	o.Print(lang.Translatef(l, "command.pinfo.info", dT.DisplayName, "<green>[Offline]</green>", dT.XUID, addr, dT.DeviceID, dT.SelfSignedID, dT.FirstLogin.Format(time.DateTime), dT.PlayTime.Round(time.Second), roles, "TODO", "TODO"))
}

// Run ...
func (b PlayerInfoOffline) Run(src cmd.Source, o *cmd.Output) {
	l := locale(src)
	dT, ok := data.LoadUser(b.Target)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	roles := strings.Join(lo.Map(dT.Roles.All(), func(r carrot.Role, _ int) string {
		return r.Colour(r.Name())
	}), ", ")
	addr := dT.Address[:9] + "..."
	// var ban, mute string
	// if b := dT.Punishments.Ban; !b.Expired() {
	// 	ban = text.Colourf("<gold>%s</gold", b.Occurrence)
	// }
	o.Print(lang.Translatef(l, "command.pinfo.info", dT.DisplayName, "<green>[Online]</green>", dT.XUID, addr, dT.DeviceID, dT.SelfSignedID, dT.FirstLogin.Format(time.DateTime), dT.PlayTime.Round(time.Second), roles, "TODO", "TODO"))
}

// Allow ...
func (PlayerInfo) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (PlayerInfoOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}
