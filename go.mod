module stash.tutu.ru/avia-search-common/frule-module

go 1.12

require (
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/pkg/errors v0.8.1
	github.com/robfig/cron v1.2.0
	github.com/stretchr/testify v1.4.0
	stash.tutu.ru/avia-search-common/contracts v1.1.0
	stash.tutu.ru/avia-search-common/repository v0.1.2
	stash.tutu.ru/avia-search-common/utils v0.2.2
	stash.tutu.ru/golang/log v0.0.0-20190925120345-16f0e05f99a0
)

replace stash.tutu.ru/avia-search-common/contracts => ../../../stash.tutu.ru/avia-search-common/contracts
