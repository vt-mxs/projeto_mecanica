package models

import "time"

type Veiculo struct {
	ID        string    `gorm:"type:uuid;primaryKey" json:"id"`
	ClienteID string    `gorm:"type:uuid;not null;index" json:"cliente_id"`
	Modelo    string    `gorm:"type:varchar(80)" json:"modelo"`
	Cor       string    `gorm:"type:varchar(30)" json:"cor"`
	Placa     string    `gorm:"type:varchar(10);uniqueIndex" json:"placa"`
	Ano       *int      `gorm:"type:int" json:"ano"`
	Km        *int64    `gorm:"type:bigint" json:"km"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Cliente      Cliente       `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	Agendamentos []Agendamento `gorm:"foreignKey:VeiculoID" json:"agendamentos,omitempty"`
	Atendimentos []Atendimento `gorm:"foreignKey:VeiculoID" json:"atendimentos,omitempty"`
}
