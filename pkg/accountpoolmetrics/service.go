package accountpoolmetrics

import "time"

// Writer put an item into the data store
type Writer interface {
	Write(i *AccountPoolMetrics, lastModifiedOn *int64) error
}

// MultipleReader reads multiple accounts from the data store
type SingletonReader interface {
	GetSingleton() (*AccountPoolMetrics, error)
}

// Reader data Layer
type Reader interface {
	SingletonReader
}

// ReaderWriter includes Reader and Writer interfaces
type ReaderWriter interface {
	Reader
	Writer
}

type Service struct{
	dataSvc ReaderWriter
}

func (s *Service) Get() (*AccountPoolMetrics, error) {
	new, err := s.dataSvc.GetSingleton()
	if err != nil {
		return nil, err
	}

	return new, err
}

func (s *Service) Save(data *AccountPoolMetrics) error {
	var lastModifiedOn *int64
	now := time.Now().Unix()
	if data.LastModifiedOn == nil {
		lastModifiedOn = nil
		data.CreatedOn = &now
		data.LastModifiedOn = &now
	} else {
		lastModifiedOn = data.LastModifiedOn
		data.LastModifiedOn = &now
	}

	//err := data.Validate()
	//if err != nil {
	//	return err
	//}
	err := s.dataSvc.Write(data, lastModifiedOn)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) Increment(name MetricName) (*AccountPoolMetrics, error) {
	panic("implement me")
}

func (s *Service) Decrement(name MetricName) (*AccountPoolMetrics, error) {
	panic("implement me")
}

// NewServiceInput Input for creating a new Service
type NewServiceInput struct {
	DataSvc ReaderWriter
}

// NewService creates a new instance of the Service
func NewService(input NewServiceInput) *Service {
	return &Service{
		dataSvc:    input.DataSvc,
	}
}
