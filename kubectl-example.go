package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	batchv1 "k8s.io/api/batch/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	// Uncomment to load all auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

var kubeClient *kubernetes.Clientset

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Only support POST")
		return
	}
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var msg map[string]string
	json.Unmarshal(b, &msg)
	// fmt.Printf("Request body: %s", msg)

	job := createJob(msg)
	// deploy a job
	jobsClient := kubeClient.BatchV1().Jobs("default")
	backofflimit := new(int32)
	*backofflimit = 5

	result1, err1 := jobsClient.Create(job)
	if err != nil {
		fmt.Println(err1)
		panic(err1)
	}
	fmt.Printf("Created job %q.\n", result1)
	fmt.Fprintf(w, "Start running a job successful !")
}

func createClient() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	kubeClient, err = kubernetes.NewForConfig(config)
	return err
}

func conertToContainerEnvVar(msg map[string]string) []apiv1.EnvVar {
	result := []apiv1.EnvVar{}
	for k, v := range msg {
		result = append(result, apiv1.EnvVar{
			Name:  k,
			Value: v,
		})
	}
	return result
}

func createJob(msg map[string]string) *batchv1.Job {
	backofflimit := new(int32)
	*backofflimit = 5
	job := &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "build-firmware-" + msg["ORDER_NUM"],
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
							Env:   conertToContainerEnvVar(msg),
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
	return job
}

func main() {
	err := createClient()
	if err != nil {
		panic(err.Error())
	}
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
