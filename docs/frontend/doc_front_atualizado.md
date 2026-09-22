# **Briefing do Frontend (Nomenclatura em Inglês)** 

## **Sistema de Gestão da Oficina** 

**Nota:** cópia do briefing original do frontend, com endpoints, JSONs, tipos/interfaces, nomes de services e estrutura de pastas traduzidos para inglês, seguindo o mesmo dicionário aplicado ao contrato da API. Textos de tela/menu voltados ao usuário final permanecem em português, conforme a regra de nomenclatura ("interface para o usuário pode continuar em português"). Base de rotas da API: `/api/v1` . 

# **PARTE 1 — TELAS E FLUXO DA APLICAÇÃO** 

## **1. Objetivo do frontend** 

O frontend deve permitir que Dono e Mecânico utilizem o sistema seguindo o fluxo real da oficina: `Login` 

```
 ↓
```

```
Dashboard
 ↓
```

```
Clientes / Veículos
```

```
 ↓
```

```
Agenda / Atendimento
```

```
 ↓
```

```
Avaliação / Diagnóstico
 ↓
```

```
Orçamento
 ↓
```

```
Aprovação
 ↓
```

```
Execução
```

```
 ↓
```

```
Valor final
 ↓
```

```
Pagamentos
 ↓
```

```
Entrega
```

```
 ↓
```

```
Caixa
```

A interface deve apresentar o estado atual do atendimento e facilitar a identificação de pendências. 

## **2. Estrutura principal de navegação** 

Sugestão de menu: 

```
Dashboard
Atendimentos
```

```
Agenda
Clientes
Veículos
Caixa
```

```
[Usuário]
```

Não criar menus separados para todos os conceitos internos do banco. 

Por exemplo: 

- "Versão do orçamento" ( `EstimateVersion` ) é parte da tela de orçamento; 

- "Alteração do orçamento" ( `EstimateChange` ) é parte do histórico do orçamento; 

- "Serviço realizado" ( `PerformedService` ) e "Peça utilizada" ( `UsedPart` ) fazem parte da execução; 

- "Saldo" ( `balance` ) faz parte da situação financeira do atendimento. 

## **3. TELA — LOGIN** 

### **Objetivo** 

Permitir autenticação do usuário. 

### **Elementos** 

- campo Login; 

- campo Senha; 

- botão Entrar. 

### **Fluxo** 

```
Usuário informa login/senha
 ↓
Frontend envia POST /api/v1/auth/login
 ↓
```

```
API autentica
 ↓
```

```
Sucesso → guardar autenticação
 ↓
```

#### `Dashboard` 

### **Erros** 

Exibir mensagem amigável quando: 

- login/senha inválidos; 

- servidor indisponível; 

- erro inesperado. 

Não mostrar detalhes técnicos do backend ao usuário. 

## **4. TELA — DASHBOARD** 

### **Objetivo** 

Ser o ponto inicial após o login. 

O dashboard deve facilitar a visualização do que precisa de atenção. 

Possíveis blocos: 

```
Atendimentos em andamento
```

```
Aguardando aprovação
Aguardando pagamento
```

```
Agendamentos do dia
```

A quantidade e os indicadores definitivos devem ser definidos conforme o uso real. 

### **Ações** 

O usuário deve conseguir acessar rapidamente: 

- atendimento; 

- agenda; 

- cliente; 

- veículo; 

- caixa, se tiver permissão. 

## **5. TELA — CLIENTES (Customers)** 

### **Objetivo** 

Cadastrar e localizar clientes. 

### **Listagem** 

Exibir: 

- nome; 

- telefone; 

- ações. 

### **Busca** 

Permitir procurar cliente existente. 

Exemplo: 

```
[ Buscar cliente... ]
```

### **Ações** 

```
Novo cliente
```

```
Visualizar
```

```
Editar
```

## **6. TELA — CADASTRO DE CLIENTE** 

Campos: 

- nome ( `name` ); 

- telefone ( `phone` ). 

### **Fluxo** 

```
Preencher
```

```
 ↓
```

```
Salvar
```

