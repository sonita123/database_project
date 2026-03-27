package handlers

import (
	"net/http"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// ProductDetailHandler — /portal/product/{id}
// Public: anyone can view. Login required to add to cart or review.
func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := mux.Vars(r)["id"]

	product, err := repository.GetProductByID(idStr)
	if err != nil {
		http.Error(w, "Product not found", 404)
		return
	}

	// Fetch all reviews for this product
	reviews, _ := repository.GetProductReviews(product.ProductID)
	if reviews == nil {
		reviews = []models.Review{}
	}

	// Track view if logged in
	userID := session.GetUserID(r)
	if userID != 0 {
		repository.RecordProductView(userID, product.ProductID)
	}

	// Can the current user leave a review?
	canReview := false
	orderItemID := 0
	alreadyReviewed := false
	if userID != 0 {
		alreadyReviewed = repository.HasReviewedProduct(userID, product.ProductID)
		if !alreadyReviewed {
			if oid, ok := repository.GetReviewableOrderItem(userID, product.ProductID); ok {
				canReview = true
				orderItemID = oid
			}
		}
	}

	// Average rating as float for display, int for star comparison
	var avgRating float64
	avgRatingInt := 0
	if len(reviews) > 0 {
		total := 0
		for _, rev := range reviews {
			total += rev.Rating
		}
		avgRating = float64(total) / float64(len(reviews))
		avgRatingInt = int(avgRating + 0.5) // round to nearest
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":           product.Name,
		"Product":         product,
		"Reviews":         reviews,
		"AvgRating":       avgRating,
		"AvgRatingInt":    avgRatingInt,
		"ReviewCount":     len(reviews),
		"CanReview":       canReview,
		"AlreadyReviewed": alreadyReviewed,
		"OrderItemID":     orderItemID,
		"IsLoggedIn":      userID != 0,
	})

	if err := Templates["product_detail"].ExecuteTemplate(w, "portal_layout", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
