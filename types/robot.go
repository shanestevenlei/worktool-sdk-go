package types

// --- Robot Info ---

// RobotInfo represents the robot information returned by GetRobotInfo.
type RobotInfo struct {
	RobotID             string `json:"robotId"`             // 机器人id
	Name                string `json:"name"`                // 企微昵称
	OpenCallback        int    `json:"openCallback"`        // 消息回调地址
	EncryptType         int    `json:"encryptType"`         // 加解密方式 0=不加密 1=AES
	CreateTime          string `json:"createTime"`          // 创建时间
	EnableAdd           bool   `json:"enableAdd"`           // 能否添加好友
	ReplyAll            int    `json:"replyAll"`            // 回复策略
	RobotKeyCheck       int    `json:"robotKeyCheck"`       // key校验 0关闭 1开启
	CallBackRequestType int    `json:"callBackRequestType"` // 1:form-data 2:json
	RobotType           int    `json:"robotType"`           // 机器人类型 0=企业微信 1=微信
}

// GetRobotInfoResponse is the response for GetRobotInfo.
type GetRobotInfoResponse struct {
	APIResponse
	Data *RobotInfo `json:"data"`
}

// IsOnlineResponse is the response for IsOnline.
type IsOnlineResponse struct {
	Code    int    `json:"code"`    // 200=在线，其他=不在线
	Message string `json:"message"`
}

// --- QA Message Callback ---

// SetQACallbackRequest configures the QA message callback URL.
// WorkTool POSTs incoming chat messages to callbackUrl; your server replies synchronously.
// See: https://doc.worktool.ymdyes.cn/doc-861677.md
type SetQACallbackRequest struct {
	OpenCallback int    `json:"openCallback"` // 0关闭 1开启
	CallbackURL  string `json:"callbackUrl"`  // QA回调URL
	ReplyAll     string `json:"replyAll"`     // 回复策略
}

// SetQACallbackResponse is the API response.
type SetQACallbackResponse = APIResponse

// --- Group List ---

// GroupItem represents a single group in the list.
type GroupItem struct {
	GroupID   string `json:"groupId"`
	GroupName string `json:"groupName"`
}

// GetGroupListResponse is the API response.
type GetGroupListResponse struct {
	APIResponse
	Data []*GroupItem `json:"data"`
}

// --- Login Logs ---

// GetLoginLogsRequest for querying robot login history.
// date: optional "yyyy-MM-dd" filter; key: optional verification code.
type GetLoginLogsRequest struct {
	Key  string `json:"key,omitempty"`
	Date string `json:"date,omitempty"`
}

// LoginLogEntry represents a single login/logout event.
type LoginLogEntry struct {
	ID        int64  `json:"id"`
	RobotID   string `json:"robotId"`
	IP        string `json:"ip"`
	LoginTime string `json:"loginTime"`
	Type      int    `json:"type"` // 1=login, 2=logout
	Success   bool   `json:"success"`
	Message   string `json:"message"`
}

// LoginLogResponse wraps APIResponse with a list of entries.
type LoginLogResponse struct {
	APIResponse
	Data []*LoginLogEntry `json:"data"`
}

// --- Corporation List (custom integration) ---

// GetCorpListRequest for fetching the robot's corporation list.
type GetCorpListRequest struct {
	Key string `json:"key,omitempty"`
}

// Corporation represents a single corporation entry.
type Corporation struct {
	CorpID   string `json:"corpId"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// CorpListResponse wraps APIResponse with a list of corporations.
type CorpListResponse struct {
	APIResponse
	Data []*Corporation `json:"data"`
}

// --- Event Callback Bindings (group QR, command exec, online/offline, etc.) ---

// Event callback types supported by the robot.
// 0=群二维码 1=指令结果 5=上线 6=下线 11=消息回调 (new).
const (
	EventCallbackTypeGroupQR     = 0
	EventCallbackTypeCommandExec = 1
	EventCallbackTypeOnline      = 5
	EventCallbackTypeOffline     = 6
	EventCallbackTypeMessageRecv = 11
)

// SetEventCallbackRequest binds an event callback URL by type.
// Path: POST /robot/robotInfo/callBack/bind
type SetEventCallbackRequest struct {
	Type        int    `json:"type"` // see EventCallbackType constants
	CallBackURL string `json:"callBackUrl"`
}

// ListEventCallbacksRequest for fetching configured event callbacks.
type ListEventCallbacksRequest struct {
	RobotKey string `json:"robotKey,omitempty"`
}

// EventCallbackBinding represents a single configured event callback.
type EventCallbackBinding struct {
	ID          int64  `json:"id"`
	Type        int    `json:"type"`
	CallBackURL string `json:"callBackUrl"`
	TypeName    string `json:"typeName"`
}

// EventCallbackListResponse wraps APIResponse with a list of bindings.
type EventCallbackListResponse struct {
	APIResponse
	Data []*EventCallbackBinding `json:"data"`
}

// DeleteEventCallbackRequest removes an event callback by type.
// Path: POST /robot/robotInfo/callBack/deleteByType
type DeleteEventCallbackRequest struct {
	Type int `json:"type"`
}
