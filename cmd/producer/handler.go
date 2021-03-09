package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/gin-gonic/gin"
)

type HandleMsgReq struct {
	Body  string `json:"body,omitempty" bson:"body,omitempty"`
	Topic string `json:"topic,omitempty" bson:"topic,omitempty"`
}

func handleReceiveMsg(ctx *gin.Context) {
	req := new(HandleMsgReq)
	err := ctx.ShouldBindJSON(req)
	if err != nil {
		fmt.Printf("handle receive error: %v\n", err)
		return
	}

	msg := &primitive.Message{
		Topic: req.Topic,
		Body:  []byte(req.Body),
	}
	for i := 0; i < 1; i++ {
		res, err := p.SendSync(context.Background(), msg)
		if err != nil {
			fmt.Printf("send message error: %s\n", err)
		} else {
			fmt.Printf("send message success: result=%s\n", res.String())
		}
	}

	ctx.String(http.StatusOK, "Send OK")
}
