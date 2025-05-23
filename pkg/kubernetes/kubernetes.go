package kubernetes

import (
	"context"
	"fmt"
	"path/filepath"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	metricsv1beta1 "k8s.io/metrics/pkg/client/clientset/versioned"
)

type Client struct {
	clientset        *kubernetes.Clientset
	metricsClientset *metricsv1beta1.Clientset
}

func NewClient(kubeConfigPath string) (*Client, error) {
	var config *rest.Config
	var err error

	if kubeConfigPath == "" {
		home := homedir.HomeDir()
		if home != "" {
			kubeConfigPath = filepath.Join(home, ".kube", "config")
		}
	}

	if kubeConfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			config, err = rest.InClusterConfig()
			if err != nil {
				return nil, fmt.Errorf("Failed to create kubernetes config %v", err)
			}
		}
	} else {
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, fmt.Errorf("failed to create in-cluster config: %v", err)
		}
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create kubernetes clientset %v", err)
	}

	metricsClientset, err := metricsv1beta1.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create metrics clientset: %v", err)
	}

	return &Client{
		clientset:        clientset,
		metricsClientset: metricsClientset,
	}, nil
}

func (c *Client) GetPods(namespace string) (*corev1.PodList, error) {
	return c.clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetNodes(namespace string) (*corev1.NodeList, error) {
	return c.clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetNamespaces() (*corev1.NamespaceList, error) {
	return c.clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetPodMetrics(namespace string) (*metricsv1beta1.PodMetricsList, error) {
	return c.metricsClientset.MetricsV1beta1().PodMetricses(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetNodeMetrics() (*metricsv1beta1.NodeMetricsList, error) {
	return c.metricsClientset.MetricsV1beta1().NodeMetricses().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetServices(namespace string) (*corev1.ServiceList, error) {
	return c.clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetPersistentVolumes() (*corev1.PersistentVolumeList, error) {
	return c.clientset.CoreV1().PersistentVolumes().List(context.TODO(), metav1.ListOptions{})
}

func (c *Client) GetPersistentVolumeClaims(namespace string) (*corev1.PersistentVolumeClaimList, error) {
	return c.clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
}
