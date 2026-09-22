# **Briefing — Contrato da API REST (Nomenclatura em Inglês)** 

## **Sistema de Gestão da Oficina** 

**Nota:** este documento é uma cópia do briefing original do contrato da API, com os nomes de endpoints, campos de request/response e valores de enum traduzidos para inglês, conforme o documento de Nomenclatura Técnica. Todo o texto explicativo permanece em português. Base de rotas: `/api/v1` . 

## **1. Objetivo** 

Definir o contrato da API REST que será utilizado pelo frontend e pelo backend do sistema da oficina. 

A API deve representar as operações do negócio e não simplesmente expor as tabelas do banco. 

O fluxo principal do sistema é: 

Cliente → Veículo → Agendamento/Atendimento → Avaliação/Diagnóstico → Orçamento → Aprovação → Execução → Valor final → Pagamento → Entrega 

O contrato deve permitir preservar a rastreabilidade: o que foi orçado → o que mudou → o que foi aprovado/recusado → o que foi executado → o que foi cobrado → o que foi pago. 

## **2. Princípios do contrato** 

### **2.1 A API não deve ser um CRUD cego do banco** 

Não criar endpoints simplesmente porque existe uma tabela. 

Exemplo: Não necessariamente deve existir: 

```
POST /api/v1/estimate-changes
```

se, para o negócio, a alteração fizer parte de uma operação maior de atualização/versionamento do orçamento. 

O endpoint deve representar uma operação que faça sentido para o sistema. 

### **2.2 Frontend não acessa o banco diretamente** 

Fluxo: 

```
Frontend
```

```
 ↓
```

```
HTTP/REST
 ↓
```

```
API
 ↓
```

```
Regras de negócio
```

```
 ↓
```

```
Banco de dados
```

O frontend não deve conhecer a estrutura interna do PostgreSQL. 

### **2.3 Regras de negócio ficam no backend** 

Exemplos: 

- cliente precisa aprovar alterações antes da execução; 

- alteração não aprovada não pode ser executada; 

- cliente não localizado → trabalho adicional aguarda; 

- valor final considera o realizado; 

- saldo considera todos os pagamentos; 

- entrega fica bloqueada com saldo pendente; 

- exceção de entrega precisa ser registrada; 

- Mecânico não pode executar operações reservadas ao Dono. 

O frontend pode impedir ações visualmente, mas a API deve validar novamente. 

## **3. Autenticação** 

A API terá usuários individuais. 

#### Perfis: 

- `OWNER` 

- `MECHANIC` 

### **Login** 

```
POST /api/v1/auth/login
```

Request: 

```
{
```

```
  "login": "owner",
```

```
  "password": "********"
}
```

Response de sucesso: 

```
{
```

```
  "accessToken": "...",
```

```
  "user": {
```

```
    "id": 1,
```

```
    "login": "owner",
```

```
    "role": "OWNER"
```

```
  }
```

```
}
```

Todas as operações protegidas devem exigir autenticação. 

## **4. Usuários** 

### **Consultar usuário autenticado** 

```
GET /api/v1/me
```

Retorna as informações do usuário atualmente autenticado. 

Não é necessário criar uma API administrativa completa de usuários no MVP. A criação/configuração dos usuários pode ser tratada separadamente conforme a decisão de implementação. 

## **5. Clientes (Customers)** 

### **Criar cliente** 

```
POST /api/v1/customers
```

Request: 

```
{
```

```
  "name": "João da Silva",
```

```
  "phone": "21999999999"
}
```

### **Consultar clientes** 

```
GET /api/v1/customers
```

Deve permitir localizar clientes existentes. 

Exemplo: 

```
GET /api/v1/customers?search=João
```

### **Consultar cliente** 

```
GET /api/v1/customers/{id}
```

### **Alterar cliente** 

```
PATCH /api/v1/customers/{id}
```

Somente dados que realmente possam ser alterados. 

## **6. Veículos (Vehicles)** 

### **Criar veículo** 

```
POST /api/v1/vehicles
```

Exemplo: 

```
{
```

```
  "model": "Gol",
```