```
 ↓
```

```
API valida
```

```
 ↓
```

```
Sucesso
```

```
 ↓
```

```
Cliente cadastrado
```

O frontend não deve decidir sozinho quais dados são válidos. A API é a autoridade para validação. 

## **7. TELA — VEÍCULOS (Vehicles)** 

### **Objetivo** 

Cadastrar e localizar veículos. 

Campos atualmente previstos: 

- modelo ( `model` ); 

- cor ( `color` ); 

- placa ( `licensePlate` ); 

- ano ( `year` ); 

- quilometragem ( `mileage` ). 

Alguns campos podem ser opcionais conforme a situação do atendimento. 

### **Ações** 

- cadastrar; 

- consultar; 

• selecionar veículo. 

## **8. TELA — AGENDA (Appointments)** 

### **Objetivo** 

Mostrar os agendamentos da oficina. 

### **Exibir** 

- data; 

- horário; 

- cliente; 

- veículo; 

- status. 

### **Status** 

- Agendado ( `SCHEDULED` ); 

- Cancelado ( `CANCELLED` ); 

- Não compareceu ( `NO_SHOW` ). 

Quando o cliente comparecer: `Agendamento ↓` 

```
Iniciar atendimento
```

## **9. TELA — NOVO AGENDAMENTO** 

Campos: 

- cliente ( `customerId` ); 

- veículo ( `vehicleId` ); 

- data/horário ( `scheduledAt` ). 

### **Fluxo** 

```
Selecionar cliente
```

```
 ↓
```

```
Selecionar veículo
```

```
 ↓
```

```
Escolher data/hora
```

```
 ↓
```

```
Salvar
```

```
 ↓
```

```
API verifica regras/conflitos
```

Se houver conflito, mostrar o erro retornado pela API. Não tentar reproduzir toda a lógica de conflito no frontend. 

## **10. TELA — ATENDIMENTOS (Service Orders)** 

### **Objetivo** 

Mostrar os atendimentos da oficina. 

Cada atendimento deve permitir identificar rapidamente: 

- cliente; 

- veículo; 

- situação; 

- pendência principal. 

Exemplo conceitual: 

```
João — Gol Prata
```

```
Aguardando aprovação
```

```
Maria — Civic Preto
```

```
Em execução
```

```
Carlos — Onix Branco
```

```
Aguardando pagamento
```

## **11. TELA — NOVO ATENDIMENTO** 

Existem dois caminhos. 

### **Caminho A — Sem agendamento** 

```
Novo atendimento
```

```
 ↓
```

```
Localizar/cadastrar cliente
```

```
 ↓
```

```
Localizar/cadastrar veículo
```

```
 ↓
```

```
Criar atendimento
```

### **Caminho B — A partir da agenda** 

```
Agendamento
```

```
 ↓
```

```
Iniciar atendimento
 ↓
```

```
Sistema reutiliza cliente + veículo
```

```
 ↓
```

```
Atendimento criado
```

Não pedir novamente informações que a API já possui. 

## **12. TELA — DETALHES DO ATENDIMENTO** 

Esta é uma das telas centrais do sistema. Ela deve funcionar como o "prontuário" daquele serviço. 

### **Cabeçalho** 

```
Atendimento #123
Cliente: João
Telefone: ...
Veículo: Gol / Prata / ABC1D23
```

### **Abas/seções** 

```
Resumo
Avaliação
Orçamento
Execução
Pagamentos
Histórico
Entrega
```

O frontend pode escolher tabs, accordion ou outra organização visual. 

## **13. SEÇÃO — AVALIAÇÃO (Assessment)** 

Permitir registrar: 

- problema relatado ( `reportedProblem` ); 

- informações da avaliação ( `assessment` ); 

- diagnóstico ( `diagnosis` ). 

Também deve ser possível indicar que o diagnóstico ainda não foi concluído. 

### **Ações** 

```
Salvar avaliação
```

```
Salvar diagnóstico
```

## **14. SEÇÃO — ORÇAMENTO (Estimate)** 

Mostrar: 

```
ORÇAMENTO
```

