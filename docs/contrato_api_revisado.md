# Contrato da API REST — Sistema de Gestão da Oficina (revisado)

Base de rotas: `/api/v1`. Esta versão ajusta o contrato original ao modelo de
banco revisado: `Estimate`, `Execution` e `EstimateDecision` deixaram de
existir como recursos com id próprio (foram fundidos em `ServiceOrder` /
`EstimateItem`), então as rotas que dependiam desses ids mudaram. As seções
não listadas aqui (customers, vehicles, appointments, autenticação, caixa,
erros, paginação) permanecem exatamente como no documento original.

## 1. O que mudou

| Rota original | Rota revisada | Motivo |
|---|---|---|
| `POST /api/v1/estimates/{id}/decisions` | `POST /api/v1/service-orders/{id}/estimate-items/{itemId}/decision` | `Estimate` não é mais um recurso com id estável (cada versão tem id próprio); a rota passa a ser ancorada no atendimento |
| `POST /api/v1/estimates/{id}/changes` | `POST /api/v1/service-orders/{id}/estimate/changes` | mesmo motivo; o backend resolve a versão vigente internamente |
| `POST /api/v1/service-orders/{id}/execution` | mantida, mas deixa de criar um recurso `Execution` com id | vira só uma ação que marca `execucao_inicio` no atendimento — nada para o frontend guardar como id |
| `POST /api/v1/executions/{id}/performed-services` | `POST /api/v1/service-orders/{id}/execution/performed-services` | `Execution` não tem mais id próprio |
| `POST /api/v1/executions/{id}/used-parts` | `POST /api/v1/service-orders/{id}/execution/used-parts` | idem |
| `POST /api/v1/executions/{id}/interrupt` | `POST /api/v1/service-orders/{id}/execution/interrupt` | idem |

`GET /api/v1/service-orders/{id}/estimate`, `GET /api/v1/service-orders/{id}/final-amount`,
pagamentos, saldo e entrega **não mudam de rota** — já eram ancorados no
atendimento.

## 2. Aprovação do orçamento (seção 12–13 do documento original)

Registrar decisão

```
POST /api/v1/service-orders/{id}/estimate-items/{itemId}/decision
```

Request:
```json
{
  "decision": "APPROVED",
  "customerId": 1
}
```

O backend continua obtendo automaticamente `recordedBy` (usuário logado) e
`createdAt`. `decision` aceita `APPROVED` ou `REJECTED`.

Isso vale tanto para um item do orçamento original quanto para um item que
nasceu de uma alteração (`newItemId` de uma alteração, ver seção 3) — é o
mesmo endpoint, porque no banco a decisão sempre vive no item, nunca na
alteração em si.

## 3. Alteração do orçamento (seção 14 do documento original)

Registrar alteração

```
POST /api/v1/service-orders/{id}/estimate/changes
```

Request:
```json
{
  "itemId": 10,
  "newValue": 320.00,
  "reason": "Peça disponível no fornecedor possui preço diferente"
}
```

O backend:
1. busca a versão vigente do orçamento pelo `service-orders/{id}`;
2. lê `previousValue` do item atual — não confia em valor enviado pelo
   frontend para o "antes";
3. cria uma nova versão, copiando os itens não afetados e criando o item
   novo com o valor alterado;
4. registra a alteração (`estimateItemId` original → `newItemId`,
   `previousValue`, `newValue`, `reason`, `recordedBy`, `createdAt`);
5. o item novo nasce com `status: PENDING` — precisa passar pelo endpoint da
   seção 2 para ser aprovado/recusado.

Response:
```json
{
  "id": 55,
  "estimateItemId": 10,
  "newItemId": 87,
  "previousValue": 250.00,
  "newValue": 320.00,
  "reason": "Peça disponível no fornecedor possui preço diferente",
  "recordedBy": { "id": 2, "login": "mecanico1" },
  "createdAt": "2026-09-20T14:32:00"
}
```

## 4. Execução (seção 16–18 do documento original)

Iniciar execução
```
POST /api/v1/service-orders/{id}/execution
```
Resposta `204 No Content` — apenas marca o início; não há id de execução
para guardar.

Registrar serviço realizado
```
POST /api/v1/service-orders/{id}/execution/performed-services
```
```json
{
  "estimateItemId": 10,
  "description": "Troca de correia",
  "actualValue": 300.00
}
```

Registrar peça utilizada
```
POST /api/v1/service-orders/{id}/execution/used-parts
```
```json
{
  "estimateItemId": 11,
  "name": "Correia dentada",
  "workshopCost": 180.00,
  "customerPrice": 250.00
}
```

