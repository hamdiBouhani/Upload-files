package utils

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type TokenResp struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetSsoToken() (string, error) {

	clientID := os.Getenv("MEERA_CLIENT_ID")
	clientSecret := os.Getenv("MEERA_CLIENT_SECRET")
	oidcAddr := os.Getenv("MEERA_CLIENT_OIDC_URI")

	// prep payload
	payload := strings.NewReader(
		"grant_type=client_credentials&client_id=" + clientID + "&client_secret=" + clientSecret + "&scope=openid email groups profile offline_access",
	)

	var tokenResp TokenResp
	reqToken, err := http.NewRequest("POST", oidcAddr, payload)
	if err != nil {
		return "", err
	}

	reqToken.Header.Add("content-type", "application/x-www-form-urlencoded")
	res, err := http.DefaultClient.Do(reqToken)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	log.Printf("%s", body)
	err = json.Unmarshal(body, &tokenResp)
	if err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}
