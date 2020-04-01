package models

import (
	"time"

	"github.com/google/uuid"
)

func NewUsers() *Users {
	m := new(Users)
	m.Id = uuid.Must(uuid.NewRandom())
	return m
}

func (m *Users) TableName() string {
	return "users"
}

type Users struct {
	Id           uuid.UUID `json:"Id" binding:"-"`
	Username     string    `json:"Username" binding:"required"`
	Password     string    `json:"Password" binding:"required"`
	Nama_lengkap string    `json:"Nama_lengkap" binding:"required"`
	Created_at   time.Time `json:"Created_at" binding:"-"`
	Updated_at   time.Time `json:"Updated_at" binding:"-"`
}
