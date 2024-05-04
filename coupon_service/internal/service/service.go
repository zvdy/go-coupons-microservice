package service

import (
	// changed dot import https://go.dev/wiki/CodeReviewComments#import-dot and added entity. prefix
	"coupon_service/internal/service/entity"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

// Repository defines methods to interact with coupons storage
type Repository interface {
	// FindByCode returns a coupon with given code or error if not found
	FindByCode(code string) (*entity.Coupon, error)

	// Save saves a coupon to storage, returns error if failed
	Save(coupon entity.Coupon) error
}

// Service struct is the main service layer for the coupon service
type Service struct {
	repo Repository // storage for coupons
}

// New creates a new service layer with given repository
func New(repository Repository) Service {
	return Service{
		repo: repository, // storage for coupons
	}
}

// ApplyCoupon applies a coupon to a basket if the basket value is greater than zero
// Returns an error if the basket value is less than or equal to zero
func (s Service) ApplyCoupon(basket entity.Basket, code string) (b *entity.Basket, e error) {
	b = &basket
	coupon, err := s.repo.FindByCode(code)
	if err != nil {
		return nil, err
	}

	if b.Value > 0 {
		b.AppliedDiscount = coupon.Discount
		b.Value = b.Value * (1 - int(coupon.Discount)/100)
		b.ApplicationSuccessful = true
		return b, nil
	}

	if b.Value <= 0 {
		return nil, fmt.Errorf("tried to apply discount to negative or zero value")
	}

	return nil, fmt.Errorf("unexpected error")
}

// CreateCoupon creates a new coupon with a given discount, code, and minimum basket value
func (s Service) CreateCoupon(discount int, code string, minBasketValue int) (*entity.Coupon, error) {
	coupon := &entity.Coupon{
		Discount:       discount,
		Code:           code,
		MinBasketValue: minBasketValue,
		ID:             uuid.NewString(),
	}

	if err := s.repo.Save(*coupon); err != nil {
		return nil, err
	}
	return coupon, nil
}

// GetCoupons returns a list of coupons for given codes
// If a coupon is not found for a code, it adds an error message to the returned error
func (s Service) GetCoupons(codes []string) ([]entity.Coupon, error) {
	coupons := make([]entity.Coupon, len(codes))
	errors := make([]error, len(codes))

	var wg sync.WaitGroup
	wg.Add(len(codes))

	for idx, code := range codes {
		go func(idx int, code string) {
			defer wg.Done()
			coupon, err := s.repo.FindByCode(code)
			if err != nil {
				errors[idx] = fmt.Errorf("code: %s, index: %d, error: %v", code, idx, err)
			}
			coupons[idx] = *coupon
		}(idx, code)
	}

	wg.Wait()

	var e error = nil
	for _, err := range errors {
		if err != nil {
			if e == nil {
				e = err
			} else {
				e = fmt.Errorf("%w; %v", e, err)
			}
		}
	}

	return coupons, e
}

// FindByCode returns a coupon for a given code
func (s Service) FindByCode(code string) (*entity.Coupon, error) {
	return s.repo.FindByCode(code)
}
