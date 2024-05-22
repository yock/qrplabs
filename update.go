package main

import (
  "time"
  "strings"
  "errors"
  "log"
  "github.com/gocolly/colly/v2"
)

func latest(order_id string) (*QRPLabsOrder, error) {
  order := QRPLabsOrder {
    order_id: order_id,
    checked_at: time.Now(),
  }

  c := colly.NewCollector(
    colly.AllowedDomains("www.qrp-labs.com"),
  )

  var last_updated_at_str string
  c.OnHTML("time[itemprop='dateModified']", func(e *colly.HTMLElement) {
    last_updated_at_str = e.Attr("datetime")
    log.Printf("Last updated at %s\n", last_updated_at_str)
  })

  var list_index string
  c.OnHTML("tr", func(e *colly.HTMLElement) {
    cells := strings.Split(e.Text, "\n")

    if cells[2] == order.order_id {
      list_index = cells[1]
    }
  })

  c.Visit("https://www.qrp-labs.com/qcxmini/assembled.html")

  if list_index == "" {
    return nil, errors.New("Order ID not present in build queue")
  }
  order.list_index = list_index

  if last_updated_at_str == ""{
    return nil, errors.New("Could not determine update date and time from HTML")
  }

  last_updated_at, err := time.Parse(time.RFC3339, last_updated_at_str)
  if err != nil {
    return nil, err
  } else {
    order.last_updated_at = last_updated_at
  }
  if order.list_index == "" {
    return nil, errors.New("not found in build queue")
  }

  return &order, nil 
}
