package services

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/xhyonline/xutil/logger"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/xhyonline/xdq/component"
)

const TopicSet = "topic-set-by-xhyonline"

// PushParams
type PushParams struct {
	TopicConfig
	Topic   string    `json:"topic"`
	Message []Message `json:"message"`
}

type Message struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Time    int    `json:"time"`
}

type TopicConfig struct {
	CallbackURL string `json:"callback_url"`
	Timeout     int32  `json:"timeout"`
}

// Check
func (s *PushParams) Check() error {
	if s.Topic == "" {
		return fmt.Errorf("主题信息不能为空")
	}
	if s.CallbackURL == "" {
		return fmt.Errorf("消息通知回调地址不能为空")
	}
	if s.Timeout == 0 {
		s.Timeout = 5
	}
	if len(s.Message) == 0 {
		return fmt.Errorf("消息为空")
	}
	for k, v := range s.Message {
		if v.ID == "" {
			return fmt.Errorf("第 %d 消息条消息 ID 为空", k+1)
		}
		if v.Content == "" {
			return fmt.Errorf("第 %d 消息条消息,消息内容为空", k+1)
		}
	}
	return nil
}

// Push 推送数据
func Push(data *PushParams) error {
	if err := data.Check(); err != nil {
		return err
	}
	payload, _ := json.Marshal(data)
	client := component.Instance.Redis
	if err := client.Kv.SAdd(TopicSet, data.Topic).Err(); err != nil {
		logger.Errorf("设置主题失败 %s payload:%s", err, string(payload))
		return err
	}
	var m = map[string]interface{}{
		"callback": data.CallbackURL,
		"timeout":  data.Timeout,
	}
	members := make([]*redis.Z, 0, len(data.Message))
	for _, item := range data.Message {
		body, _ := json.Marshal(item)
		innerID := uuid.NewString()
		m[innerID] = string(body)
		members = append(members, &redis.Z{
			Score:  float64(item.Time),
			Member: innerID,
		})
	}
	// 写 job
	if _, err := client.HMSet(data.Topic, m); err != nil {
		logger.Errorf("hmset 失败 %s payload:%s", err, string(payload))
		return err
	}
	// 写 bucket
	if err := client.Kv.ZAdd(GetScanBucketName(data.Topic), members...).Err(); err != nil {
		logger.Errorf("zadd 失败 %s payload %s", err, string(payload))
		return err
	}
	return nil
}

func GetScanBucketName(topic string) string {
	return "bucket:" + topic
}

func GetReadyListName(topic string) string {
	return "ready:" + topic
}

func GetTopicConfig(topic string) *TopicConfig {
	client := component.Instance.Redis.Kv
	fields := []string{
		"callback", "timeout",
	}
	if client.Exists(topic).Val() == 0 {
		return nil
	}
	cmd := client.HMGet(topic, fields...)
	if cmd.Err() != nil {
		return nil
	}
	// nolint
	if len(cmd.Val()) < 2 {
		logger.Errorf("获取配置信息数据错误 %+v", cmd.Val())
		return nil
	}
	timeout := cmd.Val()[1].(string)
	t, _ := strconv.ParseInt(timeout, 10, 32) // nolint
	return &TopicConfig{
		CallbackURL: cmd.Val()[0].(string),
		Timeout:     int32(t),
	}
}

// GetTopics 获取主题
func GetTopics() []string {
	// 获取所有的主题
	cmd := component.Instance.Redis.Kv.SMembers(TopicSet)
	if cmd.Err() != nil {
		logger.Errorf("scan 扫描时获取所有主题失败 %s", cmd.Err())
		return nil
	}
	return cmd.Val()
}

// DeleteTopic 删除主题
func DeleteTopic(topic string) error {
	client := component.Instance.Redis.Kv
	if err := client.SRem(TopicSet, topic).Err(); err != nil {
		logger.Errorf("删除主题失败 %s", err)
		return fmt.Errorf("删除主题失败 %s", err)
	}
	if err := client.Del(GetReadyListName(topic), GetScanBucketName(topic), topic).Err(); err != nil {
		logger.Errorf("删除主题失败 %s", err)
		return fmt.Errorf("删除主题失败 %s", err)
	}
	return nil
}

// TopicDetail 主题详情
type TopicDetail struct {
	TopicConfig
	WaitMessage []*Message `json:"wait_message"`
}

// GetWaitDataByTopic 获取正在等待被消费的数据
func GetWaitDataByTopic(topic string) (*TopicDetail, error) {
	client := component.Instance.Redis.Kv
	resp := &TopicDetail{
		TopicConfig: TopicConfig{},
		WaitMessage: make([]*Message, 0),
	}
	if client.Exists(topic).Val() == 0 {
		return resp, nil
	}
	cmd := client.HGetAll(topic)
	if cmd.Err() != nil {
		return resp, cmd.Err()
	}
	for k, item := range cmd.Val() {
		if k == "callback" {
			resp.TopicConfig.CallbackURL = item
			continue
		}
		if k == "timeout" {
			timeout, err := strconv.ParseInt(item, 10, 32) // nolint
			if err != nil {
				return resp, err
			}
			resp.TopicConfig.Timeout = int32(timeout)
			continue
		}
		tmp := new(Message)
		if err := json.Unmarshal([]byte(item), tmp); err != nil {
			return resp, err
		}
		resp.WaitMessage = append(resp.WaitMessage, tmp)
	}
	sort.Slice(resp.WaitMessage, func(i, j int) bool {
		return resp.WaitMessage[i].Time < resp.WaitMessage[j].Time
	})
	return resp, nil
}
