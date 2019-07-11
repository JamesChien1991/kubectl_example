/*
Copyright 2016 The Kubernetes Authors.
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

// Note: the example only works with the code within the same release/branch.
package main

import (
	"fmt"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func main() {
	kubeconfig := "./kube_config"
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	// do something
	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	// deploy a job
	jobsClient := clientset.BatchV1().Jobs("default")
	// result1, err1 := jobsClient.List(metav1.ListOptions{})
	// if err != nil {
	// 	fmt.Println(err1)
	// 	panic(err1)
	// }
	// fmt.Println(result1.Items[0].Spec)
	// fmt.Println(fmt.Sprintf("%T", result1.Items[0].Spec))
	backofflimit := new(int32)
	*backofflimit = 5
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "build-firmware",
			Namespace: "default",
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: backofflimit,
			Template: apiv1.PodTemplateSpec{
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  "build-firmware",
							Image: "gcr.io/ultron-cms/helloworld",
							Env: []apiv1.EnvVar{
								{
									Name:  "HELLO_MSG",
									Value: "James !!!!!!",
								},
							},
						},
					},
					RestartPolicy: "Never",
				},
				ObjectMeta: metav1.ObjectMeta{
					Name: "build-firmware",
				},
			},
		},
	}
	result1, err1 := jobsClient.Create(job)
	if err != nil {
		fmt.Println(err1)
		panic(err1)
	}
	fmt.Printf("Created job %q.\n", result1)
}
