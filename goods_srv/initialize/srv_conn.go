package initialize

import (
	"fmt"
	_ "github.com/mbobakov/grpc-consul-resolver"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"mxshop_srvs/goods_srv/global"
	"mxshop_srvs/goods_srv/proto"
)

func InitSrvConn(){
	consulInfo := global.ServerConfig.ConsulInfo
	infoConn, err := grpc.Dial(
		fmt.Sprintf("consul://%s:%d/%s?wait=14s", consulInfo.Host, consulInfo.Port, "info-srv"),
		grpc.WithInsecure(),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
		//grpc.WithUnaryInterceptor(otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer())),
	)
	if err != nil {
		zap.S().Fatal("[InitSrvConn] 连接 【信息服务失败】")
	}
	zap.S().Infof("[InitSrvConn] 连接 【信息服务】成功")
	global.InfoSrvClient = proto.NewInfoClient(infoConn)
}
