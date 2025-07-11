package service

import (
    "errors" 
    "log"
    "strings"
    "net/http"
    "io"
    "encoding/json"

    "recipe-generator/internal/api/config"
)

type AnthropicApiClient struct {
    apiKey string
    anthropicUrl string
    body map[string]interface{} 
    aiPrompt string
}

func InitializeAnthropicClient(config *config.Config) (*AnthropicApiClient, error) {

    log.Println("Initializing Anthropic API client...")

    if config.AnthropicApiKey == "" {
	return nil, errors.New("Anthropic API key is required. Please set the ANTHROPIC_API_KEY environment variable")
    }

    if config.AnthropicApiUrl == "" {
        return nil, errors.New("Anthropic API url is required. Please set the ANTHROPIC_URL environment variable")
    }

    log.Println("Succsessfully initialized Anthropic API client")
    return &AnthropicApiClient {
	apiKey: config.AnthropicApiKey,
	anthropicUrl: config.AnthropicApiUrl,	
	body: make(map[string]interface{}),
	aiPrompt: "",

    }, nil
}

func (client *AnthropicApiClient) Post(config *config.Config, payload string) (map[string]interface{}, error) {

    // anthropic logic
    url := config.AnthropicApiUrl
    
    // Create the request
    log.Println("Creating request.")
    req, err := http.NewRequest("POST", url, strings.NewReader(payload))
    
    if err != nil {
    	log.Printf("Error creating request %s\n ", url)
	return nil, err
    }
    
    req.Header.Add("anthropic-version","2023-06-01")
    req.Header.Add("x-api-key", config.AnthropicApiKey)
    req.Header.Add("Content-Type", "application/json")
    
    log.Println("Sending request.")
    res, err := http.DefaultClient.Do(req)
    
    if err != nil {
	log.Printf("Error sending request to \n %n", url)
	return nil, err
    }
    
    //close request object after sending it
    defer res.Body.Close()
    
    // grab body from response.
    responseBody, err := io.ReadAll(res.Body)
    
    if err != nil {
    	log.Println("Error reading response body from reader in response.")
	return nil, err
    }
    
    //print out response, status code and reqponse body
    log.Printf("Status Code: %d\n", res.StatusCode)
    log.Println(res)

    var responseJson map[string]interface{}
    err = json.Unmarshal(responseBody, &responseJson)

    if err != nil {
	log.Println("Error unmarshalling response body to json.")
	return nil, err 
    }

    log.Println("Response body: ", responseJson)
    
    return responseJson, nil
}


