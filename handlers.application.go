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
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func showIndexPage(c *gin.Context) {
	m := createParamMap(c)
	m["title"] = "Home Page"
	render(c, m, "index.html")
}

func showApplicationCreationPage(c *gin.Context) {
	m := createParamMap(c)
	m["title"] = "Create New Application"
	render(c, m, "create-application.html")
}

func getApplicationsPage(c *gin.Context) {
	m := createParamMap(c)
	m["title"] = "Applications"
	var applicationList = getApplications(c)
	m["applications"] = applicationList
	render(c, m, "view-applications.html")
}

func getApplications(c *gin.Context) []application {
	host, tenant := getHostAndTenant(c.Request)
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
	baseURL := fmt.Sprintf("https://%s/gateway/%s/webhdfs/v1/%s/applications/", host, tenant, tenant)

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
		} else {
			applicationList = append(applicationList, getApplication(resp))
		}

	}

	return applicationList
}

func getHostAndTenant(request *http.Request) (string, string) {
	host := request.Header.Get("X-Forwarded-Host")
	hostparts := strings.Split(host, ".")
	tenant := hostparts[1]
	return host, tenant
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
	host, tenant := getHostAndTenant(c.Request)
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
	baseURL := fmt.Sprintf("https://%s/gateway/%s/webhdfs/v1/%s/applications/", host, tenant, tenant)

	url := fmt.Sprintf("%s%s%s", baseURL, filename, "?op=CREATE")

	// url := "https://localhost:8443/gateway/unwise/webhdfs/v1/user/guest/example/" + filename + "?op=CREATE"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(content))
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
	m := createParamMap(c)
	m["title"] = "Submission Successful"
	render(c, m, "submission-successful.html")
}

func createParamMap(c *gin.Context) map[string]interface{} {
	userName, _ := c.GetQuery("user.name")
	userParts := strings.Split(userName, "_")
	userName = userParts[0]
	_, tenant := getHostAndTenant(c.Request)
	banner := "Acme Loans"
	bannerLead := "Generic loaning company"
	if tenant == "goodloans" {
		banner = "Good Loans Lending Institute"
		bannerLead = "Just a really good loan company"
	} else if tenant == "unwise" {
		banner = "Unwise Lending"
		bannerLead = "Lending to people we like"
	}
	return gin.H{"user": userName, "banner": banner, "bannerLead": bannerLead}
}
