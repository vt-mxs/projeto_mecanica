package models

import "time"

type Cliente struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	Nome      string    `gorm:"type:varchar(120);not null" json:"nome"`
	Telefone  string    `gorm:"type:varchar(20);not null" json:"telefone"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Veiculos     []Veiculo     `gorm:"foreignKey:ClienteID" json:"veiculos,omitempty"`
	Agendamentos []Agendamento `gorm:"foreignKey:ClienteID" json:"agendamentos,omitempty"`
	Atendimentos []Atendimento `gorm:"foreignKey:ClienteID" json:"atendimentos,omitempty"`
}
