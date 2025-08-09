package news

import (
	"time"

	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/news"
)

func (s *service) List(ctx core.Context) (info []*news.News, err error) {

	err = s.db.GetDbR().Table("news").Where("publish_time >= ?", time.Now().AddDate(0, 0, -7).Format("2006-01-02 15:04:05")).Find(&info).Error
	if err != nil {
		return nil, err
	}

	return
}
