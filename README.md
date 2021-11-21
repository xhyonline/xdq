# xdq 基于 Redis 与 Go 实现的延迟队列

## 特性

1. 开箱即用
2. 秒级延迟
3. 消费数据时使用 HTTP 请求回调通知

## 使用说明

如果您在 Linux 生产环境,请自行配置 redis 信息路径如下

```
/usr/local/go-micro/common/redis.toml
```

配置内容

```
[redis]
host = ""
port = 6379
password = ""
# 延迟队列默认 0 号 DB
db = 0
```

### 一、推送延迟数据

请求地址:`http://ip:port/push`

请求方式:`POST`

Content-Type:`application/json`

参数说明:

| 参数         | 类型   | 说明                             |
| ------------ | ------ | -------------------------------- |
| callback_url | string | HTTP 消息回调通知地址            |
| timeout      | int    | HTTP 消息回调时设置的超时时间    |
| topic        | string | 主题名                           |
| message      | array  | 多条消息,子集如下                            |
| id           | string | 业务的UID,请自行在业务处保持唯一 |
| content      | string | 业务内容                         |
| time         | int    | 10位时间戳                       |

请求示例如下:

```
{
    "callback_url": "http://localhost:8082/callback",
    "timeout": 3,
    "topic": "my_topic",
    "message": [
        {
            "id": "1",
            "content": "业务内容3",
            "time": 1637466645
        },
        {
            "id": "2",
            "content": "业务内容4",
            "time": 1637466645
        }
    ]
}
```

### 二、回调事件通知

通知方法:`HTTP` 请求`POST`

请求体示例如下:

```
{
	"topic": "my_topic",
	"message": [{
		"id": "2",
		"content": "业务内容4",
		"time": 1637466645
	}, {
		"id": "1",
		"content": "业务内容3",
		"time": 1637466645
	}],
	"count": 2
}
```

参数说明:

| 参数    | 类型   | 说明                             |
| ------- | ------ | -------------------------------- |
| topic   | string | 主题名                           |
| count   | int    | 消息条数                         |
| message | array  | 消息,消息具体细节请看下面三个参数                             |
| id      | string | 业务的UID,请自行在业务处保持唯一 |
| content | string | 业务内容                         |
| time    | int    | 10位时间戳                       |

请注意,如果您在一个主题下,一次推送 1W 条同时执行的数据,为了防止一次性回调的数据量太大, xdq 每次将只回调 5K 条数据,另 5K 条数据,将紧接着上次请求后 1s 发送。

如果您想要更高的实时性,您可以对单个主题进行主题拆分,例如`my_topic1`、`my_topic2`。多个主题发送数据将采取并行策略。

如果回调通知失败,xdq 将采取重复通知的策略,直至消息发送成功

### 三、获取当前所有主题

请求地址:`http://ip:port/get/topics`

请求方式:`GET`

响应结果:

```
    {
    "code": 0,
    "data": [
        "my_topic"  // 主题名
    ],
    "message": ""
}
```

### 四、删除主题

注:该操作较为危险,他将删除该主题下所有正在消费的数据,因此请确保您的主题下没有业务数据时再删除。

请求地址:`http://ip:port/delete/:name`

请求方式:`GET`

示例:`http://ip:port/delete/my_topic`

响应结果

```
{
    "code": 0,
    "data": null,
    "message": ""
}
```

### 五、获取某主题下正在等待被消费的数据

请求地址:`http://ip:port/get/topic/:name`

请求方式:`GET`

示例:`http://ip:port/get/topic/:name/my_topic`

响应示例(消息 wait_message 将根据 `time` 字段升序排列)

```
{
    "code": 0,
    "data": {
        "callback_url": "http://localhost:8082/callback",
        "timeout": 3,
        "wait_message": [
            {
                "id": "1",
                "content": "业务内容1",
                "time": 1637466642
            },
            {
                "id": "2",
                "content": "业务内容2",
                "time": 1637466645
            },
        ]
    },
    "message": ""
}
```

请注意,该操作将会返回该主题下所有等待被消费的数据,当数据量大的时候,请谨慎操作



