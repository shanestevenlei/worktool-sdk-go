package callback

import (
	"encoding/json"
	"testing"

	"github.com/shanestevenlei/worktool-sdk-go/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseQARequest(t *testing.T) {
	body := []byte(`{
		"spoken": "你好",
		"rawSpoken": "@管家 你好",
		"receivedName": "仑哥",
		"groupName": "测试群1",
		"groupRemark": "测试群1备注名",
		"roomType": 1,
		"atMe": "true",
		"textType": 1
	}`)

	msg, err := ParseQARequest(body)
	require.NoError(t, err)
	assert.Equal(t, "你好", msg.Spoken)
	assert.Equal(t, "@管家 你好", msg.RawSpoken)
	assert.Equal(t, "仑哥", msg.ReceivedName)
	assert.Equal(t, int(types.QARoomTypeExternalGroup), msg.RoomType)
	assert.True(t, msg.IsAtMe())
	assert.Equal(t, int(types.MessageTextTypeText), msg.TextType)
}

func TestParseQARequest_InvalidJSON(t *testing.T) {
	_, err := ParseQARequest([]byte(`not json`))
	assert.Error(t, err)
}

func TestQAMessage_IsAtMe(t *testing.T) {
	at := &QAMessage{}
	at.AtMe = "true"
	assert.True(t, at.IsAtMe())

	notAt := &QAMessage{}
	notAt.AtMe = "false"
	assert.False(t, notAt.IsAtMe())
}

func TestQAAck(t *testing.T) {
	resp := QAAck("参数接收成功")
	assert.Equal(t, int(types.QAResponseCodeSuccess), resp.Code)
	assert.Equal(t, "参数接收成功", resp.Message)
	assert.Nil(t, resp.Data)
}

func TestQATextReply(t *testing.T) {
	resp := QATextReply("你好，有什么可以帮您？")
	assert.Equal(t, int(types.QAResponseCodeSuccess), resp.Code)
	require.NotNil(t, resp.Data)
	assert.Equal(t, int(types.QAReplyTypeText), resp.Data.Type)
	assert.Equal(t, "你好，有什么可以帮您？", resp.Data.Info.Text)
}

func TestMarshalQAResponse(t *testing.T) {
	data, err := MarshalQAResponse(QATextReply("ok"))
	require.NoError(t, err)

	var decoded map[string]any
	require.NoError(t, json.Unmarshal(data, &decoded))
	assert.Equal(t, float64(int(types.QAResponseCodeSuccess)), decoded["code"])
	dataObj := decoded["data"].(map[string]any)
	assert.Equal(t, float64(int(types.QAReplyTypeText)), dataObj["type"])
}

func TestQAFail(t *testing.T) {
	resp := QAFail("处理失败")
	assert.Equal(t, int(types.QAResponseCodeFailure), resp.Code)
	assert.Equal(t, "处理失败", resp.Message)
}
