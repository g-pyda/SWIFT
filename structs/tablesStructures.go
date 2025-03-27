package structs

type TableRow struct {
	Name string
	Data_type string
	Addition string
}

type Table struct {
	Name string
	Rows []TableRow
	Addition string
}

var Tables = []Table{
	{"countries", []TableRow{{"iso2", "CHAR(2)", " PRIMARY KEY"},
				{"name", "VARCHAR(20)", ""},
				{"time_zone", "VARCHAR(20)", ""},
			}, ""}, 
	{"headquarters", []TableRow{{"swift", "CHAR(11)", " PRIMARY KEY"},
				{"name", "VARCHAR(50)", ""},
				{"address", "VARCHAR(255)", ""},
				{"town", "VARCHAR(20)", ""},
				{"country", "CHAR(2)", ""},
			}, ` CONSTRAINT fk_headquarters_countries
				FOREIGN KEY (country) REFERENCES countries(iso2) 
				ON DELETE CASCADE`}, 
	{"branches", []TableRow{{"swift", "CHAR(11)", " PRIMARY KEY"},
				{"headquarter", "CHAR(11)", ""},
				{"name", "VARCHAR(50)", ""},
				{"address", "VARCHAR(255)", ""},
				{"town", "VARCHAR(20)", ""},
				{"country", "CHAR(2)", ""},
			}, ` CONSTRAINT fk_branches_countries
				FOREIGN KEY (country) REFERENCES countries(iso2) 
				ON DELETE CASCADE,  
				FOREIGN KEY (headquarter) REFERENCES headquarters(swift)
				ON DELETE SET NULL`},
}
