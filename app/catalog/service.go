package catalog

import "github.com/mytheresa/go-hiring-challenge/models"

type detailsService struct{}

func newDetailsService() *detailsService {
	return &detailsService{}
}

func (s *detailsService) BuildProductDetails(product *models.Product) ProductDetailsResponse {
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

	return ProductDetailsResponse{
		Code:  product.Code,
		Price: product.Price.InexactFloat64(),
		Category: Category{
			Code: product.Category.Code,
			Name: product.Category.Name,
		},
		Variants: variants,
	}
}
