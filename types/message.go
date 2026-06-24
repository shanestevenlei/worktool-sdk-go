package types

// APIResponse is the standard WorkTool API response wrapper.
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// IsSuccess returns true if code == 200.
func (r *APIResponse) IsSuccess() bool {
	return r.Code == 200
}

// SocketType /wework/sendRawMessage 请求的 socketType。
type SocketType int

const SocketTypeWork SocketType = 2

// CmdType /wework/sendRawMessage 指令类型（MessageItem.Type）。
type CmdType int

const (
	CmdTypeSendText                CmdType = 203
	CmdTypeForwardMessage          CmdType = 205
	CmdTypeCreateGroup             CmdType = 206
	CmdTypeUpdateGroup             CmdType = 207
	CmdTypeWeDriveImage            CmdType = 208 // DissolveGroup 同为 208，靠 payload 区分
	CmdTypeDissolveGroup           CmdType = 208
	CmdTypeWeDriveFile             CmdType = 209
	CmdTypeSendTencentDoc          CmdType = 211
	CmdTypeAddFriendByPhone        CmdType = 213
	CmdTypeModifyFriend            CmdType = 214
	CmdTypeAddFriendFromGroup      CmdType = 215
	CmdTypeRecallMessage           CmdType = 216
	CmdTypeDeleteContact           CmdType = 217
	CmdTypeSendMedia               CmdType = 218 // 通过 URL 发送 image / audio / video / file
	CmdTypeModifyGroupMemberRemark CmdType = 219
	CmdTypeAddTodo                 CmdType = 221
	CmdTypeInsertCommand           CmdType = 222
	CmdTypeClearSpecificCommand    CmdType = 223
	CmdTypeClearCommands           CmdType = 224
	CmdTypeSwitchEnterprise        CmdType = 225
	CmdTypeCleanupStorage          CmdType = 226
	CmdTypeSendLink                CmdType = 302
	CmdTypeSendMiniProgram         CmdType = 303
)

// MessageTextType 消息内容类型（ForwardMessageRequest.TextType、QARequest.TextType）。
// 0=未知 1=文本 2=图片 3=语音 5=视频 7=小程序 8=链接 9=文件 13=合并聊天记录 15=回复文本。
type MessageTextType int

const (
	MessageTextTypeUnknown      MessageTextType = 0  // 未知
	MessageTextTypeText         MessageTextType = 1  // 文本
	MessageTextTypeImage        MessageTextType = 2  // 图片
	MessageTextTypeVoice        MessageTextType = 3  // 语音
	MessageTextTypeVideo        MessageTextType = 5  // 视频
	MessageTextTypeMiniProgram  MessageTextType = 7  // 小程序
	MessageTextTypeLink         MessageTextType = 8  // 链接
	MessageTextTypeFile         MessageTextType = 9  // 文件
	MessageTextTypeMergedRecord MessageTextType = 13 // 合并聊天记录
	MessageTextTypeReplyText    MessageTextType = 15 // 回复文本
)

// MediaFileType 媒体文件类型（SendFileRequest.FileType，type=218）。
type MediaFileType string

const (
	MediaFileTypeImage MediaFileType = "image" // 图片
	MediaFileTypeAudio MediaFileType = "audio" // 音频
	MediaFileTypeVideo MediaFileType = "video" // 视频
	MediaFileTypeAny   MediaFileType = "*"     // 任意类型
)

// AtEveryone 群聊 @ 所有人时的取值。
const AtEveryone = "@所有人"

// ============================================================================
// Message types
// ============================================================================

// SendTextRequest for sending a text message (type=203).
type SendTextRequest struct {
	TitleList       []string `json:"titleList"`       // recipients: friend nicknames or group names
	ReceivedContent string   `json:"receivedContent"` // text content; \n for newline
	AtList          []string `json:"atList"`          // @ 列表；群聊 @ 所有人用 AtEveryone（"@所有人"）
}

// Validate checks required fields.
func (r *SendTextRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ReceivedContent == "" {
		return ErrEmptyContent
	}
	return nil
}

// SendImageRequest for sending an image (type=218, fileType=image).
type SendImageRequest struct {
	TitleList  []string `json:"titleList"`  // recipients
	ObjectName string   `json:"objectName"` // file name with extension, e.g. "logo.png"
	FileURL    string   `json:"fileUrl"`    // publicly accessible URL
	ExtraText  string   `json:"extraText"`  // optional comment
}

