-- ============================================================
-- 000002_catalogo_auditoria.up.sql
--
-- Catálogo de serviços/peças com preços, e trilha de auditoria.
-- Ambas as coisas são pré-requisitos do painel de gestão: sem
-- catálogo não há lista de preços, e sem auditoria o painel de
-- logs do ADMIN não tem o que mostrar.
--
-- Fica separada da 000001 porque são concerns independentes: a
-- 000001 é o modelo de papéis, esta é functionality de negócio.
-- ============================================================

-- ============================================================
-- SERVICOS
-- ============================================================
CREATE TABLE servicos (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo              VARCHAR(40)  UNIQUE,
    nome                VARCHAR(160) NOT NULL,
    descricao           TEXT,
    preco_base          NUMERIC(10,2) NOT NULL CHECK (preco_base >= 0),
    tempo_estimado_min  INTEGER CHECK (tempo_estimado_min IS NULL OR tempo_estimado_min >= 0),
    ativo               BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_servicos_ativo ON servicos (ativo);

-- ============================================================
-- PECAS
-- ============================================================
-- preco_compra vs preco_venda é o que dá a margem no painel.
-- Não se impõe preco_venda >= preco_compra: inventário de
-- fundo/garantia vende-se abaixo do custo de propósito.
CREATE TABLE pecas (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    codigo        VARCHAR(40)  UNIQUE,
    nome          VARCHAR(160) NOT NULL,
    preco_compra  NUMERIC(10,2) NOT NULL CHECK (preco_compra >= 0),
    preco_venda   NUMERIC(10,2) NOT NULL CHECK (preco_venda  >= 0),
    ativo         BOOLEAN      NOT NULL DEFAULT TRUE,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pecas_ativo ON pecas (ativo);

-- ============================================================
-- LIGACAO CATALOGO <-> ITENS_ORCAMENTO
-- ============================================================
-- Opcional de propósito: o mechanic continua a poder escrever um
-- item avulso ("trocar fusível 5A") sem estar no catálogo. Mas
-- quando a origem está preenchida tem de concordar com `tipo`.
ALTER TABLE itens_orcamento
    ADD COLUMN servico_id UUID REFERENCES servicos(id) ON DELETE RESTRICT,
    ADD COLUMN peca_id    UUID REFERENCES pecas(id)    ON DELETE RESTRICT,
    ADD CONSTRAINT itens_orcamento_origem_coerente CHECK (
        (servico_id IS NULL AND peca_id IS NULL)
        OR (tipo = 'SERVICO' AND servico_id IS NOT NULL AND peca_id IS NULL)
        OR (tipo = 'PECA'    AND peca_id    IS NOT NULL AND servico_id IS NULL)
    );

CREATE INDEX idx_itens_orcamento_servico ON itens_orcamento (servico_id);
CREATE INDEX idx_itens_orcamento_peca    ON itens_orcamento (peca_id);

-- ============================================================
-- AUDITORIA
-- ============================================================
-- O que o painel de logs do ADMIN vai ler. Genérica de propósito:
-- cobrir só as operações da secção 26 do contrato com colunas
-- dedicadas obriga a uma migration sempre que aparece um caso novo.
--
-- `usuario_id` é NOT NULL: nenhuma acção entra na trilha sem actor
-- conhecido. Rotinas automáticas usam um utilizador de sistema
-- com perfil ADMIN, para não haver excepção ao CHECK.
--
-- Não há FK para `entidade_id` — a entidade alvo muda conforme o
-- `entidade`. Referenciar dez tabelas distintas aqui seria uma
-- armadilha de integridade referencial sem benefit: o registo
-- aponta para o que foi alterado, não o impede de o fazer.
CREATE TABLE auditoria (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id   UUID        NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    acao         VARCHAR(60) NOT NULL,
    entidade     VARCHAR(60) NOT NULL,
    entidade_id  UUID        NOT NULL,
    dados_antes  JSONB,
    dados_depois JSONB,
    ip           INET,
    user_agent   VARCHAR(255),
    criado_em    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auditoria_usuario_criado ON auditoria (usuario_id, criado_em DESC);
CREATE INDEX idx_auditoria_entidade        ON auditoria (entidade, entidade_id);
CREATE INDEX idx_auditoria_criado          ON auditoria (criado_em DESC);
CREATE INDEX idx_auditoria_acao            ON auditoria (acao);

COMMENT ON TABLE auditoria IS
    'Trilha de auditoria append-only. Nunca faz UPDATE nem DELETE: correccao faz-se por nova entrada.';