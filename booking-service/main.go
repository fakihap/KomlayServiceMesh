package main

import "github.com/gin-gonic/gin"
import "github.com/parnurzeal/gorequest"

import "encoding/json"
// import "fmt"
import "slices"


// if services were to be exposed outside cni
// var grandOakURL = "http://localhost:5000/grandOak/doctors/"
// var pineValleyURL = "http://localhost:9090/pineValley/doctors"

// kubernetes cni
var grandOakURL = "http://grand-oak:5000/grandOak/doctors/"
var pineValleyURL = "http://pine-valley:9090/pineValley/doctors"

type DoctorsList struct {
	Doctors struct {
		Doctor []struct {
			Name     string `json:"name"`
			Time     string `json:"time"`
			Hospital string `json:"hospital"`
		} `json:"doctor"`
	} `json:"doctors"`
}

func fetchDoctors(c *gin.Context) {
    request := gorequest.New()

    // grand-oak
    _, grandOakBody, _ := request.Get(grandOakURL + c.Param("doctorType")).End()

    // pine-valley
    _, pineValleyBody, _ := request.Post(pineValleyURL).
        Send(`{"doctorType": "` + c.Param("doctorType") + `"}`).
        End()

	var outGrandOak DoctorsList
    var outPineValley DoctorsList

	err1 := json.Unmarshal([]byte(grandOakBody), &outGrandOak)
	if err1 != nil {
        // fmt.Println(err1.Error())

        c.JSON(500, "error unmarshalling")
        return
    }	

    err2 := json.Unmarshal([]byte(pineValleyBody), &outPineValley)
	if err2 != nil {
        // fmt.Println(err2.Error())

        c.JSON(500, "error unmarshalling")
        return
    }	

    var out DoctorsList
    out.Doctors.Doctor = slices.Concat(outGrandOak.Doctors.Doctor, outPineValley.Doctors.Doctor)

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