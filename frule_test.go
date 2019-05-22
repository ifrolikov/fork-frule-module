package frule_module

import (
	"context"
	"testing"
	"time"
)

type DummyFRule struct {
	CarrierId       *int    `sql:"carrier_id"`
	Partner         *string `sql:"partner"`
	ConnectionGroup *string `sql:"connection_group"`
	Result          bool    `sql:"result"`
}

func (a DummyFRule) GetResultValue() interface{} {
	return a.Result
}

func (a DummyFRule) GetDataStorage() (map[int][]FRuler, error) {
	result := make(map[int][]FRuler)
	carrierId := 10
	connectionGroup := "test"
	connectionGroup2 := "test2"
	partner := "fake"
	partner2 := "fake2"
	result[0] = []FRuler{
		DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Partner: &partner, Result: true},
		DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup2, Partner: &partner, Result: true},
	}
	result[1] = []FRuler{
		DummyFRule{ConnectionGroup: &connectionGroup, Partner: &partner, Result: false},
		DummyFRule{ConnectionGroup: &connectionGroup2, Partner: &partner, Result: false},
	}
	result[3] = []FRuler{
		DummyFRule{Partner: &partner2, Result: true},
	}
	return result, nil
}

func (a DummyFRule) GetComparisonOrder() ComparisonOrder {
	return ComparisonOrder{
		[]string{"carrier_id", "partner", "connection_group"},
		[]string{"partner", "connection_group"},
		[]string{"carrier_id", "partner"},
		[]string{"partner"},
	}
}

func (a DummyFRule) GetComparisonOperators() ComparisonOperators {
	return ComparisonOperators{}
}

func (a DummyFRule) getStrategyKeys() []string {
	return []string{"carrier_id", "partner", "connection_group"}
}

func (a DummyFRule) getTableName() string {
	return ""
}

func (a DummyFRule) GetDefaultValue() interface{} {
	return false
}

func (a DummyFRule) GetLastUpdateTime() time.Time {
	return time.Now()
}

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

func TestCreateHash(t *testing.T) {
	definition := &FRule{
		ruleSpecificData: DummyFRule{},
		indexedKeys:      []string{"carrier_id", "partner", "connection_group"},
	}

	carrier := 15
	cgroup := "fake_group"
	partner := "fake_partner"
	testFRule := DummyFRule{
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
		hashFields := intersectSlices(definition.indexedKeys, definition.ruleSpecificData.GetComparisonOrder()[i])
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
		hashFields := intersectSlices(definition.indexedKeys, definition.ruleSpecificData.GetComparisonOrder()[i])
		hash := definition.createRuleHash(hashFields, testFRule)
		if hash != correct[i] {
			t.Errorf("Failed to calculate hash, got %s, expected %s", hash, correct[i])
		}
	}
}

func TestFRule(t *testing.T) {
	ctx := context.Background()
	frule := NewFRule(ctx, DummyFRule{})
	carrierId := 10
	carrierId2 := 5
	connectionGroup := "test"
	connectionGroup2 := "test2"
	partner := "fake"
	partner2 := "fake2"
	partner3 := "fake3"

	results := []struct {
		testRule DummyFRule
		result   bool
	}{
		{testRule: DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup, Partner: &partner}, result: true},
		{testRule: DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup2, Partner: &partner}, result: true},
		{testRule: DummyFRule{CarrierId: &carrierId2, ConnectionGroup: &connectionGroup2, Partner: &partner}, result: false},
		{testRule: DummyFRule{CarrierId: &carrierId2, ConnectionGroup: &connectionGroup2, Partner: &partner2}, result: true},
		{testRule: DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup2, Partner: &partner2}, result: true},
		{testRule: DummyFRule{CarrierId: &carrierId, ConnectionGroup: &connectionGroup2, Partner: &partner3}, result: false},
	}

	for idx, testDef := range results {
		if testDef.result != frule.GetResult(testDef.testRule).(bool) {
			t.Errorf("Failed to get frule for iteration %d", idx)
		}
	}
}
