package frule_module

import (
	"stash.tutu.ru/avia-search-common/contracts/v2/base"
	"testing"
	"time"
)

func TestCronSpec(t *testing.T) {
	testCases := []struct {
		Spec     string
		TestTime string
		Result   bool
	}{
		// min hour dom month dow
		{"* * * * *", "2016-02-18 23:55:00", true},
		{"* * 10 * *", "2016-02-18 23:55:00", false},
		{"* * 18 * *", "2016-02-18 23:55:00", true},
		{"57 * * * *", "2016-02-18 23:55:00", false},
		{"*/5 * * * *", "2016-02-18 23:55:00", true},
		{"40-59 * * * *", "2016-02-18 23:55:00", true},
		{"* * * * 4", "2016-02-18 23:55:00", true},
		{"40-59 * * * 5", "2016-02-18 23:55:00", false},
		{"40-59 10-18 * * *", "2016-02-18 23:55:00", false},
		{"40-59 10-23 1-5,17-18 * *", "2016-02-18 23:55:00", true},
		{"* * */3 * *", "2016-02-19 23:55:00", true},
	}

	for _, testCase := range testCases {
		timeToTest, err := time.Parse("2006-01-02 15:04:05", testCase.TestTime)
		if err != nil {
			t.Fatal(err)
		}
		result, err := timeMatchCronSpec(testCase.Spec, timeToTest)
		if err != nil {
			t.Errorf("Failed invokin function on palceholder %s, err = %v", testCase.Spec, err)
		}
		if result != testCase.Result {
			t.Errorf("Cronspec %s returned %t against time %s, want %t", testCase.Spec, result, testCase.TestTime, testCase.Result)
		}
	}
}

func TestPriceRange(t *testing.T) {
	testCases := []struct {
		Range      string
		BaseAmount base.Money
		Result     bool
	}{
		{Range: "400-5000", BaseAmount: base.Money{Amount: 50000}, Result: true},
		{Range: "400-5000", BaseAmount: base.Money{Amount: 40000}, Result: true},
		{Range: "1000-5000", BaseAmount: base.Money{Amount: 50000}, Result: false},
		{Range: "400-500", BaseAmount: base.Money{Amount: 50000}, Result: false},
		{Range: "0-", BaseAmount: base.Money{Amount: 50000}, Result: true},
		{Range: "500-", BaseAmount: base.Money{Amount: 50000}, Result: true},
	}
	for _, testCase := range testCases {
		if res := PriceRange(&testCase.Range, testCase.BaseAmount); res != testCase.Result {
			t.Errorf("Failed check price range %s on amount %d, want %t, got %t", testCase.Range, testCase.BaseAmount.Amount, testCase.Result, res)
		}
	}
}