```
Item              Valor   Status
------------------------------------------------
Troca de correia  R$300   Aprovado
Correia dentada   R$250   Aprovado
Limpeza           R$100   Recusado
```

```
Total: R$550
```

Cada item possui seu próprio estado ( `PENDING` / `APPROVED` / `REJECTED` ). 

## **15. TELA/COMPONENTE — CRIAÇÃO DE ORÇAMENTO** 

Permitir adicionar: 

```
Tipo:
```

```
( ) Serviço   [type = SERVICE]
( ) Peça      [type = PART]
```

```
Descrição     [description]
Valor         [unitPrice]
```

```
[Adicionar item]
```

Depois: 

```
Itens
Total
```

```
[Salvar orçamento]
```

O total exibido pode ser calculado visualmente para facilitar a interação, mas o valor oficial deve ser determinado/validado pelo backend. 

## **16. SEÇÃO — APROVAÇÃO (Decision)** 

Mostrar claramente o estado de cada item: 

- `✓ Aprovado   (APPROVED)` 

- `✗ Recusado   (REJECTED)` 

- `Pendente   (PENDING)` 

Quando uma decisão for registrada: 

`Decisão registrada Cliente: João Registrado por: Dono Data/hora: ...` A interface deve distinguir: quem tomou a decisão = cliente ( `customerId` ) de: 

quem registrou = usuário da oficina ( `recordedBy` ). 

## **17. SEÇÃO — HISTÓRICO DO ORÇAMENTO (Estimate History)** 

Mostrar cronologicamente: 

```
Versão 1
 ↓
Orçamento inicial
```

```
Versão 2
 ↓
Correia alterada de R$250 → R$320
 ↓
Motivo: ...
```

```
Decisão do cliente
 ↓
Aprovado
```

O objetivo é permitir entender o que aconteceu sem apagar versões anteriores. 

## **18. ALERTA — ALTERAÇÃO PENDENTE (Estimate Change)** 

Quando existir uma alteração aguardando decisão: 

A tela deve deixar isso evidente. 

Exemplo: 

⚠ `Alteração aguardando aprovação Peça: Correia dentada Valor anterior: R$250 Novo valor: R$320 Status: Pendente` 

A interface não deve permitir que o usuário trate a alteração como aprovada antes da decisão correspondente. 

A API continua sendo responsável pela validação definitiva. 

## **19. SEÇÃO — EXECUÇÃO (Execution)** 

Mostrar os itens efetivamente realizados. 

**Serviços** 

```
Serviços realizados
```

```
Troca de correia   R$300
```

### **Peças** 

```
Peças utilizadas
Correia dentada
Custo oficina: R$180      (workshopCost)
Valor cobrado: R$250      (customerPrice)
```

O custo da oficina ( `workshopCost` ) não deve ser confundido com o valor cobrado do cliente ( `customerPrice` ). 

## **20. COMPARAÇÃO — ORÇADO × REALIZADO** 

A interface deve facilitar a visualização: 

```
ORÇADO                REALIZADO
Correia R$250     →   Correia R$320
Filtro  R$100     →   Não utilizado
Limpeza R$80      →   Limpeza R$80
```

Não transformar o item orçado ( `estimateItem` ) no item realizado ( `performedService` / `usedPart` ). O usuário precisa conseguir identificar a diferença. 

## **21. INTERRUPÇÃO DA EXECUÇÃO** 

Permitir registrar: 

```
Status: Interrompido    (INTERRUPTED)
```

```
Observação:
```

```
[Aguardando autorização do cliente...]
```

Neste MVP não criar uma interface complexa de categorias de interrupção. 

## **22. VALOR FINAL (Final Amount)** 

Mostrar em destaque: 

```
Valor final
R$ 1.250,00
```

A tela pode detalhar: 

```
Serviços realizados: R$700
Peças utilizadas: R$550
```

```
Total: R$1.250
```

O valor oficial vem da API. 

## **23. SEÇÃO — PAGAMENTOS (Payments)** 

Mostrar: 

```
Valor final: R$1.250
```

```
Pagamentos:
```

