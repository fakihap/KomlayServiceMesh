package main

import "github.com/gin-gonic/gin"
import "github.com/parnurzeal/gorequest"

import "encoding/json"
import "fmt"
import "slices"


var grandOakURL = "http://localhost:5000/grandOak/doctors/"
var pineValleyURL = "http://localhost:9090/pineValley/doctors"

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
        Send(`{"doctorType": "cardiologist"}`).
        End()

    fmt.Printf("%s --- %s \n", grandOakBody, pineValleyBody)

	var outGrandOak DoctorsList
    var outPineValley DoctorsList

	err1 := json.Unmarshal([]byte(grandOakBody), &outGrandOak)
	if err1 != nil {
        fmt.Println(err1.Error())
        return
    }	

    err2 := json.Unmarshal([]byte(pineValleyBody), &outPineValley)
	if err2 != nil {
        fmt.Println(err2.Error())
        return
    }	

    var out DoctorsList
    out.Doctors.Doctor = slices.Concat(outGrandOak.Doctors.Doctor, outPineValley.Doctors.Doctor)

    c.JSON(200, out)
}

func main() {
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

	r.GET("/doctors/:doctorType", fetchDoctors)


    r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}




// func fetchDoctors() ([]doctor, error) {
// 	var err error
//     var client = &http.Client{}
//     var data []doctor

//     request, err := http.NewRequest("GET", baseURL+"/users", nil)
//     if err != nil {
//         return nil, err
//     }

//     response, err := client.Do(request)
//     if err != nil {
//         return nil, err
//     }
//     defer response.Body.Close()

//     err = json.NewDecoder(response.Body).Decode(&data)
//     if err != nil {
//         return nil, err
//     }

//     return data,nil

// }