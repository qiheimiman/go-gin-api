package order

import (
	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/order"
	"github.com/xinliangnote/go-gin-api/pkg/hash"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// Create 创建订单
	// @Tags API.order
	// @Router /api/order/create [post]
	Create() core.HandlerFunc

	// Cancel 取消订单
	// @Tags API.order
	// @Router /api/order/cancel [post]
	Cancel() core.HandlerFunc

	// Detail 查看详情
	// @Tags API.order
	// @Router /api/order/{id} [get]
	Detail() core.HandlerFunc
}

type handler struct {
	logger       *zap.Logger
	db           mysql.Repo
	cache        redis.Repo
	hashids      hash.Hash
	orderService order.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:       logger,
		db:           db,
		cache:        cache,
		hashids:      hash.New(configs.Get().HashIds.Secret, configs.Get().HashIds.Length),
		orderService: order.New(db, cache),
	}
}

func (h *handler) i() {}
