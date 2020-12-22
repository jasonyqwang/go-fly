package cmd

import (
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go-fly/config"
	"go-fly/docs"
	"go-fly/router"
	"go-fly/tools"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	port        string
	tcpport     string
	daemon      bool
	GoflyConfig *config.Config
)
var serverCmd = &cobra.Command{
	Use:     "server",
	Short:   "example:go-fly server port 8081",
	Example: "go-fly server -c config/",
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	serverCmd.PersistentFlags().StringVarP(&port, "port", "p", "8081", "监听端口号")
	serverCmd.PersistentFlags().StringVarP(&tcpport, "tcpport", "t", "8082", "监听tcp端口号")
	serverCmd.PersistentFlags().BoolVarP(&daemon, "daemon", "d", false, "是否为守护进程模式")
}
func run() {
	if daemon == true {
		if os.Getppid() != 1 {
			// 将命令行参数中执行文件路径转换成可用路径
			filePath, _ := filepath.Abs(os.Args[0])
			cmd := exec.Command(filePath, os.Args[1:]...)
			// 将其他命令传入生成出的进程
			cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Start() // 开始执行新进程，不等待新进程退出
			os.Exit(0)
		}
	}

	baseServer := "0.0.0.0:" + port
	//tcpBaseServer := "0.0.0.0:"+tcpport
	log.Println("start server...\r\ngo：http://" + baseServer)
	engine := gin.Default()
	engine.LoadHTMLGlob("static/html/*")
	engine.Static("/static", "./static")

	//性能监控
	pprof.Register(engine)

	//记录日志
	engine.Use(tools.LoggerToFile())
	router.InitViewRouter(engine)
	router.InitApiRouter(engine)

	//文档服务
	docs.SwaggerInfo.Title = "GO-FLY接口文档"
	docs.SwaggerInfo.Description = "go-fly即时通讯web客服管理系统 , 测试账户:kefu2 测试密码:123 类型:kefu"
	docs.SwaggerInfo.Version = "0.0.7"
	//docs.SwaggerInfo.Host = "127.0.0.1:"+port
	docs.SwaggerInfo.Host = "gofly.sopans.com"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"https"}
	engine.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//logFile, _ := os.OpenFile("./fatal.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	//tools.RedirectStderr(logFile)

	//tcp服务
	//go controller.NewTcpServer(tcpBaseServer)
	engine.Run(baseServer)
}
