package handlers

import (
	"aspire/authMiddleware"
	"aspire/manager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type LoanRequest struct {
	Amount            float64 `json:"amount"`
	LoanPeriodInWeeks int     `json:"loan_period_in_weeks"`
}

type LoanApproval struct {
	Approve bool `json:"approve"`
}

func GetToken(c *gin.Context, ls *manager.LoanService) {
	// In a prod scenario, we validate userId with password.
	userID := c.Param("userId")
	token, err := authMiddleware.CreateToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func ApproveLoan(c *gin.Context, ls *manager.LoanService) {
	loanID := c.Param("loanId")
	var loanApproval LoanApproval
	if err := c.ShouldBindJSON(&loanApproval); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loanVal, err := strconv.Atoi(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.Loans[loanVal] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid LoanId"})
		return
	}
	ls.ApproveLoan(loanVal, loanApproval.Approve)
	c.JSON(http.StatusOK, ls.Loans[loanVal])
	return
}

func ApplyLoan(c *gin.Context, ls *manager.LoanService) {
	var loanRequest LoanRequest
	if err := c.ShouldBindJSON(&loanRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := c.GetString("user_id")
	loanID := ls.ApplyLoan(userID, loanRequest.Amount, loanRequest.LoanPeriodInWeeks)
	c.JSON(200, ls.Loans[loanID])
}

func GetLoanDetails(c *gin.Context, ls *manager.LoanService) {
	loanID := c.Param("loanID")
	loanVal, err := strconv.Atoi(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.Loans[loanVal] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid LoanId"})
		return
	}
	c.JSON(http.StatusOK, ls.Loans[loanVal])
	return
}

func AddLoanRepayment(c *gin.Context, ls *manager.LoanService) {
	loanID := c.Param("loanID")
	amount := c.Param("amount")
	loanVal, err := strconv.Atoi(loanID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if ls.Loans[loanVal] == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid LoanId"})
		return
	}
	amountPay, err := strconv.ParseFloat(amount, 8)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err = ls.AddRepayments(loanVal, amountPay)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ls.Loans[loanVal])
	return
}
