package bubbles1000

// Bubbles1000 加密泡泡排名
//
//go:generate gormgen -structs Bubbles1000 -input .
type Bubbles1000 struct {
	Id             int32   //
	Name           string  //
	Slug           string  //
	Symbol         string  //
	Dominance      float64 //
	Image          string  //
	Rank           int32   //
	Price          float64 //
	Marketcap      int64   //
	Volume         int64   //
	CgId           string  //
	Symbols        string  //
	Performance    string  //
	RankDiffs      string  //
	ExchangePrices string  //
}
