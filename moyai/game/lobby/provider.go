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
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("kills"), "_", " ")))
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

func formattedDeathsLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("deaths"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.Stats.Deaths == b.Stats.Deaths {
			return 0
		}
		if a.Stats.Deaths > b.Stats.Deaths {
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
			leader.Stats.Deaths,
		))
	}
	return sb.String()
}

func formattedBestKSLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("Best KillStreak"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.Stats.BestKillStreak == b.Stats.BestKillStreak {
			return 0
		}
		if a.Stats.BestKillStreak > b.Stats.BestKillStreak {
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
			leader.Stats.BestKillStreak,
		))
	}
	return sb.String()
}

func formattedKSLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("KillStreak"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.Stats.KillStreak == b.Stats.KillStreak {
			return 0
		}
		if a.Stats.KillStreak > b.Stats.KillStreak {
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
			leader.Stats.KillStreak,
		))
	}
	return sb.String()
}

func formattedKDRLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("K/D Ratio"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.KDR() == b.KDR() {
			return 0
		}
		if a.KDR() > b.KDR() {
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
			"<grey>%v.</grey> <white>%v</white> <dark-grey>-</dark-grey> <grey>%.2f</grey>\n",
			position,
			name,
			leader.KDR(),
		))
	}
	return sb.String()
}

func formattedEloLeaderboard(g game.Game) string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper(g.Name()), "_", " ")))
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

func formattedOverallEloLeaderboard() string {
	sb := &strings.Builder{}
	sb.WriteString(text.Colourf("<bold><dark-red>TOP %v</dark-red></bold>\n", strings.ReplaceAll(strings.ToUpper("overall"), "_", " ")))
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.TotalElo() == b.TotalElo() {
			return 0
		}
		if a.TotalElo() > b.TotalElo() {
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
			leader.TotalElo(),
		))
	}
	return sb.String()
}

func formattedTotalEloLeaderboard() string {
	sb := &strings.Builder{}
	users := data.Users()

	sorter := abcsort.New("abcdefghijklmnopqrstuvwxyz123456789 ")
	sorter.Slice(users, func(i int) string {
		return users[i].Name
	})

	slices.SortFunc(users, func(a, b data.User) int {
		if a.TotalElo() == b.TotalElo() {
			return 0
		}
		if a.TotalElo() > b.TotalElo() {
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
			leader.TotalElo(),
		))
	}
	return sb.String()
}

func startLeaderBoards() {
	var gamesIndex, statsIndex int

	var games []game.Game
	for _, g := range game.Games() {
		if g.Duel() {
			games = append(games, g)
		}
	}

	statsLeaderboard := entity.NewText(formattedKillsLeaderboard(), cube.Pos{3, 60, 59}.Vec3Middle())
	eloLeaderboard := entity.NewText(formattedEloLeaderboard(games[gamesIndex]), cube.Pos{-3, 60, 59}.Vec3Middle())

	lobby.AddEntity(statsLeaderboard)
	lobby.AddEntity(eloLeaderboard)

	t := time.NewTicker(time.Second * 4)

	for range t.C {
		gamesIndex++
		statsIndex++

		if gamesIndex >= len(games) {
			eloLeaderboard.SetNameTag(formattedOverallEloLeaderboard())
			gamesIndex = 0
		}

		switch statsIndex {
		case 0:
			statsLeaderboard.SetNameTag(formattedKillsLeaderboard())
		case 1:
			statsLeaderboard.SetNameTag(formattedDeathsLeaderboard())
		case 2:
			statsLeaderboard.SetNameTag(formattedKSLeaderboard())
		case 3:
			statsLeaderboard.SetNameTag(formattedBestKSLeaderboard())
		case 4:
			statsLeaderboard.SetNameTag(formattedKDRLeaderboard())
		default:
			statsLeaderboard.SetNameTag(formattedKillsLeaderboard())
			statsIndex = 0
		}
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

	i, _ := p.Inventory().Item(2)
	_ = p.Inventory().SetItem(2, i.WithLore(strings.Split(formattedTotalEloLeaderboard(), "\n")...))

	h := newHandler(p)
	h.SendScoreBoard()
	p.Handle(h)
	p.Inventory().Handle(inventoryHandler{})
	p.Armour().Handle(inventoryHandler{})

	u, _ := data.LoadUser(p.Name())
	p.SetNameTag(u.Roles.Highest().Colour(p.Name()))
}
