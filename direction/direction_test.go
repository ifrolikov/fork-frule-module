package direction

import (
	"context"
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"stash.tutu.ru/avia-search-common/repository"
	"testing"
)

func TestDirection(t *testing.T) {
	pwd, _ := filepath.Abs("./")
	testConfig := &repository.Config{
		DataURI: "file://" + pwd + "/../testdata/direction.json",
	}
	ctx := context.Background()
	defer ctx.Done()

	frule, err := NewDirectionFRule(ctx, testConfig)
	assert.Nil(t, err)

	dataStorage := frule.GetDataStorage()
	assert.NotNil(t, dataStorage)
	assert.Len(t, (*dataStorage)[0], 2)
	assert.Len(t, (*dataStorage)[1], 1)

	maxKey := 0
	for key := range *dataStorage {
		if key > maxKey {
			maxKey = key
		}
	}
	assert.Equal(t, maxKey, 7)
}
