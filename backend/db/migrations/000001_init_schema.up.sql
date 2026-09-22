-- ============================================================
-- 000001_init_schema.up.sql
-- Criação de todas as tabelas do sistema da oficina
-- ============================================================

-- Tipos enumerados
CREATE TYPE perfil_usuario AS ENUM ('OWNER', 'MECHANIC');
CREATE TYPE status_agendamento AS ENUM ('AGENDADO', 'CANCELADO', 'NAO_COMPARECEU');
CREATE TYPE status_atendimento AS ENUM ('EM_AVALIACAO', 'AGUARDANDO_APROVACAO', 'EM_EXECUCAO', 'ENCERRADO', 'CANCELADO');
CREATE TYPE tipo_item AS ENUM ('SERVICO', 'PECA');
CREATE TYPE status_aprovacao AS ENUM ('PENDENTE', 'APROVADO', 'RECUSADO');
CREATE TYPE forma_pagamento AS ENUM ('DINHEIRO', 'PIX', 'DEBITO', 'CREDITO');
CREATE TYPE status_pagamento AS ENUM ('CONFIRMADO', 'CANCELADO');

-- ============================================================
-- USUARIO
-- ============================================================
CREATE TABLE usuarios (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login      VARCHAR(60)  NOT NULL UNIQUE,
    senha      VARCHAR(255) NOT NULL,
    perfil     perfil_usuario NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

-- ============================================================
-- CLIENTE
-- ============================================================
CREATE TABLE clientes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    nome       VARCHAR(120) NOT NULL,
    telefone   VARCHAR(20)  NOT NULL,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_clientes_telefone ON clientes (telefone);

-- ============================================================
-- VEICULO
-- ============================================================
CREATE TABLE veiculos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cliente_id UUID         NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    modelo     VARCHAR(80),
    cor        VARCHAR(30),
    placa      VARCHAR(10)  UNIQUE,
    ano        INTEGER,
    km         BIGINT,
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_veiculos_placa     ON veiculos (placa);
CREATE INDEX idx_veiculos_cliente_id ON veiculos (cliente_id);

-- ============================================================
-- AGENDAMENTO
-- ============================================================
CREATE TABLE agendamentos (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cliente_id UUID                NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    veiculo_id UUID                NOT NULL REFERENCES veiculos(id) ON DELETE RESTRICT,
    data_hora  TIMESTAMPTZ         NOT NULL,
    status     status_agendamento  NOT NULL DEFAULT 'AGENDADO',
    versao     INTEGER             NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_agendamentos_data_hora   ON agendamentos (data_hora);
CREATE INDEX idx_agendamentos_veiculo     ON agendamentos (veiculo_id, data_hora);

-- ============================================================
-- ATENDIMENTO
-- ============================================================
CREATE TABLE atendimentos (
    id                          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cliente_id                  UUID                NOT NULL REFERENCES clientes(id) ON DELETE RESTRICT,
    veiculo_id                  UUID                NOT NULL REFERENCES veiculos(id) ON DELETE RESTRICT,
    agendamento_id              UUID                UNIQUE REFERENCES agendamentos(id) ON DELETE RESTRICT,
    problema_relatado           TEXT,
    avaliacao                   TEXT,
    diagnostico                 TEXT,
    status                      status_atendimento  NOT NULL DEFAULT 'EM_AVALIACAO',
    execucao_inicio             TIMESTAMPTZ,
    execucao_fim                TIMESTAMPTZ,
    execucao_interrompida       BOOLEAN             NOT NULL DEFAULT FALSE,
    execucao_obs                TEXT,
    entrega_data                TIMESTAMPTZ,
    entrega_excecao             BOOLEAN             NOT NULL DEFAULT FALSE,
    entrega_excecao_obs         TEXT,
    entrega_excecao_usuario_id  UUID                REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at                  TIMESTAMPTZ         NOT NULL DEFAULT NOW(),
    updated_at                  TIMESTAMPTZ         NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_atendimentos_cliente_id  ON atendimentos (cliente_id);
CREATE INDEX idx_atendimentos_veiculo_id  ON atendimentos (veiculo_id);
CREATE INDEX idx_atendimentos_status      ON atendimentos (status);

-- ============================================================
-- VERSAO_ORCAMENTO
-- ============================================================
CREATE TABLE versoes_orcamento (
    id                     UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    atendimento_id         UUID        NOT NULL REFERENCES atendimentos(id) ON DELETE RESTRICT,
    numero_versao          INTEGER     NOT NULL,
    criado_em              TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    criado_por_usuario_id  UUID        NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at             TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at             TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_versoes_orcamento_atendimento ON versoes_orcamento (atendimento_id);

-- ============================================================
-- ITEM_ORCAMENTO
-- ============================================================
CREATE TABLE itens_orcamento (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    versao_orcamento_id       UUID               NOT NULL REFERENCES versoes_orcamento(id) ON DELETE RESTRICT,
    tipo                      tipo_item          NOT NULL,
    descricao                 VARCHAR(160)       NOT NULL,
    valor                     NUMERIC(10,2)      NOT NULL CHECK (valor >= 0),
    status_aprovacao          status_aprovacao   NOT NULL DEFAULT 'PENDENTE',
    decidido_por_usuario_id   UUID               REFERENCES usuarios(id) ON DELETE RESTRICT,
    decidido_em               TIMESTAMPTZ,
    created_at                TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ        NOT NULL DEFAULT NOW(),

    CHECK (
        (status_aprovacao = 'PENDENTE' AND decidido_por_usuario_id IS NULL AND decidido_em IS NULL)
        OR
        (status_aprovacao <> 'PENDENTE' AND decidido_por_usuario_id IS NOT NULL AND decidido_em IS NOT NULL)
    )
);

CREATE INDEX idx_itens_orcamento_versao ON itens_orcamento (versao_orcamento_id);

-- ============================================================
-- ALTERACAO_ORCAMENTO
-- ============================================================
CREATE TABLE alteracoes_orcamento (
    id                        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    item_orcamento_id         UUID           NOT NULL REFERENCES itens_orcamento(id) ON DELETE RESTRICT,
    novo_item_orcamento_id    UUID           REFERENCES itens_orcamento(id) ON DELETE RESTRICT,
    valor_anterior            NUMERIC(10,2)  NOT NULL,
    valor_novo                NUMERIC(10,2)  NOT NULL,
    motivo                    VARCHAR(255),
    usuario_id                UUID           NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    criado_em                 TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    created_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CHECK (valor_anterior <> valor_novo)
);

CREATE INDEX idx_alteracoes_orcamento_item ON alteracoes_orcamento (item_orcamento_id);

-- ============================================================
-- ITEM_EXECUTADO
-- ============================================================
CREATE TABLE itens_executados (
    id                  UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    atendimento_id      UUID           NOT NULL REFERENCES atendimentos(id) ON DELETE RESTRICT,
    item_orcamento_id   UUID           REFERENCES itens_orcamento(id) ON DELETE RESTRICT,
    tipo                tipo_item      NOT NULL,
    descricao           VARCHAR(160)   NOT NULL,
    valor_cobrado       NUMERIC(10,2)  NOT NULL CHECK (valor_cobrado >= 0),
    custo_oficina       NUMERIC(10,2)  CHECK (custo_oficina IS NULL OR custo_oficina >= 0),
    created_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CHECK (tipo = 'PECA' OR custo_oficina IS NULL)
);

CREATE INDEX idx_itens_executados_atendimento ON itens_executados (atendimento_id);

-- ============================================================
-- PAGAMENTO
-- ============================================================
CREATE TABLE pagamentos (
    id                 UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    atendimento_id     UUID               NOT NULL REFERENCES atendimentos(id) ON DELETE RESTRICT,
    valor              NUMERIC(10,2)      NOT NULL CHECK (valor > 0),
    forma_pagamento    forma_pagamento    NOT NULL,
    data               TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    usuario_id         UUID               NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    status             status_pagamento   NOT NULL DEFAULT 'CONFIRMADO',
    idempotency_key    VARCHAR(64)        UNIQUE,
    created_at         TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at         TIMESTAMPTZ        NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_pagamentos_atendimento   ON pagamentos (atendimento_id);
CREATE INDEX idx_pagamentos_data          ON pagamentos (data);
CREATE INDEX idx_pagamentos_idempotency   ON pagamentos (idempotency_key);
