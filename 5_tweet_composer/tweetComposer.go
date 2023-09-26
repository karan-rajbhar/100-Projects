package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/dghubble/oauth1"
	twauth "github.com/dghubble/oauth1/twitter"
)

// In your terminal please set your environment variables by running the following lines of code.
// export 'TWITTER_API_KEY'='<your_consumer_key>'
// export 'TWITTER_API_SECRET_KEY'='<your_consumer_secret>'
const outOfBand = "oob"

var config oauth1.Config

type Credentials struct {
	ConsumerKey       string
	ConsumerSecret    string
	AccessToken       string
	AccessTokenSecret string
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Error: Please provide tweet.")
		os.Exit(1)
	}

	consumerKey := os.Getenv("TWITTER_API_KEY")
	consumerSecret := os.Getenv("TWITTER_API_SECRET_KEY")

	if consumerKey == "" || consumerSecret == "" {
		fmt.Println("Please set your CONSUMER_KEY and CONSUMER_SECRET environment variables.")
		return
	}

	config := oauth1.Config{
		ConsumerKey:    consumerKey,
		ConsumerSecret: consumerSecret,
		CallbackURL:    outOfBand,
		Endpoint:       twauth.AuthorizeEndpoint,
	}

	// Get request token
	requestToken, requestSecret, err := config.RequestToken()
	if err != nil {
		fmt.Println("Failed to get request token:", err)
		return
	}

	fmt.Println("Got OAuth token:", requestToken)

	// Get authorization
	authorizationURL, err := config.AuthorizationURL(requestToken)
	if err != nil {
		fmt.Println("Failed to get authorization URL:", err)
		return
	}

	fmt.Println("Please go here and authorize:", authorizationURL.String())

	fmt.Print("Paste the PIN here: ")
	reader := bufio.NewReader(os.Stdin)
	verifier, _ := reader.ReadString('\n')
	verifier = strings.TrimSpace(verifier)

	// Get the access token
	accessToken, accessSecret, err := config.AccessToken(requestToken, requestSecret, verifier)
	if err != nil {
		fmt.Println("Failed to get access token:", err)
		return
	}

	token := oauth1.NewToken(accessToken, accessSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	// Make the request
	// payload := `{"text": "Hello world!"}`

	type payload struct {
		Text string `json:"text"`
	}
	data := &payload{Text: args[1]}

	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	resp, err := httpClient.Post("https://api.twitter.com/2/tweets", "application/json", strings.NewReader(string(out)))
	if err != nil {
		fmt.Println("Failed to make the request:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("Raw Response Body:\n%v\n", string(body))

	if resp.StatusCode != http.StatusCreated {
		fmt.Printf("Request returned an error: %d %s\n", resp.StatusCode, resp.Status)
		return
	}

	fmt.Printf("Response code: %d\n", resp.StatusCode)

}
