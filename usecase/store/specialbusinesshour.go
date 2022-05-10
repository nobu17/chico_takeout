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

type SpecialBusinessHoursUseCase interface {
	Find(id string) (*SpecialBusinessHourModel, error) 
	FindAll() ([]SpecialBusinessHourModel, error)
	Create(model *SpecialBusinessHourCreateModel) (string, error)
	Update(model *SpecialBusinessHourUpdateModel) error
	Delete(id string) error
}

type specialBusinessHoursUseCase struct {
	repository   domains.SpecialBusinessHourRepository
	bussinesHoursService domains.BusinessHoursService
}

func NewSpecialBusinessHoursUseCase(
	businessHoursRepository domains.BusinessHoursRepository,
	specialBusinessHourRepository domains.SpecialBusinessHourRepository) SpecialBusinessHoursUseCase {
	return &specialBusinessHoursUseCase{
		repository:   specialBusinessHourRepository,
		bussinesHoursService: *domains.NewBussinessHoursService(businessHoursRepository),
	}
}

func (s *specialBusinessHoursUseCase) Find(id string) (*SpecialBusinessHourModel, error) {
	item, err := s.repository.Find(id)
	if err != nil {
		return nil, err
	}

	return newSpecialBusinessHourModel(item), nil
}

func (s *specialBusinessHoursUseCase) FindAll() ([]SpecialBusinessHourModel, error) {
	items, err := s.repository.FindAll()
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

func (s *specialBusinessHoursUseCase) Create(model *SpecialBusinessHourCreateModel) (string, error) {
	item, err := domains.NewSpecialBusinessHour(model.Name, model.Date, model.Start, model.End, model.BusinessHourId)
	if err != nil {
		return "", err
	}

	err = s.validate(item)
	if err != nil {
		return "", err
	}

	return s.repository.Create(item)
}

func (s *specialBusinessHoursUseCase) Update(model *SpecialBusinessHourUpdateModel) error {
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

	err = s.validate(item)
	if err != nil {
		return err
	}

	return s.repository.Update(item)
}

func (i *specialBusinessHoursUseCase) Delete(id string) error {
	item, err := i.repository.Find(id)
	if err != nil {
		return err
	}
	if item == nil {
		return common.NewUpdateTargetNotFoundError(id)
	}

	return i.repository.Delete(id)
}

func (s *specialBusinessHoursUseCase) validate(item *domains.SpecialBusinessHour) error {
	// check business hour id exists
	exists, err := s.bussinesHoursService.ExistsBusinessHour(item.GetBusinessHourId())
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
