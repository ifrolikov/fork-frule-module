package airline_restrictions

import (
"errors"
"fmt"
"stash.tutu.ru/avia-search-common/repository"
)

type airlineRestrictionRuleRankedList [][]*AirlineRestrictionsRule

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

func (i *importer) loadRankedRules() (airlineRestrictionRuleRankedList, error) {
	data, loadErr := i.LoadURI(i.Config.DataURI)
	if loadErr != nil {
		return nil, errors.New(fmt.Sprintf("can't load AirlineRestrictionRule data from %s: %v", i.Config.DataURI, loadErr))
	}
	var rankedList airlineRestrictionRuleRankedList
	if parseErr := i.ParseData(data, &rankedList); parseErr != nil {
		return nil, errors.New(fmt.Sprintf("can't parse AirlineRestrictionRule data from %s: %v", i.Config.DataURI, parseErr))
	}
	return rankedList, nil
}