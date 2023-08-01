package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/practice/moyai/data"
	"regexp"
	"strings"
)

type Handler struct {
	player.NopHandler
	p *player.Player
}

func NewHandler(p *player.Player) *Handler {
	h := &Handler{p: p}
	return h
}

var (
	// tlds is a list of top level domains used for checking for advertisements.
	tlds = [...]string{".me", ".club", "www.", ".com", ".net", ".gg", ".cc", ".net", ".co", ".co.uk", ".ddns", ".ddns.net", ".cf", ".live", ".ml", ".gov", "http://", "https://", ",club", "www,", ",com", ",cc", ",net", ",gg", ",co", ",couk", ",ddns", ",ddns.net", ",cf", ",live", ",ml", ",gov", ",http://", "https://", "gg/"}
	// emojis is a map between emojis and their unicode representation.
	emojis = strings.NewReplacer(
		":l:", "\uE107",
		":skull:", "\uE105",
		":fire:", "\uE108",
		":eyes:", "\uE109",
		":clown:", "\uE10A",
		":100:", "\uE10B",
		":heart:", "\uE10C",
	)
)

// formatRegex is a regex used to clean color formatting on a string.
var formatRegex = regexp.MustCompile(`§[\da-gk-or]`)

// HandleChat ...
func (h *Handler) HandleChat(ctx *event.Context, message *string) {
	ctx.Cancel()

	u, ok := data.LoadUser(h.p.Name())
	if !ok {
		return
	}

	*message = emojis.Replace(*message)
	r := u.Roles.Highest()

	_, _ = chat.Global.WriteString(r.Chat(h.p.Name(), *message))
}

// HandleFoodLoss ...
func (*Handler) HandleFoodLoss(ctx *event.Context, _ int, _ *int) {
	ctx.Cancel()
}

// HandleBlockBreak ...
func (*Handler) HandleBlockBreak(ctx *event.Context, _ cube.Pos, _ *[]item.Stack, _ *int) {
	ctx.Cancel()
}

// HandleBlockPlace ...
func (*Handler) HandleBlockPlace(ctx *event.Context, _ cube.Pos, _ world.Block) {
	ctx.Cancel()
}