// Validate checks required fields.
func (r *SendImageRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ObjectName == "" {
		return ErrEmptyObjectName
	}
	if r.FileURL == "" {
		return ErrEmptyFileURL
	}
	return nil
}

// SendFileRequest for sending a file (type=218).
type SendFileRequest struct {
	TitleList  []string `json:"titleList"`
	ObjectName string   `json:"objectName"`
	FileURL    string   `json:"fileUrl"`
	FileType   string   `json:"fileType"` // audio / video / *，见 MediaFileType 常量
	ExtraText  string   `json:"extraText"`
}

// Validate checks required fields.
func (r *SendFileRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ObjectName == "" {
		return ErrEmptyObjectName
	}
	if r.FileURL == "" {
		return ErrEmptyFileURL
	}
	return nil
}

// ============================================================================
// Group types
// ============================================================================

// CreateGroupRequest for creating an external group (type=206).
type CreateGroupRequest struct {
	GroupName         string   `json:"groupName"`         // required; avoid duplicate names
	SelectList        []string `json:"selectList"`        // optional initial members
	GroupAnnouncement string   `json:"groupAnnouncement"` // optional group announcement
	GroupRemark       string   `json:"groupRemark"`       // optional group remark
	GroupTemplate     string   `json:"groupTemplate"`     // optional template name
}

// Validate checks required fields.
func (r *CreateGroupRequest) Validate() error {
	if r.GroupName == "" {
		return ErrEmptyGroupName
	}
	return nil
}

// UpdateGroupRequest for modifying group info / membership (type=207).
type UpdateGroupRequest struct {
	GroupName            string   `json:"groupName"`            // target group (use remark name if set)
	NewGroupName         string   `json:"newGroupName"`         // optional rename
	NewGroupAnnouncement string   `json:"newGroupAnnouncement"` // optional new announcement
	SelectList           []string `json:"selectList"`           // members to add
	RemoveList           []string `json:"removeList"`           // members to remove
	ShowMessageHistory   bool     `json:"showMessageHistory"`   // include chat history when adding
	GroupRemark          string   `json:"groupRemark"`          // optional group remark
	GroupTemplate        string   `json:"groupTemplate"`        // optional template name
}

// ============================================================================
// Friend / Contact types
// ============================================================================

// AddFriendByPhoneRequest for adding a friend by phone (type=213).
type AddFriendByPhoneRequest struct {
	Phone      string   `json:"phone"`      // required
	MarkName   string   `json:"markName"`   // optional remark name
	MarkExtra  string   `json:"markExtra"`  // optional extra remark info
	TagList    []string `json:"tagList"`    // optional tags
	LeavingMsg string   `json:"leavingMsg"` // optional friend request message
}

// Validate checks required fields.
func (r *AddFriendByPhoneRequest) Validate() error {
	if r.Phone == "" {
		return ErrEmptyPhone
	}
	return nil
}

// DeleteContactRequest for deleting a contact (type=217).
type DeleteContactRequest struct {
	TitleList []string `json:"titleList"` // required; names to delete
}

// Validate checks required fields.
func (r *DeleteContactRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	return nil
}

// ============================================================================
// Todo / Recall
// ============================================================================

// AddTodoRequest for creating a todo/reminder (type=221).
type AddTodoRequest struct {
	TitleList []string `json:"titleList"`
	Content   string   `json:"content"`
	AtList    []string `json:"atList"`
}

// RecallMessageRequest for recalling a sent message (type=216).
type RecallMessageRequest struct {
	TitleList []string `json:"titleList"` // optional; target chat
	MessageID string   `json:"messageId"` // required; message ID to recall
}

// Validate checks required fields.
func (r *RecallMessageRequest) Validate() error {
	if r.MessageID == "" {
		return ErrEmptyMessageID
	}
	return nil
}

// ============================================================================
// Enterprise / Storage
// ============================================================================

// SwitchEnterpriseRequest for switching robot enterprise (type=225).
type SwitchEnterpriseRequest struct {
	EnterpriseName string `json:"enterpriseName"` // exact name as shown in WeChat Enterprise
}

// CleanupStorageRequest for cleaning storage (type=226).
type CleanupStorageRequest struct {
	// No additional fields required
}

