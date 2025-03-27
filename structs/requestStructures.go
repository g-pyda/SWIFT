package structs

type ReqBranch struct {
	Address string
    BankName string
    CountryISO2 string
    CountryName string
    IsHeadquarter bool
    SwiftCode string
}

type ReqHeadquarter struct {
	Address string
    BankName string
    CountryISO2 string
    CountryName string
    IsHeadquarter bool
    SwiftCode string
    Branches []ReqBranch
}

type ReqCountry struct {
	CountryISO2 string
    CountryName string
	SwiftCodes_hq []ReqHeadquarter
    SwiftCodes_br []ReqBranch
}