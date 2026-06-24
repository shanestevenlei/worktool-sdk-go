package types

// QARequest is the payload WorkTool sends for QA message callbacks.
// Configure via Robot.SetQACallback.
// See: https://doc.worktool.ymdyes.cn/doc-861677.md
type QARequest struct {
	Spoken       string `json:"spoken"`
	RawSpoken    string `json:"rawSpoken"`
	ReceivedName string `json:"receivedName"`
	GroupName    string `json:"groupName"`
	GroupRemark  string `json:"groupRemark"`
	RoomType     int    `json:"roomType"`
	AtMe         string `json:"atMe"`
	TextType     int    `json:"textType"`
	FileBase64   string `json:"fileBase64,omitempty"`
	MessageID    string `json:"messageId,omitempty"`
}

// QAResponse is returned to WorkTool from your QA callback endpoint.
type QAResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    *QAReplyData `json:"data,omitempty"`
}

// QAReplyData carries synchronous reply content when code is 0.
type QAReplyData struct {
	Type int         `json:"type"`
	Info QAReplyInfo `json:"info"`
}

// QAReplyInfo holds reply payload; currently only text is documented.
type QAReplyInfo struct {
	Text string `json:"text"`
}
