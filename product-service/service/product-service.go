package service

import (
	"context"

	pb "be-golang-chapter-57/product-service/proto"

	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name     string
	Price    uint32
	Quantity uint32
}

type ProductServiceServer struct {
	DB *gorm.DB
	pb.UnimplementedProductServiceServer
}

func (s *ProductServiceServer) CreateProduct(ctx context.Context, req *pb.CreateProductRequest) (*pb.CreateProductResponse, error) {
	product := Product{Name: req.Name, Price: req.Price, Quantity: req.Quantity}
	s.DB.Create(&product)

	return &pb.CreateProductResponse{Product: &pb.Product{
		Id:       uint32(product.ID),
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}}, nil
}

func (s *ProductServiceServer) GetProduct(ctx context.Context, req *pb.GetProductRequest) (*pb.GetProductResponse, error) {
	var product Product
	if err := s.DB.First(&product, req.Id).Error; err != nil {
		return nil, err
	}

	return &pb.GetProductResponse{Product: &pb.Product{
		Id:       uint32(product.ID),
		Name:     product.Name,
		Price:    product.Price,
		Quantity: product.Quantity,
	}}, nil
}

func (s *ProductServiceServer) ListProducts(ctx context.Context, req *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	var products []Product
	if err := s.DB.Find(&products).Error; err != nil {
		return nil, err
	}

	var pbProducts []*pb.Product
	for _, product := range products {
		pbProducts = append(pbProducts, &pb.Product{
			Id:       uint32(product.ID),
			Name:     product.Name,
			Price:    product.Price,
			Quantity: product.Quantity,
		})
	}

	return &pb.ListProductsResponse{Products: pbProducts}, nil
}
