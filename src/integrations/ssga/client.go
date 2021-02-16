package ssga

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/models"
	httpClient "github.com/orbis-challenge/src/transport/http"
	"github.com/sirupsen/logrus"
)

const (
	topHoldingSelectors = ".fund-top-holdings tbody tr td"
	sectorWeight        = ".fund-sector-breakdown"

	percent = "%"
)

type Parser interface {
	GetEtfs() (models.SSGAResponse, error)
	GetTopHoldingsAndWeight(url string) ([]models.Holding, []models.SectorWeight, error)
}

type parserImpl struct {
	client *http.Client
	config *config.SSGA
}

// New adapter constructor.
func New(cfg *config.SSGA) Parser {
	return &parserImpl{
		client: httpClient.NewClient(),
		config: cfg,
	}
}

func (p parserImpl) processResponse(response *http.Response, body interface{}) error {
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	defer httpClient.CloseResponseBody(response)

	return json.Unmarshal(data, body)
}

func (p parserImpl) GetEtfs() (etfs models.SSGAResponse, err error) {
	resp, err := p.client.Get(p.config.URL) //nolint
	if err != nil {
		return etfs, err
	}

	err = p.processResponse(resp, &etfs)
	if err != nil {
		return etfs, err
	}

	return etfs, nil
}

func (p parserImpl) GetTopHoldingsAndWeight(url string) ([]models.Holding, []models.SectorWeight, error) {
	// Request the HTML page.
	resp, err := p.client.Get(url) //nolint
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("status code error: %d %s", resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	holdings := make([]models.Holding, 0)
	holding := models.Holding{}
	isFirst := true
	// Find the top holding
	doc.Find(topHoldingSelectors).Each(func(i int, s *goquery.Selection) {
		if len(s.Text()) > 0 {
			v := p.normalizeStr(s.Text())
			if isFirst {
				holding.Name = v
				isFirst = false
			}
			if strings.Contains(v, percent) {
				w, err := p.convertToFloat(v[:len(v)-1])
				if err != nil {
					logrus.Error("couldn't convert str to float", err)
					return
				}
				holding.Weight = w
				isFirst = true
			}
		}

		if isFirst {
			holdings = append(holdings, holding)
			holding = models.Holding{}
		}
	})

	weightsOfSector := make([]models.SectorWeight, 0)
	isFirst = true
	weight := models.SectorWeight{}
	// Find the sector weight
	doc.Find(sectorWeight).Each(func(i int, selection *goquery.Selection) {
		normalizeSW := p.normalizeSectorWeightStr(selection.Text())
		for j := 3; j < len(normalizeSW); j++ {
			sw := normalizeSW[j]
			if (j+1)%2 == 0 && j != 3 {
				weightsOfSector = append(weightsOfSector, weight)
				weight = models.SectorWeight{}
			}

			if isFirst {
				weight.Name = sw
				isFirst = false
			}
			if strings.Contains(sw, percent) {
				w, err := p.convertToFloat(sw[:len(sw)-1])
				if err != nil {
					logrus.Error("couldn't convert str to float", err)
					return
				}
				weight.Weight = w
				isFirst = true
			}
		}
	})

	return holdings, weightsOfSector, nil
}

func (p parserImpl) convertToFloat(str string) (f float64, err error) {
	return strconv.ParseFloat(str, 64)
}

func (p parserImpl) normalizeStr(str string) string {
	return strings.TrimSpace(strings.ReplaceAll(str, "\n", ""))
}

func (p parserImpl) normalizeSectorWeightStr(str string) []string {
	sw := strings.Split(str, "\n")
	normalizeSectorWeight := make([]string, 0)
	for j := range sw {
		normalizedWeight := p.normalizeStr(sw[j])
		if normalizedWeight == "" {
			continue
		}

		normalizeSectorWeight = append(normalizeSectorWeight, normalizedWeight)
	}

	return normalizeSectorWeight
}
