package manager

import (
	"testing"
)

func TestLoanService_AddRepayments(t *testing.T) {
	type fields struct {
		Loans map[int](*Loan)
	}
	type args struct {
		loanID int
		amount float64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "shikher",
			fields: fields{
				Loans: map[int](*Loan){
					1: &Loan{
						LoanId:             1,
						Requester:          "shikher",
						Amount:             100,
						InterestRate:       0.10,
						LoanPeriodInWeeks:  2,
						LoanApprovalStatus: APPROVED,
						LoanPaidStatus:     PAYMENTDUE,
						Payment: []Payment{
							Payment{
								Week:      1,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
							Payment{
								Week:      2,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
						},
					},
				},
			},
			args: args{
				loanID: 1,
				amount: 55,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoanService{
				Loans: tt.fields.Loans,
			}
			if err := ls.AddRepayments(tt.args.loanID, tt.args.amount); (err != nil) != tt.wantErr {
				t.Errorf("AddRepayments() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestLoanService_ApplyLoan(t *testing.T) {
	type fields struct {
		Loans map[int](*Loan)
	}
	type args struct {
		user              string
		amount            float64
		loanPeriodInWeeks int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    fields
		wantval int
	}{
		// TODO: Add test cases.
		{
			name: "sample loan",
			fields: fields{
				Loans: make(map[int](*Loan)),
			},
			args: args{
				user:              "123",
				amount:            100,
				loanPeriodInWeeks: 2,
			},
			want: fields{
				Loans: map[int](*Loan){
					1: &Loan{
						LoanId:             1,
						Requester:          "shikher",
						Amount:             100,
						InterestRate:       0.10,
						LoanPeriodInWeeks:  2,
						LoanApprovalStatus: PENDING,
						LoanPaidStatus:     PAYMENTDUE,
						Payment: []Payment{
							Payment{
								Week:      1,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
							Payment{
								Week:      2,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
						},
					},
				},
			},
			wantval: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoanService{
				Loans: tt.fields.Loans,
			}
			if got := ls.ApplyLoan(tt.args.user, tt.args.amount, tt.args.loanPeriodInWeeks); got != tt.wantval && ls.Loans[1].LoanApprovalStatus != tt.want.Loans[1].LoanApprovalStatus &&
				ls.Loans[1].Payment[0] != tt.want.Loans[1].Payment[0] {
				t.Errorf("ApplyLoan() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLoanService_ApproveLoan(t *testing.T) {
	type fields struct {
		Loans map[int](*Loan)
	}
	type args struct {
		loanID int
		status bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		// TODO: Add test cases.
		{
			name: "sample loan",
			fields: fields{
				Loans: map[int](*Loan){
					1: &Loan{
						LoanId:             1,
						Requester:          "shikher",
						Amount:             100,
						InterestRate:       0.10,
						LoanPeriodInWeeks:  2,
						LoanApprovalStatus: PENDING,
						LoanPaidStatus:     PAYMENTDUE,
						Payment: []Payment{
							Payment{
								Week:      1,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
							Payment{
								Week:      2,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
						},
					},
				},
			},
			args: args{
				loanID: 1,
				status: true,
			},
			want: fields{
				Loans: map[int](*Loan){
					1: &Loan{
						LoanId:             1,
						Requester:          "shikher",
						Amount:             100,
						InterestRate:       0.10,
						LoanPeriodInWeeks:  2,
						LoanApprovalStatus: APPROVED,
						LoanPaidStatus:     PAYMENTDUE,
						Payment: []Payment{
							Payment{
								Week:      1,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
							Payment{
								Week:      2,
								AmountDue: 55,
								Status:    PAYMENTDUE,
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ls := &LoanService{
				Loans: tt.fields.Loans,
			}
			ls.ApproveLoan(tt.args.loanID, tt.args.status)
			if ls.Loans[1].LoanApprovalStatus != tt.want.Loans[1].LoanApprovalStatus {
				t.Errorf("ApplyLoan() = %v, want %v", ls.Loans[1].LoanApprovalStatus, tt.want)
			}
		})
	}
}
