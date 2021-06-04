package altium

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

const uploadUrl = "https://viewer.altium.com/api/widget/set"

func newUploadProjectRequest(projectFile io.Reader) (*http.Request, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("File", "hidden")
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, projectFile); err != nil {
		return nil, err
	}
	if err := writer.WriteField("HideSrc", "true"); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", uploadUrl, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req, nil
}

type responseJson struct {
	DesignId   string   `json:"designId"`
	Modules    []string `json:"modules"`
	FaultCode  int      `json:"faultCode"`
	DesignType int      `json:"designType"`
	Status     string   `json:"status"`
	Message    string   `json:"message"`
}

func (res *responseJson) Error() string {
	return fmt.Sprintf("%s: %s", res.Status, res.Message)
}

func UploadProject(projectFile io.Reader) (id string, err error) {
	clt := &http.Client{}
	req, err := newUploadProjectRequest(projectFile)
	if err != nil {
		return "", err
	}
	res, err := clt.Do(req)
	if err != nil {
		return "", err
	}
	var resJson responseJson
	if err := json.NewDecoder(res.Body).Decode(&resJson); err != nil {
		return "", err
	}
	if resJson.Status == "Error" {
		return "", &resJson
	}
	return resJson.DesignId, nil
}
