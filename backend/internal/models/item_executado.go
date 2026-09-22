package models

import "time"

type ItemExecutado struct {
	ID              string     `gorm:"type:uuid;primaryKey" json:"id"`
	AtendimentoID   string     `gorm:"type:uuid;not null;index" json:"atendimento_id"`
	ItemOrcamentoID *string    `gorm:"type:uuid" json:"item_orcamento_id"`
	Tipo            TipoItem   `gorm:"type:tipo_item;not null" json:"tipo"`
	Descricao       string     `gorm:"type:varchar(160);not null" json:"descricao"`
	ValorCobrado    float64    `gorm:"type:numeric(10,2);not null" json:"valor_cobrado"`
	CustoOficina    *float64   `gorm:"type:numeric(10,2)" json:"custo_oficina"`
	CreatedAt       time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"not null" json:"updated_at"`

	Atendimento   Atendimento    `gorm:"foreignKey:AtendimentoID" json:"atendimento,omitempty"`
	ItemOrcamento *ItemOrcamento `gorm:"foreignKey:ItemOrcamentoID" json:"item_orcamento,omitempty"`
}
