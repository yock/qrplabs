package main

import (
  "time"
  "log"
  "os"
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
  order, err := latest(order_id)

  if err != nil {
    log.Fatalf("Could not get latest for %s: %s", order_id, err)
  }
  insert(order)
}
