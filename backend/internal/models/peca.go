package models

import "time"

// Peca é a peça do catálogo. preco_compra vs preco_venda é o que dá a
// margem no painel de gestão.
//
// O CHECK `pecas_margem_nao_negativa` (preco_venda >= preco_compra) NÃO
// existe de propósito: inventário de fundo ou garantia vende-se abaixo
// do custo sem ser erro. A verificação, se for preciso, é da aplicação.
type Peca struct {
	ID          string    `gorm:"type:uuid;primaryKey" json:"id"`
	Codigo      *string   `gorm:"type:varchar(40);uniqueIndex" json:"codigo"`
	Nome        string    `gorm:"type:varchar(160);not null" json:"nome"`
	PrecoCompra float64   `gorm:"type:numeric(10,2);not null" json:"preco_compra"`
	PrecoVenda  float64   `gorm:"type:numeric(10,2);not null" json:"preco_venda"`
	Ativo       bool      `gorm:"not null;default:true" json:"ativo"`
	CreatedAt   time.Time `gorm:"not null" json:"created_at"`
	UpdatedAt   time.Time `gorm:"not null" json:"updated_at"`
}

// Margem devolve o lucro unitário e a margem percentual sobre o custo.
// preco_compra = 0 não dá divisão; nesse caso a margem percentual
// devolve 0 para não aparecer NaN no painel.
func (p Peca) Margem() (lucro float64, percentual float64) {
	lucro = p.PrecoVenda - p.PrecoCompra
	if p.PrecoCompra > 0 {
		percentual = lucro / p.PrecoCompra * 100
	}
	return lucro, percentual
}
