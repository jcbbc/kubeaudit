package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	apiv1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/version"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp" // auth for GKE clusters
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func kubeClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig == "" && !rootConfig.localMode {
		fmt.Println("Inside cluster. Not using kube config")
		// creates the in-cluster config
		return rest.InClusterConfig()
	}

	// generate config from kubectl config current-context
	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func kubeClient(kubeconfig string) (*kubernetes.Clientset, error) {
	if rootConfig.localMode {
		kubeconfig = os.Getenv("HOME") + "/.kube/config"
	}
	if rootConfig.verbose {
		log.SetLevel(log.DebugLevel)
		log.AddHook(NewDebugHook())
	}
	if rootConfig.json {
		log.SetFormatter(&log.JSONFormatter{})
	}
	config, err := kubeClientConfig(kubeconfig)
	if err != nil {
		log.Error(err)
	}
	kube, err := kubernetes.NewForConfig(config)
	printKubernetesVersion(kube)
	return kube, err
}

func getNamespaces(clientset *kubernetes.Clientset) *NamespaceList {
	namespaceClient := clientset.Namespaces()
	namespaces, err := namespaceClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return namespaces
}

func getDeployments(clientset *kubernetes.Clientset) *DeploymentList {
	deploymentClient := clientset.AppsV1beta1().Deployments(apiv1.NamespaceAll)
	deployments, err := deploymentClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return deployments
}

func getStatefulSets(clientset *kubernetes.Clientset) *StatefulSetList {
	statefulSetClient := clientset.AppsV1beta1().StatefulSets(apiv1.NamespaceAll)
	statefulSets, err := statefulSetClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return statefulSets
}

func getDaemonSets(clientset *kubernetes.Clientset) *DaemonSetList {
	daemonSetClient := clientset.ExtensionsV1beta1().DaemonSets(apiv1.NamespaceAll)
	daemonSets, err := daemonSetClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return daemonSets
}

func getPods(clientset *kubernetes.Clientset) *PodList {
	podClient := clientset.Pods(apiv1.NamespaceAll)
	pods, err := podClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return pods
}

func getReplicationControllers(clientset *kubernetes.Clientset) *ReplicationControllerList {
	replicationControllerClient := clientset.ReplicationControllers(apiv1.NamespaceAll)
	replicationControllers, err := replicationControllerClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return replicationControllers
}

func getNetworkPolicies(clientset *kubernetes.Clientset) *networking.NetworkPolicyList {
	netPolClient := clientset.NetworkPolicies(apiv1.NamespaceAll)
	netPols, err := netPolClient.List(ListOptions{})
	if err != nil {
		log.Error(err)
	}
	return netPols
}

func printKubernetesVersion(clientset *kubernetes.Clientset) {
	discoveryClient := clientset.Discovery()
	serverInfo, err := discoveryClient.ServerVersion()
	if err != nil {
		log.Error(err)
	}
	log.WithFields(log.Fields{
		"Major":    serverInfo.Major,
		"Minor":    serverInfo.Minor,
		"Platform": serverInfo.Platform,
	}).Info("Kubernetes server version")

	clientInfo := version.Get()
	log.WithFields(log.Fields{
		"Major":    clientInfo.Major,
		"Minor":    clientInfo.Minor,
		"Platform": clientInfo.Platform,
	}).Info("Kubernetes client version")
}
