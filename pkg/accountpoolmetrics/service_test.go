package accountpoolmetrics_test

import (
	"github.com/Optum/dce/pkg/accountpoolmetrics"
	"github.com/Optum/dce/pkg/accountpoolmetrics/mocks"
	"github.com/Optum/dce/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func ptrString(s string) *string {
	ptrS := s
	return &ptrS
}

func TestGet(t *testing.T) {

	type response struct {
		data *accountpoolmetrics.AccountPoolMetrics
		err  error
	}

	tests := []struct {
		name string
		ret  response
		exp  response
	}{
		{
			name: "should get AccountPoolMetrics record",
			ret: response{
				data: &accountpoolmetrics.AccountPoolMetrics{
					ID:     ptrString("123456789012"),
					Ready: intToPtr(0),
					NotReady: intToPtr(1),
					Leased: intToPtr(2),
					Orphaned: intToPtr(3),
				},
				err: nil,
			},
			exp: response{
				data: &accountpoolmetrics.AccountPoolMetrics{
					ID:     ptrString("123456789012"),
					Ready: intToPtr(0),
					NotReady: intToPtr(1),
					Leased: intToPtr(2),
					Orphaned: intToPtr(3),
				},
				err: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocksRwd := &mocks.ReaderWriter{}

			mocksRwd.On("GetSingleton").Return(tt.ret.data, tt.ret.err)

			apmSvc := accountpoolmetrics.NewService(accountpoolmetrics.NewServiceInput{
				DataSvc: mocksRwd,
			})

			getAccountPoolMetrics, err := apmSvc.Get()
			assert.True(t, errors.Is(err, tt.exp.err), "actual error %q doesn't match expected error %q", err, tt.exp.err)

			assert.Equal(t, tt.exp.data, getAccountPoolMetrics)
		})
	}
}

func TestIncrement(t *testing.T) {
	t.Fail()
}

func TestDecrement(t *testing.T) {
	t.Fail()
}

func TestSave(t *testing.T) {
	now := time.Now().Unix()

	type response struct {
		data *accountpoolmetrics.AccountPoolMetrics
		err  error
	}

	tests := []struct {
		name      string
		returnErr error
		metrics   *accountpoolmetrics.AccountPoolMetrics
		exp       response
	}{
		{
			name:    "should save record with timestamps",
			metrics: &accountpoolmetrics.AccountPoolMetrics{
				ID:     ptrString("123456789012"),
				LastModifiedOn: &now,
				CreatedOn: &now,
				Ready: intToPtr(0),
				NotReady: intToPtr(1),
				Leased: intToPtr(2),
				Orphaned: intToPtr(3),
			},
			exp: response{
				data: &accountpoolmetrics.AccountPoolMetrics{
					ID:     ptrString("123456789012"),
					LastModifiedOn: &now,
					CreatedOn: &now,
					Ready: intToPtr(0),
					NotReady: intToPtr(1),
					Leased: intToPtr(2),
					Orphaned: intToPtr(3),
				},
				err: nil,
			},
			returnErr: nil,
		},
		{
			name: "new record should save with new created on",
			metrics: &accountpoolmetrics.AccountPoolMetrics{
				ID:     ptrString("123456789012"),
				Ready: intToPtr(0),
				NotReady: intToPtr(1),
				Leased: intToPtr(2),
				Orphaned: intToPtr(3),
			},
			exp: response{
				data: &accountpoolmetrics.AccountPoolMetrics{
					ID:     ptrString("123456789012"),
					Ready: intToPtr(0),
					NotReady: intToPtr(1),
					Leased: intToPtr(2),
					Orphaned: intToPtr(3),
					LastModifiedOn: &now,
					CreatedOn: &now,
				},
				err: nil,
			},
			returnErr: nil,
		},
		{
			name: "should fail on return err",
			metrics: &accountpoolmetrics.AccountPoolMetrics{
				ID:     ptrString("123456789012"),
				Ready: intToPtr(0),
				NotReady: intToPtr(1),
				Leased: intToPtr(2),
				Orphaned: intToPtr(3),
			},
			exp: response{
				data: &accountpoolmetrics.AccountPoolMetrics{
					ID:     ptrString("123456789012"),
					Ready: intToPtr(0),
					NotReady: intToPtr(1),
					Leased: intToPtr(2),
					Orphaned: intToPtr(3),
					LastModifiedOn: &now,
					CreatedOn: &now,
				},
				err: errors.NewInternalServer("failure", nil),
			},
			returnErr: errors.NewInternalServer("failure", nil),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mocksRwd := &mocks.ReaderWriter{}

			mocksRwd.On("Write", mock.AnythingOfType("*accountpoolmetrics.AccountPoolMetrics"), mock.AnythingOfType("*int64")).Return(tt.returnErr)

			apmSvc := accountpoolmetrics.NewService(accountpoolmetrics.NewServiceInput{
				DataSvc: mocksRwd,
			})

			err := apmSvc.Save(tt.metrics)

			assert.Truef(t, errors.Is(err, tt.exp.err), "actual error %q doesn't match expected error %q", err, tt.exp.err)
			assert.Equal(t, tt.exp.data, tt.metrics)

		})
	}
}

func intToPtr(i int16) *int16 {
	return &i
}