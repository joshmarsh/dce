package data

import (
	"errors"
	"fmt"
	"github.com/Optum/dce/pkg/accountpoolmetrics"
	errWrapper "github.com/Optum/dce/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// AccountPoolMetrics - Data Layer Struct
type AccountPoolMetrics struct {
	DynamoDB       dynamodbiface.DynamoDBAPI
	TableName      string `env:"ACCOUNT_POOL_METRICS_DB"`
	ConsistentRead bool   `env:"USE_CONSISTENT_READS" envDefault:"false"`
	Limit          int64  `env:"LIMIT" envDefault:"25"`
}

// GetSingleton retrieves the only record in the database. An error is thrown if multiple records exist.
func (a *AccountPoolMetrics) GetSingleton() (*accountpoolmetrics.AccountPoolMetrics, error) {
	res, err := a.DynamoDB.Scan(
		&dynamodb.ScanInput{
			// Query in Lease Table
			TableName: aws.String(a.TableName),
			ConsistentRead: aws.Bool(a.ConsistentRead),
		},
	)
	if err != nil {
		return nil, errWrapper.NewInternalServer(
			"getsingleton failed for accountpoolmetrics",
			err,
		)
	}

	apmArr := []*accountpoolmetrics.AccountPoolMetrics{}
	for i, item := range res.Items {
		apmArr = append(apmArr, &accountpoolmetrics.AccountPoolMetrics{})
		err = dynamodbattribute.UnmarshalMap(item, apmArr[i])
		if err != nil {
			return nil, errWrapper.NewInternalServer(
				fmt.Sprintf("failure unmarshaling accountpoolmetrics"),
				err,
			)
		}
	}

	if len(apmArr) != 1 {
		err := errors.New(fmt.Sprintf("getsingleton failed for accountpoolmetrics, expected 1 record but found  %q", len(res.Items)))
		return nil, errWrapper.NewInternalServer(
			"failure retrieving singleton record",
			err,
		)
	}

	return apmArr[0], nil
}

func (a *AccountPoolMetrics) Write(i *accountpoolmetrics.AccountPoolMetrics, lastModifiedOn *int64) error {
	return nil
}

func (a *AccountPoolMetrics) Increment(name accountpoolmetrics.MetricName) error {
	return nil
}

func (a *AccountPoolMetrics) Decrement(name accountpoolmetrics.MetricName) error {
	return nil
}
