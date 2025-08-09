package news

import "time"

// News 新闻表
//
//go:generate gormgen -structs News -input .
type News struct {
	Id          int32     // 自增主键
	Title       string    // 新闻标题
	SubTitle    string    // 新闻副标题
	WebLink     string    // 新闻网页链接
	BinanceId   string    // Binance平台新闻ID
	PublishTime time.Time `gorm:"time"` // 新闻发布时间
	CreatedAt   time.Time `gorm:"time"` // 记录创建时间
	Type        int32     // 类型 1-币安新闻
}
