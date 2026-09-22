-- ============================================================
-- 000001_init_schema.down.sql
-- Remoção de todas as tabelas e tipos enumerados
-- ============================================================

DROP TABLE IF EXISTS pagamentos;
DROP TABLE IF EXISTS itens_executados;
DROP TABLE IF EXISTS alteracoes_orcamento;
DROP TABLE IF EXISTS itens_orcamento;
DROP TABLE IF EXISTS versoes_orcamento;
DROP TABLE IF EXISTS atendimentos;
DROP TABLE IF EXISTS agendamentos;
DROP TABLE IF EXISTS veiculos;
DROP TABLE IF EXISTS clientes;
DROP TABLE IF EXISTS usuarios;

DROP TYPE IF EXISTS status_pagamento;
DROP TYPE IF EXISTS forma_pagamento;
DROP TYPE IF EXISTS status_aprovacao;
DROP TYPE IF EXISTS tipo_item;
DROP TYPE IF EXISTS status_atendimento;
DROP TYPE IF EXISTS status_agendamento;
DROP TYPE IF EXISTS perfil_usuario;
