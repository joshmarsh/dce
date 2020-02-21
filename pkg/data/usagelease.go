package data

import (
	"fmt"

	"github.com/Optum/dce/pkg/errors"
	"github.com/Optum/dce/pkg/usage"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

const usageLeaseSkPrefix string = "Usage-Lease"

type usageLeaseData struct {
	usage.Lease
	SK         string `json:"-" dynamodbav:"SK" schema:"-"`
	TimeToLive int64  `json:"timeToLive" dynamodbav:"TimeToLive,omitempty" schema:"timeToLive,omitempty"` // ttl attribute
}

// UsageLease - Data Layer Struct
type UsageLease struct {
	DynamoDB       dynamodbiface.DynamoDBAPI
	TableName      string `env:"PRINCIPAL_DB"`
	ConsistentRead bool   `env:"USE_CONSISTENT_READS" envDefault:"false"`
	Limit          int64  `env:"LIMIT" envDefault:"25"`
	TimeToLive     int    `env:"USAGE_TTL" envDefault:"30"`
	BudgetPeriod   string `env:"PRINCIPAL_BUDGET_PERIOD" envDefault:"WEEKLY"`
}

// Write the Usage record in DynamoDB
// This is an upsert operation in which the record will either
// be inserted or updated
// Returns the old record
func (a *UsageLease) Write(usg *usage.Lease) (*usage.Lease, error) {

	var err error
	returnValue := "ALL_OLD"

	usgData := usageLeaseData{
		*usg,
		fmt.Sprintf("%s-%s-%d", usageLeaseSkPrefix, *usg.LeaseID, usg.Date.Unix()),
		getTTL(*usg.Date, a.TimeToLive),
	}

	putMap, _ := dynamodbattribute.Marshal(usgData)
	input := &dynamodb.PutItemInput{
		TableName:    aws.String(a.TableName),
		Item:         putMap.M,
		ReturnValues: aws.String(returnValue),
	}

	old, err := a.DynamoDB.PutItem(input)
	if err != nil {
		return nil, errors.NewInternalServer(
			fmt.Sprintf("update failed for usage with PrincipalID %q and SK %s", *usgData.PrincipalID, usgData.SK),
			err,
		)
	}

	oldUsg := &usage.Lease{}
	err = dynamodbattribute.UnmarshalMap(old.Attributes, oldUsg)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		return nil, err
	}

	diffUsg := usage.Lease{
		PrincipalID:  usg.PrincipalID,
		Date:         usg.Date,
		CostAmount:   usg.CostAmount,
		CostCurrency: usg.CostCurrency,
		LeaseID:      usg.LeaseID,
	}
	if oldUsg.CostAmount != nil {
		diffCost := *diffUsg.CostAmount - *oldUsg.CostAmount
		diffUsg.CostAmount = &diffCost
	}

	err = a.addLeaseUsage(diffUsg)
	if err != nil {
		return nil, err
	}

	err = a.addPrincipalUsage(diffUsg)
	if err != nil {
		return nil, err
	}

	return oldUsg, nil

}

// Add to CostAmount
// Returns new values
func (a *UsageLease) addLeaseUsage(usg usage.Lease) error {

	var err error
	returnValue := "ALL_NEW"
	var expr expression.Expression
	var updateBldr expression.UpdateBuilder

	usgData := usageLeaseData{
		usg,
		fmt.Sprintf("%s-%s-Summary", usageLeaseSkPrefix, *usg.LeaseID),
		getTTL(*usg.Date, a.TimeToLive),
	}

	updateBldr = updateBldr.Add(expression.Name("CostAmount"), expression.Value(usgData.CostAmount))
	updateBldr = updateBldr.Set(expression.Name("CostCurrency"), expression.Value(usgData.CostCurrency))
	updateBldr = updateBldr.Set(expression.Name("Date"), expression.Value(usgData.Date.Unix()))
	updateBldr = updateBldr.Set(expression.Name("TimeToLive"), expression.Value(usgData.TimeToLive))
	expr, err = expression.NewBuilder().WithUpdate(updateBldr).Build()

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PrincipalId": {
				S: usgData.PrincipalID,
			},
			"SK": {
				S: aws.String(usgData.SK),
			},
		},
		TableName:                 aws.String(a.TableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String(returnValue),
	}

	old, err := a.DynamoDB.UpdateItem(input)
	if err != nil {
		return errors.NewInternalServer(
			fmt.Sprintf("update failed for usage with PrincipalID %q and SK %s", *usgData.PrincipalID, usgData.SK),
			err,
		)
	}

	newUsg := &usage.Lease{}
	err = dynamodbattribute.UnmarshalMap(old.Attributes, newUsg)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		return err
	}

	return nil

}

// Add to CostAmount
// Returns new values
func (a *UsageLease) addPrincipalUsage(usg usage.Lease) error {

	var err error
	returnValue := "ALL_NEW"
	var expr expression.Expression
	var updateBldr expression.UpdateBuilder

	periodStart := getBudgetPeriodTime(*usg.Date, a.BudgetPeriod)

	usgPrincipal := usage.Principal{
		PrincipalID:  usg.PrincipalID,
		Date:         &periodStart,
		CostAmount:   usg.CostAmount,
		CostCurrency: usg.CostCurrency,
	}
	usgData := usagePrincipalData{
		usgPrincipal,
		fmt.Sprintf("%s-%d", usagePrincipalSkPrefix, periodStart.Unix()),
		getTTL(*usg.Date, a.TimeToLive),
	}

	updateBldr = updateBldr.Add(expression.Name("CostAmount"), expression.Value(usgData.CostAmount))
	updateBldr = updateBldr.Set(expression.Name("CostCurrency"), expression.Value(usgData.CostCurrency))
	updateBldr = updateBldr.Set(expression.Name("Date"), expression.Value(usgData.Date.Unix()))
	updateBldr = updateBldr.Set(expression.Name("TimeToLive"), expression.Value(usgData.TimeToLive))
	expr, err = expression.NewBuilder().WithUpdate(updateBldr).Build()

	input := &dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"PrincipalId": {
				S: usgData.PrincipalID,
			},
			"SK": {
				S: aws.String(usgData.SK),
			},
		},
		TableName:                 aws.String(a.TableName),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ReturnValues:              aws.String(returnValue),
	}

	old, err := a.DynamoDB.UpdateItem(input)
	if err != nil {
		return errors.NewInternalServer(
			fmt.Sprintf("update failed for usage with PrincipalID %q and SK %s", *usgData.PrincipalID, usgData.SK),
			err,
		)
	}

	newUsg := &usage.Lease{}
	err = dynamodbattribute.UnmarshalMap(old.Attributes, newUsg)
	if err != nil {
		fmt.Printf("Error: %+v", err)
		return err
	}

	return nil

}
