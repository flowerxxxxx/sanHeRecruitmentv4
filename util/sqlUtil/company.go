package sqlUtil

//type companyScale struct {
//	CompanyScLevels map[string]int
//	Locker          *sync.Mutex
//}
//
//var CompanyScaleList = companyScale{
//	CompanyScLevels: map[string]int{
//		"0-20":      1,
//		"20-99":     2,
//		"100-499":   3,
//		"500-999":   4,
//		"1000-9999": 5,
//		"10000以上":   6,
//	},
//}

var CompanyScLevels = map[string]int{
	"0-20":      1,
	"20-99":     2,
	"100-499":   3,
	"500-999":   4,
	"1000-9999": 5,
	"10000以上":   6,
}