Interromper execução
```
POST /api/v1/service-orders/{id}/execution/interrupt
```
```json
{ "observation": "Aguardando autorização do cliente." }
```

## 5. `GET /api/v1/service-orders/{id}` — estrutura de resposta revisada

```json
{
  "id": 123,
  "customer": {},
  "vehicle": {},
  "appointment": {},
  "assessment": {},
  "estimate": {
    "version": 2,
    "items": [
      {
        "id": 87,
        "type": "PART",
        "description": "Correia dentada",
        "value": 320.00,
        "status": "APPROVED",
        "decidedBy": { "id": 2, "login": "mecanico1" },
        "decidedAt": "2026-09-20T15:00:00"
      }
    ],
    "changes": [
      {
        "id": 55,
        "estimateItemId": 10,
        "newItemId": 87,
        "previousValue": 250.00,
        "newValue": 320.00,
        "reason": "Peça disponível no fornecedor possui preço diferente",
        "recordedBy": { "id": 2, "login": "mecanico1" },
        "createdAt": "2026-09-20T14:32:00"
      }
    ]
  },
  "execution": {
    "startedAt": "2026-09-20T13:00:00",
    "endedAt": null,
    "interrupted": false,
    "interruptionObservation": null,
    "performedServices": [],
    "usedParts": []
  },
  "payments": [],
  "balance": 850.00,
  "delivery": null
}
```

Diferenças em relação ao documento original: `estimate.items[]` ganha
`decidedBy`/`decidedAt`; `estimate` ganha `changes[]` no mesmo payload em vez
de precisar de uma chamada separada; `execution` não tem `id` — é um objeto
embutido, montado a partir dos campos do próprio atendimento mais as listas
de itens executados.

## 6. Pagamentos — idempotência (seção 20, 31 do documento original)

```
POST /api/v1/service-orders/{id}/payments
```

O frontend deve enviar uma chave de idempotência:

```json
{
  "amount": 150.00,
  "paymentMethod": "PIX",
  "paidAt": "2026-09-20",
  "idempotencyKey": "c3f1e2a0-51b2-4e9a-9a2e-3d6f9a2b7e10"
}
```

Regra: se `idempotencyKey` já existir para o atendimento, a API **não cria
um novo pagamento** — retorna `200 OK` com o pagamento já registrado em vez
de `201 Created`. O frontend deve gerar essa chave uma única vez por
tentativa de pagamento (ex.: UUID gerado ao abrir o formulário) e reenviá-la
em caso de retry.

## 7. Agendamentos — concorrência (seção 7, 30 do documento original)

`GET /api/v1/appointments/{id}` e a listagem passam a retornar `version`:

```json
{
  "id": 42,
  "customerId": 1,
  "vehicleId": 5,
  "scheduledAt": "2026-09-20T10:00:00",
  "status": "SCHEDULED",
  "version": 3
}
```

`PATCH /api/v1/appointments/{id}` passa a exigir o campo `version` recebido
na leitura:

```json
{
  "scheduledAt": "2026-09-20T11:00:00",
  "version": 3
}
```

Se a `version` enviada não bater com a atual (porque outro usuário alterou
o registro entre a leitura e a escrita), a API responde:

```
409 Conflict
```
```json
{
  "error": {
    "code": "APPOINTMENT_VERSION_MISMATCH",
    "message": "This appointment was modified by another user. Reload and try again."
  }
}
```

## 8. Relação com o banco (seção 37 do documento original, atualizada)

```
Customer
Vehicle
Appointment          (agora com campo de versão para lock otimista)
ServiceOrder
EstimateVersion
EstimateItem          (agora carrega a decisão: status + decidedBy + decidedAt)
EstimateChange         (não carrega mais decisão — só o antes/depois do valor)
PerformedService
UsedPart
Payment                (agora com idempotencyKey)
User                    (role agora é OWNER/MECHANIC, igual ao contrato)
```

`Estimate` e `Execution` deixaram de ser tabelas/recursos próprios — o
briefing original já sugeria isso como possível ("a estrutura exata das
tabelas pode ser alterada"); a API absorve essa fusão sem expor a mudança
ao frontend, exceto pelas rotas listadas na seção 1 deste documento.

## 9. O que não muda

Autenticação, `customers`, `vehicles`, `appointments` (exceto o campo
`version`), listagem/filtros de `service-orders`, `GET .../final-amount`,
`GET .../payments`, `POST .../delivery`, caixa diário, matriz de
autorização, formato de erro, paginação — todos permanecem exatamente como
no documento original.
