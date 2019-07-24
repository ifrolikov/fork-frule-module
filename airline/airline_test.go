package airline

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/frule-module"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestNewAirlineFRule(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/airline.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewAirlineFRule(ctx, testConfig)
	assert.Nil(t, err)

	assert.Implements(t, (*frule_module.FRuler)(nil), frule)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 3)
	assert.Len(t, (*dataStorage)[3], 1)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, maxKey, 3)
}