```
  "color": "Silver",
```

```
  "licensePlate": "ABC1D23",
```

```
  "year": 2018,
```

```
  "mileage": 85000
```

```
}
```

Os campos que não são obrigatórios devem continuar opcionais conforme as regras do domínio. 

### **Consultar veículos** 

```
GET /api/v1/vehicles
```

Deve ser possível localizar veículos existentes. 

### **Consultar veículo** 

```
GET /api/v1/vehicles/{id}
```

## **7. Agendamentos (Appointments)** 

### **Criar agendamento** 

```
POST /api/v1/appointments
```

Request: 

```
{
```

```
  "customerId": 1,
```

```
  "vehicleId": 5,
```

```
  "scheduledAt": "2026-09-20T10:00:00"
```

```
}
```

### **Consultar agenda** 

```
GET /api/v1/appointments
```

Possíveis filtros: 

- `date` 

- `startDate` 

- `endDate` 

- `status` 

- `customer` 

- `vehicle` 

Os filtros exatos podem ser definidos durante a implementação. 

### **Alterar agendamento** 

```
PATCH /api/v1/appointments/{id}
```

### **Cancelar agendamento** 

```
POST /api/v1/appointments/{id}/cancel
```

O cancelamento não deve apagar o registro. 

### **Registrar não comparecimento** 

```
POST /api/v1/appointments/{id}/no-show
```

O status passa para: 

```
NO_SHOW
```

Não criar um atendimento quando o cliente simplesmente não compareceu. 

## **8. Atendimento (Service Order)** 

O atendimento é a entidade central do processo. 

Um atendimento pode: 

- nascer de um agendamento; 

- ser criado sem agendamento. 

### **Criar atendimento sem agendamento** 

```
POST /api/v1/service-orders
```

Request: 

```
{
```

```
  "customerId": 1,
```

```
  "vehicleId": 5
```

```
}
```

### **Iniciar atendimento a partir de agendamento** 

```
POST /api/v1/appointments/{id}/start-service-order
```

A API deve: 

1. verificar o agendamento; 

2. verificar seu estado; 

3. criar/vincular o atendimento; 

4. reutilizar cliente e veículo; 

5. estabelecer a relação entre agendamento e atendimento. 

### **Consultar atendimento** 

```
GET /api/v1/service-orders/{id}
```

A resposta deve permitir ao frontend obter o contexto necessário para a tela do atendimento. Idealmente: 

```
{
```

```
  "id": 123,
```

```
  "customer": {},
  "vehicle": {},
  "appointment": {},
  "assessment": {},
  "estimate": {},
  "execution": {},
  "payments": [],
  "balance": 850.00,
  "delivery": null
```

```
}
```

A estrutura final pode ser ajustada pelo backend. 

### **Listar atendimentos** 

```
GET /api/v1/service-orders
```

Filtros esperados: 

- em andamento; 

- aguardando aprovação; 

- aguardando execução; 

- aguardando pagamento; 

- finalizados. 

Os status exatos precisam ser definidos no domínio antes da implementação. 

## **9. Avaliação e diagnóstico (Assessment / Diagnosis)** 

Avaliação e diagnóstico não precisam necessariamente ser recursos independentes da API. Eles fazem parte do contexto do atendimento. 

### **Registrar avaliação** 

```
PATCH /api/v1/service-orders/{id}/assessment
```

#### Exemplo: 

```
{
```

   - `"reportedProblem": "Carro falhando ao acelerar",` 

- `"assessment": "Falha percebida durante teste" }` 

### **Registrar diagnóstico** 

```
PATCH /api/v1/service-orders/{id}/diagnosis
```

Exemplo: 

```
{
  "diagnosis": "Problema identificado no sistema de ignição"
```

```
}
```

O atendimento pode permanecer sem diagnóstico concluído enquanto houver investigação. 

## **10. Orçamento (Estimate)** 

O orçamento é uma das partes mais importantes da API. A API precisa preservar: `Orçamento original ↓ Alterações ↓ Novas versões ↓ Decisões do cliente` Nunca simplesmente substituir o orçamento anterior e perder o histórico. 