// ClearCommandsRequest for clearing all pending commands (type=224).
type ClearCommandsRequest struct {
	// No additional fields required
}

// ClearSpecificCommandRequest for clearing a specific pending command (type=223).
type ClearSpecificCommandRequest struct {
	MessageID string `json:"messageId"` // required
}

// Validate checks required fields.
func (r *ClearSpecificCommandRequest) Validate() error {
	if r.MessageID == "" {
		return ErrEmptyMessageID
	}
	return nil
}

// ============================================================================
// Batch
// ============================================================================

// BatchSendRequest wraps multiple commands in one request (max 100).
type BatchSendRequest struct {
	List []BatchItem `json:"list"`
}

// Validate checks required fields.
func (r *BatchSendRequest) Validate() error {
	if len(r.List) == 0 {
		return ErrEmptyCommandList
	}
	return nil
}

// BatchItem is a single command in a batch.
// Payload is the concrete request struct (e.g. *SendTextRequest, *CreateGroupRequest).
type BatchItem struct {
	Type    int         `json:"type"`
	Payload interface{} `json:"payload"`
}

// MessageRequest is the top-level request wrapper for /wework/sendRawMessage.
// All command types share this same envelope.
type MessageRequest struct {
	SocketType int           `json:"socketType"` // 固定为 2，见 SocketTypeWork
	List       []MessageItem `json:"list"`
}

// MessageItem is a single command inside a MessageRequest.
type MessageItem struct {
	Type    int         `json:"type"`
	Payload interface{} `json:",inline"` // one of *SendTextRequest, *CreateGroupRequest, etc.
}

// ============================================================================
// Responses
// ============================================================================

// SendMessageResponse wraps APIResponse for text message send.
type SendMessageResponse = APIResponse

// SendMediaResponse wraps APIResponse for media send.
type SendMediaResponse = APIResponse

// CreateGroupResponse wraps APIResponse.
type CreateGroupResponse = APIResponse

// UpdateGroupResponse wraps APIResponse.
type UpdateGroupResponse = APIResponse

// ============================================================================
// Additional message types (documented in API but not yet in service layer)
// ============================================================================

// SendLinkRequest for sending a link card (type=302).
type SendLinkRequest struct {
	TitleList       []string `json:"titleList"`
	ReceivedContent string   `json:"receivedContent"` // link title
	LinkURL         string   `json:"linkUrl"`         // link URL
	PictureURL      string   `json:"pictureUrl"`      // optional thumbnail
	AtList          []string `json:"atList"`
}

// SendMiniProgramRequest for sending a mini-program card (type=303).
type SendMiniProgramRequest struct {
	TitleList       []string `json:"titleList"`
	ReceivedContent string   `json:"receivedContent"` // title
	Path            string   `json:"path"`            // mini-program page path
	AppID           string   `json:"appId"`           // mini-program appid
	PictureURL      string   `json:"pictureUrl"`      // optional thumbnail
	AtList          []string `json:"atList"`
}

// DissolveGroupRequest for dissolving a group (type=208).
type DissolveGroupRequest struct {
	GroupName string `json:"groupName"` // group to dissolve
}

// Validate checks required fields.
func (r *DissolveGroupRequest) Validate() error {
	if r.GroupName == "" {
		return ErrEmptyGroupName
	}
	return nil
}

// ModifyFriendRequest for updating friend remark/info (type=214).
type ModifyFriendRequest struct {
	Friend FriendUpdate `json:"friend"`
}

// FriendUpdate contains fields to update on a friend.
type FriendUpdate struct {
	Name      string   `json:"name"`      // friend nickname (required)
	MarkName  string   `json:"markName"`  // new remark name
	MarkExtra string   `json:"markExtra"` // new extra remark info
	TagList   []string `json:"tagList"`   // new tags (replaces existing)
}

// Validate checks required fields.
func (r *ModifyFriendRequest) Validate() error {
	if r.Friend.Name == "" {
		return ErrEmptyRecipients
	}
	return nil
}

// AddFriendFromGroupRequest for adding a friend from an external group (type=215).
type AddFriendFromGroupRequest struct {
	GroupName string `json:"groupName"` // external group name
	Nickname  string `json:"nickname"`  // person's nickname in the group
}

