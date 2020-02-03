//

package accountpoolmetricsiface

import (
	"github.com/Optum/dce/pkg/accountpoolmetrics"
)

// Servicer makes working with the Account Pool Metrics Service struct easier
type Servicer interface {
	// Get returns all account pool metrics
	Get(ID string) (*accountpoolmetrics.AccountPoolMetrics, error)
	// Save writes the record to the dataSvc
	Save(data *accountpoolmetrics.AccountPoolMetrics) error
	// Update account pool metrics
	Update(ID string, data *accountpoolmetrics.AccountPoolMetrics) (*accountpoolmetrics.AccountPoolMetrics, error)
}
