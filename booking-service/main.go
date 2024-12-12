package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"sync"

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

type Doctor struct {
	Name     string `json:"name"`
	Time     string `json:"time"`
	Hospital string `json:"hospital"`
}

type DoctorsList struct {
	Doctors []Doctor `json:"doctors"`
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

	log.Printf("configMaps: %v", configMaps)

	out := DoctorsList{
		Doctors: []Doctor{},
	}
	var wg sync.WaitGroup
	res := make(chan DoctorsList, len(configMaps))
	for _, cm := range configMaps {
		host := cm.Labels["app"]
		port := cm.Data["port"]
		requestTargetTemplate := cm.Data["request_target"]
		jsonRequestBodyTemplate := cm.Data["json_request_body"]
		method := cm.Data["method"]

		requestTarget := strings.Replace(requestTargetTemplate, "<doctor_type>", c.Param("doctorType"), -1)
		jsonRequestBody := strings.Replace(jsonRequestBodyTemplate, "<doctor_type>", c.Param("doctorType"), -1)
		url := fmt.Sprintf("http://%s:%s%s", host, port, requestTarget)

		wg.Add(1)
		go func() {
			defer wg.Done()
			request := gorequest.New()

			var body string
			if method == "GET" {
				_, body, _ = request.Get(url).End()
			} else if method == "POST" {
				_, body, _ = request.Post(url).Send(jsonRequestBody).End()
			}

			log.Println(url)
			log.Println(body)

			var doctorsList DoctorsList
			err := json.Unmarshal([]byte(body), &doctorsList)
			if err == nil {
				res <- doctorsList
			} else {
				res <- DoctorsList{
					Doctors: []Doctor{},
				}
				log.Printf("error: %v\n", err)
			}
			log.Printf("doctorsList: %v\n", doctorsList)
		}()
	}

	wg.Wait()
	close(res)

	for doctorsList := range res {
		out.Doctors = append(out.Doctors, doctorsList.Doctors...)
	}
	log.Printf("out: %v\n", out)

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
