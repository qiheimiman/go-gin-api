package bubbles

import (
	"github.com/xinliangnote/go-gin-api/configs"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/repository/mysql"
	"github.com/xinliangnote/go-gin-api/internal/repository/redis"
	"github.com/xinliangnote/go-gin-api/internal/services/bubbles"
	"github.com/xinliangnote/go-gin-api/pkg/hash"

	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {
	i()

	// List 获取最新排名前1000加密货币价格
	// @Tags PUBLIC_API.bubbles
	// @Router /public_api/bubbles/list [get]
	List() core.HandlerFunc
}

type handler struct {
	logger         *zap.Logger
	db             mysql.Repo
	cache          redis.Repo
	hashids        hash.Hash
	bubblesService bubbles.Service
}

func New(logger *zap.Logger, db mysql.Repo, cache redis.Repo) Handler {
	return &handler{
		logger:         logger,
		db:             db,
		cache:          cache,
		hashids:        hash.New(configs.Get().HashIds.Secret, configs.Get().HashIds.Length),
		bubblesService: bubbles.New(db, cache),
	}
}

func (h *handler) i() {}
