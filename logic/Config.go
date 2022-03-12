package logic

import (
	"encoding/json"
	"io/ioutil"
)

type GatewayConfig struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// 程序配置
type Config struct {
	ServicePort                int             `json:"servicePort"`                // "HTTP/1服务端口": "接收业务方调用" "servicePort": 7799
	ServiceReadTimeout         int             `json:"serviceReadTimeout"`         // "接口读超时": "单位毫秒" "serviceReadTimeout": 2000
	ServiceWriteTimeout        int             `json:"serviceWriteTimeout"`        // "接口写超时": "单位毫秒" "serviceWriteTimeout": 2000
	GatewayList                []GatewayConfig `json:"gatewayList"`                // "网关列表": "推送将分发给所有网关"  {"hostname": "localhost","port": 7788}
	GatewayMaxConnection       int             `json:"gatewayMaxConnection"`       // "每个网关的最多并发连接数": "建议与gateway的CPU核数相等, 提升内部通讯吞吐"  "gatewayMaxConnection": 32
	GatewayTimeout             int             `json:"gatewayTimeout"`             // "网关单个请求的超时时间": "单位是毫秒"  "gatewayTimeout": 3000,
	GatewayIdleTimeout         int             `json:"gatewayIdleTimeout"`         // "网关连接的空闲关闭时间": "单位是秒" "gatewayIdleTimeout": 60
	GatewayDispatchWorkerCount int             `json:"gatewayDispatchWorkerCount"` // "向各个网关分发消息的协程数量": "CPU密集型, 与CPU个数相当即可" "gatewayDispatchWorkerCount": 32
	GatewayDispatchChannelSize int             `json:"gatewayDispatchChannelSize"` // "待分发消息队列长度": "分发本身很快, 队列不需要太大" "gatewayDispatchChannelSize": 100000
	GatewayMaxPendingCount     int             `json:"gatewayMaxPendingCount"`     // "每个网关的最大拥塞推送数": "当消息拥塞时, 后续发往该网关的消息将被丢弃" "gatewayMaxPendingCount": 200000
	GatewayPushRetry           int             `json:"gatewayPushRetry"`           // "每条推送的最大重试次数": "超过重试次数后, 消息将被丢弃" "gatewayPushRetry": 3
}

var (
	G_config *Config
)

func InitConfig(filename string) (err error) {
	var (
		content []byte
		conf    Config
	)

	if content, err = ioutil.ReadFile(filename); err != nil {
		return
	}

	if err = json.Unmarshal(content, &conf); err != nil {
		return
	}

	G_config = &conf
	return
}
