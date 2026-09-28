# Como inicializar a db localmente

## O que fazer

Com o serviço `migrate` no `compose.yaml`, é um comando só:

```bash
docker compose up -d
```

Isto sobe o postgres, espera que fique pronto, aplica as migrations e só depois
arranca o backend. Não é preciso expor a porta da base de dados.

### Alternativa manual (sem usar o compose)

Use isto apenas se estiveres a depurar. Requer o binário `migrate` instalado
**na tua máquina** e a porta `5432` publicada.

#### 1- Pegar as variaveis que estão no .env:

````bash
set -a; source .env; set +a
````

#### 2- utilizar essas vars para rodar as migrations

````bash
migrate -path ./db/migrations -database "postgres://$POSTGRES_USER:$POSTGRES_PASSWORD@$POSTGRES_HOST:$POSTGRES_PORT/$POSTGRES_DB?sslmode=disable" up
````

> ⚠️ Este comando **não funciona** tal como está. Ver
> [Porque é que o serviço `migrate` existe](#porque-é-que-o-serviço-migrate-existe-no-composeyaml).

---

## Porque é que o serviço `migrate` existe no `compose.yaml`

### O problema

Tentar correr as migrations à mão falhou por quatro motivos ao mesmo tempo.
Vale a pena registar porque todos vão bater noutro se não soubermos porque é que
isto é automático.

**1. O `migrate` não existe dentro do container**

```bash
docker compose exec backend migrate ...
# OCI runtime exec failed: exec: "migrate": executable file not found in $PATH
```

O `go.mod` lista a dependência, mas isso não instala o binário. Só há binário
se o `migrate` estiver na imagem.

**2. O `tool` directive aponta para o pacote errado**

O `go.mod` tem:

```go
tool github.com/golang-migrate/migrate/v4
```

Esse caminho é a **biblioteca** (`package migrate`), não o CLI. O CLI está em
`/cmd/migrate` (`package main`). Como o `tool` directive exige um pacote `main`,
esta linha está errada.

E a sintaxe também: `go run migrate` interpreta `migrate` como nome de pacote e
procura-o no stdlib.

```bash
go run migrate ...
# package migrate is not in std (/usr/local/go/src/migrate)
```

**3. `POSTGRES_HOST=localhost` aponta para o Postgres errado**

O serviço `postgres` do compose **não publica a porta `5432`**. Logo, do host,
`localhost:5432` não chega ao container. E nesta máquina havia já um Postgres
nativo do Fedora a escutar em `127.0.0.1:5432`:

```bash
ss -ltnp | grep 5432
# LISTEN 0 200 127.0.0.1:5432 0.0.0.0:*
```

Ou seja, o `migrate` estava a falar com o **Postgres do host**, não com o do
Docker, e a autenticação falhava porque são bases de dados diferentes.

Confirmou-se depois que a password no container estava correcta — testando pelo
IP interno do container, que cai na regra `scram-sha-256` do `pg_hba.conf`:

```bash
docker compose exec postgres sh -c \
  'PGPASSWORD=$POSTGRES_PASSWORD psql -h $(hostname -i) -U $POSTGRES_USER -d $POSTGRES_DB -c "select current_user"'
# vtgbrc
```

**4. Não há forma de saber que versão do schema está aplicada**

Sem tabela de controlo, ninguém sabe se a base de dados de uma máquina está em
sincronia com o repositório.

### O que o serviço `migrate` resolve

```yaml
migrate:
  image: migrate/migrate:v4.20.1
  volumes:
    - ./backend/db/migrations:/migrations:ro
  command: ["-path=/migrations", "-database=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@postgres:5432/${POSTGRES_DB}?sslmode=disable", "up"]
  depends_on:
    postgres: {condition: service_healthy}
  restart: "no"
```

Pontos importantes:

| Propriedade | Porquê |
|---|---|
| `image: migrate/migrate:v4.20.1` | O CLI vem pronto. Não precisamos de instalar nada no `go.mod` nem na imagem do backend. A versão está fixada para bater certo com a do `go.mod`. |
| `restart: "no"` | É **one-shot**: corre, aplica, e sai. Não fica a fazer loop. |
| `depends_on: service_healthy` | Não tenta ligar a uma base de dados que ainda está a arrancar. |
| `backend: service_completed_successfully` | O backend só arranca **depois** das migrations terem corrido bem. Se falharem, o backend não sobe. |
| `postgres:5432` | Hostname interno do compose. **Não** é preciso publicar a porta. |
| `/migrations:ro` | Os ficheiros vêm do repositório, montados apenas de leitura. |
| URL directo no `command` | A URL é interpolada pelo Compose **no host**, não dentro do container. Ver a nota abaixo. |

> **Porque é que a URL está no `command` e não numa `environment`?**
> A ideia natural seria `command: ["-database=$$DATABASE_URL", "up"]` com a URL
> numa variável de ambiente. Não funciona: a imagem `migrate/migrate` tem
> `ENTRYPOINT` directo ao binário, **sem shell**. Sem shell não há nada para
> expandir `$DATABASE_URL`, e o binário recebe a string literal
> `$DATABASE_URL`:
>
> ```
> error: failed to parse scheme from database URL: no scheme
> ```
>
> O `$$` só serviria se o comando passasse por um shell. Como não passa, a
> interpolamos no host e escrevemos o valor à mão no `command`.

### Migrations como passo separado, não no arranque da app

Esta é a parte que mais importa para produção. A app **não** corre migrations ao
arrancar. Corre-as o `migrate`, que é um passo próprio, one-shot.

Se as migrations corressem no boot da app, em produção teríamos várias réplicas
a subir ao mesmo tempo, todas a tentar migrar a mesma base de dados ao mesmo
tempo. O driver Postgres do golang-migrate usa `pg_advisory_lock` para
serializar isso — o que é seguro — mas mesmo assim o modelo é mais frágil do
que um passo explícito, separado, que corre uma vez e sai.

### Porquê golang-migrate e não `AutoMigrate` do GORM

O `docs/stack.md` já indicava golang-migrate, e a escolha está certa.

| | `AutoMigrate` (GORM) | Ficheiros versionados (golang-migrate) |
|---|---|---|
| Histórico de versões | ❌ | ✅ tabela `schema_migrations` |
| Rollback | ❌ | ✅ ficheiros `down` |
| SQL revisto em code review | ❌ fica implícito nos structs | ✅ é um `.sql` no repositório |
| Data migrations (`UPDATE`, `INSERT`) | ❌ | ✅ |
| Vários programadores / ambientes | dessincroniza | todos aplicam o mesmo |
| Deixar o `DROP COLUMN` para depois | fica esquecido | explícito no migration |

A própria documentação do GORM avisa: *"at some point you will need to switch to
a versioned migrations strategy"*. E o `AutoMigrate` altera tipos de coluna
quando muda o size/precision, ou quando passa de `NOT NULL` para nullable — em
tabelas grandes isso pode bloquear a tabela durante muito tempo.

**Regra:** GORM para as queries (repository), golang-migrate para o schema.
Não se misturam.

### Porque é que o postgres não publica a porta

O backend comunica com a base de dados por `postgres:5432`, dentro da rede do
compose. Não é preciso expor nada ao exterior.

Duas Razões:

- Menos superfície de ataque — a base de dados não fica acessível na rede.
- Nesta máquina, publicar `5432:5432` entraria **em conflito** com o Postgres
  nativo do Fedora, que já ocupa a `127.0.0.1:5432`.

Se precisares de uma ferramenta de GUI (DBeaver, TablePlus) para ligar
directamente à base de dados, usa uma porta diferente, por exemplo `5433`, e
publica só em `localhost`. Não faças disto a configuração por omissão.

---

## Workflow da equipa

O objectivo é: **ninguém precisa de saber nada de migrations para trabalhar.**

```bash
git pull
docker compose up -d
```

Fim. O `compose.yaml` garante a ordem: postgres pronto → migrations aplicadas →
backend arrancado.

Quando alguém cria uma migration nova:

1. Cria `backend/db/migrations/000002_*.up.sql` e `000002_*.down.sql`
2. `git push`
3. Os outros fazem `git pull` + `docker compose up -d` e recebem a base de dados
   actualizada

**O que é partilhado e o que não é:**

| | Onde vive | Partilhado? |
|---|---|---|
| Esquema da base de dados (ficheiros `.sql`) | Repositório | ✅ via git |
| Código do backend / frontend | Repositório | ✅ Todos têm o mesmo |
| `compose.yaml` | Repositório | ✅ Todos têm o mesmo |
| Dados da base de dados | Volume Docker `postgres_data` | ❌ Cada máquina tem os seus |

O volume `postgres_data` **não** é partilhado entre máquinas. Cada pessoa tem a
sua própria base de dados local, com o mesmo esquema e os seus próprios dados
de teste.

Quando aparece uma migration nova, o outro colega recebe o **esquema**
actualizado, mas **não** os teus dados — e isso é o pretendido. Os dados de
produção são um problema separado, resolvido com backup/restore, nunca com
`docker cp` de um volume para o outro.

---

## Fazer backup da base de dados

### Onde guardar os backups

**No teu disco, nunca num volume do Docker.** Isto é importante e não é
arbitrário: o `docker compose down -v` apaga *todos* os volumes do projecto. Um
backup guardado num volume desaparece precisamente quando mais precisas dele.

Uma pasta normal no repositório, já no `.gitignore`, é o sítio certo:

```bash
mkdir -p backups
```

### Backup completo (dados + schema)

```bash
docker compose exec -T postgres pg_dump -U vtgbrc -d mecanica_db -Fc > backups/db_$(date +%Y%m%d_%H%M%S).dump
```

O `-T` é obrigatório: sem ele o `docker exec` aloca uma TTY e o ficheiro fica
contaminado com sequências de controlo. O `-Fc` é o formato comprimido do
Postgres, e o `>` redirecciona para o teu disco porque o `exec` corre no teu
shell, não dentro do container.

Confirma que o ficheiro não está vazio antes de confiares nele:

```bash
ls -lh backups/
pg_restore -l backups/db_20260928_140000.dump | head
```

### Só o schema

```bash
docker compose exec -T postgres pg_dump -U vtgbrc -d mecanica_db --schema-only > backups/schema_$(date +%Y%m%d).sql
```

O schema já vive no git, nos ficheiros de migration. Este dump é uma rede de
segurança extra, útil para veres exactamente o que está na base de dados num
dado momento, independentemente do que o `schema_migrations` diz.

### Restauro

```bash
docker compose exec -T postgres pg_restore -U vtgbrc -d mecanica_db --clean --if-exists < backups/db_20260928_140000.dump
```

O `--clean --if-exists` diz ao `pg_restore` para largar os objectos que já lá
estão antes de recriar. Sem isto, o restauro falha com "already exists".

---

## Reverter migrations (`down`)

### O problema

Fazer `down` **não é desfazer uma operação**. É correr o ficheiro `down.sql`
por baixo. Duas consequências que é preciso ter sempre em cabeça:

1. **Os dados vão-se.** Se a migration adicionou uma coluna, encheu-a, e fazes
   `down` — a coluna e tudo o que lá estava desaparece. Não há "desfazer".
2. **O código da aplicação não volta atrás.** Se fizeste deploy da app v2, que
   já usa a coluna nova, e fazes `down`, tens a app v2 a correr contra o schema
   antigo → a aplicação quebra.

### O backup é o que torna isto seguro

É este o ponto central. **Com um backup feito antes, o `down` deixa de ser
arriscado**, porque passa a ser um experimento reversível:

```
 1.  Backup                    ← fotografa o estado actual
            ↓
 2.  migrate down              ← reverte a migration, pode correr mal
            ↓
 3.  testas à vontade          ← o schema está na alteração, é o que querias ver

 4.  se não gostaste            →  4.  restore do backup     ← volta tudo ao ponto 1
```

O truque que torna isto fiável: **a tabela `schema_migrations` está dentro da
base de dados, portanto está dentro do dump.** Quando fazes o restauro, não
restauras só os dados — restauras também a versão do schema. Ficas
exactamente no ponto de onde partiste, com a aplicação e a base de dados em
sincronia.

Sem backup, o passo 4 não existe. Com backup, o passo 2 é uma aposta
irreversível.

**Regra prática: nunca corras `down` sem um backup feito nos últimos minutos.**

### ⚠️ O `down` sem argumento reverte TUDO

Isto é contra-intuitivo e vale a pena interiorizar antes de usar o `down`.

No `golang-migrate` **v4.20.1**, um `down` sem número **não reverte a última
migration — reverte todas**:

```bash
docker compose run --rm migrate ... down
# Are you sure you want to apply all down migrations? [y/N]
```

O prompt aparece precisamente porque a operação é destrutiva. Se responderes
`y`, todas as tabelas e tipos da base de dados desaparecem.

| Comando | O que faz |
|---|---|
| `down` | Reverte **todas** as migrations. Pede confirmação. |
| `down 1` | Reverte **uma** — a última aplicada. Sem prompt. |
| `down 3` | Reverte as últimas três. Sem prompt. |
| `down -all` | Reverte todas, **sem** prompt. Perigoso. |
| `up` | Aplica tudo o que falta. Idempotente. |

Para quase tudo o que vais fazer, **`down 1` é o que queres**: desfazer a
última migration que acabaste de escrever.

### Correr o `down`

O `docker compose run --rm migrate <args>` substitui o `command` do serviço,
por isso tens de repetir o `-path` e o `-database`:

```bash
# 1. SEMPRE backup antes (é o que torna o down reversível)
docker compose exec -T postgres pg_dump -U vtgbrc -d mecanica_db -Fc > backups/antes_do_down.dump

# 2. reverter a ÚLTIMA migration
docker compose run --rm migrate \
  -path=/migrations \
  -database="postgres://vtgbrc:804512@postgres:5432/mecanica_db?sslmode=disable" \
  down 1

# 3. se correu mal, voltar ao ponto 1
docker compose exec -T postgres pg_restore -U vtgbrc -d mecanica_db --clean --if-exists < backups/antes_do_down.dump
```

Este ciclo foi testado e confirmado: backup → `down 1` → restauro devolve as 10
tabelas, os 7 enums, a versão `1` no `schema_migrations` e os dados intactos.

### Quando o `down` **não** é a resposta

Isto vale mesmo com backup, porque o backup só restaura a **base de dados** —
não o código.

Em produção, se algo correu mal, a prática é **roll forward**, não roll back:

```
Problema no deploy
       ↓
NÃO fazer down
       ↓
Criar a migration 000003 que corrige
       ↓
up        ← roll forward
```

Motivos:

- As down migrations raramente são testadas. Escreveste-as, nunca as correste.
  Quando precisas, há 50% de chance de estarem erradas — e com o `down` a correr
  a meio, o estado fica pior do que antes.
- Em produção não tens o luxo de "esperar pelo restore". Demora, e durante esse
  tempo a aplicação está a cair.

O `down` é uma ferramenta de **desenvolvimento**: experimentar um schema
diferente, desfazer uma migration que acabaste de escrever, ver o que muda.
Com o backup, é seguro. Sem ele, é uma aposta.

### Se uma migration falhar a meio

O `migrate` pode deixar a base de dados marcada como `dirty` e recusa-se a
correr seja o que for:

```bash
docker compose exec postgres psql -U vtgbrc -d mecanica_db -c "select * from schema_migrations"
#  version | dirty
#         2 | t        ← nenhuma migration pode correr
```

Duas saídas:

```bash
# 1) se tens backup: restaura e fica resolvido
docker compose exec -T postgres pg_restore -U vtgbrc -d mecanica_db --clean --if-exists < backups/antes.dump

# 2) se a migration a meio é aceitável, marca como aplicada
docker compose run --rm migrate \
  -path=/migrations \
  -database="postgres://vtgbrc:804512@postgres:5432/mecanica_db?sslmode=disable" \
  force 1
```

O `force 1` diz ao `migrate` "a versão 1 é a última boa, esquece a 2". Cabe-te
a ti confirmar à mão que o estado real do schema bate certo com isso — o
`force` não verifica nada, só marca.

---

## Comandos úteis

```bash
# aplicar as migrations (ou reaplicar, é idempotente)
docker compose up -d
docker compose run --rm migrate

# ver o que aconteceu
docker compose logs migrate

# entrar no psql DENTRO do container (não precisa de porta publicada)
docker compose exec postgres psql -U vtgbrc -d mecanica_db

# ver as tabelas
docker compose exec postgres psql -U vtgbrc -d mecanica_db -c "\dt"

# ver que versão do schema está aplicada
docker compose exec postgres psql -U vtgbrc -d mecanica_db -c "select * from schema_migrations"

# backup e restauro
mkdir -p backups
docker compose exec -T postgres pg_dump -U vtgbrc -d mecanica_db -Fc > backups/db_$(date +%Y%m%d_%H%M%S).dump
docker compose exec -T postgres pg_restore -U vtgbrc -d mecanica_db --clean --if-exists < backups/db_20260928_140000.dump

# recomeçar de zero (APAGA OS DADOS — só em dev)
docker compose down -v && docker compose up -d
```

---

## Notas sobre a imagem do postgres

O `compose.yaml` monta o volume em:

```yaml
volumes:
  - postgres_data:/var/lib/postgresql
```

Isto está **correcto** para a imagem `postgres:18`, mas é uma mudança recente e
vale a pena saber porquê.

A partir do PostgreSQL 18, a imagem oficial mudou o `PGDATA` para
`/var/lib/postgresql/<versão>/docker` e o `VOLUME` declarado passou a ser
`/var/lib/postgresql` (antes era `/var/lib/postgresql/data`).

Se alguma vez mudares para `postgres:17` ou inferior, **tem de ser**
`/var/lib/postgresql/data`. Com o path errado em versões antigas, os dados são
escritos num volume anónimo descartado e **não persistem** quando o container
é recriado — uma perda de dados silenciosa.

