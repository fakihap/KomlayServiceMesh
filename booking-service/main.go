package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// if services were to be exposed outside cni
// var grandOakURL = "http://localhost:5000/grandOak/doctors/"
// var pineValleyURL = "http://localhost:9090/pineValley/doctors"

// kubernetes cni
var grandOakURL = "http://grand-oak:5000/grandOak/doctors/"
var pineValleyURL = "http://pine-valley:9090/pineValley/doctors"

var clientset = getDefaultKubernetesClient()
var namespace = "default"
var aggregatorVersion = "v1"

type DoctorsList struct {
	Doctors []struct {
		Name     string `json:"name"`
		Time     string `json:"time"`
		Hospital string `json:"hospital"`
	} `json:"doctors"`
}

func getDefaultKubernetesClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(fmt.Errorf("failed to load in-cluster config: %v", err))
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create Kubernetes clientset: %v", err))
	}

	return clientset
}

func fetchDoctors(c *gin.Context) {
	configMaps := fetchConfigMaps()
	request := gorequest.New()

	log.Printf("configMaps: %v", configMaps)

	var out DoctorsList
	for _, cm := range configMaps {
		host := cm.Labels["app"]
		port := cm.Data["port"]
		requestTargetTemplate := cm.Data["request_target"]
		jsonRequestBodyTemplate := cm.Data["json_request_body"]
		method := cm.Data["method"]

		requestTarget := strings.Replace(requestTargetTemplate, "<doctor_type>", c.Param("doctorType"), -1)
		jsonRequestBody := strings.Replace(jsonRequestBodyTemplate, "<doctor_type>", c.Param("doctorType"), -1)
		url := fmt.Sprintf("http://%s:%s%s", host, port, requestTarget)

		var body string
		if method == "GET" {
			_, body, _ = request.Get(url).End()
		} else if method == "POST" {
			_, body, _ = request.Post(url).Send(jsonRequestBody).End()
		}

		var doctorsList DoctorsList
		err := json.Unmarshal([]byte(body), &doctorsList)
		if err != nil {
			c.JSON(500, "error unmarshalling")
			return
		}

		out.Doctors = slices.Concat(out.Doctors, doctorsList.Doctors)
	}

	c.JSON(200, out)
}

func fetchG(c *gin.Context) {
	request := gorequest.New()

	// grand-oak
	_, grandOakBody, errs := request.Get(grandOakURL + "surgeon").End()
	if errs != nil {
		c.JSON(501, errs)
	}

	c.JSON(200, grandOakBody)
}

func fetchP(c *gin.Context) {
	request := gorequest.New()

	// grand-oak
	_, pineValleyBody, errs := request.Post(pineValleyURL).
		Send(`{"doctorType": "cardiologist"}`).
		End()
	if errs != nil {
		c.JSON(501, errs)
	}

	c.JSON(200, pineValleyBody)
}

func fetchConfigMaps() []v1.ConfigMap {
	configMaps, err := clientset.CoreV1().ConfigMaps(namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: fmt.Sprintf("aggregatorVersion=%s", aggregatorVersion),
	})
	if err != nil {
		log.Fatalf("Error listing ConfigMaps: %v", err)
	}

	return configMaps.Items
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.GET("/doctors/:doctorType", fetchDoctors)

	r.GET("/g", fetchG)
	r.GET("/p", fetchP)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
