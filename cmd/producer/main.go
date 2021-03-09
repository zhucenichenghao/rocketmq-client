/*
Licensed to the Apache Software Foundation (ASF) under one or more
contributor license agreements.  See the NOTICE file distributed with
this work for additional information regarding copyright ownership.
The ASF licenses this file to You under the Apache License, Version 2.0
(the "License"); you may not use this file except in compliance with
the License.  You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"github.com/gin-gonic/gin"
)

var p rocketmq.Producer
var r *gin.Engine

// Package main implements a simple producer to send message.
func main() {
	p, _ = rocketmq.NewProducer(
		producer.WithNameServer([]string{"192.168.11.26:9876"}),
		producer.WithRetry(2),
	)
	err := p.Start()
	if err != nil {
		fmt.Printf("start producer error: %s", err.Error())
		os.Exit(1)
	}
	listenAndServe()
}

func listenAndServe() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		r := gin.Default()
		r.POST("/sendMsg", handleReceiveMsg)
		err := r.Run(":8000")
		if err != nil {
			fmt.Printf("gin run error: %s\n", err.Error())
			return
		}
		wg.Done()
	}()
	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		for s := range signalChan {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
				err := p.Shutdown()
				if err != nil {
					fmt.Printf("mq producer shutdown error: %s\n", err.Error())
					continue
				}
				fmt.Printf("mq producer shutdown successful\n")
			}
			break
		}
		wg.Done()
	}()
	wg.Wait()
}