```
20/09  PIX       R$150
21/09  DINHEIRO  R$300
22/09  CRÉDITO   R$400
```

```
Total pago: R$850
```

```
Saldo: R$400
```

## **24. REGISTRAR PAGAMENTO** 

Formulário: 

```
Valor            [amount]
Forma de pagamento [paymentMethod]
Data             [paidAt]
```

```
[Registrar pagamento]
```

Formas ( `paymentMethod` ): 

- Dinheiro ( `CASH` ); 

- PIX ( `PIX` ); 

- Débito ( `DEBIT_CARD` ); 

- Crédito ( `CREDIT_CARD` ). 

Não implementar quantidade de parcelas no MVP. 

## **25. PAGAMENTO ANTECIPADO** 

A mesma tela deve permitir registrar um pagamento antes da conclusão. 

Exemplo: 

```
Sinal
R$150
PIX
```

O sistema deve incorporá-lo automaticamente ao cálculo posterior do saldo ( `balance` ). 

## **26. SEÇÃO — ENTREGA (Delivery)** 

Mostrar claramente a situação financeira. 

### **Sem saldo** 

```
✓ Pagamento concluído
Saldo: R$0
[Registrar entrega]
```

### **Com saldo** 

⚠ `Pagamento pendente Saldo: R$400 Entrega bloqueada.` 

### **Exceção** 

Se o usuário autorizado puder registrar a exceção: 

```
[Registrar exceção]        (balanceException: true)
Observação:                (observation)
```

Depois disso, a API decide se a operação é válida conforme as regras. 

## **27. TELA — CAIXA (Cash)** 

Apenas usuário autorizado. 

Mostrar: 

```
CAIXA — 20/09/2026
Dinheiro   R$500
PIX        R$850
Débito     R$200
Crédito    R$300
--------------------
Total      R$1.850
```

Abaixo: 

```
Pagamentos do dia
Cliente | Atendimento | Forma | Valor
```

Os dados devem vir dos pagamentos registrados. 

## **28. CONTROLE DE PERMISSÕES NO FRONTEND** 

O frontend deve adaptar a interface ao perfil. 

### **Dono (** **`OWNER` )** 

Pode visualizar todas as funções permitidas. 

**Mecânico (** **`MECHANIC` )** 

Não deve receber ações que não pode executar. 

Por exemplo, não mostrar como ação disponível: 

```
Registrar pagamento
Liberar entrega
```

Porém: esconder o botão não é segurança. Mesmo que alguém tente chamar a API manualmente, o backend deve rejeitar a operação. 

## **29. ESTADOS IMPORTANTES DA INTERFACE** 

O frontend deve tratar pelo menos: 

```
Carregando
Sucesso
Erro
Vazio
Sem permissão
Aguardando aprovação
Aguardando pagamento
Interrompido
Finalizado
Cancelado
```

Não deixar a interface depender somente de "há dados / não há dados". 

## **30. FLUXO PRINCIPAL COMPLETO DA INTERFACE** 

```
LOGIN
 ↓
```

```
DASHBOARD
```

```
 ↓
AGENDA ou NOVO ATENDIMENTO
 ↓
CLIENTE
 ↓
VEÍCULO
 ↓
```

```
ATENDIMENTO
 ↓
AVALIAÇÃO
 ↓
DIAGNÓSTICO
```

```
 ↓
```

```
ORÇAMENTO
```

```
 ↓
APROVAÇÃO
 ↓
EXECUÇÃO
 ↓
VALOR FINAL
 ↓
```

```
PAGAMENTO
```

```
 ↓
SALDO
 ↓
```

```
ENTREGA
```

# **PARTE 2 — CONEXÃO FRONTEND ↔ BACKEND** 

## **31. Arquitetura** 

O frontend deve consumir somente a API REST. 

```
┌───────────────┐
│   FRONTEND    │
│               │
│  React/etc.   │
└───────┬───────┘
        │ HTTP/HTTPS
        ↓
┌───────────────┐
│   REST API    │
│               │
│  Regras       │
│  Autorização  │
│  Validação    │
└───────┬───────┘
        │
        ↓
```

