package main

import (
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"mxshop_srvs/order_srv/global"
	"mxshop_srvs/order_srv/utils"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/hashicorp/consul/api"
	"github.com/satori/go.uuid"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"mxshop_srvs/order_srv/handler"
	"mxshop_srvs/order_srv/initialize"
	"mxshop_srvs/order_srv/proto"
)

func main() {
	IP := flag.String("ip","0.0.0.0","ip地址")
	Port := flag.Int("port",50020,"端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	initialize.InitSrvConn()

	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0{
		*Port,_ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)


	server := grpc.NewServer()
	proto.RegisterOrderServer(server,&handler.OrderServer{})
	lis,err := net.Listen("tcp",fmt.Sprintf("%s:%d",*IP, *Port))
	if err != nil {
		panic("failed to listen" + err.Error())
	}

	grpc_health_v1.RegisterHealthServer(server,health.NewServer())

	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d",global.ServerConfig.ConsulInfo.Host,global.ServerConfig.ConsulInfo.Port)

	client,err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	check := &api.AgentServiceCheck{
		GRPC: fmt.Sprintf("192.168.1.106:%d",*Port),
		Timeout: "5s",
		Interval: "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	registration := new(api.AgentServiceRegistration)
	registration.Name = global.ServerConfig.Name
	serviceID := fmt.Sprintf("%s",uuid.NewV4())
	registration.ID = serviceID
	registration.Port = *Port
	registration.Tags = global.ServerConfig.Tags
	registration.Address = "192.168.1.106"
	registration.Check = check

	err = client.Agent().ServiceRegister(registration)
	if err != nil{
		panic(err)
	}

	go func(){
		err = server.Serve(lis)
		if err != nil {
			panic("failed to start grpc:" + err.Error())
		}
	}()

	//监听订单超时topic
	global.Producer, err = rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.1.107:9876"}),
		//producer.WithRetry(1),
	)

	err = global.Producer.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}

	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit

	_ = global.Producer.Shutdown()

	if err = client.Agent().ServiceDeregister(serviceID);err != nil{
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}

