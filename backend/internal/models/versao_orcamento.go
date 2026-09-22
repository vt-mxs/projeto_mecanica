package models

import "time"

type VersaoOrcamento struct {
	ID                   string    `gorm:"type:uuid;primaryKey" json:"id"`
	AtendimentoID        string    `gorm:"type:uuid;not null;index" json:"atendimento_id"`
	NumeroVersao         int       `gorm:"not null" json:"numero_versao"`
	CriadoEm             time.Time `gorm:"not null" json:"criado_em"`
	CriadoPorUsuarioID   string    `gorm:"type:uuid;not null" json:"criado_por_usuario_id"`
	CreatedAt            time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt            time.Time `gorm:"not null" json:"updated_at"`

	Atendimento      Atendimento     `gorm:"foreignKey:AtendimentoID" json:"atendimento,omitempty"`
	CriadoPorUsuario Usuario         `gorm:"foreignKey:CriadoPorUsuarioID" json:"criado_por_usuario,omitempty"`
	ItensOrcamento   []ItemOrcamento `gorm:"foreignKey:VersaoOrcamentoID" json:"itens_orcamento,omitempty"`
}
