package handlers

import (
	"net/http"
	"strconv"

	"unibazar/project/internal/repository"
	"unibazar/project/internal/session"
)

// PortalBase returns the common data every portal template needs.
// Lives in handlers (not session) to avoid import cycle:
//
//	handlers → session → repository   ✓
//	session  → handlers               ✗ (was the cycle)
func PortalBase(r *http.Request) map[string]any {
	role := session.GetRole(r)
	userID := session.GetUserID(r)

	isSeller := false
	cartCount := 0
	userType := "regular"

	if role == "user" && userID != 0 {
		// Check sellers table — this is the authoritative source
		isSeller, _ = repository.IsSeller(userID)

		// Also sync user_type from DB for display purposes
		if user, err := repository.GetUserByID(strconv.Itoa(userID)); err == nil {
			userType = user.UserType
		}

		// If in sellers table but user_type not updated yet — fix it now
		if isSeller && userType != "seller" {
			repository.UpdateUserType(userID, "seller")
			userType = "seller"
		}

		// Only buyers get a cart count
		if !isSeller {
			if cart, err := repository.GetOrCreateCart(userID); err == nil {
				if items, err := repository.GetCartItems(cart.CartID); err == nil {
					cartCount = len(items)
				}
			}
		}
	}

	// IsBuyer: logged in as user AND not in the sellers table
	isBuyer := role == "user" && !isSeller

	return map[string]any{
		"IsUser":      role == "user",
		"IsSupporter": role == "supporter",
		"IsSeller":    isSeller,
		"IsBuyer":     isBuyer,
		"UserType":    userType,
		"CartCount":   cartCount,
		"Role":        role,
	}
}

// mergeBase merges PortalBase data with page-specific data
func mergeBase(base, extra map[string]any) map[string]any {
	for k, v := range extra {
		base[k] = v
	}
	return base
}

// isBuyerFromBase safely extracts IsBuyer from a PortalBase map
func isBuyerFromBase(base map[string]any) bool {
	v, ok := base["IsBuyer"].(bool)
	if !ok {
		return false
	}
	return v
}
