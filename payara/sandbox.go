package payara

// SandboxDummyAccount is a test account for Payara sandbox (https://doc.payara.id/docs/1.0/sandbox-data-dummy).
type SandboxDummyAccount struct {
	BankCode      string
	BankName      string
	AccountNumber string
	AccountName   string
}

// SandboxDummyAccounts are the official sandbox test accounts. Use these for development/testing only.
var SandboxDummyAccounts = []SandboxDummyAccount{
	{BankCode: "4", BankName: "Bank Mandiri", AccountNumber: "12340995811", AccountName: "Ujang"},
	{BankCode: "5", BankName: "Bank Central Asia", AccountNumber: "12330922231", AccountName: "Asep"},
	{BankCode: "6", BankName: "Bank Jago Syariah", AccountNumber: "12389583322", AccountName: "Robert"},
	{BankCode: "281", BankName: "OVO", AccountNumber: "081212239281", AccountName: "Rudi"},
	{BankCode: "282", BankName: "DANA", AccountNumber: "081212239133", AccountName: "Zen"},
	{BankCode: "283", BankName: "GOPAY", AccountNumber: "081212239222", AccountName: "Malik"},
}

// SandboxDummyAccountByBankCode returns the first dummy account with the given bank_code, or nil.
func SandboxDummyAccountByBankCode(bankCode string) *SandboxDummyAccount {
	for i := range SandboxDummyAccounts {
		if SandboxDummyAccounts[i].BankCode == bankCode {
			return &SandboxDummyAccounts[i]
		}
	}
	return nil
}

// DefaultSandboxAccount returns the BCA dummy account (bank_code "5") for examples.
func DefaultSandboxAccount() SandboxDummyAccount {
	return SandboxDummyAccounts[1] // Asep, BCA
}
