package payment_method

import (
	"encoding/json"
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/golang/log"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(paymentMethodRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			if frule.DaysTillDeparture != nil {
				var daysTillDepartureParsed []frule_module.CronStructString
				err := json.Unmarshal([]byte(*frule.DaysTillDeparture), &daysTillDepartureParsed)
				if err != nil {
					log.Logger.Error().Stack().Err(err).Msg("Unmarshal DaysTillDeparture")
				}
				frule.DaysTillDepartureParsed = daysTillDepartureParsed
			}
			var resultParsed []frule_module.CronStructString
			err := json.Unmarshal([]byte(frule.Result), &resultParsed)
			if err != nil {
				log.Logger.Error().Stack().Err(err).Msg("Unmarshal PaymentMethod result")
			}
			frule.ResultParsed = resultParsed
			frulerList = append(frulerList, frule)
		}
		rankedFRuleStorage.Set(rank, frulerList)
	}
	container.rankedStorage = rankedFRuleStorage
}

func (container *fruleStorageContainer) GetRankedStorage() *frule_module.RankedFRuleStorage {
	return container.rankedStorage
}
