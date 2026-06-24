package service

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"

	"github.com/shanestevenlei/worktool-sdk-go/internal/client"
	"github.com/shanestevenlei/worktool-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

// transportFunc adapts a function to http.RoundTripper for tests.
type transportFunc func(*http.Request) (*http.Response, error)

func (f transportFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

// fakeFactory is a test-only HTTPClientFactory that returns a client
// wrapping the given transport.
type fakeFactory struct {
	transport client.HTTPDoer
}

func (f *fakeFactory) HTTPClient() *client.HTTPClient {
	return client.New(client.Config{HTTPDoer: f.transport})
}

// captureTransport records each request and returns the configured response.
type captureTransport struct {
	mu          sync.Mutex
	calls       []*http.Request
	bodies      []string
	respStatus  int
	respBody    string
	queryParams []map[string]string
}

func (c *captureTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.calls = append(c.calls, req)
	if req.Body != nil {
		body, _ := io.ReadAll(req.Body)
		c.bodies = append(c.bodies, string(body))
	}
	q := map[string]string{}
	for k, v := range req.URL.Query() {
		if len(v) > 0 {
			q[k] = v[0]
		}
	}
	c.queryParams = append(c.queryParams, q)
	return &http.Response{
		StatusCode: c.respStatus,
		Body:       io.NopCloser(strings.NewReader(c.respBody)),
		Header:     http.Header{"Content-Type": []string{"application/json"}},
	}, nil
}

// successResp builds the standard envelope.
func successResp(code int, message string) string {
	b, _ := json.Marshal(types.APIResponse{Code: code, Message: message})
	return string(b)
}

// =============================================================================
// MessageService tests
// =============================================================================

func TestMessageService_SendText(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	resp, err := svc.SendText(&types.SendTextRequest{
		TitleList:       []string{"仑哥"},
		ReceivedContent: "你好",
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "/wework/sendRawMessage", tr.calls[0].URL.Path)
	assert.Contains(t, tr.bodies[0], `"socketType":2`)
	assert.Contains(t, tr.bodies[0], `"type":203`)
	assert.Contains(t, tr.bodies[0], `"titleList":["仑哥"]`)
	assert.Contains(t, tr.bodies[0], `"receivedContent":"你好"`)
}

func TestMessageService_SendText_Validation(t *testing.T) {
	svc := NewMessageService()
	tr := &captureTransport{}
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.SendText(&types.SendTextRequest{TitleList: []string{}, ReceivedContent: "hi"})
	assert.Equal(t, types.ErrEmptyRecipients, err)
	_, err = svc.SendText(&types.SendTextRequest{TitleList: []string{"x"}, ReceivedContent: ""})
	assert.Equal(t, types.ErrEmptyContent, err)
	assert.Equal(t, 0, len(tr.calls), "no request should have been made on validation failure")
}

func TestMessageService_CreateGroup(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.CreateGroup(&types.CreateGroupRequest{
		GroupName:  "测试群",
		SelectList: []string{"仑哥"},
	})
	assert.NoError(t, err)
	assert.Contains(t, tr.bodies[0], `"type":206`)
	assert.Contains(t, tr.bodies[0], `"groupName":"测试群"`)
}

func TestMessageService_DissolveGroup(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.DissolveGroup(&types.DissolveGroupRequest{GroupName: "测试群"})
	assert.NoError(t, err)
	assert.Contains(t, tr.bodies[0], `"type":208`)
}

func TestMessageService_BatchSend(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.BatchSend(&types.BatchSendRequest{
		List: []types.BatchItem{
			{Type: 203, Payload: &types.SendTextRequest{
				TitleList:       []string{"仑哥"},
				ReceivedContent: "msg1",
			}},
			{Type: 206, Payload: &types.CreateGroupRequest{
				GroupName:  "g1",
				SelectList: []string{"x"},
			}},
		},
	})
	assert.NoError(t, err)
	assert.Contains(t, tr.bodies[0], `"type":203`)
	assert.Contains(t, tr.bodies[0], `"type":206`)
}

func TestMessageService_BatchSend_Empty(t *testing.T) {
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: &captureTransport{}})
	_, err := svc.BatchSend(&types.BatchSendRequest{List: nil})
	assert.Equal(t, types.ErrEmptyCommandList, err)
}

func TestMessageService_ForwardMessage(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.ForwardMessage(&types.ForwardMessageRequest{
		TitleList:       []string{"转发群"},
		ReceivedName:    "甲仑",
		OriginalContent: "美团",
		NameList:        []string{"仑哥"},
		TextType:        7,
	})
	assert.NoError(t, err)
	assert.Contains(t, tr.bodies[0], `"type":205`)
}

