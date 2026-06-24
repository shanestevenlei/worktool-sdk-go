package callback

import (
	"encoding/json"

	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// EventParser decodes event callback bodies (plain JSON).
type EventParser struct{}

// NewEventParser creates an event callback parser.
func NewEventParser() *EventParser {
	return &EventParser{}
}

// Parse decodes an event callback payload.
func (p *EventParser) Parse(data []byte) (*EventResult, error) {
	var raw types.EventCallbackRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &EventResult{EventCallbackRequest: raw}, nil
}

// ParseEventRequest is a convenience wrapper for net/http handlers.
func ParseEventRequest(data []byte) (*EventResult, error) {
	return NewEventParser().Parse(data)
}

// EventResult wraps types.EventCallbackRequest with helper methods.
type EventResult struct {
	types.EventCallbackRequest
}

// IsSuccess returns true if the event indicates successful command execution.
func (r *EventResult) IsSuccess() bool {
	return r.ErrorCode == int(EventCodeSuccess)
}

// ErrorMessage returns a human-readable message for ErrorCode, or ErrorReason if unknown.
func (r *EventResult) ErrorMessage() string {
	if msg, ok := EventErrorCodeMessages[EventErrorCode(r.ErrorCode)]; ok {
		return msg
	}
	return r.ErrorReason
}

// EventErrorCode 事件回调 errorCode 取值（指令执行结果，type=1）。
type EventErrorCode int

const (
	EventCodeSuccess           EventErrorCode = 0
	EventCodeIllegalData       EventErrorCode = 101011
	EventCodeIllegalOperation  EventErrorCode = 101012
	EventCodeIllegalPermission EventErrorCode = 101013

	EventCodeCreateGroupFail   EventErrorCode = 201011
	EventCodeGroupRenameFail   EventErrorCode = 201012
	EventCodeGroupAddFail      EventErrorCode = 201013
	EventCodeGroupRemoveFail   EventErrorCode = 201014
	EventCodeGroupAnnounceFail EventErrorCode = 201015
	EventCodeGroupRemarkFail   EventErrorCode = 201016
	EventCodeIntoRoomFail      EventErrorCode = 201101
	EventCodeSendMsgFail       EventErrorCode = 201102
	EventCodeButtonFail        EventErrorCode = 201103
	EventCodeTargetFail        EventErrorCode = 201104
	EventCodeRelayFail         EventErrorCode = 201105
	EventCodeRepeat            EventErrorCode = 201106
	EventCodeFileDownload      EventErrorCode = 201107
	EventCodeFileStorage       EventErrorCode = 201108
)

// EventErrorCodeMessages maps an event callback error code to a description.
var EventErrorCodeMessages = map[EventErrorCode]string{
	EventCodeSuccess:           "success",
	EventCodeIllegalData:       "非法数据",
	EventCodeIllegalOperation:  "非法操作",
	EventCodeIllegalPermission: "非法权限",
	EventCodeCreateGroupFail:   "创建群失败",
	EventCodeGroupRenameFail:   "修改群名失败",
	EventCodeGroupAddFail:      "群拉人失败",
	EventCodeGroupRemoveFail:   "群踢人失败",
	EventCodeGroupAnnounceFail: "修改群公告失败",
	EventCodeGroupRemarkFail:   "修改群备注失败",
	EventCodeIntoRoomFail:      "进群失败",
	EventCodeSendMsgFail:       "发送消息失败",
	EventCodeButtonFail:        "按钮失败",
	EventCodeTargetFail:        "目标失败",
	EventCodeRelayFail:         "转发失败",
	EventCodeRepeat:            "重复执行",
	EventCodeFileDownload:      "文件下载失败",
	EventCodeFileStorage:       "文件存储失败",
}
