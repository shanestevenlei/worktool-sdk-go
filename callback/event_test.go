package callback

import (
	"testing"

	"github.com/shanestevenlei/worktool-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

func TestEventParser_Parse(t *testing.T) {
	p := NewEventParser()
	result, err := p.Parse([]byte(`{"messageId":"m1","errorCode":0,"errorReason":""}`))
	assert.NoError(t, err)
	assert.Equal(t, "m1", result.MessageID)
	assert.True(t, result.IsSuccess())
	assert.Equal(t, "success", result.ErrorMessage())
}

func TestEventParser_Parse_Failure(t *testing.T) {
	p := NewEventParser()
	result, err := p.Parse([]byte(`{"messageId":"m2","errorCode":201102,"errorReason":"send fail"}`))
	assert.NoError(t, err)
	assert.Equal(t, "m2", result.MessageID)
	assert.Equal(t, 201102, result.ErrorCode)
	assert.False(t, result.IsSuccess())
	assert.Contains(t, result.ErrorMessage(), "发送消息失败")
}

func TestEventParser_Parse_InvalidJSON(t *testing.T) {
	p := NewEventParser()
	_, err := p.Parse([]byte(`not json`))
	assert.Error(t, err)
}

func TestParseEventRequest(t *testing.T) {
	result, err := ParseEventRequest([]byte(`{"messageId":"x","errorCode":0}`))
	assert.NoError(t, err)
	assert.Equal(t, "x", result.MessageID)
	assert.True(t, result.IsSuccess())
}

func TestEventResult_ErrorMessage_Unknown(t *testing.T) {
	result := &EventResult{
		EventCallbackRequest: types.EventCallbackRequest{
			MessageID:   "m",
			ErrorCode:   999999,
			ErrorReason: "custom reason",
		},
	}
	assert.Equal(t, "custom reason", result.ErrorMessage())
}
