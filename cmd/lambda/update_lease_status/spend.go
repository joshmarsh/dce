package main

import (
	"log"
	"time"

	"github.com/Optum/Redbox/pkg/awsiface"
	"github.com/Optum/Redbox/pkg/budget"
	"github.com/Optum/Redbox/pkg/common"
	"github.com/Optum/Redbox/pkg/db"
	"github.com/Optum/Redbox/pkg/usage"
	"github.com/aws/aws-sdk-go/service/costexplorer"
	"github.com/pkg/errors"
)

type calculateSpendInput struct {
	account    *db.RedboxAccount
	lease      *db.RedboxLease
	tokenSvc   common.TokenService
	budgetSvc  budget.Service
	usageSvc   usage.Service
	awsSession awsiface.AwsSession
}

func calculateSpend(input *calculateSpendInput) (float64, error) {
	adminRoleArn := input.account.AdminRoleArn
	log.Printf("Assuming role %s for budget check", adminRoleArn)
	assumedSession, err := input.tokenSvc.NewSession(input.awsSession, adminRoleArn)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to assume role %s", adminRoleArn)
	}

	// Configure the CostExplorer SDK for the Service
	input.budgetSvc.SetCostExplorer(
		costexplorer.New(assumedSession),
	)

	//Get usage for current date and add it to Usage cache db
	currentTime := time.Now()
	usageStartTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 0, 0, 0, 0, time.UTC)
	usageEndTime := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), 23, 59, 59, 0, time.UTC)

	log.Printf("usageStart: %d and usageEnd :%d", usageStartTime.Unix(), usageEndTime.Unix())
	todayCostAmount, err := input.budgetSvc.CalculateTotalSpend(usageStartTime, usageStartTime.AddDate(0, 0, 1))
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to calculate spend for account %s", input.lease.AccountID)
	}

	log.Printf("usage for today: %f", todayCostAmount)

	// Set Timetolive to one month from StartDate
	usageItem := usage.Usage{
		StartDate:    usageStartTime.Unix(),
		EndDate:      usageEndTime.Unix(),
		PrincipalID:  input.lease.PrincipalID,
		AccountID:    input.account.ID,
		CostAmount:   todayCostAmount,
		CostCurrency: "USD",
		TimeToLive:   usageStartTime.AddDate(0, 1, 0).Unix(),
	}

	input.usageSvc.PutUsage(usageItem)

	// Budget period starts last time the lease was reset.
	// We can look at the `leaseStatusModifiedOn` to know
	// when the lease status changed from `ResetLock` --> `Active`
	budgetStartTime := time.Unix(input.lease.LeaseStatusModifiedOn, 0)
	// budget's `endTime` is set to yesterday
	budgetEndTime := usageEndTime.AddDate(0, 0, -1)

	log.Printf("Retrieving usage for lease %s @ %s for period %s to %s...",
		input.lease.PrincipalID, input.lease.AccountID,
		budgetStartTime.Format("2006-01-02"), budgetEndTime.Format("2006-01-02"),
	)

	// Query Usage cache DB
	usageRecords, err := input.usageSvc.GetUsageByDateRange(budgetStartTime, budgetEndTime)
	if err != nil {
		return 0, errors.Wrapf(err, "Failed to retrieve usage for account %s", input.lease.AccountID)
	}

	// DynDB is eventually consistent. Pull cache DB for SUN-->yesterday, then add the known value for today
	spend := todayCostAmount
	for _, usage := range usageRecords {
		log.Printf("usage records retrieved: %v", usage)
		if usage.PrincipalID == input.lease.PrincipalID && usage.AccountID == input.lease.AccountID {
			spend = spend + usage.CostAmount
		}
	}

	log.Printf("Lease for %s @ %s has spent $%.2f of their $%.2f budget",
		input.lease.PrincipalID, input.lease.AccountID, spend, input.lease.BudgetAmount)

	return spend, nil
}