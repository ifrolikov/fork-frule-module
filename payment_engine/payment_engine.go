package payment_engine

import (
	"context"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"strings"
)

type EngineConfig struct {
	Engine        string
	ConfigType    *string
	ConfigSubtype *string
}

type PaymentEngineRule struct {
	Id              int     `json:"id"`
	Partner         *string `json:"partner"`
	ConnectionGroup *string `json:"connection_group"`
	CarrierId       *int64  `json:"carrier_id"`
	PaymentMethod   *string `json:"payment_method"`
	RealGds         *string `json:"real_gds"`
	Engine          string  `json:"engine"`
	ConfigType      *string `json:"config_type"`
	Subtype         *string `json:"subtype"`
	repo            *frule_module.Repository
}

func NewPaymentEngineFRule(ctx context.Context, config *repository.Config) (*PaymentEngineRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &PaymentEngineRule{repo: repo}, nil
}

func (rule *PaymentEngineRule) GetResultValue(interface{}) interface{} {
	var result []EngineConfig
	for idx, engine := range strings.Split(rule.Engine, "|") {
		var engineConfig EngineConfig
		if rule.ConfigType != nil {
			configs := strings.Split(*rule.ConfigType, "|")
			engineConfig = EngineConfig{
				Engine:        engine,
				ConfigType:    &configs[idx],
				ConfigSubtype: rule.Subtype,
			}
		} else {
			engineConfig = EngineConfig{
				Engine:        engine,
				ConfigType:    nil,
				ConfigSubtype: rule.Subtype,
			}

		}
		result = append(result, engineConfig)
	}
	return result
}

func (rule *PaymentEngineRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{"partner", "connection_group", "real_gds", "carrier_id", "payment_method"},
		[]string{"partner", "connection_group", "real_gds", "payment_method"},
		[]string{"partner", "real_gds", "carrier_id", "payment_method"},
		[]string{"partner", "real_gds", "payment_method"},
		[]string{"partner", "connection_group", "carrier_id", "payment_method"},
		[]string{"partner", "connection_group", "payment_method"},
		[]string{"partner", "carrier_id", "payment_method"},
		[]string{"partner", "payment_method"},
		[]string{"payment_method"},
	}
}

func (rule *PaymentEngineRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{}
}

func (rule *PaymentEngineRule) GetStrategyKeys() []string {
	return []string{"partner", "connection_group", "real_gds", "carrier_id", "payment_method"}
}

func (rule *PaymentEngineRule) GetDefaultValue() interface{} {
	return []EngineConfig{}
}

func (rule *PaymentEngineRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *PaymentEngineRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *PaymentEngineRule) GetRuleName() string {
	return "PaymentEngine"
}
