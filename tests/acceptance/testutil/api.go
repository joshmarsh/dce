package testutil

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go/aws/credentials"
	sigv4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type InvokeApiWithRetryInput struct {
	Method      string
	URL         string
	Creds       *credentials.Credentials
	Region      string
	JSON        interface{}
	MaxAttempts int
	// Callback function to assert API responses.
	// apiRequest() will continue to retry until this
	// function passes assertions.
	//
	// eg.
	//		F: func(r *testutil.R, apiResp *ApiResponse) {
	//			assert.Equal(r, 200, apiResp.StatusCode)
	//		},
	// or:
	//		F: statusCodeAssertion(200)
	//
	// By default, this will check that the API returns a 2XX response
	F func(r *R, apiResp *ApiResponse)
}

type ApiResponse struct {
	http.Response
	JSON interface{}
}

func InvokeApiWithRetry(t *testing.T, input *InvokeApiWithRetryInput) *ApiResponse {
	// Set defaults
	if input.Creds == nil {
		input.Creds = chainCredentials
	}
	if input.Region == "" {
		input.Region = "us-east-1"
	}
	if input.MaxAttempts == 0 {
		input.MaxAttempts = 30
	}

	// Create API request
	req, err := http.NewRequest(input.Method, input.URL, nil)
	assert.Nil(t, err)

	// Sign our API request, using sigv4
	// See https://docs.aws.amazon.com/general/latest/gr/sigv4_signing.html
	signer := sigv4.NewSigner(input.Creds)
	now := time.Now().Add(time.Duration(30) * time.Second)
	var signedHeaders http.Header
	var apiResp *ApiResponse
	Retry(t, input.MaxAttempts, 2*time.Second, func(r *R) {
		// If there's a JSON provided, add it when signing
		// Body does not matter if added before the signing, it will be overwritten
		if input.JSON != nil {
			payload, err := json.Marshal(input.JSON)
			assert.Nil(t, err)
			req.Header.Set("Content-Type", "application/JSON")
			signedHeaders, err = signer.Sign(req, bytes.NewReader(payload),
				"execute-api", input.Region, now)
			require.Nil(t, err)
		} else {
			signedHeaders, err = signer.Sign(req, nil, "execute-api",
				input.Region, now)
		}
		assert.NoError(r, err)
		assert.NotNil(r, signedHeaders)

		// Send the API requests
		// resp, err := http.DefaultClient.Do(req)
		httpClient := http.Client{
			Timeout: 60 * time.Second,
		}
		resp, err := httpClient.Do(req)
		assert.NoError(r, err)

		// Parse the JSON response
		apiResp = &ApiResponse{
			Response: *resp,
		}
		defer resp.Body.Close()
		var data interface{}

		body, err := ioutil.ReadAll(resp.Body)
		assert.NoError(r, err)

		err = json.Unmarshal([]byte(body), &data)
		if err == nil {
			apiResp.JSON = data
		}

		if input.F != nil {
			input.F(r, apiResp)
		}
	})
	return apiResp
}

func ParseResponseJSON(t *testing.T, resp *ApiResponse) map[string]interface{} {
	require.NotNil(t, resp.JSON)
	return resp.JSON.(map[string]interface{})
}