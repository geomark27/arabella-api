package services

import (
	"arabella-api/internal/app/dtos"
	"arabella-api/internal/app/models"
	"arabella-api/internal/app/repositories"
	"time"

	"github.com/shopspring/decimal"
)

// DashboardService provides dashboard data including the critical Runway calculation
type DashboardService interface {
	GetDashboard(userID uint) (*dtos.DashboardResponse, error)
	CalculateRunway(userID uint) (*dtos.RunwayCalculation, error)
	GetMonthlyStats(userID uint, month, year int) (*dtos.MonthlyStats, error)
}

type dashboardService struct {
	accountRepo     repositories.AccountRepository
	transactionRepo repositories.TransactionRepository
}

// NewDashboardService creates a new dashboard service
func NewDashboardService(
	accountRepo repositories.AccountRepository,
	transactionRepo repositories.TransactionRepository,
) DashboardService {
	return &dashboardService{
		accountRepo:     accountRepo,
		transactionRepo: transactionRepo,
	}
}

// GetDashboard returns complete dashboard data for a user
func (s *dashboardService) GetDashboard(userID uint) (*dtos.DashboardResponse, error) {
	now := time.Now()

	// Get total assets
	totalAssets, err := s.accountRepo.GetTotalAssets(userID)
	if err != nil {
		return nil, err
	}

	// Get total liabilities
	totalLiabilities, err := s.accountRepo.GetTotalLiabilities(userID)
	if err != nil {
		return nil, err
	}

	// Calculate net worth
	netWorth := totalAssets.Sub(totalLiabilities)

	// Get liquid assets (for runway)
	liquidAssets, err := s.accountRepo.GetLiquidAssets(userID)
	if err != nil {
		return nil, err
	}

	// Get current month stats
	currentMonth := int(now.Month())
	currentYear := now.Year()
	monthlyIncome, monthlyExpenses, _, err := s.transactionRepo.GetMonthlyStats(userID, currentMonth, currentYear)
	if err != nil {
		return nil, err
	}

	monthlyNetCashFlow := monthlyIncome.Sub(monthlyExpenses)

	// Calculate runway
	runway, runwayDays, avgMonthlyExpenses := s.calculateRunway(userID, liquidAssets, totalLiabilities)

	// Get account balances
	accounts, err := s.accountRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	accountBalances := make([]dtos.AccountBalanceDTO, len(accounts))
	for i, account := range accounts {
		currencyCode := "USD" // default
		currencySymbol := "$"
		if account.Currency != nil {
			currencyCode = account.Currency.Code
			currencySymbol = account.Currency.Symbol
		}

		accountBalances[i] = dtos.AccountBalanceDTO{
			ID:             account.ID,
			Name:           account.Name,
			Type:           account.AccountType,
			Balance:        account.Balance,
			CurrencyCode:   currencyCode,
			CurrencySymbol: currencySymbol,
			IsActive:       account.IsActive,
		}
	}

	return &dtos.DashboardResponse{
		TotalAssets:            totalAssets,
		TotalLiabilities:       totalLiabilities,
		NetWorth:               netWorth,
		LiquidAssets:           liquidAssets,
		MonthlyIncome:          monthlyIncome,
		MonthlyExpenses:        monthlyExpenses,
		MonthlyNetCashFlow:     monthlyNetCashFlow,
		Runway:                 runway,
		RunwayDays:             runwayDays,
		AverageMonthlyExpenses: avgMonthlyExpenses,
		AccountBalances:        accountBalances,
		AsOf:                   now,
		BaseCurrency:           "USD", // TODO: Make this configurable per user
	}, nil
}