```
┌───────────────┐
│  PostgreSQL   │
└───────────────┘
```

O frontend não acessa o PostgreSQL. 

## **32. Organização da camada de API** 

Criar uma camada responsável pela comunicação HTTP. 

Conceito: `components/ pages/ hooks/ services/ types/` 

Dentro de `services` , organizar por domínio: 

```
services/
├── auth
├── customers
├── vehicles
├── appointments
├── service-orders
├── estimates
├── executions
├── payments
├── deliveries
└── cash
```

A estrutura exata pode variar conforme o framework utilizado. 

## **33. Cliente HTTP** 

Centralizar a configuração HTTP. Exemplo conceitual: 

```
httpClient
 ↓
baseURL
 ↓
headers
 ↓
```

```
token
```

```
 ↓
```

```
tratamento comum de erros
```

Evitar espalhar chamadas fetch/Axios diretamente por todas as telas. 

## **34. URL BASE** 

A URL da API deve ser configurável por ambiente. 

Exemplo: 

```
VITE_API_URL=http://localhost:8080
```

Produção terá outra URL. 

Não colocar a URL diretamente em dezenas de arquivos. 

## **35. Autenticação** 

Após login: 

```
POST /api/v1/auth/login
```

Backend retorna autenticação. 

O frontend deve: 

1. armazenar a informação necessária para autenticação ( `accessToken` ); 

2. enviar a credencial nas requisições protegidas; 

3. identificar o usuário atual ( `user` ); 

4. controlar a navegação conforme o perfil ( `role` ). 

A forma exata de armazenamento/token deve ser definida conforme a estratégia de segurança adotada pelo backend. 

## **36. Interceptor/middleware HTTP** 

Criar mecanismo central para: 

```
Request
```

```
 ↓
```

```
Adicionar autenticação
 ↓
```

```
Enviar
 ↓
```

`Receber resposta ↓ Tratar erros comuns` Exemplo: 

```
401 → sessão inválida → voltar para login
```

```
403 → sem permissão → mostrar acesso negado
500 → erro interno → mensagem genérica
```

## **37. Tipagem dos dados** 

Criar tipos/interfaces correspondentes aos contratos da API. 

Exemplo: `interface Customer { id: number; name: string; phone: string; }` E: `interface Payment { id: number; amount: number; paymentMethod: PaymentMethod; paidAt: string; }` 

Os tipos devem acompanhar o contrato oficial da API. Não duplicar estruturas do banco desnecessariamente. 

## **38. Serviços da API** 

Exemplo conceitual: `customerService.create() customerService.list() customerService.getById() customerService.update()` 

Para atendimento: `serviceOrderService.create() serviceOrderService.getById() serviceOrderService.list() serviceOrderService.startFromAppointment()` Para orçamento: `estimateService.create() estimateService.getById() estimateService.registerDecision()` 

```
estimateService.registerChange()
```

Para pagamento: 

```
paymentService.create()
```

```
paymentService.list()
```

## **39. Fluxo de uma operação** 

Exemplo: registrar pagamento. 

```
Usuário
```

```
 ↓
```

```
Preenche formulário
 ↓
```

```
Frontend valida formato básico
 ↓
```

```
paymentService.create()
```

```
 ↓
POST /api/v1/service-orders/{id}/payments
 ↓
API valida regra de negócio
 ↓
Banco registra pagamento
 ↓
```

```
API retorna pagamento + estado atualizado
```

```
 ↓
```

```
Frontend atualiza tela
```

## **40. Regra: frontend não é autoridade** 

Exemplo: 

O frontend mostra: 

```
Saldo: R$0
```

Mas isso não significa que ele possa simplesmente liberar a entrega. 

Ao registrar entrega: 

```
POST /api/v1/service-orders/{id}/delivery
```

A API deve recalcular/verificar o estado financeiro e decidir se a operação pode acontecer. O frontend apenas apresenta o resultado. 

## **41. Atualização de estado** 

Depois de operações importantes, atualizar os dados relacionados. 

