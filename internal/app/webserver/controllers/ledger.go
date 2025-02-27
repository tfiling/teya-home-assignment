package controllers

import (
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
	//TODO - implementing pagination
	res := make([]api.Transaction, 0)
	for _, transaction := range c.ledgerService.TransactionHistory {
		res = append(res, api.FromTransactionModel(transaction))
	}
	return ctx.Status(fiber.StatusOK).JSON(res)
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
