package cluster

import (
	"clusterpulse/node"
	"fmt"
	"log"
	"sync"
	"time"
)

type Config struct {
	NodeCount      int
	ProbeInterval  time.Duration
	AckTimeout     time.Duration
	BaseLatency    time.Duration
	LatencyJitter  time.Duration
	DropRate       float64
}

type Cluster struct {
	mu      sync.Mutex
	wg      sync.WaitGroup
	log     *log.Logger
	ids     []string
	nodes   map[string]*node.Node
	started bool
	stopped bool
}

func New(cfg Config, logger *log.Logger) (*Cluster, error) {
	if cfg.NodeCount <= 0 {
		return nil, fmt.Errorf("NodeCount must be greater than 0")
	}

	ids := make([]string, cfg.NodeCount)

	for i := 0; i < cfg.NodeCount; i++ {
		ids[i] = fmt.Sprintf("node-%d", i+1)
	}

	nodes := make(map[string]*node.Node, len(ids))

	for _, nodeID := range ids {
		nodes[nodeID] = node.New(node.Config{
			ID:            nodeID,
			ProbeInterval: cfg.ProbeInterval,
			AckTimeout:    cfg.AckTimeout,
			KnownNodes:    ids,
		})
	}

	return &Cluster{
		ids:   ids,
		nodes: nodes,
		log:   logger,
	}, nil
}