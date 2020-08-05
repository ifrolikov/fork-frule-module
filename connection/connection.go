package connection

import (
	"context"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type ConnectionRule struct {
	Id              int     `json:"id"`
	Partner         *string `json:"partner"`
	ConnectionGroup *string `json:"connection_group"`
	CarrierId       *int64  `json:"carrier_id"`
	Operation       *string `json:"operation"`
	FlightDate      *string `json:"flight_date"`
	PaymentEngine   *string `json:"payment_engine"`
	Result          string  `json:"result"`
	repo            *frule_module.Repository
}

func NewConnectionFRule(ctx context.Context, config *repository.Config) (*ConnectionRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &ConnectionRule{repo: repo}, nil
}

func (rule *ConnectionRule) GetResultValue(interface{}) interface{} {
	return rule.Result
}

func (rule *ConnectionRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"partner", "connection_group", "operation", "carrier_id", "flight_date", "payment_engine"},
	[]string{"partner", "connection_group", "operation", "carrier_id", "flight_date"},
	[]string{"partner", "connection_group", "operation", "carrier_id", "payment_engine"},
	[]string{"partner", "connection_group", "operation", "carrier_id"},

	[]string{"partner", "connection_group", "carrier_id", "flight_date", "payment_engine"},
	[]string{"partner", "connection_group", "carrier_id", "flight_date"},
	[]string{"partner", "connection_group", "carrier_id", "payment_engine"},

	[]string{"partner", "connection_group", "operation", "flight_date", "payment_engine"},
	[]string{"partner", "connection_group", "operation", "flight_date"},
	[]string{"partner", "connection_group", "operation", "payment_engine"},

	[]string{"partner", "connection_group", "carrier_id"},
	[]string{"partner", "connection_group", "operation"},
	[]string{"partner", "connection_group", "payment_engine"},

	[]string{"partner", "connection_group"},
}

func (rule *ConnectionRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *ConnectionRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"partner", "connection_group", "operation", "carrier_id", "flight_date", "payment_engine"}

func (rule *ConnectionRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *ConnectionRule) GetDefaultValue() interface{} {
	return ""
}

func (rule *ConnectionRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *ConnectionRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *ConnectionRule) GetRuleName() string {
	return "Connection"
}
