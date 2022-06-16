package store

type BusinessHoursService struct {
	businessHoursRepository BusinessHoursRepository
}

func NewBusinessHoursService(businessHoursRepository BusinessHoursRepository) *BusinessHoursService {
	return &BusinessHoursService{
		businessHoursRepository: businessHoursRepository,
	}
}

func (s *BusinessHoursService) ExistsBusinessHour(businessHourId string) (bool, error) {
	businessHours, err := s.businessHoursRepository.Fetch()
	if err != nil {
		return false, err
	}
	item := businessHours.FindById(businessHourId)
	if item == nil {
		return false, nil
	}
	return true, nil
}

func (s *BusinessHoursService) FetchBusinessHours() (*BusinessHours, error) {
	businessHours, err := s.businessHoursRepository.Fetch()
	if err != nil {
		return nil, err
	}
	// create default
	if businessHours == nil {
		hours, err := NewDefaultBusinessHours()
		if err != nil {
			return nil, err
		}
		err = s.businessHoursRepository.Create(hours)
		if err != nil {
			return nil, err
		}
		businessHours, err = s.businessHoursRepository.Fetch()
		if err != nil {
			return nil, err
		}
	}
	return businessHours, nil
}

func (s *BusinessHoursService) InitDataIfNotExists() error {
	businessHours, err := s.businessHoursRepository.Fetch()
	if err != nil {
		return err
	}
	// create default
	if businessHours == nil {
		hours, err := NewDefaultBusinessHours()
		if err != nil {
			return err
		}
		err = s.businessHoursRepository.Create(hours)
		if err != nil {
			return err
		}
	}
	return nil
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
