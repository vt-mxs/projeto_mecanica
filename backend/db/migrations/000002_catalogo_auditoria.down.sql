-- ============================================================
-- 000002_catalogo_auditoria.down.sql
-- ============================================================

DROP TABLE IF EXISTS auditoria;

-- As constraints têm de cair antes das colunas, senão o DROP
-- COLUMN recusa-se por dependência de dependência.
ALTER TABLE itens_orcamento DROP CONSTRAINT IF EXISTS itens_orcamento_origem_coerente;
ALTER TABLE itens_orcamento DROP COLUMN IF EXISTS servico_id;
ALTER TABLE itens_orcamento DROP COLUMN IF EXISTS peca_id;

DROP TABLE IF EXISTS pecas;
DROP TABLE IF EXISTS servicos;