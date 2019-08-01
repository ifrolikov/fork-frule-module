package revenue

import (
	"encoding/json"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/golang/log"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(revenueRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			if frule.Revenue != nil && *frule.Revenue != "[]" {
				var revenueParsed Revenue
				if err := json.Unmarshal([]byte(*frule.Revenue), &revenueParsed); err != nil {
					log.Logger.Error().Err(err).Msg("Unmarshal revenue")
				}
				frule.RevenueParsed = &revenueParsed
			}
			if frule.Margin != nil && *frule.Margin != "[]" {
				var marginParsed Margin
				if err := json.Unmarshal([]byte(*frule.Margin), &marginParsed); err != nil {
					log.Logger.Error().Err(err).Msg("Unmarshal margin")
				}
				frule.MarginParsed = &marginParsed
			}
			frulerList = append(frulerList, frule)
		}
		rankedFRuleStorage.Set(rank, frulerList)
	}
	container.rankedStorage = rankedFRuleStorage
}

func (container *fruleStorageContainer) GetRankedStorage() *frule_module.RankedFRuleStorage {
	return container.rankedStorage
}
