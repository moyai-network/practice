package game

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/sets"
	"sync"
)

var (
	queuesMu sync.Mutex
	queues   = map[Game]sets.Set[*player.Player]{}
)

func init() {
	queuesMu.Lock()
	for _, g := range Games() {
		queues[g] = sets.New[*player.Player]()
	}
	queuesMu.Unlock()
}

func Queued(g Game) []*player.Player {
	queuesMu.Lock()
	q := queues[g]
	queuesMu.Unlock()
	return q.Values()
}

func Queue(p *player.Player, g Game) {
	queuesMu.Lock()
	queues[g].Add(p)
	queuesMu.Unlock()
}

func DeQueue(p *player.Player) {
	queuesMu.Lock()
	for _, q := range queues {
		if q.Contains(p) {
			q.Delete(p)
		}
	}
	queuesMu.Unlock()
}
