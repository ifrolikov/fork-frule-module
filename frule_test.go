package frule_module

import (
	"testing"
)

func TestIntersect(t *testing.T) {
	left := []string{"one", "two", "tree"}
	right := []string{"one", "four"}
	result := intersectSlices(left, right)
	if len(result) != 1 {
		t.Error("Invalid length")
	}
	if result[0] != "one" {
		t.Error("Invalid intersect")
	}

	left = []string{"One", "Two"}
	right = []string{"one", "two"}
	result = intersectSlices(left, right)

	if len(result) != 0 {
		t.Error("Invalid length in mixed cases")
	}

	left = []string{"dfg", "hbtrhb"}
	right = []string{"dfg", "hbtrhb"}
	result = intersectSlices(left, right)

	if len(result) != 2 {
		t.Error("Invalid length in same cases")
	}

	correct := []string{"dfg", "hbtrhb"}
	for idx, item := range result {
		if item != correct[idx] {
			t.Errorf("Error checking result at %d, got %s - expected %s", idx, item, correct[idx])
		}
	}
}

type TestFRule struct {
	CarrierId       *int    `sql:"carrier_id"`
	Partner         *string `sql:"partner"`
	ConnectionGroup *string `sql:"connection_group"`
	Result          bool    `sql:"result"`
}

func (a TestFRule) GetResultValue() interface{} {
	return 0
}

func (a TestFRule) GetContainer() FRuler {
	return TestFRule{}
}

func (a TestFRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"carrier_id", "partner", "connection_group"},
		[]string{"partner", "connection_group"},
		[]string{"carrier_id", "partner"},
		[]string{"partner"},
	}
}

func (a TestFRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a TestFRule) GetStrategyKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a TestFRule) GetIndexedKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a TestFRule) GetTableName() string {
	return "rm_frule_airline"
}

func (a TestFRule) GetDefaultValue() interface{} {
	return false
}

func TestCreateHash(t *testing.T) {
	definition := &FRule{
		ruleSpecificData: TestFRule{},
	}

	carrier := 15
	cgroup := "fake_group"
	partner := "fake_partner"
	testFRule := TestFRule{
		CarrierId:       &carrier,
		ConnectionGroup: &cgroup,
		Partner:         &partner,
		Result:          true,
	}

	correct := []string{
		"carrier_id=>15|partner=>fake_partner|connection_group=>fake_group|",
		"partner=>fake_partner|connection_group=>fake_group|",
		"carrier_id=>15|partner=>fake_partner|",
		"partner=>fake_partner|",
	}

	for i := 0; i < len(definition.ruleSpecificData.GetComparisonOrder()); i++ {
		hashFields := intersectSlices(definition.ruleSpecificData.GetIndexedKeys(), definition.ruleSpecificData.GetComparisonOrder()[i])
		hash := definition.createRuleHash(hashFields, testFRule)
		if hash != correct[i] {
			t.Errorf("Failed to calculate hash, got %s, expected %s", hash, correct[i])
		}
	}

	testFRule.ConnectionGroup = nil

	correct = []string{
		"carrier_id=>15|partner=>fake_partner|",
		"partner=>fake_partner|",
		"carrier_id=>15|partner=>fake_partner|",
		"partner=>fake_partner|",
	}

	for i := 0; i < len(definition.ruleSpecificData.GetComparisonOrder()); i++ {
		hashFields := intersectSlices(definition.ruleSpecificData.GetIndexedKeys(), definition.ruleSpecificData.GetComparisonOrder()[i])
		hash := definition.createRuleHash(hashFields, testFRule)
		if hash != correct[i] {
			t.Errorf("Failed to calculate hash, got %s, expected %s", hash, correct[i])
		}
	}
}
