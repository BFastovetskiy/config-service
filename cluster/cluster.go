package cluster

import (
	"config-service/about"
	"config-service/core"
	"config-service/utils"
	"fmt"
	"os"
	"path/filepath"

	"github.com/shaj13/raft"
	"github.com/shaj13/raft/transport"
	"github.com/shaj13/raft/transport/raftgrpc"
	"google.golang.org/grpc"
)

func InitCluster(app core.IApplication) *Cluster {
	c := &Cluster{
		app: app,
		log: app.GetLogger(),
	}

	c.raftClusterAddress = fmt.Sprintf(":%d", app.GetConfiguration().ClusterPort)
	c.startOpts = append(c.startOpts, raft.WithAddress(c.raftClusterAddress))
	raftClusterStateDirectory := filepath.Join(c.app.GetWorkDir(), about.Cluster_State_Directory)
	if !utils.ExistFileOrDir(raftClusterStateDirectory) {
		os.Mkdir(raftClusterStateDirectory, os.ModeDir)
	}
	c.opts = append(c.opts, raft.WithStateDIR(raftClusterStateDirectory))
	c.fsm = c.newstateMachine(c.app)
	c.node = raft.NewNode(c.fsm, transport.GRPC, c.opts...)
	c.raftServer = grpc.NewServer()
	raftgrpc.RegisterHandler(c.raftServer, c.node.Handler())
	return c
}

func (c Cluster) newstateMachine(app core.IApplication) *stateMachine {
	return &stateMachine{
		app: app,
		log: app.GetLogger(),
		kv:  make(map[string]string),
	}
}
