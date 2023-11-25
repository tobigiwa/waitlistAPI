package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {

	authURL := "https://accounts.zoho.com/oauth/v2/auth?scope=accounts&client_id=1000.B2IPYVZ1LWONKRLLKF5OPL30EPTJRR&response_type=code&access_type=online&redirect_uri=127.0.0.1:8090"
	zoid := "837789618"
	allOrgUserDetailsEndpoint := "http://mail.zoho.com/api/organization/" + zoid + "/accounts"

	// sendEndpoint := "https://mail.zoho.com/api/accounts/<accountId>/messages"

	resp, err := http.Get(allOrgUserDetailsEndpoint)
	if err != nil {
		log.Fatalf("error sending request %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed to get organization user details. Status code:", resp.StatusCode)
		// fmt.Println(string(resp.body))
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response body: %v", err)
	}

	fmt.Println(string(body))
}
