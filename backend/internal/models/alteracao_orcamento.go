package models

import "time"

type AlteracaoOrcamento struct {
	ID                  string        `gorm:"type:uuid;primaryKey" json:"id"`
	ItemOrcamentoID     string        `gorm:"type:uuid;not null;index" json:"item_orcamento_id"`
	NovoItemOrcamentoID *string       `gorm:"type:uuid" json:"novo_item_orcamento_id"`
	ValorAnterior       float64       `gorm:"type:numeric(10,2);not null" json:"valor_anterior"`
	ValorNovo           float64       `gorm:"type:numeric(10,2);not null" json:"valor_novo"`
	Motivo              string        `gorm:"type:varchar(255)" json:"motivo"`
	UsuarioID           string        `gorm:"type:uuid;not null" json:"usuario_id"`
	CriadoEm            time.Time     `gorm:"not null" json:"criado_em"`
	CreatedAt           time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt           time.Time     `gorm:"not null" json:"updated_at"`

	ItemOrcamento     ItemOrcamento  `gorm:"foreignKey:ItemOrcamentoID" json:"item_orcamento,omitempty"`
	NovoItemOrcamento *ItemOrcamento `gorm:"foreignKey:novo_item_orcamento_id" json:"novo_item_orcamento,omitempty"`
	Usuario           Usuario        `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}
