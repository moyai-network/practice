package command

import (
	"github.com/moyai-network/carrot/webhook"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/hako/durafmt"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/form"
	"github.com/moyai-network/practice/moyai/user"
	"go.mongodb.org/mongo-driver/bson"
)

// MuteForm is a command that is used to mute an online player through a punishment form.
type MuteForm struct{}

// MuteList is a command that outputs a list of muted players.
type MuteList struct {
	Sub cmd.SubCommand `cmd:"list"`
}

// MuteInfo is a command that displays the mute information of an online player.
type MuteInfo struct {
	Sub     cmd.SubCommand `cmd:"info"`
	Targets []cmd.Target   `cmd:"target"`
}

// MuteInfoOffline is a command that displays the mute information of an offline player.
type MuteInfoOffline struct {
	Sub    cmd.SubCommand `cmd:"info"`
	Target string         `cmd:"target"`
}

// MuteLift is a command that is used to lift the mute of an online player.
type MuteLift struct {
	Sub     cmd.SubCommand `cmd:"lift"`
	Targets []cmd.Target   `cmd:"target"`
}

// MuteLiftOffline is a command that is used to lift the mute of an offline player.
type MuteLiftOffline struct {
	Sub    cmd.SubCommand `cmd:"lift"`
	Target string         `cmd:"target"`
}

// Mute is a command that is used to mute an online player.
type Mute struct {
	Targets []cmd.Target `cmd:"target"`
	Reason  muteReason   `cmd:"reason"`
}

// MuteOffline is a command that is used to mute an offline player.
type MuteOffline struct {
	Target string     `cmd:"target"`
	Reason muteReason `cmd:"reason"`
}

// Run ...
func (MuteList) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	users, err := data.LoadUsersCond(
		bson.M{
			"$and": bson.A{
				bson.M{
					"punishments.mute.expiration": bson.M{"$ne": time.Time{}},
				}, bson.M{
					"punishments.mute.expiration": bson.M{"$gt": time.Now()},
				},
			},
		},
	)
	if err != nil {
		panic(err)
	}
	if len(users) == 0 {
		o.Error(lang.Translatef(l, "command.mute.none"))
		return
	}
	o.Print(lang.Translatef(l, "command.mute.list", len(users), strings.Join(names(users, true), ", ")))
}

// Run ...
func (m MuteInfo) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	p, ok := m.Targets[0].(*player.Player)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Punishments.Mute.Expired() {
		o.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	mute := u.Punishments.Mute
	o.Print(lang.Translatef(l, "punishment.details",
		p.Name(),
		mute.Reason,
		durafmt.Parse(mute.Remaining()),
		mute.Staff,
		mute.Occurrence.Format("01/02/2006"),
	))
}

// Run ...
func (m MuteInfoOffline) Run(s cmd.Source, o *cmd.Output) {
	l := locale(s)
	u, err := data.LoadUserOrCreate(m.Target)
	if err != nil {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Punishments.Mute.Expired() {
		o.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	o.Print(lang.Translatef(l, "punishment.details",
		u.DisplayName,
		u.Punishments.Mute.Reason,
		durafmt.Parse(u.Punishments.Mute.Remaining()),
		u.Punishments.Mute.Staff,
		u.Punishments.Mute.Occurrence.Format("01/02/2006"),
	))
}

// Run ...
func (m MuteLift) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	p, ok := m.Targets[0].(*player.Player)
	if !ok {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Punishments.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	u.Punishments.Mute = carrot.Punishment{}
	_ = data.SaveUser(u)

	user.Alert(src, "staff.alert.unmute", p.Name())

	webhook.SendPunishment(src.(cmd.NamedTarget).Name(), u.DisplayName, "Lift", webhook.UnMutePunishment())
	out.Print(lang.Translatef(l, "command.mute.lift", p.Name()))
}

// Run ...
func (m MuteLiftOffline) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	u, err := data.LoadUserOrCreate(m.Target)
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Punishments.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.not"))
		return
	}
	u.Punishments.Mute = carrot.Punishment{}
	_ = data.SaveUser(u)

	user.Alert(src, "staff.alert.unmute", u.DisplayName)

	webhook.SendPunishment(src.(cmd.NamedTarget).Name(), u.DisplayName, "Lift", webhook.UnMutePunishment())
	out.Print(lang.Translatef(l, "command.mute.lift", u.DisplayName))
}

// Run ...
func (m MuteForm) Run(s cmd.Source, _ *cmd.Output) {
	p := s.(*player.Player)
	p.SendForm(form.NewMute(p))
}

// Run ...
func (m Mute) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	if len(m.Targets) > 1 {
		out.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	t, ok := m.Targets[0].(*player.Player)
	if !ok {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if t == src {
		out.Error(lang.Translatef(l, "command.mute.self"))
		return
	}
	u, err := data.LoadUserOrCreate(t.Name())
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	if u.Roles.Contains(role.Operator{}) {
		out.Error(lang.Translatef(l, "command.mute.operator"))
		return
	}
	if !u.Punishments.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.already"))
		return
	}
	sn := src.(cmd.NamedTarget)
	reason, length := parseMuteReason(m.Reason)
	u.Punishments.Mute = carrot.Punishment{
		Staff:      sn.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	_ = data.SaveUser(u)

	user.Alert(src, "staff.alert.mute", t.Name(), reason)

	webhook.SendPunishment(src.(cmd.NamedTarget).Name(), u.DisplayName, reason, webhook.MutePunishment())
	out.Print(lang.Translatef(l, "command.mute.success", t.Name(), reason))
}

// Run ...
func (m MuteOffline) Run(src cmd.Source, out *cmd.Output) {
	l := locale(src)
	sn := src.(cmd.NamedTarget)

	u, err := data.LoadUserOrCreate(m.Target)
	if err != nil {
		out.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}

	if strings.EqualFold(u.Name, m.Target) {
		out.Error(lang.Translatef(l, "command.mute.self"))
		return
	}

	if u.Roles.Contains(role.Operator{}) {
		out.Error(lang.Translatef(l, "command.mute.operator"))
		return
	}
	if !u.Punishments.Mute.Expired() {
		out.Error(lang.Translatef(l, "command.mute.already"))
		return
	}

	reason, length := parseMuteReason(m.Reason)
	u.Punishments.Mute = carrot.Punishment{
		Staff:      sn.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	_ = data.SaveUser(u)

	user.Alert(src, "staff.alert.mute", u.DisplayName, reason)

	webhook.SendPunishment(src.(cmd.NamedTarget).Name(), u.DisplayName, reason, webhook.MutePunishment())
	out.Print(lang.Translatef(l, "command.mute.success", u.DisplayName, reason))
}

// Allow ...
func (MuteList) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteInfo) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteInfoOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteForm) Allow(s cmd.Source) bool {
	return allow(s, false, role.Trial{})
}

// Allow ...
func (Mute) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteLift) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

// Allow ...
func (MuteLiftOffline) Allow(s cmd.Source) bool {
	return allow(s, true, role.Trial{})
}

type (
	muteReason string
)

// Type ...
func (muteReason) Type() string {
	return "muteReason"
}

// Options ...
func (muteReason) Options(cmd.Source) []string {
	return []string{
		"spam",
		"toxic",
		"advertising",
	}
}

// parseMuteReason returns the formatted muteReason and mute duration.
func parseMuteReason(r muteReason) (string, time.Duration) {
	switch r {
	case "spam":
		return "Spam", time.Hour * 6
	case "toxic":
		return "Toxicity", time.Hour * 9
	case "advertising":
		return "Advertising", time.Hour * 24 * 3
	}
	panic("should never happen")
}
