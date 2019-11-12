package interline

import (
	"context"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type InterlineRule struct {
	Id                   int     `json:"id"`
	Partner              *string `json:"partner"`
	ConnectionGroup      *string `json:"connection_group"`
	CarrierPlating       *int64  `json:"carrier_plating"`
	PureInterline        *bool   `json:"pure_interline"`
	CarriersForbid       string  `json:"carriers_forbid"`
	CarriersForbidParsed []int64
	CarriersNeed         string `json:"carrier_need"`
	CarriersNeedParsed   []int64
	Carriers             []int64
	Result               bool `json:"result"`
	repo                 *frule_module.Repository
}

func NewInterlineFRule(ctx context.Context, config *repository.Config) (*InterlineRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &InterlineRule{repo: repo}, nil
}

func (rule *InterlineRule) GetResultValue(testRule interface{}) interface{} {
	params := testRule.(InterlineRule)

	if len(params.Carriers) > 0 {
		if len(rule.CarriersNeedParsed) > 0 {
			s := false
			for _, carrierId := range rule.CarriersNeedParsed {
				if frule_module.InSliceInt64(carrierId, params.Carriers) {
					s = true
				}
			}
			if !s {
				return false
			}
		}

		if len(rule.CarriersForbidParsed) > 0 {
			for _, carrierId := range rule.CarriersForbidParsed {
				if frule_module.InSliceInt64(carrierId, params.Carriers) {
					return false
				}
			}
		}
	}

	return rule.Result
}

var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"partner", "connection_group", "carrier_plating", "pure_interline"},
	[]string{"partner", "connection_group", "carrier_plating"},
	[]string{"partner", "connection_group", "pure_interline"},
	[]string{"partner", "connection_group"},
	[]string{"partner", "carrier_plating", "pure_interline"},
	[]string{"partner", "carrier_plating"},
	[]string{"partner", "pure_interline"},
	[]string{"partner"},
}

func (rule *InterlineRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *InterlineRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"partner", "connection_group", "carrier_plating", "pure_interline"}

func (rule *InterlineRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *InterlineRule) GetDefaultValue() interface{} {
	return false
}

func (rule *InterlineRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *InterlineRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *InterlineRule) GetRuleName() string {
	return "Interline"
}
