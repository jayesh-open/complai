package categorizer

import "github.com/complai/complai/services/go/gst-service/internal/domain"

func Categorize(e *domain.SalesRegisterEntry) string {
	if e.DocumentType == "CRN" || e.DocumentType == "DBN" {
		if e.BuyerGSTIN != "" && e.BuyerGSTIN != "URP" {
			return domain.SectionCDNR
		}
		return domain.SectionCDNUR
	}

	if e.SupplyType == "EXP" {
		return domain.SectionEXP
	}

	if isNILRated(e) {
		return domain.SectionNIL
	}

	if e.SupplyType == "B2CL" {
		return domain.SectionB2CL
	}

	if e.SupplyType == "B2CS" {
		return domain.SectionB2CS
	}

	if e.BuyerGSTIN != "" && e.BuyerGSTIN != "URP" {
		return domain.SectionB2B
	}

	return domain.SectionB2CS
}

func isNILRated(e *domain.SalesRegisterEntry) bool {
	return e.CGSTRate.IsZero() && e.SGSTRate.IsZero() && e.IGSTRate.IsZero() &&
		e.TaxableValue.IsPositive()
}
