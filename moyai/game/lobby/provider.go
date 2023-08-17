package lobby

import (
	"github.com/moyai-network/practice/moyai/game"
	"math"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/icza/abcsort"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game/duel"
	"github.com/moyai-network/practice/moyai/game/ffa"
	"github.com/moyai-network/practice/moyai/game/kit"
	"github.com/restartfu/roman"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"golang.org/x/exp/slices"
)

func init() {
	ffa.InitializeLobby(AddPlayer)
	duel.InitializeLobby(AddPlayer)
}

func New(w *world.World) {
	go world.NewLoader(32, w, world.NopViewer{}).Load(math.MaxInt)
	lobby = w
	go startLeaderBoards()
}

func formattedKillsLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><redstone>TOP %v</redstone></bold>\n", strings.ReplaceAll(strings.ToUpper("kills"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.Stats.Kills == b.Stats.Kills {
			return 0
		}
		if a.Stats.Kills > b.Stats.Kills {
			return -1
		}
		return 1
	})

	for i := 0; i < 10; i++ {
		if len(users) < i+1 {
			break
		}
		leader := users[i]
		name := leader.DisplayName
		if leader.Roles.Contains(role.Plus{}) {
			name = text.Colourf("<red>%s</red>", name)
		}

		position, _ := roman.Itor(i + 1)
		sb.WriteString(text.Colourf(
			"<grey>%v.</grey> <white>%v</white> <dark-grey>-</dark-grey> <grey>%v</grey>\n",
			position,
			name,
			leader.Stats.Kills,
		))
	}
	return sb.String()
}

func formattedEloLeaderboard(g game.Game) string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><redstone>TOP %v</redstone></bold>\n", strings.ReplaceAll(strings.ToUpper(g.Name()), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.GameElo(g) == b.GameElo(g) {
			return 0
		}
		if a.GameElo(g) > b.GameElo(g) {
			return -1
		}
		return 1
	})

	for i := 0; i < 10; i++ {
		if len(users) < i+1 {
			break
		}
		leader := users[i]
		name := leader.DisplayName
		if leader.Roles.Contains(role.Plus{}) {
			name = text.Colourf("<red>%s</red>", name)
		}

		position, _ := roman.Itor(i + 1)
		sb.WriteString(text.Colourf(
			"<grey>%v.</grey> <white>%v</white> <dark-grey>-</dark-grey> <grey>%v</grey>\n",
			position,
			name,
			leader.GameElo(g),
		))
	}
	return sb.String()
}

func startLeaderBoards() {
	var gamesIndex int
	games := game.Games()

	killsLeaderboard := entity.NewText(formattedKillsLeaderboard(), cube.Pos{3, 60, 59}.Vec3Middle())
	eloLeaderboard := entity.NewText(formattedEloLeaderboard(games[gamesIndex]), cube.Pos{-3, 60, 59}.Vec3Middle())

	lobby.AddEntity(killsLeaderboard)
	lobby.AddEntity(eloLeaderboard)

	t := time.NewTicker(time.Second * 4)

	for range t.C {
		gamesIndex++

		if gamesIndex >= len(games) {
			gamesIndex = 0
		}

		killsLeaderboard.SetNameTag(formattedKillsLeaderboard())
		eloLeaderboard.SetNameTag(formattedEloLeaderboard(games[gamesIndex]))
	}
}

var lobby *world.World

func Contains(p *player.Player) bool {
	return lobby == p.World()
}

func AddPlayer(p *player.Player) {
	if c, closeable := p.Handler().(interface{ Close() }); closeable {
		c.Close()
	}

	lobby.AddEntity(p)
	p.Teleport(lobby.Spawn().Vec3Middle())

	kit.Apply(kit.Lobby{}, p)
	h := newHandler(p)
	h.SendScoreBoard()
	p.Handle(h)
	p.Inventory().Handle(inventoryHandler{})
	p.Armour().Handle(inventoryHandler{})

	u, _ := data.LoadUser(p.Name())
	p.SetNameTag(u.Roles.Highest().Colour(p.Name()))
}
