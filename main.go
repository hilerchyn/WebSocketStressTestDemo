// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ignore

package main

import (
	"flag"
	"log"
	_ "net/url"
	"time"

	"github.com/gorilla/websocket"
	"sync"
	"strings"
)

var addr = flag.String("addr", "localhost:8000", "服务器连接地址")
var con = flag.Int("con", 40, "并发总数")
var interval = flag.Int("interval", 1, "重连的时间间隔,以秒为单位")
var content = flag.String("content", "42[\"login\",{\"client_id\":\"1\",\"token\":\"a433709204fadb3f97acb386cf3bb1b9\",\"client_type\":\"user\"}]", "向服务器发送的数据")

var tickets chan int

func main() {

	flag.Parse()
	log.SetFlags(0)


	tickets = make(chan int, *con)

	var wg sync.WaitGroup

	for {

		select {
		case tickets<-1:
			wg.Add(1)
			log.Println("tickets capacity:", cap(tickets))
			log.Println("create connection")
			go connection(&wg)
		case <-time.After(time.Duration(100) * time.Millisecond):

		}

	}

	wg.Wait()

}

func connection(wg *sync.WaitGroup) {

	ticker := time.NewTicker(time.Duration(*interval) * time.Second)
	defer ticker.Stop()

	//u := url.URL{Scheme: "ws", Host: *addr, Path: /socket.io/?EIO=3&transport=websocket"}
	u := "ws://"+*addr+"/socket.io/?EIO=3&transport=websocket"
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/socket.io/1/?EIO=2&transport=polling&t=1446710999349"}
	log.Printf("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Println("dial:", err)

		ticker.Stop()
		c.Close()
		wg.Done()
		<-tickets
		return
	}
	defer c.Close()

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)

				ticker.Stop()
				c.Close()
				wg.Done()
				<-tickets
				return
			}
			//log.Printf("recv: %s", message)
			if strings.Contains(string(message), "login_resp") {

				ticker.Stop()
				c.Close()
				wg.Done()
				<-tickets

				log.Println("done")
				break
			}
		}
	}()



	for t := range ticker.C {
		//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//req	__NSCFString *	@"42[\"login\",{\"client_id\":\"1\",\"token\":\"a433709204fadb3f97acb386cf3bb1b9\",\"client_type\":\"user\"}]"	0x00007fcc78eac2f0
		t.String()
		//err := c.WriteMessage(websocket.TextMessage, []byte("42[\"message\",\"3:2\"]"))
		err := c.WriteMessage(websocket.TextMessage, []byte(*content))
		if err != nil {
			log.Println("write:", err)
			wg.Done()
			break
		}
	}
}
