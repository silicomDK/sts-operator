/*
Copyright 2021.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	stsv1alpha1 "github.com/silicomdk/sts-operator/api/v1alpha1"
	"github.com/silicomdk/sts-operator/controllers"
	//+kubebuilder:scaffold:imports
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(stsv1alpha1.AddToScheme(scheme))
	//+kubebuilder:scaffold:scheme
}

func main() {
	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "ccd94a0b.silicom.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	if err = (&controllers.StsDeviceNodeReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("StsDeviceNode"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StsDeviceNode")
		os.Exit(1)
	}
	if err = (&controllers.StsConfigReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),
		Log:    ctrl.Log.WithName("controllers").WithName("StsConfig"),
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "StsConfig")
		os.Exit(1)
	}
	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	labelNode(mgr.GetClient())

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}

}

func labelNode(client client.Client) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// node-role.sts.silicom.com/master
	// node-role.sts.silicom.com/boundary
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "feature.node.kubernetes.io/pci-0200_8086_1591_1374_02d8.present"})
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
				setupLog.Info(fmt.Sprintf("Removing %s\n", name))
				delete(node.Labels, name)
			}
		}

		//node.Annotations["mode.sts.silicom.com/profile1"] = "enp2s0f0"
		//node.Annotations["mode.sts.silicom.com/profile2"] = "enp2s0f1,enp2s0f02"

		_, err := clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
		if err != nil {
			return
		}

		out, err := exec.Command("lspci", "-n", "-d", "8086:1591").Output()
		if err != nil {
			setupLog.Error(err, "Error with lspci")
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
				setupLog.Error(err, "Error with ip link")
				return
			}

			if strings.Contains(string(out), "state UP") {
				node.Labels[fmt.Sprintf("iface.sts.silicom.com/%s", iface)] = "up"
			} else {
				node.Labels[fmt.Sprintf("iface.sts.silicom.com/%s", iface)] = "down"
			}

			setupLog.Info(fmt.Sprintf("Adding %s\n", iface))
		}

		_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
		if err != nil {
			return
		}
	}
}