### **Criar orçamento** 

```
POST /api/v1/service-orders/{id}/estimate
```

Exemplo: `{ "items": [ { "type": "SERVICE", "description": "Troca de correia", "quantity": 1, "unitPrice": 300.00 }, { "type": "PART", "description": "Correia dentada", "quantity": 1, "unitPrice": 250.00 } ] }` A API calcula o total. Não confiar no total enviado pelo frontend. 

## **11. Consultar orçamento** 

```
GET /api/v1/service-orders/{id}/estimate
```

A resposta deve permitir visualizar: 

- versão atual; 

- itens; 

- valores; 

- status de aprovação; 

- histórico; 

- alterações; 

- decisões. 

## **12. Aprovação do orçamento** 

A decisão é do cliente. O usuário da oficina apenas registra essa decisão no sistema. 

Portanto, não tratar o usuário da oficina como "cliente aprovador". 

### **Registrar decisão** 

Endpoint sugerido: 

```
POST /api/v1/estimates/{id}/decisions
```

Exemplo: 

```
{
```

```
  "itemId": 10,
```

```
  "decision": "APPROVED",
```

```
  "customerId": 1
```

```
}
```

O backend deve obter automaticamente: 

- usuário que registrou ( `recordedBy` ); 

- data/hora ( `createdAt` ). 

Não confiar no frontend para enviar o usuário que realizou o registro. 

## **13. Aprovação parcial** 

Cada item pode possuir decisão própria. 

#### Exemplo: 

```
Troca de correia   APPROVED
Troca de filtro     REJECTED
Limpeza             PENDING
```

A API deve permitir representar essa situação. 

## **14. Alteração do orçamento (Estimate Change)** 

Durante a execução pode surgir: 

- peça diferente; 

- alteração de preço; 

- serviço adicional; 

- necessidade descoberta após desmontagem. 

A alteração precisa preservar o valor anterior. 

### **Registrar alteração** 

Endpoint sugerido: 

```
POST /api/v1/estimates/{id}/changes
```

Exemplo: 

```
{
```

```
  "itemId": 10,
```

```
  "previousValue": 250.00,
```

```
  "newValue": 320.00,
```

```
  "reason": "Peça disponível no fornecedor possui preço diferente"
}
```

O backend deve registrar automaticamente: 

- usuário ( `recordedBy` ); 

- data/hora ( `createdAt` ); 

- versão anterior ( `previousVersion` ); 

- nova versão ( `newVersion` ); 

- relação com o orçamento ( `estimateId` ). 

Os valores enviados pelo frontend devem ser validados contra o estado atual do orçamento. 

## **15. Regra de aprovação de alteração** 

Uma alteração relevante precisa passar novamente pela decisão do cliente. 

Fluxo: 

```
Alteração identificada
 ↓
```

`Nova versão do orçamento ↓ Pendente ↓ Cliente aprova?` ↙ ↘ 

```
SIM       NÃO
```

```
 ↓         ↓
```

```
Executa   Não executa
```

Se o cliente não puder ser localizado: 

```
Pendente
```

```
 ↓
```

```
Aguardando decisão
```

```
 ↓
```

```
Não executar alteração
```

## **16. Execução (Execution)** 

A execução representa o que realmente aconteceu na oficina. Não deve simplesmente alterar o orçamento original. 

### **Iniciar execução** 

```
POST /api/v1/service-orders/{id}/execution
```

### **Registrar serviço realizado** 

```
POST /api/v1/executions/{id}/performed-services
```

Exemplo: 

```
{
```

```
  "estimateItemId": 10,
```

```
  "description": "Troca de correia",
```

```
  "actualValue": 300.00
```

```
}
```

### **Registrar peça utilizada** 

```
POST /api/v1/executions/{id}/used-parts
```

Exemplo: 

```
{
```

```
  "estimateItemId": 11,
```

```
  "name": "Correia dentada",
```

```
  "workshopCost": 180.00,
```

```
  "customerPrice": 250.00
```

```
}
```

