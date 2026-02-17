package models

import (
	"errors"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"github.com/joho/godotenv"
	"github.com/mytheresa/go-hiring-challenge/app/database"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupDBWithSeed(t *testing.T) *gorm.DB {
	t.Helper()

	_ = godotenv.Load("../.env")

	db, closeFn := database.New(
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
	)
	t.Cleanup(func() {
		require.NoError(t, closeFn())
	})

	sqlFiles, err := os.ReadDir("../sql")
	require.NoError(t, err)

	var names []string
	for _, entry := range sqlFiles {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".sql") {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)

	for _, name := range names {
		content, readErr := os.ReadFile(filepath.Join("..", "sql", name))
		require.NoError(t, readErr)
		require.NoError(t, db.Exec(string(content)).Error)
	}

	return db
}

func TestTableNames(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "products", (&Product{}).TableName())
	assert.Equal(t, "product_variants", (&Variant{}).TableName())
	assert.Equal(t, "categories", (&Category{}).TableName())
}

func TestCategoriesRepositoryCreateAndList(t *testing.T) {
	db := setupDBWithSeed(t)
	repo := NewCategoriesRepository(db)

	list, err := repo.GetAllCategories()
	require.NoError(t, err)
	assert.Len(t, list, 3)

	created, err := repo.CreateCategory(Category{Code: "BAGS", Name: "Bags"})
	require.NoError(t, err)
	assert.NotZero(t, created.ID)

	list, err = repo.GetAllCategories()
	require.NoError(t, err)
	assert.Len(t, list, 4)

	_, err = repo.CreateCategory(Category{Code: "BAGS", Name: "Bags Duplicate"})
	assert.Error(t, err)
}

func TestCategoriesRepositoryErrorBranches(t *testing.T) {
	db := setupDBWithSeed(t)
	repo := NewCategoriesRepository(db)

	require.NoError(t, db.Exec("DROP TABLE categories CASCADE").Error)

	_, err := repo.GetAllCategories()
	assert.Error(t, err)

	_, err = repo.CreateCategory(Category{Code: "X", Name: "X"})
	assert.Error(t, err)
}

func TestProductsRepositoryListProductsAndFilters(t *testing.T) {
	db := setupDBWithSeed(t)
	repo := NewProductsRepository(db)

	products, total, err := repo.ListProducts(ProductCatalogFilter{Offset: 0, Limit: 10})
	require.NoError(t, err)
	assert.Equal(t, int64(8), total)
	assert.Len(t, products, 8)

	products, total, err = repo.ListProducts(ProductCatalogFilter{Offset: 1, Limit: 2})
	require.NoError(t, err)
	assert.Equal(t, int64(8), total)
	assert.Len(t, products, 2)

	products, total, err = repo.ListProducts(ProductCatalogFilter{Offset: 0, Limit: 10, Category: "Shoes"})
	require.NoError(t, err)
	assert.Equal(t, int64(2), total)
	assert.Len(t, products, 2)

	price := decimal.RequireFromString("10")
	products, total, err = repo.ListProducts(ProductCatalogFilter{Offset: 0, Limit: 10, PriceLessThan: &price})
	require.NoError(t, err)
	assert.Equal(t, int64(3), total)
	assert.Len(t, products, 3)
}

func TestProductsRepositoryErrorBranches(t *testing.T) {
	db := setupDBWithSeed(t)
	repo := NewProductsRepository(db)

	require.NoError(t, db.Exec("DROP TABLE products CASCADE").Error)

	_, _, err := repo.ListProducts(ProductCatalogFilter{Offset: 0, Limit: 10})
	assert.Error(t, err)

	_, err = repo.GetProductByCode("PROD001")
	assert.Error(t, err)
}

func TestProductsRepositoryGetByCode(t *testing.T) {
	db := setupDBWithSeed(t)
	repo := NewProductsRepository(db)

	product, err := repo.GetProductByCode("PROD001")
	require.NoError(t, err)
	assert.Equal(t, "PROD001", product.Code)
	assert.Equal(t, "CLOTHING", product.Category.Code)
	assert.NotEmpty(t, product.Variants)

	_, err = repo.GetProductByCode("MISSING")
	assert.Error(t, err)
	assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
}
