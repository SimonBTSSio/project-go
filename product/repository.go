package product

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"go-project/model"
)

type Repository interface {
	Create(product model.Product) (model.Product, error)
	GetAll() ([]model.Product, error)
	GetById(id int) (model.Product, error)
	Update(id int, inputProduct InputProductUpdate) (model.Product, error)
	Delete(id int) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(product model.Product) (model.Product, error) {
	var prod model.Product
	r.db.Where("name = ?", product.Name).First(&prod)
	if prod.ID != 0 {
		err := errors.New("Vous ne pouvez pas créer un produit avec deux fois le même nom !")
		return product, err
	}

	err := r.db.Create(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) GetAll() ([]model.Product, error) {
	var products []model.Product
	err := r.db.Find(&products).Error
	if err != nil {
		return products, err
	}

	return products, nil
}

func (r *repository) GetById(id int) (model.Product, error) {
	var product model.Product

	err := r.db.Where(&model.Product{ID: id}).First(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Update(id int, inputProduct InputProductUpdate) (model.Product, error) {
	product, err := r.GetById(id)
	if err != nil {
		return product, err
	}

	fmt.Println("INPUT : ", inputProduct.Name)
	fmt.Println("PRODUCT : ", product.Name)

	if inputProduct.Name != "" && inputProduct.Name != product.Name {
		var prod model.Product
		r.db.Where("name = ?", inputProduct.Name).First(&prod)
		if prod.ID != 0 {
			err := errors.New("Vous ne pouvez pas créer un produit avec deux fois le même nom !")
			return product, err
		}
		product.Name = inputProduct.Name
	}
	if inputProduct.Price != 0 {
		product.Price = inputProduct.Price
	}

	err = r.db.Save(&product).Error
	if err != nil {
		return product, err
	}

	return product, nil
}

func (r *repository) Delete(id int) error {
	product := &model.Product{ID: id}
	tx := r.db.Delete(product)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
