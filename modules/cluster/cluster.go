package cluster

import (
	"github.com/huminghe/infini-framework/core/cluster"
	pb "github.com/huminghe/infini-framework/core/cluster/pb"
	"github.com/huminghe/infini-framework/core/config"
	"github.com/huminghe/infini-framework/core/global"
	"github.com/huminghe/infini-framework/core/rpc"
	"github.com/huminghe/infini-framework/modules/cluster/demo/server"
	"github.com/huminghe/infini-framework/modules/cluster/discovery"
)

type ClusterModule struct {
}

func (module ClusterModule) Name() string {
	return "Cluster"
}

func (module ClusterModule) Setup(cfg *config.Config) {
	if !global.Env().SystemConfig.ClusterConfig.Enabled {
		return
	}

	rpc.Setup(&global.Env().SystemConfig.ClusterConfig.RPCConfig)
	cluster.New(&global.Env().SystemConfig.ClusterConfig)
}

func (module ClusterModule) Start() error {
	if !global.Env().SystemConfig.ClusterConfig.Enabled {
		return nil
	}

	server.Init()

	pb.RegisterDiscoveryServer(rpc.GetRPCServer(), &discovery.Discovery{})

	pb.RegisterRaftServer(rpc.GetRPCServer(), &discovery.RaftServer{})

	rpc.StartRPCServer()

	if err := cluster.Open(); err != nil {
		panic(err)
	}

	return nil
}

func (module ClusterModule) Stop() error {
	if !global.Env().SystemConfig.ClusterConfig.Enabled {
		return nil
	}

	cluster.SnapshotClusterState()
	return nil
}
