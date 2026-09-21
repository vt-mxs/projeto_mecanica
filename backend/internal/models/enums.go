package models

type PerfilUsuario string

const (
	PerfilOWNER    PerfilUsuario = "OWNER"
	PerfilMECHANIC PerfilUsuario = "MECHANIC"
	PerfilADMIN    PerfilUsuario = "ADMIN"
)

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

type TipoItem string

const (
	Servico TipoItem = "SERVICO"
	Peca    TipoItem = "PECA"
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
