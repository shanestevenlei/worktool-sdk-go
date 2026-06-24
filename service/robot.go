package service

import (
	"github.com/shanestevenlei/worktool-sdk-go/internal/client"
	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// RobotService handles robot configuration and management.
// It carries no HTTP state; each method builds its own HTTP client
// via the injected factory.
type RobotService struct {
	factory HTTPClientFactory
}

// NewRobotService creates a new RobotService.
func NewRobotService() *RobotService {
	return &RobotService{}
}

// SetHTTPFactory wires up the HTTP client factory.
func (s *RobotService) SetHTTPFactory(f HTTPClientFactory) {
	s.factory = f
}

// http returns a fresh HTTP client for the current request.
func (s *RobotService) http() *client.HTTPClient {
	return s.factory.HTTPClient()
}

// GetInfo returns the robot's configuration and status.
// No request body required (GET).
func (s *RobotService) GetInfo() (*types.GetRobotInfoResponse, error) {
	var resp types.GetRobotInfoResponse
	err := s.http().DoGET("/robot/robotInfo/get", nil, &resp)
	return &resp, err
}

// IsOnline checks whether the robot client is currently online.
// No request body required (GET).
func (s *RobotService) IsOnline() (*types.IsOnlineResponse, error) {
	var resp types.IsOnlineResponse
	err := s.http().DoGET("/robot/robotInfo/online", nil, &resp)
	return &resp, err
}

// SetQACallback configures the QA message callback URL and reply strategy.
func (s *RobotService) SetQACallback(req *types.SetQACallbackRequest) (*types.SetQACallbackResponse, error) {
	var resp types.SetQACallbackResponse
	err := s.http().DoPOST("/robot/robotInfo/update", req, &resp)
	return &resp, err
}

// GetGroupList returns the list of groups managed by this robot.
// No request body required (GET).
func (s *RobotService) GetGroupList() (*types.GetGroupListResponse, error) {
	var resp types.GetGroupListResponse
	err := s.http().DoGET("/robot/robotInfo/groupList", nil, &resp)
	return &resp, err
}

// GetLoginLogs retrieves the robot's login history.
// key: optional verification code; date: optional "yyyy-MM-dd" filter.
func (s *RobotService) GetLoginLogs(req *types.GetLoginLogsRequest) (*types.LoginLogResponse, error) {
	params := map[string]string{}
	if req.Key != "" {
		params["key"] = req.Key
	}
	if req.Date != "" {
		params["date"] = req.Date
	}
	var resp types.LoginLogResponse
	err := s.http().DoGET("/robot/robotInfo/onlineInfos", params, &resp)
	return &resp, err
}

// GetCorpList retrieves the list of corporations available to this robot (custom integration).
func (s *RobotService) GetCorpList(req *types.GetCorpListRequest) (*types.CorpListResponse, error) {
	params := map[string]string{}
	if req.Key != "" {
		params["key"] = req.Key
	}
	var resp types.CorpListResponse
	err := s.http().DoGET("/robot/robotInfo/corpList", params, &resp)
	return &resp, err
}

// SetEventCallback binds an event callback of the given type to a URL.
func (s *RobotService) SetEventCallback(req *types.SetEventCallbackRequest) (*types.APIResponse, error) {
	var resp types.APIResponse
	err := s.http().DoPOST("/robot/robotInfo/callBack/bind", req, &resp)
	return &resp, err
}

// ListEventCallbacks lists all event callbacks configured for the robot.
func (s *RobotService) ListEventCallbacks(req *types.ListEventCallbacksRequest) (*types.EventCallbackListResponse, error) {
	params := map[string]string{}
	if req.RobotKey != "" {
		params["robotKey"] = req.RobotKey
	}
	var resp types.EventCallbackListResponse
	err := s.http().DoGET("/robot/robotInfo/callBack/get", params, &resp)
	return &resp, err
}

// DeleteEventCallback removes an event callback by type.
func (s *RobotService) DeleteEventCallback(req *types.DeleteEventCallbackRequest) (*types.EventCallbackListResponse, error) {
	var resp types.EventCallbackListResponse
	err := s.http().DoPOST("/robot/robotInfo/callBack/deleteByType", req, &resp)
	return &resp, err
}

