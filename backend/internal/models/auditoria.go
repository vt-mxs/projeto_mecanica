package models

import (
	"encoding/json"
	"time"
)

// Auditoria é a trilha append-only que o painel de logs do ADMIN lê.
//
// Nunca faz UPDATE nem DELETE: uma correcção é uma entrada nova.
// UsuarioID é NOT NULL — nenhuma acção entra na trilha sem actor
// conhecido. Rotinas automáticas usam um utilizador de sistema com
// perfil ADMIN.
type Auditoria struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	UsuarioID   string    `gorm:"type:uuid;not null" json:"usuario_id"`
	Acao        string    `gorm:"type:varchar(60);not null" json:"acao"`
	Entidade    string    `gorm:"type:varchar(60);not null" json:"entidade"`
	EntidadeID  string    `gorm:"type:uuid;not null" json:"entidade_id"`
	DadosAntes  []byte    `gorm:"type:jsonb" json:"dados_antes"`
	DadosDepois []byte    `gorm:"type:jsonb" json:"dados_depois"`
	IP          *string   `gorm:"type:inet" json:"ip"`
	UserAgent   *string   `gorm:"type:varchar(255)" json:"user_agent"`
	CriadoEm    time.Time `gorm:"not null" json:"criado_em"`

	Usuario Usuario `gorm:"foreignKey:UsuarioID" json:"usuario,omitempty"`
}

// Ações registadas. Values em minúsculas para o filtro do painel
// poder ser exacto.
const (
	AcaoCriarUsuario       = "criar_usuario"
	AcaoDesativarUsuario   = "desativar_usuario"
	AcaoReativarUsuario    = "reativar_usuario"
	AcaoAlterarPerfil      = "alterar_perfil"
	AcaoDecisaoOrcamento   = "decisao_orcamento"
	AcaoAlteracaoOrcamento = "alteracao_orcamento"
	AcaoPagamento          = "registar_pagamento"
	AcaoCancelarPagamento  = "cancelar_pagamento"
	AcaoExecucao           = "execucao"
	AcaoEntrega            = "entrega"
	AcaoEntregaExcecao     = "entrega_excecao"
)

// NewAuditoria monta uma entrada. before/after podem ser structs,
// maps ou qualquer coisa que o json.Marshal aceite; nil fica NULL.
//
// ip e userAgent são opcionais porque nem toda a acção vem de um
// request (rotinas automáticas não têm IP).
func NewAuditoria(usuarioID, acao, entidade, entidadeID string, antes, depois any, ip, userAgent *string) (*Auditoria, error) {
	a := Auditoria{
		UsuarioID:  usuarioID,
		Acao:       acao,
		Entidade:   entidade,
		EntidadeID: entidadeID,
		IP:         ip,
		UserAgent:  userAgent,
	}

	if antes != nil {
		b, err := json.Marshal(antes)
		if err != nil {
			return nil, err
		}
		a.DadosAntes = b
	}
	if depois != nil {
		b, err := json.Marshal(depois)
		if err != nil {
			return nil, err
		}
		a.DadosDepois = b
	}

	return &a, nil
}
