package callback

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
	"encoding/json"

	"github.com/shanestevenlei/worktool-sdk-go/types"
)

// Parser decrypts and decodes callback bodies.
//
// WorkTool uses AES-256-CBC with zero-IV and PKCS7 padding. The key is
// interpreted as either a hex string or a raw byte string and zero-padded
// to 32 bytes.
type Parser struct {
	secretKey string
}

// NewParser creates a callback parser.
// secretKey: the AES key configured on the robot (empty = no decryption).
func NewParser(secretKey string) *Parser {
	return &Parser{secretKey: secretKey}
}

// Parse decrypts (if needed) and decodes a callback payload.
func (p *Parser) Parse(data []byte) (*Callback, error) {
	if p.secretKey != "" {
		decrypted, err := p.decryptAES(data)
		if err != nil {
			return nil, err
		}
		data = decrypted
	}

	var raw types.CallbackRequest
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	return &Callback{CallbackRequest: raw}, nil
}

// ParseFromRequest is a convenience wrapper for net/http handlers.
func ParseFromRequest(data []byte, secretKey string) (*Callback, error) {
	return NewParser(secretKey).Parse(data)
}

// Callback wraps types.CallbackRequest with helper methods.
type Callback struct {
	types.CallbackRequest
}

// IsSuccess returns true if the callback indicates a successful execution.
func (c *Callback) IsSuccess() bool {
	return c.ErrorCode == CodeSuccess
}

// ErrorMessage returns a human-readable message for ErrorCode, or "" if unknown.
func (c *Callback) ErrorMessage() string {
	if msg, ok := ErrorCodeMessages[c.ErrorCode]; ok {
		return msg
	}
	return c.ErrorReason
}

func (p *Parser) decryptAES(data []byte) ([]byte, error) {
	key, err := hex.DecodeString(p.secretKey)
	if err != nil {
		key = []byte(p.secretKey)
	}
	if len(key) < 32 {
		padded := make([]byte, 32)
		copy(padded, key)
		key = padded
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(data) < aes.BlockSize {
		return nil, errShortCiphertext
	}

	iv := data[:aes.BlockSize]
	ciphertext := data[aes.BlockSize:]
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, errCiphertextAlignment
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertext, ciphertext)

	pad := ciphertext[len(ciphertext)-1]
	n := int(pad)
	if n == 0 || n > len(ciphertext) || n > aes.BlockSize {
		return nil, errPadding
	}
	return ciphertext[:len(ciphertext)-n], nil
}

var (
	errShortCiphertext     = &aesError{"ciphertext too short"}
	errCiphertextAlignment = &aesError{"ciphertext is not a multiple of the block size"}
	errPadding             = &aesError{"invalid PKCS7 padding"}
)

type aesError struct{ msg string }

func (e *aesError) Error() string { return "worktool/callback: " + e.msg }

// =============================================================================
// Error code constants and human-readable messages
// =============================================================================

const (
	CodeSuccess           = 0
	CodeIllegalData       = 101011
	CodeIllegalOperation  = 101012
	CodeIllegalPermission = 101013

	CodeCreateGroupFail   = 201011
	CodeGroupRenameFail   = 201012
	CodeGroupAddFail      = 201013
	CodeGroupRemoveFail   = 201014
	CodeGroupAnnounceFail = 201015
	CodeGroupRemarkFail   = 201016
	CodeIntoRoomFail      = 201101
	CodeSendMsgFail       = 201102
	CodeButtonFail        = 201103
	CodeTargetFail        = 201104
	CodeRelayFail         = 201105
	CodeRepeat            = 201106
	CodeFileDownload      = 201107
	CodeFileStorage       = 201108
)

// ErrorCodeMessages maps a callback error code to a human-readable description.
var ErrorCodeMessages = map[int]string{
	CodeSuccess:           "success",
	CodeIllegalData:       "非法数据",
	CodeIllegalOperation:  "非法操作",
	CodeIllegalPermission: "非法权限",
	CodeCreateGroupFail:   "创建群失败",
	CodeGroupRenameFail:   "修改群名失败",
	CodeGroupAddFail:      "群拉人失败",
	CodeGroupRemoveFail:   "群踢人失败",
	CodeGroupAnnounceFail: "修改群公告失败",
	CodeGroupRemarkFail:   "修改群备注失败",
	CodeIntoRoomFail:      "进群失败",
	CodeSendMsgFail:       "发送消息失败",
	CodeButtonFail:        "按钮失败",
	CodeTargetFail:        "目标失败",
	CodeRelayFail:         "转发失败",
	CodeRepeat:            "重复执行",
	CodeFileDownload:      "文件下载失败",
	CodeFileStorage:       "文件存储失败",
}
