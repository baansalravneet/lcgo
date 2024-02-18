package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Generate Leetcode Solution File")
	url := getInput("Enter the URL to the LeetCode problem: ")
	fileExtension := getInput("Enter file extension: ")
	languageSlug := getInput("Enter language slug (java, cpp, golang, etc): ")
	fileName, problemSlug, err := getFileName(url)
	if err != nil {
		fmt.Println("Error calling LeetCode:", err)
		os.Exit(1)
	}
	boilerPlate, err := getBoilerPlate(problemSlug, languageSlug)
	if err != nil {
		fmt.Println("Error getting starter code:", err)
		os.Exit(1)
	}
	file := fmt.Sprintf("%s.%s", fileName, fileExtension)
	err = createFile(file, boilerPlate)
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	fmt.Println("File created successfully:", file)
}
