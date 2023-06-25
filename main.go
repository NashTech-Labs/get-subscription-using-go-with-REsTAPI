package main

 

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os/exec"
    "strings"
)

 

type Subscription struct {
    ID          string `json:"id"`
    DisplayName string `json:"displayName"`
    State       string `json:"state"`
    TenantID    string `json:"tenantId"`
    QuotaID     string `json:"quotaId"`
    // Include additional fields as needed
}

 

func main() {
    subscriptionID := "<subscription-id>"
    apiVersion := "2020-01-01"
    url := fmt.Sprintf("https://management.azure.com/subscriptions/%s?api-version=%s", subscriptionID, apiVersion)

 

    accessToken, err := getAccessToken()
    if err != nil {
        fmt.Printf("Oh, We are failed to get access token: %s\n", err.Error())
        return
    }

 

    subscriptionJSON, err := getSubscriptionDetails(url, accessToken)
    if err != nil {
        fmt.Printf("Oh, We are failed to get subscription details: %s\n", err.Error())
        return
    }

 

    sample_data, err := printSubscriptionDetails(subscriptionJSON)

    fmt.Println("\nSubscription Details:")
    fmt.Printf("Subscription ID: %s\n", strings.TrimPrefix(sample_data.ID, "/subscriptions/")    )
    // fmt.Println("This is the id without trim")
    // fmt.Printf("Subscription ID: %s\n", subscriptionID)
    // fmt.Println("This is the id without trim with the sample_data")
    // fmt.Printf("Subscription ID: %s\n", sample_data.ID)
    fmt.Printf("Display Name: %s\n", sample_data.DisplayName)
    fmt.Printf("State: %s\n", sample_data.State)
    fmt.Printf("Tenant ID: %s\n", sample_data.TenantID)
    // fmt.Printf("Quota ID: %s\n", sample_data.QuotaID)
}

 

func getAccessToken() (string, error) {
    cmd := exec.Command("az", "account", "get-access-token", "--query", "accessToken", "--output", "tsv")
    output, err := cmd.Output()
    if err != nil {
        return "", err
    }

 

    return strings.TrimSpace(string(output)), nil
}

 

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