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
					log.Logger.Error().Stack().Err(err).Msg("Unmarshal revenue")
				}
				for idx := range revenueParsed.Full {
					revenueParsed.Full[idx].Result.SegmentParsed = parseMoneySpec(revenueParsed.Full[idx].Result.Segment)
					revenueParsed.Full[idx].Result.TicketParsed = parseMoneySpec(revenueParsed.Full[idx].Result.Ticket)
				}
				for idx := range revenueParsed.Child {
					revenueParsed.Child[idx].Result.SegmentParsed = parseMoneySpec(revenueParsed.Child[idx].Result.Segment)
					revenueParsed.Child[idx].Result.TicketParsed = parseMoneySpec(revenueParsed.Child[idx].Result.Ticket)
				}
				for idx := range revenueParsed.Infant {
					revenueParsed.Infant[idx].Result.SegmentParsed = parseMoneySpec(revenueParsed.Infant[idx].Result.Segment)
					revenueParsed.Infant[idx].Result.TicketParsed = parseMoneySpec(revenueParsed.Infant[idx].Result.Ticket)
				}
				frule.RevenueParsed = &revenueParsed
			}
			if frule.Margin != nil && *frule.Margin != "[]" {
				var marginParsed Margin
				if err := json.Unmarshal([]byte(*frule.Margin), &marginParsed); err != nil {
					log.Logger.Error().Stack().Err(err).Msg("Unmarshal margin")
				}
				for idx := range marginParsed.Full {
					marginParsed.Full[idx].ResultParsed = parseMoneySpec(&marginParsed.Full[idx].Result).Money
				}
				for idx := range marginParsed.Child {
					marginParsed.Child[idx].ResultParsed = parseMoneySpec(&marginParsed.Child[idx].Result).Money
				}
				for idx := range marginParsed.Infant {
					marginParsed.Infant[idx].ResultParsed = parseMoneySpec(&marginParsed.Infant[idx].Result).Money
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
