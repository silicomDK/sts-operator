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

package main

/*
#cgo CFLAGS: -I.
#cgo LDFLAGS: -L. -ltsync -lrt
#include <libtsync.h>
#include <mqueue.h>
#include <string.h>
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/silicomdk/sts-grpc-tsyncd/tsynctl"

	grpc "google.golang.org/grpc"
)

const (
	port = ":50051"
)

var (
	tsc C.struct_tsync_connect
)

// server is used to implement tsynctl.TsynctlServer.
type server struct {
	pb.UnimplementedTsynctlServer
}

// GetMode
func (s *server) GetMode(ctx context.Context, in *pb.Empty) (*pb.ModeReply, error) {
	return &pb.ModeReply{Mode: int32(C.tsync_get_mode(&tsc))}, nil
}

// GetStatus
func (s *server) GetStatus(ctx context.Context, in *pb.Empty) (*pb.StatusReply, error) {
	return &pb.StatusReply{Status: int32(C.tsync_get_status(&tsc))}, nil
}

// GetTime
func (s *server) GetTime(ctx context.Context, in *pb.Empty) (*pb.TimeReply, error) {
	return &pb.TimeReply{Time: int32(C.tsync_get_time(&tsc))}, nil
}

// GetClass
func (s *server) GetClass(ctx context.Context, in *pb.Empty) (*pb.ClassReply, error) {
	return &pb.ClassReply{Class: int32(C.tsync_get_clk_class(&tsc))}, nil
}

// GetClass
func (s *server) StreamAlarms(in *pb.Empty, stream pb.Tsynctl_StreamAlarmsServer) error {
	data := C.alarm_data{}

	// this is blocking call
	for {
		C.tsync_get_alarm(&tsc, &data)
		alarm := &pb.AlarmReply{Type: int32(data._type), Prev: int32(data.prev), Val: int32(data.val)}
		if err := stream.Send(alarm); err != nil {
			return err
		}
	}
}

func connectTsyncd() {
	cs := C.GoStringN(C.CString("gRPC IPC Client 0.0.1\x00"), 20)
	cs2 := C.CString(cs)
	for {
		ret := C.tsync_connect_open(&tsc, cs2)
		fmt.Printf("Opened (%s) ... TX=%d, RX=%d, AL=%d\n", tsc.name, tsc.mqTx, tsc.mqRx, tsc.mqAlr)
		if ret == 0 {
			fmt.Printf("sent msg CONN_REQ, server detected, connecting...\n")

			ret = C.tsync_connect_check(&tsc)
			if ret < 0 {
				time.Sleep(3 * time.Second)
				continue
			}

			if tsc.msgRx.msg_type == C.TSYNC_CONN_RSP {
				fmt.Printf("recv msg CONN_RSP, connected to server\n")
				break
			}
		} else {
			C.tsync_connect_close(&tsc)
			time.Sleep(3 * time.Second)
		}
	}

	fmt.Printf("Success ...TX=%d, RX=%d, AL=%d\n", tsc.mqTx, tsc.mqRx, tsc.mqAlr)
}

func main() {

	fmt.Println("Starting gRPC")
	time.Sleep(5 * time.Second)
	connectTsyncd()

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterTsynctlServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
