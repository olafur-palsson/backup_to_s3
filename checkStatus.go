/*
package main

import (
    "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"encoding/json"
	"os"
	"io/ioutil"
    "log"
	"net/http"
	"strconv"
	"fmt"
)

// should be environment vars
const BUCKET = "nottestbucket"
const REGION = "us-west-2"

var ACCESS_TOKEN = os.Getenv("GITLAB_ACCESS_TOKEN")
var gitlabProjectsUrl = "https://www.gitlab.com/api/v4/projects"
var listGitlabProjectsUrl = gitlabProjectsUrl + "?owned=true&access_token=" + ACCESS_TOKEN

var s3service *s3.S3
var awsSession *session.Session

type Project struct {
	Id int `json:id`
}

func main() {
	fmt.Println(ACCESS_TOKEN)
	projects := getProjects()
	for _, project := range projects {
		fmt.Println(project)
		checkProject(project)
	}
}

func getProjects() []Project {
	resp, err := http.Get(listGitlabProjectsUrl)
	if err != nil {
		fmt.Println("err1")
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err2")
		log.Fatal(err)
	}

	var projects []Project
	json.Unmarshal(body, &projects)
	return projects
}

func checkProject(project Project) {

	exportUrl := createCheckStatusUrl(project)
	fmt.Println(exportUrl)
	resp, err := http.Get(exportUrl)
	if err != nil {
		fmt.Println("err3")
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("err4")
		log.Fatal(err)
	}
	fmt.Println("Response:")
	fmt.Println(string(body))
}

func createCheckStatusUrl(project Project) string {
	projectId := strconv.Itoa(project.Id)
	return gitlabProjectsUrl + "/" + projectId + "/export?access_token=" + ACCESS_TOKEN
}

 */