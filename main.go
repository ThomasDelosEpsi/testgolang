package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	tokenURL  = "https://cloud.uipath.com/identity_/connect/token"
	uipathAPI = "https://cloud.uipath.com/lyrecomanagement/DefaultTenant/orchestrator_/odata/Jobs"
	orgUnitID = "1966354"
)

var (
	clientID     = os.Getenv("UIPATH_CLIENT_ID")
	clientSecret = os.Getenv("UIPATH_CLIENT_SECRET")
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
}

func getAccessToken() (string, error) {
	body := "grant_type=client_credentials&client_id=" + clientID + "&client_secret=" + clientSecret + "&scope=OR.Jobs.Read"
	resp, err := http.Post(tokenURL, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(body)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

func getUiPathJobs(accessToken string) ([]byte, error) {
	req, err := http.NewRequest("GET", uipathAPI, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-UIPATH-OrganizationUnitId", orgUnitID)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func main() {
	accessToken, err := getAccessToken()
	if err != nil {
		fmt.Println("Erreur lors de la récupération du token:", err)
		return
	}

	fmt.Println("Access Token:", accessToken)

	jobsData, err := getUiPathJobs(accessToken)
	if err != nil {
		fmt.Println("Erreur lors de la récupération des jobs:", err)
		return
	}

	filePath := "C:/Users/tdelos/OneDrive - LYRECO MANAGEMENT/Bureau/uipath_jobs.txt"
	err = os.WriteFile(filePath, jobsData, 0644)
	if err != nil {
		fmt.Println("Erreur lors de l'écriture du fichier:", err)
		return
	}

	fmt.Println("Résultats enregistrés dans", filePath)
}
