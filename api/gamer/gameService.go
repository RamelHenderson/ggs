package gamer

import (
	"gorm.io/gorm"
)

type IService interface {
	GetUserByID()
}
type Service struct {
	DB *gorm.DB
}

func (service *Service) GetUserByID(id uint) *Gamer {
	return &Gamer{}
}
