package initialize

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"mxshop_srvs/goods_srv/global"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func GetEnvInfo(env string) bool{
	viper.AutomaticEnv()
	return viper.GetBool(env)
}

func InitConfig(){
	debug := GetEnvInfo("MXSHOP_DEBUG")
	configFilePrefix := "config"
	configFileName := fmt.Sprintf("goods_srv/%s-pro.yaml",configFilePrefix)
	if debug {
		configFileName = fmt.Sprintf("goods_srv/%s-debug.yaml",configFilePrefix)
	}

	v :=viper.New()
	v.SetConfigFile(configFileName)
	if err := v.ReadInConfig();err != nil{
		panic(err)
	}

	if err := v.Unmarshal(&global.NacosConfig);err != nil{
		panic(err)
	}
	zap.S().Infof("配置信息: %v",global.NacosConfig)

	sc := []constant.ServerConfig{
		{
			IpAddr: global.NacosConfig.Host,
			Port: global.NacosConfig.Port,
		},
	}

	cc := constant.ClientConfig{
		NamespaceId: global.NacosConfig.Namespace,
		TimeoutMs: 5000,
		NotLoadCacheAtStart: true,
		LogDir: "tmp/nacos/log",
		CacheDir: "tmp/nacos/cache",
		RotateTime: "1h",
		MaxAge: 3,
		LogLevel: "debug",
	}

	configClient,err := clients.CreateConfigClient(map[string]interface{}{
		"serverConfigs": sc,
		"clientConfig":cc,
	})
	if err != nil {
		panic(err)
	}

	content,err := configClient.GetConfig(vo.ConfigParam{
		DataId: global.NacosConfig.DataId,
		Group: global.NacosConfig.Group,
	})

	if err != nil {
		panic(err)
	}


	err = json.Unmarshal([]byte(content),&global.ServerConfig)

	if err != nil {
		zap.S().Fatalf("读取Nacos配置失败: %s",err.Error())
	}
	fmt.Println(global.ServerConfig)

}