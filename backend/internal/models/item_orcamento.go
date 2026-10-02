package models

import "time"

type ItemOrcamento struct {
	ID                string          `gorm:"type:uuid;primaryKey" json:"id"`
	VersaoOrcamentoID string          `gorm:"type:uuid;not null;index" json:"versao_orcamento_id"`
	Tipo              TipoItem        `gorm:"type:tipo_item;not null" json:"tipo"`
	Descricao         string          `gorm:"type:varchar(160);not null" json:"descricao"`
	Valor             float64         `gorm:"type:numeric(10,2);not null" json:"valor"`
	StatusAprovacao   StatusAprovacao `gorm:"type:status_aprovacao;not null;default:PENDENTE" json:"status_aprovacao"`
	// DecididoPorUsuarioID é o CLIENTE que aprovou ou recusou.
	DecididoPorUsuarioID *string    `gorm:"type:uuid" json:"decidido_por_usuario_id"`
	DecididoEm           *time.Time `json:"decidido_em"`
	// DecisaoRegistradaPorUsuarioID é o MECHANIC que anotou a
	// decisão do cliente. O CHECK do Postgres exige os dois juntos
	// quando status_aprovacao <> 'PENDENTE'.
	DecisaoRegistradaPorUsuarioID *string   `gorm:"type:uuid" json:"decisao_registrada_por_usuario_id"`
	ServicoID                     *string   `gorm:"type:uuid" json:"servico_id"`
	PecaID                        *string   `gorm:"type:uuid" json:"peca_id"`
	CreatedAt                     time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt                     time.Time `gorm:"not null" json:"updated_at"`

	VersaoOrcamento             VersaoOrcamento      `gorm:"foreignKey:VersaoOrcamentoID" json:"versao_orcamento,omitempty"`
	DecididoPorUsuario          *Usuario             `gorm:"foreignKey:DecididoPorUsuarioID" json:"decidido_por_usuario,omitempty"`
	DecisaoRegistradaPorUsuario *Usuario             `gorm:"foreignKey:DecisaoRegistradaPorUsuarioID" json:"decisao_registrada_por_usuario,omitempty"`
	Servico                     *Servico             `gorm:"foreignKey:ServicoID" json:"servico,omitempty"`
	Peca                        *Peca                `gorm:"foreignKey:PecaID" json:"peca,omitempty"`
	AlteracoesOrigem            []AlteracaoOrcamento `gorm:"foreignKey:ItemOrcamentoID" json:"alteracoes_origem,omitempty"`
	AlteracoesNovo              []AlteracaoOrcamento `gorm:"foreignKey:NovoItemOrcamentoID" json:"alteracoes_novo,omitempty"`
	ItensExecutados             []ItemExecutado      `gorm:"foreignKey:ItemOrcamentoID" json:"itens_executados,omitempty"`
}
