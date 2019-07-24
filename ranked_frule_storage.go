package frule_module

type RankedFRuleStorage map[int][]FRuler

func NewRankedFRuleStorage() *RankedFRuleStorage {
	s := make(RankedFRuleStorage)
	return &s
}

func (s *RankedFRuleStorage) Set(rank int, ruleSet []FRuler) {
	(*s)[rank] = ruleSet
}