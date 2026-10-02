-- ============================================================
-- 000001_init_schema.up.sql
-- Criação de todas as tabelas do sistema da oficina
--
-- Modelo de papéis (ver plan-aply/modDB.md):
--   ADMIN    -> gestor do sistema. Gere os perfis dos MECHANIC
--               e as suas contas. Não opera a oficina.
--   MECHANIC -> staff de balcão. Atende, executa, entrega, caixa.
--   CLIENTE  -> dono do carro. Tem login, pede agendamento,
--               pede atendimento e confirma/recusa o orçamento.
--
-- NOTA: o valor 'OWNER' foi removido de propósito. Antes significava
-- "dono da oficina" e as suas permissões (pagamento, entrega, caixa)
-- hoje estão repartidas entre MECHANIC (entrega, caixa) e CLIENTE
-- (pagamento). Ver plan-aply/fix-db-1.md secção 1.
-- ============================================================

-- Tipos enumerados
CREATE TYPE perfil_usuario AS ENUM ('ADMIN', 'MECHANIC', 'CLIENTE');
CREATE TYPE status_agendamento AS ENUM ('AGENDADO', 'CANCELADO', 'NAO_COMPARECEU');
CREATE TYPE status_atendimento AS ENUM ('EM_AVALIACAO', 'AGUARDANDO_APROVACAO', 'EM_EXECUCAO', 'ENCERRADO', 'CANCELADO');
CREATE TYPE tipo_item AS ENUM ('SERVICO', 'PECA');
CREATE TYPE status_aprovacao AS ENUM ('PENDENTE', 'APROVADO', 'RECUSADO');
CREATE TYPE forma_pagamento AS ENUM ('DINHEIRO', 'PIX', 'DEBITO', 'CREDITO');
CREATE TYPE status_pagamento AS ENUM ('CONFIRMADO', 'CANCELADO');

-- ============================================================
-- USUARIO
-- ============================================================
-- DELETE é proibido sobre esta tabela: 6 chaves estrangeiras a
-- referenciam com ON DELETE RESTRICT e o histórico financeiro não
-- pode ser reescrito. Para desligar um utilizador usa-se
-- `ativo = FALSE` + `desativado_em`.
--
-- `perfil = 'ADMIN'` é o único que pode gerir os perfis MECHANIC.
CREATE TABLE usuarios (
    id                       UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    login                    VARCHAR(60)  NOT NULL UNIQUE,
    senha                    VARCHAR(255) NOT NULL,
    perfil                   perfil_usuario NOT NULL,
    nome                     VARCHAR(120) NOT NULL DEFAULT '',
    ativo                    BOOLEAN      NOT NULL DEFAULT TRUE,
    ultimo_login_em          TIMESTAMPTZ,
    desativado_em            TIMESTAMPTZ,
    desativado_por_usuario_id UUID         REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at               TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at               TIMESTAMPTZ  NOT NULL DEFAULT NOW(),

    -- ADMIN e MECHANIC têm de ter nome preenchido: aparecem em
    -- listas do painel. CLIENTE pode ficar vazio, porque o nome
    -- canónico vive em clientes.nome.
    CONSTRAINT usuarios_staff_exige_nome CHECK (perfil = 'CLIENTE' OR btrim(nome) <> ''),

    -- Um utilizador desativado tem sempre data de desativação, e
    -- vice-versa. Impede o estado "meio desativado" que ninguém
    -- consegue depois reconciliar.
    CONSTRAINT usuarios_ativo_coerente CHECK (
        (ativo AND desativado_em IS NULL)
        OR (NOT ativo AND desativado_em IS NOT NULL)
    )
);

CREATE INDEX idx_usuarios_perfil_ativo ON usuarios (perfil, ativo);
CREATE INDEX idx_usuarios_ativo       ON usuarios (ativo);

