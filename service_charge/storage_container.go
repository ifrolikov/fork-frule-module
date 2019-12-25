package service_charge

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
	for rank, ruleSet := range data.(serviceChargeRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			if frule.Margin != nil && *frule.Margin != "[]" {
				var marginParsed Margin
				if err := json.Unmarshal([]byte(*frule.Margin), &marginParsed); err != nil {
					log.Logger.Error().Stack().Err(err).Msg("Unmarshal margin")
				}
				for idx := range marginParsed.Full {
					marginParsed.Full[idx].ResultParsed = parseMoneySpec(&marginParsed.Full[idx].Result)
				}
				for idx := range marginParsed.Child {
					marginParsed.Child[idx].ResultParsed = parseMoneySpec(&marginParsed.Child[idx].Result)
				}
				for idx := range marginParsed.Infant {
					marginParsed.Infant[idx].ResultParsed = parseMoneySpec(&marginParsed.Infant[idx].Result)
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
