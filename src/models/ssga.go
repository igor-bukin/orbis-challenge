package models

type SSGAResponse struct {
	Data struct {
		FundType []struct {
			Key  string `json:"key"`
			Name string `json:"name"`
			Size int    `json:"size"`
		} `json:"fundType"`
		Funds struct {
			Etfs struct {
				ViewBy struct {
					Overview struct {
						Name string `json:"name"`
					} `json:"overview"`
					Performance struct {
						Name     string `json:"name"`
						Children struct {
							MonthEnd   string `json:"monthEnd"`
							QuarterEnd string `json:"quarterEnd"`
						} `json:"children"`
					} `json:"performance"`
					Pricing struct {
						Name string `json:"name"`
					} `json:"pricing"`
					Documents struct {
						Name string `json:"name"`
					} `json:"documents"`
				} `json:"viewBy"`
				Labels []struct {
					Key       string   `json:"key"`
					Name      string   `json:"name"`
					Tab       []string `json:"tab"`
					Type      string   `json:"type,omitempty"`
					SubLabels []struct {
						Key  string `json:"key"`
						Name string `json:"name"`
					} `json:"subLabels,omitempty"`
					NoSort bool `json:"noSort,omitempty"`
				} `json:"labels"`
				Datas []struct {
					Domicile        string        `json:"domicile"`
					FundName        string        `json:"fundName"`
					FundTicker      string        `json:"fundTicker"`
					FundURI         string        `json:"fundUri"`
					Ter             []interface{} `json:"ter"`
					Nav             []interface{} `json:"nav"`
					Aum             []interface{} `json:"aum"`
					AsOfDate        []string      `json:"asOfDate"`
					PerfAsOf        []string      `json:"PerfAsOf"`
					Mo1             []interface{} `json:"mo1"`
					Qtd             []interface{} `json:"qtd"`
					Ytd             []interface{} `json:"ytd"`
					Yr1             []interface{} `json:"yr1"`
					Yr3             []interface{} `json:"yr3"`
					Yr5             []interface{} `json:"yr5"`
					Yr10            []interface{} `json:"yr10"`
					SinceInception  []interface{} `json:"sinceInception"`
					InceptionDate   []string      `json:"inceptionDate"`
					PerfAsOf1       []string      `json:"PerfAsOf_1"`
					Mo11            []interface{} `json:"mo1_1"`
					Qtd1            []interface{} `json:"qtd_1"`
					Ytd1            []interface{} `json:"ytd_1"`
					Yr11            []interface{} `json:"yr1_1"`
					Yr31            []interface{} `json:"yr3_1"`
					Yr51            []interface{} `json:"yr5_1"`
					Yr101           []interface{} `json:"yr10_1"`
					SinceInception1 []interface{} `json:"sinceInception_1"`
					PrimaryExchange string        `json:"primaryExchange"`
					ClosePrice      []interface{} `json:"closePrice"`
					BidAsk          []interface{} `json:"bidAsk"`
					PremiumDiscount []interface{} `json:"premiumDiscount"`
					DocumentPdf     []struct {
						DocType string `json:"docType"`
						Docs    []struct {
							Language    string `json:"language"`
							Name        string `json:"name"`
							Path        string `json:"path"`
							CanDownload bool   `json:"canDownload"`
						} `json:"docs"`
					} `json:"documentPdf"`
					FundFilter string `json:"fundFilter"`
					PopUp      bool   `json:"popUp"`
					Keywords   string `json:"keywords"`
					PerfIndex  []struct {
						FundName        string `json:"fundName"`
						FundTicker      string `json:"fundTicker"`
						Ter             string `json:"ter"`
						PerfAsOf        string `json:"PerfAsOf"`
						Mo1             string `json:"mo1"`
						Qtd             string `json:"qtd"`
						Ytd             string `json:"ytd"`
						Yr1             string `json:"yr1"`
						Yr3             string `json:"yr3"`
						Yr5             string `json:"yr5"`
						Yr10            string `json:"yr10"`
						SinceInception  string `json:"sinceInception"`
						InceptionDate   string `json:"inceptionDate"`
						PerfAsOf1       string `json:"PerfAsOf_1"`
						Mo11            string `json:"mo1_1"`
						Qtd1            string `json:"qtd_1"`
						Ytd1            string `json:"ytd_1"`
						Yr11            string `json:"yr1_1"`
						Yr31            string `json:"yr3_1"`
						Yr51            string `json:"yr5_1"`
						Yr101           string `json:"yr10_1"`
						SinceInception1 string `json:"sinceInception_1"`
						Num             int    `json:"num"`
					} `json:"perfIndex,omitempty"`
				} `json:"datas"`
				Categories []struct {
					Key           string `json:"key"`
					Name          string `json:"name"`
					SubCategories []struct {
						Key           string `json:"key"`
						Name          string `json:"name"`
						SubCategories []struct {
							Key           string `json:"key"`
							Name          string `json:"name"`
							CheckBox      bool   `json:"checkBox"`
							SubCategories []struct {
								Key   string `json:"key"`
								Name  string `json:"name"`
								Funds string `json:"funds"`
								Size  int    `json:"size"`
							} `json:"subCategories"`
							Funds string `json:"funds"`
							Size  int    `json:"size"`
						} `json:"subCategories,omitempty"`
						Funds string `json:"funds"`
						Size  int    `json:"size"`
					} `json:"subCategories"`
					Funds string `json:"funds"`
					Size  int    `json:"size"`
				} `json:"categories"`
				QuickLinks []struct {
					Name       string `json:"name"`
					Path       string `json:"path"`
					IsExternal bool   `json:"isExternal"`
					Target     bool   `json:"target"`
				} `json:"quickLinks"`
				ExpenseRatio string `json:"expenseRatio"`
				FormatInfo   struct {
					DecimalPoint string `json:"decimalPoint"`
					ThousandsSep string `json:"thousandsSep"`
				} `json:"formatInfo"`
			} `json:"etfs"`
		} `json:"funds"`
	} `json:"data"`
	Msg    string `json:"msg"`
	Status int    `json:"status"`
}
