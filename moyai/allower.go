package moyai

import (
	"net"
	"strings"

	"github.com/hako/durafmt"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/rcrowley/go-bson"
	"github.com/sandertv/gophertunnel/minecraft/protocol/login"
	"github.com/unickorn/strutils"
	"golang.org/x/text/language"
)

// allower ensures that all players who join are whitelisted if whitelisting is enabled.
type Allower struct {
	whitelisted bool
}

// Allow ...
func (a *Allower) Allow(ip net.Addr, identity login.IdentityData, client login.ClientData) (string, bool) {
	users, err := data.LoadUsersCond(bson.M{
		"xuid": identity.XUID,
	})
	if err != nil {
		panic(err)
	}
	for _, u := range users {
		if !u.Punishments.Ban.Expired() {
			if u.Punishments.Ban.Permanent {
				description := lang.Translatef(language.English, "user.blacklist.description", strings.TrimSpace(u.Punishments.Ban.Reason))
				if u.XUID == identity.XUID {
					return strutils.CenterLine(lang.Translatef(language.English, "user.blacklist.header") + "\n" + description), false
				}
				return strutils.CenterLine(lang.Translatef(language.English, "user.blacklist.header.alt") + "\n" + description), false
			}
			description := lang.Translatef(language.English, "user.ban.description", strings.TrimSpace(u.Punishments.Ban.Reason), durafmt.ParseShort(u.Punishments.Ban.Remaining()))
			if u.XUID == identity.XUID {
				return strutils.CenterLine(lang.Translatef(language.English, "user.ban.header") + "\n" + description), false
			}
			return strutils.CenterLine(lang.Translatef(language.English, "user.ban.header.alt") + "\n" + description), false
		}
	}

	// if a.whitelisted {
	// 	u, ok := data.LoadUser(identity.DisplayName)
	// 	if !ok {
	// 		return strutils.CenterLine(lang.Translatef(language.English, "user.server.whitelist")), false
	// 	}
	// 	return strutils.CenterLine(lang.Translate("user.server.whitelist")), u.Whitelisted || u.Roles.Contains(role.Trial{}, role.Operator{})
	// }
	return "", true
}
