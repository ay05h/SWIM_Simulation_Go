package node

import (
	"clusterpulse/protocol"
	"math/rand"
	"sort"
	"time"
)

type Config struct {
	ID            string
	ProbeInterval time.Duration
	AckTimeout    time.Duration
	KnownNodes    []string
}

type Node struct {
	id            string
	peers         []string
	probeInterval time.Duration
	ackTimeout    time.Duration
	rng           *rand.Rand
	inbox         chan protocol.Message
}

func New(cfg Config) *Node {
	if cfg.ProbeInterval <= 0 {
		cfg.ProbeInterval = 750 * time.Millisecond
	}

	if cfg.AckTimeout <= 0 {
		cfg.AckTimeout = 350 * time.Millisecond
	}

	peers := make([]string, 0, len(cfg.KnownNodes))
	seen := make(map[string]struct{}, len(cfg.KnownNodes))

	for _, nodeID := range cfg.KnownNodes {
		if nodeID == "" || nodeID == cfg.ID {
			continue
		}

		if _, exists := seen[nodeID]; exists {
			continue
		}

		seen[nodeID] = struct{}{}
		peers = append(peers, nodeID)
	}

	sort.Strings(peers)

	bufferSize := len(peers)*4+8
	if bufferSize < 16 {
		bufferSize = 16
	}

	return &Node{
		id:            cfg.ID,
		peers:         peers,
		probeInterval: cfg.ProbeInterval,
		ackTimeout:    cfg.AckTimeout,
		rng: rand.New(
			rand.NewSource(
				time.Now().UnixNano() + int64(len(cfg.ID))*7919,
			),
		),
		inbox: make(chan protocol.Message, bufferSize),
	}
}

func (n *Node) ID() string {
	return n.id
}

func (n *Node) Inbox() chan protocol.Message {
	return n.inbox
}
