package services

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/xhyonline/xutil/logger"
)

type callbackParams struct {
	Topic   string     `json:"topic"`
	Message []*Message `json:"message"`
	Count   int        `json:"count"`
}

// CallBack 回调,回调并不关心是否成功,失败则记录日志
func CallBack(url, topic string, message []*Message, timeout time.Duration) bool {
	params := &callbackParams{
		Topic:   topic,
		Message: message,
		Count:   len(message),
	}
	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
		Timeout: timeout * time.Second,
	}
	body, _ := json.Marshal(params)
	r, err := http.NewRequestWithContext(context.Background(), "POST", url, bytes.NewReader(body))
	if err != nil {
		logger.Errorf("构造消费数据失败 %s 元数据:%s", url, string(body))
		return false
	}
	r.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(r)
	if err != nil {
		logger.Errorf("发送消费数据失败地址:%s 错误:%s 元数据:%s", url, err, string(body))
		return false
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		logger.Errorf("状态码: %d 请求地址%s 元数据:%s", resp.StatusCode, url, string(body))
		return false
	}
	return true
}
