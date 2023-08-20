package game

import (
	"github.com/df-mc/dragonfly/server/player"
	"github.com/moyai-network/carrot/sets"
	"sync"
)

var (
	queuesMu          sync.Mutex
	casualQueues      = map[Game]sets.Set[*player.Player]{}
	competitiveQueues = map[Game]sets.Set[*player.Player]{}
)

func init() {
	queuesMu.Lock()
	for _, g := range Games() {
		casualQueues[g] = sets.New[*player.Player]()
		competitiveQueues[g] = sets.New[*player.Player]()
	}
	queuesMu.Unlock()
}

func Queued(g Game, competitive bool) []*player.Player {
	queuesMu.Lock()
	defer queuesMu.Unlock()

	if competitive {
		return competitiveQueues[g].Values()
	}
	return casualQueues[g].Values()
}

func Queue(p *player.Player, g Game, competitive bool) {
	queuesMu.Lock()
	if competitive {
		competitiveQueues[g].Add(p)
	} else {
		casualQueues[g].Add(p)
	}
	queuesMu.Unlock()
}

func DeQueue(p *player.Player) {
	queuesMu.Lock()
	for _, q := range casualQueues {
		if q.Contains(p) {
			q.Delete(p)

		}
	}
	for _, q := range competitiveQueues {
		if q.Contains(p) {
			q.Delete(p)
		}
	}
	queuesMu.Unlock()
}
