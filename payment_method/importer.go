package payment_method

import (
	"errors"
	"fmt"
	"stash.tutu.ru/avia-search-common/repository"
)

type paymentMethodRuleRankedList [][]*PaymentMethodRule

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

func (i *importer) loadRankedRules() (paymentMethodRuleRankedList, error) {
	data, loadErr := i.LoadURI(i.Config.DataURI)
	if loadErr != nil {
		return nil, errors.New(fmt.Sprintf("can't load FareRule data from %s: %v", i.Config.DataURI, loadErr))
	}
	var rankedList paymentMethodRuleRankedList
	if parseErr := i.ParseData(data, &rankedList); parseErr != nil {
		return nil, errors.New(fmt.Sprintf("can't parse FareRule data from %s: %v", i.Config.DataURI, parseErr))
	}
	return rankedList, nil
}
