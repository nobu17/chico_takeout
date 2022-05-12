package store

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/store"
)

type SpecialHolidayModel struct {
	Id    string
	Name  string
	Start string
	End   string
}

func newSpecialHolidayModel(item *domains.SpecialHoliday) *SpecialHolidayModel {
	return &SpecialHolidayModel{
		Id:    item.GetId(),
		Name:  item.GetName(),
		Start: item.GetStart(),
		End:   item.GetEnd(),
	}
}

type SpecialHolidayCreateModel struct {
	Name  string
	Start string
	End   string
}

type SpecialHolidayUpdateModel struct {
	Id    string
	Name  string
	Start string
	End   string
}

type SpecialHolidayUseCase interface {
	Find(id string) (*SpecialHolidayModel, error)
	FindAll() ([]SpecialHolidayModel, error)
	Create(model *SpecialHolidayCreateModel) (string, error)
	Update(model *SpecialHolidayUpdateModel) error
	Delete(id string) error
}

type specialHolidayUseCase struct {
	repository     domains.SpecialHolidayRepository
	holidayService domains.HolidayService
}

func NewSpecialHolidayUseCase(repository domains.SpecialHolidayRepository) SpecialHolidayUseCase {
	return &specialHolidayUseCase{
		repository:     repository,
		holidayService: *domains.NewHolidayService(repository),
	}
}

func (i *specialHolidayUseCase) Find(id string) (*SpecialHolidayModel, error) {
	item, err := i.repository.Find(id)
	if err != nil {
		return nil, err
	}

	if item == nil {
		return nil, common.NewNotFoundError(id)
	}

	return newSpecialHolidayModel(item), nil
}

func (i *specialHolidayUseCase) FindAll() ([]SpecialHolidayModel, error) {
	items, err := i.repository.FindAll()
	if err != nil {
		return nil, err
	}

	models := []SpecialHolidayModel{}
	for _, item := range items {
		model := newSpecialHolidayModel(&item)
		models = append(models, *model)
	}

	return models, nil
}

func (s *specialHolidayUseCase) Create(model *SpecialHolidayCreateModel) (string, error) {
	item, err := domains.NewSpecialHoliday(model.Name, model.Start, model.End)
	if err != nil {
		return "", err
	}
	err = s.checkOverlap(item)
	if err != nil {
		return "", err
	}

	return s.repository.Create(item)
}

func (s *specialHolidayUseCase) Update(model *SpecialHolidayUpdateModel) error {
	item, err := s.repository.Find(model.Id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(model.Id)
	}

	err = item.Set(model.Name, model.Start, model.End)
	if err != nil {
		return err
	}

	err = s.checkOverlap(item)
	if err != nil {
		return err
	}

	return s.repository.Update(item)
}

func (i *specialHolidayUseCase) Delete(id string) error {
	item, err := i.repository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.repository.Delete(id)
}

func (s *specialHolidayUseCase) checkOverlap(item *domains.SpecialHoliday) error {
	isOverwrap, err := s.holidayService.CheckOverWrap(item)
	if err != nil {
		return err
	}
	if isOverwrap {
		return common.NewValidationError("special holiday", "date is overwrapped")
	}
	return nil
}
