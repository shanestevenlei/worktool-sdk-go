package types

// --- History Messages ---

// GetHistoryRequest for querying historical messages (GET, query params only).
type GetHistoryRequest struct {
	Title     string `json:"title,omitempty"`      // filter by chat name
	Page      string `json:"page,omitempty"`       // page number
	Size      string `json:"size,omitempty"`       // page size
	Sort      string `json:"sort,omitempty"`        // e.g. "create_time,desc"
	StartTime string `json:"startTime,omitempty"`  // "2006-01-02 15:04:05"
	EndTime   string `json:"endTime,omitempty"`    // "2006-01-02 15:04:05"
}

// GetRawMessagesRequest for querying sent raw commands (GET /wework/listRawMessage).
type GetRawMessagesRequest struct {
	MessageID string `json:"messageId,omitempty"` // filter by message id
	Page      string `json:"page,omitempty"`      // page number
	Size      string `json:"size,omitempty"`      // page size
	Sort      string `json:"sort,omitempty"`      // e.g. "create_time,desc"
}

// MessageRecord represents a single message in history.
type MessageRecord struct {
	WorkType    string `json:"workType"`    // 工作类型
	TitleList   string `json:"titleList"`   // 群聊或私聊名
	NameList    string `json:"nameList"`    // 消息发送人
	Sender      int64  `json:"sender"`
	Type        int    `json:"type"`        // 消息类型
	ItemMsgList string `json:"itemMsgList"` // 消息内容
	CreateTime  string `json:"createTime"`  // 创建时间
	RobotID     string `json:"robotId"`
}

// MessagePage represents a paginated message list.
type MessagePage struct {
	PageNum  int              `json:"pageNum"`
	PageSize int              `json:"pageSize"`
	TotalPage int             `json:"totalPage"`
	Total    int              `json:"total"`
	List     []*MessageRecord `json:"list"`
}

// GetHistoryResponse is the response for GetHistoryMessages.
type GetHistoryResponse struct {
	APIResponse
	Data *MessagePage `json:"data"`
}

// --- Event Callback Log ---

// GetEventCallbackLogRequest queries event callback logs (GET).
type GetEventCallbackLogRequest struct {
	Name      string `json:"name,omitempty"`
	Page      string `json:"page,omitempty"`
	Size      string `json:"size,omitempty"`
	Sort      string `json:"sort,omitempty"`
	StartTime string `json:"startTime,omitempty"`
	EndTime   string `json:"endTime,omitempty"`
}

// EventCallbackLogRecord is a single event callback log entry.
type EventCallbackLogRecord struct {
	ID          int64    `json:"id"`
	MessageID   string   `json:"messageId"`
	ErrorCode   int      `json:"errorCode"`
	ErrorReason string   `json:"errorReason"`
	RunTime     int64    `json:"runTime"`
	TimeCost    float64  `json:"timeCost"`
	Type        int      `json:"type"`
	RawMsg      string   `json:"rawMsg"`
	SuccessList []string `json:"successList"`
	FailList    []string `json:"failList"`
	GroupName   string   `json:"groupName"`
	QRCode      string   `json:"qrCode"`
	CreateTime  string   `json:"createTime"`
}

// GetEventCallbackLogResponse is the response for GetEventCallbackLog.
type GetEventCallbackLogResponse struct {
	APIResponse
}

// --- Command Log ---

// CommandLogItem represents a sent command log entry.
type CommandLogItem struct {
	ID         int64  `json:"id"`
	MessageID  string `json:"messageId"`
	Type       int    `json:"type"`
	Content    string `json:"content"`
	Status     int    `json:"status"`
	CreateTime string `json:"createTime"`
}

// GetCommandLogResponse is the response for GetCommandLog.
type GetCommandLogResponse struct {
	APIResponse
}
