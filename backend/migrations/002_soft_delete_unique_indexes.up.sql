ALTER TABLE products DROP CONSTRAINT IF EXISTS products_sku_key;
CREATE UNIQUE INDEX IF NOT EXISTS idx_products_sku_unique
  ON products (sku)
  WHERE deleted_at IS NULL;
