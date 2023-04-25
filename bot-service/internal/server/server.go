package server

import (
	"fmt"
	r "github.com/go-resty/resty/v2"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"
	"github.com/serhiq/skye-trading-bot/internal/bot/performer"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	orderController "github.com/serhiq/skye-trading-bot/internal/contorller/order"
	productController "github.com/serhiq/skye-trading-bot/internal/contorller/product"
	"github.com/serhiq/skye-trading-bot/internal/delivery"
	evoClient "github.com/serhiq/skye-trading-bot/internal/delivery/evotor"
	"github.com/serhiq/skye-trading-bot/internal/delivery/file"
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
	cfg               config.Config
	store             *mysql.Store
	delivery          *performer.Performer
	startFunc         []func()
	stopFunc          []func()
	orderController   contorller.OrderController
	productController contorller.ProductController
	sessionRepository *repositoryChat.Repository
}

func Serve(cfg config.Config) (*Server, error) {
	var s = &Server{
		cfg:       cfg,
		store:     nil,
		delivery:  nil,
		startFunc: nil,
		stopFunc:  nil,
	}

	for _, init := range []func() error{
		s.initDb,
		s.initOrderController,
		s.initProductController,
		s.initSessionRepository,
		s.initBot,
	} {
		if err := init(); err != nil {
			return nil, errors.Wrap(err, "serve failed")
		}
	}
	return s, nil
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

func (s *Server) initSessionRepository() error {

	s.sessionRepository = repositoryChat.New(s.store.Db)
	return nil
}

func (s *Server) initProductController() error {

	productRepo := repositoryProduct.New(s.store.Db)
	client := r.New()

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
	case config.FileKind:
		productProvider = file.FileClient{}
	default:
		return fmt.Errorf("unknown ProductApiKind: %s", s.cfg.ProductAPI.Kind)
	}

	productCtrl := productController.New(productRepo, productProvider, r.New())
	s.addStartDelegate(func() {
		//logger.SugaredLogger.Info(worker.LOG_TAG + "  start")
		productCtrl.StartSync()
	})

	s.addStopDelegate(func() {
		productCtrl.StopSync()
		//logger.SugaredLogger.Info(worker.LOG_TAG + "  stop")
	})

	s.productController = productCtrl
	return nil
}

func (s *Server) initOrderController() error {

	client := r.New()
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
	case config.FileOrderKind:
		orderProvider = file.FileClient{}

	default:
		return fmt.Errorf("unknown order API %s", s.cfg.ProductAPI.Kind)
	}

	s.orderController = orderController.New(orderRepo, orderProvider, orderRepo)
	return nil
}

func (s *Server) initDb() error {
	store, err := mysql.New(s.cfg.DBConfig)

	if err != nil {
		return err
	}

	s.store = store

	s.addStopDelegate(func() {
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
	return err
}

func (s *Server) addStartDelegate(delegate func()) {
	s.startFunc = append(s.startFunc, delegate)
}

func (s *Server) addStopDelegate(delegate func()) {
	s.stopFunc = append(s.stopFunc, delegate)
}

func (s *Server) initBot() error {
	sBot, err := performer.New(performer.Options{
		Token: s.cfg.Telegram.Token,
	}, s.productController, s.sessionRepository, s.orderController)

	if err != nil {
		return errors.Wrap(err, "cannot initialize Bot")
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
