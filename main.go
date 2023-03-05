package main

import (
	"context"
	"github.com/gin-gonic/gin"
	openai "github.com/sashabaranov/go-openai"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	openAIClient *openai.Client
	log          *logrus.Logger
)

func main() {
	initLog()
	openAIClient = openai.NewClient("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")
	r := gin.Default()
	r.LoadHTMLGlob("tmpl/*")
	r.Static("/static", "static")
	r.GET("/chat", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	r.GET("/chat/get", func(c *gin.Context) {
		msg := c.Query("msg")
		if strings.TrimSpace(msg) == "" {
			c.String(http.StatusOK, "empty input")
			return
		}
		log.Info(msg)
		rsp, err := genResponse(msg)
		if err != nil {
			log.Error(err.Error())
			rsp = "服务器繁忙，请稍后再试"
		}
		rsp = strings.TrimLeft(rsp, "\n")
		c.String(http.StatusOK, rsp)
	})
	r.Run()
}
func genResponse(msg string) (string, error) {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*20)
	resp, err := openAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: msg,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}
func initLog() {
	// 写入文件
	f, err := os.OpenFile("./log.txt", os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	// 实例化
	log = logrus.New()
	// 设置输出
	log.Out = io.MultiWriter(f, os.Stdout)
	// 设置日志级别
	log.SetLevel(logrus.DebugLevel)
}
