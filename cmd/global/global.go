package global

import (
	"net/http"

	"github.com/xinliangnote/go-gin-api/cmd/config"
	"github.com/xinliangnote/go-gin-api/cmd/pkg/logger"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var (
	ServerConfig       *config.ServerConfig
	Logger             *logger.Logger
	HttpClient         *http.Client
	HttpLongClient     *http.Client
	DB                 *gorm.DB
	ErpSlaveDB         *gorm.DB
	ErpQualityDB       *gorm.DB
	ErpReportDB        *gorm.DB
	ErpReportHuoshanDB *gorm.DB
	WatchWorkDB        *gorm.DB
	FscBookingDB       *gorm.DB
	OpenApiDB          *gorm.DB
	ClickHouseDB       *gorm.DB
	Env                string
	Validate           *validator.Validate
)
