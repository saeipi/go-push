package gateway

import (
	"encoding/json"
	"io/ioutil"
)

// 程序配置
type Config struct {
	WsPort               int    `json:"wsPort"`               // "websocket监听端口": "建议nginx做代理转发" "wsPort": 7777
	WsReadTimeout        int    `json:"wsReadTimeout"`        // "websocket HTTP握手读超时": "单位毫秒" "wsReadTimeout": 2000
	WsWriteTimeout       int    `json:"wsWriteTimeout"`       // "websocket HTTP握手写超时": "单位毫秒" "wsWriteTimeout": 2000
	WsInChannelSize      int    `json:"wsInChannelSize"`      // "websocket读队列长度": "一般不需要修改" "wsInChannelSize": 1000
	WsOutChannelSize     int    `json:"wsOutChannelSize"`     // "WebSocket写队列长度": "一般不需要修改" "wsOutChannelSize": 1000
	WsHeartbeatInterval  int    `json:"wsHeartbeatInterval"`  // "WebSocket心跳检查间隔": "单位秒, 超过时间没有收到心跳, 服务端将主动断开链接" "wsHeartbeatInterval": 60
	MaxMergerDelay       int    `json:"maxMergerDelay"`       // "合并推送的最大延迟时间": "单位毫秒, 在抵达maxPushBatchSize之前超时则发送" "maxMergerDelay": 1000
	MaxMergerBatchSize   int    `json:"maxMergerBatchSize"`   // "合并最多消息条数": "消息推送频次越高, 应该使用更大的合并批次, 得到更高的吞吐收益" "maxMergerBatchSize": 100
	MergerWorkerCount    int    `json:"mergerWorkerCount"`    // "消息合并协程的数量": "消息合并与json编码耗费CPU, 注意一个房间的消息只会由同一个协程处理." "MergerWorkerCount": 32
	MergerChannelSize    int    `json:"mergerChannelSize"`    // "消息合并队列的容量": "每个房间消息合并线程有一个队列, 推送量超过队列将被丢弃" "mergerChannelSize": 1000
	ServicePort          int    `json:"servicePort"`          // "内部通讯HTTP2端口": "严禁该端口暴露到外网" "servicePort": 7788
	ServiceReadTimeout   int    `json:"serviceReadTimeout"`   // "内部通讯HTTP2读超时": "单位毫秒" "serviceReadTimeout": 2000
	ServiceWriteTimeout  int    `json:"serviceWriteTimeout"`  // "内部通讯HTTP2写超时": "单位毫秒" "serviceWriteTimeout": 2000
	ServerPem            string `json:"serverPem"`            // "内部通讯HTTP2 TLS证书": "私有证书,默认有效期10年",
	ServerKey            string `json:"serverKey"`            // "内部通讯HTTP2 TLS密钥": "与证书配对"
	BucketCount          int    `json:"bucketCount"`          // "连接分桶的数量": "桶越多, 推送的锁粒度越小, 推送并发度越高" "bucketCount": 512
	BucketWorkerCount    int    `json:"bucketWorkerCount"`    // "每个桶的处理协程数量": "影响同一时刻可以有多少个不同消息被分发出去" "bucketWorkerCount": 32
	MaxJoinRoom          int    `json:"maxJoinRoom"`          // "每个连接最多加入房间数量": "目前房间ID没有校验, 所以先做简单的数量控制" "maxJoinRoom": 5
	DispatchChannelSize  int    `json:"dispatchChannelSize"`  // "待分发队列的长度": "分发队列缓冲所有待推送的消息, 等待被分发到Bucket"  "dispatchChannelSize": 100000
	DispatchWorkerCount  int    `json:"dispatchWorkerCount"`  // "分发协程的数量": "分发协程用于将待推送消息扇出给各个Bucket" "dispatchWorkerCount": 32
	BucketJobChannelSize int    `json:"bucketJobChannelSize"` // "Bucket工作队列长度": "每个Bucket的分发任务放在一个独立队列中"  "bucketJobChannelSize": 1000
	BucketJobWorkerCount int    `json:"bucketJobWorkerCount"` // "Bucket发送协程的数量": "每个Bucket有多个协程并发的推送消息" "bucketJobWorkerCount": 32
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