O custo da oficina ( `workshopCost` ) e o valor cobrado ( `customerPrice` ) são informações diferentes. 

## **17. Relação orçamento × realizado** 

A API deve permitir identificar: 

```
ITEM ORÇADO (estimateItem)
```

```
 ↓
```

```
ITEM REALIZADO (performedService / usedPart)
```

ou: 

```
ITEM ORÇADO (estimateItem)
```

```
 ↓
```

```
NÃO REALIZADO
```

Isso permite comparar planejamento e execução. O item orçado não deve ser sobrescrito para virar o item realizado. 

## **18. Interrupção** 

O MVP possui apenas: 

```
INTERRUPTED
```

com uma observação livre. 

Endpoint sugerido: 

```
POST /api/v1/executions/{id}/interrupt
```

Request: 

```
{
  "observation": "Aguardando autorização do cliente."
```

```
}
```

Não criar categorias detalhadas de interrupção neste momento. 

## **19. Valor final (Final Amount)** 

O valor final deve ser calculado com base no que realmente foi executado/utilizado. 

```
Serviços realizados
```

```
+
```

```
Peças utilizadas
```

```
=
```

```
Valor final
```

Não utilizar simplesmente o total do orçamento original. 

### **Consultar valor final** 

```
GET /api/v1/service-orders/{id}/final-amount
```

O backend deve ser a fonte da verdade para esse cálculo. 

## **20. Pagamentos (Payments)** 

Um atendimento pode possuir vários pagamentos. 

Exemplo: 

```
Pagamento 1: R$ 150 — PIX
Pagamento 2: R$ 300 — DINHEIRO
Pagamento 3: R$ 400 — CRÉDITO
```

### **Registrar pagamento** 

```
POST /api/v1/service-orders/{id}/payments
```

Request: 

```
{
  "amount": 150.00,
```

- `"paymentMethod": "PIX",` 

```
  "paidAt": "2026-09-20"
}
```

Formas de pagamento ( `paymentMethod` ): 

```
CASH
```

```
PIX
```

```
DEBIT_CARD
```

```
CREDIT_CARD
```

O usuário autenticado deve ser identificado pelo backend. 

## **21. Consultar pagamentos** 

```
GET /api/v1/service-orders/{id}/payments
```

Retornar todos os pagamentos vinculados ao atendimento. 

Response: 

```
{
  "data": [
    {
      "id": 200,
      "amount": 150.00,
      "paymentMethod": "PIX",
      "paidAt": "2026-09-20"
    }
  ]
```

```
}
```

## **22. Saldo (Balance)** 

Regra: 

```
BALANCE =
FINAL AMOUNT
```

```
SUM OF ALL PAYMENTS
```

Exemplo: 

```
Valor final: R$ 1.000
Pagamento 1: R$ 150
Pagamento 2: R$ 200
Total pago: R$ 350
Saldo: R$ 650
```

O saldo não deve depender de cálculo realizado pelo frontend. 

## **23. Entrega (Delivery)** 

### **Registrar entrega** 

`POST /api/v1/service-orders/{id}/delivery` Antes de permitir a entrega, o backend deve verificar o saldo. `Saldo = 0` 

```
Entrega permitida.
```

```
Saldo > 0
```

```
Entrega bloqueada, salvo quando for registrada a exceção prevista pelo negócio.
```

### **Entrega com exceção** 

Quando houver exceção: 

```
{
  "balanceException": true,
  "observation": "Cliente retirou o veículo com pagamento pendente."
}
```

A exceção deve ficar registrada. 

## **24. Caixa diário (Daily Cash)** 

O caixa diário é derivado dos pagamentos. Não criar necessariamente uma entidade "Caixa" apenas para armazenar esses totais. 

### **Consultar caixa** 

```
GET /api/v1/cash
```

Filtros: 

• `date` 

Exemplo: 

```
GET /api/v1/cash?date=2026-09-20
```

Exemplo de resposta: `{ "date": "2026-09-20", "totals": { "cash": 500.00, "pix": 850.00, "debitCard": 200.00, "creditCard": 300.00 }, "total": 1850.00, "payments": [] }` 

