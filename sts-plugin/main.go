// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020-2021 Intel Corporation

package main

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	stsv1alpha1 "github.com/silicomdk/sts-operator/api/v1alpha1"
	pb "github.com/silicomdk/sts-operator/grpc/tsynctl"
	grpc "google.golang.org/grpc"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type GPSStatusRsp struct {
	Tpvs []TPV `json:"tpv"`
}

type TPV struct {
	Time string  `json:"time"`
	Lat  float32 `json:"lat"`
	Lon  float32 `json:"lon"`
}

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(stsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func query_host(stsNode *stsv1alpha1.StsNode) {
	var idx int

	out, err := exec.Command("lspci", "-n", "-d", "8086:1591").Output()
	if err != nil {
		fmt.Errorf("error with lspci %v", err)
		return
	}

	stsNode.Status.EthInterfaces = []stsv1alpha1.StsNodeInterfaceStatus{}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		var nodeInterface stsv1alpha1.StsNodeInterfaceStatus

		id := strings.Split(scanner.Text(), " ")
		path := fmt.Sprintf("/sys/bus/pci/devices/0000:%s/net/*", id[0])
		res, _ := filepath.Glob(path)
		iface := string(filepath.Base(res[0]))

		nodeInterface.EthName = iface
		nodeInterface.PciAddr = id[0]
		nodeInterface.EthPort = idx
		idx = idx + 1

		out, err := exec.Command("ip", "link", "show", "dev", iface).Output()
		if err != nil {
			fmt.Printf("Error with ip link %v\n", err)
			return
		}

		if strings.Contains(string(out), "state UP") {
			nodeInterface.Status = "up"
		} else {
			nodeInterface.Status = "down"
		}

		stsNode.Status.EthInterfaces = append(stsNode.Status.EthInterfaces, nodeInterface)
	}
}

func query_tsyncd(svc_str string, stsNode *stsv1alpha1.StsNode) {

	conn, err := grpc.Dial(svc_str, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		fmt.Printf("Could not connect: %v\n", err)
		return
	}
	defer conn.Close()

	gRpcClient := pb.NewTsynctlGrpcClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	message, err := gRpcClient.GetStatus(ctx, &pb.Empty{})
	if err != nil {
		fmt.Printf("could not get status: %v\n", err)
	} else {
		stsNode.Status.TsyncStatus.Status = message.GetMessage()
	}

	message, err = gRpcClient.GetMode(ctx, &pb.Empty{})
	if err != nil {
		fmt.Printf("could not get mode: %v\n", err)
	} else {
		stsNode.Status.TsyncStatus.Mode = message.GetMessage()
	}

	timeReply, err := gRpcClient.GetTime(ctx, &pb.Empty{})
	if err != nil {
		fmt.Printf("could not get time: %v\n", err)
	} else {
		stsNode.Status.TsyncStatus.Time = timeReply.GetMessage()

	}
	cancel()
}

func query_gpsd(svc_str string, stsNode *stsv1alpha1.StsNode) {
	var status GPSStatusRsp

	conn, err := net.Dial("tcp", svc_str)
	if err != nil {
		fmt.Println(fmt.Sprintf("Dial failed: %s: %v\n", svc_str, err))
		return
	}
	defer conn.Close()

	fmt.Fprintf(conn, "?POLL;")
	rsp, _ := bufio.NewReader(conn).ReadString('\n')
	if len(rsp) < 1 {
		fmt.Printf("Bad GPS Read: %s\n", rsp)
		return
	}

	err = json.Unmarshal([]byte(rsp), &status)
	if err != nil {
		fmt.Println("Error occured during unmarshaling.")
	}
}

func main() {
	stsNode := &stsv1alpha1.StsNode{}

	nodeName := os.Getenv("NODE_NAME")
	namespace := os.Getenv("NAMESPACE")

	tsyncPort, _ := strconv.Atoi(os.Getenv("TSYNC_PORT"))
	tsync_svc_str := fmt.Sprintf("%s:%d", nodeName, tsyncPort)

	gpsPort, _ := strconv.Atoi(os.Getenv("GPS_PORT"))
	gpsd_svc_str := fmt.Sprintf("%s:%d", nodeName, gpsPort)

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		panic(err.Error())
	}

	for {
		err = k8sClient.Get(context.Background(),
			client.ObjectKey{
				Namespace: namespace,
				Name:      nodeName,
			}, stsNode)
		if err != nil {
			fmt.Println("Can't get stsnode yet, waiting... 5 seconds\n")
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}

	for {
		query_host(stsNode)
		query_tsyncd(tsync_svc_str, stsNode)
		query_gpsd(gpsd_svc_str, stsNode)

		if err := k8sClient.Status().Update(context.TODO(), stsNode); err != nil {
			fmt.Printf("Update failed: %v\n", err)
		}
		time.Sleep(30 * time.Second)
	}
	os.Exit(0)
}