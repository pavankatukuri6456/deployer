package server

import (
	"context"
	"deployer/model"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func DeployApp(c *gin.Context) {

	var req model.AppRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	config, err := rest.InClusterConfig()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to load in-cluster config"})
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create clientset"})
		return
	}
	ctx := context.Background()

	if err := ensureNamespaceExists(ctx, clientset, req.Namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := ensureImagePullSecret(ctx, clientset, req.Namespace); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := createDeployment(ctx, clientset, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := createService(ctx, clientset, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if err := createIngress(ctx, clientset, req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Application '%s' deployed in namespace '%s'", req.AppName, req.Namespace)})
}

func ensureNamespaceExists(ctx context.Context, clientset *kubernetes.Clientset, ns string) error {
	_, err := clientset.CoreV1().Namespaces().Get(ctx, ns, v1.GetOptions{})
	if errors.IsNotFound(err) {
		_, err := clientset.CoreV1().Namespaces().Create(ctx, &corev1.Namespace{
			ObjectMeta: v1.ObjectMeta{Name: ns},
		}, v1.CreateOptions{})
		return err
	}
	return err
}

func ensureImagePullSecret(ctx context.Context, clientset *kubernetes.Clientset, ns string) error {
	secret, err := clientset.CoreV1().Secrets("default").Get(ctx, "jfrog-docker-config", v1.GetOptions{})
	if err != nil {
		return err
	}
	secret.ObjectMeta = v1.ObjectMeta{Name: "jfrog-docker-config", Namespace: ns}
	_, err = clientset.CoreV1().Secrets(ns).Create(ctx, secret, v1.CreateOptions{})
	if errors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func createDeployment(ctx context.Context, clientset *kubernetes.Clientset, req model.AppRequest) error {
	deploy := &appsv1.Deployment{
		ObjectMeta: v1.ObjectMeta{Name: req.AppName, Namespace: req.Namespace},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &v1.LabelSelector{MatchLabels: map[string]string{"app": req.AppName}},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: v1.ObjectMeta{Labels: map[string]string{"app": req.AppName}},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  req.AppName,
						Image: req.Image,
						Ports: []corev1.ContainerPort{{ContainerPort: 8080}},
					}},
					ImagePullSecrets: []corev1.LocalObjectReference{{Name: "jfrog-docker-config"}},
				},
			},
		},
	}
	_, err := clientset.AppsV1().Deployments(req.Namespace).Create(ctx, deploy, v1.CreateOptions{})
	return err
}

func createService(ctx context.Context, clientset *kubernetes.Clientset, req model.AppRequest) error {
	service := &corev1.Service{
		ObjectMeta: v1.ObjectMeta{Name: req.AppName, Namespace: req.Namespace},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{"app": req.AppName},
			Ports: []corev1.ServicePort{{
				Port:       80,
				TargetPort: intstrFromInt(8080),
			}},
		},
	}
	_, err := clientset.CoreV1().Services(req.Namespace).Create(ctx, service, v1.CreateOptions{})
	return err

}

func createIngress(ctx context.Context, clientset *kubernetes.Clientset, req model.AppRequest) error {
	pathType := networkingv1.PathTypePrefix
	ing := &networkingv1.Ingress{
		ObjectMeta: v1.ObjectMeta{
			Name:      req.AppName + "-ingress",
			Namespace: req.Namespace,
			Annotations: map[string]string{
				"nginx.ingress.kubernetes.io/rewrite-target": "/",
			},
		},
		Spec: networkingv1.IngressSpec{
			Rules: []networkingv1.IngressRule{{
				Host: req.AppName + ".local",
				IngressRuleValue: networkingv1.IngressRuleValue{
					HTTP: &networkingv1.HTTPIngressRuleValue{
						Paths: []networkingv1.HTTPIngressPath{{
							Path:     "/",
							PathType: &pathType,
							Backend: networkingv1.IngressBackend{
								Service: &networkingv1.IngressServiceBackend{
									Name: req.AppName,
									Port: networkingv1.ServiceBackendPort{Number: 80},
								},
							},
						}},
					},
				},
			}},
		},
	}
	_, err := clientset.NetworkingV1().Ingresses(req.Namespace).Create(ctx, ing, v1.CreateOptions{})
	return err
}

func int32Ptr(i int32) *int32 { return &i }

func intstrFromInt(i int) intstr.IntOrString {
	return intstr.IntOrString{Type: intstr.Int, IntVal: int32(i)}
}
