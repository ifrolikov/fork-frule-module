package frule_module

type RankedFRuleStorage map[int][]FRuler

func NewRankedFRuleStorage() *RankedFRuleStorage {
	s := make(RankedFRuleStorage)
	return &s
}

func (s *RankedFRuleStorage) Set(rank int, ruleSet []FRuler) {
	(*s)[rank] = ruleSet
}

func (s *RankedFRuleStorage) GetMaxRank() int {
	maxKey := 0
	for key := range *s {
		if key > maxKey {
			maxKey = key
		}
	}
	return maxKey
}