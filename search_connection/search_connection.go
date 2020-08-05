package search_connection

import (
	"context"
	"encoding/json"
	"regexp"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"stash.tutu.ru/golang/log"
	"strconv"
	"time"
)

type SearchConnectionRule struct {
	Id                     int     `json:"id"`
	Partner                *string `json:"partner"`
	ConnectionGroup        *string `json:"connection_group"`
	DepartureDate          time.Time
	MinDepartureDate       *string `json:"min_departure_date"`
	MinDepartureDateParsed []frule_module.CronStructString
	MaxDepartureDate       *string `json:"max_departure_date"`
	MaxDepartureDateParsed []frule_module.CronStructString
	repo                   *frule_module.Repository
}

var specSearchConnectionRegexp = regexp.MustCompile(`\+(\d+)([dwmy])`)

func NewSearchConnectionFRule(ctx context.Context, config *repository.Config) (*SearchConnectionRule, error) {
	repo, err := frule_module.NewFRuleRepository(
		ctx,
		&fruleStorageContainer{},
		&importer{repository.BasicImporter{Config: config}})
	if err != nil {
		return nil, err
	}
	return &SearchConnectionRule{repo: repo}, nil
}

func (rule *SearchConnectionRule) GetResultValue(testRule interface{}) interface{} {
	nowLocal := time.Now()
	nowUtc := time.Now().UTC()

	if minDepartureDateBorder := rule.getSpecInterval(rule.MinDepartureDateParsed, nowLocal); minDepartureDateBorder != "" {
		if minDepartureDate := rule.getDateToCompare(minDepartureDateBorder, nowUtc); minDepartureDate != nil {
			if !testRule.(SearchConnectionRule).DepartureDate.After(*minDepartureDate) {
				return false
			}
		}
	}

	if maxDepartureDateBorder := rule.getSpecInterval(rule.MaxDepartureDateParsed, nowLocal); maxDepartureDateBorder != "" {
		if maxDepartureDate := rule.getDateToCompare(maxDepartureDateBorder, nowUtc); maxDepartureDate != nil {
			if !testRule.(SearchConnectionRule).DepartureDate.Before(*maxDepartureDate) {
				return false
			}
		}
	}

	return true
}

func (rule *SearchConnectionRule) GetCompareDynamicFieldsFunction() *frule_module.CompareDynamicFieldsFunction {
	return nil
}


var comparisonOrder = frule_module.ComparisonOrder{
	[]string{"partner", "connection_group"},
}

func (rule *SearchConnectionRule) GetComparisonOrder() frule_module.ComparisonOrder {
	return comparisonOrder
}

var comparisonOperators = frule_module.ComparisonOperators{}

func (rule *SearchConnectionRule) GetComparisonOperators() frule_module.ComparisonOperators {
	return comparisonOperators
}

var strategyKeys = []string{"partner", "connection_group"}

func (rule *SearchConnectionRule) GetStrategyKeys() []string {
	return strategyKeys
}

func (rule *SearchConnectionRule) GetDefaultValue() interface{} {
	return false
}

/*
func (rule *SearchConnectionRule) GetDataStorage() (map[int][]frule_module.FRuler, error) {
	result := make(map[int][]frule_module.FRuler)
	repo := createRepository(rule.config)
	for rank, ruleList := range repo.GetStorage() {
		for _, ruleItem := range ruleList {
			var err error
			ruleItem.MinDepartureDateParsed, err = rule.parseCronSpecField(*ruleItem.MinDepartureDate)
			if err != nil {
				return nil, err
			}
			ruleItem.MaxDepartureDateParsed, err = rule.parseCronSpecField(*ruleItem.MaxDepartureDate)
			if err != nil {
				return nil, err
			}
			result[rank] = append(result[rank], ruleItem)
		}
	}
	return result, nil
}*/

func (rule *SearchConnectionRule) GetDataStorage() *frule_module.RankedFRuleStorage {
	return rule.repo.GetRankedFRuleStorage()
}

func (rule *SearchConnectionRule) GetNotificationChannel() chan repository.Notification {
	return rule.repo.NotificationChannel
}

func (rule *SearchConnectionRule) GetRuleName() string {
	return "SearchConnection"
}

func (rule *SearchConnectionRule) parseCronSpecField(value string) ([]frule_module.CronStructString, error) {
	var resultParsed []frule_module.CronStructString
	err := json.Unmarshal([]byte(value), &resultParsed)

	if err != nil {
		log.Logger.Error().Stack().Err(err).Msg("Unmarshal")
	}

	return resultParsed, nil
}

func (rule *SearchConnectionRule) getSpecInterval(specs []frule_module.CronStructString, t time.Time) string {
	for i := range specs {
		if frule_module.CronSpec(&specs[i].Spec, t) {
			return specs[i].Value
		}
	}
	return ""
}

func (rule *SearchConnectionRule) getDateToCompare(intervalString string, t time.Time) *time.Time {
	intervalParsed := specSearchConnectionRegexp.FindStringSubmatch(intervalString)

	if len(intervalParsed) >= 3 {
		intervalInt, err := strconv.Atoi(intervalParsed[1])
		if err != nil {
			return nil
		}

		var years, months, days int

		switch intervalParsed[2] {
		case "d":
			days = intervalInt
		case "w":
			days = intervalInt * 7
		case "m":
			months = intervalInt
		case "y":
			years = intervalInt
		default:
			return nil
		}

		minDepartureDate := t.AddDate(years, months, days)

		return &minDepartureDate
	}

	return nil
}
