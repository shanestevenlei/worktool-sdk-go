package callback

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"testing"

	"github.com/shanestevenlei/worktool-sdk-go/types"
	"github.com/stretchr/testify/assert"
)

// pkcs7Pad pads the plaintext to AES block size using PKCS7.
func pkcs7Pad(src []byte, blockSize int) []byte {
	pad := blockSize - len(src)%blockSize
	return append(src, bytes.Repeat([]byte{byte(pad)}, pad)...)
}

func TestParser_Parse_NoEncryption(t *testing.T) {
	p := NewParser("")
	cb, err := p.Parse([]byte(`{"messageId":"m1","errorCode":0,"errorReason":""}`))
	assert.NoError(t, err)
	assert.Equal(t, "m1", cb.MessageID)
	assert.True(t, cb.IsSuccess())
	assert.Equal(t, "success", cb.ErrorMessage())
}

func TestParser_Parse_WithEncryption(t *testing.T) {
	// Use a non-hex key so decryptAES falls back to using the raw bytes.
	key := "this-is-a-32-byte-secret-key!!!!"
	p := NewParser(key)

	plaintext := pkcs7Pad([]byte(`{"messageId":"m2","errorCode":201102,"errorReason":"send fail"}`), aes.BlockSize)
	block, err := aes.NewCipher([]byte(key))
	assert.NoError(t, err)
	iv := make([]byte, aes.BlockSize) // zero-IV matches WorkTool convention
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	encrypted := append(iv, ciphertext...)

	cb, err := p.Parse(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, "m2", cb.MessageID)
	assert.Equal(t, 201102, cb.ErrorCode)
	assert.False(t, cb.IsSuccess())
	assert.Contains(t, cb.ErrorMessage(), "发送消息失败")
}

func TestParser_Parse_InvalidJSON(t *testing.T) {
	p := NewParser("")
	_, err := p.Parse([]byte(`not json`))
	assert.Error(t, err)
}

func TestParser_Parse_ShortCiphertext(t *testing.T) {
	p := NewParser("this-is-a-32-byte-secret-key!!!!")
	_, err := p.Parse([]byte("short"))
	assert.Error(t, err)
}

func TestParser_Parse_BadPadding(t *testing.T) {
	key := "this-is-a-32-byte-secret-key!!!!"
	p := NewParser(key)

	// 32 bytes of plaintext = 2 cipher blocks, but last byte is 0 → invalid padding.
	block, _ := aes.NewCipher([]byte(key))
	iv := make([]byte, aes.BlockSize)
	plaintext := make([]byte, aes.BlockSize*2) // all zeros
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	_, err := p.Parse(append(iv, ciphertext...))
	assert.Error(t, err)
}

func TestParseFromRequest(t *testing.T) {
	cb, err := ParseFromRequest([]byte(`{"messageId":"x","errorCode":0}`), "")
	assert.NoError(t, err)
	assert.Equal(t, "x", cb.MessageID)
	assert.True(t, cb.IsSuccess())
}

func TestCallback_ErrorMessage_Unknown(t *testing.T) {
	cb := &Callback{
		CallbackRequest: types.CallbackRequest{
			MessageID:   "m",
			ErrorCode:   999999,
			ErrorReason: "custom reason",
		},
	}
	assert.Equal(t, "custom reason", cb.ErrorMessage())
}

func TestErrorCodeMessages_HasAllCodes(t *testing.T) {
	for code, want := range map[int]string{
		CodeSuccess:         "success",
		CodeCreateGroupFail: "创建群失败",
		CodeGroupAddFail:    "群拉人失败",
		CodeSendMsgFail:     "发送消息失败",
	} {
		got, ok := ErrorCodeMessages[code]
		assert.True(t, ok, "missing code %d", code)
		assert.Equal(t, want, got)
	}
}

func TestNewParser(t *testing.T) {
	p := NewParser("mykey")
	assert.NotNil(t, p)
	assert.Equal(t, "mykey", p.secretKey)
}
