package game

import (
	"regexp"
	"strings"

	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/game/kit"
)

func Games() []Game {
	return []Game{NoDebuff(), Boxing()}
}

func NoDebuff() Game {
	return Game{name: "NoDebuff", texture: "zeqa/textures/ui/gm/nodebuff.png", kit: kit.NoDebuff{}, ffa: true}
}

func Boxing() Game {
	return Game{name: "Boxing", texture: "zeqa/textures/ui/gm/boxing.png", kit: kit.Boxing{}}
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

func ByName(name string) Game {
	switch strings.ToLower(formatRegex.ReplaceAllString(name, "")) {
	case "nodebuff":
		return NoDebuff()
	case "boxing":
		return Boxing()
	}
	panic("should never happen: unknown game name: '" + name + "'")
}

type Game struct {
	name    string
	texture string

	ffa bool
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

func (g Game) Kit() carrot.Kit {
	return g.kit
}
