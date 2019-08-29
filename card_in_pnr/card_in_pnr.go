package card_in_pnr

import (
	"context"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
)

type CardInPnrRule struct {
	Id              int     `json:"id"`
	Partner         *string `json:"partner"`
	ConnectionGroup *string `json:"connection_group"`
	CarrierId       *int64  `json:"carrier_id"`
	Result          bool    `json:"result"`
	repo            *frule_module.Repository
}

func NewCardInPnrRuleFRule(ctx context.Context, config *repository.Config) (*CardInPnrRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &CardInPnrRule{repo: repo}, nil
}

func (rule *CardInPnrRule) GetResultValue(interface{}) interface{} {
	return rule.Result
}

func (rule *CardInPnrRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{"partner", "carrier_id", "connection_group"},
		[]string{"partner", "connection_group"},
		[]string{"partner", "carrier_id"},
		[]string{"partner"},
	}
}

func (rule *CardInPnrRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{}
}

func (rule *CardInPnrRule) GetStrategyKeys() []string {
	return []string{"partner", "carrier_id", "connection_group"}
}

func (rule *CardInPnrRule) GetDefaultValue() interface{} {
	return false
}

func (rule *CardInPnrRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *CardInPnrRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *CardInPnrRule) GetRuleName() string {
	return "CardInPnr"
}