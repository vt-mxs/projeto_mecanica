package models

import "time"

type Atendimento struct {
	ID                        string            `gorm:"type:uuid;primaryKey" json:"id"`
	ClienteID                 string            `gorm:"type:uuid;not null;index" json:"cliente_id"`
	VeiculoID                 string            `gorm:"type:uuid;not null;index" json:"veiculo_id"`
	AgendamentoID             *string           `gorm:"type:uuid;uniqueIndex" json:"agendamento_id"`
	ProblemaRelatado          string            `gorm:"type:text" json:"problema_relatado"`
	Avaliacao                 string            `gorm:"type:text" json:"avaliacao"`
	Diagnostico               string            `gorm:"type:text" json:"diagnostico"`
	Status                    StatusAtendimento `gorm:"type:status_atendimento;not null;default:EM_AVALIACAO" json:"status"`
	ExecucaoInicio            *time.Time        `json:"execucao_inicio"`
	ExecucaoFim               *time.Time        `json:"execucao_fim"`
	ExecucaoInterrompida      bool              `gorm:"not null;default:false" json:"execucao_interrompida"`
	ExecucaoObs               string            `gorm:"type:text" json:"execucao_obs"`
	EntregaData               *time.Time        `json:"entrega_data"`
	EntregaExcecao            bool              `gorm:"not null;default:false" json:"entrega_excecao"`
	EntregaExcecaoObs         string            `gorm:"type:text" json:"entrega_excecao_obs"`
	EntregaExcecaoUsuarioID   *string           `gorm:"type:uuid" json:"entrega_excecao_usuario_id"`
	CreatedAt                 time.Time         `gorm:"not null" json:"created_at"`
	UpdatedAt                 time.Time         `gorm:"not null" json:"updated_at"`

	Cliente                Cliente           `gorm:"foreignKey:ClienteID" json:"cliente,omitempty"`
	Veiculo                Veiculo           `gorm:"foreignKey:VeiculoID" json:"veiculo,omitempty"`
	Agendamento            *Agendamento      `gorm:"foreignKey:AgendamentoID" json:"agendamento,omitempty"`
	EntregaExcecaoUsuario  *Usuario          `gorm:"foreignKey:EntregaExcecaoUsuarioID" json:"entrega_excecao_usuario,omitempty"`
	VersoesOrcamento       []VersaoOrcamento `gorm:"foreignKey:AtendimentoID" json:"versoes_orcamento,omitempty"`
	ItensExecutados        []ItemExecutado   `gorm:"foreignKey:AtendimentoID" json:"itens_executados,omitempty"`
	Pagamentos             []Pagamento       `gorm:"foreignKey:AtendimentoID" json:"pagamentos,omitempty"`
}
