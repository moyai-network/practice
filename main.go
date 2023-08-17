package main

import (
	"os"

	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	_ "github.com/flonja/multiversion/protocols" // VERY IMPORTANT
	v486 "github.com/flonja/multiversion/protocols/v486"
	_ "github.com/moyai-network/carrot/console"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/worlds"
	"github.com/moyai-network/practice/moyai"
	"github.com/moyai-network/practice/moyai/command"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game/lobby"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/oomph-ac/oomph"
	"github.com/oomph-ac/oomph/utils"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
)

func main() {
	lang.Register(language.English)

	log := logrus.New()
	log.Formatter = &logrus.TextFormatter{ForceColors: true}
	log.Level = logrus.InfoLevel

	config, err := readConfig()
	if err != nil {
		log.Fatalln(err)
	}
	chat.Global.Subscribe(chat.StdoutSubscriber{})

	c, err := config.Config(log)
	if err != nil {
		panic(err)
	}

	c.ReadOnlyWorld = true
	c.Entities = ent.Registry
	c.Allower = &moyai.Allower{}

	c.Name = text.Colourf("<bold><purple>MOYAI</purple></bold>") + "§8"

	ac := oomph.New(log, ":19132")
	ac.Listen(&c, c.Name, []minecraft.Protocol{v486.New()}, false, false)
	go func() {
		for {
			p, err := ac.Accept()
			if err != nil {
				return
			}
			p.SetCombatMode(utils.AuthorityType(config.Oomph.CombatMode))
			p.SetMovementMode(utils.AuthorityType(config.Oomph.MovementMode))
			p.SetCombatCutoff(3)    // 2 ticks => 100ms
			p.SetKnockbackCutoff(3) // 2 ticks => 100ms
			p.Handle(user.NewOomphHandler(p))
		}
	}()

	srv := c.New()
	srv.CloseOnProgramEnd()

	w := srv.World()
	w.Handle(&worlds.Handler{})
	w.StopWeatherCycle()
	w.SetDefaultGameMode(world.GameModeAdventure)
	w.SetTime(6000)
	w.StopTime()
	w.SetTickRange(0)
	w.StopThundering()
	w.StopRaining()
	lobby.New(w)

	registerCommands()

	srv.Listen()
	for srv.Accept(accept) {
		// Do nothing
	}
}

func accept(p *player.Player) {
	lobby.AddPlayer(p)
	user.Add(p)
	p.ShowCoordinates()
}

func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("spawn", text.Colourf("<redstone>Teleport to spawn.</redstone>"), []string{"hub"}, command.Spawn{}),
		cmd.New("role", text.Colourf("<redstone>Role management commands.</redstone>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}),
		cmd.New("duel", text.Colourf("<redstone>Duel other players or parties.</redstone>"), nil, command.DuelAccept{}, command.Duel{}),
		cmd.New("ban", text.Colourf("<redstone>Ban other players</redstone>"), nil, command.BanList{}, command.BanLiftOffline{}, command.BanInfoOffline{}, command.Ban{}, command.BanOffline{}, command.BanForm{}),
		cmd.New("kick", text.Colourf("<redstone>Kick other players</redstone>"), nil, command.Kick{}),
		cmd.New("mute", text.Colourf("<redstone>Mute other players</redstone>"), nil, command.MuteList{}, command.MuteLiftOffline{}, command.MuteInfoOffline{}, command.Mute{}, command.MuteOffline{}, command.MuteForm{}),
		cmd.New("rekit", text.Colourf("<redstone>re-apply your kit</redstone>"), nil, command.ReKit{}),
		cmd.New("pprof", text.Colourf("<redstone>You shouldn't have access to this</redstone>"), nil, command.Pprof{}),
		cmd.New("status", text.Colourf("<redstone>View technical stats of the server.</redstone>"), nil, command.Status{}),
		//cmd.New("replay", text.Colourf("<redstone>View replay of duels.</redstone>"), nil, command.ReplayRecent{}),
	} {
		cmd.Register(c)
	}
}

func readConfig() (moyai.Config, error) {
	c := moyai.DefaultConfig()
	g := gophig.NewGophig("./config", "toml", 0777)

	err := g.GetConf(&c)
	if os.IsNotExist(err) {
		err = g.SetConf(c)
	}
	return c, err
}
