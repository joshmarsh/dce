package tests

import (
	"fmt"
	"github.com/Optum/dce/tests/acceptance/testutil"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestAccountPoolMetrics(t *testing.T) {
	tfOpts := &terraform.Options{
		TerraformDir: "../../modules",
	}
	tfOut := terraform.OutputAll(t, tfOpts)
	apiURL := tfOut["api_url"].(string)
	// TotalAccounts = NotReady + Ready + Leased + Orphaned
	// When account created then TotalAccounts should increment by one
	// When account deleted then TotalAccounts should decrement by one
	// When account goes from !NotReady to NotReady then NotReady should increment by one and other status should decrement by one AND TotalAccounts should not change
	// When account goes from !Ready to Ready then NotReady should increment by one and other status should decrement by one AND TotalAccounts should not change
	// When account goes from !Leased to Leased then NotReady should increment by one and other status should decrement by one AND TotalAccounts should not change
	// When account goes from !Orphaned to Orphaned then NotReady should increment by one and other status should decrement by one AND TotalAccounts should not change

	// TotalAccounts = NotReady + Ready + Leased + Orphaned
	t.Run("When account created then TotalAccounts should increment by one", func(t *testing.T) {
		awsSession, err := session.NewSession()
		require.Nil(t, err)


		// Check TotalAccounts

		// Create account
		role := testutil.CreateAdminAPIInvokeRole(t, awsSession)
		createAccountRes := testutil.InvokeApiWithRetry(t, &testutil.InvokeApiWithRetryInput{
			Method: "POST",
			URL:    apiURL + "/accounts",
			JSON: createAccountRequest{
				ID:           role.AccountID,
				AdminRoleArn: role.AdminRoleArn,
			},
			MaxAttempts: 15,
			F: func(r *testutil.R, apiResp *testutil.ApiResponse) {
				assert.Equal(r, 201, apiResp.StatusCode)
			},
		})

		// Check the response
		postResJSON := testutil.ParseResponseJSON(t, createAccountRes)
		fmt.Println("@@@@@", postResJSON)
		require.Equal(t, role.AccountID, postResJSON["id"])
		require.Equal(t, "NotReady", postResJSON["accountStatus"])
		require.Equal(t, role.AdminRoleArn, postResJSON["adminRoleArn"])
	})
}