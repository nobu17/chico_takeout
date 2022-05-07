package memory

import (
	"fmt"
	"sync"

	domains "chico/takeout/domains/store"
)

var businessMux sync.Mutex
var businessHoursMemory *domains.BusinessHours

type BusinessHoursMemoryRepository struct {
	inMemory *domains.BusinessHours
}

func NewBusinessHoursMemoryRepository() *BusinessHoursMemoryRepository {
	businessMux.Lock()
	if businessHoursMemory == nil {
		mem, err := domains.NewDefaultBusinessHours()
		if err != nil {
			fmt.Println("%w", err)
			panic("failed to init businessHoursMemory")
		}
		businessHoursMemory = mem
	}
	businessMux.Unlock()
	return &BusinessHoursMemoryRepository{inMemory: businessHoursMemory}
}

func (b *BusinessHoursMemoryRepository) Fetch() (*domains.BusinessHours, error) {
	return b.inMemory, nil
}

func (b *BusinessHoursMemoryRepository) Update(target *domains.BusinessHours) error {
	*b.inMemory = *target
	return nil
}
