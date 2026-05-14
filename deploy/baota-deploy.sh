#!/bin/bash
# 宝塔面板部署脚本 - 方丈机器人
# 使用方法: bash baota-deploy.sh

set -e

echo "=========================================="
echo "  方丈机器人 - 宝塔面板部署脚本"
echo "=========================================="
echo ""

# 检查是否在宝塔环境
if [ ! -d "/www/server" ]; then
    echo "⚠️  警告: 未检测到宝塔环境，是否继续? (y/n)"
    read -r response
    if [ "$response" != "y" ]; then
        echo "部署已取消"
        exit 1
    fi
fi

# 1. 安装 Go 环境
echo "1️⃣  检查 Go 环境..."
if ! command -v go &> /dev/null; then
    echo "   Go 未安装，正在安装..."
    # 下载并安装 Go 1.21
    cd /tmp
    wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
    tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    source ~/.bashrc
    rm go1.21.6.linux-amd64.tar.gz
    echo "   ✅ Go 安装完成"
else
    echo "   ✅ Go 已安装: $(go version)"
fi

# 2. 创建项目目录
echo ""
echo "2️⃣  创建项目目录..."
DEPLOY_DIR="/www/wwwroot/fangzhang-bot"
mkdir -p $DEPLOY_DIR
cd $DEPLOY_DIR
echo "   ✅ 项目目录: $DEPLOY_DIR"

# 3. 复制项目文件（假设项目已在当前目录）
echo ""
echo "3️⃣  部署项目文件..."
# 如果是在本地开发后上传，可以跳过此步骤
# 这里假设你已经通过宝塔上传了项目文件

if [ ! -f "go.mod" ]; then
    echo "   ⚠️  未找到 go.mod，请确保项目文件已上传到 $DEPLOY_DIR"
    echo "   你可以通过宝塔面板的文件管理上传项目压缩包"
    exit 1
fi

# 4. 下载依赖
echo ""
echo "4️⃣  下载项目依赖..."
go mod download
echo "   ✅ 依赖下载完成"

# 5. 编译项目
echo ""
echo "5️⃣  编译项目..."
go build -o fangzhang-bot cmd/bot/main.go
chmod +x fangzhang-bot
echo "   ✅ 编译完成: $DEPLOY_DIR/fangzhang-bot"

# 6. 创建数据目录
echo ""
echo "6️⃣  创建数据目录..."
mkdir -p data
chmod 755 data
echo "   ✅ 数据目录创建完成"

# 7. 配置向导
echo ""
echo "7️⃣  配置 Bot..."
if [ ! -f "config/config.yaml" ]; then
    mkdir -p config
    cat > config/config.yaml << EOF
bot:
  token: "YOUR_BOT_TOKEN_HERE"
  debug: false

database:
  driver: "sqlite3"
  dsn: "./data/bot.db"

admin:
  super_users:
    - 123456789
EOF
    echo "   ✅ 配置文件已创建: config/config.yaml"
    echo ""
    echo "   ⚠️  请编辑配置文件，填入你的 Bot Token:"
    echo "      vi config/config.yaml"
    echo ""
    echo "   获取 Bot Token:"
    echo "     1. 在 Telegram 中搜索 @BotFather"
    echo "     2. 发送 /newbot 创建机器人"
    echo "     3. 复制获得的 Token"
    echo ""
    echo "   获取你的 Telegram ID:"
    echo "     1. 在 Telegram 中搜索 @userinfobot"
    echo "     2. 发送任意消息，它会返回你的 ID"
else
    echo "   ✅ 配置文件已存在"
fi

# 8. 创建 systemd 服务
echo ""
echo "8️⃣  创建系统服务..."
SERVICE_FILE="/etc/systemd/system/fangzhang-bot.service"

sudo tee $SERVICE_FILE > /dev/null << EOF
[Unit]
Description=Fangzhang Telegram Bot
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=$DEPLOY_DIR
ExecStart=$DEPLOY_DIR/fangzhang-bot
Restart=always
RestartSec=10
StandardOutput=append:/www/wwwlogs/fangzhang-bot.log
StandardError=append:/www/wwwlogs/fangzhang-bot.error.log

[Install]
WantedBy=multi-user.target
EOF

sudo systemctl daemon-reload
echo "   ✅ 系统服务已创建"

# 9. 创建日志文件
echo ""
echo "9️⃣  创建日志文件..."
touch /www/wwwlogs/fangzhang-bot.log
touch /www/wwwlogs/fangzhang-bot.error.log
chmod 755 /www/wwwlogs/fangzhang-bot*.log
echo "   ✅ 日志文件已创建"

# 10. 防火墙提示
echo ""
echo "🔟  配置防火墙..."
echo "   如果你的服务器启用了防火墙，请开放端口 8080 (Web 后台):"
echo "   firewall-cmd --add-port=8080/tcp --permanent"
echo "   firewall-cmd --reload"
echo ""
echo "   或者在宝塔面板中："
echo "   1. 进入【安全】页面"
echo "   2. 添加端口 8080"

# 完成
echo ""
echo "=========================================="
echo "  ✅ 部署完成！"
echo "=========================================="
echo ""
echo "📝  后续步骤:"
echo "   1. 编辑配置文件: vi config/config.yaml"
echo "   2. 启动服务: systemctl start fangzhang-bot"
echo "   3. 设置开机自启: systemctl enable fangzhang-bot"
echo "   4. 查看状态: systemctl status fangzhang-bot"
echo "   5. 查看日志: tail -f /www/wwwlogs/fangzhang-bot.log"
echo ""
echo "🌐  Web 后台访问:"
echo "   http://你的服务器IP:8080"
echo ""
echo "📋  常用命令:"
echo "   启动:   systemctl start fangzhang-bot"
echo "   停止:   systemctl stop fangzhang-bot"
echo "   重启:   systemctl restart fangzhang-bot"
echo "   状态:   systemctl status fangzhang-bot"
echo "   日志:   journalctl -u fangzhang-bot -f"
echo ""
