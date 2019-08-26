
package main

import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/aws/request"
	"github.com/aws/aws-sdk-go/service/s3"

	// "bytes"
	"encoding/json"
	"os"
    "os/exec"
	"io/ioutil"
    "log"
	"time"
	"net/http"
	"strconv"
	"fmt"

	"./zipFolder"
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
	Id int `json:"id"`
	UrlToRepo string `json:"web_url"`
}

func main() {
	fmt.Println(ACCESS_TOKEN)
	projects := getProjects()
	for _, project := range projects {
		exportProject2(project)
	}
}

func getProjects() []Project {
	resp, err := http.Get(listGitlabProjectsUrl)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(body))

	var projects []Project
	json.Unmarshal(body, &projects)
	return projects
}

func exportProject2(project Project) {
	pullProject(project)
	zipProject(project)
	upload(project)
}

/*
func exportProject(project Project) {
	awsPresignedUrl, _ := createPresignedUrlFor(project)
	fmt.Println(awsPresignedUrl)

	payload := createExportRequestPayload(awsPresignedUrl)
	exportUrl := createExportUrl(project)
	req, err := http.NewRequest("POST", exportUrl, payload)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	req.Header.Add("PRIVATE-TOKEN", ACCESS_TOKEN)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Response:")
	fmt.Println(string(body))
}


 */
func upload(project Project) {
	data, err := os.Open(strconv.Itoa(project.Id))
	if err != nil {
		log.Fatal(err)
	}

	defer data.Close()
	url, err := createPresignedUrlFor(project)
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("PUT", url, data)
	if err != nil {
		log.Fatal(err)
	}
	client := http.Client{}
	client.Do(req)
}

func pullProject(project Project) {
	fmt.Println(project)
	err := exec.Command("git", "clone", project.UrlToRepo, strconv.Itoa(project.Id))
	if err != nil {
		log.Fatal(err)
	}
}

func zipProject(project Project) {
	zipFolder.ZipWriter(strconv.Itoa(project.Id))
}

/*
func createExportRequestPayload(presignedUrl string) *bytes.Buffer {
	payload := 
		"upload[http_method]=PUT" + "\n" +
		"upload[url]=" + presignedUrl
	fmt.Println(payload)
	return bytes.NewBuffer([]byte(payload))
}

 */

func createPresignedUrlFor(project Project) (string, error) {
	// Create S3 service client
	awsSession := createSessionWithCredentialsAndRegion()
	s3service = s3.New(awsSession)
	req := createS3Request(strconv.Itoa(project.Id))
	return req.Presign(15 * time.Minute)
}

func createSessionWithCredentialsAndRegion() *session.Session {
    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(REGION)},
	)
	if err != nil {
		log.Fatal(err)
	}
	return sess
}

func createS3Request(projectId string) *request.Request {
    req, _ := s3service.PutObjectRequest(&s3.PutObjectInput{
        Bucket: aws.String(BUCKET),
        Key:    aws.String(projectId),
	})
	return req
}

/*
func createExportUrl(project Project) string {
	projectId := strconv.Itoa(project.Id)
	return gitlabProjectsUrl + "/" + projectId + "/export"
}

 */