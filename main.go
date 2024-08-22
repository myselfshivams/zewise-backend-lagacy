package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/Kirisakiii/neko-micro-blog-backend/configs"
	"github.com/Kirisakiii/neko-micro-blog-backend/consts"
	"github.com/Kirisakiii/neko-micro-blog-backend/controllers"
	"github.com/Kirisakiii/neko-micro-blog-backend/loggers"
	"github.com/Kirisakiii/neko-micro-blog-backend/middlewares"
	"github.com/Kirisakiii/neko-micro-blog-backend/models"
	"github.com/Kirisakiii/neko-micro-blog-backend/rontines"
	"github.com/Kirisakiii/neko-micro-blog-backend/services"
	"github.com/Kirisakiii/neko-micro-blog-backend/stores"
)

var (
	logger            *logrus.Logger
	cfg               *configs.Config
	db                *gorm.DB
	storeFactory      *stores.Factory
	controllerFactory *controllers.Factory
	middlewareFactory *middlewares.Factory
)

func init() {
	// logger
	logger = loggers.NewLogger()
	logger.Infoln("正在执行程序初始化...")

	var err error

	// 加载配置文件
	cfg, err = configs.NewConfig()
	if err != nil {
		logger.Panicln(err.Error())
	}

	// 设置日志等级
	var (
		logLevel logrus.Level
		logMode  gormLogger.LogLevel
	)
	switch cfg.Env.Type {
	case "development":
		logLevel = logrus.DebugLevel
		logMode = gormLogger.Error
	case "production":
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	default:
		logLevel = logrus.InfoLevel
		logMode = gormLogger.Silent
	}

	// 设置logrus日志等级
	logger.SetLevel(logLevel)
	logger.Debugln("日志记录等级设定为:", strings.ToUpper(logLevel.String()))

	// 连接数据库
	logger.Debugln("尝试连接至数据库...")
	db, err = gorm.Open(
		postgres.Open(fmt.Sprintf(
			"postgres://%s:%s@%s:%d/%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.DBName,
		)),
		&gorm.Config{
			Logger: gormLogger.Default.LogMode(logMode),
		},
	)
	if err != nil {
		logger.Panicln(err.Error())
	}
	logger.Debugln("数据库连接成功")

	// 迁移模型
	logger.Debugln("正在迁移数据表模型...")
	err = models.Migrate(db)
	if err != nil {
		logger.Panicln("迁移数据库模型失败：", err.Error())
	}

	// 建立数据访问层工厂
	storeFactory = stores.NewFactory(db)

	// 建立控制器层工厂
	controllerFactory = controllers.NewFactory(
		services.NewFactory(storeFactory),
	)

	// 建立中间件工厂
	middlewareFactory = middlewares.NewFactory(storeFactory)
}

func main() {
	// 创建定时任务
	crontab := cron.New()
	crontab.AddJob(
		"@every 5min",
		cron.NewChain(
			cron.SkipIfStillRunning(cron.DefaultLogger),
		).Then(
			rontines.NewAvatarCleanerJob(logger, db),
		),
	)
	crontab.Start()

	// 创建 fiber 实例
	var fiberConfig fiber.Config
	// 如果是生产环境，则开启 Prefork
	if cfg.Env.Type == "production" {
		fiberConfig = fiber.Config{
			Prefork: true,
		}
	}
	fiberConfig.BodyLimit = 10 * consts.POST_IMAGE_MAX_FILE_SIZE
	app := fiber.New(fiberConfig)

	// 设置中间件
	app.Use(fiberLogger.New(fiberLogger.Config{
		Format: "[${time}][${latency}][${status}][${method}] ${path}\n",
	}))
	app.Use(compress.New(compress.Config{
		Level: cfg.Compress.Level,
	}))

	// Auth 中间件
	authMiddleware := middlewareFactory.NewTokenAuthMiddleware()

	// 静态资源路由
	resource := app.Group("/resources")
	// 头像资源路由
	resource.Static("/avatar", "./public/avatars", fiber.Static{
		Compress: true,
	})
	// 博文图片资源路由
	resource.Static("/image", "./public/images", fiber.Static{
		Compress: true,
	})

	// api 路由
	api := app.Group("/api")

	// User 路由
	userController := controllerFactory.NewUserController()
	user := api.Group("/user")
	user.Get("/profile", userController.NewProfileHandler())                                             // 查询用户信息
	user.Post("/register", userController.NewRegisterHandler())                                          // 用户注册
	user.Post("/login", userController.NewLoginHandler())                                                // 用户登录
	user.Post("/upload-avatar", authMiddleware.NewMiddleware(), userController.NewUploadAvatarHandler()) // 上传头像
	user.Post("/update-psw", userController.NewUpdatePasswordHandler())                                  // 修改密码
	user.Post("/edit", authMiddleware.NewMiddleware(), userController.NewUpdateProfileHandler())         // 修改用户资料

	//post 路由
	postController := controllerFactory.NewPostController()
	post := api.Group("/post")
	post.Post("/new", authMiddleware.NewMiddleware(), postController.NewCreatePostHandler())            // 创建文章
	post.Get("/list", postController.NewPostListHandler())                                              // 获取文章列表
	post.Get("/detail", postController.NewPostDetailHandler())                                          // 获取文章信息
	post.Delete("/delete/:post", authMiddleware.NewMiddleware(), postController.NewDeletePostHandler()) // 删除文章

	// Comment 路由
	commentController := controllerFactory.NewCommentController()
	comment := api.Group("/comment")
	comment.Post("/new", authMiddleware.NewMiddleware(), commentController.NewCreateCommentHandler(storeFactory.NewPostStore(), storeFactory.NewUserStore())) // 创建评论
	comment.Post("/edit", authMiddleware.NewMiddleware(), commentController.NewUpdateCommentHandler())                                                        // 修改评论
	comment.Post("/delete", authMiddleware.NewMiddleware(), commentController.DeleteCommentHandler())                                                         // 删除评论
	comment.Get("/list", commentController.NewCommentListHandler())                                                                                           // 获取评论列表
	comment.Get("/detail", commentController.NewCommentDetailHandler())

	// 启动服务器
	log.Fatal(app.Listen(fmt.Sprintf("%s:%d", cfg.Database.Host, cfg.Server.Port)))
}
