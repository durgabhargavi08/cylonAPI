package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

var default_message string

const (
	production  = "production"
	staging     = "staging"
	testinG     = "testing"
	development = "development"
	performace  = "performance"
)

type Responsejson struct {
	ID              string      `json:"id"`
	RevisionUUID    string      `json:"revisionUuid"`
	ProjectName     string      `json:"projectName"`
	DisplayName     string      `json:"displayName"`
	ApplicationType string      `json:"applicationType"`
	Performance     Development `json:"performance"`
	Development     Development `json:"development"`
	Testing         Development `json:"testing"`
	Staging         Development `json:"staging"`
	Production      Development `json:"production"`
	UpdatedAt       string      `json:"updatedAt"`
	UpdatedBy       string      `json:"updatedBy"`
}

type Development struct {
	LastDeploymentTimestamp string            `json:"lastDeploymentTimestamp"`
	UpdatedAt               string            `json:"updatedAt"`
	Values                  DevelopmentValues `json:"values"`
}

type DevelopmentValues struct {
	Image Image `json:"image"`
}

type Image struct {
	Tag string `json:"tag"`
}

type Output struct {
	imageTagFound string
	matchFound    bool
	projectId     string
	projectName   string
}

// Arguements that are needed to pass:
// 1 - token
// 2 - projectId
// 3 - ImageId
// 4 - environment

func main() {
	if len(os.Args) < 5 {
		printErr("not enough arguements to make the call")
		return
	}
	validator(os.Args[1], os.Args[2], os.Args[3], os.Args[4])
	return
}

func validator(token, projectId, imageId, environment string) {

	endpoint := "https://cylon-api.cisco.com/middleware/api/project/" + projectId

	client := new(http.Client)
	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		printErr("error creating request (" + err.Error() + ")")
		return
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		printErr("error calling cylon (" + err.Error() + ")")
		return
	}
	if response.StatusCode != http.StatusOK {
		printErr("project not found")
		return
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		printErr("error occured in reading response (" + err.Error() + ")")
		return
	}
	object := new(Responsejson)
	err = json.Unmarshal(bytes, object)
	if err != nil {
		printErr("error occured in unmarshalling the response (" + err.Error() + ")")
		return
	}
	output := new(Output)
	output.projectId = projectId
	output.projectName = object.ProjectName

	if strings.EqualFold(environment, production) {
		output.imageTagFound = object.Production.Values.Image.Tag
	}
	if strings.EqualFold(environment, development) {
		output.imageTagFound = object.Development.Values.Image.Tag
	}
	if strings.EqualFold(environment, performace) {
		output.imageTagFound = object.Performance.Values.Image.Tag
	}
	if strings.EqualFold(environment, testinG) {
		output.imageTagFound = object.Testing.Values.Image.Tag
	}
	if strings.EqualFold(environment, staging) {
		output.imageTagFound = object.Staging.Values.Image.Tag
	}

	output.matchFound = output.imageTagFound == imageId

	bytes, err = json.MarshalIndent(output, " ", "\t")
	if err != nil {
		printErr("error occured in marshalling the output (" + err.Error() + ")")
		return
	}

	fmt.Print(string(bytes))
	return
}

func printErr(str string) {
	default_message = `{"error": ` + str + `}`
	fmt.Print(default_message)
}
