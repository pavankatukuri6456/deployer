package docker

import (
	"context"
	"deployer/model"
	"fmt"
	"math/rand"
	"time"

	tektonv1 "github.com/tektoncd/pipeline/pkg/apis/pipeline/v1"
	tektonclient "github.com/tektoncd/pipeline/pkg/client/clientset/versioned"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func generateImageTag(app, env string) string {
	rand.Seed(time.Now().UnixNano())
	randomID := rand.Intn(100000)
	return fmt.Sprintf("trialqju370.jfrog.io/pavankatukuri6456-docker/%s/%s/%d", app, env, randomID)
}

func TriggerPipeline(ctx context.Context, req model.DeployRequest) (string, string, error) {
	ns := "deployer-ns"
	imageTag := generateImageTag(req.AppName, req.Instance)

	config, err := rest.InClusterConfig()
	if err != nil {
		return "", "", err
	}
	client, err := tektonclient.NewForConfig(config)
	if err != nil {
		return "", "", err
	}

	pipelineRun := &tektonv1.PipelineRun{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: fmt.Sprintf("%s-%s-build-", req.AppName, req.Instance),
			Namespace:    ns,
		},
		Spec: tektonv1.PipelineRunSpec{
			PipelineRef: &tektonv1.PipelineRef{
				Name: "build-pipeline",
			},
			TaskRunTemplate: tektonv1.PipelineTaskRunTemplate{
				ServiceAccountName: "tekton-bot",
			},
			Params: []tektonv1.Param{
				{
					Name: "repo-url",
					Value: tektonv1.ParamValue{
						Type:      tektonv1.ParamTypeString,
						StringVal: req.RepoURL,
					},
				},
				{
					Name: "revision",
					Value: tektonv1.ParamValue{
						Type:      tektonv1.ParamTypeString,
						StringVal: req.Branch, // <--- You must include this in your DeployRequest model
					},
				},
				{
					Name: "image-url",
					Value: tektonv1.ParamValue{
						Type:      tektonv1.ParamTypeString,
						StringVal: imageTag,
					},
				},
			},
			Workspaces: []tektonv1.WorkspaceBinding{
				{
					Name: "shared-workspace",
					VolumeClaimTemplate: &corev1.PersistentVolumeClaim{
						ObjectMeta: metav1.ObjectMeta{
							GenerateName: "shared-workspace-pvc-",
						},
						Spec: corev1.PersistentVolumeClaimSpec{
							AccessModes: []corev1.PersistentVolumeAccessMode{
								corev1.ReadWriteOnce,
							},
							Resources: corev1.VolumeResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceStorage: resource.MustParse("1Gi"),
								},
							},
						},
					},
				},
			},
		},
	}

	pr, err := client.TektonV1().PipelineRuns(ns).Create(ctx, pipelineRun, metav1.CreateOptions{})
	if err != nil {
		return "", "", err
	}

	return pr.Name, imageTag, nil
}
