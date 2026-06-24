package types

// EventCallbackRequest is the payload WorkTool sends to an event callback URL.
// For type=1 (EventCallbackTypeCommandExec), this is the command execution result.
// Configure via Robot.SetEventCallback.
// See: https://doc.worktool.ymdyes.cn/api-44952776.md
type EventCallbackRequest struct {
	MessageID   string   `json:"messageId"`
	ErrorCode   int      `json:"errorCode"` // 0=成功 其他=失败
	ErrorReason string   `json:"errorReason"`
	RunTime     int64    `json:"runTime"`
	TimeCost    float64  `json:"timeCost"`
	Type        int      `json:"type"`
	RawMsg      string   `json:"rawMsg"`
	SuccessList []string `json:"successList"`
	FailList    []string `json:"failList"`
	GroupName   string   `json:"groupName"`
	QRCode      string   `json:"qrCode"`
}
