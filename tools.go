package frule_module

import (
	"reflect"
	"stash.tutu.ru/golang/log"
	"stash.tutu.ru/golang/resources/db"
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
		if field.Tag == reflect.StructTag("sql:\""+tagValue+"\"") {
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
