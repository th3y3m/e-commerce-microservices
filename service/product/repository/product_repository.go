package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"th3y3m/e-commerce-microservices/pkg/elasticsearch_server"
	"th3y3m/e-commerce-microservices/service/product/model"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type productRepository struct {
	log           *logrus.Logger
	db            *gorm.DB
	redis         *redis.Client
	elasticClient *elasticsearch.Client // Add Elasticsearch client
}

type IProductRepository interface {
	Get(ctx context.Context, productID int64) (*Product, error)
	GetAll(ctx context.Context) ([]*Product, error)
	Create(ctx context.Context, product *Product) (*Product, error)
	Update(ctx context.Context, product *Product) (*Product, error)
	Delete(ctx context.Context, productID int64) error
	GetQuerySearch(db *gorm.DB, req *model.GetProductsRequest) *gorm.DB
	GetList(ctx context.Context, req *model.GetProductsRequest) ([]*Product, error)
}

func NewProductRepository(db *gorm.DB, redis *redis.Client, log *logrus.Logger, elasticClient *elasticsearch.Client) IProductRepository {
	return &productRepository{
		db:            db,
		redis:         redis,
		log:           log,
		elasticClient: elasticClient,
	}
}

func (pr *productRepository) Get(ctx context.Context, productID int64) (*Product, error) {
	pr.log.Infof("Fetching product with ID: %d", productID)
	var product Product
	cacheKey := fmt.Sprintf("product:%d", productID)

	// Try to get the product from Redis cache
	if pr.redis != nil {
		cachedProduct, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedProduct), &product); err == nil {
				pr.log.Infof("Product found in cache: %d", productID)
				return &product, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get product from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).First(&product, productID).Error; err != nil {
		pr.log.Errorf("Error fetching product from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		productJSON, _ := json.Marshal(product)
		if err := pr.redis.Set(ctx, cacheKey, productJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save product to Redis: %v", err)
		} else {
			pr.log.Infof("Product saved to cache: %d", productID)
		}
	}

	return &product, nil
}

func (pr *productRepository) GetAll(ctx context.Context) ([]*Product, error) {
	pr.log.Info("Fetching all products")
	var products []*Product
	cacheKey := "all_products"

	// Try to get the products from Redis cache
	if pr.redis != nil {
		cachedProducts, err := pr.redis.Get(ctx, cacheKey).Result()
		if err == nil {
			if err := json.Unmarshal([]byte(cachedProducts), &products); err == nil {
				pr.log.Info("Products found in cache")
				return products, nil
			}
		} else if err != redis.Nil {
			pr.log.Warnf("Failed to get products from Redis: %v", err)
		}
	} else {
		pr.log.Warn("Redis client is not initialized")
	}

	// If not found in cache, get from database
	if err := pr.db.WithContext(ctx).Find(&products).Error; err != nil {
		pr.log.Errorf("Error fetching products from database: %v", err)
		return nil, err
	}

	// Save to cache if Redis is available
	if pr.redis != nil {
		productsJSON, _ := json.Marshal(products)
		if err := pr.redis.Set(ctx, cacheKey, productsJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save products to Redis: %v", err)
		} else {
			pr.log.Info("Products saved to cache")
		}
	}

	return products, nil
}

