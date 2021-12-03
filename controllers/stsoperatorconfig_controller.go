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

package controllers

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	nfdv1 "github.com/openshift/cluster-nfd-operator/api/v1"
	stsv1alpha1 "github.com/silicomdk/sts-operator/api/v1alpha1"
)

// StsOperatorConfigReconciler reconciles a StsOperatorConfig object
type StsOperatorConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

const defaultName = "sts-operator-config"
const defaultNamespace = "sts-silicom"

//+kubebuilder:rbac:groups=sts.silicom.com,resources=stsoperatorconfigs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=sts.silicom.com,resources=stsoperatorconfigs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=sts.silicom.com,resources=stsoperatorconfigs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the StsOperatorConfig object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.10.0/pkg/reconcile
func (r *StsOperatorConfigReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	reqLogger := r.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling StsOperatorConfig")

	// Fetch the StsOperatorConfig instance
	defaultCfg := &stsv1alpha1.StsOperatorConfig{}
	err := r.Get(context.TODO(), types.NamespacedName{
		Name: defaultName, Namespace: defaultNamespace}, defaultCfg)
	if err != nil {
		reqLogger.Error(err, "Failed to get Deafult Operator CR")
		return ctrl.Result{}, err
	}

	err = r.DeployNfd(defaultCfg)
	if err != nil {
		reqLogger.Error(err, "Failed to create NFD CR")
		return ctrl.Result{}, err
	}

	err = r.DeployPlugin(defaultCfg)
	if err != nil {
		reqLogger.Error(err, "Failed to deploy plugin daemons")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *StsOperatorConfigReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&stsv1alpha1.StsOperatorConfig{}).
		Complete(r)
}

func (r *StsOperatorConfigReconciler) DeployNfd(defaultCfg *stsv1alpha1.StsOperatorConfig) error {

	nfdOperand := &nfdv1.NodeFeatureDiscovery{}
	nfdOperand.Name = "nfd-sts-silicom"
	nfdOperand.Namespace = defaultCfg.Spec.Namespace
	nfdOperand.Spec.Operand.Namespace = defaultCfg.Spec.Namespace

	content, err := ioutil.ReadFile("/assets/nfd-discovery.yaml")
	if err != nil {
		fmt.Println("ERROR: Loading nfd-discovery.yaml file")
		return err
	}

	workerConfig := &nfdv1.ConfigMap{}
	workerConfig.ConfigData = string(content)
	nfdOperand.Spec.WorkerConfig = workerConfig

	if err := r.Get(context.TODO(), client.ObjectKey{
		Namespace: nfdOperand.Namespace,
		Name:      nfdOperand.Name,
	}, nfdOperand); err != nil {

		err = r.Create(context.TODO(), nfdOperand)
		if err != nil {
			panic(err)
		}
	} else {
		r.Update(context.TODO(), nfdOperand)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (r *StsOperatorConfigReconciler) DeployPlugin(defaultCfg *stsv1alpha1.StsOperatorConfig) error {
	var buff bytes.Buffer
	var objects []client.Object

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

	err = t.Execute(&buff, defaultCfg)
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

		if err := r.Get(context.TODO(), key, old); err != nil {
			if err := r.Create(context.TODO(), obj); err != nil {
				if err != nil {
					panic(err)
				}
			} else {
				if !equality.Semantic.DeepDerivative(obj, old) {
					obj.SetResourceVersion(old.GetResourceVersion())
					if err := r.Update(context.TODO(), obj); err != nil {
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
