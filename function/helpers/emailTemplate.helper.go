package helpers

import "fmt"

const (
	htmlBody = `
		<style>
			.header {
				background-color: #F3F5F6;
				text-align: center;
				padding: 10px;
				font-family: Arial, sans-serif;
				color: #003A40;
			}

			.header img {
				width: 130px;
			}

			.body {
				background-color: #ffffff;
				padding: 20px;
				font-family: Arial, sans-serif;
				color: #333333;
			}

			.footer {
				background-color: #f3f3f3;
				text-align: center;
				padding: 10px;
				font-size: 12px;
				font-family: Arial, sans-serif;
			}

			.footer a {
				color: #0073e6;
				text-decoration: none;
				margin: 0 5px;
			}

			.footer a:hover {
				text-decoration: underline;
			}
		</style>

		<div class="header">
			<img src="https://upload.wikimedia.org/wikipedia/commons/thumb/b/b0/Stori_Logo_2023.svg/1920px-Stori_Logo_2023.svg.png" alt="Stori" width="130"/>
			<h1>Hello %s!, we present your account summary.</h1>
			<h2>Summary Account Report</h1>
		</div>
	
		<div class="body">
			<h3>Total balance is: %v</h2> 
			<h4>Transactions per Month</h2> 
			<ul>%s</ul>
			<p>Average debit amount: %v</p>
			<p>Average credit amount: %v</p>
		</div>
	
		<div class="footer">
			<p>Follow us on our social networks:</p>
			<a href="https://www.facebook.com/Stori.MX" target="_blank">Facebook</a>
			<a href="https://twitter.com/mi_stori" target="_blank">Twitter</a>
			<a href="https://twitter.com/mi_stori" target="_blank">Instagram</a>
		</div>
`
)

type SummaryTemplate struct {
	ClientName          string
	TotalBalance        float64
	MonthlyTransactions map[string]int
	AverageDebitAmount  float64
	AverageCreditAmount float64
}

func GenerateSummaryTemplate(summary SummaryTemplate) string {
	transactionsByMonth := generateTransactionsPerMonth(summary.MonthlyTransactions)
	return fmt.Sprintf(
		htmlBody,
		summary.ClientName,
		summary.TotalBalance,
		transactionsByMonth,
		summary.AverageDebitAmount,
		summary.AverageCreditAmount,
	)
}

func generateTransactionsPerMonth(transactionsByMonth map[string]int) string {
	lis := ""
	for month, transactions := range transactionsByMonth {
		lis += fmt.Sprintf("<li>%s: %v</li>", month, transactions)
	}
	return lis
}
