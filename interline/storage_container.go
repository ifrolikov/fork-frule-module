package interline

import (
	"stash.tutu.ru/avia-search-common/frule-module"
	"strconv"
	"strings"
)

type fruleStorageContainer struct {
	rankedStorage *frule_module.RankedFRuleStorage
}

func (container *fruleStorageContainer) Update(data interface{}) {
	rankedFRuleStorage := frule_module.NewRankedFRuleStorage()
	for rank, ruleSet := range data.(interlineRuleRankedList) {
		frulerList := make([]frule_module.FRuler, 0, len(ruleSet))
		for _, frule := range ruleSet {
			if carriersNeedParsed, err := container.splitCarriersString(frule.CarriersNeed); err != nil {
				continue
			} else {
				frule.CarriersNeedParsed = carriersNeedParsed
			}
			if carriersForbidParsed, err := container.splitCarriersString(frule.CarriersForbid); err != nil {
				continue
			} else {
				frule.CarriersForbidParsed = carriersForbidParsed
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


func (container *fruleStorageContainer) splitCarriersString(carriersString string) ([]int64, error) {
	var carriersStringParsed []int64

	if carriersString != "" {
		for _, s := range strings.Split(carriersString, "|") {
			d, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, err
			}
			carriersStringParsed = append(carriersStringParsed, d)
		}
	}

	return carriersStringParsed, nil
}