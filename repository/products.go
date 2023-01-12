package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return ProductRepository{db}
}

func (p *ProductRepository) AddProduct(product model.Product) error {
	p.db.Create(&product)
	return nil
}

func (p *ProductRepository) ReadProducts() ([]model.Product, error) {
	results := []model.Product{}
	rows, err := p.db.Table("products").Select("*").Where("deleted_at is null").Rows()
	defer rows.Close()
	for rows.Next() {
		p.db.ScanRows(rows, &results)
	}

	if err != nil {
		return []model.Product{}, err
	}
	return results, nil
}

func (p *ProductRepository) DeleteProduct(id uint) error {
	tx := p.db.Where("id = ?", id).Delete(&model.Product{})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

func (p *ProductRepository) UpdateProduct(id uint, product model.Product) error {
	p.db.Model(&model.Product{}).Where("id = ?", id).Updates(product)
	return nil
}
