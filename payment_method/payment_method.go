package payment_method

import (
	"context"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"strings"
	"time"
)

type PaymentMethodRule struct {
	Id                      int     `json:"id"`
	Partner                 *string `json:"partner"`
	ConnectionGroup         *string `json:"connection_group"`
	CarrierId               *int64  `json:"carrier_id"`
	DaysTillDeparture       *string `json:"days_till_departure"`
	DaysTillDepartureParsed []frule_module.CronStructString
	Autoticketing           *bool  `json:"autoticketing"`
	Result                  string `json:"result"`
	ResultParsed            []frule_module.CronStructString
	TestDaysTillDeparture   int
	repo                    *frule_module.Repository
}

func NewPaymentMethodFRule(ctx context.Context, config *repository.Config) (*PaymentMethodRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &PaymentMethodRule{repo: repo}, nil
}

func (rule *PaymentMethodRule) GetResultValue(interface{}) interface{} {
	for _, daysTillDeparture := range rule.DaysTillDepartureParsed {
		if daysTillDeparture.Value == "" {
			return rule.parseResult()
		} else if frule_module.CronSpec(&daysTillDeparture.Spec, time.Now()) {
			days, err := strconv.Atoi(daysTillDeparture.Value)
			if err != nil {
				log.Logger.Error().Stack().Err(err).Msg("Parsing days")
			}
			if rule.TestDaysTillDeparture <= days {
				return rule.GetDefaultValue()
			} else {
				return rule.parseResult()
			}
		}
	}
	return rule.GetDefaultValue()
}

func (rule *PaymentMethodRule) parseResult() []string {
	var resultSlice []string
	for _, result := range rule.ResultParsed {
		if frule_module.CronSpec(&result.Spec, time.Now()) {
			for _, paymentMethodName := range strings.Split(result.Value, ",") {
				resultSlice = append(resultSlice, strings.TrimSpace(paymentMethodName))
			}
		}
	}
	return resultSlice
}

func (rule *PaymentMethodRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return frule_module.ComparisonOrder{
		[]string{"partner", "connection_group", "carrier_id", "autoticketing"},
		[]string{"partner", "connection_group", "carrier_id"},
		[]string{"partner", "connection_group", "autoticketing"},
		[]string{"partner", "connection_group"},
		[]string{"partner", "carrier_id", "autoticketing"},
		[]string{"partner", "carrier_id"},
		[]string{"partner", "autoticketing"},
		[]string{"partner"},
	}
}

func (rule *PaymentMethodRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return frule_module.ComparisonOperators{}
}

func (rule *PaymentMethodRule) GetStrategyKeys() []string {
	return []string{"partner", "connection_group", "carrier_id", "autoticketing"}
}

func (rule *PaymentMethodRule) GetDefaultValue() interface{} {
	return []string{}
}

func (rule *PaymentMethodRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *PaymentMethodRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *PaymentMethodRule) GetRuleName() string {
	return "PaymentMethod"
}
