package main

import (
	"context"
	"log"

	"github.com/Optum/dce/pkg/common"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// Start the Lambda Handler
func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, event events.DynamoDBEvent) error {

	awsSession := session.Must(session.NewSession())
	accountsTopic := common.RequireEnv("ACCOUNT_TOPIC_ARN")
	snsSvc := &common.SNS{Client: sns.New(awsSession)}

	records, err := common.PrepareSNSMessageJSON(event.Records)
	_, err = snsSvc.PublishMessage(&accountsTopic, stringPtr(records), true)
	if err != nil {
		log.Fatalf(err.Error())
	}

	//// We get a stream of DynDB records, representing changes to the table
	//for _, record := range event.Records {
	//	fmt.Println(record)
	//
	//}

	return nil
}


func stringPtr(str string) *string {
	return &str
}