package callback

import (
	"encoding/json"
	"strconv"

	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// QA room type values in types.QARequest.RoomType.
const (
	QARoomTypeExternalGroup   = 1
	QARoomTypeExternalContact = 2
	QARoomTypeInternalGroup   = 3
	QARoomTypeInternalContact = 4
)

// QA message text type values in types.QARequest.TextType.
const (
	QATextTypeUnknown      = 0
	QATextTypeText         = 1
	QATextTypeImage        = 2
	QATextTypeVoice        = 3
	QATextTypeVideo        = 5
	QATextTypeMiniProgram  = 7
	QATextTypeLink         = 8
	QATextTypeFile         = 9
	QATextTypeMergedRecord = 13
	QATextTypeReplyText    = 15
)

// QAReplyTypeText is the synchronous text reply type in QAResponse.Data.
const QAReplyTypeText = 5000

const (
	qaCodeSuccess = 0
	qaCodeFailure = -1
)

// ParseQARequest decodes a QA message callback request body.
func ParseQARequest(data []byte) (*QAMessage, error) {
	var raw types.QARequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &QAMessage{QARequest: raw}, nil
}

// QAMessage wraps types.QARequest with helper methods.
type QAMessage struct {
	types.QARequest
}

// IsAtMe reports whether the robot was @mentioned in a group chat.
func (m *QAMessage) IsAtMe() bool {
	v, err := strconv.ParseBool(m.AtMe)
	if err != nil {
		return m.AtMe == "true"
	}
	return v
}

// QAAck returns a successful QA response without synchronous reply content.
func QAAck(message string) *types.QAResponse {
	return &types.QAResponse{Code: qaCodeSuccess, Message: message}
}

// QAFail returns a failed QA response.
func QAFail(message string) *types.QAResponse {
	return &types.QAResponse{Code: qaCodeFailure, Message: message}
}

// QATextReply returns a successful QA response with synchronous text reply.
func QATextReply(text string) *types.QAResponse {
	return &types.QAResponse{
		Code:    qaCodeSuccess,
		Message: "success",
		Data: &types.QAReplyData{
			Type: QAReplyTypeText,
			Info: types.QAReplyInfo{Text: text},
		},
	}
}

// MarshalQAResponse encodes a QA response for net/http handlers.
func MarshalQAResponse(resp *types.QAResponse) ([]byte, error) {
	return json.Marshal(resp)
}
