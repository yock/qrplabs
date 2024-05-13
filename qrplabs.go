package main

import (
  "strings"
  "time"
  "log"
  "os"
  "github.com/gocolly/colly/v2"
)

type QRPLabsOrder struct {
  list_index string
  order_id string
  checked_at time.Time
  last_updated_at time.Time
}

func main() {
  order_id := os.Args[1]
  log.Printf("Getting latest build order for %s", order_id)
  c := colly.NewCollector(
    colly.AllowedDomains("www.qrp-labs.com"),
  )

  order := QRPLabsOrder {
    order_id: order_id,
    checked_at: time.Now(),
  }
  
  c.OnHTML("time", func(e *colly.HTMLElement) {
    last_updated_at, err := time.Parse(time.RFC3339, e.Attr("datetime"))

    if err != nil {
      log.Fatal(err)
    } else {
      order.last_updated_at = last_updated_at
    }
  })

  c.OnHTML("tr", func(e *colly.HTMLElement) {
    cells := strings.Split(e.Text, "\n")

    if cells[2] == order.order_id {
      order.list_index = cells[1]
    }
  })

  c.Visit("https://www.qrp-labs.com/qcxmini/assembled.html")
  if order.list_index == "" {
    log.Fatalf("Order %s not found in build queue", order_id)
  }
  insert(order)
}
