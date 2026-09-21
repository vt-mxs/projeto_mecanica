package models

import "time"

type Pagamento struct {
	ID             string          `gorm:"type:uuid;primaryKey" json:"id"`
	AtendimentoID  string          `gorm:"type:uuid;not null;index" json:"atendimento_id"`
	Valor          float64         `gorm:"type:numeric(10,2);not null" json:"valor"`
	FormaPagamento FormaPagamento  `gorm:"type:forma_pagamento;not null" json:"forma_pagamento"`
	Data           time.Time       `gorm:"not null" json:"data"`
	UsuarioID      string          `gorm:"type:uuid;not null" json:"usuario_id"`
	Status         StatusPagamento `gorm:"type:status_pagamento;not null;default:CONFIRMADO" json:"status"`
	IdempotencyKey *string         `gorm:"type:varchar(64);uniqueIndex" json:"idempotency_key"`
	CreatedAt      time.Time       `gorm:"not null" json:"created_at"`
	UpdatedAt      time.Time       `gorm:"not null" json:"updated_at"`

	Atendimento Atendimento `gorm:"foreignKey:AtendimentoID" json:"atendimento,omitempty"`
	Usuario     Usuario     `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}
