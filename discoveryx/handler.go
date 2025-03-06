package registryx

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/asmile1559/dyshop/utils/balancerx"
	"github.com/jiny3/gopkg/hookx"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// 启动多个服务实例并注册到 Etcd
func StartEtcdServices[T any](
	endpoints []string,
	services []any,
	prefix string,
	registerFunc func(grpc.ServiceRegistrar, T),
	serverFactory func(instanceID string, etcdService *EtcdService) T,
) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	go hookx.Exit(cancel)
	for _, raw := range services {
		serviceMap := raw.(map[string]any)
		id := serviceMap["id"].(string)
		address := serviceMap["address"].(string)

		wg.Add(1)
		go func(id, addr string) {
			defer wg.Done()
			startEtcdServiceInstance(ctx, endpoints, id, addr, prefix, registerFunc, serverFactory)
		}(id, address)
	}
	wg.Wait()
}

// 启动单个服务实例并注册到 Etcd
func startEtcdServiceInstance[T any](
	ctx context.Context,
	endpoints []string,
	instanceID, address, prefix string,
	registerFunc func(grpc.ServiceRegistrar, T),
	serverFactory func(instanceID string, etcdService *EtcdService) T,
) {
	client, err := NewEtcdClient(endpoints)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create etcd client")
	}
	defer client.Close()

	etcdService, err := NewEtcdService(client, instanceID, prefix, address, 10*time.Second)
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create Etcd service")
	}

	err = etcdService.Register()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to register service")
	}
	defer etcdService.DeRegister()

	_, port, err := net.SplitHostPort(address)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to split address %s", address)
	}
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		logrus.WithError(err).Fatalf("Failed to listen on port %s", port)
	}

	grpcServer := grpc.NewServer()
	serverInstance := serverFactory(instanceID, etcdService)
	registerFunc(grpcServer, serverInstance)

	logrus.WithFields(logrus.Fields{
		"instanceID": instanceID,
		"address":    address,
	}).Info("running...")
	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			logrus.WithError(err).WithField("instanceID", instanceID).Error("Failed to serve instance")
		}
	}()
	<-ctx.Done()
	logrus.WithField("instanceID", instanceID).Debug("Shutting down...")
}

// 从 Etcd 中发现服务
func DiscoverEtcdServices[T any](
	endpoints []string,
	prefix string,
	newClientFunc func(grpc.ClientConnInterface) T,
) (T, *grpc.ClientConn, error) {
	var zero T

	// 初始化 Etcd 客户端
	client, err := NewEtcdClient(endpoints)
	if err != nil {
		logrus.WithError(err).Error("Failed to create etcd client")
		return zero, nil, err
	}

	// 从 Etcd 中发现服务
	services, err := discoverService(client, prefix)
	if err != nil {
		logrus.WithError(err).Error("Failed to discover service")
		return zero, nil, err
	}
	if len(services) == 0 {
		logrus.WithField("prefix", prefix).Error("No services found for key")
		return zero, nil, fmt.Errorf("no services found for key")
	}

	// 初始化负载均衡策略

	// balancer := balancerx.NewRandomBalancer() // 随机策略

	// if err := balancerx.InitRoundRobinKey(client, "/config/hello-service/round_robin_index"); err != nil {
	// 	logrus.Fatalf("Failed to init round robin key: %v", err)
	// }
	// balancer := balancerx.NewRoundRobinBalancer(client, "/config/hello-service/round_robin_index") // 轮询策略

	balancer := balancerx.NewLeastConnBalancer(client, prefix) // 最小连接数策略

	// 使用负载均衡器选择服务地址
	serviceAddress := balancer.Select(services)
	logrus.WithField("serviceAddr", serviceAddress).Debug("Selected service via balancer")

	// 连接到 gRPC 服务
	conn, err := grpc.NewClient(serviceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.WithError(err).Error("Failed to connect to gRPC server")
		return zero, nil, err
	}

	// 创建 gRPC 客户端
	clientInstance := newClientFunc(conn)
	return clientInstance, conn, nil
}
