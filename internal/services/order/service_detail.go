package order

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/order"
)

type SearchOneData struct {
	Id      int32  // 用户ID
	OrderNo string // 订单编号
}

func (s *service) Detail(ctx core.Context, searchOneData *SearchOneData) (info *order.Order, err error) {

	err = s.db.GetDbR().Table("order").Where("id = ?", searchOneData.Id).First(&info).Error
	if err != nil {
		return nil, err
	}

	return
}
