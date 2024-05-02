package api

import (
	"coupon_service/internal/api/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Apply is a handler function that applies a coupon to a basket
func (a *API) Apply(c *gin.Context) {
	apiReq := entity.ApplicationRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		a.httpRequestsTotal.WithLabelValues("400", "/apply").Inc()
		return
	}

	coupon, err := a.svc.FindByCode(apiReq.Code)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Coupon with code " + apiReq.Code + " not found"})
		a.httpRequestsTotal.WithLabelValues("404", "/apply").Inc()
		return
	}

	if apiReq.Basket.Value < coupon.MinBasketValue {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Basket value (%d) inferior to the minimum basket value for this code (%d)", apiReq.Basket.Value, coupon.MinBasketValue)})
		a.httpRequestsTotal.WithLabelValues("400", "/apply").Inc()
		return
	}

	basket, err := a.svc.ApplyCoupon(apiReq.Basket, apiReq.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		a.httpRequestsTotal.WithLabelValues("500", "/apply").Inc()
		return
	}

	a.httpRequestsTotal.WithLabelValues("200", "/apply").Inc()

	c.JSON(http.StatusOK, basket)
}

// Create is a handler function that creates a new coupon
func (a *API) Create(c *gin.Context) {
	apiReq := entity.Coupon{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		a.httpRequestsTotal.WithLabelValues("400", "/create").Inc()
		return
	}

	// Check if a coupon with the same code already exists
	_, err := a.svc.FindByCode(apiReq.Code)
	if err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Coupon with code %s already exists", apiReq.Code)})
		a.httpRequestsTotal.WithLabelValues("400", "/create").Inc()
		return
	}

	coupon, err := a.svc.CreateCoupon(apiReq.Discount, apiReq.Code, apiReq.MinBasketValue)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		a.httpRequestsTotal.WithLabelValues("500", "/create").Inc()
		return
	}

	a.httpRequestsTotal.WithLabelValues("201", "/create").Inc()

	c.JSON(http.StatusCreated, coupon) // Return Created 201 with the coupon details
}

// Get is a handler function that retrieves a list of coupons based on the provided codes
func (a *API) Get(c *gin.Context) {
	apiReq := entity.CouponRequest{}
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		a.httpRequestsTotal.WithLabelValues("400", "/coupons").Inc()
		return
	}

	for _, code := range apiReq.Codes {
		_, err := a.svc.FindByCode(code)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon with code " + code + " not found"})
			a.httpRequestsTotal.WithLabelValues("404", "/coupons").Inc()
			return
		}
	}

	coupons, err := a.svc.GetCoupons(apiReq.Codes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		a.httpRequestsTotal.WithLabelValues("500", "/coupons").Inc()
		return
	}

	a.httpRequestsTotal.WithLabelValues("200", "/coupons").Inc()

	c.JSON(http.StatusOK, coupons)
}
