package manual_exchange_refund

import (
	"github.com/kelseyhightower/envconfig"
	"stash.tutu.ru/avia-search-common/repository"
	"time"
)

type Specs struct {
	DataURI                     string         `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_DATA_URI" required:"true"`
	StatusURI                   *string        `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_STATUS_URI"`
	ComparisonOrderURI          *string        `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_COMPARISON_ORDER_URI"`
	ComparisonOrderUpdatePeriod *time.Duration `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_COMPARISON_ORDER_UPDATE_PERIOD"`
	UpdatePeriod                *time.Duration `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_UPDATE_PERIOD"`
	InsecureSkipVerify          bool           `envconfig:"FRULE_MANUAL_EXCHANGE_REFUND_REPOSITORY_INSECURE_SKIP_VERIFY" default:"false"`
}

func GetConfigFromEnv() (*repository.Config, error) {
	var s Specs
	err := envconfig.Process("", &s)
	if err != nil {
		return nil, err
	}
	return s.toConfig(), nil
}

func (s *Specs) toConfig() *repository.Config {
	return &repository.Config{
		DataURI:            s.DataURI,
		StatusURI:          s.StatusURI,
		UpdatePeriod:       s.UpdatePeriod,
		InsecureSkipVerify: s.InsecureSkipVerify,
	}
}
