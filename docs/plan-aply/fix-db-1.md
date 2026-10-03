# fix-db-1 — alterações estruturais na base de dados

Complemento de [`modDB.md`](./modDB.md). Cada secção diz **o quê**,
**porquê** e **como foi testado**.

Tudo isto já está aplicado e validado contra o Postgres 18 real.
Os comandos de teste estão no fim, em [secção 9](#9-como-foi-validado).

---

## 1. O enum `perfil_usuario`

```sql
-- ANTES
CREATE TYPE perfil_usuario AS ENUM ('OWNER', 'MECHANIC');
-- DEPOIS
CREATE TYPE perfil_usuario AS ENUM ('ADMIN', 'MECHANIC', 'CLIENTE');
```

`OWNER` era ambiguo e é o que motivou a troca. No código Go já existia
um `PerfilADMIN` (`backend/internal/models/enums.go:8`) que **nunca
existiu na base de dados** — qualquer `INSERT` com esse perfil rebentava.
Agora os dois lados concordam.

### Quem faz o quê

| Papel | Login | Operações |
|---|---|---|
| `ADMIN` | staff | gere perfis de `MECHANIC` e as contas. Não atende, não executa, não entrega, não mexe no caixa. |
| `MECHANIC` | staff | atende, avaliação, orçamento, execução, **entrega**, **caixa**. |
| `CLIENTE` | sim | pede agendamento, pede atendimento, **confirma/recusa o orçamento**, **paga**. |

### Repartição das permissões antigas

O `OWNER` antigo tinha tudo na matriz do contrato (§36): clientes,
veículos, pagamentos, entrega, caixa. Esobrou ao ser substituído:

| Operação | Antes | Agora |
|---|---|---|
| Pagamento | `OWNER` | `CLIENTE` (o dono do carro paga; registo simulado) |
| Entrega | `OWNER` | `MECHANIC` |
| Caixa | `OWNER` | `MECHANIC` |
| Gestão de utilizadores | implícita | `ADMIN` |
| Logs do sistema | inexistente | `ADMIN` (via tabela `auditoria`, secção 6) |

---

## 2. Porquê que o `DELETE` de mecânicos não funciona

Isto foi pedido explicitamente e **não é possível**. Não é uma escolha
de desenho minha, é uma consequência do schema.

O `000001` original tinha **6 chaves estrangeiras** a `usuarios`, todas
com `ON DELETE RESTRICT`. O schema inteiro tem **18 `RESTRICT`**:

```bash
rg -o 'ON DELETE \w+' backend/db/migrations/000001_init_schema.up.sql | sort | uniq -c
#  18 ON DELETE RESTRICT
```

O `RESTRICT` está onde é útil. `pagamentos.usuario_id` é `NOT NULL`:

```sql
usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT
```

Um mecânico que já recebeu um pagamento fica **preso na base de dados
para sempre**. E `CASCADE` não é saída: apagar o mecânico apagava o
histórico do caixa — exactamente o que o painel de logs vai mostrar.

**Verificado na base de dados real:**

```
ERROR:  update or delete on table "usuarios" violates RESTRICT setting
        of foreign key constraint "pagamentos_usuario_id_fkey"
```

### A saída: soft delete

```sql
ativo                     BOOLEAN     NOT NULL DEFAULT TRUE,
desativado_em             TIMESTAMPTZ,
desativado_por_usuario_id UUID        REFERENCES usuarios(id) ON DELETE RESTRICT,
```

Com `CHECK (usuarios_ativo_coerente)` a impedir o estado a meio:

```sql
CONSTRAINT usuarios_ativo_coerente CHECK (
    (ativo AND desativado_em IS NULL)
    OR (NOT ativo AND desativado_em IS NOT NULL)
)
```

`desativado_por_usuario_id` existe para o painel de logs mostrar **quem**
desligou a conta — é o `ADMIN`.

---

## 3. O `CLIENTE` com login

O dono do carro já existia como `clientes` (nome + telefone), sem conta
e sem ligação a `usuarios`. Agora pode ter login.

```sql
ALTER TABLE clientes ADD COLUMN usuario_id UUID UNIQUE REFERENCES usuarios(id) ON DELETE RESTRICT;
```

`NULL` = cliente de balcão, ainda sem portal. Preenchido = tem conta e
pode pedir agendamento e confirmar o orçamento.

O `UNIQUE` garante que uma conta não fica associada a dois clientes.

### Limitação que fica

`clientes.usuario_id → usuarios.perfil = 'CLIENTE'` **não pode ser
garantido pela base de dados**. O PostgreSQL proíbe subqueries em
constraints `CHECK`:

```sql
-- ISTO NÃO COMPILA no Postgres:
CONSTRAINT ... CHECK (EXISTS (SELECT 1 FROM usuarios WHERE id = usuario_id AND perfil = 'CLIENTE'))
```

Fica como invariante da camada de aplicação. O mesmo se aplica a
`decisao_registrada_por_usuario_id` ter de ser `MECHANIC`.

---

## 4. Quem decide vs. quem regista

O `CLIENTE` confirma o orçamento no portal. O `MECHANIC` anota. Uma
coluna não chegava para os dois, e sem isso perdia-se a informação de
quem escreveu a decisão — que é o que o painel de logs quer.

```sql
decido_por_usuario_id           UUID REFERENCES usuarios(id) ON DELETE RESTRICT,  -- CLIENTE
decisao_registrada_por_usuario_id UUID REFERENCES usuarios(id) ON DELETE RESTRICT, -- MECHANIC
```

O `CHECK` foi alargado para que uma decisão **nunca** possa existir sem
os dois:

```sql
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
```

---

## 5. Pagamento: agora é o `CLIENTE` que paga

```sql
usuario_id UUID NOT NULL REFERENCES usuarios(id) ON DELETE RESTRICT,  -- agora = CLIENTE
```

O nome da coluna **não** mudou, de propósito: o contrato da API
(documento que ficou intacto) descreve a entidade `Payment` sem listar
colunas, e renomear criaria drift silencioso. Em vez disso há
autodocumentação em SQL:

```sql
COMMENT ON COLUMN pagamentos.usuario_id IS
    'CLIENTE que efetuou o pagamento (perfil = CLIENTE). Registo simulado, sem integracao a gateway.';
```

Se mais tarde quiseres o nome explícito, é `ALTER TABLE pagamentos RENAME
usuario_id TO pago_por_usuario_id` — e aí sim actualizar o contrato.

### Buraco de auditoria fechado

`status` podia ir a `CANCELADO` sem forma de saber **quem** cancelou.
Este painel de logs precisava disso:

```sql
cancelado_por_usuario_id UUID REFERENCES usuarios(id) ON DELETE RESTRICT,
CONSTRAINT pagamentos_check CHECK (
    (status = 'CONFIRMADO' AND cancelado_por_usuario_id IS NULL)
    OR (status = 'CANCELADO' AND cancelado_por_usuario_id IS NOT NULL)
)
```

---

## 6. Tabela `auditoria` (nova, `000002`)

O contrato pede auditoria na §26. Não havia nada. Agora existe:

```sql
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
```

Decisões que valem registar:

- **Genérica, não colunas dedicadas.** Cobrir só as 8 operações da §26
  com colunas próprias obriga a uma migration sempre que aparece um
  caso novo.
- **`usuario_id` é `NOT NULL`.** Nenhuma acção entra na trilha sem actor
  conhecido. Rotinas automáticas usam um utilizador de sistema com
  perfil `ADMIN` — assim não há excepção ao CHECK.
- **Sem FK para `entidade_id`.** A entidade alvo muda com o `entidade`;
  referenciar dez tabelas aqui seria armadilha sem benefício — o registo
  aponta para o que foi alterado, não o impede de o fazer.
- **Append-only.** Correcção faz-se por entrada nova, nunca por `UPDATE`.

Índices: `(usuario_id, criado_em DESC)`, `(entidade, entidade_id)`,
`(criado_em DESC)`, `(acao)`.

---

## 7. Catálogo de serviços e peças (nova, `000002`)

Não existia **nenhuma** tabela de catálogo. `itens_orcamento.descricao`
era texto livre e `itens_executados.descricao` era uma cópia. Sem lista
de preços não há como montar o painel de margem.

```sql
CREATE TABLE servicos (
    id, codigo UNIQUE, nome, descricao,
    preco_base NUMERIC(10,2) CHECK (preco_base >= 0),
    tempo_estimado_min, ativo, created_at, updated_at
);

CREATE TABLE pecas (
    id, codigo UNIQUE, nome,
    preco_compra NUMERIC(10,2) CHECK (preco_compra >= 0),
    preco_venda  NUMERIC(10,2) CHECK (preco_venda  >= 0),
    ativo, created_at, updated_at
);
```

Ligação ao orçamento, **opcional de propósito** — o mechanic continua a
poder escrever "trocar fusível 5A" sem estar no catálogo:

```sql
CONSTRAINT itens_orcamento_origem_coerente CHECK (
    (servico_id IS NULL AND peca_id IS NULL)
    OR (tipo = 'SERVICO' AND servico_id IS NOT NULL AND peca_id IS NULL)
    OR (tipo = 'PECA'    AND peca_id    IS NOT NULL AND servico_id IS NULL)
)
```

### CHECK do custo, invertido

O original **proibia** custo em serviços:

```sql
-- ANTES: nada de custo em serviços, logo zero dados de margem
CHECK (tipo = 'PECA' OR custo_oficina IS NULL)
-- DEPOIS: custo livre em serviços, obrigatório em peças
CHECK (tipo <> 'PECA' OR custo_oficina IS NOT NULL)
```

A peça passou a exigir `custo_oficina` porque sem ele a margem não é
calculável. O serviço ficou livre porque muitos são executados sem custo
de material adicional.

**Isto é uma mudança de comportamento**, não só de schema: uma peça
executada sem custo passa a ser rejeitada.

---

## 8. O que ficou desalinhado (a corrigir depois)

Nada disto foi tocado, conforme combinado:

| Onde | O que está errado |
|---|---|
| `backend/internal/models/enums.go:6-8` | `PerfilOWNER` não existe mais na DB. `PerfilADMIN` agora existe mas o `ADMIN` está mal — falta `PerfilCLIENTE`. Vai rebentar em runtime. |
| `backend/internal/models/usuario.go` | `UserRepsonse` (typo), **expõe `senha`** no JSON e não tem `Login`. Não tem `nome`/`ativo`. |
| `docs/backend/contrato_api_atualizado.md` §3 | Perfis ainda `OWNER`/`MECHANIC`. |
| idem §4 / `:170` | Diz que não é necessária API de utilizadores. O painel exige-a. |
| idem §25, §36 | Matriz ainda dá Pagamento/Entrega/Caixa ao `OWNER`. |
| idem §3 `:147` | Login response devolve `"role"`; o schema diz `perfil`. |
| idem §26 | Pede auditoria — a tabela agora existe, o endpoint não. |
| `docs/modelo_banco.md:120` | `"perfil": "DONO"` — quarta variante de nome, já nem existe no schema. |
| `docs/modelo_banco_oficina.md:38,50` | Tabela de precedência da mudança do enum. |
| `docs/oficina_modelo_conceitual_erd.html:58` | `string perfil` sem valores. |
| `docs/frontend/doc_front_atualizado.md` | Sem ecrãs de portal do cliente nem painel de admin. |
| `backend/go.mod` | Fiber **v2 e v3** ao mesmo tempo, ambos `// indirect`. Escolher um antes de escrever código. |

---

## 9. Como foi validado

Todas as migrations foram aplicadas contra o Postgres 18 real e cada
constraint foi testada com o insert que a deve fazer falhar.

| # | Teste | Resultado |
|---|---|---|
| 1 | `CLIENTE` sem nome | passa |
| 2 | `MECHANIC` sem nome | falha `usuarios_staff_exige_nome` |
| 3 | `ativo=FALSE` sem `desativado_em` | falha `usuarios_ativo_coerente` |
| 4 | `INSERT` com perfil `OWNER` | falha `invalid input value for enum` |
| 5 | `DELETE` de mecânico com pagamento | falha `RESTRICT setting of FK pagamentos_usuario_id_fkey` |
| 6 | `pagamentos.status='CANCELADO'` sem quem cancelou | falha `pagamentos_check` |
| 7 | mesmo `UPDATE` com quem cancelou | passa |
| 8 | decisão sem `decisao_registrada_por` | falha `itens_orcamento_check` |
| 9 | decisão com os dois campos | passa |
| 10 | `tipo='PECA'` ligado a `servico_id` | falha `itens_orcamento_origem_coerente` |
| 11 | `tipo='PECA'` ligado a `peca_id` | passa |
| 12 | item sem catálogo (texto livre) | passa |
| 13 | peça executada sem `custo_oficina` | falha `itens_executados_check` |
| 14 | `auditoria` sem `entidade_id` | falha `NOT NULL` |

Ciclo de migrations:

```
ANTES    -> schema_migrations=[2] tabelas=14
DOWN 2   -> schema_migrations=[sem_linhas] tabelas=1   (1 = o schema_migrations do próprio migrate)
UP       -> schema_migrations=[2] tabelas=14
DOWN 1   -> schema_migrations=[1] tabelas=11
UP FINAL -> schema_migrations=[2] tabelas=14
```

Repara que depois de `down 2` fica 1 tabela, não 0 — é o
`schema_migrations` do próprio `golang-migrate`, que nunca é dropado.