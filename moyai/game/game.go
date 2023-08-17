package game

import (
	"regexp"
	"strings"

	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/game/kit"
)

func Games() []Game {
	return []Game{NoDebuff(), Boxing(), Fist(), Gapple()}
}

func NoDebuff() Game {
	return Game{name: "NoDebuff", texture: "zeqa/textures/ui/gm/nodebuff.png", kit: kit.NoDebuff{}, ffa: true, duel: true}
}

func Fist() Game {
	return Game{name: "Fist", texture: "zeqa/textures/ui/gm/fist.png", kit: kit.Fist{}, ffa: true}
}

func Gapple() Game {
	return Game{name: "Gapple", texture: "zeqa/textures/ui/gm/gapple.png", kit: kit.Gapple{}, ffa: true, duel: true}
}

func Boxing() Game {
	return Game{name: "Boxing", texture: "zeqa/textures/ui/gm/boxing.png", kit: kit.Boxing{}, duel: true}
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

func ByName(name string) Game {
	switch strings.ToLower(formatRegex.ReplaceAllString(name, "")) {
	case "nodebuff":
		return NoDebuff()
	case "boxing":
		return Boxing()
	case "fist":
		return Fist()
	case "gapple":
		return Gapple()
	}
	panic("should never happen: unknown game name: '" + name + "'")
}

type Game struct {
	name    string
	texture string

	ffa  bool
	duel bool

	kit carrot.Kit
}

func (g Game) Name() string {
	return g.name
}

func (g Game) Texture() string {
	return g.texture
}

func (g Game) FFA() bool {
	return g.ffa
}
func (g Game) Duel() bool {
	return g.duel
}

func (g Game) Kit() carrot.Kit {
	return g.kit
}
