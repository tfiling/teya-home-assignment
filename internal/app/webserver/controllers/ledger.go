package controllers

import (
	"fmt"
	"strconv"
	"teya_home_assignment/internal/app/webserver/api"
	"teya_home_assignment/internal/pkg/ledger"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/pkg/errors"
	"github.com/shopspring/decimal"
)

type LedgerController struct {
	ledgerService *ledger.Ledger
}

func NewLedgerController() (*LedgerController, error) {
	service, err := ledger.NewLedger()
	if err != nil {
		return nil, errors.Wrap(err, "could not create ledger controller")
	}
	return &LedgerController{ledgerService: service}, nil
}

func (c *LedgerController) RegisterRoutes(router fiber.Router) error {
	router.Post(TransactionRoute, c.createTransaction)
	router.Get(TransactionRoute, c.getAllTransaction)
	router.Get(AccountRoute, c.getBalance)
	return nil
}

func (c *LedgerController) createTransaction(ctx *fiber.Ctx) error {
	reqBody := api.NewTransactionReqBody{}
	if err := ctx.BodyParser(&reqBody); err != nil {
		fmt.Println("invalid request body on transaction create")
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}
	if err := validator.New().Struct(reqBody); err != nil {
		fmt.Printf("invalid request on transaction create: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	transactionAmount, err := decimal.NewFromString(reqBody.Amount)
	if err != nil {
		fmt.Printf("invalid request on transaction create: %v\n", err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid transaction amount")
	}
	if err := c.ledgerService.AddTransaction(transactionAmount); err != nil {
		fmt.Printf("failed to add transaction: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("could not add transaction")
	}
	fmt.Printf("successfully add transaction: %v\n", transactionAmount)
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *LedgerController) getAllTransaction(ctx *fiber.Ctx) error {
	// Default limit. It is optional
	limit := 10

	offsetParam := ctx.Query("offset")
	if offsetParam == "" {
		fmt.Println("invalid request on transaction getAllTransaction - missing offset parameter")
		return ctx.Status(fiber.StatusBadRequest).SendString("missing offset query parameter")
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		fmt.Printf("invalid request on transaction getAllTransaction - invalid offset parameter: %v(error: %v)\n",
			offsetParam, err)
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid offset parameter")
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 || limit > 100 {
			fmt.Printf("invalid request on transaction getAllTransaction - invalid limit parameter: %v(error: %v)\n",
				limitParam, err)
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid limit query parameter: must be between 1 and 100")
		}
	}
	transactionsHistory, err := c.ledgerService.GetTransactionHistory(offset, limit)
	if err != nil {
		fmt.Printf("failed to get transaction history: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("could not get transactions")
	}

	transactions := make([]api.Transaction, len(transactionsHistory))
	for i, transaction := range transactionsHistory {
		transactions[i] = api.FromTransactionModel(transaction)
	}

	response := api.PaginatedTransactionsResponse{
		Transactions: transactions,
		Pagination: api.Pagination{
			Offset: offset,
			Limit:  limit,
		},
	}

	fmt.Printf("successfully returned transactions: %v\n", transactions)
	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *LedgerController) getBalance(ctx *fiber.Ctx) error {
	balance, err := c.ledgerService.GetBalance()
	if err != nil {
		fmt.Printf("failed to get balance: %v\n", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString("could not get balance")
	}
	resp := api.GetBalanceRespBody{
		Balance: balance.String(),
	}
	fmt.Printf("successfully calculated balance: %v\n", balance)
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
