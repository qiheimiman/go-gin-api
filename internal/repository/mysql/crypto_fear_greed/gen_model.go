package crypto_fear_greed

// CryptoFearGreed 加密市场贪婪恐惧指数
//
//go:generate gormgen -structs CryptoFearGreed -input .
type CryptoFearGreed struct {
	Id            int32  //
	Date          string // 所属日期
	CurrentValue  int32  //
	BearishValue  int32  //
	BullishValue  int32  //
	LastWeekValue int32  //
}
