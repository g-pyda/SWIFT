package structs

type ReqBranch struct {
	Address string `json:"address" binding:"max=255"`
    BankName string `json:"bankName" binding:"required,max=255"`
    CountryISO2 string `json:"countryISO2" binding:"required,min=2,max=2"`
    CountryName string `json:"countryName" binding:"required,max=20"`
    IsHeadquarter *bool `json:"isHeadquarter" binding:"required"`
    SwiftCode string `json:"swiftCode" binding:"required,min=11,max=11"`
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

type ReqHeadquarterBranch struct {
	Address string `json:"address"`
    BankName string `json:"bankName"`
    CountryISO2 string `json:"countryISO2"`
    CountryName string `json:"countryName"`
    IsHeadquarter bool `json:"isHeadquarter"`
    SwiftCode string `json:"swiftCode"`
    Branches *[]ReqBranch `json:"branches,omitempty"`
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


type ReqAll struct {
    Entries []ReqHeadBranInCountry `json:entries`
}

type ReqErr struct {
    Message string `json:message`
}