// calculateRunway calculates how many months a user can survive without income
// Formula: (Liquid Assets - Short-term Liabilities) / Average Monthly Expenses (last 3 months)
func (s *dashboardService) calculateRunway(userID uint, liquidAssets, liabilities decimal.Decimal) (float64, int, decimal.Decimal) {
	// Available funds = liquid assets - liabilities
	availableFunds := liquidAssets.Sub(liabilities)

	// Calculate average monthly expenses over last 3 months
	now := time.Now()
	var totalExpenses decimal.Decimal
	var monthsWithData int

	for i := 0; i < 3; i++ {
		targetDate := now.AddDate(0, -i, 0)
		month := int(targetDate.Month())
		year := targetDate.Year()

		_, expenses, _, err := s.transactionRepo.GetMonthlyStats(userID, month, year)
		if err == nil && !expenses.IsZero() {
			totalExpenses = totalExpenses.Add(expenses)
			monthsWithData++
		}
	}

	// If no expense data, return 0
	if monthsWithData == 0 || totalExpenses.IsZero() {
		return 0, 0, decimal.Zero
	}

	// Calculate average
	avgMonthlyExpenses := totalExpenses.Div(decimal.NewFromInt(int64(monthsWithData)))

	// Calculate runway in months
	runwayMonths, _ := availableFunds.Div(avgMonthlyExpenses).Float64()

	// Calculate runway in days
	runwayDays := int(runwayMonths * 30)

	return runwayMonths, runwayDays, avgMonthlyExpenses
}

// CalculateRunway returns detailed runway calculation
func (s *dashboardService) CalculateRunway(userID uint) (*dtos.RunwayCalculation, error) {
	now := time.Now()

	liquidAssets, err := s.accountRepo.GetLiquidAssets(userID)
	if err != nil {
		return nil, err
	}

	liabilities, err := s.accountRepo.GetTotalLiabilities(userID)
	if err != nil {
		return nil, err
	}

	availableFunds := liquidAssets.Sub(liabilities)

	runwayMonths, runwayDays, avgExpenses := s.calculateRunway(userID, liquidAssets, liabilities)

	// Determine status
	status := "HEALTHY"
	message := "Your financial runway is healthy. Keep it up!"

	if runwayMonths < 3 {
		status = "CRITICAL"
		message = "⚠️ CRITICAL: Your runway is less than 3 months. Consider increasing income or reducing expenses immediately."
	} else if runwayMonths < 6 {
		status = "WARNING"
		message = "⚠️ WARNING: Your runway is below 6 months. Consider building a larger emergency fund."
	}

	// Get account breakdowns
	bankAccounts, _ := s.accountRepo.FindByUserAndType(userID, "BANK")
	cashAccounts, _ := s.accountRepo.FindByUserAndType(userID, "CASH")
	creditCardAccounts, _ := s.accountRepo.FindByUserAndType(userID, "CREDIT_CARD")

	return &dtos.RunwayCalculation{
		LiquidAssets:           liquidAssets,
		ShortTermLiabilities:   liabilities,
		AvailableFunds:         availableFunds,
		AverageMonthlyExpenses: avgExpenses,
		RunwayMonths:           runwayMonths,
		RunwayDays:             runwayDays,
		CalculationDate:        now,
		BaseCurrency:           "USD",
		BankAccounts:           convertAccountsToBalanceDTO(bankAccounts),
		CashAccounts:           convertAccountsToBalanceDTO(cashAccounts),
		CreditCardAccounts:     convertAccountsToBalanceDTO(creditCardAccounts),
		Status:                 status,
		Message:                message,
	}, nil
}

// GetMonthlyStats returns income/expense statistics for a specific month
func (s *dashboardService) GetMonthlyStats(userID uint, month, year int) (*dtos.MonthlyStats, error) {
	income, expenses, count, err := s.transactionRepo.GetMonthlyStats(userID, month, year)
	if err != nil {
		return nil, err
	}

	netCashFlow := income.Sub(expenses)

	return &dtos.MonthlyStats{
		Month:            month,
		Year:             year,
		Income:           income,
		Expenses:         expenses,
		NetCashFlow:      netCashFlow,
		TransactionCount: int(count),
	}, nil
}

// Helper function to convert accounts to balance DTOs
func convertAccountsToBalanceDTO(accounts []*models.Account) []dtos.AccountBalanceDTO {
	result := make([]dtos.AccountBalanceDTO, len(accounts))
	for i, account := range accounts {
		currencyCode := "USD"
		currencySymbol := "$"
		if account.Currency != nil {
			currencyCode = account.Currency.Code
			currencySymbol = account.Currency.Symbol
		}

		result[i] = dtos.AccountBalanceDTO{
			ID:             account.ID,
			Name:           account.Name,
			Type:           account.AccountType,
			Balance:        account.Balance,
			CurrencyCode:   currencyCode,
			CurrencySymbol: currencySymbol,
			IsActive:       account.IsActive,
		}
	}
	return result
}
