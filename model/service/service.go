package service

import (
	"sw/global"

	"gorm.io/gorm"
)

type ServiceModel struct {
	gorm.Model
	Opc string `json:"opc" yaml:"opc"`
}

type AddService struct {
	Opc string `json:"opc" yaml:"opc"`
}

type UpdateService struct {
	Id uint `json:"id"`
	AddService
}

func LoadAddService(add AddService) *ServiceModel {
	return &ServiceModel{
		Opc: add.Opc,
	}
}

func LoadUpdateService(update UpdateService) *ServiceModel {
	var s ServiceModel
	global.DB.First(&s, update.Id)
	s.Opc = update.Opc
	return &s
}

func (s *ServiceModel) Create() {
	global.DB.Create(s)
}

func (s *ServiceModel) Update() {
	global.DB.Save(s)
}

func (s *ServiceModel) Delete() {
	global.DB.Delete(s)
}