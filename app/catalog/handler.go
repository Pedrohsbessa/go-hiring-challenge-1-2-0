package catalog

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/mytheresa/go-hiring-challenge/app/api"
	"github.com/mytheresa/go-hiring-challenge/models"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Response struct {
	Products []Product `json:"products"`
	Total    int64     `json:"total"`
}

type Product struct {
	Code     string   `json:"code"`
	Price    float64  `json:"price"`
	Category Category `json:"category"`
}

type Category struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type ProductReader interface {
	ListProducts(filter models.ProductCatalogFilter) ([]models.Product, int64, error)
	GetProductByCode(code string) (*models.Product, error)
}

type CatalogHandler struct {
	repo ProductReader
}

func NewCatalogHandler(r ProductReader) *CatalogHandler {
	return &CatalogHandler{
		repo: r,
	}
}

func (h *CatalogHandler) HandleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	offset := parseOffset(query.Get("offset"))
	limit := parseLimit(query.Get("limit"))

	var priceLessThan *decimal.Decimal
	if rawPriceLt := query.Get("price_lt"); rawPriceLt != "" {
		parsed, err := decimal.NewFromString(rawPriceLt)
		if err != nil {
			api.ErrorResponse(w, http.StatusBadRequest, "invalid query parameter: price_lt")
			return
		}
		priceLessThan = &parsed
	}

	res, total, err := h.repo.ListProducts(models.ProductCatalogFilter{
		Offset:        offset,
		Limit:         limit,
		Category:      query.Get("category"),
		PriceLessThan: priceLessThan,
	})
	if err != nil {
		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch products")
		return
	}

	products := make([]Product, len(res))
	for i, p := range res {
		products[i] = Product{
			Code:  p.Code,
			Price: p.Price.InexactFloat64(),
			Category: Category{
				Code: p.Category.Code,
				Name: p.Category.Name,
			},
		}
	}

	response := Response{
		Products: products,
		Total:    total,
	}

	api.OKResponse(w, response)
}

type ProductDetailsResponse struct {
	Code     string           `json:"code"`
	Price    float64          `json:"price"`
	Category Category         `json:"category"`
	Variants []ProductVariant `json:"variants"`
}

type ProductVariant struct {
	Name  string  `json:"name"`
	SKU   string  `json:"sku"`
	Price float64 `json:"price"`
}

func (h *CatalogHandler) HandleGetByCode(w http.ResponseWriter, r *http.Request) {
	code := r.PathValue("code")
	if code == "" {
		api.ErrorResponse(w, http.StatusBadRequest, "missing product code")
		return
	}

	product, err := h.repo.GetProductByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			api.ErrorResponse(w, http.StatusNotFound, "product not found")
			return
		}

		api.ErrorResponse(w, http.StatusInternalServerError, "failed to fetch product details")
		return
	}

	variants := make([]ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		price := product.Price
		if variant.Price != nil {
			price = *variant.Price
		}

		variants[i] = ProductVariant{
			Name:  variant.Name,
			SKU:   variant.SKU,
			Price: price.InexactFloat64(),
		}
	}

	api.OKResponse(w, ProductDetailsResponse{
		Code:  product.Code,
		Price: product.Price.InexactFloat64(),
		Category: Category{
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
		Variants: variants,
	})
}

func parseOffset(raw string) int {
	offset, err := strconv.Atoi(raw)
	if err != nil || offset < 0 {
		return 0
	}

	return offset
}

func parseLimit(raw string) int {
	if raw == "" {
		return 10
	}

	limit, err := strconv.Atoi(raw)
	if err != nil {
		return 10
	}

	if limit < 1 {
		return 1
	}

	if limit > 100 {
		return 100
	}

	return limit
}
