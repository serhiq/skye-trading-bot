package product

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	"github.com/serhiq/skye-trading-bot/internal/delivery"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	domainProduct "github.com/serhiq/skye-trading-bot/pkg/type/product"
	"image/jpeg"
	"os"
	"path/filepath"
	"time"
)

type ProductControllerImpl struct {
	productRepo contorller.ProductRepository
	//orderHistoryRepo OrderHistoryRepository
	productProvider delivery.ProductProvider
	scheduler       *gocron.Scheduler
	client          *resty.Client
}

func New(productRepo contorller.ProductRepository, productProvider delivery.ProductProvider, client *resty.Client) contorller.ProductController {
	return &ProductControllerImpl{
		productRepo:     productRepo,
		productProvider: productProvider,
		client:          client,
	}
}

func (c ProductControllerImpl) GetProductByParent(parentUuid string) ([]*domainProduct.Product, error) {
	return c.productRepo.GetProductByParent(parentUuid)

}

func (c ProductControllerImpl) GetProductByUuid(uuid string) (*domainProduct.Product, error) {
	return c.productRepo.GetProductByUuid(uuid)
}

func (c ProductControllerImpl) StartSync() {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(int(10)).Minutes().Do(updateMenu, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Error: updateMenu: %s", err)
		//fmt.Println()
	}

	c.scheduler = scheduler
	scheduler.StartAsync()
}

var updateMenu = func(s ProductControllerImpl) {
	err := s.start()
	if err != nil {
		logger.SugaredLogger.Errorf("sync: err, %s", err)
	}
}

func (c ProductControllerImpl) start() error {

	products, err := c.productProvider.GetProducts()

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	//clearOldPreviews()

	for _, product := range products {

		if product.Image != "" {
			thumbnail, err := createThumbnail(c.client, product.UUID, product.Image)
			if err != nil {
				logger.SugaredLogger.Errorf("SyncWorker: createThumbnail, %s", err)
			} else {
				product.Image = thumbnail
			}
		}
	}

	err = c.productRepo.Import(products)
	if err != nil {
		return fmt.Errorf("sync: err, %s", err)
	}
	return nil
}

func (p ProductControllerImpl) StopSync() {
	if p.scheduler != nil {
		p.scheduler.Clear()
	} else {
		println("scheduler nil")
	}
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
