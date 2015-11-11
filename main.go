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
)

var addr = flag.String("addr", "localhost:8000", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go connection(&wg)
	}

	wg.Wait()

}

func connection(wg *sync.WaitGroup) {

	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/socket.io/?EIO=3&transport=websocket"}
	u := "ws://localhost:8000/socket.io/?EIO=3&transport=websocket"
	//u := url.URL{Scheme: "ws", Host: *addr, Path: "/socket.io/1/?EIO=2&transport=polling&t=1446710999349"}
	log.Printf("connecting to %s", u)

	c, _, err := websocket.DefaultDialer.Dial(u, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				wg.Done()
				break
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		//err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
		//req	__NSCFString *	@"42[\"login\",{\"client_id\":\"1\",\"token\":\"a433709204fadb3f97acb386cf3bb1b9\",\"client_type\":\"user\"}]"	0x00007fcc78eac2f0
		log.Println(t.String())
		//err := c.WriteMessage(websocket.TextMessage, []byte("42[\"message\",\"3:2\"]"))
		err := c.WriteMessage(websocket.TextMessage, []byte("42[\"login\",{\"client_id\":\"1\",\"token\":\"a433709204fadb3f97acb386cf3bb1b9\",\"client_type\":\"user\"}]"))
		if err != nil {
			log.Println("write:", err)
			wg.Done()
			break
		}
	}
}
