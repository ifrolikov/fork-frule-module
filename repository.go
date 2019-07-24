package frule_module

import (
	"context"
	"stash.tutu.ru/avia-search-common/repository"
)

type Repository repository.Repository

func NewFRuleRepository(ctx context.Context, storages fruleStorageContainer, importer repository.Importer) (*Repository, error) {
	if repo, err := repository.NewRepository(ctx, importer, storages); err != nil {
		return nil, err
	} else {
		return (*Repository)(repo), nil
	}
}

func (repo *Repository) GetRankedFRuleStorage() *RankedFRuleStorage {
	return repo.Storages.(fruleStorageContainer).GetRankedStorage()
}