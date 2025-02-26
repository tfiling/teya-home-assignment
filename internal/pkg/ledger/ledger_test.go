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
	balance, exact := ledgerInstance.GetBalance().Float64()
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
	ledgerInstance.AddTransaction(amount)

	// Assert
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
	ledgerInstance.AddTransaction(amount)

	// Assert
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
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(100))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(200))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(300))

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
	balance, exact := ledgerInstance.GetBalance().Float64()

	// Assert
	assert.True(t, exact)
	assert.Equal(t, float64(0), balance)
}

func TestLedger_GetBalance__ReturnsSumOfTransactions(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.50))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(200.75))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(-50.25))

	expectedBalance := 100.50 + 200.75 - 50.25

	// Act
	balance, exact := ledgerInstance.GetBalance().Float64()

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
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.1))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.2))
	ledgerInstance.AddTransaction(decimal.NewFromFloat32(0.3))

	// Act
	balance, exact := ledgerInstance.GetBalance().Float64()

	// Assert
	assert.Equal(t, 0.6, balance)
	// Assert that 'exact' is false, which indicates that the decimal-to-float64 conversion
	// couldn't be done with exact precision. This is expected behavior when working with
	// decimal values (0.1, 0.2, 0.3) that don't have exact binary floating-point representations.
	assert.False(t, exact)
}
