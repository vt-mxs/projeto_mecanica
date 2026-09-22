package models

import "time"

type Usuario struct {
	ID        string        `gorm:"type:uuid;primaryKey" json:"id"`
	Login     string        `gorm:"type:varchar(60);uniqueIndex;not null" json:"login"`
	Senha     string        `gorm:"type:varchar(255);not null" json:"senha"`
	Perfil    PerfilUsuario `gorm:"type:perfil_usuario;not null" json:"perfil"`
	CreatedAt time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time     `gorm:"not null" json:"updated_at"`
}

type UserRepsonse struct {
	ID        string        `gorm:"type:uuid;primaryKey" json:"id"`
	Senha     string        `gorm:"type:varchar(255);not null" json:"senha"`
	Perfil    PerfilUsuario `gorm:"type:perfil_usuario;not null" json:"perfil"`
	CreatedAt time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time     `gorm:"not null" json:"updated_at"`
}
