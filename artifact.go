// Helper codes for Steam Community Market: Artifact

package artifact

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"time"
)

const (
	appID = 583950 // Artifact

	apiURL        = "https://steamcommunity.com/market/search/render/"
	marketBaseURL = "https://steamcommunity.com/market/listings/%d/%s"
	imageBaseURL  = "https://steamcommunity-a.akamaihd.net/economy/image/"

	searchMaxCount = 100
)

// Lang type for language
type Lang string

// Rarity type for rarity
type Rarity string

// SortDirection type for sort direction
type SortDirection string

// SortColumn type for sort column
type SortColumn string

// iconURL type for icon URLs
type iconURL string

// constants
const (
	// languages
	LangEnglish Lang = "english"
	LangKorean  Lang = "koreana"
	LangDefault      = LangEnglish

	// rarities
	RarityCommon   Rarity = "tag_Rarity_Common"
	RarityUncommon Rarity = "tag_Rarity_Uncommon"
	RarityRare     Rarity = "tag_Rarity_Rare"
	RarityAll      Rarity = ""

	// sort directions
	SortDirectionAsc     SortDirection = "asc"
	SortDirectionDesc    SortDirection = "desc"
	SortDirectionDefault               = SortDirectionAsc

	// sort columns
	SortColumnName     SortColumn = "name"
	SortColumnQuantity SortColumn = "quantity"
	SortColumnPrice    SortColumn = "price"
	SortColumnDefault             = SortColumnName
)

// MarketSearchResult struct for search result of community market
type MarketSearchResult struct {
	Success    bool             `json:"success"`
	Start      int              `json:"start"`
	PageSize   int              `json:"pagesize"`
	TotalCount int              `json:"total_count"`
	SearchData MarketSearchData `json:"searchdata"`
	Results    []MarketItem     `json:"results"`
}

// MarketSearchData struct for data from search result of community market
type MarketSearchData struct {
	Query              string `json:"query"`
	SearchDescriptions bool   `json:"search_descriptions"`
	TotalCount         int    `json:"total_count"`
	PageSize           int    `json:"pagesize"`
	Prefix             string `json:"prefix"`
	ClassPrefix        string `json:"class_prefix"`
}

// MarketItem struct for items of community market
type MarketItem struct {
	Name             string                `json:"name"`
	HashName         string                `json:"hash_name"`
	SellListings     int                   `json:"sell_listings"`
	SellPrice        int                   `json:"sell_price"`
	SellPriceText    string                `json:"sell_price_text"`
	AppIcon          string                `json:"app_icon"`
	AppName          string                `json:"app_name"`
	AssetDescription MarketItemDescription `json:"asset_description"`
	SalePriceText    string                `json:"sale_price_text"`
}

// StoreURL generates a URL which directs to the Steam Communit Market page of this item
func (i *MarketItem) StoreURL() string {
	return fmt.Sprintf(marketBaseURL, appID, i.HashName)
}

// ToJSON prettifies MarketItem to JSON string
func (i *MarketItem) ToJSON() string {
	bytes, err := json.MarshalIndent(*i, "", " ")

	if err == nil {
		return string(bytes)
	}

	return err.Error()
}

// MarketItemDescription struct for description of items of community market
type MarketItemDescription struct {
	AppID                       int     `json:"appid"`
	ClassID                     string  `json:"classid"`
	InstanceID                  string  `json:"instanceid"`
	Currency                    int     `json:"currency"`
	BackgroundColor             string  `json:"background_color"`
	Icon                        iconURL `json:"icon_url"`
	IconLarge                   iconURL `json:"icon_url_large"`
	Tradable                    int     `json:"tradable"`
	Name                        string  `json:"name"`
	Type                        string  `json:"type"`
	MarketName                  string  `json:"market_name"`
	MarketHashName              string  `json:"market_hash_name"`
	Commodity                   int     `json:"commodity"`
	MarketTradableRestriction   int     `json:"market_tradable_restriction"`
	MarketMarketableRestriction int     `json:"market_marketable_restriction"`
	Marketable                  int     `json:"marketable"`
}

// IconURL returns the URL of icon
func (d *MarketItemDescription) IconURL() string {
	return imageURL(d.Icon)
}

// LargeIconURL returns the URL of large icon
func (d *MarketItemDescription) LargeIconURL() string {
	return imageURL(d.IconLarge)
}

// FetchAll fetches all items with given rarity and language
func FetchAll(rarity Rarity, language Lang, sortColumn SortColumn, sortDirection SortDirection) (items []MarketItem, err error) {
	offset := 0
	items = []MarketItem{}

	var results MarketSearchResult
	for {
		results, err = searchFor(rarity, language, searchMaxCount, offset, sortColumn, sortDirection)

		if err != nil {
			return nil, fmt.Errorf("Failed to search all with error: %s", err)
		}

		if !results.Success {
			return nil, fmt.Errorf("Failed to search all: %v", results)
		}

		// no more results, then stop the loop
		if len(results.Results) <= 0 {
			break
		}

		items = append(items, results.Results...)

		offset += searchMaxCount
	}

	return items, nil
}

// search items for given rarity, language, count, and offset
func searchFor(rarity Rarity, language Lang, count, offset int, sortColumn SortColumn, sortDirection SortDirection) (MarketSearchResult, error) {
	var err error
	var req *http.Request
	if req, err = http.NewRequest("GET", apiURL, nil); err == nil {
		// set HTTP headers
		//req.Header.Set("Authorization", "xxxx") // set auth header

		// set parameters
		queries := req.URL.Query()
		queries.Add("appid", strconv.Itoa(appID))
		queries.Add("search_descriptions", strconv.Itoa(0))
		queries.Add("norender", strconv.Itoa(1)) // render to JSON
		queries.Add(fmt.Sprintf("category_%d_Rarity[]", appID), string(rarity))
		queries.Add("count", strconv.Itoa(count))
		queries.Add("start", strconv.Itoa(offset))
		queries.Add("sort_column", string(sortColumn))
		queries.Add("sort_dir", string(sortDirection))
		req.URL.RawQuery = queries.Encode()

		// set cookies
		req.AddCookie(&http.Cookie{Name: "Steam_Language", Value: string(language)})

		httpClient := &http.Client{
			Transport: &http.Transport{
				Dial: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 300 * time.Second,
				}).Dial,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ResponseHeaderTimeout: 10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
		}

		var resp *http.Response
		resp, err = httpClient.Do(req)

		if resp != nil {
			defer resp.Body.Close()
		}

		if err == nil {
			if resp.StatusCode == 200 {
				var bytes []byte
				if bytes, err = ioutil.ReadAll(resp.Body); err == nil {
					var results MarketSearchResult
					if err = json.Unmarshal(bytes, &results); err == nil {
						return results, nil
					}
				}
			} else {
				err = fmt.Errorf("HTTP %d (%s)", resp.StatusCode, resp.Status)
			}
		}
	}

	return MarketSearchResult{}, err
}

// get full image URL from given icon URL
func imageURL(iconURL iconURL) string {
	return fmt.Sprintf("%s%s", imageBaseURL, iconURL)
}
