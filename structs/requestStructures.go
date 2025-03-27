package structs

type ReqBranch struct {
	Address string `json:"address"`
    BankName string `json:"bankName"`
    CountryISO2 string `json:"countryISO2"`
    CountryName string `json:"countryName"`
    IsHeadquarter bool `json:"isHeadquarter"`
    SwiftCode string `json:"swiftCode"`
}

type ReqHeadquarter struct {
	Address string `json:"address"`
    BankName string `json:"bankName"`
    CountryISO2 string `json:"countryISO2"`
    CountryName string `json:"countryName"`
    IsHeadquarter bool `json:"isHeadquarter"`
    SwiftCode string `json:"swiftCode"`
    Branches []ReqBranch `json:"branches"`
}

type ReqCountry struct {
	CountryISO2 string `json:"countryISO2"`
    CountryName string `json:"countryName"`
	SwiftCodes []ReqHeadBranInCountry `json:"swiftCodes"`
}

type ReqHeadBranInCountry struct {
    Address string `json:"address"`
    BankName string `json:"bankName"`
    CountryISO2 string `json:"countryISO2"`
    IsHeadquarter bool `json:"isHeadquarter"`
    SwiftCode string `json:"swiftCode"`
}

