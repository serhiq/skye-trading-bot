package worker

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	repositoryProduct "github.com/serhiq/skye-trading-bot/internal/repository/product"
	"github.com/serhiq/skye-trading-bot/pkg/restoClient"
	evotorrestogo "github.com/softc24/evotor-resto-go"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

const LOG_TAG = "[sync_worker]"

type SyncWorker struct {
	r         *repositoryProduct.Repository
	evoClient *restoClient.RestoClient
	client    *resty.Client
	scheduler *gocron.Scheduler
}

func (s *SyncWorker) Start() error {

	resp, err := s.evoClient.GetMenu()
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	//clearOldPreviews()

	result := []*repositoryProduct.Product{}

	for _, product := range resp {

		if !isProductSupported(product) {
			continue
		}

		menuItem := &repositoryProduct.Product{
			Name: product.Name,
			//StoreID:    item.StoreID,
			UUID:        product.UUID,
			ParentUUID:  product.ParentUUID,
			Group:       product.Group,
			Image:       "",
			MeasureName: product.MeasureName,
			Price:       uint64(product.Price),
		}

		if product.ImageURL != "" {
			thumbnail, err := createThumbnail(s.client, product.UUID, product.ImageURL)
			if err != nil {
				logger.SugaredLogger.Errorf("SyncWorker: createThumbnail, %s", err)
			} else {
				menuItem.Image = thumbnail
			}
		}
		result = append(result, menuItem)
	}

	err = s.r.ImportMenu(result)
	if err != nil {
		return fmt.Errorf("sync: err, %s", err)
	}
	return nil
}

func isProductSupported(p evotorrestogo.MenuItem) bool {

	if p.IsUnavailable {
		return false
	}

	if p.Group {
		return true
	}

	if !p.Group && p.Type == "NORMAL" {
		return true
	}
	return true
}

func createThumbnail(r *resty.Client, fileName string, url string) (string, error) {
	tmpImage := filepath.Join(config.TempPatch, fileName)
	defer os.Remove(tmpImage)

	previewImage := filepath.Join(config.PreviewCachePatch, fileName)

	resp, err := r.R().
		SetOutput(tmpImage).
		Get(url)

	if err != nil {
		return "", fmt.Errorf("external_api: error get image, url=%s, err=%s", url, err)
	}

	if !resp.IsSuccess() {
		return "", fmt.Errorf("sync: error get image, , url=%s, code == %d", url, resp.StatusCode())
	}

	err = resizeTmpFile(tmpImage, previewImage)
	if err != nil {
		return "", fmt.Errorf("sync: error resize image, %s", err)
	}

	return previewImage, nil
}

func clearOldPreviews() {
	err := os.RemoveAll(config.PreviewCachePatch)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	fmt.Println("Previews cleared successfully")
}

func (s *SyncWorker) Stop() {
	if s.scheduler != nil {
		s.scheduler.Clear()
	} else {
		println("scheduler nil")
	}
}

func (s *SyncWorker) EnqueueUniquePeriodicWork() {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(int(10)).Minutes().Do(updateMenu, s)
	if err != nil {
		logger.SugaredLogger.Errorf("Error: updateMenu: %s", err)
		//fmt.Println()
	}

	s.scheduler = scheduler

	scheduler.StartAsync()
}

var updateMenu = func(s *SyncWorker) {
	err := s.Start()
	if err != nil {
		logger.SugaredLogger.Errorf("sync: err, %s", err)
	}
}

func New(r *repositoryProduct.Repository, client *restoClient.RestoClient, resty *resty.Client) *SyncWorker {
	return &SyncWorker{
		r:         r,
		evoClient: client,
		client:    resty,
	}
}

func resizeTmpFile(tmpDest string, filename string) error {
	file, err := os.Open(tmpDest)

	img, err := jpeg.Decode(file)
	if err != nil {
		return err
	}

	err = file.Close()
	if err != nil {
		return err
	}

	m := resize.Resize(0, 300, img, resize.Lanczos3)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, m, nil)
}
