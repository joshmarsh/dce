package accountpoolmetrics

// AccountPoolMetrics - Handles importing and exporting AccountPoolMetrics
type AccountPoolMetrics struct {
	ID     				   *string                `json:"id,omitempty" dynamodbav:"Id" schema:"id,omitempty"`                                                              // AccountPoolMetrics ID
	LastModifiedOn         *int64                 `json:"lastModifiedOn,omitempty" dynamodbav:"LastModifiedOn" schema:"lastModifiedOn,omitempty"`                          // Last Modified Epoch Timestamp
	CreatedOn              *int64                 `json:"createdOn,omitempty"  dynamodbav:"CreatedOn,omitempty" schema:"createdOn,omitempty"`                              // Account CreatedOn
	Ready                  *int16                 `json:"ready,omitempty" dynamodbav:"Ready,omitempty"`                                                                    // Number of accounts in "Ready" status
	NotReady               *int16                 `json:"notReady,omitempty" dynamodbav:"NotReady,omitempty"`                                                              // Number of accounts in "NotReady" status
	Leased                 *int16                 `json:"leased,omitempty" dynamodbav:"Leased,omitempty"`                                                                  // Number of accounts in "Leased" status
	Orphaned               *int16                 `json:"orphaned,omitempty" dynamodbav:"Orphaned,omitempty"`                                                              // Number of accounts in "Orphaned" status
}

// MetricName is a name of an AccountPoolMetrics metric
type MetricName string

const (
	// Ready is the number of accounts in "Ready" status
	Ready MetricName = "Ready"
	// NotReady is the number of accounts in "NotReady" status
	NotReady MetricName = "NotReady"
	// Leased is the number of accounts in "Leased" status
	Leased MetricName = "Leased"
	// Orphaned is the number of accounts in "Orphaned" status
	Orphaned MetricName = "Orphaned"
)