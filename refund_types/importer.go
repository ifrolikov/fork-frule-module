package refund_types

import (
	"errors"
	"fmt"
	"stash.tutu.ru/avia-search-common/repository"
)

type commonRefundTypesRuleRankedList [][]*RefundTypesRule

type importer struct {
	repository.BasicImporter
}

func (i *importer) LoadData() (interface{}, error) {
	if rankedList, err := i.loadRankedRules(); err != nil {
		return nil, err
	} else {
		return rankedList, nil
	}
}

func (i *importer) loadRankedRules() (commonRefundTypesRuleRankedList, error) {
	data, loadErr := i.LoadURI(i.Config.DataURI)
	if loadErr != nil {
		return nil, errors.New(fmt.Sprintf("can't load ManualExchangeRefundRule data from %s: %v", i.Config.DataURI, loadErr))
	}
	var rankedList commonRefundTypesRuleRankedList
	if parseErr := i.ParseData(data, &rankedList); parseErr != nil {
		return nil, errors.New(fmt.Sprintf("can't parse ManualExchangeRefundRule data from %s: %v", i.Config.DataURI, parseErr))
	}
	return rankedList, nil
}