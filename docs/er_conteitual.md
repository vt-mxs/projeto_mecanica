# Schema: projeto_mecanica
Database: pgsql
Tables: 11
Columns are NOT NULL unless marked nullable.

## clients
id: uuid PK
user_id: uuid FK
name: varchar
email: varchar unique nullable
phone: bigint unique

## vehicle
id: uuid PK
client_id: uuid FK
v_model: varchar
v_color: varchar
license_plate: char unique -- 10 chars
v_year: date -- Use date.Year
v_km: bigint

## meet
id: uuid PK
v_id: uuid FK
client_id: uuid
meet_date: date -- usar date com time
meet_status: enum default:not_appear

## treatment
id: uuid PK
v_id: uuid FK
meet_id: uuid FK nullable
problem: varchar
rate: varchar
billing_id: uuid FK
serv_execution_id: uuid FK
payments: bigint[]
balance: bigint
is_Delivered: boolean default:false

## billing
id: uuid PK
services_id: uuid
total_value: bigint FK
approvement_status: enum default:not_aproved
changes_id: uuid nullable
decisions_id: uuid

## services
id: uuid PK FK
s_type: varchar
description: varchar
value: bigint
vehicle_id: bigint FK
vehicle_part: bigint
vehiclePart_billing: bigint nullable
vehiclePart_MCost: bigint nullable

## operation_history
id: bigint PK

## users
id: uuid PK
email: varchar unique
password: varchar
role: enum default:owner

## decision
id: uuid PK FK
services_id: uuid FK
decision_status: enum default:not_aproved
decision_description: varchar

## changes
id: uuid PK FK
services_id: uuid FK
old_value: bigint
new_billing_value: bigint nullable
changes_description: varchar nullable
decision_id: uuid

## payments
id: uuid PK
orderId: bigint
valueInCents: bigint
currency: char -- char de 5, max 8
services_id: uuid FK
pay_method: varchar
payment_status: enum default:uncomplete

## Relationships
clients.id → vehicle.client_id (one-to-many)
meet.v_id → vehicle.id (one-to-one)
meet.client_id → vehicle.client_id (one-to-one)
treatment.v_id → vehicle.id (one-to-one)
treatment.meet_id → meet.id (one-to-one)
billing.id → treatment.billing_id (one-to-one)
services.id → treatment.serv_execution_id (one-to-one)
billing.services_id → services.id (one-to-many)
services.vehicle_id → vehicle.id (many-to-one)
users.id → clients.user_id (one-to-one)
decision.services_id → services.id (one-to-one)
billing.decisions_id → decision.id (one-to-many)
changes.id → billing.changes_id (many-to-one)
changes.services_id → services.id (one-to-one)
changes.old_value → billing.total_value (one-to-one)
payments.services_id → services.id (one-to-one)

## Enums
meet.meet_status: appear, not_appear
billing.approvement_status: not_aproved, in_process, aproved
users.role: admin, owner, mechanic
decision.decision_status: not_aproved, aproved
payments.payment_status: uncomplete, inProgress, complete

## Indexes
clients.ca4da1d6-e6b4-4361-9d7a-7bd71d341314 (index): email
clients.5e97ec84-41ae-4d35-a38a-07537eaa9e46 (index): phone
vehicle.bc357813-79eb-49bd-9454-75ed79e4dad5 (index): license_plate
vehicle.0132f9b5-f512-4ec5-bf1f-938d6c2f0740 (index): v_model
meet.f8db64a7-6189-4105-b795-1c186944e2d0 (index): meet_status
treatment.e0d396d6-56eb-49b0-93f7-f3d76079e995 (index): is_Delivered
billing.b26dd399-8b8b-4441-b9a6-c7b7855f87c2 (index): approvement_status
services.3692263e-d7eb-4657-b664-ad796934bffc (index): s_type
services.d5638eb4-635b-493c-beca-1fbb78260d70 (index): vehicle_id
users.8b9ca874-060a-493b-926c-e0d3dee2e1f5 (index): email
users.9f91058d-59f7-4ad2-b29d-57c7e277932c (index): role
decision.a1e419a7-c463-49d7-ab98-defa3c8183ca (index): decision_status
changes.074c3f20-3dcb-4378-a768-f9eb711c8bd5 (index): services_id

## Notes
"Fiquei em duvida como seria a implementação da execução dos serviços, se add os campos como not null em services ou se crio outra tabela relacionada com o services." (near: services)
"Não é preciso implementar nada de gateway de pagamento, mas vou deixar a opção de escolher o metodo de pagamento." (near: payments)

---
If you suggest changes to this schema, ALSO return them as a DrawSQL patch: one fenced ```json code block, so I can apply them to my diagram. If you are only answering a question, skip the patch.

Patch rules:
- One JSON object with "strategy": "merge"; combine every change into that single patch.
- Reference tables and columns by name. Only include what changes — omit unchanged columns and unused optional fields.
- Column: {"name","type","length","is_primary_key","is_auto_increment","is_nullable","is_unique_key","is_index","is_unsigned","default","enum_values"} — use exact type names for pgsql as shown in the schema above.
- Relationship: {"type":"one-to-many","source_table":"users","source_column":"id","target_table":"posts","target_column":"user_id"} — for "one-to-many" the source is the "one" side; for "many-to-one" the source is the "many" (foreign-key) side. Prefer "one-to-many".
- Add "default_type" ("string" | "function" | "number" | "boolean") when a default could be misread — a literal string "CURRENT_TIMESTAMP" or "0" needs "default_type":"string".
- SET columns (MySQL) list their members in "set_values", not "enum_values".
- Index: {"indexes":[{"name":"idx_posts_user_created","type":"index","columns":["user_id","created_at"]}]} on its table — "type" is "index" or "unique"; delete via "deletions".
- Group: {"name":"Billing"}. A group never lists its members — put "parent_name":"Billing" on each table that belongs to it.
- Sticky note: {"sticky_notes":[{"content":"..."}]}.
- Delete via "deletions"; rename via "old_name". Do not include left/top coordinates.

Examples:
Add table: `{"strategy":"merge","tables":[{"name":"posts","comment":"Blog posts","columns":[{"name":"id","type":"bigint","is_primary_key":true,"is_auto_increment":true},{"name":"user_id","type":"bigint","is_index":true},{"name":"title","type":"varchar","length":255},{"name":"body","type":"text","is_nullable":true},{"name":"created_at","type":"timestamp","is_nullable":true},{"name":"updated_at","type":"timestamp","is_nullable":true}]}]}`
Rename table and column (use old_name only when renaming): `{"strategy":"merge","tables":[{"name":"articles","old_name":"posts","columns":[{"name":"content","old_name":"body"}]}]}`
Delete: `{"strategy":"merge","deletions":{"tables":["tmp"],"columns":[{"table":"users","columns":["legacy"]}],"relationships":[{"source_table":"users","source_column":"id","target_table":"posts","target_column":"user_id"}]}}`
Multiple groups with nested tables (every table sets parent_name): `{"strategy":"merge","groups":[{"name":"Auth"},{"name":"Content"}],"tables":[{"name":"users","parent_name":"Auth","columns":[{"name":"id","type":"bigint","is_primary_key":true,"is_auto_increment":true}]},{"name":"roles","parent_name":"Auth","columns":[{"name":"id","type":"bigint","is_primary_key":true,"is_auto_increment":true}]},{"name":"posts","parent_name":"Content","columns":[{"name":"id","type":"bigint","is_primary_key":true,"is_auto_increment":true}]}]}`
