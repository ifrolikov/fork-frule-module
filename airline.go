package frule_module

type AirlineRule struct {
	Id              int     `sql:"id"`
	CarrierId       *int    `sql:"carrier_id"`
	Partner         *string `sql:"partner"`
	ConnectionGroup *string `sql:"connection_group"`
	Result          bool    `sql:"result"`
}

func (a AirlineRule) GetResultValue() interface{} {
	return a.Result
}

func (a AirlineRule) GetContainer() FRuler {
	return AirlineRule{}
}

func (a AirlineRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"carrier_id", "partner", "connection_group"},
		[]string{"partner", "connection_group"},
		[]string{"carrier_id", "partner"},
		[]string{"partner"},
	}
}

func (a AirlineRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a AirlineRule) GetStrategyKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a AirlineRule) GetIndexedKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a AirlineRule) GetTableName() string {
	return "rm_frule_airline"
}

func (a AirlineRule) GetDefaultValue() interface{} {
	return false
}
