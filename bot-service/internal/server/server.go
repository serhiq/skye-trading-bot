package server

import (
	"fmt"
	r "github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/serhiq/skye-trading-bot/internal/bot/performer"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	orderController "github.com/serhiq/skye-trading-bot/internal/contorller/order"
	productController "github.com/serhiq/skye-trading-bot/internal/contorller/product"
	"github.com/serhiq/skye-trading-bot/internal/delivery"
	evoClient "github.com/serhiq/skye-trading-bot/internal/delivery/evotor"
	orderClient "github.com/serhiq/skye-trading-bot/internal/delivery/orders"
	restoranClient "github.com/serhiq/skye-trading-bot/internal/delivery/resto"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	repositoryChat "github.com/serhiq/skye-trading-bot/internal/repository/chat"
	repositoryOrder "github.com/serhiq/skye-trading-bot/internal/repository/order"
	repositoryProduct "github.com/serhiq/skye-trading-bot/internal/repository/product"
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
	/////////////////////////////////////////////////////////////////////////////
	orderRepo := repositoryOrder.New(s.store.Db)

	var orderProvider contorller.OrderSender
	switch s.cfg.OrderAPI.Kind {
	case config.OrderAPIKind:
		//logger.SugaredLogger.Infow("Using Evotor Product API")
		orderProvider = orderClient.New(client, &orderClient.Options{
			Auth:    s.cfg.OrderAPI.Auth,
			BaseUrl: s.cfg.OrderAPI.BaseURL,
			Store:   s.cfg.OrderAPI.Store,
		},
		)
	case config.RestoOrderAPIKind:
		//logger.SugaredLogger.Infow("Using Restoran Order API")

		orderProvider = restoranClient.New(client, &restoranClient.Options{
			Auth:    s.cfg.ProductAPI.Auth,
			BaseUrl: s.cfg.ProductAPI.BaseURL,
			Store:   s.cfg.ProductAPI.Store,
		},
		)
	default:
		return fmt.Errorf("unknown order APIr: %s", s.cfg.ProductAPI.Kind)
	}

	var orderCtrl = orderController.New(orderRepo, orderProvider)

	//////////////////////////////////////////////////////////////////
	productRepo := repositoryProduct.New(s.store.Db)

	var productProvider delivery.ProductProvider

	switch s.cfg.ProductAPI.Kind {
	case config.EvotorAPIKind:
		//logger.SugaredLogger.Info("Using Evotor Product API")
		productProvider = evoClient.New(client, &evoClient.Options{
			Auth:     s.cfg.ProductAPI.Auth,
			BaseUrl:  s.cfg.ProductAPI.BaseURL,
			Store:    s.cfg.ProductAPI.Store,
			MenuUuid: s.cfg.ProductAPI.MenuUuid,
		},
		)
	case config.RestoAPIKind:
		//logger.SugaredLogger.Info("Using Resto Product API")

		productProvider = restoranClient.New(client, &restoranClient.Options{
			Auth:    s.cfg.ProductAPI.Auth,
			BaseUrl: s.cfg.ProductAPI.BaseURL,
			Store:   s.cfg.ProductAPI.Store,
		},
		)
	default:
		return fmt.Errorf("unknown ProductApiKind: %s", s.cfg.ProductAPI.Kind)
	}

	productCtrl := productController.New(productRepo, productProvider, r.New())
	var repoChat = repositoryChat.New(s.store.Db)
	s.addStartDelegate(func() {
		//logger.SugaredLogger.Info(worker.LOG_TAG + "  start")
		productCtrl.StartSync()
	})

	s.addStopDelegate(func() {
		productCtrl.StopSync()
		//logger.SugaredLogger.Info(worker.LOG_TAG + "  stop")
	})

	sBot, err := performer.New(performer.Options{
		Token: s.cfg.Telegram.Token,
	}, productCtrl, repoChat, orderCtrl)

	if err != nil {
		logger.SugaredLogger.Panicw("initApp: cannot initialize Bot", "err", err)
	}

	s.delivery = sBot

	u := tgbotapi.NewUpdate(0)
	u.Timeout = longPollTimeout

	updates := s.delivery.App.Bot.Api.GetUpdatesChan(u)
	time.Sleep(time.Millisecond * 500)
	updates.Clear()

	s.addStartDelegate(func() {
		logger.SugaredLogger.Infof("Bot online %s", s.delivery.App.Bot.Api.Self.UserName)
		for update := range updates {
			go s.delivery.Dispatch(&update)
		}
	})

	s.addStopDelegate(func() {
		logger.SugaredLogger.Info("Bot is stopping...")
		s.delivery.App.Bot.Api.StopReceivingUpdates()
	})

	return nil
}

func (s *Server) addStartDelegate(delegate func()) {
	s.startFunc = append(s.startFunc, delegate)
}

func (s *Server) addStopDelegate(delegate func()) {
	s.stopFunc = append(s.stopFunc, delegate)
}
