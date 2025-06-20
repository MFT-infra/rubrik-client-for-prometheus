package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/rubrikinc/rubrik-sdk-for-go/rubrikcdm"
)

type serviceAccountResponse struct {
	Token string `json:"token"`
	TTL   int    `json:"ttl"`
}

func connectRubrik(cfg *Config) (*rubrikcdm.Credentials, error) {
	if cfg.RubrikIP == "" {
		return nil, errors.New("rubrik_ip is required")
	}

	if cfg.ServiceID != "" && cfg.ServiceSecret != "" {
		token, err := getServiceAccountToken(cfg.RubrikIP, cfg.ServiceID, cfg.ServiceSecret)
		if err != nil {
			return nil, fmt.Errorf("service account auth failed: %w", err)
		}
		log.Println("Using service account for authentication")
		return rubrikcdm.ConnectAPIToken(cfg.RubrikIP, token), nil
	}

	if cfg.ApiToken != "" {
		log.Println("Using API token for authentication")
		return rubrikcdm.ConnectAPIToken(cfg.RubrikIP, cfg.ApiToken), nil
	}

	if cfg.Username != "" && cfg.Password != "" {
		log.Println("Using username/password for authentication")
		return rubrikcdm.Connect(cfg.Username, cfg.Password, cfg.RubrikIP), nil
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
