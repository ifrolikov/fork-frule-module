package manual_exchange_refund

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"io/ioutil"
	"net/http"
	frule_module "stash.tutu.ru/avia-search-common/frule-module"
	"sync"
	"time"
)

// Описание импортера, в котором логика по импорту, сам никуда не ходит

type ComparisonOrderImporterInterface interface {
	getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error)
}

type comparisonOrderContainer struct {
	comparisonOrder frule_module.ComparisonOrder
}

type ComparisonOrderImporter struct {
	updateDuration           time.Duration
	lastUpdateTime           *time.Time
	comparisonOrderContainer *comparisonOrderContainer
	updater                  ComparisonOrderUpdaterInterface
}

func NewComparisonOrderImporter(
	updateDuration time.Duration,
	updater ComparisonOrderUpdaterInterface,
	defaultComparisonOrder *comparisonOrderContainer,
) *ComparisonOrderImporter {
	var lastUpdateTime *time.Time = nil
	if defaultComparisonOrder != nil {
		currentTime := time.Now()
		lastUpdateTime = &currentTime
	}
	return &ComparisonOrderImporter{
		updateDuration:           updateDuration,
		comparisonOrderContainer: defaultComparisonOrder,
		lastUpdateTime:           lastUpdateTime,
		updater:                  updater,
	}
}

func (importer *ComparisonOrderImporter) getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error) {
	if importer.comparisonOrderContainer == nil || importer.lastUpdateTime == nil {
		container, err := importer.updater.update(logger)
		if err != nil {
			return nil, err
		}
		currentTime := time.Now()
		importer.lastUpdateTime = &currentTime
		importer.comparisonOrderContainer = container
	} else {
		if time.Since(*importer.lastUpdateTime) > importer.updateDuration {
			currentTime := time.Now()
			importer.lastUpdateTime = &currentTime
			go func() {
				container, _ := importer.updater.update(logger)
				if container != nil {
					importer.comparisonOrderContainer = container
				}
			}()
		}
	}
	return importer.comparisonOrderContainer.comparisonOrder, nil
}

// Описание updater-а, который ходит в монолит

type ComparisonOrderUpdaterInterface interface {
	update(logger zerolog.Logger) (*comparisonOrderContainer, error)
}

type ComparisonOrderUpdater struct {
	apiUrl string
	client *http.Client
	mtx    *sync.Mutex
}

func NewComparisonOrderUpdater(apiUrl string) *ComparisonOrderUpdater {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
		Timeout: time.Duration(5 * time.Second),
	}

	mtx := &sync.Mutex{}
	return &ComparisonOrderUpdater{
		apiUrl: apiUrl,
		client: client,
		mtx:    mtx,
	}
}

func (updater *ComparisonOrderUpdater) update(logger zerolog.Logger) (*comparisonOrderContainer, error) {
	defer updater.mtx.Unlock()
	updater.mtx.Lock()

	resp, err := updater.client.Get(updater.apiUrl)
	if err != nil {
		logger.Err(err).Send()
		return nil, err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	logger.Debug().Msgf("response on importComparisonOrder\n:code %d\nbody: %s",
		resp.StatusCode,
		resp.Body)

	if resp.StatusCode != 200 {
		err := errors.New(fmt.Sprintf(
			"Invalid response from server on importComparisonOrder\ncode: %d\nbody: %s",
			resp.StatusCode,
			respBody))
		logger.Err(err).Send()
		return nil, err
	}

	comparisonOrder := frule_module.ComparisonOrder{}
	err = json.Unmarshal(respBody, &comparisonOrder)
	if err != nil {
		logger.Err(errors.Wrap(err, "Error by unmarshal response on importComparisonOrder")).Send()
		return nil, err
	}

	return &comparisonOrderContainer{comparisonOrder: comparisonOrder}, nil
}
