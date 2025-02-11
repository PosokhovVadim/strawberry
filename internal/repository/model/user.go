package model

import (
	"time"

	"github.com/PosokhovVadim/stawberry/internal/domain/entity"
)

type User struct {
	ID        uint `gorm:"primaryKey"`
	Name      string
	Email     string `gorm:"unique"`
	Password  string
	IsStore   bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ConvertUserFromSvc(u entity.Register) User {
	return User{
		Name:     u.Username,
		Email:    u.Email,
		Password: u.Password,
		// IsStore:  u.IsStore,
	}
}

func ConvertUserToEntity(m User) entity.User {
	return entity.User{
		Id:        int(m.ID),
		Name:      m.Name,
		Email:     m.Email,
		Password:  m.Password,
		IsStore:   m.IsStore,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
