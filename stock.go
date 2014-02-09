/*
  Package stock allows to use yahoo finance api to get historical quotes data
  It uses YQL pointing to yahoo.finance.historicaldata to get the necessary info
*/
package yahoofinance

import (
  "fmt"
  "time"
  "net/url"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

const (
  YQL_HISTORICAL_DATA = "select Symbol, Date, Close from yahoo.finance.historicaldata where symbol=\"%s\" and startDate=\"%s\" and endDate=\"%s\""
  YQL_WEB_SERVICE_URL = "http://query.yahooapis.com/v1/public/yql?format=json&diagnostics=false&callback=&%s"
)

// Aux struct to unmarshal API Response
type rootResult struct {
  Query queryInfo
}

// Aux struct to unmarshal API Response
type queryInfo struct {
  Results resultsInfo
}

// Aux struct to unmarshal API Response
type resultsInfo struct {
  Quote []QuoteInfo
}

// Aux struct to unmarshal API Response
type QuoteInfo struct {
  Symbol string
  Close string 
  Date string 
}

func perror(err error) {
  if err != nil {
    panic(err)
  }
}

func HistoricalPrices(symbol string, start, end time.Time) []QuoteInfo{
  yql := fmt.Sprintf(YQL_HISTORICAL_DATA, symbol, start.Format("2006-01-02"), end.Format("2006-01-02")) // Replace dynamic info for q Param
  v := url.Values{} //I use v to hack reserved symbols to be encoded before API call
  v.Set("q", yql)
  v.Add("env", "store://datatables.org/alltableswithkeys") 
  yqlUrl := fmt.Sprintf(YQL_WEB_SERVICE_URL, v.Encode())

  resp, err := http.Get(yqlUrl) //Call to Yahoo API
  perror(err)
  defer resp.Body.Close() //Automatic close resp after use it
  body, err := ioutil.ReadAll(resp.Body) //Get reponse
  perror(err)

  var data rootResult
  err = json.Unmarshal(body, &data) //Fill aux struct with API response
  perror(err)

  return data.Query.Results.Quote //Return quotes
}
