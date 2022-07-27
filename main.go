package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
)

type event struct {
	URL  string `json:"url"`
	Data string `json:"data"`
}

func newClient() lambda.Lambda {
	var s *session.Session

	//Used NewSharedCredentials, which will take the required AWS keyID and secret key from ~/.aws/credentials file
	s = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String("us-west-2"),
		Credentials: credentials.NewSharedCredentials("", ""),
	}))
	var lambdaClient *lambda.Lambda = lambda.New(s)
	return *lambdaClient
}

func send(evt event, client interface{}) int {
	data, err := json.Marshal(evt)
	if err != nil {
		fmt.Println("Error while marshalling:", err)
		return 500
	}
	fmt.Println(string(data))

	lambdaClient, _ := client.(lambda.Lambda)
	res, err := lambdaClient.Invoke(&lambda.InvokeInput{
		FunctionName: aws.String("testFunction"),
		Payload:      data,
	})
	if err != nil {
		fmt.Println("Error calling testFunction:", err)
		return 501
	}
	fmt.Println("res.Payload:", string(res.Payload))

	var statusCode int
	err = json.Unmarshal(res.Payload, &statusCode)
	return statusCode
}

func main() {
	evt := event{"https://webhook.site/59ca3a13-3cca-4012-a716-94d1ac717be6", "Hi, this is lambda poc!"}
	//data, _ := json.Marshal(evt)
	//fmt.Println(string(data))

	client := newClient()
	res := send(evt, client)
	fmt.Println("response:", res)
}