Os totais devem ser derivados dos pagamentos registrados. 

## **25. Autorização** 

A API deve validar o perfil do usuário. 

#### **Dono (** **`OWNER` )** 

Pode: 

- clientes; 

- veículos; 

- agendamentos; 

- atendimentos; 

- avaliação; 

- diagnóstico; 

- orçamento; 

- aprovação/registro de decisão; 

- execução; 

- valor final; 

- pagamentos; 

- entrega; 

- caixa. 

#### **Mecânico (** **`MECHANIC` )** 

#### Pode: 

- clientes; 

- veículos; 

- agendamentos; 

- atendimentos; 

- avaliação; 

- diagnóstico; 

- execução; 

- registro das informações permitidas. 

Não pode: 

- registrar pagamento; 

- liberar entrega; 

- executar operações reservadas ao Dono; 

- exercer funções administrativas reservadas ao Dono. 

As permissões devem ser implementadas no backend, não apenas escondidas no frontend. 

## **26. Auditoria** 

Operações relevantes devem registrar: 

- usuário; 

- data/hora; 

- operação; 

- registro afetado. 

Especialmente: 

- aprovação/recusa; 

- alteração de orçamento; 

- alteração de valor; 

- execução; 

- pagamento; 

- entrega; 

- exceção; 

- operações administrativas relevantes. 

## **27. Erros HTTP** 

A API deve utilizar códigos HTTP coerentes. 

Exemplos: 

```
200 OK
```

```
201 Created
```

- `204 No Content` 

- `400 Bad Request` 

```
401 Unauthorized
```

```
403 Forbidden
```

```
404 Not Found
```

```
409 Conflict
```

```
422 Unprocessable Entity
```

```
500 Internal Server Error
```

Exemplos de situações: 

**401** — Usuário não autenticado. **403** — Usuário autenticado, mas sem permissão. 

**404** — Recurso inexistente. **409** — Conflito de estado, por exemplo, tentativa de criar agendamento incompatível com uma regra de conflito. 

**422** — Dados semanticamente inválidos para a operação. 

## **28. Formato padrão de erro** 

Sugestão: 

```
{
  "error": {
    "code": "ESTIMATE_NOT_APPROVED",
    "message": "The estimate item has not been approved."
  }
}
```

Os códigos definitivos devem ser padronizados durante a implementação. 

## **29. Paginação e filtros** 

Listagens potencialmente grandes devem aceitar paginação. 

Exemplo: 

```
GET /api/v1/customers?page=1&limit=20
```

A resposta pode conter: `{ "data": [], "pagination": { "page": 1, "limit": 20, "total": 120 }` 

```
}
```

A estrutura final pode ser definida pelo backend. 

## **30. Concorrência** 

Como Dono e Mecânico utilizam a mesma agenda, a API deve tratar conflitos de alteração. 

Exemplo: 

```
Dono consulta horário 10:00
```

```
 +
```

```
Mecânico consulta horário 10:00
```

```
 ↓
```

```
Ambos tentam criar/alterar
```

```
 ↓
```

```
Backend valida novamente
```

```
 ↓
```

```
Somente operação válida é aceita
```

Não confiar somente no frontend para impedir conflitos. 

## **31. Idempotência** 

Operações sensíveis, principalmente pagamentos, devem considerar o risco de uma mesma requisição ser enviada duas vezes. 

Exemplo: 

```
Usuário registra pagamento de R$500
```

```
 ↓
```

```
requisição demora
```

```
 ↓
```

```
frontend tenta novamente
```

A API não deve criar dois pagamentos acidentalmente. 

A estratégia de idempotência deve ser definida tecnicamente pelo backend. 

## **32. Transações** 

Operações que alteram várias informações relacionadas devem ser atômicas. 

Exemplo: 

```
Registrar alteração do orçamento
```

```
pode envolver:
```

```
alteração
```

```
+
```

```
nova versão
```

```
+
```

```
novo estado dos itens
```

```
+
```

```
registro da decisão pendente
```

```
+
```

```
auditoria
```

Se uma parte falhar, o backend não deve deixar o orçamento em estado inconsistente. 

