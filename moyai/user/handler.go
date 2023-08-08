package user

import (
	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/player"
	"github.com/df-mc/dragonfly/server/player/chat"
	"github.com/df-mc/dragonfly/server/player/scoreboard"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/moyai-network/carrot"
	"github.com/moyai-network/carrot/lang"
	"github.com/moyai-network/practice/moyai/data"
	"github.com/sandertv/gophertunnel/minecraft/text"
	"regexp"
	"strings"
	"time"
)

type Handler struct {
	player.NopHandler
	p *player.Player

	chatCoolDown carrot.CoolDown
}

func NewHandler(p *player.Player) *Handler {
	h := &Handler{p: p}
	h.SendScoreBoard()
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

	u, ok := data.LoadUser(h.p.Name())
	if !ok {
		return
	}

	*message = emojis.Replace(*message)
	r := u.Roles.Highest()
	msg := r.Chat(h.p.Name(), *message)

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

func (h *Handler) SendScoreBoard() {
	l := h.p.Locale()
	u, _ := data.LoadUser(h.p.Name())

	var kdr float64
	if u.Stats.Deaths > 0 {
		kdr = float64(u.Stats.Kills / u.Stats.Deaths)
	} else {
		kdr = float64(u.Stats.Kills)
	}

	sb := scoreboard.New(carrot.GlyphFont("PRACTICE"))
	sb.RemovePadding()
	_, _ = sb.WriteString("§r\uE000")

	_, _ = sb.WriteString("\uE142\uE143\uE144\uE143\uE142")
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>K<grey>:</grey> <black>%d</black> D<grey>:</grey> <black>%d</black>", u.Stats.Kills, u.Stats.Deaths))
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>KDR<grey>:</grey> <black>%.2f</black>", kdr))
	_, _ = sb.WriteString(text.Colourf("<black>\uE141 </black>KS<grey>:</grey> <black>%d</black>", u.Stats.KillStreak))

	_, _ = sb.WriteString("\uE000")
	for i, li := range sb.Lines() {
		if !strings.Contains(li, "\uE000") {
			sb.Set(i, "  "+li)
		}
	}
	_, _ = sb.WriteString(lang.Translatef(l, "scoreboard.footer"))
	h.p.SendScoreboard(sb)
}
