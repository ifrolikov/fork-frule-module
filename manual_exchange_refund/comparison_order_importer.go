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

type ComparisonOrderImporterInterface interface {
	getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error)
}

type comparisonOrderContainer struct {
	comparisonOrder frule_module.ComparisonOrder
}

type ComparisonOrderImporter struct {
	apiUrl                   string
	updateDuration           time.Duration
	client                   *http.Client
	mtx                      *sync.Mutex
	lastUpdateTime           *time.Time
	comparisonOrderContainer *comparisonOrderContainer
}

func NewComparisonOrderImporter(apiUrl string, updateDuration time.Duration) *ComparisonOrderImporter {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{
		Transport: transport,
	}

	mtx := &sync.Mutex{}
	return &ComparisonOrderImporter{apiUrl: apiUrl, updateDuration: updateDuration, client: client, mtx: mtx}
}

func (importer *ComparisonOrderImporter) getComparisonOrder(logger zerolog.Logger) (frule_module.ComparisonOrder, error) {
	if importer.comparisonOrderContainer == nil || importer.lastUpdateTime == nil {
		if err := importer.importComparisonOrder(logger); err != nil {
			return nil, err
		}
	} else {
		if time.Since(*importer.lastUpdateTime) > importer.updateDuration {
			go func() {
				_ = importer.importComparisonOrder(logger)
			}()
		}
	}
	return importer.comparisonOrderContainer.comparisonOrder, nil
}

func (importer *ComparisonOrderImporter) importComparisonOrder(logger zerolog.Logger) error {
	defer importer.mtx.Unlock()
	importer.mtx.Lock()

	resp, err := importer.client.Get(importer.apiUrl)
	if err != nil {
		logger.Err(err).Send()
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
		return err
	}

	comparisonOrder := frule_module.ComparisonOrder{}
	err = json.Unmarshal(respBody, &comparisonOrder)
	if err != nil {
		logger.Err(errors.Wrap(err, "Error by unmarshal response on importComparisonOrder")).Send()
		return err
	}

	importer.comparisonOrderContainer = &comparisonOrderContainer{comparisonOrder: comparisonOrder}

	currentTime := time.Now()
	importer.lastUpdateTime = &currentTime

	return nil
}
