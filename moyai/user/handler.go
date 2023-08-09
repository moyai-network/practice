package user

import (
	"regexp"
	"strings"
	"time"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/entity"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/data"
)

type Handler struct {
	player.NopHandler
	p *player.Player

	chatCoolDown carrot.CoolDown
}

func NewHandler(p *player.Player) *Handler {
	h := &Handler{p: p}
	return h
}

type UserHandler interface {
	UserHandler() *Handler
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

	if *message == "die" {
		h.p.Hurt(20, entity.VoidDamageSource{})
	}

	u, ok := data.LoadUser(h.p.Name())
	if !ok {
		return
	}

	*message = emojis.Replace(*message)
	r := u.Roles.Highest()
	msg := r.Chat(h.p.Name(), *message)

	if !u.Punishments.Mute.Expired() {
		h.p.Message(lang.Translatef(h.p.Locale(), "user.message.mute"))
		return
	}

	if h.chatCoolDown.Active() {
		h.p.Message(msg)
		return
	}
	h.chatCoolDown.Set(time.Second)

	_, _ = chat.Global.WriteString(msg)
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
