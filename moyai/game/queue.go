package game

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/sets"
	"sync"
)

var (
	queuesMu     sync.Mutex
	queues       = map[Game]sets.Set[*player.Player]{}
	rankedQueues = map[Game]sets.Set[*player.Player]{}
)

func init() {
	queuesMu.Lock()
	for _, g := range Games() {
		queues[g] = sets.New[*player.Player]()
		rankedQueues[g] = sets.New[*player.Player]()
	}
	queuesMu.Unlock()
}

func Queued(g Game, ranked bool) []*player.Player {
	queuesMu.Lock()
	defer queuesMu.Unlock()

	if ranked {
		return rankedQueues[g].Values()
	}
	return queues[g].Values()
}

func Queue(p *player.Player, g Game, ranked bool) {
	queuesMu.Lock()
	if ranked {
		rankedQueues[g].Add(p)
	} else {
		queues[g].Add(p)
	}
	queuesMu.Unlock()
}

func DeQueue(p *player.Player) {
	queuesMu.Lock()
	for _, q := range queues {
		if q.Contains(p) {
			q.Delete(p)

		}
	}
	for _, q := range rankedQueues {
		if q.Contains(p) {
			q.Delete(p)
		}
	}
	queuesMu.Unlock()
}
