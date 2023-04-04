package server

import (
	"fmt"
	r "github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/performer"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	repositoryChat "github.com/serhiq/skye-trading-bot/internal/repository/chat"
	repositoryOrder "github.com/serhiq/skye-trading-bot/internal/repository/order"
	repositoryProduct "github.com/serhiq/skye-trading-bot/internal/repository/product"
	"github.com/serhiq/skye-trading-bot/internal/worker"
	"github.com/serhiq/skye-trading-bot/pkg/restoClient"
	"github.com/serhiq/skye-trading-bot/pkg/store/mysql"
	"log"
	"time"
)

const (
	longPollTimeout = 15
)

type Server struct {
	cfg       config.Config
	store     *mysql.Store
	delivery  *performer.Performer
	startFunc []func()
	stopFunc  []func()
}

func Serve(cfg config.Config) (*Server, error) {
	var server = &Server{
		cfg:       cfg,
		store:     nil,
		delivery:  nil,
		startFunc: nil,
		stopFunc:  nil,
	}

	err := server.initApp()
	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) Start() error {
	fmt.Println("Server is starting...")

	for _, start := range s.startFunc {
		start()
	}

	return nil
}

func (s *Server) Stop() {
	for _, stop := range s.stopFunc {
		stop()
	}
}

func (s *Server) initApp() (err error) {

	store, err := mysql.New(s.cfg.DBConfig)

	if err != nil {
		panic(err)
	}

	s.store = store

	err = store.Db.AutoMigrate(&repositoryChat.Chat{}, &repositoryProduct.Product{}, &repositoryOrder.Order{}, &repositoryOrder.OrderPosition{})
	if err != nil {
		return err
	}

	s.addStopDelegate(func() {
		log.Println("db stop func")

		db, err := s.store.Db.DB()
		if err != nil {
			log.Printf("database: error close database, %s", err)
			return
		}
		err = db.Close()
		if err != nil {
			log.Printf("database: error close database, %s", err)
			return
		}
		log.Print("database: close")
	})

	client := r.New()
	evoClient := restoClient.New(client, &restoClient.Options{
		Auth:    s.cfg.RestaurantAPI.Auth,
		Store:   s.cfg.RestaurantAPI.Store,
		BaseUrl: s.cfg.RestaurantAPI.BaseURL,
	})

	var repoProduct = repositoryProduct.New(s.store.Db)
	var repoChat = repositoryChat.New(s.store.Db)
	var repoOrder = repositoryOrder.New(evoClient, s.store.Db)

	syncWorker := worker.New(repoProduct, evoClient, r.New())

	s.addStartDelegate(func() {
		logger.SugaredLogger.Infow(worker.LOG_TAG + "  start")
		syncWorker.EnqueueUniquePeriodicWork()
	})

	s.addStopDelegate(func() {
		logger.SugaredLogger.Infow(worker.LOG_TAG + "  stop")
		syncWorker.Stop()
	})

	sBot, err := performer.New(performer.Options{
		Token: s.cfg.Telegram.Token,
	}, repoProduct, repoChat, repoOrder)

	if err != nil {
		logger.SugaredLogger.Panicw("initApp: cannot initialize Bot", "err", err)
	}

	s.delivery = sBot

	u := tgbotapi.NewUpdate(0)
	u.Timeout = longPollTimeout

	updates := s.delivery.App.Bot.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	s.addStartDelegate(func() {
		logger.SugaredLogger.Infof("Bot online %s", s.delivery.App.Bot.Self.UserName)
		for update := range updates {
			go s.delivery.Dispatch(&update)
		}
	})

	s.addStopDelegate(func() {
		logger.SugaredLogger.Info("Bot is stopping...")
		s.delivery.App.Bot.StopReceivingUpdates()
	})

	return nil
}

func (s *Server) addStartDelegate(delegate func()) {
	s.startFunc = append(s.startFunc, delegate)
}

func (s *Server) addStopDelegate(delegate func()) {
	s.stopFunc = append(s.stopFunc, delegate)
}
