package codeshare

import (
	"context"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type CodeshareRule struct {
	Id               int     `json:"id"`
	Partner          *string `json:"partner"`
	ConnectionGroup  *string `json:"connection_group"`
	CarrierOperating *int64  `json:"carrier_operating"`
	CarrierMarketing *int64  `json:"carrier_marketing"`
	ServiceClass     *string `json:"service_class"`
	Result           bool    `json:"result"`
	repo             *frule_module.Repository
}

func NewCodeshareFRule(ctx context.Context, config *repository.Config) (*CodeshareRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &CodeshareRule{repo: repo}, nil
}

func (rule *CodeshareRule) GetResultValue(testRule interface{}) interface{} {
	return rule.Result
}

func (rule *CodeshareRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}

func (rule *CodeshareRule) GetCreateRuleHashForIndexedFieldsFunction() *frule_module.CreateRuleHashForIndexedFieldsFunction {
	return nil
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"partner", "connection_group", "carrier_operating", "carrier_marketing", "service_class"},
	[]string{"partner", "connection_group", "carrier_operating", "carrier_marketing"},
	[]string{"partner", "connection_group", "service_class"},
	[]string{"partner", "connection_group"},
	[]string{"partner", "carrier_operating", "carrier_marketing", "service_class"},
	[]string{"partner", "carrier_operating", "carrier_marketing"},
	[]string{"partner", "service_class"},
	[]string{"partner"},
}

func (rule *CodeshareRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *CodeshareRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"partner", "connection_group", "carrier_operating", "carrier_marketing", "service_class"}

func (rule *CodeshareRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *CodeshareRule) GetDefaultValue() interface{} {
	return false
}

func (rule *CodeshareRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *CodeshareRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *CodeshareRule) GetRuleName() string {
	return "Codeshare"
}
