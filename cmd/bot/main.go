package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
	"github.com/BoyQcb/fangzhang-bot/internal/handler"
	"github.com/BoyQcb/fangzhang-bot/internal/middleware"
	"github.com/BoyQcb/fangzhang-bot/internal/store"
	"github.com/BoyQcb/fangzhang-bot/internal/web"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Bot struct {
		Token string `yaml:"token"`
		Debug bool   `yaml:"debug"`
	} `yaml:"bot"`
	Database struct {
		Driver string `yaml:"driver"`
		DSN    string `yaml:"dsn"`
	} `yaml:"database"`
	Admin struct {
		SuperUsers []int64 `yaml:"super_users"`
	} `yaml:"admin"`
}

func main() {
	// 加载配置
	config := loadConfig()

	// 初始化数据库
	db := store.InitDB(config.Database.Driver, config.Database.DSN)
	defer db.Close()

	// 创建bot
	ctx := context.Background()
	opts := []bot.Option{
		bot.WithDefaultHandler(handler.DefaultHandler),
	}

	if config.Bot.Debug {
		opts = append(opts, bot.WithDebug())
	}

	b, err := bot.New(config.Bot.Token, opts...)
	if err != nil {
		log.Fatalf("创建bot失败: %v", err)
	}

	// 注册中间件
	b.RegisterMiddleware(middleware.Logger)
	b.RegisterMiddleware(middleware.Auth(config.Admin.SuperUsers))

	// 注册处理器
	registerHandlers(b, db)

	// 启动Web后台
	go startWebServer(b)

	// 启动bot
	log.Println("Bot 启动成功...")
	b.Start(ctx)

	// 优雅关闭
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("正在关闭 bot...")
	b.Close()
}

// startWebServer 启动Web后台服务器
func startWebServer(b *bot.Bot) {
	gin.SetMode(gin.ReleaseMode)
	r := web.SetupRouter(b)

	log.Println("Web 后台启动在 :8080...")
	if err := r.Run(":8080"); err != nil {
		log.Printf("Web 服务器启动失败: %v", err)
	}
}

func loadConfig() *Config {
	data, err := os.ReadFile("config/config.yaml")
	if err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		log.Fatalf("解析配置文件失败: %v", err)
	}

	return &config
}

func registerHandlers(b *bot.Bot, db *store.DB) {
	// 消息管理
	b.RegisterHandler(bot.HandlerTypeMessage, "/delete", bot.HandlerFunc(handler.DeleteMessage))
	b.RegisterHandler(bot.HandlerTypeMessage, "/edit", bot.HandlerFunc(handler.EditMessage))
	b.RegisterHandler(bot.HandlerTypeMessage, "/forward", bot.HandlerFunc(handler.ForwardMessage))

	// 群组管理
	b.RegisterHandler(bot.HandlerTypeMessage, "/kick", bot.HandlerFunc(handler.KickUser))
	b.RegisterHandler(bot.HandlerTypeMessage, "/mute", bot.HandlerFunc(handler.MuteUser))
	b.RegisterHandler(bot.HandlerTypeMessage, "/unmute", bot.HandlerFunc(handler.UnmuteUser))
	b.RegisterHandler(bot.HandlerTypeMessage, "/promote", bot.HandlerFunc(handler.PromoteAdmin))

	// 内容过滤
	b.RegisterHandler(bot.HandlerTypeMessage, "/addword", bot.HandlerFunc(handler.AddSensitiveWord))
	b.RegisterHandler(bot.HandlerTypeMessage, "/delword", bot.HandlerFunc(handler.DeleteSensitiveWord))
	b.RegisterHandler(bot.HandlerTypeMessage, "/listwords", bot.HandlerFunc(handler.ListSensitiveWords))

	// 定时任务
	b.RegisterHandler(bot.HandlerTypeMessage, "/addschedule", bot.HandlerFunc(handler.AddSchedule))
	b.RegisterHandler(bot.HandlerTypeMessage, "/delschedule", bot.HandlerFunc(handler.DeleteSchedule))
	b.RegisterHandler(bot.HandlerTypeMessage, "/listschedules", bot.HandlerFunc(handler.ListSchedules))

	// 数据统计
	b.RegisterHandler(bot.HandlerTypeMessage, "/stats", bot.HandlerFunc(handler.ShowStats))
	b.RegisterHandler(bot.HandlerTypeMessage, "/topusers", bot.HandlerFunc(handler.TopUsers))

	// 抽奖
	b.RegisterHandler(bot.HandlerTypeMessage, "/lottery", bot.HandlerFunc(handler.StartLottery))
	b.RegisterHandler(bot.HandlerTypeMessage, "/draw", bot.HandlerFunc(handler.DrawLottery))
}