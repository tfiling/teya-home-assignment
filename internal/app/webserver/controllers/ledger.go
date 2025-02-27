package controllers

import (
	"strconv"
	"teya_home_assignment/internal/app/webserver/api"
	"teya_home_assignment/internal/pkg/ledger"

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
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid request body")
	}
	transactionAmount, err := decimal.NewFromString(reqBody.Amount)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid transaction amount")
	}
	if err := c.ledgerService.AddTransaction(transactionAmount); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("could not add transaction")
	}
	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *LedgerController) getAllTransaction(ctx *fiber.Ctx) error {
	// Default limit. It is optional
	limit := 10

	offsetParam := ctx.Query("offset")
	if offsetParam == "" {
		return ctx.Status(fiber.StatusBadRequest).SendString("missing offset query parameter")
	}
	offset, err := strconv.Atoi(offsetParam)
	if err != nil || offset < 0 {
		return ctx.Status(fiber.StatusBadRequest).SendString("invalid offset parameter")
	}

	if limitParam := ctx.Query("limit"); limitParam != "" {
		limit, err = strconv.Atoi(limitParam)
		if err != nil || limit <= 0 || limit > 100 {
			return ctx.Status(fiber.StatusBadRequest).SendString("invalid limit query parameter: must be between 1 and 100")
		}
	}
	transactionsHistory, err := c.ledgerService.GetTransactionHistory(offset, limit)
	if err != nil {
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

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (c *LedgerController) getBalance(ctx *fiber.Ctx) error {
	balance, err := c.ledgerService.GetBalance()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString("could not get balance")
	}
	resp := api.GetBalanceRespBody{
		Balance: balance.String(),
	}
	return ctx.Status(fiber.StatusOK).JSON(resp)
}
