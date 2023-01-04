package payment

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"go-project/model"
)

type Repository interface {
	Create(productName string) (model.Payment, error)
	GetAll() ([]model.Payment, error)
	GetById(id int) (model.Payment, error)
	Update(id int, inputPayment InputPayment) (model.Payment, error)
	Delete(id int) error
	Stream() (<-chan model.Payment, <-chan error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db}
}

func (r *repository) Create(productName string) (model.Payment, error) {
	var product model.Product
	r.db.Where("name = ?", productName).First(&product)

	var payment model.Payment
	payment.PricePaid = product.Price
	payment.ProductId = product.ID
	payment.Product = product

	err := r.db.Create(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) GetAll() ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Find(&payments).Error

	for index, payment := range payments {
		r.db.Where(&model.Product{ID: payment.ProductId}).First(&payments[index].Product)
	}

	if err != nil {
		return payments, err
	}

	return payments, nil
}

func (r *repository) GetById(id int) (model.Payment, error) {
	var payment model.Payment

	err := r.db.Where(&model.Payment{ID: id}).First(&payment).Error

	r.db.Where(&model.Product{ID: payment.ProductId}).First(&payment.Product)
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Update(id int, inputPayment InputPayment) (model.Payment, error) {
	var product model.Product
	r.db.Where("name = ?", inputPayment.ProductName).First(&product)

	payment, err := r.GetById(id)
	if err != nil {
		return payment, err
	}

	if payment.ProductId == product.ID {
		return payment, nil
	}

	payment.PricePaid = product.Price
	payment.ID = product.ID
	payment.Product = product

	err = r.db.Save(&payment).Error
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (r *repository) Delete(id int) error {
	payment := &model.Payment{ID: id}
	tx := r.db.Delete(payment)
	if tx.Error != nil {
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		return errors.New("Payment not found")
	}

	return nil
}

func (r *repository) Stream() (<-chan model.Payment, <-chan error) {
	payments := make(chan model.Payment)
	errc := make(chan error)

	go func() {
		defer close(payments)
		rows, err := r.db.Model(&model.Payment{}).Rows()
		if err != nil {
			fmt.Println("HERE : ")
			errc <- err
			return
		}
		defer rows.Close()

		for rows.Next() {
			var payment model.Payment
			err := r.db.ScanRows(rows, &payment)
			if err != nil {
				fmt.Println("HERE : ")
				errc <- err
				return
			}
			payments <- payment
			errc <- nil
		}
	}()

	return payments, errc
}