-- ============================================================
-- CLIENTE
-- ============================================================
-- `usuario_id` é a ligação opcional ao login do dono do carro.
-- NULL = cliente de balcão, ainda sem portal. Quando existir,
-- permite ao CLIENTE pedir agendamento e confirmar o orçamento.
--
-- A coerência `clientes.usuario_id -> usuarios.perfil = 'CLIENTE'`
-- NÃO pode ser garantida aqui: o PostgreSQL proíbe subqueries em
-- constraints CHECK. É invariante da camada de aplicação.
CREATE TABLE clientes (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    usuario_id  UUID UNIQUE REFERENCES usuarios(id) ON DELETE RESTRICT,
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
-- Aprovação do orçamento:
--   decidido_por_usuario_id          -> o CLIENTE que aprovou/recusou
--   decisao_registrada_por_usuario_id -> o MECHANIC/ADMIN que registou
--
-- Uma coluna não reachia para os dois: o cliente confirma no portal,
-- mas é o staff que anota. O CHECK obriga os dois a existirem em
-- simultâneo, para nunca ficar uma decisão sem rastro de quem a
-- tomou nem de quem a escreveu.
--
-- `servico_id` / `peca_id` são preenchidos pela 000002 e são
-- opcionais: item de texto livre continua válido.
CREATE TABLE itens_orcamento (
    id                             UUID               PRIMARY KEY DEFAULT gen_random_uuid(),
    versao_orcamento_id            UUID               NOT NULL REFERENCES versoes_orcamento(id) ON DELETE RESTRICT,
    tipo                           tipo_item          NOT NULL,
    descricao                      VARCHAR(160)       NOT NULL,
    valor                          NUMERIC(10,2)      NOT NULL CHECK (valor >= 0),
    status_aprovacao               status_aprovacao   NOT NULL DEFAULT 'PENDENTE',
    decidido_por_usuario_id        UUID               REFERENCES usuarios(id) ON DELETE RESTRICT,
    decidido_em                    TIMESTAMPTZ,
    decisao_registrada_por_usuario_id UUID            REFERENCES usuarios(id) ON DELETE RESTRICT,
    created_at                     TIMESTAMPTZ        NOT NULL DEFAULT NOW(),
    updated_at                     TIMESTAMPTZ        NOT NULL DEFAULT NOW(),

    CHECK (
        (status_aprovacao = 'PENDENTE'
            AND decidido_por_usuario_id IS NULL
            AND decidido_em IS NULL
            AND decisao_registrada_por_usuario_id IS NULL)
        OR
        (status_aprovacao <> 'PENDENTE'
            AND decidido_por_usuario_id IS NOT NULL
            AND decidido_em IS NOT NULL
            AND decisao_registrada_por_usuario_id IS NOT NULL)
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

    -- Antes: CHECK (tipo = 'PECA' OR custo_oficina IS NULL)
    -- Proibia registar custo de mão-de-obra, o que deixava o
    -- painel de margem sem dados para serviços. Agora o custo é
    -- obrigatório em peças (sem ele não há margem calculável) e
    -- opcional em serviços (muitos são executados sem custo
    -- adicional de material).
    CHECK (tipo <> 'PECA' OR custo_oficina IS NOT NULL)
);

CREATE INDEX idx_itens_executados_atendimento ON itens_executados (atendimento_id);

-- ============================================================
-- PAGAMENTO
-- ============================================================
-- SEMÂNTICA ALTERADA: `usuario_id` deixou de ser "funcionário que
-- recebeu o dinheiro" e passou a ser "CLIENTE que pagou". O
-- pagamento é registado de forma simulada (sem gateway), mas quem
-- paga é o dono do carro.
--
-- `idempotency_key` é o que impede pagamento duplicado quando o
-- mechanic regrista e o cliente reenvia o formulário.
--
-- `cancelado_por_usuario_id` fecha um buraco de auditoria:
-- `status` podia ir a CANCELADO mas não havia forma de saber quem
-- o cancelou.
CREATE TABLE pagamentos (
    id                        UUID           PRIMARY KEY DEFAULT gen_random_uuid(),
    atendimento_id            UUID           NOT NULL REFERENCES atendimentos(id) ON DELETE RESTRICT,
    valor                     NUMERIC(10,2)  NOT NULL CHECK (valor > 0),
    forma_pagamento           forma_pagamento NOT NULL,
    data                      TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    usuario_id                UUID           NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,
    status                    status_pagamento   NOT NULL DEFAULT 'CONFIRMADO',
    cancelado_por_usuario_id  UUID           REFERENCES usuarios(id) ON DELETE RESTRICT,
    idempotency_key           VARCHAR(64)    UNIQUE,
    created_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),
    updated_at                TIMESTAMPTZ    NOT NULL DEFAULT NOW(),

    CHECK (
        (status = 'CONFIRMADO' AND cancelado_por_usuario_id IS NULL)
        OR
        (status = 'CANCELADO' AND cancelado_por_usuario_id IS NOT NULL)
    )
);

CREATE INDEX idx_pagamentos_atendimento   ON pagamentos (atendimento_id);
CREATE INDEX idx_pagamentos_data          ON pagamentos (data);
CREATE INDEX idx_pagamentos_idempotency   ON pagamentos (idempotency_key);

COMMENT ON COLUMN pagamentos.usuario_id IS
    'CLIENTE que efetuou o pagamento (perfil = CLIENTE). Registo simulado, sem integracao a gateway.';
COMMENT ON COLUMN pagamentos.cancelado_por_usuario_id IS
    'Staff que cancelou o pagamento. Obrigatorio quando status = CANCELADO.';

-- ============================================================
-- ATENDIMENTO — entrega
-- ============================================================
-- `entrega_excecao_usuario_id` passou de "autorizada pelo Dono"
-- para "autorizada pelo MECHANIC". O campo não muda de nome nem
-- de tipo; só muda quem o preenche.
COMMENT ON COLUMN atendimentos.entrega_excecao_usuario_id IS
    'MECHANIC que autorizou a entrega com saldo pendente. A entrega e o caixa sao do MECHANIC.';
