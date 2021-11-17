package services

// PushParams
type PushParams struct {
	Topic   string `json:"topic"`
	Message []struct {
		ID      string `json:"id"`
		Content string `json:"content"`
		Time    int    `json:"time"`
	} `json:"message"`
}

// Push 推送数据
func Push(data *PushParams) error {
	return nil
}
