package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/models"
	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

// ── User: view and write reviews ──────────────────────────────

func UserReviewsHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	reviewable, _ := repository.GetReviewableProducts(userID)
	written, _ := repository.GetUserReviews(userID)

	if reviewable == nil {
		reviewable = []models.ReviewableItem{}
	}
	if written == nil {
		written = []models.Review{}
	}

	data := mergeBase(PortalBase(r), map[string]any{
		"Title":      "My Reviews",
		"Reviewable": reviewable,
		"Written":    written,
	})
	Templates["user_reviews"].ExecuteTemplate(w, "portal_layout", data)
}

func SubmitReviewHandler(w http.ResponseWriter, r *http.Request) {
	userID := session.GetUserID(r)

	productID, err := strconv.Atoi(r.FormValue("product_id"))
	if err != nil {
		http.Error(w, "invalid product", 400)
		return
	}
	orderItemID, err := strconv.Atoi(r.FormValue("order_item_id"))
	if err != nil {
		http.Error(w, "invalid order item", 400)
		return
	}
	rating, err := strconv.Atoi(r.FormValue("rating"))
	if err != nil || rating < 1 || rating > 5 {
		http.Error(w, "rating must be 1–5", 400)
		return
	}
	comment := r.FormValue("comment")

	if err := repository.CreateReview(userID, productID, orderItemID, rating, comment); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	http.Redirect(w, r, "/portal/reviews", http.StatusSeeOther)
}

// ── Admin: manage all reviews ─────────────────────────────────

func ReviewsHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	reviews, err := repository.GetAllReviewsPaginated(settingsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total, _ := repository.CountAllReviews()
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	data := mergeBase(PortalBase(r), map[string]any{
		"Title": "Reviews", "CurrentPath": r.URL.Path,
		"Reviews": reviews, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
		"AdminUsername": AdminUsername(r),
	})
	Templates["settings_reviews"].ExecuteTemplate(w, "layout", data)
}

func DeleteReviewHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteReview(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/reviews", http.StatusSeeOther)
}
