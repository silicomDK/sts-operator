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

	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/yaml"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/go-logr/logr"
	srov1beta1 "github.com/openshift-psap/special-resource-operator/api/v1beta1"
	helmerv1beta1 "github.com/openshift-psap/special-resource-operator/pkg/helmer/api/v1beta1"
	nfdv1 "github.com/openshift/cluster-nfd-operator/api/v1"
	stsv1alpha1 "github.com/silicomdk/sts-operator/api/v1alpha1"
)

// StsOperatorConfigReconciler reconciles a StsOperatorConfig object
type StsOperatorConfigReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Log    logr.Logger
}

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
	operatorCfgList := &stsv1alpha1.StsOperatorConfigList{}

	opts := (&client.ListOptions{}).ApplyOptions([]client.ListOption{client.InNamespace(req.NamespacedName.Namespace)})
	err := r.List(ctx, operatorCfgList, opts)
	if err != nil {
		reqLogger.Info("No Operator CR found in this namespace")
		return ctrl.Result{}, err
	}

	//	finalizer := "sts.silicom.com/finalizer"

	operatorCfg := &operatorCfgList.Items[0]

	err = r.DeployNfd(operatorCfg)
	if err != nil {
		reqLogger.Error(err, "Failed to create NFD CR")
		return ctrl.Result{}, err
	}

	err = r.DeploySro(operatorCfg)
	if err != nil {
		reqLogger.Error(err, "Failed to deploy SRO requirements")
		return ctrl.Result{}, err
	}

	err = r.DeployPlugin(operatorCfg)
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

func (r *StsOperatorConfigReconciler) DeploySro(operatorCfg *stsv1alpha1.StsOperatorConfig) error {

	if len(operatorCfg.Spec.Sro.Chart.Repository.URL) < 1 {
		operatorCfg.Spec.Sro.Chart.Repository.URL = "http://ice-driver-src"
	}

	svc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ice-driver-src",
			Namespace: operatorCfg.Spec.Sro.Namespace,
			Labels: map[string]string{
				"app": "ice-driver-src",
			},
			Annotations: map[string]string{
				"openshift.io/scc": "sts-silicom",
			},
		},
		Spec: v1.ServiceSpec{
			Type:     "NodePort",
			Selector: map[string]string{"app": "ice-driver-src"},
			Ports: []v1.ServicePort{
				{
					Name:       "ice-driver-src",
					Port:       int32(operatorCfg.Spec.Sro.SrcSvcPort),
					TargetPort: intstr.FromInt(operatorCfg.Spec.Sro.SrcSvcPort),
					Protocol:   "TCP",
				},
			},
		},
	}

	if err := r.Get(context.TODO(), client.ObjectKey{
		Namespace: operatorCfg.Spec.Sro.Namespace,
		Name:      svc.Name,
	}, svc); err != nil {

		err = r.Create(context.TODO(), svc)
		if err != nil {
			panic(err)
		}
	} else {
		r.Update(context.TODO(), svc)
		if err != nil {
			panic(err)
		}
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ice-driver-src",
			Namespace: operatorCfg.Spec.Sro.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": "ice-driver-src",
				},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": "ice-driver-src",
					},
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "ice-driver-src",
							Image:           operatorCfg.Spec.Sro.SrcImage,
							ImagePullPolicy: "Always",
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									Protocol:      v1.ProtocolTCP,
									ContainerPort: int32(operatorCfg.Spec.Sro.SrcSvcPort),
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Get(context.TODO(), client.ObjectKey{
		Namespace: operatorCfg.Spec.Sro.Namespace,
		Name:      deployment.Name,
	}, deployment); err != nil {

		err = r.Create(context.TODO(), deployment)
		if err != nil {
			panic(err)
		}
	} else {
		r.Update(context.TODO(), deployment)
		if err != nil {
			panic(err)
		}
	}

	sr := &srov1beta1.SpecialResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "ice-special-resource",
			Namespace: operatorCfg.Spec.Sro.Namespace,
		},
		Spec: srov1beta1.SpecialResourceSpec{
			Debug:        false,
			Namespace:    operatorCfg.Spec.Sro.Namespace,
			NodeSelector: map[string]string{"feature.node.kubernetes.io/custom-intel.e810_c.devices": "true"},
			Chart: helmerv1beta1.HelmChart{
				Version: "0.0.1",
				Name:    "ice-special-resource",
				Repository: helmerv1beta1.HelmRepo{
					Name: "ice-special-resource",
					URL:  fmt.Sprintf("%s:%d", operatorCfg.Spec.Sro.Chart.Repository.URL, operatorCfg.Spec.Sro.SrcSvcPort),
				},
			},
			Set: unstructured.Unstructured{
				Object: map[string]interface{}{
					"kind":           "Values",
					"apiVersion":     "sro.openshift.io/v1beta1",
					"driverRegistry": fmt.Sprintf("%s/%s", operatorCfg.Spec.Sro.DriverRegistry, operatorCfg.Spec.Sro.Namespace),
					"buildArgs": []map[string]interface{}{
						{
							"name":  "ICE_VERSION",
							"value": operatorCfg.Spec.Sro.IceVersion,
						},
						{
							"name":  "ICE_SRC",
							"value": operatorCfg.Spec.Sro.SrcImage,
						},
					},
					"runArgs": map[string]interface{}{
						"platform": "openshift-container-platform",
						"buildIce": operatorCfg.Spec.Sro.Build,
					},
				},
			},
		},
	}

	if err := r.Get(context.TODO(), client.ObjectKey{
		Namespace: operatorCfg.Spec.Sro.Namespace,
		Name:      sr.Name,
	}, sr); err != nil {

		err = r.Create(context.TODO(), sr)
		if err != nil {
			panic(err)
		}
	} else {
		r.Update(context.TODO(), sr)
		if err != nil {
			panic(err)
		}
	}

	return nil
}

func (r *StsOperatorConfigReconciler) DeployNfd(operatorCfg *stsv1alpha1.StsOperatorConfig) error {

	nfdOperand := &nfdv1.NodeFeatureDiscovery{}
	nfdOperand.Name = "nfd-sts-silicom"
	nfdOperand.Namespace = operatorCfg.Namespace
	nfdOperand.Spec.Operand.Namespace = operatorCfg.Namespace

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

func (r *StsOperatorConfigReconciler) DeployPlugin(operatorCfg *stsv1alpha1.StsOperatorConfig) error {
	var buff bytes.Buffer
	var objects []client.Object

	fmt.Printf("Starting plugin in ns: %s\n", operatorCfg.Namespace)

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

	err = t.Execute(&buff, operatorCfg)
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
