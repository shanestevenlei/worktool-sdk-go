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
- ✅ 机器人配置：信息、在线状态、回调、登录日志、企业列表
- ✅ 历史记录：原始消息、指令回调日志
- ✅ 回调：`callback` 包统一处理 QA 消息回调与事件回调（明文 JSON）
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

    info, _ := c.Robot.GetInfo()
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
    AtList:          []string{types.AtEveryone},
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
    FileType:   string(types.MediaFileTypeAny),
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
    TextType:        int(types.MessageTextTypeMiniProgram),
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
        {Type: int(types.CmdTypeSendText), Payload: &types.SendTextRequest{
            TitleList: []string{"仑哥"}, ReceivedContent: "第一条",
        }},
        {Type: int(types.CmdTypeCreateGroup), Payload: &types.CreateGroupRequest{
            GroupName: "新群", SelectList: []string{"仑哥"},
        }},
        {Type: int(types.CmdTypeSendMedia), Payload: &types.SendImageRequest{
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

// QA 消息回调（用户发消息 → 你的服务回复）
c.Robot.SetQACallback(&types.SetQACallbackRequest{
    OpenCallback: int(types.OpenCallbackEnabled),
    CallbackURL:  "https://your-server.com/qa",
    ReplyAll:     string(types.ReplyAllStrategyEnabled),
})

// 事件回调（指令结果、群二维码、上下线等）
c.Robot.SetEventCallback(&types.SetEventCallbackRequest{
    Type:        int(types.EventCallbackTypeCommandExec),
    CallBackURL: "https://your-server.com/event",
})
c.Robot.ListEventCallbacks(&types.ListEventCallbacksRequest{})
c.Robot.DeleteEventCallback(&types.DeleteEventCallbackRequest{
    Type: int(types.EventCallbackTypeCommandExec),
})

// 其它事件类型示例
c.Robot.SetEventCallback(&types.SetEventCallbackRequest{
    Type:        int(types.EventCallbackTypeOnline),
    CallBackURL: "https://your-server.com/online",
})

// 登录日志 / 企业列表
c.Robot.GetLoginLogs(&types.GetLoginLogsRequest{Date: "2025-01-01"})
c.Robot.GetCorpList(&types.GetCorpListRequest{})
```

## 历史记录

```go
c.History.GetRawMessages(&types.GetRawMessagesRequest{MessageID: "msg_1"})
c.History.GetEventCallbackLog(&types.GetEventCallbackLogRequest{Name: "测试群"})
c.History.GetHistoryMessages(&types.GetHistoryRequest{Title: "仑哥"})
```

## 回调处理

WorkTool 有两种回调协议，均在 `callback` 包中处理（`qa.go` / `event.go`）：

| 类型 | 场景 | 配置 API | 解析 API |
|------|------|----------|----------|
| QA 消息回调 | 用户发消息 → 你回复 | `Robot.SetQACallback` | `callback.ParseQARequest` |
| 事件回调 | 指令结果、上下线等 | `Robot.SetEventCallback` | `callback.NewEventParser` |

### QA 消息回调

WorkTool 将用户消息 POST 到你的 URL，需在 **3 秒内**响应。
文档：[消息回调接口规范](https://doc.worktool.ymdyes.cn/doc-861677.md)

```go
import (
    "io"
    "net/http"

    "github.com/shanestevenlei/worktool-sdk-go/callback"
)

func handleQA(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := io.ReadAll(r.Body)

    msg, err := callback.ParseQARequest(body)
    if err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    // 异步处理：先 QAAck，再调用 Message.SendText 回复
    // 同步回复：QATextReply
    resp := callback.QATextReply("收到：" + msg.Spoken)
    data, _ := callback.MarshalQAResponse(resp)
    w.Header().Set("Content-Type", "application/json")
    w.Write(data)
}
```

### 事件回调

文档：[机器人回调接口标准](https://doc.worktool.ymdyes.cn/api-44952776.md)（type=1 为指令执行结果）

```go
import (
    "fmt"
    "io"
    "net/http"

    "github.com/shanestevenlei/worktool-sdk-go/callback"
)

var parser = callback.NewEventParser()

func handleEvent(w http.ResponseWriter, r *http.Request) {
    defer r.Body.Close()
    body, _ := io.ReadAll(r.Body)

    result, err := parser.Parse(body)
    if err != nil {
        http.Error(w, "bad request", 400)
        return
    }

    if result.IsSuccess() {
        fmt.Printf("指令执行成功: %s\n", result.MessageID)
    } else {
        fmt.Printf("指令失败: %s, code=%d, reason=%s\n",
            result.MessageID, result.ErrorCode, result.ErrorMessage())
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
`ErrEmptyForwardRecipients`, `ErrEmptyCallbackURL`。

## 错误码常量（事件回调）

```go
callback.EventCodeSuccess           // 0
callback.EventCodeIllegalData       // 101011
callback.EventCodeCreateGroupFail   // 201011
callback.EventCodeGroupAddFail      // 201013
callback.EventCodeSendMsgFail       // 201102
callback.EventCodeFileDownload      // 201107

// 中文错误信息
callback.EventErrorCodeMessages[code] // map[int]string
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
├── callback/                      # 回调处理（QA + 事件）
│   ├── doc.go                     # 包说明
│   ├── qa.go                      # QA 消息回调
│   ├── event.go                   # 事件回调
│   ├── qa_test.go
│   └── event_test.go
├── service/
│   ├── service.go                 # HTTPClientFactory + 路径常量
│   ├── message.go                 # 消息相关（25 个方法）
│   ├── message_test.go            # 单元测试
│   ├── robot.go                   # 机器人配置（12 个方法）
│   └── history.go                 # 历史查询（3 个方法）
├── types/
│   ├── message.go                 # 请求/响应类型 + Validate
│   ├── robot.go                   # 机器人类型 + 回调配置
│   ├── eventcallback.go           # 事件回调 payload
│   ├── qacallback.go              # QA 消息回调 payload
│   ├── history.go                 # 查询请求 + 记录类型
│   └── errors.go                  # SDK 错误定义
└── internal/
    └── client/
        ├── client.go              # HTTPClient（每次请求新建）
        └── client_test.go
```

## 注意事项

1. **QPM = 60**（每分钟 60 次）。建议使用 `BatchSend` 合并指令。
2. **回调**：QA 与指令执行回调均为明文 `application/json`。
3. **回调配置**：QA 用 `SetQACallback` + `callback.ParseQARequest`；事件用 `SetEventCallback` + `callback.NewEventParser`。
4. **@所有人**：仅群主或群管理员可触发。
5. **添加好友**：每天 ≤100 人次，新号勿用。

## License

[Apache License 2.0](./LICENSE)
