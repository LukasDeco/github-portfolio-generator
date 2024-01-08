package netlify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	netlifyAPIURL    = "https://api.netlify.com/api/v1"
	netlifyDeployURL = netlifyAPIURL + "/sites"
)

func DeployNetlify(htmlFile string) error {
	// Get file paths in the current directory
	filePaths, err := filepath.Glob(htmlFile)
	if err != nil {
		fmt.Println("Error finding HTML files:", err)
		return err
	}

	// Create a map to store file shasums
	fileShasums := make(map[string]string)

	// Calculate shasum for each file
	for _, filePath := range filePaths {
		shasum, err := calculateShasum(filePath)
		if err != nil {
			fmt.Printf("Error calculating shasum for %s: %v\n", filePath, err)
			return err
		}
		fileShasums[filePath] = shasum
	}

	// Create a deploy request payload
	deployPayload := map[string]interface{}{
		"files": fileShasums,
	}

	// Make a POST request to create a deploy
	siteID, err := createSite(deployPayload)
	if err != nil {
		fmt.Println("Error creating deploy:", err)
		return err
	}

	// Make a POST request to create a file upload
	deployID, err := createDeploy(siteID, fileShasums)
	if err != nil {
		fmt.Println("Error creating deploy:", err)
		return err
	}

	// Upload each file to the deploy
	for _, filePath := range filePaths {
		err := uploadFile(deployID, filePath)
		if err != nil {
			fmt.Printf("Error uploading %s: %v\n", filePath, err)
			return err
		}
	}

	fmt.Println("Deployment successful!")
	return nil
}

func calculateShasum(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", content), nil
}

type CreateSiteResponse struct {
	SiteID string `json:"site_id"`
	// Add other fields from the response as needed
}

func createSite(payload map[string]interface{}) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", netlifyDeployURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("NETLIFY_ACCESS_TOKEN"))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return "", err
		}
		fmt.Println(string(body))
		return "", fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	// Parse the response JSON
	var createDeployResponse CreateSiteResponse
	err = json.NewDecoder(resp.Body).Decode(&createDeployResponse)
	if err != nil {
		return "", err
	}

	return createDeployResponse.SiteID, nil
}

type CreateDeployRes struct {
	ID string `json:"id"`
	// Add other fields from the response as needed
}

func createDeploy(siteID string, fileShasums map[string]string) (string, error) {
	// Create a deploy request payload
	deployPayload := map[string]interface{}{
		"files": fileShasums,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(deployPayload)
	if err != nil {
		return "", err
	}

	// Make a POST request to create a deploy
	client := &http.Client{}
	req, err := http.NewRequest("POST", fmt.Sprintf("https://api.netlify.com/api/v1/sites/%s/deploys", siteID), bytes.NewBuffer(payloadBytes))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("NETLIFY_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response JSON
	var createDeployResponse CreateDeployRes
	err = json.NewDecoder(resp.Body).Decode(&createDeployResponse)
	if err != nil {
		return "", err
	}

	return createDeployResponse.ID, nil
}

const netlifyUploadURL = "https://api.netlify.com/api/v1/deploys/%s/files/%s"

func uploadFile(deployID, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Build the upload URL
	uploadURL := fmt.Sprintf(netlifyUploadURL, deployID, filepath.Base(filePath))

	// Create a PUT request with the file binary
	req, err := http.NewRequest("PUT", uploadURL, file)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("NETLIFY_ACCESS_TOKEN"))
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode > 299 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println(string(body))
		return fmt.Errorf("received non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
