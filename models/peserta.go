package models

type Peserta struct {
	MemberId        string  `json:"memberId"`
	LoanNo          string  `json:"loanNo"`
	IdentityNo      string  `json:"identityNo"`
	Name            string  `json:"name"`
	BirthDate       string  `json:"birthDate"`
	Gender          string  `json:"gender"`
	Address         *string `json:"address"`
	InsuranceStart  string  `json:"insuranceStart"`
	InsuranceEnd    string  `json:"insuranceEnd"`
	InsurancePeriod int     `json:"insurancePeriod"`
	InsuredValue    float32 `json:"insuredValue"`
	Rate            float32 `json:"rate"`
	Premium         float32 `json:"premium"`
	Status          string  `json:"status"`
}
