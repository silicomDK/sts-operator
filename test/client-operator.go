package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

type StsConfig struct {
	Mode          string
	Interfaces    []string
	EnableGPS     bool
	Name          string
	NodeLabel     string
	SilTsyncImage string
	Namespace     string
}

func _template() {
	content, err := ioutil.ReadFile("../assets/sts-deployment.yaml")
	if err != nil {
		return
	}

	stsConfig := StsConfig{"master",
		[]string{"enp2s0f1", "enp2s0f2", "enp2s0f3"},
		true,
		"test-sts1",
		"feature.node.kubernetes.io/usb-ff_1374_0001.present",
		"quay.io/silicom/siltsync:1.2.0.1",
		"sts-silicom"}

	t, err := template.New("asset").Parse(string(content))
	if err != nil {
		fmt.Println("ERROR2: Reconciling StsConfig")
		return
	}

	var buff bytes.Buffer
	err = t.Execute(&buff, stsConfig)
	if err != nil {
		fmt.Println("ERROR3: Reconciling StsConfig")
		return
	}

	rx := regexp.MustCompile("\n-{3}")
	objectsDefs := rx.Split(buff.String(), -1)

	for _, objectDef := range objectsDefs {
		obj := unstructured.Unstructured{}
		r := strings.NewReader(objectDef)
		decoder := yaml.NewYAMLOrJSONDecoder(r, 4096)
		err := decoder.Decode(&obj)
		if len(objectDef) == 0 {
			fmt.Println("emtpy")
			continue
		}
		if err != nil {
			fmt.Printf("ERROR4: Reconciling StsConfig %s/%d\n", objectDef, len(objectDef))
			return
		}

		fmt.Println(obj.GetName())
		//objects, _ = append(objects, &obj)
	}
}

func main() {
	_template()
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// node-role.sts.silicom.com/master
	// node-role.sts.silicom.com/boundary
	nodes, err := clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{LabelSelector: "feature.node.kubernetes.io/usb-ff_1374_0001.present"})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d nodes in the cluster\n", len(nodes.Items))

	for _, node := range nodes.Items {
		for name, value := range node.Labels {
			fmt.Println(name, value)
			if strings.Contains(name, "sts.silicom.com") {
				fmt.Println(fmt.Sprint("Removing ", name))
				delete(node.Labels, name)
			}
		}
		_, err := clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
		if err != nil {
			return
		}

		out, err := exec.Command("lspci", "-n", "-d", "8086:02f0").Output()
		id := strings.Split(string(out), " ")[0]
		path := fmt.Sprintf("/sys/bus/pci/devices/0000:%s/net/*", id)
		res, _ := filepath.Glob(path)

		node.Labels["iface.sts.silicom.com"] = string(filepath.Base(res[0]))
		fmt.Println(fmt.Sprintf("***** ADDING %s \n", filepath.Base(res[0])))

		_, err = clientset.CoreV1().Nodes().Update(context.TODO(), &node, metav1.UpdateOptions{})
		if err != nil {
			return
		}
	}

}

/*
// UpdateNode updates node
func UpdateNode(client *clientset, node *v1.Node) error {
	_, err := client.Core().Nodes().Update(node)
	if err != nil {
		return err
	}

	return nil
}

// RemoveCPUModelNodeLabels removes labels from node which were created by kubevirt-node-labeller
func RemoveCPUModelNodeLabels(node *v1.Node, oldLabels map[string]bool) {
	for label := range node.Labels {
		if ok := oldLabels[label]; ok || strings.Contains(label, labelNamespace+"/cpu-model-") {
			delete(node.Labels, label)
		}
	}
}

// GetNode gets node by name
func GetNode(client *clientset) (*v1.Node, error) {
	nodeName := os.Getenv("NODE_NAME")

	node, err := client.Core().Nodes().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	return node, nil
}

// AddNodeLabels adds labels to node
func AddNodeLabels(node *v1.Node, labels map[string]string) {
	for name, value := range labels {
		node.Labels[labelNamespace+name] = value
		node.Annotations[labellerNamespace+"-"+labelNamespace+name] = value
	}
}

// GetNodeLabellerLabels gets all labels which were created by kubevirt-node-labeller
func GetNodeLabellerLabels(node *v1.Node) map[string]bool {
	labellerLabels := make(map[string]bool)
	for key := range node.Annotations {
		if strings.Contains(key, labellerNamespace) {
			delete(node.Annotations, key)
			labellerLabels[key] = true
		}
	}
	return labellerLabels
}
*/
