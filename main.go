/*
 * Copyright 2023 Nhan Nguyen.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/nhaancs/go-realworld/logger"
	"github.com/nhaancs/go-realworld/pgx"
	"net/http"
	"os"
)

func main() {
	var (
		log *logger.Logger
		ctx = context.Background()
	)

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT ******")
		},
	}

	traceIDFunc := func(ctx context.Context) string {
		return "trace-id-here"
	}
	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, "Realworld", traceIDFunc, events)

	log.Info(ctx, "Connect to database")
	db, err := pgx.Open(pgx.Config{
		User:         "postgres",
		Password:     "postgres",
		Host:         "localhost:5432",
		Name:         "postgres",
		MaxIdleConns: 2,
		MaxOpenConns: 0,
		DisableTLS:   true,
	})
	if err != nil {
		log.Error(ctx, "connecting to db: %w", err)
		os.Exit(1)
	}
	defer func() {
		log.Info(ctx, "shutdown", "status", "stopping database support", "host", "localhost:5432")
		db.Close()
	}()

	gin.SetMode(gin.ReleaseMode)
	ginEngine := gin.New()
	server := http.Server{
		Handler: ginEngine,
		Addr:    ":3000",
	}

	const prefix = "api"
	router := ginEngine.Group(prefix)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, struct {
			Status string
		}{"OK"})
	})

	router.GET("/db", func(c *gin.Context) {
		err := pgx.StatusCheck(ctx, db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
		c.JSON(http.StatusOK, struct {
			Status string
		}{"OK"})
	})

	log.Info(ctx, "Server start at port 3000")
	if err := server.ListenAndServe(); err != nil {
		log.Error(ctx, "Server error: %+v\n", err)
		os.Exit(1)
	}

}
