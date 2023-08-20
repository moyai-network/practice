package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/game"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"strings"
)

// competitiveStats is a form that displays the competitive stats of a player.
type competitiveStats struct {
	// id is the xuid or name of the target player.
	id string
}

// NewCompetitiveStats creates a new competitive stats form to send to a player.
func NewCompetitiveStats(id string) form.Form {
	u, _ := data.LoadUser(id)
	displayName := u.DisplayName
	stats := u.Stats

	var games []string
	for _, g := range game.Games() {
		if !g.Duel() {
			continue
		}
		games = append(games, text.Colourf("<red>%s:</red> <white>%v</white>", g.Name(), u.GameElo(g)))
	}

	return form.NewMenu(competitiveStats{id: id}, text.Colourf("<dark-red>» <red>%v's Competitive Stats</red> «</dark-red>", displayName)).WithButtons(
		form.NewButton("View Casual Stats", ""),
	).WithBody(
		text.Colourf("<red>Wins:</red> <white>%v</white>", stats.RankedWins),
		text.Colourf("\n<red>Losses:</red> <white>%v</white>\n", stats.RankedLosses),
		text.Colourf("\n<dark-red>Elo</dark-red>\n"),
		text.Colourf("<red>Global:</red> <white>%v</white>\n", u.TotalElo()),
		strings.Join(games, "\n "),
	)
}

// Submit ...
func (c competitiveStats) Submit(s form.Submitter, _ form.Button) {
	p := s.(*player.Player)
	f := NewCasualStats(c.id)
	p.SendForm(f)
}
