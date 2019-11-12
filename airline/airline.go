package airline

import (
	"context"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type AirlineRule struct {
	Id              int     `json:"id"`
	CarrierId       *int64  `json:"carrier_id"`
	Partner         *string `json:"partner"`
	ConnectionGroup *string `json:"connection_group"`
	Result          bool    `json:"result"`
	repo            *frule_module.Repository
}

func NewAirlineFRule(ctx context.Context, config *repository.Config) (*AirlineRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &AirlineRule{repo: repo}, nil
}

func (rule *AirlineRule) GetResultValue(testRule interface{}) interface{} {
	return rule.Result
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"carrier_id", "partner", "connection_group"},
	[]string{"partner", "connection_group"},
	[]string{"carrier_id", "partner"},
	[]string{"partner"},
}

func (rule *AirlineRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *AirlineRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"carrier_id", "partner", "connection_group"}

func (rule *AirlineRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *AirlineRule) GetDefaultValue() interface{} {
	return false
}

func (rule *AirlineRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *AirlineRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *AirlineRule) GetRuleName() string {
	return "Airline"
}
