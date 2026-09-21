package models

import "time"

type Agendamento struct {
	ID        string            `gorm:"type:uuid;primaryKey" json:"id"`
	ClienteID string            `gorm:"type:uuid;not null;index" json:"cliente_id"`
	VeiculoID string            `gorm:"type:uuid;not null;index" json:"veiculo_id"`
	DataHora  time.Time         `gorm:"not null" json:"data_hora"`
	Status    StatusAgendamento `gorm:"type:status_agendamento;not null;default:AGENDADO" json:"status"`
	Versao    int               `gorm:"not null;default:1" json:"versao"`
	CreatedAt time.Time         `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time         `gorm:"not null" json:"updated_at"`

	Cliente Cliente `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	Veiculo Veiculo `gorm:"foreignKey:VeiculoID" json:"veiculo,omitempty"`
}
