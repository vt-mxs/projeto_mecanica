# ModDB — o que mudou na base de dados e como aplicar

Ficheiro de entrada. O detalle de cada alteração está em
[`fix-db-1.md`](./fix-db-1.md).

---

## 1. Resumo em 5 linhas

O `OWNER` deixou de existir. Antes significava "dono da oficina" e
misturava duas coisas muito diferentes — gerir o sistema **e** operar a
oficina. Agora há três papéis separados e essa mistura acabou:

| Papel | Quem é | O que faz |
|---|---|---|
| `ADMIN` | Gestor do sistema | Gere os perfis dos `MECHANIC` e as suas contas. **Não** opera a oficina. |
| `MECHANIC` | Staff de balcão | Atende, executa, faz a entrega e o caixa. |
| `CLIENTE` | Dono do carro | Faz login, pede agendamento, pede atendimento, confirma/recusa o orçamento. |

O `CLIENTE` entra com login pela primeira vez. Antes o dono do carro era
só uma linha em `clientes` (nome + telefone) e não tinha conta nenhuma.

O `ADMIN` **não** apaga mecânicos. Tinha pedido isso, mas o schema
proíbe e o motivo está em [`fix-db-1.md` secção 2](./fix-db-1.md#2-por-que-o-delete-de-mecânicos-não-funciona).
Desliga-se com `ativo = FALSE`.

O pagamento passa a ser do `CLIENTE` (é o dono do carro que paga) e é
registado de forma simulada, sem gateway. Entrega e caixa ficam com o
`MECHANIC`.

---

## 2. Ficheiros tocados

```
backend/db/migrations/000001_init_schema.up.sql      EDITADO
backend/db/migrations/000001_init_schema.down.sql     EDITADO
backend/db/migrations/000002_catalogo_auditoria.up.sql   NOVO
backend/db/migrations/000002_catalogo_auditoria.down.sql NOVO
```

Nada mais foi alterado. Os contratos de API e os docs de modelo
ficaram **intactos**, conforme combinado — mas isso deixa os docs
desalinhados do schema. A lista do que ficou por alinhar está em
[`fix-db-1.md` secção 8](./fix-db-1.md#8-o-que-ficou-desalinhado-a-corrigir-depois).

---

## 3. ⚠️ O que vai dar erro no Docker

### 3.1 Editar a `000001` não chega — o migrate não vai re-correr

Este é o ponto principal e a razão de o comando parecer "não fazer
nada". O `golang-migrate` guarda a versão aplicada na tabela
`schema_migrations`. Num postgres que já correu a `000001`, a tabela
diz `version = 1`. Quando encontra a `000002`, ele diz "a 1 já está aplicada,
não há nada a fazer" — **e não re-executa a `000001`**, por mais que tenhas
mudado o ficheiro.

Ou seja: editar a `000001` num postgres já inicializado é um no-op
silencioso. Não dá erro — dá *nada*, que é pior.

### 3.2 Estado actual desta máquina

Durante a validação destas migrations eu levantei o postgres e corri
`up` + `down`. O teu postgres local **já está na versão 2** com as 14
tabelas novas. Não precisas de fazer mais nada.

```bash
docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  -c "select version from schema_migrations"'
```

### 3.3 Para as máquinas dos colegas (que já tinham a `000001`)

Escolhe uma das duas. As duas **destroem os dados** — que são só de
teste, cada um tem o seu volume.

#### Opção A — Recomeçar de zero (recomendado)

```bash
docker compose down -v
docker compose up -d
```

`down -v` apaga o volume `postgres_data` inteiro. É o caminho mais
limpo: o `000001` novo aplica-se como se fosse a primeira vez.

#### Opção B — Reverter só a `000001` e reaplicar

Se não quiseres mexer no volume:

```bash
set -a; source .env; set +a

# 1. reverter a 000001 (isso faz DROP TABLE — os dados vão-se)
docker compose run --rm migrate \
  -path=/migrations \
  -database="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/$POSTGRES_DB?sslmode=disable" \
  down 1

# 2. reaplicar tudo
docker compose up -d
```

Repara que o `down 1` faz `DROP TABLE`. Portanto a Opção B **também**
perde os dados — a diferença é só não deitar o volume fora. Se tens
dados a preservar, nenhuma das duas serve.

### 3.4 O erro que vais ver se ignorares o 3.1

Se correres apenas `docker compose up -d` num postgres já parado na
versão 1, o migrate não complain. Quando reparares e tentares o
`000002` à mão, ou o `up`, aparece:

```
error: Dirty database version 1. Fix and force version.
```

Saídas, ambas explicadas em `docs/config-DB.md`:

```bash
# (a) se tens backup
pg_restore ... < backups/antes.dump

# (b) marcar a 1 como a última boa e seguir em frente
docker compose run --rm migrate \
  -path=/migrations \
  -database="postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@postgres:5432/$POSTGRES_DB?sslmode=disable" \
  force 1
```

### 3.5 ⚠️ Não uses `force 0`

Parece a solução para o 3.1 mas **não é**. O `force 0` diz ao migrate
"a versão 0 é a actual" sem correr nada. As tabelas antigas continuam
lá, e o `up` seguinte rebenta:

```
error: relation "usuarios" already exists
```

### 3.6 Ordem de serviços

O `compose.yaml` já garante `postgres healthy` → `migrate` → `backend`.
Como a `000002` é nova, não muda nada aí. Se o backend subir e falhar a
ligar, é porque o migrate correu mal — ver `docker compose logs migrate`.

---

## 4. Verificar que ficou bem

```bash
# versão aplicada (deve dizer 2)
docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  -c "select version from schema_migrations"'

# tabelas (deve dizer 14)
docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  -c "select count(*) from information_schema.tables where table_schema='"'"'public'"'"'"'

# os três papéis (deve devolver ADMIN, MECHANIC, CLIENTE)
docker compose exec -T postgres sh -c 'psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" \
  -c "select unnest(enum_range(null::perfil_usuario))"'
```

> O `sh -c` deixa o `psql` ler `$POSTGRES_USER` e `$POSTGRES_DB` do `env`
> do próprio container, em vez de escrever o utilizador à mão no comando.
> Ver `docs/config-DB.md`, secção "Comandos úteis".

---

## 5. O que NÃO está feito (de propósito)

- **Nenhum endpoint novo.** Não há `POST /api/v1/users`, nem painel de
  logs, nem portal do cliente. Isto é só schema.
- **Nenhum DTO, handler, service ou repository.** Continua tudo em
  `b.go` com 7 linhas.
- **`models/enums.go` ainda tem `PerfilOWNER` e `PerfilADMIN` errado.**
  Vai rebentar em runtime quando o repositório tentar mapear
  `ADMIN`. É o próximo passo.
- **Contratos não foram tocados.** O §25 do contrato ainda diz que o
  `OWNER` regista pagamentos, o que agora é falso.

Nada disto foi feito porque o pedido era mapeamento + schema.