// handlers.article_test.go

package main

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestApplicationListJSON(t *testing.T) {
	response := []byte(`{"FileStatuses":{"FileStatus":[{"accessTime":1489768158179,"blockSize":134217728,"childrenNum":0,"fileId":16841,"group":"hdfs","length":40,"modificationTime":1489768158204,"owner":"guest","pathSuffix":"asdga","permission":"755","replication":3,"storagePolicy":0,"type":"FILE"},{"accessTime":1489828691082,"blockSize":134217728,"childrenNum":0,"fileId":17394,"group":"hdfs","length":59,"modificationTime":1489828691101,"owner":"guest","pathSuffix":"s12a","permission":"755","replication":3,"storagePolicy":0,"type":"FILE"},{"accessTime":1489721080736,"blockSize":134217728,"childrenNum":0,"fileId":16548,"group":"hdfs","length":65,"modificationTime":1489721080788,"owner":"guest","pathSuffix":"test","permission":"755","replication":3,"storagePolicy":0,"type":"FILE"}]}}`)

	var dat map[string]interface{}
	if err := json.Unmarshal(response, &dat); err != nil {
		panic(err)
	}
	fileStatuses := dat["FileStatuses"].(map[string]interface{})
	if fileStatuses != nil {
		files := fileStatuses["FileStatus"].([]interface{})
		if files != nil {
			for _, v := range files {
				file := v.(map[string]interface{})
				fmt.Println(file["pathSuffix"])
			}
		}
	}
	// r := getRouter(true)

	// // Define the route similar to its definition in the routes file
	// r.GET("/applications", getApplications)

	// // Create a request to send to the above route
	// req, _ := http.NewRequest("GET", "/applications", nil)
	// req.Header.Add("Accept", "application/json")

	// testHTTPResponse(t, r, req, func(w *httptest.ResponseRecorder) bool {
	// 	// Test that the http status code is 200
	// 	statusOK := w.Code == http.StatusOK

	// 	// Test that the response is JSON which can be converted to
	// 	// an array of application structs
	// 	p, err := ioutil.ReadAll(w.Body)
	// 	if err != nil {
	// 		return false
	// 	}
	// 	var applications []application
	// 	err = json.Unmarshal(p, &applications)

	// 	return err == nil && len(applications) >= 2 && statusOK
	// })
}
