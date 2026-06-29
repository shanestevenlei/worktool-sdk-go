package types

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBatchItem_MarshalJSON_OfficialFormat(t *testing.T) {
	data, err := json.Marshal(BatchItem{
		Type: int(CmdTypeSendText),
		Payload: &SendTextRequest{
			TitleList:       []string{"仑哥"},
			ReceivedContent: "你好",
		},
	})
	require.NoError(t, err)

	body := string(data)
	assert.Contains(t, body, `"type":203`)
	assert.Contains(t, body, `"titleList":["仑哥"]`)
	assert.Contains(t, body, `"receivedContent":"你好"`)
	assert.NotContains(t, body, "Payload")
	assert.NotContains(t, body, "payload")
}

func TestMessageRequest_MarshalJSON_OfficialFormat(t *testing.T) {
	data, err := json.Marshal(&MessageRequest{
		SocketType: int(SocketTypeWork),
		List: []BatchItem{
			{
				Type: int(CmdTypeSendMedia),
				Payload: &SendFileRequest{
					TitleList:  []string{"仑哥"},
					ObjectName: "logo.png",
					FileURL:    "https://example.com/logo.png",
					FileType:   string(MediaFileTypeImage),
				},
			},
		},
	})
	require.NoError(t, err)

	body := string(data)
	assert.Contains(t, body, `"socketType":2`)
	assert.Contains(t, body, `"type":218`)
	assert.Contains(t, body, `"fileType":"image"`)
	assert.NotContains(t, body, "Payload")
	assert.NotContains(t, body, "payload")
}

func TestBatchItem_MarshalJSON_NilPayload(t *testing.T) {
	_, err := json.Marshal(BatchItem{Type: int(CmdTypeSendText)})
	assert.ErrorIs(t, err, ErrEmptyPayload)
}
