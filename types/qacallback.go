package types

// QARoomType QA 回调 roomType 取值。
// 1=外部群 2=外部联系人 3=内部群 4=内部联系人。
type QARoomType int

const (
	QARoomTypeExternalGroup   QARoomType = 1 // 外部群
	QARoomTypeExternalContact QARoomType = 2 // 外部联系人
	QARoomTypeInternalGroup   QARoomType = 3 // 内部群
	QARoomTypeInternalContact QARoomType = 4 // 内部联系人
)

// QAReplyType QA 同步回复类型（QAResponse.Data.Type）。
type QAReplyType int

// QAReplyTypeText 当前仅支持文本回复。
const QAReplyTypeText QAReplyType = 5000

// QAResponseCode 返回给 WorkTool 的 QA 回调响应码。
type QAResponseCode int

const (
	QAResponseCodeSuccess QAResponseCode = 0  // 成功
	QAResponseCodeFailure QAResponseCode = -1 // 失败
)

// QARequest is the payload WorkTool sends for QA message callbacks.
// Configure via Robot.SetQACallback.
// See: https://doc.worktool.ymdyes.cn/doc-861677.md
type QARequest struct {
	Spoken       string `json:"spoken"`       // 问题文本（已去除 @ 机器人部分）
	RawSpoken    string `json:"rawSpoken"`    // 原始问题文本
	ReceivedName string `json:"receivedName"` // 提问者昵称
	GroupName    string `json:"groupName"`    // 群名（私聊为空）
	GroupRemark  string `json:"groupRemark"`  // 群备注名（私聊为空）
	RoomType     int    `json:"roomType"`     // 见 QARoomType 常量
	AtMe         string `json:"atMe"`         // 是否 @ 机器人，"true"/"false"
	TextType     int    `json:"textType"`     // 见 MessageTextType 常量
	FileBase64   string `json:"fileBase64,omitempty"`
	MessageID    string `json:"messageId,omitempty"`
}

// QAResponse is returned to WorkTool from your QA callback endpoint.
type QAResponse struct {
	Code    int          `json:"code"` // 0=成功 -1=失败，见 QAResponseCode 常量
	Message string       `json:"message"`
	Data    *QAReplyData `json:"data,omitempty"`
}

// QAReplyData carries synchronous reply content when code is 0.
type QAReplyData struct {
	Type int         `json:"type"` // 见 QAReplyType 常量（当前仅 QAReplyTypeText）
	Info QAReplyInfo `json:"info"`
}

// QAReplyInfo holds reply payload; currently only text is documented.
type QAReplyInfo struct {
	Text string `json:"text"`
}
