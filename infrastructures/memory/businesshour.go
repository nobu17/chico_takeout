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
		resetBusinessHoursMemory()
	}
	businessMux.Unlock()
	return &BusinessHoursMemoryRepository{inMemory: businessHoursMemory}
}

func resetBusinessHoursMemory() {
	mem, err := domains.NewDefaultBusinessHours()
	if err != nil {
		fmt.Println("%w", err)
		panic("failed to init businessHoursMemory")
	}
	businessHoursMemory = mem
}

func (b *BusinessHoursMemoryRepository) Reset() {
	resetBusinessHoursMemory()
}

func (s *BusinessHoursMemoryRepository) GetMemory() *domains.BusinessHours {
	return s.inMemory
}

func (b *BusinessHoursMemoryRepository) Fetch() (*domains.BusinessHours, error) {
	return b.inMemory, nil
}

func (b *BusinessHoursMemoryRepository) Update(target *domains.BusinessHours) error {
	*b.inMemory = *target
	return nil
}
