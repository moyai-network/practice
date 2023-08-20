package form

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"time"
)

// casualStats is a form that displays the casual stats of a player.
type casualStats struct {
	// id is the xuid or name of the target player.
	id string
}

// NewCasualStats creates a new casual stats form to send to a player. An error will be returned if the offline user
// or if the offline user has their stats hidden.
func NewCasualStats(id string) form.Form {
	u, _ := data.LoadUser(id)

	displayName := u.DisplayName
	playtimeTotal := u.PlayTime.Round(time.Second)
	stats := u.Stats
	var playtimeSession time.Duration
	if p, ok := user.Lookup(u.Name); ok {
		playtimeSession = time.Since(p.Handler().(user.UserHandler).UserHandler().JoinTime()).Round(time.Second)
	}

	return form.NewMenu(casualStats{id: u.Name}, text.Colourf("<dark-red>» <red>%v's Casual Stats</red> «</dark-red>", displayName)).WithButtons(
		form.NewButton("View Competitive Stats", ""),
	).WithBody(
		text.Colourf(" <red>Playtime (Session):</red> <white>%s</white>\n", playtimeSession),
		text.Colourf("<red>Playtime (All Time):</red> <white>%s</white>\n", playtimeTotal+playtimeSession),
		text.Colourf("<red>Kills:</red> <white>%v</white>\n", stats.Kills),
		text.Colourf("<red>Killstreak:</red> <white>%v</white>\n", stats.KillStreak),
		text.Colourf("<red>Best Killstreak:</red> <white>%v</white>\n", stats.BestKillStreak),
		text.Colourf("<red>Deaths:</red> <white>%v</white>\n", stats.Deaths),
	)
}

// Submit ...
func (c casualStats) Submit(s form.Submitter, _ form.Button) {
	p, ok := s.(*player.Player)
	if !ok {
		return
	}
	p.SendForm(NewCompetitiveStats(c.id))
}
