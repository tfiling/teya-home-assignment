package ledger_test

import (
	"testing"
	"teya_home_assignment/internal/pkg/ledger"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLedger_NewLedger__CreatesEmptyLedger(t *testing.T) {
	// Act
	ledgerInstance, err := ledger.NewLedger()

	// Assert
	require.NoError(t, err)
	assert.NotNil(t, ledgerInstance)
	assert.Empty(t, ledgerInstance.TransactionHistory)
	balanceDecimal, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	balance, exact := balanceDecimal.Float64()
	assert.True(t, exact)
	assert.Equal(t, float64(0), balance)
}

func TestLedger_AddTransaction__AddsPositiveAmountTransaction(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)
	amount := decimal.NewFromFloat32(100.50)

	// Act
	err = ledgerInstance.AddTransaction(amount)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, ledgerInstance.TransactionHistory, 1)
	assert.Equal(t, uint64(1), ledgerInstance.TransactionHistory[0].ID)
	assert.Equal(t, amount, ledgerInstance.TransactionHistory[0].Amount)
}

func TestLedger_AddTransaction__AddsNegativeAmountTransaction(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)
	amount := decimal.NewFromFloat32(-50.25)

	// Act
	err = ledgerInstance.AddTransaction(amount)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, ledgerInstance.TransactionHistory, 1)
	assert.Equal(t, uint64(1), ledgerInstance.TransactionHistory[0].ID)
	assert.Equal(t, amount, ledgerInstance.TransactionHistory[0].Amount)
}

func TestLedger_AddTransaction__AssignsUniqueIDs(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	// Act
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(200))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(300))
	assert.NoError(t, err)

	// Assert
	assert.Len(t, ledgerInstance.TransactionHistory, 3)
	assert.Equal(t, uint64(1), ledgerInstance.TransactionHistory[0].ID)
	assert.Equal(t, uint64(2), ledgerInstance.TransactionHistory[1].ID)
	assert.Equal(t, uint64(3), ledgerInstance.TransactionHistory[2].ID)
}

func TestLedger_GetBalance__ReturnsZeroForEmptyLedger(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	// Act
	balanceDecimal, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	balance, exact := balanceDecimal.Float64()

	// Assert
	assert.True(t, exact)
	assert.Equal(t, float64(0), balance)
}

func TestLedger_GetBalance__ReturnsSumOfTransactions(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.50))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(200.75))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(-50.25))
	assert.NoError(t, err)

	expectedBalance := 100.50 + 200.75 - 50.25

	// Act
	balanceDecimal, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	balance, exact := balanceDecimal.Float64()

	// Assert
	assert.True(t, exact)
	assert.Equal(t, expectedBalance, balance)
}

func TestLedger_GetBalance__HandlesFloatingPointPrecision(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	// Add transactions with potential floating point precision issues
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.1))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.2))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.3))
	assert.NoError(t, err)

	// Act
	balanceDecimal, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	balance, exact := balanceDecimal.Float64()

	// Assert
	assert.Equal(t, 0.6, balance)
	// Assert that 'exact' is false, which indicates that the decimal-to-float64 conversion
	// couldn't be done with exact precision. This is expected behavior when working with
	// decimal values (0.1, 0.2, 0.3) that don't have exact binary floating-point representations.
	assert.False(t, exact)
}