Exemplo: 

```
Registrar pagamento
```

```
 ↓
```

```
API responde
```

```
 ↓
```

```
Atualizar:
```

- `payments` 

- `paidAmount` 

- `balance` 

- `delivery (estado da entrega)` 

Evitar manter valores financeiros duplicados em estados independentes que possam ficar inconsistentes. 

## **42. Formulários** 

Todos os formulários devem possuir: 

```
Estado inicial
```

```
 ↓
```

```
Preenchimento
```

```
 ↓
Validação
```

```
 ↓
```

```
Loading
 ↓
```

`Sucesso / erro` Durante envio: `[Salvando...]` Evitar múltiplos envios acidentais. 

## **43. Tratamento de erros da API** 

O frontend deve interpretar os códigos HTTP e, quando disponíveis, os códigos de erro do backend. Exemplo: 

```
{
  "error": {
    "code": "ESTIMATE_NOT_APPROVED",
```

```
    "message": "The estimate item has not been approved."
  }
```

```
}
```

Mostrar ao usuário a mensagem apropriada. 

Não exibir: 

```
NullPointerException
```

```
SQLSTATE...
```

```
stack trace...
```

## **44. Conflitos** 

Exemplo: 

Dois usuários tentam agendar o mesmo horário. 

O frontend pode verificar a agenda previamente, mas a API é a autoridade final. 

Se receber: 

```
409 Conflict
```

o frontend deve: 

1. informar que o horário não está mais disponível; 

2. atualizar a agenda; 

3. permitir escolher outro horário. 

## **45. Loading** 

Toda operação assíncrona deve possuir estado visual. 

Exemplo: 

```
Carregando clientes...
```

```
Salvando orçamento...
```

```
Registrando pagamento...
```

```
Calculando valor final...
```

Evitar telas aparentemente travadas. 

## **46. Cache e consultas** 

Listagens que são reutilizadas frequentemente podem utilizar cache/query management conforme a tecnologia escolhida. 

Exemplos: 

- clientes ( `customers` ); 

- veículos ( `vehicles` ); 

- agenda ( `appointments` ); 

- atendimento atual ( `service-order` ). 

Porém, dados que podem mudar rapidamente devem ser atualizados após mutações relevantes. 

## **47. Dados do atendimento** 

A tela de atendimento provavelmente será a tela mais importante do sistema. 

O frontend deve evitar fazer várias chamadas desnecessárias para montar a mesma tela se a API fornecer uma resposta agregada adequada. 

Exemplo conceitual: 

`GET /api/v1/service-orders/{id}` pode retornar o contexto: `{ "customer": {}, "vehicle": {}, "appointment": {}, "assessment": {}, "estimate": {}, "execution": {}, "payments": [], "balance": 0, "delivery": null }` A estrutura definitiva será definida pelo contrato da API. 

## **48. Sincronização entre usuários** 

Como Dono e Mecânico utilizam o mesmo sistema: `Dono altera agenda ↓` 

```
Backend
 ↓
Mecânico consulta agenda
 ↓
Novo estado
```

O frontend não deve presumir que seus dados locais são sempre atuais. 

A estratégia de atualização pode ser: 

- refetch; 

- polling; 

- atualização após ações; 

- WebSocket/SSE, se posteriormente houver necessidade. 

Não implementar WebSocket apenas por antecipação. 

## **49. Responsividade** 

O sistema deve funcionar em: 

- desktop; 

- notebook; 

- smartphone. 

### **Priorizar:** 

Desktop Uso mais amplo de tabelas, painéis e múltiplas informações. Smartphone Priorizar: 

- informações essenciais; 

- ações principais; 

- formulários; 

- consulta rápida. 

Não simplesmente diminuir uma tela desktop até caber no celular. 

## **50. Critérios de aceite do frontend** 

### **Fluxo 1** 

Usuário consegue: 

```
Login
```

- `→ Dashboard` 

- `→ localizar cliente` 

- `→ localizar veículo` 

- `→ criar atendimento` 

### **Fluxo 2** 

Usuário consegue: 

```
Agenda
```

