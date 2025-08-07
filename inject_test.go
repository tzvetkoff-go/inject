package inject_test

import (
	"fmt"

	"github.com/tzvetkoff-go/inject"
)

type DatabaseConnection struct {
	DSN string
}

type ItemRepository interface {
	GetItem() error
}

type StandardItemRepository struct {
	DatabaseConnection *DatabaseConnection `inject:"DatabaseConnection"`
}

func (r *StandardItemRepository) GetItem() error {
	return nil
}

type CachingItemRepository struct {
	StandardItemRepository ItemRepository `inject:"ItemRepository"`
}

func (r *CachingItemRepository) GetItem() error {
	return nil
}

type ItemService struct {
	ItemRepository ItemRepository `inject:"ItemRepository"`
}

func ExampleInjector_Inject() {
	injector := inject.New()

	// Setup some stuff ...
	databaseConnection := &DatabaseConnection{
		DSN: "host=127.0.0.1 username=root password= database=w00t schema=public",
	}
	standardItemRepository := &StandardItemRepository{}
	cachingItemRepository := &CachingItemRepository{}
	itemService := &ItemService{}

	// Register everything to the injector ...
	injector.ProvideObject("DatabaseConnection", databaseConnection)
	// For ItemRepository we want to do some test shenanigans ...
	injector.Provide("ItemRepository", func(injectInto interface{}) interface{} {
		if _, ok := injectInto.(*CachingItemRepository); ok {
			return standardItemRepository
		}

		return cachingItemRepository
	})
	injector.ProvideObject("ItemService", itemService)

	// Inject them all ...
	if err := injector.Inject(
		databaseConnection,
		standardItemRepository,
		cachingItemRepository,
		itemService,
	); err != nil {
		panic(err)
	}

	// Test ...
	if _, ok := itemService.ItemRepository.(*CachingItemRepository); !ok {
		panic("itemService.itemRepository should be of type *CachingItemRepository")
	}

	fmt.Printf("%T\n", injector.GetObject("DatabaseConnection", nil))
	fmt.Printf("%T\n", cachingItemRepository.StandardItemRepository)
	fmt.Printf("%T\n", itemService.ItemRepository)

	// Output:
	// *inject_test.DatabaseConnection
	// *inject_test.StandardItemRepository
	// *inject_test.CachingItemRepository
}
