package testutil

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/client"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/service/iam"
	aws2 "github.com/gruntwork-io/terratest/modules/aws"
	"github.com/stretchr/testify/require"
	"time"

	"testing"
)

type AssumableRole struct {
	AccountID    string
	RoleName     string
	AdminRoleArn string
	Policies 	[]string
}

var chainCredentials = credentials.NewChainCredentials([]credentials.Provider{
	&credentials.EnvProvider{},
	&credentials.SharedCredentialsProvider{Filename: "", Profile: ""},
})

// CreateAssumableRole creates an assumable role in the account referred to by awsSession. Policies may be attached
// via the policies parameter.
func CreateAssumableRole(t *testing.T, awsSession client.ConfigProvider, adminRoleName string, policies []string) *AssumableRole {
	currentAccountID := aws2.GetAccountId(t)

	// Create an Admin Role that can be assumed
	// within this account
	iamSvc := iam.New(awsSession)
	assumeRolePolicy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
					"AWS": "arn:aws:iam::%s:root"
					},
					"Action": "sts:AssumeRole",
					"Condition": {}
				}
			]
		}`, currentAccountID)
	roleRes, err := iamSvc.CreateRole(&iam.CreateRoleInput{
		AssumeRolePolicyDocument: aws.String(assumeRolePolicy),
		Path:                     aws.String("/"),
		RoleName:                 aws.String(adminRoleName),
	})
	require.Nil(t, err)

	adminRoleArn := *roleRes.Role.Arn

	for _, p := range policies {
		_, err = iamSvc.AttachRolePolicy(&iam.AttachRolePolicyInput{
			RoleName:  aws.String(adminRoleName),
			PolicyArn: aws.String(p),
		})
		require.Nil(t, err)
	}
	// IAM Role takes a while to propagate....
	//time.Sleep(10 * time.Second)

	return &AssumableRole{
		AdminRoleArn: adminRoleArn,
		RoleName:     adminRoleName,
		AccountID:    currentAccountID,
		Policies: policies,
	}
}

func CreateAdminAPIInvokeRole(t *testing.T, awsSession client.ConfigProvider) *AssumableRole {
	adminRoleName := "dce-api-test-admin-role-" + fmt.Sprintf("%v", time.Now().Unix())
	policies := []string{
		"arn:aws:iam::aws:policy/IAMFullAccess",
		"arn:aws:iam::aws:policy/AmazonAPIGatewayInvokeFullAccess",
	}
	return CreateAssumableRole(t, awsSession, adminRoleName, policies)
}















