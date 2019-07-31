package frule_module

import (
	"github.com/robfig/cron"
	"reflect"
	"stash.tutu.ru/avia-search-common/contracts/base"
	"strconv"
	"strings"
	"time"
)

func InSlice(item string, list []string) bool {
	for _, check := range list {
		if check == item {
			return true
		}
	}
	return false
}

func InSliceInt64(item int64, list []int64) bool {
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
	var workingValue reflect.Type
	if reflect.TypeOf(input).Kind() == reflect.Ptr {
		workingValue = reflect.TypeOf(input).Elem()
	} else {
		workingValue = reflect.TypeOf(input)
	}
	for i := 0; i < workingValue.NumField(); i++ {
		field := workingValue.Field(i)
		if field.Tag == reflect.StructTag("json:\""+tagValue+"\"") {
			if reflect.TypeOf(input).Kind() == reflect.Ptr {
				return reflect.ValueOf(input).Elem().FieldByName(field.Name)
			} else {
				return reflect.ValueOf(input).FieldByName(field.Name)
			}
		}
	}
	return reflect.Zero(reflect.TypeOf(input))
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

func PriceRange(rangeSpec *string, testPrice base.Money) bool {
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

func CronSpec(cronSpec *string, testTime time.Time) bool {
	if cronSpec == nil {
		return false
	}
	result, err := timeMatchCronSpec(*cronSpec, testTime)
	if err != nil {
		return false
	}
	return result
}

type CronStructString struct {
	Spec  string `json:"key"`
	Value string `json:"value"`
}

type CronStructBool struct {
	Spec  string
	Value bool
}

