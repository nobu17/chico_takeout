package store

type StoreService struct {
	businessHoursRepository       BusinessHoursRepository
	specialBusinessHourRepository SpecialBusinessHourRepository
}

func NewStoreService(businessHoursRepository BusinessHoursRepository, specialBusinessHourRepository SpecialBusinessHourRepository) *StoreService {
	return &StoreService{
		businessHoursRepository:       businessHoursRepository,
		specialBusinessHourRepository: specialBusinessHourRepository,
	}
}

func (s *StoreService) FetchBusinessHours() (*BusinessHours, error) {
	businessHours, err := s.businessHoursRepository.Fetch()
	if err != nil {
		return nil, err
	}
	// create default
	if businessHours == nil {
		return NewDefaultBusinessHours()
	}
	return businessHours, nil
}

func (s *StoreService) ExistsBusinessHour(businessHourId string) (bool, error) {
	businessHours, err := s.businessHoursRepository.Fetch()
	if err != nil {
		return false, err
	}
	_, err = businessHours.FindById(businessHourId)
	if err != nil {
		return false, nil
	}
	return true, nil
}

type HolidayService struct {
	specialHolidayRepository SpecialHolidayRepository
}

func NewHolidayService(specialHolidayRepository SpecialHolidayRepository) *HolidayService {
	return &HolidayService{
		specialHolidayRepository: specialHolidayRepository,
	}
}

func (h *HolidayService) CheckOverWrap(item *SpecialHoliday) (bool, error) {
	holidays, err := h.specialHolidayRepository.FindAll()
	if err != nil {
		return false, err
	}
	for _, holiday := range holidays {
		// skip self
		if item.Equals(holiday) {
			continue
		}
		if item.IsOverlap(holiday) {
			return true, nil
		}
	}
	return false, nil
}
