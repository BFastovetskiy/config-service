package cluster

import (
	"config-service/core"
	"sync"

	"github.com/shaj13/raft"
	"google.golang.org/grpc"
)

type stateMachine struct {
	app core.IApplication
	log core.ILogger
	mu  sync.Mutex
	kv  map[string]string
}

type Cluster struct {
	app                core.IApplication
	log                core.ILogger
	raftServer         *grpc.Server
	opts               []raft.Option
	startOpts          []raft.StartOption
	node               *raft.Node
	fsm                *stateMachine
	raftClusterAddress string
}

type entry struct {
	Key   string
	Value string
}
