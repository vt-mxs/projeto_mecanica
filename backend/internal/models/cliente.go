package models

import "time"

type Cliente struct {
	ID string `gorm:"type:uuid;primaryKey" json:"id"`
	// UsuarioID liga o cliente à sua conta de login. NULL = cliente
	// de balcão, ainda sem portal.
	//
	// O Postgres não valida que este usuário tem perfil CLIENTE
	// (CHECK não aceita subquery). É invariante da camada de
	// aplicação — ver plan-aply/fix-db-1.md secção 3.
	UsuarioID *string   `gorm:"type:uuid;uniqueIndex" json:"usuario_id"`
	Nome      string    `gorm:"type:varchar(120);not null" json:"nome"`
	Telefone  string    `gorm:"type:varchar(20);not null" json:"telefone"`
	CreatedAt time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`

	Usuario      *Usuario      `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
	Veiculos     []Veiculo     `gorm:"foreignKey:ClienteID" json:"veiculos,omitempty"`
	Agendamentos []Agendamento `gorm:"foreignKey:ClienteID" json:"agendamentos,omitempty"`
	Atendimentos []Atendimento `gorm:"foreignKey:ClienteID" json:"atendimentos,omitempty"`
}