func TestMessageService_RecallMessage(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.RecallMessage(&types.RecallMessageRequest{MessageID: "msg_1"})
	assert.NoError(t, err)
	assert.Contains(t, tr.bodies[0], `"type":216`)
}

func TestMessageService_RecallMessage_Empty(t *testing.T) {
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: &captureTransport{}})
	_, err := svc.RecallMessage(&types.RecallMessageRequest{MessageID: ""})
	assert.Equal(t, types.ErrEmptyMessageID, err)
}

func TestMessageService_ConcurrencySafety(t *testing.T) {
	// Verify that multiple goroutines sharing the same service do not
	// interfere with each other — they each get their own HTTPClient.
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewMessageService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = svc.SendText(&types.SendTextRequest{
				TitleList:       []string{"仑哥"},
				ReceivedContent: "concurrent",
			})
		}()
	}
	wg.Wait()
	assert.Equal(t, 50, len(tr.calls))
}

// =============================================================================
// RobotService tests
// =============================================================================

func TestRobotService_GetInfo(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: `{"code":0,"message":"ok","data":{"robotId":"r1","name":"bot"}}`}
	svc := NewRobotService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	resp, err := svc.GetInfo()
	assert.NoError(t, err)
	assert.Equal(t, "r1", resp.Data.RobotID)
	assert.Equal(t, "/robot/robotInfo/get", tr.calls[0].URL.Path)
}

func TestRobotService_SetEncryption(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewRobotService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.SetEncryption(&types.SetEncryptionRequest{
		SecretKey:   "16bytekey1234567",
		EncryptType: 1,
	})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/robotInfo/update", tr.calls[0].URL.Path)
	assert.Contains(t, tr.bodies[0], `"secretKey":"16bytekey1234567"`)
	assert.Contains(t, tr.bodies[0], `"encryptType":1`)
}

func TestRobotService_BindCallback(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewRobotService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.BindCallback(&types.BindCallbackRequest{
		Type:        types.CallbackTypeCommandExec,
		CallBackURL: "https://example.com/cb",
	})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/robotInfo/callBack/bind", tr.calls[0].URL.Path)
}

func TestRobotService_DeleteCallbackLegacy(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewRobotService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.DeleteCallbackLegacy(&types.DeleteCallbackLegacyRequest{IDs: []int64{1, 2, 3}})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/robotInfo/callBack/del", tr.calls[0].URL.Path)
	assert.Equal(t, "[1,2,3]", tr.bodies[0])
}

func TestRobotService_GetLoginLogs(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: `{"code":0,"data":[]}`}
	svc := NewRobotService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.GetLoginLogs(&types.GetLoginLogsRequest{Date: "2025-01-01"})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/robotInfo/onlineInfos", tr.calls[0].URL.Path)
	assert.Equal(t, "2025-01-01", tr.queryParams[0]["date"])
}

// =============================================================================
// HistoryService tests
// =============================================================================

func TestHistoryService_GetHistoryMessages(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: `{"code":0,"data":{"list":[]}}`}
	svc := NewHistoryService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.GetHistoryMessages(&types.GetHistoryRequest{
		Title: "仑哥",
		Page:  "1",
		Size:  "20",
	})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/wework/message", tr.calls[0].URL.Path)
	assert.Equal(t, "仑哥", tr.queryParams[0]["title"])
	assert.Equal(t, "20", tr.queryParams[0]["size"])
}

func TestHistoryService_GetQALog(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: successResp(0, "ok")}
	svc := NewHistoryService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.GetQALog(&types.GetQALogRequest{Name: "测试群"})
	assert.NoError(t, err)
	assert.Equal(t, "/robot/qaLog/list", tr.calls[0].URL.Path)
	assert.Equal(t, "测试群", tr.queryParams[0]["name"])
}

func TestHistoryService_GetRawMessages(t *testing.T) {
	tr := &captureTransport{respStatus: 200, respBody: `{"code":0,"data":{"list":[]}}`}
	svc := NewHistoryService()
	svc.SetHTTPFactory(&fakeFactory{transport: tr})

	_, err := svc.GetRawMessages(&types.GetRawMessagesRequest{MessageID: "msg_1"})
	assert.NoError(t, err)
	assert.Equal(t, "/wework/listRawMessage", tr.calls[0].URL.Path)
	assert.Equal(t, "msg_1", tr.queryParams[0]["messageId"])
}
