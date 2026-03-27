package models

import "time"

// ── existing models stay in models.go ──────────────────────
// These are the additional types needed for the missing features

type SellerTopProduct struct {
	ProductID   int
	ProductName string
	StallName   string
	TotalSold   int
	Revenue     float64
}
type SellerStats struct {
	TotalRevenue   float64
	TotalOrders    int
	TotalItemsSold int
	TotalProducts  int
	ActiveProducts int
	TotalStalls    int
	ActiveStalls   int
}

// SupporterKPI holds all 5 performance indicators for a supporter
type SupporterKPI struct {
	// KPI 1: avg hours between stall approval and 2nd fraud badge
	AvgHoursToSecondFraudBadge float64

	// KPI 2: % of approved stalls that ended up fully suspended
	SuspendedStallsPercent float64

	// KPI 3: % of approved stalls in the bottom sales decile
	BottomDecilePercent float64

	// KPI 4: total operations (stall approvals + discount codes created) in last 7 days
	OperationsLast7Days int

	// KPI 5: top user's share % of all discounts created by this supporter
	TopUserDiscountShare float64
}

// DiscountCodeUsage combines a discount code with whether the user has used it
type DiscountCodeUsage struct {
	DiscountID     int
	Code           string
	DiscountType   string
	Percentage     *float64
	FixedAmount    *float64
	ExpirationDate time.Time
	IsActive       bool
	UsedByUser     bool
}
type Admin struct {
	AdminID   int
	Username  string
	Password  string
	CreatedAt time.Time
}

// RecommendedProduct is a product recommended via the similarity algorithm
type RecommendedProduct struct {
	ProductID   int
	Name        string
	Price       float64
	StallID     int
	BuyersCount int // how many similar users bought this
}
