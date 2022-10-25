package recommendModule

// CompanyWeight 公司规模权重
var companyWeight = map[string]float64{
	"0-20":      1.0,
	"20-99":     1.2,
	"100-499":   1.4,
	"500-999":   1.6,
	"1000-9999": 2.0,
	"10000以上":   2.4,
}
