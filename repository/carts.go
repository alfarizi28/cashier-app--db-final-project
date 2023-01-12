package repository

import (
	"a21hc3NpZ25tZW50/model"

	"gorm.io/gorm"
)

type CartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return CartRepository{db}
}

func (c *CartRepository) ReadCart() ([]model.JoinCart, error) {
	results := []model.JoinCart{}
	tx := c.db.Table("carts").Select("carts.id, carts.product_id, products.name, carts.quantity,carts.total_price").Joins("join products on carts.product_id = products.id").Scan(&results)
	if tx.Error != nil {
		return []model.JoinCart{}, tx.Error
	}

	return results, nil
}

func (c *CartRepository) AddCart(product model.Product) error {
	result := model.Cart{}
	tx := c.db.Raw("SELECT * FROM carts where id = ?", product.ID).Scan(&result)
	discount := product.Price - ((product.Price / 100) * product.Discount)
	carts := model.Cart{
		ProductID:  product.ID,
		Quantity:   1,
		TotalPrice: discount,
	}
	product.Stock -= 1
	if tx.RowsAffected == 0 {
		c.db.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&carts).Error; err != nil {
				return err
			}

			tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(product)

			return nil
		})
		return nil
	}

	c.db.Transaction(func(tx *gorm.DB) error {

		result.Quantity += carts.Quantity
		result.TotalPrice += carts.TotalPrice

		tx.Model(&model.Product{}).Where("id = ?", product.ID).Updates(product)

		tx.Model(&model.Cart{}).Where("id = ?", product.ID).Updates(result)

		return nil
	})

	return nil
}

func (c *CartRepository) DeleteCart(id uint, productID uint) error {
	c.db.Transaction(func(tx *gorm.DB) error {
		var carts model.Cart
		if err := tx.Raw("Select * from carts where id = ?", id).Scan(&carts).Error; err != nil {
			return err
		}

		if err := tx.Delete(&model.Cart{}, id).Error; err != nil {
			return err
		}

		var products model.Product
		if err := tx.Raw("Select * FROM products where id = ?", productID).Scan(&products).Error; err != nil {
			return err
		}

		products.Stock += int(carts.Quantity)

		tx.Model(&model.Product{}).Where("id = ?", productID).Updates(products)

		return nil
	})

	return nil
}

func (c *CartRepository) UpdateCart(id uint, cart model.Cart) error {
	c.db.Model(&model.Cart{}).Where("id = ?", id).Updates(cart)
	return nil
}
