package memory

import (
	domains "chico/takeout/domains/store"
)

var businessHoursMemory *domains.BusinessHours

type BusinessHoursMemoryRepository struct {
	inMemory *domains.BusinessHours
}

func NewBusinessHoursMemoryRepository() *BusinessHoursMemoryRepository {
	return &BusinessHoursMemoryRepository{inMemory: businessHoursMemory}
}

func (b *BusinessHoursMemoryRepository) Fetch() (*domains.BusinessHours, error) {
	return b.inMemory, nil
}

func (b *BusinessHoursMemoryRepository) Update(target domains.BusinessHours) error {
	*b.inMemory = target
	return nil
}
