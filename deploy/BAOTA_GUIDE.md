# 宝塔面板部署指南 - 方丈机器人

本文档详细介绍如何将方丈机器人部署到宝塔面板环境。

## 📋 前置要求

- 已安装宝塔面板 (www.bt.cn)
- 服务器系统: CentOS 7+ / Ubuntu 18.04+
- 至少 512MB 内存
- 已开放防火墙端口（如需 Web 后台）

## 🚀 部署步骤

### 方法一：使用自动部署脚本（推荐）

1. **下载部署脚本**

在宝塔面板的终端中执行：

```bash
cd /www/wwwroot
wget https://your-repo-url/deploy/baota-deploy.sh
bash baota-deploy.sh
```

或者，如果你已经手动上传了项目文件：

```bash
cd /www/wwwroot/fangzhang-bot
bash deploy/baota-deploy.sh
```

2. **按提示操作**

脚本会自动完成：
- 安装 Go 环境
- 下载项目依赖
- 编译项目
- 创建 systemd 服务
- 创建日志文件

3. **编辑配置文件**

```bash
vi /www/wwwroot/fangzhang-bot/config/config.yaml
```

填入你的 Bot Token 和管理员 ID（参考下文）。

4. **启动服务**

```bash
systemctl start fangzhang-bot
systemctl enable fangzhang-bot
```

---

### 方法二：手动部署（适合定制化需求）

#### 步骤 1: 安装 Go 环境

在宝塔面板的【终端】中执行：

```bash
# 下载 Go 1.21
cd /tmp
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz

# 解压到 /usr/local
tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz

# 配置环境变量
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc

# 验证安装
go version
```

#### 步骤 2: 创建项目目录

```bash
mkdir -p /www/wwwroot/fangzhang-bot
cd /www/wwwroot/fangzhang-bot
```

#### 步骤 3: 上传项目文件

**方法 A: 通过宝塔文件管理上传**

1. 在宝塔面板中点击【文件】
2. 进入 `/www/wwwroot/fangzhang-bot`
3. 上传项目压缩包（zip 或 tar.gz）
4. 右键解压

**方法 B: 使用 Git 克隆**

```bash
cd /www/wwwroot/fangzhang-bot
git clone https://github.com/xxx/fangzhang-bot.git .
```

**方法 C: 从本地上传**

1. 在本地压缩项目文件夹
2. 通过宝塔【文件】上传
3. 解压到 `/www/wwwroot/fangzhang-bot`

#### 步骤 4: 下载依赖并编译

```bash
cd /www/wwwroot/fangzhang-bot
go mod download
go build -o fangzhang-bot cmd/bot/main.go
chmod +x fangzhang-bot
```

#### 步骤 5: 配置 Bot

创建/编辑配置文件：

```bash
mkdir -p config
vi config/config.yaml
```

填入以下内容（**替换成你自己的值**）：

```yaml
bot:
  token: "123456789:ABCdefGHIjklMNOpqrsTUVwxyz"  # 从 @BotFather 获取
  debug: false

database:
  driver: "sqlite3"
  dsn: "./data/bot.db"

admin:
  super_users:
    - 123456789  # 你的 Telegram ID（从 @userinfobot 获取）
```

**获取 Bot Token:**

1. 在 Telegram 中搜索 `@BotFather`
2. 发送 `/newbot`
3. 按提示设置机器人名称和用户名
4. 复制获得的 Token（格式：`123456789:ABCdef...`）

**获取你的 Telegram ID:**

1. 在 Telegram 中搜索 `@userinfobot`
2. 发送任意消息
3. 它会返回你的 ID（一串数字）

#### 步骤 6: 创建 systemd 服务

创建服务文件：

```bash
vi /etc/systemd/system/fangzhang-bot.service
```

填入以下内容：

```ini
[Unit]
Description=Fangzhang Telegram Bot
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/www/wwwroot/fangzhang-bot
ExecStart=/www/wwwroot/fangzhang-bot/fangzhang-bot
Restart=always
RestartSec=10
StandardOutput=append:/www/wwwlogs/fangzhang-bot.log
StandardError=append:/www/wwwlogs/fangzhang-bot.error.log

[Install]
WantedBy=multi-user.target
```

