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

	// node-role.sts.silicom.com/master
	// node-role.sts.silicom.com/boundary
	nodes, err := clientset.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{LabelSelector: "feature.node.kubernetes.io/pci-0200_8086_1591_1374_02d8.present"})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("There are %d nodes in the cluster with STS2 card\n", len(nodes.Items))

	for _, node := range nodes.Items {
		if node.Name != os.Getenv("NODE_NAME") {
			continue
		}
		for name, _ := range node.Labels {
			if strings.Contains(name, "iface.sts.silicom.com") {
				fmt.Printf("Removing %s\n", name)
				delete(node.Labels, name)
			}
		}

		//node.Annotations["mode.sts.silicom.com/profile1"] = "enp2s0f0"
		//node.Annotations["mode.sts.silicom.com/profile2"] = "enp2s0f1,enp2s0f02"

		_, err := clientset.CoreV1().Nodes().Update(context.Background(), &node, metav1.UpdateOptions{})
		if err != nil {
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

		_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
		if err != nil {
			return
		}
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
