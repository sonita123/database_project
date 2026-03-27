package router

import (
	"net/http"

	"unibazar/project/internal/handlers"
	"unibazar/project/internal/session"

	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {

	r := mux.NewRouter()

	// ── Static files ─────────────────────────────────────────
	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("./web/static"))),
	)

	// ── Public auth routes ───────────────────────────────────
	r.HandleFunc("/admin/login", handlers.AdminLoginHandler).Methods("GET", "POST")
	r.HandleFunc("/admin/register", handlers.AdminRegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/register", handlers.UserRegisterHandler).Methods("GET", "POST")
	r.HandleFunc("/login", handlers.UserLoginHandler).Methods("GET", "POST")
	r.HandleFunc("/supporter/login", handlers.SupporterLoginHandler).Methods("GET", "POST")
	r.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET", "POST")

	// ── Shop — PUBLIC (no login required to browse) ───────────
	r.HandleFunc("/portal/shop", handlers.ShopHandler).Methods("GET")
	r.HandleFunc("/portal/product/{id}", handlers.ProductDetailHandler).Methods("GET")

	// ── User portal ──────────────────────────────────────────
	r.HandleFunc("/portal", session.RequireUser(handlers.UserPortalHandler)).Methods("GET")
	r.HandleFunc("/portal/balance/add", session.RequireUser(handlers.AddBalanceHandler)).Methods("POST")
	r.HandleFunc("/portal/cart", session.RequireUser(handlers.CartHandler)).Methods("GET")
	r.HandleFunc("/portal/cart/add", handlers.AddToCartHandler).Methods("POST") // self-redirects to login
	r.HandleFunc("/portal/cart/update", session.RequireUser(handlers.UpdateCartItemHandler)).Methods("POST")
	r.HandleFunc("/portal/cart/remove", session.RequireUser(handlers.RemoveCartItemHandler)).Methods("POST")
	r.HandleFunc("/portal/checkout", session.RequireUser(handlers.CheckoutHandler)).Methods("GET", "POST")
	r.HandleFunc("/portal/orders", session.RequireUser(handlers.UserOrdersHandler)).Methods("GET")

	// User: addresses
	r.HandleFunc("/portal/addresses", session.RequireUser(handlers.AddressesHandler)).Methods("GET")
	r.HandleFunc("/portal/addresses/add", session.RequireUser(handlers.AddAddressHandler)).Methods("GET", "POST")
	r.HandleFunc("/portal/addresses/{id}/edit", session.RequireUser(handlers.EditAddressHandler)).Methods("GET", "POST")
	r.HandleFunc("/portal/addresses/{id}/delete", session.RequireUser(handlers.DeleteAddressHandler)).Methods("POST")
	r.HandleFunc("/portal/addresses/{id}/default", session.RequireUser(handlers.SetDefaultAddressHandler)).Methods("POST")

	// User: reviews
	r.HandleFunc("/portal/reviews", session.RequireUser(handlers.UserReviewsHandler)).Methods("GET")
	r.HandleFunc("/portal/reviews/submit", session.RequireUser(handlers.SubmitReviewHandler)).Methods("POST")

	// User: seller & VIP requests (any logged-in user can request)
	r.HandleFunc("/portal/become-seller", session.RequireUser(handlers.BecomeSellerPageHandler)).Methods("GET", "POST")
	r.HandleFunc("/portal/request-vip", session.RequireUser(handlers.RequestVIPHandler)).Methods("GET", "POST")
	r.HandleFunc("/portal/stall/create", session.RequireUser(handlers.CreateStallPageHandler)).Methods("GET", "POST")

	// ── Seller portal (sellers only) ─────────────────────────
	r.HandleFunc("/seller/dashboard", session.RequireUser(handlers.SellerDashboardHandler)).Methods("GET")
	r.HandleFunc("/seller/stalls/{stall_id}/products", session.RequireUser(handlers.SellerProductsHandler)).Methods("GET")
	r.HandleFunc("/seller/stalls/{stall_id}/products/add", session.RequireUser(handlers.SellerAddProductHandler)).Methods("GET", "POST")
	r.HandleFunc("/seller/stalls/{stall_id}/products/{product_id}/edit", session.RequireUser(handlers.SellerEditProductHandler)).Methods("GET", "POST")
	r.HandleFunc("/seller/stalls/{stall_id}/products/{product_id}/delete", session.RequireUser(handlers.SellerDeleteProductHandler)).Methods("POST")
	r.HandleFunc("/seller/reviews", session.RequireUser(handlers.SellerReviewsHandler)).Methods("GET")

	// ── Supporter portal ─────────────────────────────────────
	r.HandleFunc("/supporter/portal", session.RequireSupporter(handlers.SupporterPortalHandler)).Methods("GET")
	r.HandleFunc("/supporter/portal/orders", session.RequireSupporter(handlers.SupporterOrdersHandler)).Methods("GET")
	r.HandleFunc("/supporter/portal/requests", session.RequireSupporter(handlers.SupporterRequestsHandler)).Methods("GET")
	r.HandleFunc("/supporter/portal/stalls", session.RequireSupporter(handlers.StallsHandlerForSupporter)).Methods("GET")

	// Supporter actions
	r.HandleFunc("/supporter/orders/{id}/status", session.RequireSupporter(handlers.SupporterUpdateOrderHandler)).Methods("POST")
	r.HandleFunc("/supporter/portal/requests/{id}/handle", session.RequireSupporter(handlers.SupporterHandleRequestHandler)).Methods("POST")
	r.HandleFunc("/supporter/portal/stalls/{id}/action", session.RequireSupporter(handlers.SupporterApproveStallHandler)).Methods("POST")

	// ── Admin dashboard ──────────────────────────────────────
	r.HandleFunc("/", session.RequireAdmin(handlers.Dashboard)).Methods("GET")

	// ── Users ────────────────────────────────────────────────
	r.HandleFunc("/users",
		session.RequireAdmin(handlers.UsersPageHandler)).Methods("GET")
	r.HandleFunc("/users/create",
		session.RequireAdmin(handlers.CreateUserHandler)).Methods("GET", "POST")
	r.HandleFunc("/users/edit/{id}",
		session.RequireAdmin(handlers.EditUserHandler)).Methods("GET", "POST")
	r.HandleFunc("/users/delete/{id}",
		session.RequireAdmin(handlers.DeleteUserHandler)).Methods("GET")
	r.HandleFunc("/users/{id}/set-password",
		session.RequireAdmin(handlers.SetUserPasswordHandler)).Methods("GET", "POST")
	r.HandleFunc("/users/{id}",
		session.RequireAdmin(handlers.UserPageHandler)).Methods("GET")

	// ── Supporters ───────────────────────────────────────────
	r.HandleFunc("/supporters",
		session.RequireAdmin(handlers.SupportersPageHandler)).Methods("GET")
	r.HandleFunc("/supporters/create",
		session.RequireAdmin(handlers.CreateSupporterHandler)).Methods("GET", "POST")
	r.HandleFunc("/supporters/{id}/set-password",
		session.RequireAdmin(handlers.SetSupporterPasswordHandler)).Methods("GET", "POST")
	r.HandleFunc("/supporters/{id}/edit",
		session.RequireAdmin(handlers.EditSupporterHandler)).Methods("GET", "POST")
	r.HandleFunc("/supporters/{id}/delete",
		session.RequireAdmin(handlers.DeleteSupporterHandler)).Methods("POST")
	r.HandleFunc("/supporters/{id}",
		session.RequireAdmin(handlers.ViewSupporterHandler)).Methods("GET")

	// ── Settings ─────────────────────────────────────────────
	r.HandleFunc("/settings",
		session.RequireAdmin(handlers.SettingsHandler)).Methods("GET")

	// Products
	r.HandleFunc("/settings/products",
		session.RequireAdmin(handlers.ProductsHandler)).Methods("GET")
	r.HandleFunc("/settings/products/create",
		session.RequireAdmin(handlers.CreateProductHandler)).Methods("GET", "POST")
	r.HandleFunc("/settings/products/{id}/edit",
		session.RequireAdmin(handlers.EditProductHandler)).Methods("GET", "POST")
	r.HandleFunc("/settings/products/{id}/delete",
		session.RequireAdmin(handlers.DeleteProductHandler)).Methods("POST")

	// Stalls
	r.HandleFunc("/settings/stalls",
		session.RequireAdmin(handlers.StallsHandler)).Methods("GET")
	r.HandleFunc("/settings/stalls/{id}/status",
		session.RequireAdmin(handlers.UpdateStallStatusHandler)).Methods("POST")
	r.HandleFunc("/settings/stalls/{id}/delete",
		session.RequireAdmin(handlers.DeleteStallHandler)).Methods("POST")

	// Discounts
	r.HandleFunc("/settings/discounts",
		session.RequireAdmin(handlers.DiscountCodesHandler)).Methods("GET")
	r.HandleFunc("/settings/discounts/create",
		session.RequireAdmin(handlers.CreateDiscountCodeHandler)).Methods("GET", "POST")
	r.HandleFunc("/settings/discounts/{id}/toggle",
		session.RequireAdmin(handlers.ToggleDiscountCodeHandler)).Methods("POST")
	r.HandleFunc("/settings/discounts/{id}/delete",
		session.RequireAdmin(handlers.DeleteDiscountCodeHandler)).Methods("POST")

	// Orders (admin view)
	r.HandleFunc("/settings/orders",
		session.RequireAdmin(handlers.OrdersHandler)).Methods("GET")

	// Reviews
	r.HandleFunc("/settings/reviews",
		session.RequireAdmin(handlers.ReviewsHandler)).Methods("GET")
	r.HandleFunc("/settings/reviews/{id}/delete",
		session.RequireAdmin(handlers.DeleteReviewHandler)).Methods("POST")

	// Fraud Reports
	r.HandleFunc("/settings/fraud-reports",
		session.RequireAdmin(handlers.FraudReportsHandler)).Methods("GET")
	r.HandleFunc("/settings/fraud-reports/{id}/delete",
		session.RequireAdmin(handlers.DeleteFraudReportHandler)).Methods("POST")

	return r
}
