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
	"context"
	"fmt"
	"os"
	"os/signal"
	"rocketmq-client/utils"
	"sync"
	"syscall"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/rlog"
)

var c rocketmq.PushConsumer
var count int

func main() {
	rlog.SetLogger(XRocketMQLog)
	c, _ = rocketmq.NewPushConsumer(
		consumer.WithGroupName(utils.ConsumerGroup),
		consumer.WithNameServer([]string{utils.RocketMQAddr}),
	)
	err := c.Subscribe(utils.TestTopic, consumer.MessageSelector{}, func(ctx context.Context,
		msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		count += len(msgs)
		fmt.Printf("receive %d msg\n", count)

		//time.Sleep(time.Hour)
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	// Note: start after subscribe
	err = c.Start()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	listenAndServe()
}

func listenAndServe() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		signalChan := make(chan os.Signal)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGQUIT, syscall.SIGHUP)
		for s := range signalChan {
			switch s {
			case syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP:
				err := c.Shutdown()
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
