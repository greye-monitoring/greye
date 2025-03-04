package repositories

//
//import (
//	"greye/pkg/importProcess/domain/ports"
//	"flag"
//	"k8s.io/client-go/kubernetes"
//	"k8s.io/client-go/tools/clientcmd"
//	"k8s.io/client-go/util/homedir"
//	"path/filepath"
//)
//
//type KubernetesRepository struct {
//	serviceName string
//}
//
//var _ ports.ImportProcessRepository = (*KubernetesRepository)(nil)
//
//func NewKubernetesRepository(serviceName string) *KubernetesRepository {
//	return &KubernetesRepository{
//		serviceName: serviceName,
//	}
//}
//
//// GetConfig returns a Kubernetes Clientset based on the service name.
//func (repo *KubernetesRepository) GetConfig() (*kubernetes.Clientset, error) {
//	if repo.serviceName == "localhost" {
//		return getLocalClientset()
//	}
//	return getInClusterClientset()
//}
//
//func getInClusterClientset() (*kubernetes.Clientset, error) {
//	var kubeconfig *string
//	if home := homedir.HomeDir(); home != "" {
//		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
//	} else {
//		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
//	}
//	flag.Parse()
//
//	// Use the current context in kubeconfig
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		return nil, err
//	}
//
//	return kubernetes.NewForConfig(config)
//}
//
//func getLocalClientset() (*kubernetes.Clientset, error) {
//	var kubeconfig *string
//	if home := homedir.HomeDir(); home != "" {
//		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
//	} else {
//		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
//	}
//	flag.Parse()
//
//	// use the current context in kubeconfig
//	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
//	if err != nil {
//		panic(err.Error())
//	}
//
//	// create the clientset
//	clientset, err := kubernetes.NewForConfig(config)
//	if err != nil {
//		panic(err.Error())
//	}
//	return clientset, nil
//}
