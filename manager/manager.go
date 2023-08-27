package manager

import (
	"errors"
	"fmt"
	"time"
)

type LoanApprovalStatus string

const (
	PENDING  LoanApprovalStatus = "PENDING"
	APPROVED LoanApprovalStatus = "APPROVED"
	REJECTED LoanApprovalStatus = "REJECTED"
)

type PaymentStatus string

const (
	PAYMENTPAID PaymentStatus = "PAID"
	PAYMENTDUE  PaymentStatus = "DUE"
)

const InterestRate float64 = 0.05 // 5% interest rate

type Loan struct {
	LoanId             int
	Requester          string
	Amount             float64
	InterestRate       float64
	LoanPeriodInWeeks  int
	LoanApprovalStatus LoanApprovalStatus
	LoanPaidStatus     PaymentStatus
	Payment            []Payment
	UserPayment        []float64
}

type Payment struct {
	Week           int
	AmountDue      float64
	PaymentDay     time.Time
	PaymentDueDate time.Time
	Status         PaymentStatus
}

type LoanService struct {
	Loans map[int](*Loan)
}

func NewLoanService() *LoanService {
	return &LoanService{
		Loans: make(map[int](*Loan)),
	}
}

func (ls *LoanService) ApplyLoan(user string, amount float64, loanPeriodInWeeks int) int {
	loanID := len(ls.Loans) + 1
	loan := &Loan{
		LoanId:             loanID,
		Requester:          user,
		Amount:             amount,
		InterestRate:       InterestRate,
		LoanPeriodInWeeks:  loanPeriodInWeeks,
		LoanApprovalStatus: PENDING,
		LoanPaidStatus:     PAYMENTDUE,
	}
	ls.Loans[loanID] = loan
	return loanID
}

func (ls *LoanService) ApproveLoan(loanID int, status bool) {
	loan := ls.Loans[loanID]
	if loan.LoanApprovalStatus != PENDING {
		fmt.Println("Can't update the loan status Now")
	}
	if status {
		loan.LoanApprovalStatus = APPROVED
		loan.GeneratePayments()
	} else {
		loan.LoanApprovalStatus = REJECTED
	}
}

func (ls *LoanService) AddRepayments(loanID int, amount float64) error {
	loan := ls.Loans[loanID]
	if loan.LoanApprovalStatus != APPROVED {
		return errors.New("Loan is not approved")
	}
	loan.UserPayment = append(loan.UserPayment, amount)
	for i, payment := range loan.Payment {
		if payment.Status == PAYMENTDUE {
			if amount > payment.AmountDue {
				loan.Payment[i].AmountDue = 0
				loan.Payment[i].Status = PAYMENTPAID
				loan.Payment[i].PaymentDay = time.Now()
				amount -= payment.AmountDue
			} else {
				loan.Payment[i].AmountDue -= amount
				break
			}
		}
	}
	if loan.Payment[len(loan.Payment)-1].Status == PAYMENTPAID {
		loan.LoanPaidStatus = PAYMENTPAID
	}
	return nil
}

func (ls *LoanService) GetPayments(loanID int) []Payment {
	return ls.Loans[loanID].Payment
}

func (ls *Loan) GeneratePayments() {
	payments := []Payment{}
	totalAmountDue := ls.Amount
	weeklyPayment := (ls.Amount * (1 + ls.InterestRate)) / float64(ls.LoanPeriodInWeeks)

	for week := 1; week <= ls.LoanPeriodInWeeks; week++ {
		payment := Payment{
			Week:           week,
			AmountDue:      weeklyPayment,
			PaymentDueDate: time.Now().AddDate(0, 0, week*7),
			Status:         PAYMENTDUE, // Assuming payments are due weekly.
		}
		payments = append(payments, payment)
		totalAmountDue -= weeklyPayment
	}

	// Adding a final payment for the remaining amount (if any).
	if totalAmountDue > 0 {
		finalPayment := Payment{
			Week:           ls.LoanPeriodInWeeks + 1,
			AmountDue:      totalAmountDue,
			PaymentDueDate: time.Now().AddDate(0, 0, (ls.LoanPeriodInWeeks+1)*7),
			Status:         PAYMENTDUE, // Assuming payments are due weekly.
		}
		payments = append(payments, finalPayment)
	}
	ls.Payment = payments
}
