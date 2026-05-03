package domain

type TableDef struct {
	Part        int
	Table       string
	Description string
}

var GSTR9Tables = []TableDef{
	{1, "4A", "B2B supplies (taxable)"},
	{1, "4B", "B2C supplies (taxable)"},
	{1, "4C", "Exports (with payment)"},
	{1, "4D", "Exports (without payment / SEZ)"},
	{1, "4E", "Non-GST outward supplies"},
	{2, "5A", "Imports (goods)"},
	{2, "5B", "Imports (services)"},
	{2, "5C", "Inward supplies under reverse charge"},
	{2, "5D", "Inward supplies from ISD"},
	{2, "5E", "All other inward supplies"},
	{3, "6A", "ITC availed — imports"},
	{3, "6B", "ITC availed — inward RCM"},
	{3, "6C", "ITC availed — ISD"},
	{3, "6D", "ITC availed — all other"},
	{3, "6E", "ITC reversed"},
	{3, "6F", "Net ITC available"},
	{4, "9", "Tax paid (cash + ITC)"},
	{5, "10-14", "Prior-year amendments"},
	{6, "15-19", "HSN-wise summary of outward + inward"},
}

func ReturnPeriodsForFY(fy string) []string {
	parts := []string{}
	if len(fy) < 7 {
		return parts
	}
	startYear := fy[:4]
	months := []string{
		"04", "05", "06", "07", "08", "09",
		"10", "11", "12",
	}
	for _, m := range months {
		parts = append(parts, startYear+m)
	}
	endYear := fy[5:]
	for _, m := range []string{"01", "02", "03"} {
		parts = append(parts, endYear+m)
	}
	return parts
}