## **33. Regras que a API obrigatoriamente deve garantir** 

- **R01** — Serviço depende da aprovação do cliente. 

- **R02** — Alteração relevante de serviço/peça/valor exige nova decisão. 

- **R03** — Cliente não localizado → alteração não é executada. 

- **R04** — Recusa após desmontagem → trabalho adicional não executado e veículo reassemblado conforme processo da oficina. 

- **R05** — Valor final considera o realizado. 

- **R06** — Peças e serviços realizados compõem o valor final. 

- **R07** — Aprovação parcial é permitida. 

- **R08** — Pagamentos antecipados são considerados no fechamento. 

- **R09** — Saldo = valor final − todos os pagamentos. 

- **R10** — Um atendimento pode possuir vários pagamentos. 

- **R11** — Alteração do valor final recalcula o saldo. 

- **R12** — Agendamento cancelado não é apagado. 

- **R13** — Não comparecimento possui estado próprio. 

- **R14** — Conflitos de agendamento devem ser impedidos pelo backend. 

- **R15** — Dono e Mecânico utilizam a mesma agenda. 

- **R16** — Itens do orçamento possuem estado próprio de aprovação. 

- **R17** — Alteração de orçamento preserva histórico. 

- **R18** — A decisão do cliente e o usuário que registrou a decisão são informações diferentes. 

- **R19** — Pagamento sempre pertence a um atendimento. 

- **R20** — Entrega com saldo pendente é bloqueada, salvo exceção registrada. 

## **34. Recursos que NÃO devem ser criados no MVP** 

Não implementar endpoints apenas para recursos que ficaram fora do escopo: 

- estoque ( `Inventory` ); 

- catálogo permanente de peças ( `PartCatalog` ); 

- catálogo permanente de serviços ( `ServiceCatalog` ); 

- peças compradas mas não utilizadas; 

- controle detalhado de parcelas ( `Installment` ); 

- estorno/refund ( `Refund` ); 

- múltiplos responsáveis pelo atendimento; 

- lógica completa de motorista/terceiro ( `Driver` ); 

- categorias detalhadas de interrupção; 

- entidade explícita de box ( `Box` ); 

- funcionamento completo offline. 

## **35. Critérios de aceite da API** 

#### **Cenário 1 — Atendimento sem agendamento** 

A API deve permitir: 

```
POST /api/v1/customers
```

```
POST /api/v1/vehicles
```

```
POST /api/v1/service-orders
```

```
PATCH /api/v1/service-orders/{id}/assessment
```

```
PATCH /api/v1/service-orders/{id}/diagnosis
```

```
POST /api/v1/service-orders/{id}/estimate
```

```
POST /api/v1/estimates/{id}/decisions
```

```
POST /api/v1/service-orders/{id}/execution
```

```
POST /api/v1/executions/{id}/performed-services (ou /used-parts)
```

```
GET /api/v1/service-orders/{id}/final-amount
```

```
POST /api/v1/service-orders/{id}/payments
```

```
POST /api/v1/service-orders/{id}/delivery
```

#### **Cenário 2 — Atendimento agendado** 

```
POST /api/v1/appointments
```

```
 ↓
```

```
POST /api/v1/appointments/{id}/start-service-order
```

```
 ↓
```

```
Atendimento criado com cliente + veículo + agendamento
```

#### **Cenário 3 — Aprovação parcial** 

A API deve permitir que diferentes itens possuam diferentes decisões. 

#### **Cenário 4 — Alteração de preço** 

```
Orçamento R$ 500
```

```
 ↓
```

```
Alteração
 ↓
```

```
Nova versão
 ↓
```

```
Cliente decide
 ↓
```

```
Histórico anterior permanece consultável
```

#### **Cenário 5 — Cliente não localizado** 

```
Alteração
```

```
 ↓
```

```
Pendente
```

```
 ↓
```

```
Cliente não localizado
```

```
 ↓
```

```
Não executar alteração
```

#### **Cenário 6 — Pagamento antecipado** 

```
Valor final = R$ 1.000
Pagamento = R$ 150
```

