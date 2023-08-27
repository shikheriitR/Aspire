package main

import (
	"aspire/authMiddleware"
	"aspire/handlers"
	"aspire/manager"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	loanservice := manager.NewLoanService()
	r.GET("/getToken/:userId", func(c *gin.Context) {
		handlers.GetToken(c, loanservice)
	})
	// Protected route that requires authentication
	r.POST("/apply", authMiddleware.AuthMiddleware(""), func(c *gin.Context) {
		handlers.ApplyLoan(c, loanservice)
	})
	r.POST("/approve/:loanId", authMiddleware.AuthMiddleware("approve"), func(c *gin.Context) {
		handlers.ApproveLoan(c, loanservice)
	}) // ONLY ADMIN can approve which is user 789

	r.GET("/getLoanDetails/:loanID", authMiddleware.AuthMiddleware(""), func(c *gin.Context) {
		handlers.GetLoanDetails(c, loanservice)
	})
	r.POST("/addRepayment/:loanID/:amount", authMiddleware.AuthMiddleware(""), func(c *gin.Context) {
		handlers.AddLoanRepayment(c, loanservice)
	})

	r.Run(":8080")
}
