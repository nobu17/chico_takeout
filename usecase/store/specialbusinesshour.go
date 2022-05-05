package store

import (
	"chico/takeout/common"
	domains "chico/takeout/domains/store"
	"fmt"
)

type SpecialBusinessHourModel struct {
	Id             string
	Name           string
	Date           string
	Start          string
	End            string
	BusinessHourId string
}

func newSpecialBusinessHourModel(item *domains.SpecialBusinessHour) *SpecialBusinessHourModel {
	return &SpecialBusinessHourModel{
		Id:             item.GetId(),
		Name:           item.GetName(),
		Date:           item.GetDate(),
		Start:          item.GetStart(),
		End:            item.GetEnd(),
		BusinessHourId: item.GetBusinessHourId(),
	}
}

type SpecialBusinessHourCreateModel struct {
	Name           string
	Date           string
	Start          string
	End            string
	BusinessHourId string
}

type SpecialBusinessHourUpdateModel struct {
	Id             string
	Name           string
	Date           string
	Start          string
	End            string
	BusinessHourId string
}

type SpecialBusinessHoursUseCase struct {
	repository   domains.SpecialBusinessHourRepository
	storeService domains.StoreService
}

func NewSpecialBusinessHoursUseCase(
	businessHoursRepository domains.BusinessHoursRepository,
	specialBusinessHourRepository domains.SpecialBusinessHourRepository) *SpecialBusinessHoursUseCase {
	return &SpecialBusinessHoursUseCase{
		repository:   specialBusinessHourRepository,
		storeService: *domains.NewStoreService(businessHoursRepository, specialBusinessHourRepository),
	}
}

func (i *SpecialBusinessHoursUseCase) Find(id string) (*SpecialBusinessHourModel, error) {
	item, err := i.repository.Find(id)
	if err != nil {
		return nil, err
	}

	return newSpecialBusinessHourModel(item), nil
}

func (i *SpecialBusinessHoursUseCase) FindAll() ([]SpecialBusinessHourModel, error) {
	items, err := i.repository.FindAll()
	if err != nil {
		return nil, err
	}

	models := []SpecialBusinessHourModel{}
	for _, item := range items {
		model := newSpecialBusinessHourModel(&item)
		models = append(models, *model)
	}

	return models, nil
}

func (s *SpecialBusinessHoursUseCase) Create(model SpecialBusinessHourCreateModel) (string, error) {
	item, err := domains.NewSpecialBusinessHour(model.Name, model.Date, model.Start, model.End, model.BusinessHourId)
	if err != nil {
		return "", err
	}

	err = s.validate(*item)
	if err != nil {
		return "", err
	}

	return s.repository.Create(*item)
}

func (s *SpecialBusinessHoursUseCase) Update(model SpecialBusinessHourUpdateModel) error {
	item, err := s.repository.Find(model.Id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError("specialBusinessHour")
	}

	err = item.Set(model.Name, model.Date, model.Start, model.End, model.BusinessHourId)
	if err != nil {
		return err
	}

	err = s.validate(*item)
	if err != nil {
		return err
	}

	return s.repository.Update(*item)
}

func (i *SpecialBusinessHoursUseCase) Delete(id string) error {
	item, err := i.repository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.repository.Delete(id)
}

func (s *SpecialBusinessHoursUseCase) validate(item domains.SpecialBusinessHour) error {
	// check business hour id exists
	exists, err := s.storeService.ExistsBusinessHour(item.GetBusinessHourId())
	if err != nil {
		return err
	}
	if !exists {
		return common.NewValidationError("businesHourId", fmt.Sprintf("id not exists:%s", item.GetBusinessHourId()))
	}

	all, err := s.repository.FindAll()
	if err != nil {
		return err
	}

	spec := domains.NewSpecialBusinessHourSpecification(all)
	err = spec.Validate(item)
	if err != nil {
		return err
	}
	return nil
}
