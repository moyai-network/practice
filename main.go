package main

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/moyai-network/carrot/tebex"
	"github.com/moyai-network/practice/moyai/data"

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
	"github.com/moyai-network/practice/moyai/user"
	"github.com/restartfu/gophig"
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

	c.Name = text.Colourf("<bold><dark-red>MOYAI</dark-red></bold>") + "§8"
	c.JoinMessage = "<green>[+] %s</green>"
	c.QuitMessage = "<red>[-] %s</red>"

	// ac := oomph.New(log, ":19132")
	// ac.Listen(&c, c.Name, []minecraft.Protocol{}, true, true)
	// go func() {
	// 	for {
	// 		p, err := ac.Accept()
	// 		if err != nil {
	// 			return
	// 		}
	// 		logrus.Info("LOL OK")
	// 		p.Handle(user.NewOomphHandler(p))
	// 	}
	// }()

	srv := c.New()

	ch := make(chan os.Signal, 2)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-ch
		if err := srv.Close(); err != nil {
			log.Errorf("close server: %v", err)
		}
		if err := data.Close(); err != nil {
			log.Errorf("close data: %v", err)
		}
	}()

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

	store := loadStore(config.Moyai.Tebex, log)

	srv.Listen()
	for srv.Accept(acceptFunc(store, log)) {
		// Do nothing
	}
}

// acceptFunc returns a function for handling players joining.
func acceptFunc(store *tebex.Client, log *logrus.Logger) func(p *player.Player) {
	return func(p *player.Player) {
		user.Add(p)
		lobby.AddPlayer(p)
		p.SetGameMode(world.GameModeCreative)

		p.Message(text.Colourf("<green>Make sure to join our discord server at discord.gg/moyai!</green>"))
		store.ExecuteCommands(p)
	}
}

// loadStore initializes the Tebex store connection.
func loadStore(key string, log *logrus.Logger) *tebex.Client {
	store := tebex.NewClient(log, time.Second*5, key)
	name, domain, err := store.Information()
	if err != nil {
		log.Fatalf("tebex: %v", err)
		return nil
	}
	log.Infof("Connected to Tebex under %v (%v).", name, domain)
	return store
}

func registerCommands() {
	for _, c := range []cmd.Command{
		cmd.New("alias", text.Colourf("<dark-red>Get aliases of a player.</dark-red>"), nil, command.AliasOnline{}, command.AliasOffline{}),
		cmd.New("spawn", text.Colourf("<dark-red>Teleport to spawn.</dark-red>"), []string{"hub"}, command.Spawn{}),
		cmd.New("role", text.Colourf("<dark-red>Role management commands.</dark-red>"), nil, command.RoleAdd{}, command.RoleRemove{}, command.RoleAddOffline{}, command.RoleRemoveOffline{}),
		cmd.New("duel", text.Colourf("<dark-red>Duel other players or parties.</dark-red>"), nil, command.DuelAccept{}, command.Duel{}),
		cmd.New("ban", text.Colourf("<dark-red>Ban other players</dark-red>"), nil, command.BanList{}, command.BanLiftOffline{}, command.BanInfoOffline{}, command.Ban{}, command.BanOffline{}, command.BanForm{}),
		cmd.New("kick", text.Colourf("<dark-red>Kick other players</dark-red>"), nil, command.Kick{}),
		cmd.New("mute", text.Colourf("<dark-red>Mute other players</dark-red>"), nil, command.MuteList{}, command.MuteLiftOffline{}, command.MuteInfoOffline{}, command.Mute{}, command.MuteOffline{}, command.MuteForm{}),
		cmd.New("rekit", text.Colourf("<dark-red>re-apply your kit</dark-red>"), nil, command.ReKit{}),
		cmd.New("pprof", text.Colourf("<dark-red>You shouldn't have access to this</dark-red>"), nil, command.Pprof{}),
		cmd.New("pinfo", text.Colourf("<dark-red>You shouldn't have access to this</dark-red>"), nil, command.PlayerInfo{}, command.PlayerInfoOffline{}),
		cmd.New("status", text.Colourf("<dark-red>View technical stats of the server.</dark-red>"), nil, command.NewStatus(time.Now())),
		cmd.New("settings", text.Colourf("<dark-red>Manage your settings.</dark-red>"), []string{"parameters"}, command.Settings{}),
		cmd.New("stats", text.Colourf("<dark-red>See your or other people's stats.</dark-red>"), []string{"statistics"}, command.StatsOffline{}, command.Stats{}),
		cmd.New("whisper", text.Colourf("<dark-red>Send a private message to a player.</dark-red>"), []string{"w", "msg", "tell"}, command.Whisper{}),
		cmd.New("reply", text.Colourf("<dark-red>Reply to your last messenger.</dark-red>"), []string{"r"}, command.Reply{}),
		//cmd.New("replay", text.Colourf("<dark-red>View replay of duels.</dark-red>"), nil, command.ReplayRecent{}),
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
