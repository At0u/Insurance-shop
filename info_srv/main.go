package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"mxshop_srvs/info_srv/global"
	"mxshop_srvs/info_srv/model"
	"mxshop_srvs/info_srv/utils"
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

	"mxshop_srvs/info_srv/handler"
	"mxshop_srvs/info_srv/initialize"
	"mxshop_srvs/info_srv/proto"
)

func main() {
	IP := flag.String("ip","0.0.0.0","ip地址")
	Port := flag.Int("port",50015,"端口号")

	initialize.InitLogger()
	initialize.InitConfig()
	initialize.InitDB()
	zap.S().Info(global.ServerConfig)

	flag.Parse()
	zap.S().Info("ip: ", *IP)
	if *Port == 0{
		*Port,_ = utils.GetFreePort()
	}

	zap.S().Info("port: ", *Port)


	server := grpc.NewServer()
	proto.RegisterInfoServer(server,&handler.InfoServer{})
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

	c, _ := rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{"192.168.1.107:9876"}),
	)

	if err := c.Subscribe("orderInfo", consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {

		for i := range msgs {
			order := &model.Order{}
			fmt.Printf("获取到值： %v \n", msgs[i])
			_ = json.Unmarshal(msgs[i].Body,order)
			updateErr := handler.UpdateInfo(order)
			if updateErr!=nil{
				return consumer.ConsumeRetryLater, nil
			}
		}
		return consumer.ConsumeSuccess, nil
	}); err != nil {
		zap.S().Infof(err.Error() + "从mq获取消息失败")
	}

	_ = c.Start()




	//接收终止信号
	quit := make(chan os.Signal)
	signal.Notify(quit,syscall.SIGINT,syscall.SIGTERM)
	<-quit
	_ = c.Shutdown()
	if err = client.Agent().ServiceDeregister(serviceID);err != nil{
		zap.S().Info("注销失败")
	}
	zap.S().Info("注销成功")

}

