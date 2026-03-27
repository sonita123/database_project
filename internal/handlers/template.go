package handlers

import (
	"html/template"
	"path/filepath"
)

var Templates map[string]*template.Template

var funcMap = template.FuncMap{
	"add":  func(a, b int) int { return a + b },
	"sub":  func(a, b int) int { return a - b },
	"subf": func(a, b float64) float64 { return a - b },
	"max": func(a, b int) int {
		if a > b {
			return a
		}
		return b
	},
	"min": func(a, b int) int {
		if a < b {
			return a
		}
		return b
	},
	"mul": func(a float64, b int) float64 { return a * float64(b) },
	"int": func(f float64) int { return int(f) },
	"seq": func(start, end int) []int {
		s := make([]int, end-start+1)
		for i := range s {
			s[i] = start + i
		}
		return s
	},
	"deref": func(f *float64) float64 {
		if f == nil {
			return 0
		}
		return *f
	},
	"nearbyPages": func(current, total int) []int {
		pages := []int{}
		var start, end int
		if current-2 > 1 {
			start = current - 2
		} else {
			start = 1
		}
		if current+2 < total {
			end = current + 2
		} else {
			end = total
		}
		for i := start; i <= end; i++ {
			pages = append(pages, i)
		}
		return pages
	},
}

func LoadTemplates(root string) {
	Templates = make(map[string]*template.Template)

	layout := filepath.Join(root, "web", "templates", "layouts", "layout.gohtml")
	authLayout := filepath.Join(root, "web", "templates", "layouts", "auth_layout.gohtml")
	portalLayout := filepath.Join(root, "web", "templates", "layouts", "portal_layout.gohtml")

	// ── Main admin pages (sidebar layout) ────────────────────
	pages := map[string]string{
		"dashboard":                filepath.Join(root, "web", "templates", "pages", "dashboard.gohtml"),
		"users":                    filepath.Join(root, "web", "templates", "pages", "users.gohtml"),
		"create_user":              filepath.Join(root, "web", "templates", "pages", "create_user.gohtml"),
		"edit_user":                filepath.Join(root, "web", "templates", "pages", "edit_user.gohtml"),
		"supporters":               filepath.Join(root, "web", "templates", "pages", "supporters.gohtml"),
		"user":                     filepath.Join(root, "web", "templates", "pages", "user.gohtml"),
		"view_supporter":           filepath.Join(root, "web", "templates", "pages", "supporter.gohtml"),
		"create_supporter":         filepath.Join(root, "web", "templates", "pages", "create_supporter.gohtml"),
		"edit_supporter":           filepath.Join(root, "web", "templates", "pages", "update_supporter.gohtml"),
		"settings_products":        filepath.Join(root, "web", "templates", "pages", "setting_product.gohtml"),
		"settings_create_product":  filepath.Join(root, "web", "templates", "pages", "setting_product_create.gohtml"),
		"settings_edit_product":    filepath.Join(root, "web", "templates", "pages", "setting_product_edit.gohtml"),
		"settings_stalls":          filepath.Join(root, "web", "templates", "pages", "stall.gohtml"),
		"settings":                 filepath.Join(root, "web", "templates", "pages", "setting.gohtml"),
		"settings_discounts":       filepath.Join(root, "web", "templates", "pages", "discounts.gohtml"),
		"settings_create_discount": filepath.Join(root, "web", "templates", "pages", "create_discount.gohtml"),
		"settings_fraud":           filepath.Join(root, "web", "templates", "pages", "fraud.gohtml"),
		"set_user_password":        filepath.Join(root, "web", "templates", "pages", "set_password_user.gohtml"),
		"set_supporter_password":   filepath.Join(root, "web", "templates", "pages", "set_password_supporter.gohtml"),
		// Both keys point to the same file
		"admin_orders":    filepath.Join(root, "web", "templates", "pages", "orders.gohtml"),
		"settings_orders": filepath.Join(root, "web", "templates", "pages", "orders.gohtml"),
		// Reviews — matches your filename setting_reviews.gohtml
		"settings_reviews": filepath.Join(root, "web", "templates", "pages", "setting_reviews.gohtml"),
	}

	for name, page := range pages {
		Templates[name] = template.Must(
			template.New("").Funcs(funcMap).ParseFiles(layout, page),
		)
	}

	// ── Auth pages (no sidebar) ───────────────────────────────
	authPages := map[string]string{
		"login_user":      filepath.Join(root, "web", "templates", "pages", "login_user.gohtml"),
		"login_supporter": filepath.Join(root, "web", "templates", "pages", "login_supporter.gohtml"),
		"login_admin":     filepath.Join(root, "web", "templates", "pages", "admin_login.gohtml"),
		"register_admin":  filepath.Join(root, "web", "templates", "pages", "admin_register.gohtml"),
		"register_user":   filepath.Join(root, "web", "templates", "pages", "register_user.gohtml"),
	}

	for name, page := range authPages {
		Templates[name] = template.Must(
			template.New("").Funcs(funcMap).ParseFiles(authLayout, page),
		)
	}

	// ── Portal pages (user / seller / supporter) ──────────────
	portalPages := map[string]string{

		// ── Buyer / user
		"user_portal":    filepath.Join(root, "web", "templates", "pages", "user_portal.gohtml"),
		"shop":           filepath.Join(root, "web", "templates", "pages", "shop.gohtml"),
		"cart":           filepath.Join(root, "web", "templates", "pages", "cart.gohtml"),
		"checkout":       filepath.Join(root, "web", "templates", "pages", "checkout.gohtml"),
		"orders":         filepath.Join(root, "web", "templates", "pages", "user_orders.gohtml"),
		"become_seller":  filepath.Join(root, "web", "templates", "pages", "become_seller.gohtml"),
		"request_vip":    filepath.Join(root, "web", "templates", "pages", "request_vip.gohtml"),
		"create_stall":   filepath.Join(root, "web", "templates", "pages", "create_stall.gohtml"),
		"product_detail": filepath.Join(root, "web", "templates", "pages", "product_detail.gohtml"),
		// matches your filename: address.gohtml
		"addresses":    filepath.Join(root, "web", "templates", "pages", "address.gohtml"),
		"address_form": filepath.Join(root, "web", "templates", "pages", "address_form.gohtml"),
		// matches your filename: user_review.gohtml
		"user_reviews": filepath.Join(root, "web", "templates", "pages", "user_review.gohtml"),

		// ── Seller
		"seller_dashboard":    filepath.Join(root, "web", "templates", "pages", "seller_dashboard.gohtml"),
		"seller_products":     filepath.Join(root, "web", "templates", "pages", "seller_products.gohtml"),
		"seller_product_form": filepath.Join(root, "web", "templates", "pages", "seller_product_form.gohtml"),
		"seller_reviews":      filepath.Join(root, "web", "templates", "pages", "seller_reviews.gohtml"),

		// ── Supporter
		"supporter_portal": filepath.Join(root, "web", "templates", "pages", "supporter_portal.gohtml"),
		"supporter_orders": filepath.Join(root, "web", "templates", "pages", "support_order.gohtml"),
		// matches your filename: supporter_request.gohtml
		"supporter_requests": filepath.Join(root, "web", "templates", "pages", "supporter_request.gohtml"),
		"supporter_stalls":   filepath.Join(root, "web", "templates", "pages", "supporter_stalls.gohtml"),
	}

	for name, page := range portalPages {
		Templates[name] = template.Must(
			template.New("portal_layout").Funcs(funcMap).ParseFiles(portalLayout, page),
		)
	}
}
