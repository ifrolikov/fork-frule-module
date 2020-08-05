package warning

import (
	"github.com/ifrolikov/fork-frule-module"
	"stash.tutu.ru/golang/log"
	"time"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(warningRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			if frule.StartDate != nil {
				if parsedTime, err := time.Parse("2006-01-02", *frule.StartDate); err != nil {
					log.For("warning").Error().Stack().Err(err).Msg("parsing StartDate")
				} else {
					frule.ParsedStartDate = &parsedTime
				}
			}

			if frule.FinishDate != nil {
				if parsedTime, err := time.Parse("2006-01-02", *frule.FinishDate); err != nil {
					log.For("warning").Error().Stack().Err(err).Msg("parsing FinishDate")
				} else {
					frule.ParsedFinishDate = &parsedTime
				}
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