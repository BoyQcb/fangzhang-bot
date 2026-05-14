# Fangzhang Bot (方丈机器人复刻版)

一个功能完整的 Telegram 机器人，完全复刻方丈机器人的功能，使用 Go 语言开发。

## ✨ 功能特性

### 1. 消息管理
- `/delete <message_id>` - 删除消息
- `/edit <message_id> <新内容>` - 编辑消息
- `/forward <来源聊天ID> <消息ID>` - 转发消息

### 2. 群组管理
- `/kick <user_id>` - 踢人（可回复消息）
- `/mute <user_id> [时长(秒)]` - 禁言用户
- `/unmute <user_id>` - 解除禁言
- `/promote <user_id>` - 设置管理员

### 3. 内容过滤
- `/addword <敏感词>` - 添加敏感词
- `/delword <敏感词>` - 删除敏感词
- `/listwords` - 列出所有敏感词
- 自动过滤包含敏感词的消息

### 4. 定时任务
- `/addschedule <cron表达式> <消息内容>` - 添加定时任务
- `/delschedule <任务ID>` - 删除定时任务
- `/listschedules` - 列出所有定时任务

### 5. 数据统计
- `/stats` - 显示消息统计
- `/topusers` - 显示活跃用户排行（最近7天）

### 6. 抽奖
- `/lottery <奖品名称> [持续时间(秒)]` - 开始抽奖
- `/join` - 参与抽奖（回复抽奖消息）
- `/draw <抽奖ID>` - 手动开奖（管理员）

## 🚀 安装部署

### 环境要求
- Go 1.21+
- SQLite3

### 1. 克隆项目
```bash
git clone https://github.com/xxx/fangzhang-bot.git
cd fangzhang-bot
```

### 2. 安装依赖
```bash
go mod download
```

### 3. 配置 Bot
1. 在 Telegram 中找到 [@BotFather](https://t.me/BotFather)
2. 发送 `/newbot` 创建新机器人
3. 按提示设置机器人名称和用户名
4. 复制获得的 **Bot Token**

### 4. 修改配置文件
编辑 `config/config.yaml`：
```yaml
bot:
  token: "YOUR_BOT_TOKEN_HERE"  # 替换为你的 Bot Token
  debug: false

database:
  driver: "sqlite3"
  dsn: "./data/bot.db"

admin:
  super_users:
    - 123456789  # 替换为你的 Telegram 用户 ID
```

**获取你的 Telegram ID：**
- 在 Telegram 中搜索 [@userinfobot](https://t.me/userinfobot)
- 发送任意消息，它会返回你的 ID

### 5. 编译运行
```bash
# 编译
go build -o fangzhang-bot cmd/bot/main.go

# 运行
./fangzhang-bot
```

## 🐳 Docker 部署

### 1. 创建 Dockerfile
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o bot cmd/bot/main.go

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/bot .
COPY config/ config/
CMD ["./bot"]
```

### 2. 构建镜像
```bash
docker build -t fangzhang-bot .
```

### 3. 运行容器
```bash
docker run -d \
  --name fangzhang-bot \
  -v $(pwd)/config:/app/config \
  -v $(pwd)/data:/app/data \
  fangzhang-bot
```

## 🔧 系统服务部署（Linux）

### 1. 创建 systemd 服务文件
创建 `/etc/systemd/system/fangzhang-bot.service`：
```ini
[Unit]
Description=Fangzhang Telegram Bot
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/fangzhang-bot
ExecStart=/opt/fangzhang-bot/fangzhang-bot
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
```

### 2. 启用服务
```bash
sudo systemctl daemon-reload
sudo systemctl enable fangzhang-bot
sudo systemctl start fangzhang-bot
sudo systemctl status fangzhang-bot
```

## 📝 使用说明

### 命令列表

#### 消息管理（需要管理员权限）
| 命令 | 说明 | 示例 |
|------|------|------|
| `/delete` | 删除消息 | `/delete 123` |
| `/edit` | 编辑消息 | `/edit 123 新内容` |
| `/forward` | 转发消息 | `/forward -100123456 789` |

#### 群组管理（需要管理员权限）
| 命令 | 说明 | 示例 |
|------|------|------|
| `/kick` | 踢人 | `/kick 123456` 或回复消息 |
| `/mute` | 禁言 | `/mute 123456 300` |
| `/unmute` | 解除禁言 | `/unmute 123456` |
| `/promote` | 设置管理员 | `/promote 123456` |

#### 内容过滤（需要管理员权限）
| 命令 | 说明 | 示例 |
|------|------|------|
| `/addword` | 添加敏感词 | `/addword 敏感词` |
| `/delword` | 删除敏感词 | `/delword 敏感词` |
| `/listwords` | 列出敏感词 | `/listwords` |

#### 定时任务（需要管理员权限）
| 命令 | 说明 | 示例 |
|------|------|------|
| `/addschedule` | 添加定时任务 | `/addschedule "0 9 * * *" "早安"` |
| `/delschedule` | 删除定时任务 | `/delschedule 1` |
| `/listschedules` | 列出定时任务 | `/listschedules` |

**Cron 表达式示例：**
- `"0 9 * * *"` - 每天早上 9 点
- `"0 */2 * * *"` - 每 2 小时
- `"0 0 * * 0"` - 每周日零点

#### 数据统计
| 命令 | 说明 |
|------|------|
| `/stats` | 显示消息统计 |
| `/topusers` | 显示活跃用户排行 |

#### 抽奖
| 命令 | 说明 | 示例 |
|------|------|------|
| `/lottery` | 开始抽奖 | `/lottery "iPhone 15" 300` |
| `/join` | 参与抽奖 | 回复抽奖消息并输入 `/join` |
| `/draw` | 手动开奖 | `/draw 1`（需要管理员权限）|

## 🛠 开发指南

### 项目结构
```
fangzhang-bot/
├── cmd/bot/main.go          # 程序入口
├── internal/
│   ├── handler/             # 消息处理器
│   ├── middleware/          # 中间件
│   ├── model/              # 数据模型
│   ├── store/              # 数据存储
│   └── service/            # 业务逻辑
├── config/config.yaml      # 配置文件
├── go.mod
└── README.md
```

### 添加新功能
1. 在 `internal/handler/` 中添加处理器函数
2. 在 `cmd/bot/main.go` 中注册处理器
3. 如需要，在 `internal/model/` 中添加数据模型
4. 在 `internal/store/` 中添加数据库操作
5. 在 `internal/service/` 中添加业务逻辑

### 测试
```bash
# 运行测试
go test ./...

# 运行特定测试
go test -v ./internal/handler
```

## 📋 TODO
- [ ] 完善权限验证（使用 Telegram API 查询管理员状态）
- [ ] 添加 Redis 支持（用于频率限制和缓存）
- [ ] 添加 Web UI 管理界面
- [ ] 支持更多数据库（PostgreSQL、MySQL）
- [ ] 添加单元测试
- [ ] 添加 CI/CD 配置
- [ ] 支持多语言
- [ ] 添加日志轮转
- [ ] 支持 Docker Compose 部署

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License

## 🙏 致谢

- [go-telegram/bot](https://github.com/go-telegram/bot) - Telegram Bot API 库
- [robfig/cron](https://github.com/robfig/cron) - 定时任务库
- [gorm](https://gorm.io/) - ORM 库

## 📧 联系方式

如有问题，请提交 Issue 或联系开发者。

---

**⚠️ 免责声明：**
本程序仅供学习和研究使用，请遵守 Telegram 的使用条款和相关规定。
