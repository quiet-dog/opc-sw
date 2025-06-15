package service

import (
	"fmt"
	"sw/global"
	"sw/opc"
	"time"

	"gorm.io/gorm"
)

type ServiceModel struct {
	gorm.Model
	Opc      string `json:"opc" yaml:"opc"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

func (s *ServiceModel) AfterCreate(tx *gorm.DB) (err error) {
	err = global.OpcGateway.AddClinet(fmt.Sprintf("%d", s.ID), opc.OpcClient{
		Endpoint: s.Opc,
		// 不断地读取数据
		Duration: time.Second * 600000000,
		Username: s.Username,
		Password: s.Password,
	})
	return
}

type AddService struct {
	Opc      string `json:"opc" yaml:"opc"`
	Username string `json:"username" yaml:"username"`
	Password string `json:"password" yaml:"password"`
}

type UpdateService struct {
	Id uint `json:"id"`
	AddService
}

func LoadAddService(add AddService) *ServiceModel {
	return &ServiceModel{
		Opc:      add.Opc,
		Username: add.Username,
		Password: add.Password,
	}
}

func LoadUpdateService(update UpdateService) *ServiceModel {
	var s ServiceModel
	global.DB.First(&s, update.Id)
	s.Opc = update.Opc
	s.Username = update.Username
	s.Password = update.Password
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
