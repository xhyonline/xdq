package services

import (
	"encoding/json"
	"fmt"

	"github.com/xhyonline/xutil/logger"

	"github.com/go-redis/redis/v7"
	"github.com/google/uuid"
	"github.com/xhyonline/xdq/component"
)

const TopicSet = "topic"

// PushParams
type PushParams struct {
	Topic       string `json:"topic"`
	CallbackURL string `json:"callback_url"`
	Timeout     int32  `json:"timeout"`
	Message     []struct {
		ID      string `json:"id"`
		Content string `json:"content"`
		Time    int    `json:"time"`
	} `json:"message"`
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
	if err := client.Kv.ZAdd(data.Topic, members...).Err(); err != nil {
		logger.Errorf("zadd 失败 %s payload %s", err, string(payload))
		return err
	}
	return nil
}

// GetTopics 获取主题
func GetTopics() {

}