重新加载 systemd：

```bash
systemctl daemon-reload
```

#### 步骤 7: 启动服务

```bash
# 启动服务
systemctl start fangzhang-bot

# 设置开机自启
systemctl enable fangzhang-bot

# 查看状态
systemctl status fangzhang-bot
```

#### 步骤 8: 配置防火墙

**方法 A: 通过宝塔面板**

1. 进入宝塔面板【安全】页面
2. 在【防火墙】部分，添加端口 `8080`（Web 后台使用）
3. 协议选择 `TCP`，备注填写 `fangzhang-bot-web`

**方法 B: 命令行**

```bash
# firewalld
firewall-cmd --add-port=8080/tcp --permanent
firewall-cmd --reload

# 或 ufw
ufw allow 8080/tcp
ufw reload
```

---

## 🌐 访问 Web 后台

部署成功后，在浏览器中访问：

```
http://你的服务器IP:8080
```

**功能：**

- 📊 统计面板 - 查看用户数、消息数、定时任务数
- 🔧 敏感词管理 - 添加/删除敏感词
- ⏰ 定时任务管理 - 添加/删除定时任务
- 📨 消息发送 - 通过 Web 界面发送消息

---

## 📋 常用命令

```bash
# 启动服务
systemctl start fangzhang-bot

# 停止服务
systemctl stop fangzhang-bot

# 重启服务
systemctl restart fangzhang-bot

# 查看状态
systemctl status fangzhang-bot

# 查看日志（实时）
journalctl -u fangzhang-bot -f

# 查看日志文件
tail -f /www/wwwlogs/fangzhang-bot.log
```

---

## 🔧 故障排除

### 1. Bot 无响应

**检查：**
- Token 是否正确
- 服务是否运行：`systemctl status fangzhang-bot`
- 查看日志：`tail -f /www/wwwlogs/fangzhang-bot.log`

### 2. Web 后台无法访问

**检查：**
- 端口 8080 是否开放：`netstat -tulpn | grep 8080`
- 防火墙规则：`firewall-cmd --list-ports`
- 查看错误日志：`tail -f /www/wwwlogs/fangzhang-bot.error.log`

### 3. 编译失败

**可能原因：**
- Go 版本过低（需要 1.21+）
- 依赖下载失败（网络问题）
- 代码错误

**解决：**
```bash
# 更新 Go
rm -rf /usr/local/go
# 重新下载并安装最新版

# 清理并重新下载依赖
go clean -modcache
go mod download
```

### 4. 数据库错误

**检查：**
- 数据目录是否存在：`ls -la data/`
- 权限是否正确：`chmod 755 data`
- 手动初始化数据库：运行 bot 会自动创建

---

## 📝 更新项目

当项目有更新时：

```bash
cd /www/wwwroot/fangzhang-bot

# 停止服务
systemctl stop fangzhang-bot

# 备份配置文件
cp config/config.yaml /tmp/config.yaml.bak

# 更新代码（根据你的部署方式）
git pull
# 或重新上传文件

# 重新编译
go mod download
go build -o fangzhang-bot cmd/bot/main.go

# 恢复配置文件
cp /tmp/config.yaml.bak config/config.yaml

# 启动服务
systemctl start fangzhang-bot
```

---

## 🔒 安全建议

1. **修改 Web 后台端口**

编辑 `cmd/bot/main.go`，修改这一行：

```go
if err := r.Run(":8080"); err != nil {
```

改成其他端口，比如 `:9090`。

2. **添加 Web 后台认证**

在 `internal/web/router.go` 中添加登录验证中间件。

3. **限制 Web 后台访问 IP**

在宝塔面板【安全】中，仅允许信任的 IP 访问端口 8080。

4. **定期备份数据库**

```bash
# 每天凌晨 2 点备份
crontab -e
0 2 * * * cp /www/wwwroot/fangzhang-bot/data/bot.db /backup/bot.db.$(date +\%Y\%m\%d)
```

---

## 📞 技术支持

如遇到问题，请：
1. 查看日志文件
2. 检查配置文件
3. 提交 Issue 到项目仓库
4. 或联系开发者

---

**🎉 部署完成！祝你使用愉快！**
