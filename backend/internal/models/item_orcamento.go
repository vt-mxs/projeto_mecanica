package models

import "time"

type ItemOrcamento struct {
	ID                   string           `gorm:"type:uuid;primaryKey" json:"id"`
	VersaoOrcamentoID    string           `gorm:"type:uuid;not null;index" json:"versao_orcamento_id"`
	Tipo                 TipoItem         `gorm:"type:tipo_item;not null" json:"tipo"`
	Descricao            string           `gorm:"type:varchar(160);not null" json:"descricao"`
	Valor                float64          `gorm:"type:numeric(10,2);not null" json:"valor"`
	StatusAprovacao      StatusAprovacao  `gorm:"type:status_aprovacao;not null;default:PENDENTE" json:"status_aprovacao"`
	DecididoPorUsuarioID *string          `gorm:"type:uuid" json:"decidido_por_usuario_id"`
	DecididoEm           *time.Time       `json:"decidido_em"`
	CreatedAt            time.Time        `gorm:"not null" json:"created_at"`
	UpdatedAt            time.Time        `gorm:"not null" json:"updated_at"`

	VersaoOrcamento    VersaoOrcamento      `gorm:"foreignKey:VersaoOrcamentoID" json:"versao_orcamento,omitempty"`
	DecididoPorUsuario *Usuario             `gorm:"foreignKey:DecididoPorUsuarioID" json:"decidido_por_usuario,omitempty"`
	AlteracoesOrigem   []AlteracaoOrcamento `gorm:"foreignKey:ItemOrcamentoID" json:"alteracoes_origem,omitempty"`
	AlteracoesNovo     []AlteracaoOrcamento `gorm:"foreignKey:novo_item_orcamento_id" json:"alteracoes_novo,omitempty"`
	ItensExecutados    []ItemExecutado      `gorm:"foreignKey:ItemOrcamentoID" json:"itens_executados,omitempty"`
}
