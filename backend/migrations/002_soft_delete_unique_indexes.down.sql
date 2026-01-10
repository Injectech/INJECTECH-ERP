DROP INDEX IF EXISTS idx_products_sku_unique;
ALTER TABLE products ADD CONSTRAINT products_sku_key UNIQUE (sku);
