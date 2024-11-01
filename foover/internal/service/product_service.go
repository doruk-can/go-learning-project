package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"foover/internal/models"
	"foover/internal/store/mongo"
)

type ProductService interface {
	FetchAndStoreProducts(ctx context.Context, productAPIURL string) error
	IsValidProductID(ctx context.Context, productID string) (bool, error)
}

type productService struct {
	store mongo.Store
}

func NewProductService(store mongo.Store) ProductService {
	return &productService{
		store: store,
	}
}

func (p *productService) FetchAndStoreProducts(ctx context.Context, productAPIURL string) error {
	resp, err := http.Get(productAPIURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to fetch products: status code %d", resp.StatusCode)
	}

	var apiResponse struct {
		Data struct {
			MachineProducts []struct {
				ID string `json:"id"`
			} `json:"machineProducts"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		return err
	}

	var products []models.Product
	for _, mp := range apiResponse.Data.MachineProducts {
		products = append(products, models.Product{
			ProductID: mp.ID,
		})
	}

	if err := p.store.SaveProducts(ctx, products); err != nil {
		return err
	}

	return nil
}

func (p *productService) IsValidProductID(ctx context.Context, productID string) (bool, error) {
	return p.store.IsValidProductID(ctx, productID)
}
