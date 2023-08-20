package command

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/form"
)

type Stats struct {
	Targets cmd.Optional[[]cmd.Target] `cmd:"target"`
}

type StatsOffline struct {
	Target string `cmd:"target"`
}

// Run ...
func (st Stats) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}

	l := p.Locale()
	targets := st.Targets.LoadOr(nil)
	if len(targets) <= 0 {
		f := form.NewCasualStats(p.Name())
		p.SendForm(f)
		return
	}
	if len(targets) > 1 {
		o.Error(lang.Translatef(l, "command.targets.exceed"))
		return
	}
	tp := targets[0].(*player.Player)
	t, ok := data.LoadUser(tp.Name())
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	f := form.NewCasualStats(t.Name)
	p.SendForm(f)
}

// Run ...
func (st StatsOffline) Run(s cmd.Source, o *cmd.Output) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	l := p.Locale()

	t, ok := data.LoadUser(st.Target)
	if !ok {
		o.Error(lang.Translatef(l, "command.target.unknown"))
		return
	}
	f := form.NewCasualStats(t.Name)
	p.SendForm(f)
}
