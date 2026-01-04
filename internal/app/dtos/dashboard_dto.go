package dtos

import (
	"time"

	"github.com/shopspring/decimal"
)

// DashboardResponse represents the main dashboard data for a user
type DashboardResponse struct {
	// Financial Overview
	TotalAssets      decimal.Decimal `json:"total_assets"`      // Sum of all BANK + CASH + SAVINGS + INVESTMENT accounts
	TotalLiabilities decimal.Decimal `json:"total_liabilities"` // Sum of all CREDIT_CARD accounts
	NetWorth         decimal.Decimal `json:"net_worth"`         // Assets - Liabilities
	LiquidAssets     decimal.Decimal `json:"liquid_assets"`     // BANK + CASH only

	// Monthly Stats (current month)
	MonthlyIncome      decimal.Decimal `json:"monthly_income"`        // Total INCOME transactions this month
	MonthlyExpenses    decimal.Decimal `json:"monthly_expenses"`      // Total EXPENSE transactions this month
	MonthlyNetCashFlow decimal.Decimal `json:"monthly_net_cash_flow"` // Income - Expenses

	// Runway Calculation (Critical feature from DOC.md)
	Runway                 float64         `json:"runway"`                   // Months of expenses covered by liquid assets
	RunwayDays             int             `json:"runway_days"`              // Runway in days
	AverageMonthlyExpenses decimal.Decimal `json:"average_monthly_expenses"` // Last 3 months average

	// Account Breakdown
	AccountBalances []AccountBalanceDTO `json:"account_balances"` // Current balance of each account

	// Metadata
	AsOf         time.Time `json:"as_of"`         // When this data was calculated
	BaseCurrency string    `json:"base_currency"` // User's base currency (default: USD)
}

// AccountBalanceDTO represents the current balance of an account for dashboard display
type AccountBalanceDTO struct {
	ID             uint            `json:"id"`
	Name           string          `json:"name"`
	Type           string          `json:"type"`
	Balance        decimal.Decimal `json:"balance"`
	CurrencyCode   string          `json:"currency_code"`
	CurrencySymbol string          `json:"currency_symbol"`
	IsActive       bool            `json:"is_active"`
}

// MonthlyStats represents income/expense stats for a specific month
type MonthlyStats struct {
	Month            int             `json:"month"` // 1-12
	Year             int             `json:"year"`
	Income           decimal.Decimal `json:"income"`
	Expenses         decimal.Decimal `json:"expenses"`
	NetCashFlow      decimal.Decimal `json:"net_cash_flow"` // Income - Expenses
	TransactionCount int             `json:"transaction_count"`
}

// MonthlyStatsResponse represents a time series of monthly statistics
type MonthlyStatsResponse struct {
	Stats         []MonthlyStats  `json:"stats"`
	TotalIncome   decimal.Decimal `json:"total_income"`
	TotalExpenses decimal.Decimal `json:"total_expenses"`
	StartDate     time.Time       `json:"start_date"`
	EndDate       time.Time       `json:"end_date"`
}

// RunwayCalculation represents detailed runway calculation breakdown
type RunwayCalculation struct {
	LiquidAssets           decimal.Decimal `json:"liquid_assets"`            // BANK + CASH
	ShortTermLiabilities   decimal.Decimal `json:"short_term_liabilities"`   // CREDIT_CARD balances
	AvailableFunds         decimal.Decimal `json:"available_funds"`          // LiquidAssets - ShortTermLiabilities
	AverageMonthlyExpenses decimal.Decimal `json:"average_monthly_expenses"` // Last 3 months
	RunwayMonths           float64         `json:"runway_months"`            // AvailableFunds / AvgExpenses
	RunwayDays             int             `json:"runway_days"`
	CalculationDate        time.Time       `json:"calculation_date"`
	BaseCurrency           string          `json:"base_currency"`

	// Breakdown by account type
	BankAccounts       []AccountBalanceDTO `json:"bank_accounts"`
	CashAccounts       []AccountBalanceDTO `json:"cash_accounts"`
	CreditCardAccounts []AccountBalanceDTO `json:"credit_card_accounts"`

	// Warning levels
	Status  string `json:"status"` // "HEALTHY", "WARNING", "CRITICAL"
	Message string `json:"message"`
}

// CategoryExpenseBreakdown represents expense breakdown by category
type CategoryExpenseBreakdown struct {
	CategoryID       uint            `json:"category_id"`
	CategoryName     string          `json:"category_name"`
	Amount           decimal.Decimal `json:"amount"`
	Percentage       float64         `json:"percentage"` // Percentage of total expenses
	TransactionCount int             `json:"transaction_count"`
}

// CategoryExpenseBreakdownResponse represents expenses grouped by category
type CategoryExpenseBreakdownResponse struct {
	Breakdown     []CategoryExpenseBreakdown `json:"breakdown"`
	TotalExpenses decimal.Decimal            `json:"total_expenses"`
	Month         int                        `json:"month"`
	Year          int                        `json:"year"`
}

// DashboardFilters represents filters for dashboard data queries
type DashboardFilters struct {
	StartDate    *time.Time `json:"start_date" validate:"omitempty"`
	EndDate      *time.Time `json:"end_date" validate:"omitempty"`
	BaseCurrency string     `json:"base_currency" validate:"omitempty,len=3"` // ISO 4217 code
}
