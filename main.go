package main

 // import the required package and its dependencies
import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "strings"
)

// define the structure which will be used in the request

type Subscription struct {
    ID          string `json:"id"`
    DisplayName string `json:"displayName"`
    State       string `json:"state"`
    TenantID    string `json:"tenantId"`
    QuotaID     string `json:"quotaId"`
    // Include additional fields as needed
}

 
// main function

func main() {
    subscriptionID := "<subscription-id>"
    apiVersion := "2020-01-01"
    url := fmt.Sprintf("https://management.azure.com/subscriptions/%s?api-version=%s", subscriptionID, apiVersion)

 
    // calling the getAccessToken method to get the access token
    accessToken, err := getAccessToken()
    if err != nil {
        fmt.Printf("Oh, We are failed to get access token: %s\n", err.Error())
        return
    }

 
    // calling the getSubscriptionDetails method to get the subscription details.
    subscriptionJSON, err := getSubscriptionDetails(url, accessToken)
    if err != nil {
        fmt.Printf("Oh, We are failed to get subscription details: %s\n", err.Error())
        return
    }

 
    // calling the printSubscriptionDetails method and assign the values to the sample_data variable

    sample_data, err := printSubscriptionDetails(subscriptionJSON)

    fmt.Println("\nSubscription Details:")
    fmt.Printf("Subscription ID: %s\n", strings.TrimPrefix(sample_data.ID, "/subscriptions/")    )

    fmt.Printf("Display Name: %s\n", sample_data.DisplayName)
    fmt.Printf("State: %s\n", sample_data.State)
    fmt.Printf("Tenant ID: %s\n", sample_data.TenantID)
}

 
// function to get the access token 
func getAccessToken() (string, error) {
    cmd := exec.Command("az", "account", "get-access-token", "--query", "accessToken", "--output", "tsv")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

 

    return strings.TrimSpace(string(output)), nil
}

// function to fetch the subcription information

func getSubscriptionDetails(url, accessToken string) ([]byte, error) {
    client := &http.Client{}
    req, err := http.NewRequest(http.MethodGet, url, nil)
    if err != nil {
        return nil, err
    }

 

    req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))

 

    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

 

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return nil, err
    }

 

    return body, nil
}

 
// print the values of subcription into a normal text from the JSON response
func printSubscriptionDetails(subscriptionJSON []byte) (Subscription, error) {
    fmt.Println("Complete JSON Response:")
    fmt.Println(string(subscriptionJSON))

 

    var subscription Subscription
    err := json.Unmarshal(subscriptionJSON, &subscription)
    if err != nil {
        fmt.Printf("Failed to unmarshal JSON response: %s\n", err.Error())
        return Subscription{}, err
    }
    return subscription, err

 

}