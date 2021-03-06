package payment_engine

import (
	"github.com/kelseyhightower/envconfig"
	"stash.tutu.ru/avia-search-common/repository"
	"time"
)

type Specs struct {
	DataURI            string         `envconfig:"FRULE_PAYMENT_ENGINE_REPOSITORY_DATA_URI" required:"true"`
	StatusURI          *string        `envconfig:"FRULE_PAYMENT_ENGINE_REPOSITORY_STATUS_URI"`
	UpdatePeriod       *time.Duration `envconfig:"FRULE_PAYMENT_ENGINE_REPOSITORY_UPDATE_PERIOD"`
	InsecureSkipVerify bool           `envconfig:"FRULE_PAYMENT_ENGINE_REPOSITORY_INSECURE_SKIP_VERIFY" default:"false"`
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
