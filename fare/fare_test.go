package fare

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestNewFareFRule(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/fare.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewFareFRule(ctx, testConfig)
	assert.Nil(t, err)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 1)
	assert.Len(t, (*dataStorage)[16], 3)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, maxKey, 17)
}
