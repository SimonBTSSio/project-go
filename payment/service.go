package payment

import (
	model "go-project/model"
)

type Service interface {
	Create(input InputPayment) (model.Payment, error)
	GetAll() ([]model.Payment, error)
	GetById(id int) (model.Payment, error)
	Update(id int, inputProduct InputPayment) (model.Payment, error)
	Delete(id int) error
	Stream() (<-chan model.Payment, <-chan error)
}

type service struct {
	repository Repository
}

func NewService(r Repository) *service {
	return &service{r}
}

func (s *service) Create(input InputPayment) (model.Payment, error) {
	var payment model.Payment
	var productName string = input.ProductName

	newPayment, err := s.repository.Create(productName)
	if err != nil {
		return payment, err
	}

	return newPayment, nil
}

func (s *service) GetAll() ([]model.Payment, error) {
	payments, err := s.repository.GetAll()
	if err != nil {
		return payments, err
	}

	return payments, nil
}

func (s *service) GetById(id int) (model.Payment, error) {
	payment, err := s.repository.GetById(id)
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (s *service) Update(id int, input InputPayment) (model.Payment, error) {
	payment, err := s.repository.Update(id, input)
	if err != nil {
		return payment, err
	}

	return payment, nil
}

func (s *service) Delete(id int) error {
	err := s.repository.Delete(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) Stream() (<-chan model.Payment, <-chan error) { // Stream all payments from the database
	payments, errs := s.repository.Stream()
	return payments, errs
}
