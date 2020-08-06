package service

import (
	"context"
	"github.com/sirupsen/logrus"
	jewerly "github.com/zhashkevych/jewelry-shop-backend"
	"github.com/zhashkevych/jewelry-shop-backend/pkg/repository"
	"github.com/zhashkevych/jewelry-shop-backend/storage"
	"io"
	"math/rand"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	fileNameLength = 16
)

type ProductService struct {
	repo        repository.Product
	fileStorage storage.Storage
}

func NewProductService(repo repository.Product, fileStorage storage.Storage) *ProductService {
	return &ProductService{repo: repo, fileStorage: fileStorage}
}

func (s *ProductService) Create(product jewerly.CreateProductInput) error {
	return s.repo.Create(product)
}

// todo: add total calculation, additional filters & pagination
func (s *ProductService) GetAll(filters jewerly.GetAllProductsFilters) (jewerly.ProductsList, error) {
	productList, err := s.repo.GetAll(filters)
	if err != nil {
		return productList, err
	}

	for i, product := range productList.Products {
		images, err := s.repo.GetProductImages(product.Id)
		if err != nil {
			logrus.Errorf("failed to get images for product id %d: %s", product.Id, err.Error())
			continue
		}

		productList.Products[i].Images = images
	}

	return productList, nil
}

func (s *ProductService) Delete(id int) error {
	return s.repo.Delete(id)
}

func (s *ProductService) GetById(id int, language string) (jewerly.ProductResponse, error) {
	product, err := s.repo.GetById(id, language)
	if err != nil {
		return product, err
	}

	images, err := s.repo.GetProductImages(product.Id)
	if err != nil {
		logrus.Errorf("failed to get images for product id %d: %s", product.Id, err.Error())
		return product, err
	}

	product.Images = images

	return product, nil
}

func (s *ProductService) UploadImage(ctx context.Context, file io.Reader, size int64, contentType string) (int, error) {
	url, err := s.fileStorage.Upload(ctx, storage.UploadInput{
		File:        file,
		Name:        generateFileName(),
		Size:        size,
		ContentType: contentType,
	})
	if err != nil {
		return 0, err
	}

	return s.repo.CreateImage(url, "")
}

func generateFileName() string {
	b := make([]byte, fileNameLength)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}