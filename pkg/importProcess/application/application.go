package application

import (
	"flag"
	"golang.org/x/net/context"
	"greye/pkg/importProcess/domain/ports"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

type ImportProcessApplication struct {
	clientSet *kubernetes.Clientset
}

var _ ports.ImportProcessApplication = (*ImportProcessApplication)(nil)

func NewImportProcessApplication(url string) *ImportProcessApplication {

	clientSet, _ := CreateConfig(url)
	return &ImportProcessApplication{
		clientSet: clientSet,
	}
}

func (app *ImportProcessApplication) GetKubernetesMonitoringObject(resourceVersion string) watch.Interface {
	//val := func(i int64) *int64 { return &i }(10)
	//svcWatch, err := app.clientSet.CoreV1().Services("").Watch(context.TODO(), metav1.ListOptions{TimeoutSeconds: func(i int64) *int64 { return &i }(20), ResourceVersion: resourceVersion}) //svcList, err := app.clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	svcWatch, err := app.clientSet.CoreV1().Services("").Watch(context.TODO(), metav1.ListOptions{ResourceVersion: resourceVersion}) //svcList, err := app.clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

	//var services []v1.Service
	//for _, svc := range svcList.Items {
	//	if svc.Annotations["cm-enabled"] == "true" {
	//		services = append(services, svc)
	//	}
	//}
	//
	////pods, err := app.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("There are %d pods in the cluster\n", len(services))

	return svcWatch
}

func (app *ImportProcessApplication) GetKubernetesServices() *v1.ServiceList {

	svcList, err := app.clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})

	//svcList, err := app.clientSet.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	//var services []v1.Service
	//for _, svc := range svcList.Items {
	//	if svc.Annotations["cm-enabled"] == "true" {
	//		services = append(services, svc)
	//	}
	//}
	//
	////pods, err := app.clientSet.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	//fmt.Printf("There are %d pods in the cluster\n", len(services))

	return svcList
}

// // GetConfig returns a Kubernetes Clientset based on the service name.
func CreateConfig(url string) (*kubernetes.Clientset, error) {
	if url == "localhost" {
		return getLocalClientset()
	}
	return getInClusterClientset()
}

func getInClusterClientset() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset, nil
}

func getLocalClientset() (*kubernetes.Clientset, error) {
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
	return clientset, nil
}
