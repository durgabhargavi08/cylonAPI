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

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

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
	ImageTagFound string `json:"-"`
	MatchFound    bool   `json:"matchFound"`
	ProjectId     string `json:"-"`
	ProjectName   string `json:"-"`
}

// Arguements that are needed to pass:
// 1 - token
// 2 - projectId
// 3 - ImageId
// 4 - environment

var possibleEnvs = []string{production, staging, testinG, development}

func main() {
	if len(os.Args) < 5 {
		printErr("not enough arguements to make the call")
		return
	}
	env := getLowerEnv(os.Args[4])
	validator(new(http.Client), os.Args[1], os.Args[2], os.Args[3], env)
}

func getLowerEnv(env string) string {
	for i := range possibleEnvs {
		if strings.EqualFold(env, possibleEnvs[i]) {
			if i+1 == len(possibleEnvs) {
				return possibleEnvs[i]
			} else {
				return possibleEnvs[i+1]
			}
		}
	}
	return ""
}

func validator(client HTTPClient, token, projectId, imageId, environment string) bool {

	// move it config
	endpoint := "https://cylon-api.cisco.com/middleware/api/project/" + projectId

	request, err := http.NewRequest(http.MethodGet, endpoint, nil)
	if err != nil {
		printErr("error creating request (" + err.Error() + ")")
		return false
	}
	request.Header.Add("Authorization", "Bearer "+token)
	response, err := client.Do(request)
	if err != nil {
		printErr("error calling cylon (" + err.Error() + ")")
		return false
	}
	if response.StatusCode != http.StatusOK {
		printErr("project not found")
		return false
	}
	defer response.Body.Close()
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		printErr("error occured in reading response (" + err.Error() + ")")
		return false
	}
	object := new(Responsejson)
	err = json.Unmarshal(bytes, object)
	if err != nil {
		printErr("error occured in unmarshalling the response (" + err.Error() + ")")
		return false
	}
	output := Output{}
	output.ProjectId = projectId
	output.ProjectName = object.ProjectName

	if strings.EqualFold(environment, production) {
		output.ImageTagFound = object.Production.Values.Image.Tag
	}
	if strings.EqualFold(environment, development) {
		output.ImageTagFound = object.Development.Values.Image.Tag
	}
	if strings.EqualFold(environment, performace) {
		output.ImageTagFound = object.Performance.Values.Image.Tag
	}
	if strings.EqualFold(environment, testinG) {
		output.ImageTagFound = object.Testing.Values.Image.Tag
	}
	if strings.EqualFold(environment, staging) {
		output.ImageTagFound = object.Staging.Values.Image.Tag
	}

	output.MatchFound = output.ImageTagFound == imageId

	bytesOutput, err := json.MarshalIndent(output, " ", "\t")
	if err != nil {
		printErr("error occured in marshalling the output (" + err.Error() + ")")
		return false
	}

	fmt.Print(string(bytesOutput))
	return true
}

func printErr(str string) {
	default_message = `{"error": "` + str + `"}`
	fmt.Print(default_message)
}
