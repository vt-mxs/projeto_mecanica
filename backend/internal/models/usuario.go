package models

import "time"

// IDs são string, não uuid.UUID, e isso é uma divergência real
// entre os models e o schema — ver o bloco "BUG CONHECIDO" em
// servico.go. Não Changing aqui porque o repository ainda não
// existe; corrigir enquanto se escreve o repository.
type Usuario struct {
	ID                     string        `gorm:"type:uuid;primaryKey" json:"id"`
	Login                  string        `gorm:"type:varchar(60);uniqueIndex;not null" json:"login"`
	Senha                  string        `gorm:"type:varchar(255);not null" json:"-"` // hash bcrypt, nunca sai em JSON
	Perfil                 PerfilUsuario `gorm:"type:perfil_usuario;not null" json:"perfil"`
	Nome                   string        `gorm:"type:varchar(120);not null;default:''" json:"nome"`
	Ativo                  bool          `gorm:"not null;default:true" json:"ativo"`
	UltimoLoginEm          *time.Time    `json:"ultimo_login_em"`
	DesativadoEm           *time.Time    `json:"desativado_em"`
	DesativadoPorUsuarioID *string       `gorm:"type:uuid" json:"desativado_por_usuario_id"`
	CreatedAt              time.Time     `gorm:"not null" json:"created_at"`
	UpdatedAt              time.Time     `gorm:"not null" json:"updated_at"`

	DesativadoPor *Usuario `gorm:"foreignKey:DesativadoPorUsuarioID" json:"desativado_por,omitempty"`
	Cliente       *Cliente `gorm:"foreignKey:UsuarioID" json:"cliente,omitempty"`
}

// UsuarioResponse é o DTO de saída. Nunca inclui a senha.
//
// O antigo `UserRepsonse` (typo) expunha `senha` em JSON e não tinha
// `Login` — numa tabela de login. Este substitui-o.
type UsuarioResponse struct {
	ID                     string        `json:"id"`
	Login                  string        `json:"login"`
	Perfil                 PerfilUsuario `json:"perfil"`
	Nome                   string        `json:"nome"`
	Ativo                  bool          `json:"ativo"`
	UltimoLoginEm          *time.Time    `json:"ultimo_login_em"`
	DesativadoEm           *time.Time    `json:"desativado_em"`
	DesativadoPorUsuarioID *string       `json:"desativado_por_usuario_id"`
	CreatedAt              time.Time     `json:"created_at"`
	UpdatedAt              time.Time     `json:"updated_at"`
}

func NewUsuarioResponse(u Usuario) UsuarioResponse {
	return UsuarioResponse{
		ID:                     u.ID,
		Login:                  u.Login,
		Perfil:                 u.Perfil,
		Nome:                   u.Nome,
		Ativo:                  u.Ativo,
		UltimoLoginEm:          u.UltimoLoginEm,
		DesativadoEm:           u.DesativadoEm,
		DesativadoPorUsuarioID: u.DesativadoPorUsuarioID,
		CreatedAt:              u.CreatedAt,
		UpdatedAt:              u.UpdatedAt,
	}
}
