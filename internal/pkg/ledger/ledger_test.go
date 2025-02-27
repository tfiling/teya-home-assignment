package ledger_test

import (
	"fmt"
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

func TestLedger_GetTransactionHistory__ReturnsEmptyForEmptyLedger(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	// Act
	history, err := ledgerInstance.GetTransactionHistory(0, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, history)
}

func TestLedger_GetTransactionHistory__ReturnsAllTransactionsWithinLimit(t *testing.T) {
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

	// Act
	history, err := ledgerInstance.GetTransactionHistory(0, 3)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, history, 3)
	assert.Equal(t, uint64(1), history[0].ID)
	assert.Equal(t, uint64(2), history[1].ID)
	assert.Equal(t, uint64(3), history[2].ID)
}

func TestLedger_GetTransactionHistory__RespectsLimit(t *testing.T) {
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

	// Act
	history, err := ledgerInstance.GetTransactionHistory(0, 2)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, uint64(1), history[0].ID)
	assert.Equal(t, uint64(2), history[1].ID)
}

func TestLedger_GetTransactionHistory__RespectsOffset(t *testing.T) {
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

	// Act
	history, err := ledgerInstance.GetTransactionHistory(1, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, history, 2)
	assert.Equal(t, uint64(2), history[0].ID)
	assert.Equal(t, uint64(3), history[1].ID)
}

func TestLedger_GetTransactionHistory__HandlesOffsetBeyondSize(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.50))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(200.75))
	assert.NoError(t, err)

	// Act
	history, err := ledgerInstance.GetTransactionHistory(5, 10)

	// Assert
	assert.NoError(t, err)
	assert.Empty(t, history)
}

func TestLedger_GetTransactionHistory__PreservesTransactionOrder(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)
	amounts := []float32{100.50, -25.75, 300.00, -150.25, 50.00}
	for _, amount := range amounts {
		err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(amount))
		assert.NoError(t, err)
	}

	// Act
	history, err := ledgerInstance.GetTransactionHistory(0, 10)

	// Assert
	assert.NoError(t, err)
	assert.Len(t, history, 5)
	for i, amount := range amounts {
		expectedAmount := decimal.NewFromFloat32(amount)
		assert.Equal(t, uint64(i+1), history[i].ID)
		assert.Equal(t, expectedAmount, history[i].Amount)
	}
}

func TestCachingLedger_GetBalance__UsesCachedBalanceForSubsequentCalls(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.50))
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(200.75))
	assert.NoError(t, err)
	expectedBalance := decimal.NewFromFloat32(100.50 + 200.75)
	initialBalance, err := ledgerInstance.GetBalance()

	// Act
	cachedBalance, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, expectedBalance.Equal(initialBalance), fmt.Sprintf("%+v != %+v", expectedBalance, initialBalance))
	assert.True(t, expectedBalance.Equal(cachedBalance), fmt.Sprintf("%+v != %+v", expectedBalance, cachedBalance))
}

func TestCachingLedger_GetBalance__UpdatesCacheAfterNewTransactions(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.00))
	assert.NoError(t, err)
	initialBalance, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(100.00).Equal(initialBalance),
		fmt.Sprintf("%+v != 100.00", initialBalance))

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(50.00))
	assert.NoError(t, err)

	// Act
	updatedBalance, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(150.00).Equal(updatedBalance),
		fmt.Sprintf("%+v != 150.00", updatedBalance))
}

func TestCachingLedger_GetBalance__CachesIncrementalUpdates(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.00))
	assert.NoError(t, err)
	_, err = ledgerInstance.GetBalance()
	assert.NoError(t, err)
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(50.00))
	assert.NoError(t, err)
	balanceAfterSecond, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(150.00).Equal(balanceAfterSecond),
		fmt.Sprintf("%+v != 150.00", balanceAfterSecond))

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(25.00))
	assert.NoError(t, err)

	// Act
	finalBalance, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(175.00).Equal(finalBalance),
		fmt.Sprintf("%+v != 150.00", finalBalance))
}

func TestCachingLedger_GetBalance__HandlesEmptyLedgerCaching(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	initialBalance, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(0).Equal(initialBalance),
		fmt.Sprintf("%+v != 0.00", initialBalance))

	// Act
	cachedBalance, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(0).Equal(cachedBalance),
		fmt.Sprintf("%+v != 0.00", cachedBalance))
}

func TestCachingLedger_GetBalance__PreservesCacheWithNegativeTransactions(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(100.00))
	assert.NoError(t, err)
	initialBalance, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(100.00).Equal(initialBalance),
		fmt.Sprintf("%+v != 100.00", initialBalance))
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(-30.00))
	assert.NoError(t, err)

	// Act
	updatedBalance, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(70.00).Equal(updatedBalance),
		fmt.Sprintf("%+v != 70.00", updatedBalance))
}

func TestCachingLedger_GetBalance__CacheWorksWithMultipleCallsAndTransactions(t *testing.T) {
	// Arrange
	ledgerInstance, err := ledger.NewLedger()
	require.NoError(t, err)
	require.NotNil(t, ledgerInstance)

	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(10.00))
	assert.NoError(t, err)
	firstBalance, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(10.00).Equal(firstBalance),
		fmt.Sprintf("%+v != 10.00", firstBalance))
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(20.00))
	assert.NoError(t, err)
	secondBalance, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(30.00).Equal(secondBalance),
		fmt.Sprintf("%+v != 30.00", secondBalance))
	err = ledgerInstance.AddTransaction(decimal.NewFromFloat32(30.00))
	assert.NoError(t, err)

	// Act
	thirdBalanceFirstCall, err := ledgerInstance.GetBalance()
	assert.NoError(t, err)
	thirdBalanceSecondCall, err := ledgerInstance.GetBalance()

	// Assert
	assert.NoError(t, err)
	assert.True(t, decimal.NewFromFloat32(60.00).Equal(thirdBalanceFirstCall),
		fmt.Sprintf("%+v != 60.00", thirdBalanceFirstCall))
	assert.True(t, decimal.NewFromFloat32(60.00).Equal(thirdBalanceSecondCall),
		fmt.Sprintf("%+v != 60.00", thirdBalanceSecondCall))
}
