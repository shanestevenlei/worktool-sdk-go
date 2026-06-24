package service

import (
	"github.com/shanestevenlei/worktool-sdk-go/internal/client"
	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// MessageService handles messaging and group operations.
// It carries no HTTP state; each method builds its own HTTP client.
type MessageService struct {
	factory HTTPClientFactory
}

// NewMessageService creates a new MessageService.
func NewMessageService() *MessageService {
	return &MessageService{}
}

// SetHTTPFactory wires up the HTTP client factory.
func (s *MessageService) SetHTTPFactory(f HTTPClientFactory) {
	s.factory = f
}

// http returns a fresh HTTP client for the current request.
func (s *MessageService) http() *client.HTTPClient {
	return s.factory.HTTPClient()
}

// send is the shared helper for single-command /wework/sendRawMessage calls.
func (s *MessageService) send(cmdType int, payload, resp interface{}) error {
	body := &types.MessageRequest{
		SocketType: 2,
		List:       []types.MessageItem{{Type: cmdType, Payload: payload}},
	}
	return s.http().DoPOST(sendRawMessagePath, body, resp)
}

// -----------------------------------------------------------------------------
// Text & content
// -----------------------------------------------------------------------------

// SendText sends a text message (type=203).
func (s *MessageService) SendText(req *types.SendTextRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(203, req, &resp)
	return &resp, err
}

// LogDebug is identical to SendText but flagged for debug purposes (type=203).
func (s *MessageService) LogDebug(req *types.SendTextRequest) (*types.SendMessageResponse, error) {
	return s.SendText(req)
}

// ForwardMessage forwards a message from one chat to another (type=205).
// Requires a special "xxx小程序转发群" set up beforehand.
func (s *MessageService) ForwardMessage(req *types.ForwardMessageRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(205, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Media (image / file / audio / video) — all share type=218
// -----------------------------------------------------------------------------

// SendImage sends an image (type=218, fileType=image).
func (s *MessageService) SendImage(req *types.SendImageRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(218, req, &resp)
	return &resp, err
}

// SendFile sends a file: audio, video, or other (type=218).
func (s *MessageService) SendFile(req *types.SendFileRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(218, req, &resp)
	return &resp, err
}

// SendAnyMedia is the generic entry for image/audio/video/file (type=218).
// Equivalent to SendImage when fileType=="image", SendFile otherwise.
func (s *MessageService) SendAnyMedia(req *types.SendFileRequest) (*types.SendMessageResponse, error) {
	return s.SendFile(req)
}

// -----------------------------------------------------------------------------
// WeDrive (微盘)
// -----------------------------------------------------------------------------

// SendWeDriveImage pushes an image from WeDrive (type=208). Note: distinct from
// SendImage which downloads from a public URL.
func (s *MessageService) SendWeDriveImage(req *types.SendWeDriveRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(208, req, &resp)
	return &resp, err
}

// SendWeDriveFile pushes a file from WeDrive (type=209).
func (s *MessageService) SendWeDriveFile(req *types.SendWeDriveRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(209, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Tencent docs & forms (deprecated by WeChat, but kept for completeness)
// -----------------------------------------------------------------------------

// SendTencentDoc pushes a Tencent Doc (type=211).
func (s *MessageService) SendTencentDoc(req *types.SendDocRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(211, req, &resp)
	return &resp, err
}

// SendCollector pushes a Tencent Docs collection form (type=211).
func (s *MessageService) SendCollector(req *types.SendDocRequest) (*types.SendMessageResponse, error) {
	return s.SendTencentDoc(req)
}

// -----------------------------------------------------------------------------
// Groups (206 create, 207 update, 208 dissolve)
// -----------------------------------------------------------------------------

// CreateGroup creates an external group (type=206).
func (s *MessageService) CreateGroup(req *types.CreateGroupRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(206, req, &resp)
	return &resp, err
}

// UpdateGroup modifies group properties and membership (type=207).
func (s *MessageService) UpdateGroup(req *types.UpdateGroupRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(207, req, &resp)
	return &resp, err
}

// DissolveGroup dissolves (deletes) an external group (type=208).
// Note: 208 is shared with SendWeDriveImage — the API distinguishes by payload.
func (s *MessageService) DissolveGroup(req *types.DissolveGroupRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(208, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Friends (213 add by phone, 214 modify, 215 add from group, 217 delete)
// -----------------------------------------------------------------------------

// AddFriendByPhone sends a friend request by phone number (type=213).
func (s *MessageService) AddFriendByPhone(req *types.AddFriendByPhoneRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(213, req, &resp)
	return &resp, err
}

// ModifyFriend updates remark, tags, or extra info (type=214).
func (s *MessageService) ModifyFriend(req *types.ModifyFriendRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(214, req, &resp)
	return &resp, err
}

// AddFriendFromGroup sends a friend request to someone in an external group (type=215).
func (s *MessageService) AddFriendFromGroup(req *types.AddFriendFromGroupRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(215, req, &resp)
	return &resp, err
}

// DeleteContact deletes a contact by name list (type=217).
func (s *MessageService) DeleteContact(req *types.DeleteContactRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(217, req, &resp)
	return &resp, err
}

// ModifyGroupMemberRemark sets a remark for a group member (type=219).
func (s *MessageService) ModifyGroupMemberRemark(req *types.ModifyGroupMemberRemarkRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(219, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Message lifecycle (216 recall, 222 insert, 223 clear specific, 224 clear all)
// -----------------------------------------------------------------------------

// RecallMessage recalls a previously sent message by messageId (type=216).
func (s *MessageService) RecallMessage(req *types.RecallMessageRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(216, req, &resp)
	return &resp, err
}

// AddTodo creates a todo/reminder item (type=221).
func (s *MessageService) AddTodo(req *types.AddTodoRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(221, req, &resp)
	return &resp, err
}

// InsertCommand inserts a command at the front of the queue (type=222).
func (s *MessageService) InsertCommand(req *types.InsertCommandRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(222, req, &resp)
	return &resp, err
}

// ClearSpecificCommand clears a specific pending command by messageId (type=223).
func (s *MessageService) ClearSpecificCommand(req *types.ClearSpecificCommandRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	var resp types.SendMessageResponse
	err := s.send(223, req, &resp)
	return &resp, err
}

// ClearCommands clears all pending client commands (type=224).
func (s *MessageService) ClearCommands(req *types.ClearCommandsRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(224, req, &resp)
	return &resp, err
}

// SwitchEnterprise switches the robot's enterprise to the target (type=225).
func (s *MessageService) SwitchEnterprise(req *types.SwitchEnterpriseRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(225, req, &resp)
	return &resp, err
}

// CleanupStorage cleans up WeChat storage (cache + resource files) (type=226).
func (s *MessageService) CleanupStorage(req *types.CleanupStorageRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(226, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Custom (paid) integrations
// -----------------------------------------------------------------------------

// SendLink sends a custom-styled link card (type=302, requires paid authorization).
func (s *MessageService) SendLink(req *types.SendLinkRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(302, req, &resp)
	return &resp, err
}

// SendMiniProgram sends a custom-styled mini-program card (type=303, requires paid authorization).
func (s *MessageService) SendMiniProgram(req *types.SendMiniProgramRequest) (*types.SendMessageResponse, error) {
	var resp types.SendMessageResponse
	err := s.send(303, req, &resp)
	return &resp, err
}

// -----------------------------------------------------------------------------
// Batch
// -----------------------------------------------------------------------------

// BatchSend sends multiple commands in a single request (max 100).
func (s *MessageService) BatchSend(req *types.BatchSendRequest) (*types.SendMessageResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	body := &types.MessageRequest{
		SocketType: 2,
		List:       toMessageItems(req.List),
	}
	var resp types.SendMessageResponse
	err := s.http().DoPOST(sendRawMessagePath, body, &resp)
	return &resp, err
}

// toMessageItems converts BatchItem list to the generic MessageItem list.
func toMessageItems(items []types.BatchItem) []types.MessageItem {
	result := make([]types.MessageItem, len(items))
	for i, item := range items {
		result[i] = types.MessageItem{Type: item.Type, Payload: item.Payload}
	}
	return result
}

// =============================================================================
// Response type aliases — all /wework/sendRawMessage commands share the same
// {code, message, data} envelope. Aliases make call sites self-documenting.
// =============================================================================

type (
	SendLinkResponse                = types.SendMessageResponse
	SendMiniProgramResponse         = types.SendMessageResponse
	SendImageResponse               = types.SendMessageResponse
	SendFileResponse                = types.SendMessageResponse
	SendAnyMediaResponse            = types.SendMessageResponse
	SendWeDriveImageResponse        = types.SendMessageResponse
	SendWeDriveFileResponse         = types.SendMessageResponse
	SendTencentDocResponse          = types.SendMessageResponse
	SendCollectorResponse           = types.SendMessageResponse
	ForwardMessageResponse          = types.SendMessageResponse
	LogDebugResponse                = types.SendMessageResponse
	ModifyFriendResponse            = types.SendMessageResponse
	AddFriendFromGroupResponse      = types.SendMessageResponse
	ModifyGroupMemberRemarkResponse = types.SendMessageResponse
	RecallMessageResponse           = types.SendMessageResponse
	AddTodoResponse                 = types.SendMessageResponse
	InsertCommandResponse           = types.SendMessageResponse
	ClearSpecificCommandResponse    = types.SendMessageResponse
	ClearCommandsResponse           = types.SendMessageResponse
	SwitchEnterpriseResponse        = types.SendMessageResponse
	CleanupStorageResponse          = types.SendMessageResponse
	CreateGroupResponse             = types.SendMessageResponse
	UpdateGroupResponse             = types.SendMessageResponse
	DissolveGroupResponse           = types.SendMessageResponse
	AddFriendByPhoneResponse        = types.SendMessageResponse
	DeleteContactResponse           = types.SendMessageResponse
	BatchSendResponse               = types.SendMessageResponse
)
