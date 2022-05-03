package store

type StoreService struct {
	businessHoursRepository BusinessHoursRepository
}

func NewStoreService(businessHoursRepository BusinessHoursRepository) *StoreService {
	return &StoreService{
		businessHoursRepository: businessHoursRepository,
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
