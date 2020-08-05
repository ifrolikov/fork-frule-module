package search_connection

import (
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/golang/log"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(searchConnectionRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			var err error
			frule.MinDepartureDateParsed, err = frule.parseCronSpecField(*frule.MinDepartureDate)
			if err != nil {
				log.For("search_connection").Error().Stack().Err(err).Msg("parsing minDepartureDate")
			}
			frule.MaxDepartureDateParsed, err = frule.parseCronSpecField(*frule.MaxDepartureDate)
			if err != nil {
				log.For("search_connection").Error().Stack().Err(err).Msg("parsing maxDepartureDate")
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
