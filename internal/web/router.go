package web

import (
	" context"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
)

// SetupRouter 设置路由
func SetupRouter(bot *bot.Bot) *gin.Engine {
	r := gin.Default()

	// 加载模板
	r.LoadHTMLGlob("templates/*")

	// 创建处理器
	s := NewWebHandler(bot)

	// 页面路由
	r.GET("/", h.Index)

	// API 路由
	api := r.Group("/api")
	{
		// 统计
		api.GET("/stats", h.GetStats)

		// 敏感词管理
		api.GET("/sensitive-words", h.GetSensitiveWords)
		api.POST("/sensitive-words", h.AddSensitiveWord)
		api.DELETE("/sensitive-words/:id", h.DeleteSensitiveWord)

		// 定时任务管理
		api.GET("/schedules", h.GetSchedules)
		api.POST("/schedules", h.CreateSchedule)
		api.DELETE("/schedules/:id", h.DeleteSchedule)

		// 消息发送
		api.POST("/send-message", h.SendMessage)
	}

	return r
}
