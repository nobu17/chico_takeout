package message

import (
	domains "chico/takeout/domains/message"
	"chico/takeout/infrastructures/rdbms"
	"errors"

	"gorm.io/gorm"
)

type StoreMessageRepository struct {
	db *gorm.DB
}

func NewStoreMessageRepository(db *gorm.DB) *StoreMessageRepository {
	return &StoreMessageRepository{
		db: db,
	}
}

type StoreMessageModel struct {
	rdbms.BaseModel
	Content string
}

func (s *StoreMessageModel) toDomain() (*domains.StoreMessage, error) {
	model, err := domains.NewMessageForOrm(s.ID, s.Content, s.UpdatedAt, s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return model, nil
}

func newStoreMessageModel(s *domains.StoreMessage) *StoreMessageModel {
	model := StoreMessageModel{}
	model.ID = s.GetId()
	model.Content = s.GetContent()
	model.CreatedAt = s.GetCreated()
	model.UpdatedAt = s.GetEdited()
	return &model
}

func (s *StoreMessageRepository) Find(id string) (*domains.StoreMessage, error) {
	model := StoreMessageModel{}

	err := s.db.First(&model, "ID=?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	dom, err := model.toDomain()
	if err != nil {
		return nil, err
	}
	return dom, nil
}

func (s *StoreMessageRepository) Create(message *domains.StoreMessage) (string, error) {
	model := newStoreMessageModel(message)
	err := s.db.Create(&model).Error
	if err != nil {
		return "", err
	}
	return message.GetId(), nil
}

func (s *StoreMessageRepository) Update(message *domains.StoreMessage) error {
	model := newStoreMessageModel(message)
	err := s.db.Save(&model).Error
	return err
}