// Validate checks required fields.
func (r *AddFriendFromGroupRequest) Validate() error {
	if r.GroupName == "" {
		return ErrEmptyGroupName
	}
	if r.Nickname == "" {
		return ErrEmptyRecipients
	}
	return nil
}

// ModifyGroupMemberRemarkRequest for setting group member remark (type=219).
type ModifyGroupMemberRemarkRequest struct {
	GroupName    string `json:"groupName"`    // target group
	MemberName   string `json:"memberName"`   // member's nickname
	MemberRemark string `json:"memberRemark"` // remark to set for this member
}

// Validate checks required fields.
func (r *ModifyGroupMemberRemarkRequest) Validate() error {
	if r.GroupName == "" {
		return ErrEmptyGroupName
	}
	if r.MemberName == "" {
		return ErrEmptyRecipients
	}
	return nil
}

// InsertCommandRequest for inserting a command ahead of the queue (type=222).
type InsertCommandRequest struct {
	// Command is the command to insert; must be one of the standard command types
	// (e.g. *SendTextRequest, *CreateGroupRequest, etc.).
	Command interface{} `json:"command"`
}

// AddFriendByPhoneResponse wraps APIResponse.
type AddFriendByPhoneResponse = APIResponse

// DeleteContactResponse wraps APIResponse.
type DeleteContactResponse = APIResponse

// AddTodoResponse wraps APIResponse.
type AddTodoResponse = APIResponse

// RecallMessageResponse wraps APIResponse.
type RecallMessageResponse = APIResponse

// SwitchEnterpriseResponse wraps APIResponse.
type SwitchEnterpriseResponse = APIResponse

// CleanupStorageResponse wraps APIResponse.
type CleanupStorageResponse = APIResponse

// ClearCommandsResponse wraps APIResponse.
type ClearCommandsResponse = APIResponse

// ClearSpecificCommandResponse wraps APIResponse.
type ClearSpecificCommandResponse = APIResponse

// BatchSendResponse wraps APIResponse.
type BatchSendResponse = APIResponse

// ============================================================================
// Additional message types (filled in for full API coverage)
// ============================================================================

// SendWeDriveRequest for pushing an image (type=208) or file (type=209) from WeDrive.
// Note: the WorkTool backend distinguishes these two by command type, not by
// payload. Use SendWeDriveImage (208) or SendWeDriveFile (209) on MessageService.
type SendWeDriveRequest struct {
	TitleList  []string `json:"titleList"`
	ObjectName string   `json:"objectName"` // file name as it appears in WeDrive
	ExtraText  string   `json:"extraText"`  // optional comment
}

// Validate checks required fields.
func (r *SendWeDriveRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ObjectName == "" {
		return ErrEmptyObjectName
	}
	return nil
}

// SendDocRequest for pushing a Tencent Doc or collection form (type=211).
type SendDocRequest struct {
	TitleList  []string `json:"titleList"`
	ObjectName string   `json:"objectName"` // doc/form name as it appears in Tencent Docs
	ExtraText  string   `json:"extraText"`
}

// Validate checks required fields.
func (r *SendDocRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ObjectName == "" {
		return ErrEmptyObjectName
	}
	return nil
}

// ForwardMessageRequest for forwarding a message (type=205).
//
// Requires a special "xxx小程序转发群" to be set up beforehand with the robot.
//
// textType: 0=未知 1=文本 2=图片 5=视频 7=小程序 8=链接 9=文件，见 MessageTextType 常量。
type ForwardMessageRequest struct {
	TitleList       []string `json:"titleList"`       // forwarding group name
	ReceivedName    string   `json:"receivedName"`    // nickname of the original sender
	OriginalContent string   `json:"originalContent"` // original content (e.g. mini-program name)
	NameList        []string `json:"nameList"`        // recipients (nicknames or group names)
	ExtraText       string   `json:"extraText"`       // optional comment
	TextType        int      `json:"textType"`        // 见 MessageTextType 常量
}

// Validate checks required fields.
func (r *ForwardMessageRequest) Validate() error {
	if len(r.TitleList) == 0 {
		return ErrEmptyRecipients
	}
	if r.ReceivedName == "" {
		return ErrEmptyFriendName
	}
	if r.OriginalContent == "" {
		return ErrEmptyContent
	}
	if len(r.NameList) == 0 {
		return ErrEmptyForwardRecipients
	}
	return nil
}
