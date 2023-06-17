package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	elog "github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"github.com/vvinokurshin/DBCourseVK/cmd/router"
	"github.com/vvinokurshin/DBCourseVK/cmd/server"
	forumDelivery "github.com/vvinokurshin/DBCourseVK/internal/forum/delivery"
	forumRepo "github.com/vvinokurshin/DBCourseVK/internal/forum/repository/postgres"
	forumUC "github.com/vvinokurshin/DBCourseVK/internal/forum/usecase"
	postDelivery "github.com/vvinokurshin/DBCourseVK/internal/post/delivery"
	postRepo "github.com/vvinokurshin/DBCourseVK/internal/post/repository/postgres"
	postUC "github.com/vvinokurshin/DBCourseVK/internal/post/usecase"
	serviceDelivery "github.com/vvinokurshin/DBCourseVK/internal/service/delivery"
	serviceRepo "github.com/vvinokurshin/DBCourseVK/internal/service/repository/postgres"
	serviceUC "github.com/vvinokurshin/DBCourseVK/internal/service/usecase"
	threadDelivery "github.com/vvinokurshin/DBCourseVK/internal/thread/delivery"
	threadRepo "github.com/vvinokurshin/DBCourseVK/internal/thread/repository/postgres"
	threadUC "github.com/vvinokurshin/DBCourseVK/internal/thread/usecase"
	userDelivery "github.com/vvinokurshin/DBCourseVK/internal/user/delivery"
	userRepo "github.com/vvinokurshin/DBCourseVK/internal/user/repository/postgres"
	userUC "github.com/vvinokurshin/DBCourseVK/internal/user/usecase"
	"log"
)

var connString = fmt.Sprintf("host=localhost port=5432 user=valeriy password=valeriy_pw sslmode=disable dbname=db_forum")

func main() {
	db, err := sqlx.Connect("postgres", connString)
	if err != nil {
		log.Fatal(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)

	userRepository := userRepo.NewRepository(db)
	forumRepository := forumRepo.NewRepository(db)
	threadRepository := threadRepo.NewRepository(db)
	postRepository := postRepo.NewRepository(db)
	serviceRepository := serviceRepo.NewRepository(db)

	userUseCase := userUC.NewUseCase(userRepository)
	forumUseCase := forumUC.NewUseCase(forumRepository, userRepository)
	threadUseCase := threadUC.NewUseCase(threadRepository, forumRepository, userRepository)
	postUseCase := postUC.NewUseCase(postRepository, threadRepository, userRepository, forumRepository)
	serviceUseCase := serviceUC.NewUseCase(serviceRepository)

	userHandler := userDelivery.NewDelivery(userUseCase)
	forumHandler := forumDelivery.NewDelivery(forumUseCase)
	threadHandler := threadDelivery.NewDelivery(threadUseCase)
	postHandler := postDelivery.NewDelivery(postUseCase)
	serviceHandler := serviceDelivery.NewDelivery(serviceUseCase)

	e := echo.New()

	e.Logger.SetHeader(`time=${time_rfc3339} level=${level} prefix=${prefix} ` +
		`file=${short_file} line=${line} message:`)
	e.Logger.SetLevel(elog.INFO)

	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `time=${time_custom} remote_ip=${remote_ip} ` +
			`host=${host} method=${method} uri=${uri} user_agent=${user_agent} ` +
			`status=${status} error="${error}" ` +
			`bytes_in=${bytes_in} bytes_out=${bytes_out}` + "\n",
		CustomTimeFormat: "2006-01-02 15:04:05",
	}))

	router.AddRoutes(e, userHandler, forumHandler, threadHandler, postHandler, serviceHandler)

	s := server.NewServer(e)

	if err := s.Start(); err != nil {
		e.Logger.Fatal(err)
	}
}
