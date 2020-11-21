package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"time"
)

type Service struct {
	TelegramBotApi *tgbotapi.BotAPI
}

func main() {
	s := &Service{}
	s.TelegramBotApi = getTelegramBotApiConn()

	go s.GetMessageFromTelegram()

	msg := tgbotapi.NewMessage(-494747213, "service started")
	if _, err := s.TelegramBotApi.Send(msg); err != nil {
		log.Println(err)
	}

	router := gin.New()
	router.Use(gin.Recovery())

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("[%s] - %s %s path=%s, status_code=%d, latency=%s, user_agent=%s, error_message=%s\n",
			param.TimeStamp.Format(time.RFC1123),
			param.ClientIP,
			param.Method,
			param.Path,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	router.GET("/", s.Test)

	log.Printf("start web server on %s", os.Getenv("PORT"))

	if err := router.Run(":" + os.Getenv("PORT")); err != nil {
		log.Println("router run :"+os.Getenv("PORT"), err)
	}
}

func (s *Service) GetMessageFromTelegram() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := s.TelegramBotApi.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		// check Message, group, user
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Chat == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		if _, err := s.TelegramBotApi.Send(msg); err != nil {
			log.Println(err)
		}
	}
}

func (s *Service) Test(c *gin.Context) {
	text := fmt.Sprintf("%s - %s", os.Getenv("TEST_TEXT"), uuid.New().String())

	files := []string{
		"AgACAgIAAxkBAAIHT1-1qFEYD4ZrVnhRZVKZjGNFwne9AAIEsTEbMOKwSbh5zGI2DbgssSJzly4AAwEAAwIAA20AA9hCAwABHgQ",
		"AgACAgIAAxkBAAIHT1-1qFEYD4ZrVnhRZVKZjGNFwne9AAIEsTEbMOKwSbh5zGI2DbgssSJzly4AAwEAAwIAA3gAA9dCAwABHgQ",
	}

	publishDate := time.Now().UTC()
	messageText := fmt.Sprintf("``` %s ``` %s", publishDate.Format("01.02.2006 15:04:05 MST"), "New Test text")
	var medias []interface{}
	for i, file := range files {
		f := tgbotapi.NewInputMediaPhoto(file)

		if i == 0 {
			f.ParseMode = "markdown"
			f.Caption = messageText
		}

		medias = append(medias, f)
	}

	msg := tgbotapi.NewMediaGroup(-494747213, medias)

	if _, err := s.TelegramBotApi.Send(msg); err != nil {
		log.Println(err)
	}

	c.String(http.StatusOK, text)
}

func getTelegramBotApiConn() *tgbotapi.BotAPI {
	bot := &tgbotapi.BotAPI{}
	var err error
	for {
		bot, err = tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_SECRET"))
		if err != nil {
			log.Println("new bot api connection error:", err)
			continue
		}

		bot.Debug = false
		break
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return bot
}
