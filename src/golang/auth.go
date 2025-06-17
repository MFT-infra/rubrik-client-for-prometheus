// auth.go
package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

var (
	rubrikIP      = flag.String("rubrik-ip", os.Getenv("RUBRIK_CDM_NODE_IP"), "Rubrik CDM node IP address")
	username      = flag.String("username", os.Getenv("RUBRIK_CDM_USERNAME"), "Username for Rubrik")
	password      = flag.String("password", os.Getenv("RUBRIK_CDM_PASSWORD"), "Password for Rubrik")
	apiToken      = flag.String("api-token", os.Getenv("RUBRIK_CDM_API_TOKEN"), "Rubrik API Token")
	serviceID     = flag.String("service-id", os.Getenv("RUBRIK_SERVICE_ID"), "Service Account ID")
	serviceSecret = flag.String("service-secret", os.Getenv("RUBRIK_SERVICE_SECRET"), "Service Account Secret")
)

type serviceAccountResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}

func connectRubrik() (*rubrikcdm.Credentials, error) {
	if *rubrikIP == "" {
		return nil, errors.New("rubrik-ip is required")
	}

	if *serviceID != "" && *serviceSecret != "" {
		token, err := getServiceAccountToken(*rubrikIP, *serviceID, *serviceSecret)
		if err != nil {
			return nil, fmt.Errorf("service account auth failed: %w", err)
		}
		log.Println("Using service account for authentication")
		return rubrikcdm.ConnectAPIToken(*rubrikIP, token), nil
	}

	if *apiToken != "" {
		log.Println("Using API token for authentication")
		return rubrikcdm.ConnectAPIToken(*rubrikIP, *apiToken), nil
	}

	if *username != "" && *password != "" {
		log.Println("Using username/password for authentication")
		return rubrikcdm.Connect(*username, *password, *rubrikIP), nil
	}

	return nil, errors.New("no valid authentication method provided")
}

func getServiceAccountToken(ip, id, secret string) (string, error) {
	url := fmt.Sprintf("https://%s/api/service_account/session", ip)
	payload := map[string]interface{}{
		"serviceAccountId":  id,
		"secret":            secret,
		"sessionTtlMinutes": 1440,
	}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("unexpected status code %d: %s", resp.StatusCode, string(data))
	}

	var saResp serviceAccountResponse
	if err := json.NewDecoder(resp.Body).Decode(&saResp); err != nil {
		return "", err
	}

	return saResp.Token, nil
}
