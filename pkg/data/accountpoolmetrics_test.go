package data

import (
	"errors"
	"fmt"
	"github.com/Optum/dce/pkg/accountpoolmetrics"
	awsmocks "github.com/Optum/dce/pkg/awsiface/mocks"

	errWrapper "github.com/Optum/dce/pkg/errors"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"strconv"
	"testing"
)

func TestGetSingleton(t *testing.T) {

	err := errors.New(fmt.Sprintf("getsingleton failed for accountpoolmetrics, expected 1 record but found  %q", 0))
	zeroRecordsError :=  errWrapper.NewInternalServer(
		"failure retrieving singleton record",
		err,
	)

	err = errors.New(fmt.Sprintf("getsingleton failed for accountpoolmetrics, expected 1 record but found  %q", 2))
	twoRecordsError :=  errWrapper.NewInternalServer(
		"failure retrieving singleton record",
		err,
	)

	dynamoErr := errors.New("dynamo error")
	dynamoErrorWrapped := errWrapper.NewInternalServer(
		"getsingleton failed for accountpoolmetrics",
		dynamoErr,
	)

	tests := []struct {
		name           string
		dynamoErr      error
		dynamoOutput   *dynamodb.ScanOutput
		expectedErr    error
		expectedObject *accountpoolmetrics.AccountPoolMetrics
	}{
		{
			name:      "should return a single metric object when one exists",
			expectedObject: &accountpoolmetrics.AccountPoolMetrics{
					ID:             ptrString("123456789012"),
					LastModifiedOn: ptrInt64(1573592058),
					CreatedOn:      ptrInt64(1573592058),
					Ready:          ptrInt16(1),
					NotReady:       ptrInt16(2),
					Leased:         ptrInt16(3),
					Orphaned:       ptrInt16(4),
				},
			dynamoErr: nil,
			dynamoOutput: &dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"Id": {
							S: aws.String("123456789012"),
						},
						"LastModifiedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"CreatedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"Ready": {
							N: aws.String("1"),
						},
						"NotReady": {
							N: aws.String("2"),
						},
						"Leased": {
							N: aws.String("3"),
						},
						"Orphaned": {
							N: aws.String("4"),
						},
					},
				},
			},
			expectedErr: nil,
		},
		{
			name:           "should error when no objects exists",
			expectedObject: nil,
			dynamoErr:      nil,
			dynamoOutput: &dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{},
			},
			expectedErr: zeroRecordsError,
		},
		{
			name:           "should error when two objects exists",
			expectedObject: nil,
			dynamoErr:      nil,
			dynamoOutput: &dynamodb.ScanOutput{
				Items: []map[string]*dynamodb.AttributeValue{
					{
						"Id": {
							S: aws.String("123456789012"),
						},
						"LastModifiedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"CreatedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"Ready": {
							N: aws.String("1"),
						},
						"NotReady": {
							N: aws.String("2"),
						},
						"Leased": {
							N: aws.String("3"),
						},
						"Orphaned": {
							N: aws.String("4"),
						},
					},					{
						"Id": {
							S: aws.String("123456789012"),
						},
						"LastModifiedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"CreatedOn": {
							N: aws.String(strconv.Itoa(1573592058)),
						},
						"Ready": {
							N: aws.String("1"),
						},
						"NotReady": {
							N: aws.String("2"),
						},
						"Leased": {
							N: aws.String("3"),
						},
						"Orphaned": {
							N: aws.String("4"),
						},
					},
				},
			},
			expectedErr: twoRecordsError,
		},
		{
			name:            "should return nil when dynamodb err",
			expectedObject:  nil,
			dynamoErr:       dynamoErr,
			dynamoOutput:    &dynamodb.ScanOutput{
				Items:       []map[string]*dynamodb.AttributeValue{},
			},
			expectedErr: dynamoErrorWrapped,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange
			mockDynamo := awsmocks.DynamoDBAPI{}

			mockDynamo.On("Scan", mock.MatchedBy(func(input *dynamodb.ScanInput) bool {
				return *input.TableName == "AccountPoolMetrics"
			})).Return(
				tt.dynamoOutput, tt.dynamoErr,
			)
			apmData := &AccountPoolMetrics{
				DynamoDB:  &mockDynamo,
				TableName: "AccountPoolMetrics",
			}

			// Act
			result, err := apmData.GetSingleton()

			// Assert
			assert.True(t, errWrapper.Is(err, tt.expectedErr))
			assert.Equal(t, tt.expectedObject, result)
		})
	}
}

func TestWrite(t *testing.T) {
	t.Fail()
}

func TestIncrement(t *testing.T) {
	t.Fail()
}

func TestDecrement(t *testing.T) {
	t.Fail()
}