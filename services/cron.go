package services

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v7"

	"github.com/xhyonline/xdq/component"
	"github.com/xhyonline/xutil/logger"
)

var scanLock sync.Mutex

// ScanBucketForReady 扫描
func ScanBucketForReady() {
	scanLock.Lock()
	defer scanLock.Unlock()
	client := component.Instance.Redis.Kv
	// 获取所有的主题
	topics := GetTopics()
	var wg sync.WaitGroup
	now := time.Now().Unix()
	wg.Add(len(topics))
	for _, topic := range topics {
		go func(topic string) {
			defer wg.Done()
			bucket := GetScanBucketName(topic)
			cmd := client.ZRangeByScore(bucket, &redis.ZRangeBy{
				Min: "0",
				Max: fmt.Sprintf("%d", now),
			})
			if cmd.Err() != nil {
				logger.Errorf("遍历主题失败 %s", cmd.Err())
				return
			}
			if len(cmd.Val()) == 0 {
				return
			}
			ready := make([]interface{}, 0, len(cmd.Val()))
			for _, uid := range cmd.Val() {
				ready = append(ready, uid)
			}
			if err := client.ZRem(bucket, ready...).Err(); err != nil {
				logger.Errorf("移除 bucket 成员失败 %s", err)
				return
			}
			if err := client.LPush(GetReadyListName(topic), ready...).Err(); err != nil {
				logger.Errorf("写入准备区域失败 %s", err)
				return
			}
		}(topic)
	}
	wg.Wait()
}

var consumeLock sync.Mutex

// ConsumeReadyJob 消费
func ConsumeReadyJob() {
	consumeLock.Lock()
	defer consumeLock.Unlock()
	topics := GetTopics()
	var wg sync.WaitGroup
	client := component.Instance.Redis.Kv
	wg.Add(len(topics))
	for _, topic := range topics {
		go func(topic string) {
			defer wg.Done()
			cfg := GetTopicConfig(topic)
			if cfg == nil {
				return
			}
			var flag = false
			for {
				var n int64 = 1
				ready := GetReadyListName(topic)
				readyCmd := client.LRange(ready, 0, n-1)
				if readyCmd.Err() != nil {
					logger.Errorf("扫描 ready 队列失败 %s", readyCmd.Err())
					return
				}
				if len(readyCmd.Val()) == 0 {
					return
				}
				if len(readyCmd.Val()) < int(n) {
					flag = true
				}
				jobs := client.HMGet(topic, readyCmd.Val()...)
				if jobs.Err() != nil {
					logger.Errorf("获取 jobs 失败 %s", jobs.Err())
					return
				}
				messages := make([]*Message, 0)
				for _, job := range jobs.Val() {
					item := new(Message)
					if err := json.Unmarshal([]byte(job.(string)), item); err != nil {
						logger.Error("消费数据失败,元数据为:", job.(string))
						continue
					}
					messages = append(messages, item)
				}
				// 回到失败等着下一次
				if !CallBack(cfg.CallbackURL, topic, messages, time.Duration(cfg.Timeout)) {
					break
				}
				// TODO 删除业务数据
				client.LTrim(ready, n, -1)
				if flag {
					break
				}
			}
		}(topic)
	}
	wg.Wait()
}
