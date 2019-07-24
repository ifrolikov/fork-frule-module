package search_request

import (
	"errors"
	"fmt"
	"stash.tutu.ru/avia-search-common/repository"
)

type searchRequestRuleRankedList [][]*SearchRequestRule

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

func (i *importer) loadRankedRules() (searchRequestRuleRankedList, error) {
	data, loadErr := i.LoadURI(i.Config.DataURI)
	if loadErr != nil {
		return nil, errors.New(fmt.Sprintf("can't load SearchRequestRule data from %s: %v", i.Config.DataURI, loadErr))
	}
	var rankedList searchRequestRuleRankedList
	if parseErr := i.ParseData(data, &rankedList); parseErr != nil {
		return nil, errors.New(fmt.Sprintf("can't parse SearchRequestRule data from %s: %v", i.Config.DataURI, parseErr))
	}
	return rankedList, nil
}