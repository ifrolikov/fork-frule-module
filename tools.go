package frule_module

import (
	"github.com/robfig/cron"
	"reflect"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/golang/resources/db"
	"strconv"
	"strings"
	"time"
)

func inSlice(item string, list []string) bool {
	for _, check := range list {
		if check == item {
			return true
		}
	}
	return false
}

func inSliceInt64(item int64, list []int64) bool {
	for _, check := range list {
		if check == item {
			return true
		}
	}
	return false
}

func intersectSlices(left, right []string) []string {
	var result []string
	hash := make(map[string]struct{})

	for _, item := range left {
		hash[item] = struct{}{}
	}
	for _, item := range right {
		if _, ok := hash[item]; ok {
			result = append(result, item)
		}
	}
	return result
}

func getFieldValueByTag(input interface{}, tagValue string) reflect.Value {
	for i := 0; i < reflect.TypeOf(input).NumField(); i++ {
		field := reflect.TypeOf(input).Field(i)
		if field.Tag == reflect.StructTag("gorm:\"column:"+tagValue+"\"") {
			return reflect.ValueOf(input).FieldByName(field.Name)
		}
	}
	return reflect.Zero(reflect.TypeOf(input))
}

func getLastUpdateTime(fRuleName string, db *db.Database) time.Time {
	result := struct{ Cdate string }{}
	db.Table("rm_frule_history").
		Select("cdate").
		Order("cdate DESC").
		Where("rule_type = '" + fRuleName + "'").
		Limit(1).
		First(&result)
	cdate, err := time.ParseInLocation("2006-01-02 15:04:05", result.Cdate, time.Local)
	if err != nil {
		log.Logger.Error().Err(err).Msg("Time parsing")
	}
	return cdate
}

func timeMatchCronSpec(spec string, testTime time.Time) (bool, error) {
	if spec == "* * * * *" {
		return true, nil
	}
	schedule, err := cron.ParseStandard(spec)
	if err != nil {
		return false, err
	}
	converted := schedule.(*cron.SpecSchedule)
	monthTest := converted.Month&(1<<uint(testTime.Month())) != 0
	dayTest := converted.Dom&(1<<uint(testTime.Day())) != 0
	hourTest := converted.Hour&(1<<uint(testTime.Hour())) != 0
	minuteTest := converted.Minute&(1<<uint(testTime.Minute())) != 0
	weekDayTest := converted.Dow&(1<<uint(testTime.Weekday())) != 0
	result := monthTest && dayTest && hourTest && minuteTest && weekDayTest
	return result, nil
}

func priceRange(rangeSpec *string, testPrice base.Money) bool {
	if rangeSpec == nil {
		return false
	}
	rangeParts := strings.Split(*rangeSpec, "-")
	if rangeParts[1] == "" {
		value, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return false
		}
		return int64(value)*100 <= testPrice.Amount
	} else {
		lvalue, err := strconv.Atoi(rangeParts[0])
		if err != nil {
			return false
		}
		rvalue, err := strconv.Atoi(rangeParts[1])
		if err != nil {
			return false
		}
		return int64(lvalue)*100 <= testPrice.Amount && testPrice.Amount < int64(rvalue)*100
	}

}
func cronSpec(cronSpec *string, testTime time.Time) bool {
	if cronSpec == nil {
		return false
	}
	result, err := timeMatchCronSpec(*cronSpec, testTime)
	if err != nil {
		return false
	}
	return result
}

type cronStrucString struct {
	spec  string
	value string
}

type cronStrucBool struct {
	spec  string
	value bool
}