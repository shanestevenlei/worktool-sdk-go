package types

// --- Robot Info ---

// RobotInfo represents the robot information returned by GetRobotInfo.
type RobotInfo struct {
	RobotID             string `json:"robotId"`             // 机器人 id
	Name                string `json:"name"`                // 企微昵称
	OpenCallback        int    `json:"openCallback"`        // 消息回调 0=关闭 1=开启，见 OpenCallback 常量
	EncryptType         int    `json:"encryptType"`         // 加解密方式 0=不加密 1=AES，见 EncryptType 常量
	CreateTime          string `json:"createTime"`          // 创建时间
	EnableAdd           bool   `json:"enableAdd"`           // 能否添加好友
	ReplyAll            int    `json:"replyAll"`            // 回复策略 0=关闭 1=开启，见 ReplyAll 常量
	RobotKeyCheck       int    `json:"robotKeyCheck"`       // key 校验 0=关闭 1=开启，见 RobotKeyCheck 常量
	CallBackRequestType int    `json:"callBackRequestType"` // 回调请求格式 1=form-data 2=json，见 CallbackRequestType 常量
	RobotType           int    `json:"robotType"`           // 机器人类型 0=企业微信 1=微信，见 RobotType 常量
}

// GetRobotInfoResponse is the response for GetRobotInfo.
type GetRobotInfoResponse struct {
	APIResponse
	Data *RobotInfo `json:"data"`
}

// IsOnlineResponse is the response for IsOnline.
type IsOnlineResponse struct {
	Code    int    `json:"code"` // 200=在线，其他=不在线，见 RobotOnlineCode
	Message string `json:"message"`
}

// RobotOnlineStatus 机器人在线状态码（IsOnlineResponse.Code）。
type RobotOnlineStatus int

// RobotOnlineCode 机器人在线时 IsOnline 返回的 code 值。
const RobotOnlineCode RobotOnlineStatus = 200

// OpenCallback 消息回调开关（SetQACallbackRequest.OpenCallback）。
type OpenCallback int

const (
	OpenCallbackDisabled OpenCallback = 0 // 关闭
	OpenCallbackEnabled  OpenCallback = 1 // 开启
)

// ReplyAll 回复策略（RobotInfo.ReplyAll）。
type ReplyAll int

// ReplyAllStrategy 回复策略（SetQACallbackRequest.ReplyAll，API 使用 string）。
type ReplyAllStrategy string

const (
	ReplyAllDisabled ReplyAll = 0 // 关闭
	ReplyAllEnabled  ReplyAll = 1 // 开启

	ReplyAllStrategyDisabled ReplyAllStrategy = "0" // 关闭
	ReplyAllStrategyEnabled  ReplyAllStrategy = "1" // 开启
)

// EncryptType 加解密方式（RobotInfo.EncryptType）。
type EncryptType int

const (
	EncryptTypeNone EncryptType = 0 // 不加密
	EncryptTypeAES  EncryptType = 1 // AES
)

// RobotKeyCheck key 校验开关（RobotInfo.RobotKeyCheck）。
type RobotKeyCheck int

const (
	RobotKeyCheckDisabled RobotKeyCheck = 0 // 关闭
	RobotKeyCheckEnabled  RobotKeyCheck = 1 // 开启
)

// CallbackRequestType 回调请求格式（RobotInfo.CallBackRequestType）。
type CallbackRequestType int

const (
	CallbackRequestTypeFormData CallbackRequestType = 1 // form-data
	CallbackRequestTypeJSON     CallbackRequestType = 2 // json
)

// RobotType 机器人类型（RobotInfo.RobotType）。
type RobotType int

const (
	RobotTypeWeCom  RobotType = 0 // 企业微信
	RobotTypeWeChat RobotType = 1 // 微信
)

// LoginLogType 登录日志类型（LoginLogEntry.Type）。
type LoginLogType int

const (
	LoginLogTypeLogin  LoginLogType = 1 // 登录
	LoginLogTypeLogout LoginLogType = 2 // 登出
)

// --- QA Message Callback ---

// SetQACallbackRequest configures the QA message callback URL.
// WorkTool POSTs incoming chat messages to callbackUrl; your server replies synchronously.
// See: https://doc.worktool.ymdyes.cn/doc-861677.md
type SetQACallbackRequest struct {
	OpenCallback int    `json:"openCallback"` // 0=关闭 1=开启，见 OpenCallback 常量
	CallbackURL  string `json:"callbackUrl"`  // QA 回调 URL
	ReplyAll     string `json:"replyAll"`     // 回复策略 "0"/"1"，见 ReplyAllStrategy 常量
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
	Type      int    `json:"type"` // 1=登录 2=登出，见 LoginLogType 常量
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

// EventCallbackType 事件回调类型（SetEventCallbackRequest.Type 等）。
// 0=群二维码 1=指令结果 5=上线 6=下线 11=消息回调。
type EventCallbackType int

const (
	EventCallbackTypeGroupQR     EventCallbackType = 0  // 群二维码
	EventCallbackTypeCommandExec EventCallbackType = 1  // 指令结果
	EventCallbackTypeOnline      EventCallbackType = 5  // 上线
	EventCallbackTypeOffline     EventCallbackType = 6  // 下线
	EventCallbackTypeMessageRecv EventCallbackType = 11 // 消息回调
)

// SetEventCallbackRequest binds an event callback URL by type.
// Path: POST /robot/robotInfo/callBack/bind
type SetEventCallbackRequest struct {
	Type        int    `json:"type"` // 见 EventCallbackType 常量
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
