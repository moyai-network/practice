package form

import (
	"github.com/sandertv/gophertunnel/minecraft/text"
	"math/rand"
	"sort"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/form"
	"github.com/hako/durafmt"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/carrot/role"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/moyai-network/practice/moyai/user"
	"golang.org/x/exp/maps"
)

// ban is a form that allows a user to issue a ban.
type ban struct {
	// Reason is a dropdown that allows a user to select a ban reason.
	Reason form.Dropdown
	// OnlinePlayer is a dropdown that allows a user to select an online player.
	OnlinePlayer form.Dropdown
	// OfflinePlayer is an input field that allows a user to enter an offline player.
	OfflinePlayer form.Input
	// online is a list of online players' XUIDs indexed by their names.
	online map[string]string
}

// NewBan creates a new form to issue a ban.
func NewBan() form.Form {
	online := make(map[string]string)
	for _, u := range user.All() {
		online[u.Name()] = u.Name()
	}
	names := [...]string{"Steve Harvey", "Elon Musk", "Bill Gates", "Mark Zuckerberg", "Jeff Bezos", "Warren Buffet", "Larry Page", "Sergey Brin", "Larry Ellison", "Tim Cook", "Steve Ballmer", "Daniel Larson", "Steve"}
	list := maps.Keys(online)
	sort.Strings(list)
	return form.New(ban{
		Reason:        form.NewDropdown("Reason", []string{"Unfair Advantage", "Unfair Advantage in Ranked", "Interference", "Exploitation", "Permission Abuse", "Invalid Skin", "Evasion", "Advertising"}, 0),
		OnlinePlayer:  form.NewDropdown("Online Player", list, 0),
		OfflinePlayer: form.NewInput("Offline Player", "", names[rand.Intn(len(names)-1)]),
		online:        online,
	}, text.Colourf("<orange>» <black>Ban</black> «</orange>"))
}

// Submit ...
func (b ban) Submit(s form.Submitter) {
	p := s.(*player.Player)
	u, err := data.LoadUserOrCreate(p.Name())
	if err != nil {
		// User somehow left midway through the form.
		return
	}

	if !u.Roles.Contains(role.Mod{}, role.Operator{}) {
		// In case the user's role was removed while the form was open.
		return
	}
	var length time.Duration
	reason := b.Reason.Options[b.Reason.Value()]
	switch reason {
	case "Unfair Advantage":
		length = time.Hour * 24 * 30
	case "Unfair Advantage in Ranked":
		length = time.Hour * 24 * 90
	case "Interference":
		length = time.Hour * 12
	case "Exploitation":
		length = time.Hour * 24 * 9
	case "Permission Abuse":
		length = time.Hour * 24 * 30
	case "Invalid Skin":
		length = time.Hour * 24 * 3
	case "Evasion":
		length = time.Hour * 24 * 120
	case "Advertising":
		length = time.Hour * 24 * 6
	}

	punishment := carrot.Punishment{
		Staff:      p.Name(),
		Reason:     reason,
		Occurrence: time.Now(),
		Expiration: time.Now().Add(length),
	}
	var name string
	if offlineName := strings.TrimSpace(b.OfflinePlayer.Value()); offlineName != "" {
		if strings.EqualFold(offlineName, p.Name()) {
			p.Message("command.ban.self")
			return
		}
		t, err := data.LoadUserOrCreate(offlineName)
		if err != nil {
			p.Message("command.target.unknown")
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			p.Message("command.ban.operator")
			return
		}
		if !t.Punishments.Ban.Expired() {
			p.Message("command.ban.already")
			return
		}
		t.Punishments.Ban = punishment
		_ = data.SaveUser(t)

		name = t.DisplayName
	} else {
		t, err := data.LoadUserOrCreate(b.online[b.OnlinePlayer.Options[b.OnlinePlayer.Value()]])
		if err != nil {
			p.Message("command.target.unknown")
			return
		}
		if t.Roles.Contains(role.Operator{}) {
			p.Message("command.ban.operator`")
			return
		}

		tH, ok := user.Lookup(t.Name)
		t.Punishments.Ban = punishment

		if ok {
			tH.Disconnect(strings.Join([]string{
				lang.Translatef(tH.Locale(), "user.ban.header"),
				lang.Translatef(tH.Locale(), "user.ban.description", reason, durafmt.ParseShort(length)),
			}, "\n"))
		}
		name = t.DisplayName

		_ = data.SaveUser(t)
	}

	_ = data.SaveUser(u) // Save in case of a server crash or anything that may cause the data to not get saved.

	user.Alert(p, "staff.alert.ban", name, reason)
	user.Broadcast("command.ban.broadcast", p.Name(), name, reason)
	//webhook.SendPunishment(p.Name(), name, reason, "Ban")
	p.Message("command.ban.success", name, reason)
}
