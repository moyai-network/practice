package game

import (
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/practice/moyai/game/kit"
	"regexp"
	"strings"
)

func Games() []Game {
	return []Game{NoDebuff()}
}

func NoDebuff() Game {
	return Game{name: "NoDebuff", texture: "textures/items/potion_bottle_splash_heal", kit: kit.NoDebuff{}}
}

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

func ByName(name string) Game {
	switch strings.ToLower(formatRegex.ReplaceAllString(name, "")) {
	case "nodebuff":
		return NoDebuff()
	}
	panic("should never happen: unknown game name: '" + name + "'")
}

type Game struct {
	name    string
	texture string
	kit     carrot.Kit
}

func (g Game) Name() string {
	return g.name
}

func (g Game) Texture() string {
	return g.texture
}

func (g Game) Kit() carrot.Kit {
	return g.kit
}
