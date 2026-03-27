package repository

import (
	"context"
	"database/sql"
	"unibazar/project/internal/db"
	"unibazar/project/internal/models"
)

// GetSupporterKPIs fetches all 5 performance indicators for a supporter
func GetSupporterKPIs(supporterID string) (models.SupporterKPI, error) {
	var kpi models.SupporterKPI

	err := db.Conn.QueryRowContext(context.Background(), `
		WITH approved_stalls AS (
			SELECT stall_id, approved_at 
			FROM stalls 
			WHERE approved_by = @supporter_id AND approved_at IS NOT NULL
		),
		second_fraud AS (
			SELECT fr.stall_id, s.approved_at, MIN(fr.reported_at) AS second_report_at
			FROM (
				SELECT stall_id, reported_at, 
					   ROW_NUMBER() OVER (PARTITION BY stall_id ORDER BY reported_at) rn 
				FROM fraud_reports
			) fr
			JOIN approved_stalls s ON s.stall_id = fr.stall_id 
			WHERE fr.rn = 2
			GROUP BY fr.stall_id, s.approved_at
		),
		suspension_stats AS (
			SELECT 
				SUM(CASE WHEN st.status = 'suspended' THEN 1 ELSE 0 END) suspended_count,
			    COUNT(*) total_count 
			FROM approved_stalls a 
			JOIN stalls st ON st.stall_id = a.stall_id
		),
		stall_sales AS (
			SELECT a.stall_id, ISNULL(SUM(oi.quantity * oi.price), 0) total_revenue
			FROM approved_stalls a 
			LEFT JOIN products p ON p.stall_id = a.stall_id 
			LEFT JOIN order_items oi ON oi.product_id = p.product_id 
			GROUP BY a.stall_id
		),
		stall_decile_stats AS (
			SELECT total_revenue, NTILE(10) OVER (ORDER BY total_revenue) decile 
			FROM stall_sales
		),
		discount_usage_amounts AS (
			SELECT du.user_id, 
				   SUM(CASE 
						WHEN dc.discount_type = 'fixed' THEN ISNULL(dc.fixed_amount, 0)
						WHEN dc.discount_type = 'percentage' THEN ISNULL(dc.percentage, 0) 
						ELSE 0 
				   END) user_total
			FROM discount_usage du 
			JOIN discount_codes dc ON dc.discount_id = du.discount_id 
			WHERE dc.supporter_id = @supporter_id 
			GROUP BY du.user_id
		),
		stats AS (
			SELECT 
				ISNULL(AVG(DATEDIFF(HOUR, sf.approved_at, sf.second_report_at)), 0) avg_hours,
				ss.suspended_count * 100.0 / NULLIF(ss.total_count, 0) suspended_pct,
				100.0 * COUNT(*) / NULLIF((SELECT COUNT(*) FROM stall_sales), 0) bottom_decile_pct,
				ISNULL(MAX(d.user_total), 0) * 100.0 / NULLIF(SUM(d.user_total), 0) top_user_share
			FROM second_fraud sf 
			CROSS JOIN suspension_stats ss 
			CROSS JOIN stall_decile_stats sds 
			CROSS JOIN discount_usage_amounts d
			WHERE sds.decile = 1
		)
		SELECT avg_hours, suspended_pct, bottom_decile_pct, top_user_share 
		FROM stats
	`,
		sql.Named("supporter_id", supporterID),
	).Scan(
		&kpi.AvgHoursToSecondFraudBadge,
		&kpi.SuspendedStallsPercent,
		&kpi.BottomDecilePercent,
		&kpi.TopUserDiscountShare,
	)

	if err != nil {
		kpi = models.SupporterKPI{}
	}

	// KPI 4: Operations last 7 days (FIXED QUERY)
	var operations int
	err = db.Conn.QueryRowContext(context.Background(), `
		SELECT 
			(SELECT COUNT(*) 
			 FROM stalls 
			 WHERE approved_by = @supporter_id 
			   AND approved_at >= DATEADD(DAY, -7, GETDATE()))
			+
			(SELECT COUNT(*) 
			 FROM discount_codes 
			 WHERE supporter_id = @supporter_id 
			   AND created_at >= DATEADD(DAY, -7, GETDATE()))
	`,
		sql.Named("supporter_id", supporterID),
	).Scan(&operations)

	if err != nil {
		kpi.OperationsLast7Days = 0
	} else {
		kpi.OperationsLast7Days = operations
	}

	return kpi, nil
}
