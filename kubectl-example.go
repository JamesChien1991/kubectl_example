package p

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"

	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func KubectlExample(w http.ResponseWriter, r *http.Request) {
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

	var d struct {
		Message string `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello World!")
		return
	}
	if d.Message == "" {
		fmt.Fprint(w, "Hello World!")
		return
	}
	fmt.Fprint(w, html.EscapeString(d.Message))
}
