/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a client for Greeter service.
package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	pb "client/grpc/helloworld"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultName = "world"
)

var (
	addr = flag.String("addr", "localhost:50051", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("did not connect: ", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// --------普通请求方式------- //
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: *name})
	if err != nil {
		fmt.Println("could not greet: ", err)
	}
	fmt.Println("Greeting: ", r.GetMessage())

	// --------双向流式请求方式------- //
	bstream, err := c.BidirectionalStreamingSayHello(ctx)
	if err != nil {
		fmt.Println("could not open stream: ", err)
		return
	}
	go func() {
		// 异步发送信息
		for i := 0; i < 10; i++ {
			bstream.Send(&pb.HelloRequest{Name: "bidirection stream client"})
		}
		bstream.CloseSend()
	}()
	fmt.Printf("response:\n")
	for {
		r, err := bstream.Recv()
		if err != nil {
			break
		}
		fmt.Printf(" - %s\n", r.Message)
	}
}
