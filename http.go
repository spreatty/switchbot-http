package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

var bot = &Bot{address: botAddress}

type BotReply struct {
	Code int
}

type ErrReply struct {
	Error error
}

func startHTTP() {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.Default()

	engine.POST("/press", func(c *gin.Context) {
		err := bot.Open()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		botReplyCode, err := bot.Press()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		c.JSON(200, BotReply{Code: botReplyCode})
	})

	engine.POST("/on", func(c *gin.Context) {
		err := bot.Open()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		botReplyCode, err := bot.On()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		c.JSON(200, BotReply{Code: botReplyCode})
	})

	engine.POST("/off", func(c *gin.Context) {
		err := bot.Open()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		botReplyCode, err := bot.Off()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		c.JSON(200, BotReply{Code: botReplyCode})
	})

	engine.POST("/open", func(c *gin.Context) {
		err := bot.Open()
		if err != nil {
			c.JSON(500, ErrReply{Error: err})
		}
		bot.StartKeepAlive()
		c.Status(200)
	})

	engine.POST("/close", func(c *gin.Context) {
		bot.StopKeepAlive()
		c.Status(200)
	})

	err := engine.Run(":" + fmt.Sprint(config.HttpPort))
	if err != nil {
		log.Fatalln("Failed to start server", err)
	}
}
