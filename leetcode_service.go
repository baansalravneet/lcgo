package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

var URL string = "https://leetcode.com/graphql"
//go:embed embed/*.json
var content embed.FS

func getBoilerPlate(problemSlug, languageSlug string) (string, error) {
	requestBody, err := getRequestBody(problemSlug, "embed/code_request.json")
	if err != nil {
		return "", err
	}
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", err
	}
	starterCode, err := getStarterCode(responseBody, languageSlug)
	if err != nil || starterCode == "" {
		return "", errors.New("error getting starter code")
	}
	return starterCode, err
}

func getStarterCode(responseBody []byte, languageSlug string) (string, error) {
	var outputData map[string]interface{}
	err := json.Unmarshal(responseBody, &outputData)
	if err != nil {
		return "", err
	}
	var starterCode string
	if jsonData, ok := outputData["data"].(map[string]interface{}); ok {
		if jsonQuestion, ok := jsonData["question"].(map[string]interface{}); ok {
			if codeSnippets, ok := jsonQuestion["codeSnippets"].([]interface{}); ok {
				for _, snippet := range codeSnippets {
					if codeSnippet, ok := snippet.(map[string]interface{}); ok {
						slug := codeSnippet["langSlug"].(string)
						if slug == languageSlug {
							starterCode = codeSnippet["code"].(string)
						}
					}
				}
			}
		}
	}
	return starterCode, nil
}

func getProblemSlug(url string) string {
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		return parts[4]
	}
	return ""
}

func getRequestBody(problemSlug string, request_file string) ([]byte, error) {
	requestBody, err := content.ReadFile(request_file)
	if err != nil {
		return nil, err
	}
	var data map[string]interface{}
	err = json.Unmarshal(requestBody, &data)
	if err != nil {
		return nil, err
	}
	data["variables"].(map[string]interface{})["titleSlug"] = problemSlug
	postRequest, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return postRequest, nil
}

func getQuestionId(responseBody []byte) (string, string) {
	var outputData map[string]interface{}
	err := json.Unmarshal(responseBody, &outputData)
	if err != nil {
		return "", ""
	}
	var questionId string
	var questionTitle string
	if jsonData, ok := outputData["data"].(map[string]interface{}); ok {
		if jsonQuestion, ok := jsonData["question"].(map[string]interface{}); ok {
			if qId, ok := jsonQuestion["questionId"].(string); ok {
				questionId = qId
			}
			if title, ok := jsonQuestion["title"].(string); ok {
				questionTitle = title
			}
		}
	}
	return questionId, questionTitle
}

func getFileName(url string) (string, string, error) {
	problemSlug := getProblemSlug(url)
	if problemSlug == "" {
		return "", "", errors.New("error getting problem slug from the url")
	}
	requestBody, err := getRequestBody(problemSlug, "embed/request.json")
	if err != nil {
		fmt.Println(err)
		return "", "", errors.New("error forming request body")
	}
	response, err := http.Post(URL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}
	questionId, questionTitle := getQuestionId(responseBody)
	if questionId == "" || questionTitle == "" {
		return "", "", errors.New("error getting question id and title")
	}
	return fmt.Sprintf("%s. %s", questionId, questionTitle), problemSlug, nil
}
