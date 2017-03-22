// handlers.article.go

package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	userName, _ := c.GetQuery("user.name")
	render(c, gin.H{"title": "Home Page", "user": userName}, "index.html")
}

func showApplicationCreationPage(c *gin.Context) {
	userName, _ := c.GetQuery("user.name")
	render(c, gin.H{"title": "Create New Application", "user": userName}, "create-application.html")
}

func getApplicationsPage(c *gin.Context) {
	userName, _ := c.GetQuery("user.name")
	var applicationList = getApplications(c)
	render(c, gin.H{"title": "Applications", "applications": applicationList, "user": userName}, "view-applications.html")
}

func getApplications(c *gin.Context) []application {
	// var applicationList = getAllApplications()
	var applicationList = []application{}
	jwttoken, err := c.Cookie("hadoop-jwt")

	if err != nil && jwttoken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// make rest Call
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	baseURL := "https://localhost:8443/gateway/sandbox/webhdfs/v1/user/guest/example/"
	url := fmt.Sprintf("%s%s", baseURL, "?OP=LISTSTATUS")
	req, err := http.NewRequest("GET", url, nil)
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "hadoop-jwt", Value: jwttoken, Expires: expiration}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
	fileList := getFileList(resp)
	for _, v := range fileList {
		url = fmt.Sprintf("%s%s%s", baseURL, v, "?OP=OPEN")
		req, err = http.NewRequest("GET", url, nil)
		req.AddCookie(&cookie)
		resp, err = client.Do(req)
		if err != nil {
			fmt.Printf("Error : %s", err)
		}
		applicationList = append(applicationList, getApplication(resp))
	}

	return applicationList
}

func getApplication(response *http.Response) application {
	defer response.Body.Close()
	var application application
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&application)
	return application
}

func getFileList(response *http.Response) []string {
	defer response.Body.Close()
	var fileList []string
	var dat map[string]interface{}
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&dat)

	fileStatuses := dat["FileStatuses"].(map[string]interface{})
	if fileStatuses != nil {
		files := fileStatuses["FileStatus"].([]interface{})
		if files != nil {
			for _, v := range files {
				file := v.(map[string]interface{})
				path := file["pathSuffix"].(string)
				if path != "" {
					fileList = append(fileList, path)
				}
			}
		}
	}
	return fileList
}

func submitApplication(application *application, c *gin.Context) {
	// userName, _ := c.GetQuery("user.name")
	content, err := json.Marshal(application)
	if err != nil {
		log.Fatal(err)
	}
	jwttoken, err := c.Cookie("hadoop-jwt")

	if err != nil && jwttoken == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
	}

	// make rest Call
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	filename := strconv.Itoa(application.ID)

	url := "https://localhost:8443/gateway/sandbox/webhdfs/v1/user/guest/example/" + filename + "?op=CREATE"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(content))
	// req.SetBasicAuth("guest", "guest-password")
	expiration := time.Now().Add(365 * 24 * time.Hour)
	cookie := http.Cookie{Name: "hadoop-jwt", Value: jwttoken, Expires: expiration}
	req.AddCookie(&cookie)
	_, err = client.Do(req)
	if err != nil {
		fmt.Printf("Error : %s", err)
	}
}

func createApplication(c *gin.Context) {
	name := c.PostForm("name")
	address := c.PostForm("address")
	loan, _ := strconv.Atoi(c.PostForm("loan"))

	a := createNewApplication(name, address, loan)
	submitApplication(a, c)
	render(c, gin.H{
		"title":   "Submission Successful",
		"payload": a}, "submission-successful.html")

}
