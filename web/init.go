package web

import (
    "fmt"
    "net/http"
    //"io/ioutil"
    "bytes"
    "os"
    "strings"
    "io"
    "encoding/json"
    "bufio"
)

func Test() {
    fmt.Println("In main")
    
    fmt.Println("Beginning Setup Process")
    fmt.Println("Visit https://www.strava.com/settings/api and create an application, then copy and paste the Client ID and Client Secret into the following prompts:")
    
    fmt.Printf("Client ID: ")
	in := bufio.NewReader(os.Stdin)
	clientID, _ := in.ReadString('\n')
	clientID = strings.ReplaceAll(clientID, "\n", "")
	
	fmt.Printf("Client Secret: ")
	in = bufio.NewReader(os.Stdin)
	clientSecret, _ := in.ReadString('\n')
	clientSecret = strings.ReplaceAll(clientSecret, "\n", "")
	
	fmt.Println("Go the the following url and grant access, then copy the code from the URL and paste here:")
	fmt.Printf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http%%3A%%2F%%2Flocalhost&scope=activity:read_all,activity:write&state=mystate&approval_prompt=force\n", clientID)
	
	fmt.Printf("Code: ")
	in = bufio.NewReader(os.Stdin)
	code, _ := in.ReadString('\n')
	code = strings.ReplaceAll(code, "\n", "")
	
	// Get refresh and access tokens and save to file
	
	values := map[string]string{"client_id": clientID, "client_secret": clientSecret, "code": code, "grant_type": "authorization_code"}
	
	
	jsonValue, _ := json.Marshal(values)
	
	res, err := http.Post("https://www.strava.com/oauth/token", "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
	    panic(err)
	}
    defer res.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode >= 400 && res.StatusCode <= 500 {
		fmt.Println("Error response. Status Code: ", res.StatusCode)
		os.Exit(1)
	}

	fmt.Println("Response:", string(responseBody))
	var finalJSON map[string]interface{}
	if err := json.Unmarshal([]byte(responseBody), &finalJSON); err != nil {
		panic(err)
	}
	var refreshToken string = finalJSON["refresh_token"].(string)
	var accessToken string = finalJSON["access_token"].(string)
	
	err = os.WriteFile("save.txt", []byte(refreshToken + "\n" + accessToken), 0666)
}
