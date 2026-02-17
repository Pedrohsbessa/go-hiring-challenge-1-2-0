ALTER TABLE products
ADD COLUMN IF NOT EXISTS category_id INTEGER;

DO $$
BEGIN
	IF NOT EXISTS (
		SELECT 1
		FROM pg_constraint
		WHERE conname = 'fk_products_category'
	) THEN
		ALTER TABLE products
		ADD CONSTRAINT fk_products_category
		FOREIGN KEY (category_id)
		REFERENCES categories(id)
		ON DELETE RESTRICT;
	END IF;
END $$;
