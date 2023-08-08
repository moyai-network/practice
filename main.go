package main

import (
	"github.com/df-mc/dragonfly/server/cmd"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	_ "github.com/moyai-network/carrot/console"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/worlds"
	"github.com/moyai-network/practice/moyai"
	"github.com/moyai-network/practice/moyai/command"
	ent "github.com/moyai-network/practice/moyai/entity"
	"github.com/moyai-network/practice/moyai/game/lobby"
	"github.com/restartfu/gophig"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"github.com/sirupsen/logrus"
	"golang.org/x/text/language"
	"os"
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

	c.Name = text.Colourf("<bold><quartz>MOYAI</quartz></bold>") + "§8"

	/*ac := oomph.New(log, ":19132")
	ac.Listen(&c, c.Name, []minecraft.Protocol{}, false, false)
	go func() {
		for {
			p, err := ac.Accept()
			if err != nil {
				return
			}
			p.SetCombatMode(2)
			p.SetMovementMode(2)
			p.Handle(user.NewOomphHandler(p))
		}
	}()*/

	srv := c.New()

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
	srv.CloseOnProgramEnd()
	for srv.Accept(accept) {
		// Do nothing
	}
}

func accept(p *player.Player) {
	lobby.AddPlayer(p)
}

func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("spawn", text.Colourf("<orange>Teleport to spawn</orange>"), []string{"hub"}, command.Spawn{}),
		cmd.New("role", text.Colourf("<orange>Role management commands</orange>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}),
		//cmd.New("duel", carrot.GlyphFont("duel other players or parties", item.ColourOrange()), nil, command.Duel{}),
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
