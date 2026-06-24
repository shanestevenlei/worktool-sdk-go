# WorkTool SDK for Go

Go SDK for the [WorkTool WeChat Enterprise API](https://worktool.apifox.cn/).

## 功能特性

- ✅ 发送文本 / 图片 / 音视频 / 文件（结构化请求体）
- ✅ 推送微盘图片 / 微盘文件 / 腾讯文档 / 收集表
- ✅ 转发消息、自定义链接、小程序卡片
- ✅ 群管理：创建、修改、踢人、拉人、修改备注、修改群成员备注、设置群模板
- ✅ 好友管理：按手机号添加、从群添加、修改、删除、修改群成员备注
- ✅ 消息生命周期：撤回、插队、批量、清空、清空指定
- ✅ 切换企业、清理存储、添加待办
- ✅ 机器人配置：信息、在线状态、加密、回调（旧/新）、登录日志、企业列表
- ✅ 历史记录：原始消息、QA 日志、回调日志
- ✅ AES-256-CBC 解密回调（自动）
- ✅ 请求参数校验（Validate）
- ✅ **Client 无状态**：每次请求构建独立 HTTPClient，多 goroutine 并发安全

## 安装

```bash
go get github.com/shanestevenlei/worktool-sdk-go
```

## 快速开始

```go
package main

import (
	"fmt"

	worktool "github.com/shanestevenlei/worktool-sdk-go"
	"github.com/shanestevenlei/worktool-sdk-go/types"
)

func main() {
    c := worktool.NewClient("your_robot_id")

    // 发送文本消息
    resp, err := c.Message.SendText(&types.SendTextRequest{
        TitleList:       []string{"仑哥"},
        ReceivedContent: "你好~",
    })
    if err != nil {
        panic(err)
    }
    fmt.Printf("code=%d, message=%s\n", resp.Code, resp.Message)

    // 带加密
    c2 := worktool.NewClientWithSecret("robot_id", "16bytekey1234567", 1)
    info, _ := c2.Robot.GetInfo()
    fmt.Println(info.Data.Name)
}
```

## 架构设计

```
worktool.Client               (stateless facade)
   │
   ├── .Message ──► MessageService ──► HTTPClientFactory
   ├── .Robot   ──► RobotService   ──► HTTPClientFactory  (same pointer as Client)
   └── .History ──► HistoryService ──► HTTPClientFactory  (same pointer as Client)
                                              │
                                              ▼
                                    Client.HTTPClient()      // called per request
                                              │
                                              ▼
                                    client.New(Config{...})   // fresh per request
```

关键设计：

- **`worktool.Client` 不持有任何 HTTP 状态**。每个 service 方法通过注入的 `HTTPClientFactory` 在调用时构建全新的 `HTTPClient`，确保多 goroutine 并发安全。
- **service 是无状态的**。多个 goroutine 共享同一 service 不会冲突。
- **测试通过 `http.RoundTripper` mock 注入**（参见 `Config.HTTPDoer`）。

## 消息发送

### 文本 / 富文本

```go
c.Message.SendText(&types.SendTextRequest{
    TitleList:       []string{"仑哥"},
    ReceivedContent: "你好",
    AtList:          []string{"@所有人"},
})
```

### 图片 / 文件 / 音视频（type=218）

```go
// 网络图片
c.Message.SendImage(&types.SendImageRequest{
    TitleList:  []string{"仑哥"},
    ObjectName: "logo.png",
    FileURL:    "https://example.com/logo.png",
})

// 网络文件
c.Message.SendFile(&types.SendFileRequest{
    TitleList:  []string{"仑哥"},
    ObjectName: "report.pdf",
    FileURL:    "https://example.com/report.pdf",
    FileType:   "*", // image / audio / video / *
})

// 微盘文件 / 图片
c.Message.SendWeDriveImage(&types.SendWeDriveRequest{
    TitleList:  []string{"仑哥"},
    ObjectName: "微盘里的图片名",
})
c.Message.SendWeDriveFile(&types.SendWeDriveRequest{
    TitleList:  []string{"仑哥"},
    ObjectName: "微盘里的文件名",
})

// 腾讯文档 / 收集表（已废弃但保留）
c.Message.SendTencentDoc(&types.SendDocRequest{
    TitleList:  []string{"仑哥"},
    ObjectName: "文档名",
})

// 转发消息（需要先创建「小程序转发群」）
c.Message.ForwardMessage(&types.ForwardMessageRequest{
    TitleList:       []string{"转发群名"},
    ReceivedName:    "原始发送人昵称",
    OriginalContent: "原始内容",
    NameList:        []string{"仑哥"},
    TextType:        7, // 0=未知 1=文本 2=图片 5=视频 7=小程序 8=链接 9=文件
})

// 自定义链接/小程序（付费）
c.Message.SendLink(&types.SendLinkRequest{
    TitleList:       []string{"仑哥"},
    ReceivedContent: "标题",
    LinkURL:         "https://example.com",
    PictureURL:      "https://example.com/icon.png",
})
```

### 群管理

```go
// 创建群
c.Message.CreateGroup(&types.CreateGroupRequest{
    GroupName:  "测试群01",
    SelectList: []string{"仑哥", "小明"},
    GroupAnnouncement: "欢迎大家",
})

// 修改群
c.Message.UpdateGroup(&types.UpdateGroupRequest{
    GroupName:    "测试群01",
    NewGroupName: "测试群02",
    SelectList:   []string{"小王"},
    RemoveList:   []string{"小明"},
})

// 解散群
c.Message.DissolveGroup(&types.DissolveGroupRequest{GroupName: "测试群01"})

// 修改群成员备注
c.Message.ModifyGroupMemberRemark(&types.ModifyGroupMemberRemarkRequest{
    GroupName:    "测试群01",
    MemberName:   "张三",
    MemberRemark: "项目经理",
})
```

### 好友管理

```go
c.Message.AddFriendByPhone(&types.AddFriendByPhoneRequest{Phone: "13800138000"})
c.Message.AddFriendFromGroup(&types.AddFriendFromGroupRequest{
    GroupName: "外部群",
    Nickname:  "张三",
})
c.Message.ModifyFriend(&types.ModifyFriendRequest{
    Friend: types.FriendUpdate{Name: "仑哥", MarkName: "VIP客户"},
})
c.Message.DeleteContact(&types.DeleteContactRequest{TitleList: []string{"张三"}})
```

### 消息生命周期

```go
c.Message.RecallMessage(&types.RecallMessageRequest{MessageID: "msg_1"})
c.Message.InsertCommand(&types.InsertCommandRequest{Command: &types.SendTextRequest{...}})
c.Message.ClearCommands(&types.ClearCommandsRequest{})
c.Message.ClearSpecificCommand(&types.ClearSpecificCommandRequest{MessageID: "msg_1"})
c.Message.AddTodo(&types.AddTodoRequest{...})
c.Message.SwitchEnterprise(&types.SwitchEnterpriseRequest{EnterpriseName: "新企业"})
c.Message.CleanupStorage(&types.CleanupStorageRequest{})
```

### 批量发送（合并请求，节省 QPM）

```go
c.Message.BatchSend(&types.BatchSendRequest{
    List: []types.BatchItem{
        {Type: 203, Payload: &types.SendTextRequest{
            TitleList: []string{"仑哥"}, ReceivedContent: "第一条",
        }},
        {Type: 206, Payload: &types.CreateGroupRequest{
            GroupName: "新群", SelectList: []string{"仑哥"},
        }},
        {Type: 218, Payload: &types.SendImageRequest{
            TitleList: []string{"仑哥"}, ObjectName: "x.png", FileURL: "https://x.png",
        }},
    },
})
```

## 机器人配置

```go
// 基础
c.Robot.GetInfo()
c.Robot.IsOnline()
c.Robot.SetEncryption(&types.SetEncryptionRequest{SecretKey: "k", EncryptType: 1})

// 消息回调
c.Robot.SetCallback(&types.SetCallbackRequest{
    OpenCallback: 1,
    CallbackURL:  "https://your-server.com/cb",
    ReplyAll:     "1",
})

// 事件回调（新）
c.Robot.BindCallback(&types.BindCallbackRequest{
    Type:        callback.CallbackTypeCommandExec,
    CallBackURL: "https://your-server.com/exec",
})
c.Robot.ListCallbacks(&types.ListCallbacksRequest{})
c.Robot.DeleteCallback(&types.DeleteCallbackRequest{Type: 1})

// 登录日志 / 企业列表
c.Robot.GetLoginLogs(&types.GetLoginLogsRequest{Date: "2025-01-01"})
c.Robot.GetCorpList(&types.GetCorpListRequest{})
```

## 历史记录

```go
c.History.GetRawMessages(&types.GetRawMessagesRequest{MessageID: "msg_1"})
c.History.GetQALog(&types.GetQALogRequest{Name: "测试群"})
c.History.GetHistoryMessages(&types.GetHistoryRequest{Title: "仑哥"})
```

## 回调处理

```go
package main

import (
    "fmt"

    "github.com/shanestevenlei/worktool-sdk-go/callback"
)

var parser = callback.NewParser("your_secret_key_if_encrypted")

func handle(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := io.ReadAll(r.Body)

    cb, err := parser.Parse(body)
    if err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    if cb.IsSuccess() {
        fmt.Printf("指令执行成功: %s\n", cb.MessageID)
    } else {
        fmt.Printf("指令失败: %s, code=%d, reason=%s\n",
            cb.MessageID, cb.ErrorCode, cb.ErrorMessage())
    }

    w.WriteHeader(200)
}
```

## 请求校验

每个 request 结构体都有 `Validate()`：

```go
_, err := c.Message.SendText(&types.SendTextRequest{TitleList: []string{}})
// err == types.ErrEmptyRecipients
```

可用错误：`ErrEmptyRecipients`, `ErrEmptyContent`, `ErrEmptyObjectName`,
`ErrEmptyFileURL`, `ErrEmptyGroupName`, `ErrEmptyPhone`, `ErrEmptyMessageID`,
`ErrEmptyEnterpriseName`, `ErrEmptyCommandList`, `ErrEmptyFriendName`,
`ErrEmptyForwardRecipients`, `ErrEmptyCallbackURL`, `ErrInvalidEncryptType`。

## 错误码常量（callback 包）

```go
callback.CodeSuccess           // 0
callback.CodeIllegalData       // 101011
callback.CodeCreateGroupFail   // 201011
callback.CodeGroupAddFail      // 201013
callback.CodeSendMsgFail       // 201102
callback.CodeFileDownload      // 201107

// 中文错误信息
callback.ErrorCodeMessages[code] // map[int]string
```

## 测试

```bash
go test ./...
```

所有 service 通过 `http.RoundTripper` 注入 mock，零网络依赖：

```go
import (
    "github.com/shanestevenlei/worktool-sdk-go"
    "github.com/shanestevenlei/worktool-sdk-go/internal/client"
)

// mockRoundTripper implements client.HTTPDoer
c := worktool.New(worktool.Config{
    RobotID:  "robot_test",
    HTTPDoer: myMockRoundTripper,
})
```

## 项目结构

```
worktool-sdk-go/
├── client.go                      # SDK 入口（无状态）
├── callback/
│   ├── callback.go                # 解析 + AES 解密 + 错误码
│   └── callback_test.go
├── service/
│   ├── service.go                 # HTTPClientFactory + 路径常量
│   ├── message.go                 # 消息相关（25 个方法）
│   ├── message_test.go            # 单元测试
│   ├── robot.go                   # 机器人配置（12 个方法）
│   └── history.go                 # 历史查询（3 个方法）
├── types/
│   ├── message.go                 # 请求/响应类型 + Validate
│   ├── robot.go                   # 机器人类型 + Callback 类型
│   ├── history.go                 # 查询请求 + 记录类型
│   └── errors.go                  # SDK 错误定义
└── internal/
    └── client/
        ├── client.go              # HTTPClient（每次请求新建）
        └── client_test.go
```

## 注意事项

1. **QPM = 60**（每分钟 60 次）。建议使用 `BatchSend` 合并指令。
2. **AES 加密**：robot 配置 `encryptType=1` 后，body 需 AES-256-CBC 加密（zero-IV, PKCS7）。
3. **回调**：优先使用消息回调（`SetCallback`）而非轮询历史消息。
4. **@所有人**：仅群主或群管理员可触发。
5. **添加好友**：每天 ≤100 人次，新号勿用。

## License

[Apache License 2.0](./LICENSE)
