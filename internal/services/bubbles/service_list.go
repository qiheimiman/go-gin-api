package bubbles

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql/bubbles1000"
)

func (s *service) List(ctx core.Context) (info []*bubbles1000.Bubbles1000, err error) {

	err = s.db.GetDbR().Table("bubbles1000").Where("`rank` <= ?", 30).Find(&info).Error
	if err != nil {
		return nil, err
	}

	return
}
