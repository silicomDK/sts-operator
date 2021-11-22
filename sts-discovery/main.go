// SPDX-License-Identifier: Apache-2.0
// Copyright (c) 2020-2021 Intel Corporation

package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func labelNode() {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := clientset.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	node, err2 := clientset.CoreV1().Nodes().Get(context.Background(), os.Getenv("NODE_NAME"), metav1.GetOptions{})
	if err2 != nil {
		panic(err2.Error())
		return
	}

	for name, _ := range node.Labels {
		if strings.Contains(name, "iface.sts.silicom.com") {
			fmt.Printf("Removing %s\n", name)
			delete(node.Labels, name)
		}
	}

	_, err = clientset.CoreV1().Nodes().Update(context.Background(), node, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
		return
	}

	out, err := exec.Command("lspci", "-n", "-d", "8086:1591").Output()
	if err != nil {
		fmt.Errorf("error with lspci %v", err)
		return
	}

	scanner := bufio.NewScanner(strings.NewReader(string(out)))
	for scanner.Scan() {
		id := strings.Split(scanner.Text(), " ")
		path := fmt.Sprintf("/sys/bus/pci/devices/0000:%s/net/*", id[0])

		res, _ := filepath.Glob(path)
		iface := string(filepath.Base(res[0]))

		out, err := exec.Command("ip", "link", "show", "dev", iface).Output()
		if err != nil {
			fmt.Errorf("Error with ip link %v", err)
			return
		}

		if strings.Contains(string(out), "state UP") {
			node.Labels[fmt.Sprintf("iface.sts.silicom.com/%s", iface)] = "up"
		} else {
			node.Labels[fmt.Sprintf("iface.sts.silicom.com/%s", iface)] = "down"
		}

		fmt.Printf("Adding %s\n", iface)
	}

	_, err = clientset.CoreV1().Nodes().Update(context.TODO(), node, metav1.UpdateOptions{})
	if err != nil {
		return
	}
}

func main() {
	labelNode()
	fmt.Printf("Accelerator discovery finished successfully\n")
	for {
		time.Sleep(time.Duration(1000000000) * time.Second)
	}
	os.Exit(0)
}
