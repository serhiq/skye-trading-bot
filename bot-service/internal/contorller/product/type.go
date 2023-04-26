package product

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"github.com/go-resty/resty/v2"
	"github.com/nfnt/resize"
	"github.com/pkg/errors"
	"github.com/serhiq/skye-trading-bot/internal/config"
	"github.com/serhiq/skye-trading-bot/internal/contorller"
	"github.com/serhiq/skye-trading-bot/internal/delivery"
	"github.com/serhiq/skye-trading-bot/internal/logger"
	"github.com/serhiq/skye-trading-bot/internal/utils"
	domainProduct "github.com/serhiq/skye-trading-bot/pkg/type/product"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ProductControlerImpl struct {
	productRepo contorller.ProductRepository
	//orderHistoryRepo OrderHistoryRepository
	productProvider delivery.ProductProvider
	scheduler       *gocron.Scheduler
	client          *resty.Client
}

func New(productRepo contorller.ProductRepository, productProvider delivery.ProductProvider, client *resty.Client) contorller.ProductController {
	return &ProductControlerImpl{
		productRepo:     productRepo,
		productProvider: productProvider,
		client:          client,
	}
}

func (c ProductControlerImpl) GetProductByParent(parentUuid string) ([]*domainProduct.Product, error) {
	return c.productRepo.GetProductByParent(parentUuid)

}

func (c ProductControlerImpl) GetProductByUuid(uuid string) (*domainProduct.Product, error) {
	return c.productRepo.GetProductByUuid(uuid)
}

func (c ProductControlerImpl) StartSync() {
	scheduler := gocron.NewScheduler(time.UTC)
	_, err := scheduler.Every(int(10)).Minutes().Do(updateMenu, c)
	if err != nil {
		logger.SugaredLogger.Errorf("Error: updateMenu: %s", err.Error())
	}

	c.scheduler = scheduler
	scheduler.StartAsync()
}

var updateMenu = func(s ProductControlerImpl) {
	err := s.start()
	if err != nil {
		logger.SugaredLogger.Errorf("sync: err, %s", err.Error())
	}
}

func (c ProductControlerImpl) start() error {

	products, err := c.productProvider.GetProducts()

	if err != nil {
		return err
	}

	err = clearDirectory(config.PreviewCachePatch)
	if err != nil {
		fmt.Printf("sync: clear dir error, %s", err.Error())
	}

	err = clearDirectory(config.TempPatch)
	if err != nil {
		fmt.Printf("sync: clear dir error, %s", err.Error())
	}

	for _, product := range products {

		if product.Image != "" {
			thumbnail, err := createThumbnail(c.client, product.UUID, product.Image)
			if err != nil {
				logger.SugaredLogger.Errorf("sync: createThumbnail, %s", err.Error())
			} else {
				product.Image = thumbnail
			}
		}
	}

	err = c.productRepo.Import(products)
	if err != nil {
		return fmt.Errorf("sync: err, %s", err.Error())
	}
	return nil
}

func clearDirectory(folderPath string) error {
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			err = os.Remove(folderPath + "/" + file.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (p ProductControlerImpl) StopSync() {
	if p.scheduler != nil {
		p.scheduler.Clear()
	} else {
		println("scheduler nil")
	}
}

func createThumbnail(r *resty.Client, fileName string, image string) (string, error) {

	previewImage := filepath.Join(config.PreviewCachePatch, fileName)
	var imagePath string

	if strings.HasPrefix(image, "http") {
		tmpImage := filepath.Join(config.TempPatch, fileName)
		defer os.Remove(tmpImage)

		resp, err := r.R().
			SetOutput(tmpImage).
			Get(image)

		if err != nil {
			return "", fmt.Errorf("external_api: error get image, url=%s, err=%s", image, err.Error())
		}

		if !resp.IsSuccess() {
			return "", fmt.Errorf("sync: error get image, , url=%s, code == %d", image, resp.StatusCode())
		}
		imagePath = tmpImage
	} else {
		imagePath = filepath.Join(config.FileProviderPatch, fileName)
	}

	err := resizeTmpFile(imagePath, previewImage)
	if err != nil {
		return "", fmt.Errorf("sync: error resize image, %s", err.Error())
	}

	return fileName, nil
}

func resizeTmpFile(tmpDest string, filename string) error {
	file, err := os.Open(tmpDest)
	defer file.Close()

	if err != nil {
		if os.IsNotExist(err) {
			return errors.Wrap(err, fmt.Sprintf("failed to read file: %s not found in current directory %s", tmpDest, utils.GetCurrentDirectory()))
		}
		return err
	}

	img, err := jpeg.Decode(file)
	if err != nil {
		return fmt.Errorf("%s while decoding", err)
	}

	m := resize.Resize(0, 300, img, resize.Lanczos3)

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	return jpeg.Encode(out, m, nil)
}
