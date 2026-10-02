package models

// PerfilUsuario espelha o enum `perfil_usuario` do Postgres.
// Ver plan-aply/fix-db-1.md secção 1.
//
// O valor 'OWNER' foi removido: significava "dono da oficina" e
// misturava gestão de sistema com operação de oficina.
type PerfilUsuario string

const (
	// PerfilADMIN gere os perfis dos MECHANIC e as suas contas.
	// Não atende, não executa, não entrega e não mexe no caixa.
	PerfilADMIN PerfilUsuario = "ADMIN"
	// PerfilMECHANIC é o staff de balcão: atendimento, execução,
	// entrega e caixa.
	PerfilMECHANIC PerfilUsuario = "MECHANIC"
	// PerfilCLIENTE é o dono do carro. Tem login, pede agendamento,
	// pede atendimento, confirma o orçamento e paga.
	PerfilCLIENTE PerfilUsuario = "CLIENTE"
)

// StaffPerfil é true para quem tem conta de funcionário da oficina.
// Serve de guarda para middleware e serviços.
func (p PerfilUsuario) StaffPerfil() bool {
	return p == PerfilADMIN || p == PerfilMECHANIC
}

type StatusAgendamento string

const (
	Agendado      StatusAgendamento = "AGENDADO"
	Cancelado     StatusAgendamento = "CANCELADO"
	NaoCompareceu StatusAgendamento = "NAO_COMPARECEU"
)

type StatusAtendimento string

const (
	EmAvaliacao          StatusAtendimento = "EM_AVALIACAO"
	AguardandoAprovacao  StatusAtendimento = "AGUARDANDO_APROVACAO"
	EmExecucao           StatusAtendimento = "EM_EXECUCAO"
	Encerrado            StatusAtendimento = "ENCERRADO"
	CanceladoAtendimento StatusAtendimento = "CANCELADO"
)

// TipoServico e TipoPeca chamam-se assim (e não Servico/Peca) para
// não colidirem com os models Servico e Peca, que são o catálogo.
type TipoItem string

const (
	TipoServico TipoItem = "SERVICO"
	TipoPeca    TipoItem = "PECA"
)

type StatusAprovacao string

const (
	Pendente StatusAprovacao = "PENDENTE"
	Aprovado StatusAprovacao = "APROVADO"
	Recusado StatusAprovacao = "RECUSADO"
)

type FormaPagamento string

const (
	Dinheiro FormaPagamento = "DINHEIRO"
	Pix      FormaPagamento = "PIX"
	Debito   FormaPagamento = "DEBITO"
	Credito  FormaPagamento = "CREDITO"
)

type StatusPagamento string

const (
	Confirmado   StatusPagamento = "CONFIRMADO"
	CanceladoPag StatusPagamento = "CANCELADO"
)
