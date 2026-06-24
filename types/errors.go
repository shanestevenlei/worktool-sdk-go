package types

import "errors"

var (
	ErrEmptyRecipients       = errors.New("worktool: recipients (titleList) cannot be empty")
	ErrEmptyContent          = errors.New("worktool: content (receivedContent) cannot be empty")
	ErrEmptyObjectName       = errors.New("worktool: objectName cannot be empty")
	ErrEmptyFileURL          = errors.New("worktool: fileUrl cannot be empty")
	ErrEmptyGroupName        = errors.New("worktool: groupName cannot be empty")
	ErrEmptyPhone            = errors.New("worktool: phone cannot be empty")
	ErrEmptyMessageID        = errors.New("worktool: messageId cannot be empty")
	ErrEmptyEnterpriseName   = errors.New("worktool: enterpriseName cannot be empty")
	ErrEmptyCommandList      = errors.New("worktool: command list cannot be empty")
	ErrEmptyFriendName       = errors.New("worktool: friend nickname cannot be empty")
	ErrEmptyForwardRecipients = errors.New("worktool: forward nameList cannot be empty")
	ErrEmptyCallbackURL      = errors.New("worktool: callbackUrl cannot be empty")
)