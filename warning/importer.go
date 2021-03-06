package warning

import (
	"fmt"
	"stash.tutu.ru/avia-search-common/repository"
)

type warningRuleRankedList [][]*WarningRule

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

func (i *importer) loadRankedRules() (warningRuleRankedList, error) {
	data, loadErr := i.LoadURI(i.Config.DataURI)
	if loadErr != nil {
		return nil, fmt.Errorf("can't load WarningRule data from %s: %v", i.Config.DataURI, loadErr)
	}
	var rankedList warningRuleRankedList
	if parseErr := i.ParseData(data, &rankedList); parseErr != nil {
		return nil, fmt.Errorf("can't parse WarningRule data from %s: %v", i.Config.DataURI, parseErr)
	}
	return rankedList, nil
}