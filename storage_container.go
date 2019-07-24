package frule_module

type fruleStorageContainer interface {
	Update(data interface{})
	GetRankedStorage() *RankedFRuleStorage
}