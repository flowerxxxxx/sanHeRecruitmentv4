package sqlUtil

func DealSliToSql(love []string) string {
	sql := "("
	for _, v := range love {
		sql = sql + "'" + v + "',"
	}
	sql = sql[:len(sql)-1] + ")"
	return sql
}

//func main(){
//	love := []string{"工商类","管理类","管理类"}
//	fmt.Println(love)
//
//	sqlx := dealSliToSql(love)
//	fmt.Println(sqlx)
//	//lei := "('工商类','管理类')"
//	sql := "SELECT * FROM `articles` where type1 in "+ sqlx
//	fmt.Println(sql)
//	fimilar := "SELECT * FROM `articles` where type1 in ('工商类','管理类')"
//	fmt.Println(fimilar)
//}
