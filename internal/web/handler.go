package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-telegram/bot"
)

// WebHandler Web 处理器
type WebHandler struct {
	bot *bot.Bot
}

// NewWebHandler 创建 Web 处理器
func NewWebHandler(b *bot.Bot) *WebHandler {
	return &WebHandler{bot: b}
}

// Index 首页
func (h *WebHandler) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "方丈机器人管理后台",
	})
}

// GetStats 获取统计信息
func (h *WebHandler) GetStats(c *gin.Context) {
	// 统计用户数
	userCount, _ := store.CountUsers()

	// 统计消息数
	var messageCount int64
	store.GetDB().Model(&model.Message{}).Count(&messageCount)

	// 统计定时任务数（简化）
	scheduleCount := 0

	// 统计抽奖数
	var lotteryCount int64
	store.GetDB().Model(&model.Lottery{}).Count(&lotteryCount)

	c.JSON(http.StatusOK, gin.H{
		"user_count":    userCount,
		"message_count": messageCount,
		"schedule_count": scheduleCount,
		"lottery_count": lotteryCount,
	})
}

// GetSensitiveWords 获取敏感词列表
func (h *WebHandler) GetSensitiveWords(c *gin.Context) {
	var words []model.SensitiveWord
	store.GetDB().Find(&words)
	c.JSON(http.StatusOK, words)
}

// AddSensitiveWord 添加敏感词
func (h *WebHandler) AddSensitiveWord(c *gin.Context) {
	var req struct {
		Word string `json:"word"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	word := &model.SensitiveWord{Word: req.Word}
	if err := store.CreateSensitiveWord(word); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "添加成功"})
}

// DeleteSensitiveWord 删除敏感词
func (h *WebHandler) DeleteSensitiveWord(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	if err := store.DeleteSensitiveWord(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetSchedules 获取定时任务列表
func (h *WebHandler) GetSchedules(c *gin.Context) {
	schedules, err := store.ListAllSchedules()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

// CreateSchedule 创建定时任务
func (h *WebHandler) CreateSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "功能暂时禁用"})
}

// DeleteSchedule 删除定时任务
func (h *WebHandler) DeleteSchedule(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "功能暂时禁用"})
}

// SendMessage 发送消息
func (h *WebHandler) SendMessage(c *gin.Context) {
	var req struct {
		ChatID  int64  `json:"chat_id"`
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if h.bot == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Bot 实例未初始化"})
		return
	}

	// 调用 Telegram Bot API 发送消息
	_, err := h.bot.SendMessage(c.Request.Context(), &bot.SendMessageParams{
		ChatID: req.ChatID,
		Text:   req.Message,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "消息发送成功",
		"chat_id": req.ChatID,
	})
}
