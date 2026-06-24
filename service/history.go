package service

import (
	"github.com/shanestevenlei/worktool-sdk-go/internal/client"
	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// HistoryService handles historical message and log queries.
// It carries no HTTP state; each method builds its own HTTP client.
type HistoryService struct {
	factory HTTPClientFactory
}

// NewHistoryService creates a new HistoryService.
func NewHistoryService() *HistoryService {
	return &HistoryService{}
}

// SetHTTPFactory wires up the HTTP client factory.
func (s *HistoryService) SetHTTPFactory(f HTTPClientFactory) {
	s.factory = f
}

// http returns a fresh HTTP client for the current request.
func (s *HistoryService) http() *client.HTTPClient {
	return s.factory.HTTPClient()
}

// GetHistoryMessages returns historical messages (deprecated; use callback instead).
// Path: GET /robot/wework/message
func (s *HistoryService) GetHistoryMessages(req *types.GetHistoryRequest) (*types.GetHistoryResponse, error) {
	params := map[string]string{}
	if req.Title != "" {
		params["title"] = req.Title
	}
	if req.Page != "" {
		params["page"] = req.Page
	}
	if req.Size != "" {
		params["size"] = req.Size
	}
	if req.Sort != "" {
		params["sort"] = req.Sort
	}
	if req.StartTime != "" {
		params["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		params["endTime"] = req.EndTime
	}
	var resp types.GetHistoryResponse
	err := s.http().DoGET("/robot/wework/message", params, &resp)
	return &resp, err
}

// GetEventCallbackLog retrieves event callback log entries.
// Path: GET /robot/qaLog/list
func (s *HistoryService) GetEventCallbackLog(req *types.GetEventCallbackLogRequest) (*types.GetEventCallbackLogResponse, error) {
	params := map[string]string{}
	if req.Name != "" {
		params["name"] = req.Name
	}
	if req.Page != "" {
		params["page"] = req.Page
	}
	if req.Size != "" {
		params["size"] = req.Size
	}
	if req.Sort != "" {
		params["sort"] = req.Sort
	}
	if req.StartTime != "" {
		params["startTime"] = req.StartTime
	}
	if req.EndTime != "" {
		params["endTime"] = req.EndTime
	}
	var resp types.GetEventCallbackLogResponse
	err := s.http().DoGET("/robot/qaLog/list", params, &resp)
	return &resp, err
}

// GetRawMessages retrieves sent raw command log entries.
// Path: GET /wework/listRawMessage
func (s *HistoryService) GetRawMessages(req *types.GetRawMessagesRequest) (*types.GetHistoryResponse, error) {
	params := map[string]string{}
	if req.MessageID != "" {
		params["messageId"] = req.MessageID
	}
	if req.Page != "" {
		params["page"] = req.Page
	}
	if req.Size != "" {
		params["size"] = req.Size
	}
	if req.Sort != "" {
		params["sort"] = req.Sort
	}
	var resp types.GetHistoryResponse
	err := s.http().DoGET(listRawMessagePath, params, &resp)
	return &resp, err
}
