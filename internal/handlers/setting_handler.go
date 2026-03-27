package handlers

import (
	"net/http"
	"strconv"
	"unibazar/project/internal/repository"

	"github.com/gorilla/mux"
)

const settingsPerPage = 10

func settingsParsePage(r *http.Request) (int, int) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}
	return page, (page - 1) * settingsPerPage
}

func settingsBuildPages(total, perPage int) (int, []int) {
	totalPages := total / perPage
	if total%perPage != 0 {
		totalPages++
	}
	if totalPages < 1 {
		totalPages = 1
	}
	pages := make([]int, totalPages)
	for i := range pages {
		pages[i] = i + 1
	}
	return totalPages, pages
}

func settingsRender(w http.ResponseWriter, key string, data map[string]any) {
	if err := Templates[key].ExecuteTemplate(w, "layout", data); err != nil {
		http.Error(w, err.Error(), 500)
	}
}

// ── Settings overview ─────────────────────────────────────────

func SettingsHandler(w http.ResponseWriter, r *http.Request) {
	totalProducts, _ := repository.CountProducts()
	totalStalls, _ := repository.CountStalls()
	totalOrders, _ := repository.CountOrders()
	totalDiscounts, _ := repository.CountDiscountCodes()
	totalReports, _ := repository.CountFraudReports()
	totalReviews, _ := repository.CountReviews()

	settingsRender(w, "settings", map[string]any{
		"Title":          "Settings",
		"CurrentPath":    r.URL.Path,
		"TotalProducts":  totalProducts,
		"TotalStalls":    totalStalls,
		"TotalOrders":    totalOrders,
		"TotalDiscounts": totalDiscounts,
		"TotalReports":   totalReports,
		"TotalReviews":   totalReviews,
	})
}

// ── Products ─────────────────────────────────────────────────

func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	products, err := repository.GetProductsPaginated(settingsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total, _ := repository.CountProducts()
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	settingsRender(w, "settings_products", map[string]any{
		"Title": "Products", "CurrentPath": r.URL.Path,
		"Products": products, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
	})
}

func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		stalls, _ := repository.GetAllStalls()
		settingsRender(w, "settings_create_product", map[string]any{
			"Title": "Add Product", "CurrentPath": "/settings/products", "Stalls": stalls,
		})
		return
	}
	if err := repository.CreateProduct(
		r.FormValue("stall_id"), r.FormValue("name"),
		r.FormValue("price"), r.FormValue("stock"), r.FormValue("status"),
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/products", http.StatusSeeOther)
}

func EditProductHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if r.Method == "GET" {
		product, err := repository.GetProductByID(id)
		if err != nil {
			http.Error(w, "Product not found", 404)
			return
		}
		stalls, _ := repository.GetAllStalls()
		settingsRender(w, "settings_edit_product", map[string]any{
			"Title": "Edit Product", "CurrentPath": "/settings/products",
			"Product": product, "Stalls": stalls,
		})
		return
	}
	if err := repository.UpdateProduct(id,
		r.FormValue("name"), r.FormValue("price"),
		r.FormValue("stock"), r.FormValue("status"),
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/products", http.StatusSeeOther)
}

func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteProduct(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/products", http.StatusSeeOther)
}

// ── Stalls ───────────────────────────────────────────────────

func StallsHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	stalls, err := repository.GetStallsPaginated(settingsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total, _ := repository.CountStalls()
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	settingsRender(w, "settings_stalls", map[string]any{
		"Title": "Stalls", "CurrentPath": r.URL.Path,
		"Stalls": stalls, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
	})
}

func UpdateStallStatusHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.UpdateStallStatus(mux.Vars(r)["id"], r.FormValue("status")); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/stalls", http.StatusSeeOther)
}

func DeleteStallHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteStall(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/stalls", http.StatusSeeOther)
}

// ── Discount Codes ───────────────────────────────────────────

func DiscountCodesHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	codes, err := repository.GetDiscountCodesPaginated(settingsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total, _ := repository.CountDiscountCodes()
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	settingsRender(w, "settings_discounts", map[string]any{
		"Title": "Discount Codes", "CurrentPath": r.URL.Path,
		"Codes": codes, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
	})
}

func CreateDiscountCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		supporters, _ := repository.GetAllSupportersSimple()
		settingsRender(w, "settings_create_discount", map[string]any{
			"Title": "Add Discount Code", "CurrentPath": "/settings/discounts",
			"Supporters": supporters,
		})
		return
	}
	if err := repository.CreateDiscountCode(
		r.FormValue("code"), r.FormValue("discount_type"),
		r.FormValue("percentage"), r.FormValue("fixed_amount"),
		r.FormValue("expiration_date"), r.FormValue("max_uses"),
		r.FormValue("supporter_id"),
	); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/discounts", http.StatusSeeOther)
}

func ToggleDiscountCodeHandler(w http.ResponseWriter, r *http.Request) {
	isActive := r.FormValue("action") == "enable"
	if err := repository.ToggleDiscountCode(mux.Vars(r)["id"], isActive); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/discounts", http.StatusSeeOther)
}

func DeleteDiscountCodeHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteDiscountCode(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/discounts", http.StatusSeeOther)
}

// ── Fraud Reports ────────────────────────────────────────────

func FraudReportsHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	reports, err := repository.GetFraudReportsPaginated(settingsPerPage, offset)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total, _ := repository.CountFraudReports()
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	settingsRender(w, "settings_fraud", map[string]any{
		"Title": "Fraud Reports", "CurrentPath": r.URL.Path,
		"Reports": reports, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
	})
}

func DeleteFraudReportHandler(w http.ResponseWriter, r *http.Request) {
	if err := repository.DeleteFraudReport(mux.Vars(r)["id"]); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	http.Redirect(w, r, "/settings/fraud-reports", http.StatusSeeOther)
}

// ── Orders (admin view) ───────────────────────────────────────

func OrdersHandler(w http.ResponseWriter, r *http.Request) {
	page, offset := settingsParsePage(r)
	orders, err := repository.GetAllOrdersForSupporter()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	total := len(orders)
	totalPages, pages := settingsBuildPages(total, settingsPerPage)

	// paginate in memory
	start := offset
	if start > total {
		start = total
	}
	end := start + settingsPerPage
	if end > total {
		end = total
	}
	paged := orders[start:end]

	settingsRender(w, "settings_orders", map[string]any{
		"Title": "Orders", "CurrentPath": r.URL.Path,
		"Orders": paged, "Total": total,
		"CurrentPage": page, "TotalPages": totalPages,
		"PrevPage": page - 1, "NextPage": page + 1, "Pages": pages,
		"AdminUsername": AdminUsername(r),
	})
}
