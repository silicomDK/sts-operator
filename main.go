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
	"bytes"
	"context"
	"flag"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	nfdv1 "github.com/openshift/cluster-nfd-operator/api/v1"
	stsv1alpha1 "github.com/silicomdk/sts-operator/api/v1alpha1"
	"github.com/silicomdk/sts-operator/controllers"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

type StsPlugin struct {
	Namespace     string
	ImageRegistry string
	TsyncPort     int
	GpsPort       int
	StsVersion    string
	Version       string
}

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(stsv1alpha1.AddToScheme(scheme))
	utilruntime.Must(nfdv1.AddToScheme(scheme))
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
		LeaderElectionID:       "sts.silicom.com",
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
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

	deployPlugin()
	deployNfd()

	setupLog.Info("starting StsConfig manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func deployNfd() error {

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		panic(err.Error())
	}

	nfdOperand := &nfdv1.NodeFeatureDiscovery{}
	nfdOperand.Name = "nfd-sts-silicom"
	nfdOperand.Namespace = "sts-silicom"
	nfdOperand.Spec.Operand.Namespace = "sts-silicom"

	content, err := ioutil.ReadFile("/assets/nfd-discovery.yaml")
	if err != nil {
		fmt.Println("ERROR: Loading nfd-discovery.yaml file")
		return err
	}

	workerConfig := &nfdv1.ConfigMap{}
	workerConfig.ConfigData = string(content)
	nfdOperand.Spec.WorkerConfig = workerConfig

	if err := k8sClient.Get(context.TODO(), client.ObjectKey{
		Namespace: nfdOperand.Namespace,
		Name:      nfdOperand.Name,
	}, nfdOperand); err != nil {

		err = k8sClient.Create(context.TODO(), nfdOperand)
		if err != nil {
			panic(err)
		}
	} else {
		k8sClient.Update(context.TODO(), nfdOperand)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func deployPlugin() error {
	var buff bytes.Buffer
	var objects []client.Object

	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	k8sClient, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		panic(err.Error())
	}

	configPlugin := &StsPlugin{}
	configPlugin.Namespace = os.Getenv("NAMESPACE")
	configPlugin.ImageRegistry = "quay.io/silicom"
	configPlugin.Version = "v0.0.1-dev"
	configPlugin.TsyncPort = 50051
	configPlugin.GpsPort = 2947

	content, err := ioutil.ReadFile("/assets/sts-plugin.yaml")
	if err != nil {
		fmt.Println("ERROR: Loading sts-plugin.yaml file")
		return err
	}

	t, err := template.New("asset").Option("missingkey=error").Parse(string(content))
	if err != nil {
		fmt.Println("ERROR: New template")
		return err
	}

	err = t.Execute(&buff, configPlugin)
	if err != nil {
		fmt.Println("ERROR: Template execute failure")
		return err
	}

	rx := regexp.MustCompile("\n-{3}")
	objectsDefs := rx.Split(buff.String(), -1)

	for _, objectDef := range objectsDefs {
		obj := unstructured.Unstructured{}
		r := strings.NewReader(objectDef)
		decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
		err := decoder.Decode(&obj)
		if err != nil {
			fmt.Println("ERROR: Decoding YAML failure")
			return err
		}
		objects = append(objects, &obj)
	}

	for _, obj := range objects {
		gvk := obj.GetObjectKind().GroupVersionKind()
		old := &unstructured.Unstructured{}
		old.SetGroupVersionKind(gvk)
		key := client.ObjectKeyFromObject(obj)

		if err := k8sClient.Get(context.TODO(), key, old); err != nil {
			if err := k8sClient.Create(context.TODO(), obj); err != nil {
				if err != nil {
					panic(err)
				}
			} else {
				if !equality.Semantic.DeepDerivative(obj, old) {
					obj.SetResourceVersion(old.GetResourceVersion())
					if err := k8sClient.Update(context.TODO(), obj); err != nil {
						fmt.Println("ERROR: Update failed", "key", key, "GroupVersionKind", gvk)
						return err
					}
					fmt.Println("Object updated", "key", key, "GroupVersionKind", gvk)
				} else {
					fmt.Println("Object has not changed", "key", key, "GroupVersionKind", gvk)
				}
			}
		}
	}
	return nil
}