```
Saldo = R$ 850
```

#### **Cenário 7 — Múltiplos pagamentos** 

A API deve somar todos os pagamentos associados ao atendimento. 

#### **Cenário 8 — Alteração do valor final** 

Depois de um pagamento parcial: 

```
Valor anterior
```

```
 → pagamento
```

```
 → novo valor
```

```
 → novo saldo
```

Os pagamentos antigos não podem ser alterados. 

#### **Cenário 9 — Entrega bloqueada** 

```
Saldo > 0
 ↓
```

```
POST /api/v1/service-orders/{id}/delivery
 ↓
```

```
API rejeita
```

#### **Cenário 10 — Entrega com exceção** 

```
Saldo > 0
 ↓
exceção registrada (balanceException)
 ↓
```

```
entrega permitida
 ↓
```

```
exceção preservada no histórico
```

#### **Cenário 11 — Reconstrução do atendimento** 

A API deve permitir consultar informações suficientes para reconstruir: 

```
Orçamento original
 ↓
Alterações
 ↓
Decisões
```

```
 ↓
Execução
```

```
 ↓
```

```
Realizado
```

```
 ↓
Valor final
```

```
 ↓
```

```
Pagamentos
```

```
 ↓
```

```
Saldo
```

```
 ↓
```

```
Entrega
```

## **36. O que o backend deve entregar a partir deste briefing** 

Antes da implementação definitiva, o responsável pelo backend deve produzir: 

#### **1. Lista definitiva de endpoints** 

Com: 

- método HTTP; 

- rota; 

- autenticação; 

- perfil necessário; 

- • request; 

- response; 

- erros possíveis. 

#### **2. OpenAPI/Swagger** 

Documentar o contrato formalmente. 

#### **3. DTOs** 

Separar: 

```
Request DTO
```

```
Response DTO
```

da estrutura interna do banco. 

#### **4. Regras de validação** 

Documentar quais validações acontecem no backend. 

#### **5. Matriz de autorização** 

Exemplo: 

|**Operação**|**Dono (OWNER) **|**Mecânico (MECHANIC)**|
|---|---|---|
|Cliente|✓|✓|
|Veículo|✓|✓|
|Agendamento|✓|✓|
|Atendimento|✓|✓|
|Avaliação|✓|✓|
|Diagnóstico|✓|✓|
|Orçamento|✓|conforme regra|
|Decisão/aprovação|✓|restrito|
|Execução|✓|✓|
|Pagamento|✓|✗|
|Entrega|✓|✗|
|Caixa|✓|✗|



A matriz deve ser validada antes de ser considerada definitiva. 

## **37. Relação com o banco** 

A API não deve ser uma cópia das tabelas. 

O banco representa a persistência: 

```
Customer
Vehicle
Appointment
ServiceOrder
Estimate
EstimateVersion
EstimateItem
EstimateChange
EstimateDecision
Execution
PerformedService
UsedPart
Payment
Delivery
```

```
User
```

A API representa operações do domínio: 

```
/appointments/{id}/start-service-order
/estimates/{id}/decisions
/estimates/{id}/changes
/executions/{id}/interrupt
/service-orders/{id}/payments
```

```
/service-orders/{id}/delivery
```

O backend é responsável por transformar uma operação da API nas alterações necessárias no banco. 

## **38. Entregável final esperado** 

O contrato final deve permitir que frontend e backend trabalhem em paralelo. 

O documento final deve conter: 

```
ENDPOINT
```

```
 ↓
```

```
REQUEST
 ↓
```

```
VALIDAÇÕES
```

```
 ↓
```

```
REGRA DE NEGÓCIO
```

```
 ↓
```

```
ALTERAÇÕES NO BANCO
```

```
 ↓
```

```
RESPONSE
```

```
 ↓
```

```
ERROS
```

```
 ↓
```

##### `PERMISSÕES` 

A principal preocupação não é ter muitos endpoints. É garantir que o contrato consiga representar corretamente o processo real da oficina sem perder: histórico + aprovação + execução + valor final + pagamentos + rastreabilidade. 

