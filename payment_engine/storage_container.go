package payment_engine

import (
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(paymentEngineRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			frulerList = append(frulerList, frule)
		}
		rankedFRuleStorage.Set(rank, frulerList)
	}
	container.rankedStorage = rankedFRuleStorage
}

func (container *fruleStorageContainer) GetRankedStorage() *frule_module.RankedFRuleStorage {
	return container.rankedStorage
}