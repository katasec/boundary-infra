package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	v1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func InitalAuthInfoFromJob(namespace string, jobName *string) (authMethod string, loginName string, password string) {
	/*
		-------------------------------------------
		Ouput from Boundary DB init job looks like this:
		--------------------------------------------

		Initial auth information:
		  Auth Method ID:     ampw_Xxxxxxxxxx
		  Auth Method Name:   Generated global scope initial password auth method
		  Login Name:         admin
		  Password:           xxxxxxxxxxxxxxxxxxxx
		  Scope ID:           global
	*/

	// Get pod name for K8s job
	label := fmt.Sprintf("job-name=%s", *jobName)
	podName := GetPodByLabel(namespace, label)

	logs := GetPodLogs(namespace, podName)

	// Find index where auth info is displayed in the pod logs
	index := Find(logs, "Initial auth information:")
	if index == -1 {
		log.Fatal("Initial auth info was not found, exitting !")
	}

	// Grab AuthMethod from line 2 in the output
	authMethod = strings.Trim(strings.Split(logs[index+1], ":")[1], " ")
	if authMethod == "" {
		log.Fatalf(errors.New("could not get auth method from logs").Error())
	}
	// Grab AuthMethod from line 4 in the output
	loginName = strings.Trim(strings.Split(logs[index+3], ":")[1], " ")
	if loginName == "" {
		log.Fatalf(errors.New("could not get loginName from logs").Error())
	}
	// Grab AuthMethod from line 4 in the output
	password = strings.Trim(strings.Split(logs[index+4], ":")[1], " ")
	if password == "" {
		log.Fatalf(errors.New("could not get password from logs").Error())
	}
	// if len(password) != 20 {
	// 	log.Fatal("Password was not 20 chars, parsing error, exiting !")
	// }

	return authMethod, loginName, password
}

func GetPodInfoByLabel(namespace string, label string) corev1.Pod {

	client := GetK8sClientSet()

	// Find pod with label
	pods, err := client.CoreV1().Pods(namespace).List(
		context.Background(),
		metav1.ListOptions{LabelSelector: label},
	)

	// Fail if pod was not found
	if err != nil {
		log.Fatalf("Could not find pod with '%s':\n %v ", label, err.Error())
	}
	if len(pods.Items) == 0 {
		log.Fatalf("Could not find pod with '%s':\n", label)
	}

	// Return first pod with given label
	return pods.Items[0]
}

func GetPodByLabel(namespace string, label string) string {

	client := GetK8sClientSet()

	// Find pod with label
	pods, err := client.CoreV1().Pods(namespace).List(
		context.Background(),
		metav1.ListOptions{LabelSelector: label},
	)

	// Fail if pod was not found
	if err != nil {
		log.Fatalf("Could not find pod with '%s':\n %v ", label, err.Error())
	}
	if len(pods.Items) == 0 {
		log.Fatalf("Could not find pod with '%s':\n", label)
	}

	// Return first pod with given label
	return pods.Items[0].Name
}

func GetPodLogs(namespace string, podName string) []string {

	// Get k8s client
	client := GetK8sClientSet()

	// Request log stream
	req := client.CoreV1().Pods(namespace).GetLogs(podName, &corev1.PodLogOptions{})
	ctx := context.Background()
	logStream, err := req.Stream(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer logStream.Close()

	// Create byte buffer from Log stream
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logStream)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Return logs as string
	return strings.Split(buf.String(), "\n")
}

func GetK8sClientSet() *kubernetes.Clientset {
	// Check k8s config file exits
	configFile := filepath.Join(homedir.HomeDir(), ".kube", "config")
	if _, err := os.Stat(configFile); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("%s does not exist, exitting !", configFile)
	}

	// Use current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", configFile)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Create K8s client
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	return clientset
}

func GetJobStatus(namespace string, jobName string) (status int, err error) {

	// Get k8s client
	client := GetK8sClientSet()

	job, err := client.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return 5, nil
		} else {
			log.Fatal(err.Error())
			//return -1, err // Error
		}
	}

	if job.Status.Active == 0 && job.Status.Succeeded == 0 && job.Status.Failed == 0 {
		return 0, nil // not started
	}

	if job.Status.Active > 0 {
		return 1, nil // still running
	}

	if job.Status.Succeeded > 0 {
		return 2, nil // Job ran successfully
	}

	return 4, nil // Unknown
}

func GetJobDetails(namespace string, jobName string) *v1.Job {

	// Get k8s client
	client := GetK8sClientSet()

	job, err := client.BatchV1().Jobs(namespace).Get(context.Background(), jobName, metav1.GetOptions{})
	if err != nil {
		log.Fatal(err.Error())
	}

	return job
}
