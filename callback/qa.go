package callback

import (
	"encoding/json"
	"strconv"

	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// ParseQARequest decodes a QA callback request body.
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
	return &types.QAResponse{Code: int(types.QAResponseCodeSuccess), Message: message}
}

// QAFail returns a failed QA response.
func QAFail(message string) *types.QAResponse {
	return &types.QAResponse{Code: int(types.QAResponseCodeFailure), Message: message}
}

// QATextReply returns a successful QA response with synchronous text reply.
func QATextReply(text string) *types.QAResponse {
	return &types.QAResponse{
		Code:    int(types.QAResponseCodeSuccess),
		Message: "success",
		Data: &types.QAReplyData{
			Type: int(types.QAReplyTypeText),
			Info: types.QAReplyInfo{Text: text},
		},
	}
}

// MarshalQAResponse encodes a QA response for net/http handlers.
func MarshalQAResponse(resp *types.QAResponse) ([]byte, error) {
	return json.Marshal(resp)
}