func (pr *productRepository) Create(ctx context.Context, product *Product) (*Product, error) {
	pr.log.Infof("Creating product: %+v", product)
	if err := pr.db.WithContext(ctx).Create(product).Error; err != nil {
		pr.log.Errorf("Error creating product: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("product:%d", product.ProductID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		productJSON, _ := json.Marshal(product)
		if err := pr.redis.Set(ctx, cacheKey, productJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save product to Redis: %v", err)
		} else {
			pr.log.Infof("Product saved to cache: %d", product.ProductID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_products").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all products cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all products")
		}
	}

	// Return the newly created product (with any updated fields)
	return product, nil
}

func (pr *productRepository) Update(ctx context.Context, product *Product) (*Product, error) {
	pr.log.Infof("Updating product: %+v", product)
	if err := pr.db.WithContext(ctx).Save(product).Error; err != nil {
		pr.log.Errorf("Error updating product: %v", err)
		return nil, err
	}
	cacheKey := fmt.Sprintf("product:%d", product.ProductID)

	// Save to cache if Redis is available
	if pr.redis != nil {
		productJSON, _ := json.Marshal(product)
		if err := pr.redis.Set(ctx, cacheKey, productJSON, 0).Err(); err != nil {
			pr.log.Warnf("Failed to save product to Redis: %v", err)
		} else {
			pr.log.Infof("Product saved to cache: %d", product.ProductID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_products").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all products cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all products")
		}
	}

	// Return the updated product
	return product, nil
}

func (pr *productRepository) Delete(ctx context.Context, productID int64) error {
	pr.log.Infof("Deleting product with ID: %d", productID)
	if err := pr.db.WithContext(ctx).Delete(&Product{}, productID).Error; err != nil {
		pr.log.Errorf("Error deleting product: %v", err)
		return err
	}

	cacheKey := fmt.Sprintf("product:%d", productID)

	// Delete from cache if Redis is available
	if pr.redis != nil {
		if err := pr.redis.Del(ctx, cacheKey).Err(); err != nil {
			pr.log.Warnf("Failed to delete product from Redis: %v", err)
		} else {
			pr.log.Infof("Product deleted from cache: %d", productID)
		}

		// Invalidate the cache for all records
		if err := pr.redis.Del(ctx, "all_products").Err(); err != nil {
			pr.log.Warnf("Failed to invalidate all products cache: %v", err)
		} else {
			pr.log.Info("Invalidated cache for all products")
		}
	}

	return nil
}

func (pr *productRepository) GetQuerySearch(db *gorm.DB, req *model.GetProductsRequest) *gorm.DB {
	pr.log.Infof("Building query for product search: %+v", req)

	if req.IsDeleted != nil {
		db = db.Where("is_deleted = ?", *req.IsDeleted)
	}

	if req.SellerID != nil {
		db = db.Where("seller_id = ?", *req.SellerID)
	}

	if req.CategoryID != nil {
		db = db.Where("category_id = ?", *req.CategoryID)
	}

	if pr.elasticClient == nil && req.ProductName != "" {
		db = db.Where("product_name LIKE ?", "%"+req.ProductName+"%")
	}

	if req.Description != "" {
		db = db.Where("description LIKE ?", "%"+req.Description+"%")
	}

	if req.MinPrice != nil {
		db = db.Where("price >= ?", *req.MinPrice)
	}

	if req.MaxPrice != nil {
		db = db.Where("price <= ?", *req.MaxPrice)
	}

	if req.MinQuantity != nil {
		db = db.Where("quantity >= ?", *req.MinQuantity)
	}

	if req.MaxQuantity != nil {
		db = db.Where("quantity <= ?", *req.MaxQuantity)
	}

	if !req.FromDate.IsZero() {
		db = db.Where("created_at >= ?", req.FromDate)
	}

	if !req.ToDate.IsZero() {
		db = db.Where("created_at <= ?", req.ToDate)
	}

	return db
}

func (pr *productRepository) GetList(ctx context.Context, req *model.GetProductsRequest) ([]*Product, error) {
	pr.log.Infof("Fetching product list with request: %+v", req)
	var products []*Product

	// Handle the database query and filtering first
	db := pr.GetQuerySearch(pr.db.WithContext(ctx), req)

	// Handle sorting
	var sort string
	var order string

	if req.Paging.Sort == "" {
		sort = "created_at"
	} else {
		sort = req.Paging.Sort
	}

	if req.Paging.SortDirection == "" {
		order = "desc"
	} else {
		order = req.Paging.SortDirection
	}

	db = db.Order(fmt.Sprintf("%s %s", sort, order))

	// Pagination
	result := db.Table("products").
		Offset(int(req.Paging.PageIndex-1) * int(req.Paging.PageSize)).
		Limit(int(req.Paging.PageSize)).
		Find(&products)

	if result.Error != nil {
		pr.log.Errorf("Error fetching product list: %v", result.Error)
		return nil, result.Error
	}

	pr.log.Infof("Fetched %d products from the database", len(products))

	// If Elasticsearch is available, index the products and perform a search
	if pr.elasticClient == nil {
		return products, nil
	}

	// Ensure the index exists
	if err := elasticsearch_server.CreateIndex(pr.elasticClient, "products"); err != nil {
		pr.log.Errorf("Error creating index or index already exists: %v", err)
	}

	err := pr.indexProducts(products)
	if err != nil {
		pr.log.Errorf("Error indexing products to Elasticsearch: %v", err)
	}

	if req.ProductName != "" {
		products, err = pr.searchProducts(ctx, req.ProductName)
		if err != nil {
			pr.log.Errorf("Error searching products in Elasticsearch: %v", err)
		}
	}

	return products, nil
}

func (pr *productRepository) indexProducts(products []*Product) error {
	for _, product := range products {
		productJSON, err := json.Marshal(product)
		if err != nil {
			pr.log.Errorf("Error marshaling product %d: %v", product.ProductID, err)
			continue
		}

		pr.log.Infof("Indexing productID %d to Elasticsearch: %s", product.ProductID, string(productJSON))

		res, err := pr.elasticClient.Index(
			"products", // Use your Elasticsearch index name
			bytes.NewReader(productJSON),
			pr.elasticClient.Index.WithDocumentID(fmt.Sprintf("%d", product.ProductID)),
		)
		if err != nil {
			pr.log.Errorf("Error indexing product %d: %v", product.ProductID, err)
			continue
		}

		// Read the full response body
		bodyBytes, err := io.ReadAll(res.Body)
		res.Body.Close() // Close the response body immediately
		if err != nil {
			pr.log.Errorf("Failed to read Elasticsearch indexing response body: %v", err)
			continue
		}

		// Log the full response body
		pr.log.Infof("Elasticsearch indexing response body: %s", string(bodyBytes))

		if res.IsError() {
			pr.log.Errorf("Error response from Elasticsearch: %s", string(bodyBytes))
			continue
		}

		pr.log.Infof("ProductID %d indexed to Elasticsearch", product.ProductID)
	}
	return nil
}

func (pr *productRepository) searchProducts(ctx context.Context, name string) ([]*Product, error) {
	pr.log.Infof("Searching for products with name: %s", name)

	searchQuery := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"ProductName": map[string]interface{}{
								"query":     name,
								"fuzziness": "AUTO",
								"operator":  "OR",
							},
						},
					},
					{
						"match": map[string]interface{}{
							"ProductName.phonetic": map[string]interface{}{
								"query":     name,
								"fuzziness": "AUTO",
							},
						},
					},
					{
						"match_phrase": map[string]interface{}{
							"ProductName": map[string]interface{}{
								"query": name,
								"slop":  1,
							},
						},
					},
					{
						"wildcard": map[string]interface{}{
							"ProductName": fmt.Sprintf("*%s*", name),
						},
					},
				},
				"minimum_should_match": 1,
			},
		},
	}

	log.Printf("Elasticsearch search query: %v", searchQuery)

	// Serialize the search query to JSON
	queryJSON, err := json.Marshal(searchQuery)
	if err != nil {
		pr.log.Errorf("Error marshaling search query: %v", err)
		return nil, err
	}

	// Log the search query
	pr.log.Infof("Elasticsearch search query: %s", string(queryJSON))

	// Execute the search request to Elasticsearch
	res, err := pr.elasticClient.Search(
		pr.elasticClient.Search.WithContext(ctx),
		pr.elasticClient.Search.WithIndex("products"), // Elasticsearch index name
		pr.elasticClient.Search.WithBody(bytes.NewReader(queryJSON)),
		pr.elasticClient.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		pr.log.Errorf("Error executing Elasticsearch search: %v", err)
		return nil, err
	}
	defer res.Body.Close()

	// Log the Elasticsearch response status
	pr.log.Infof("Elasticsearch response status: %s", res.Status())

	// Read the full response body
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		pr.log.Errorf("Failed to read Elasticsearch response body: %v", err)
		return nil, err
	}

	// Log the full response body
	pr.log.Infof("Elasticsearch response body: %s", string(bodyBytes))

	// Parse the Elasticsearch response
	var esResponse struct {
		Hits struct {
			Hits []struct {
				Source Product `json:"_source"`
			} `json:"hits"`
		} `json:"hits"`
	}
	if err := json.Unmarshal(bodyBytes, &esResponse); err != nil {
		pr.log.Errorf("Error decoding Elasticsearch response: %v", err)
		return nil, err
	}

	// Collect products from Elasticsearch response
	products := make([]*Product, 0)
	for _, hit := range esResponse.Hits.Hits {
		products = append(products, &hit.Source)
	}

	pr.log.Infof("Found %d products from Elasticsearch", len(products))
	return products, nil
}
