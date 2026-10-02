package models

import "time"

// ╔══════════════════════════════════════════════════════════════╗
// ║ BUG CONHECIDO — ID string em coluna uuid, em TODOS os models  ║
// ╚══════════════════════════════════════════════════════════════╝
//
// O schema define `id UUID PRIMARY KEY DEFAULT gen_random_uuid()`.
// Os models definem `ID string`. O GORM, ao não ver valor, insere
// a string vazia em vez de omitir a coluna, e o Postgres rejeita:
//
//	ERROR: invalid input syntax for type uuid: "" (SQLSTATE 22P02)
//
// Isto afecta TODOS os models, não só este, e não aparece em
// `go build` nem em `go vet` — só quando se tenta um INSERT real.
// Não é introduzido pela 000001 nem pela 000002: já existia no
// primeiro commit dos models.
//
// Duas saídas, para quando o repository for escrito:
//
//  1. Trocar `ID string` por `ID uuid.UUID` + `github.com/google/uuid`,
//     que já está no go.mod. Gera o UUID na aplicação e o tipo
//     bate certo com o schema.
//
//  2. Manter `string` e marcar `gorm:"default:gen_random_uuid()"` em
//     todos os campos ID, para o GORM omita a coluna e deixe o
//     DEFAULT do Postgres actuar.
//
// A (2) foi TESTADA contra o Postgres real e funciona: o GORM passa
// a emitir `INSERT INTO servicos (nome, preco_base, ativo) VALUES
// (...) RETURNING id`, sem a coluna id, e o Postgres devolve o UUID.
// Note-se que sem a tag, o GORM emite a coluna mesmo assim, a com
// `''`.
//
// A (1) é a preferível: deixa o schema ser a única fonte de verdade
// dos defaults e o tipo do Go alinhado com o tipo do Postgres.
//
// NOTA: isto NÃO foi corrigido agora porque não fazia parte do
// plano e o ficheiro de repository ainda não existe — a decisão de
// onde os IDs são gerados pertence ao repository.

type Servico struct {
	ID               string    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Codigo           *string   `gorm:"type:varchar(40);uniqueIndex" json:"codigo"`
	Nome             string    `gorm:"type:varchar(160);not null" json:"nome"`
	Descricao        *string   `json:"descricao"`
	PrecoBase        float64   `gorm:"type:numeric(10,2);not null" json:"preco_base"`
	TempoEstimadoMin *int      `json:"tempo_estimado_min"`
	Ativo            bool      `gorm:"not null;default:true" json:"ativo"`
	CreatedAt        time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt        time.Time `gorm:"not null" json:"updated_at"`
}