- `→ selecionar agendamento` 

- `→ iniciar atendimento` 

- `→ reutilizar cliente/veículo` 

### **Fluxo 3** 

Usuário consegue: 

```
Atendimento
```

- `→ avaliação` 

- `→ diagnóstico` 

- `→ orçamento` 

- `→ visualizar itens` 

```
→ visualizar aprovação
```

### **Fluxo 4** 

Usuário consegue visualizar: 

```
Orçamento original
```

- `→ alterações` 

- `→ decisões` 

- `→ execução` 

sem perder o histórico. 

### **Fluxo 5** 

Usuário consegue: 

```
Execução
```

- `→ serviços realizados` 

- `→ peças utilizadas` 

- `→ valor final` 

### **Fluxo 6** 

Usuário consegue: 

```
Pagamento antecipado
```

- `→ pagamentos posteriores` 

- `→ total pago` 

- `→ saldo` 

### **Fluxo 7** 

Usuário autorizado consegue: 

```
Saldo = 0
```

- `→ registrar entrega` 

### **Fluxo 8** 

Usuário não autorizado recebe: 

```
403
```

- `→ ação bloqueada` 

mesmo que tente acessar a operação diretamente. 

## **51. Ordem recomendada de implementação do frontend** 

`1. Estrutura do projeto` 

```
 ↓
```

`2. Cliente HTTP` 

```
 ↓
```

`3. Autenticação` 

```
 ↓
```

`4. Layout + navegação` 

```
 ↓
```

`5. Clientes` 

```
 ↓
```

`6. Veículos` 

```
 ↓
```

`7. Agenda` 

```
 ↓
```

`8. Atendimentos` 

```
 ↓
```

`9. Avaliação/Diagnóstico` 

```
 ↓
```

`10. Orçamento` 

```
 ↓
```

`11. Aprovação/Histórico` 

```
 ↓
```

`12. Execução` 

```
 ↓
```

`13. Valor final` 

```
 ↓
```

`14. Pagamentos` 

```
 ↓
```

`15. Entrega` 

```
 ↓
```

`16. Caixa` 

```
 ↓
```

`17. Permissões` 

```
 ↓
```

`18. Responsividade` 

```
 ↓
```

`19. Testes de integração` 

## **52. Dependência entre as equipes** 

O trabalho pode acontecer em paralelo: 

```
 ┌──→ Ricardo: Banco
```

```
 │
Requisitos ──────┼──→ Backend: API
 │
 └──→ Frontend: Telas
```

Mas existe um contrato entre frontend e backend: 

```
       CONTRATO DA API
```

```
              ↓
   ┌─────────┴─────────┐
   ↓                   ↓
BACKEND            FRONTEND
Implementa         Consome
   ↓                   ↓
   └─────── HTTP ──────┘
```

O frontend não deve esperar o backend terminar para começar. Ele pode começar com mocks baseados no contrato. 

## **53. Estratégia de desenvolvimento paralelo** 

### **Etapa 1** 

Backend define: 

- endpoints; 

- requests; 

- responses; 

- erros; 

- autenticação; 

- permissões. 

### **Etapa 2** 

Frontend cria: 

- telas; 

- componentes; 

- navegação; 

- estados; 

- formulários; 

usando dados mockados conforme o contrato. 

### **Etapa 3** 

Backend disponibiliza os endpoints. 

### **Etapa 4** 

Frontend substitui mocks pelas chamadas reais. 

### **Etapa 5** 

Os dois lados executam os cenários de aceite. 

## **54. Regra de ouro para o frontend** 

O frontend deve responder: 

"Como o usuário realiza esta operação?" 

O backend deve responder: 

"Esta operação é válida?" 

E o banco deve responder: 

"Como os dados necessários para isso são persistidos?" 

Portanto: 

```
FRONT
```

```
Interface + interação + estado visual
```

```
 ↓
```

```
API
```

```
Regras + validações + autorização
 ↓
```

```
BANCO
```

```
Persistência + integridade
```

Essa separação deve ser mantida durante todo o desenvolvimento do MVP. 

