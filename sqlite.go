package main

import (
  "time"
  "log"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB, error) {
  db, err := sql.Open("sqlite3", "./updates.db")
  if err != nil {
    log.Fatal(err)
  }

  sql := `
  CREATE TABLE IF NOT EXISTS updates (
    id integer not null primary key,
    list_index integer not null,
    order_id not null,
    last_updated_at datetime not null,
    checked_at datetime not null
  );
  `
  _, err = db.Exec(sql)

  return db, err
}

func insert(record *QRPLabsOrder) {
  db, err := openDB()
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  var last_updated_at time.Time
  sql := "select last_updated_at from updates order by checked_at desc limit 1"
  db.QueryRow(sql).Scan(&last_updated_at)

  log.Printf("Most recent tracked update at %s\n", last_updated_at)
  
  if last_updated_at.Compare(record.last_updated_at) >= 1 {
    log.Println("No update since last checked")
    return
  }

  tx, err := db.Begin()
  if err != nil {
    log.Fatal(err)
  }
  
  sql = "INSERT INTO updates (list_index, order_id, last_updated_at, checked_at) values(?, ?, ?, ?)"
  stmt, err := tx.Prepare(sql)
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()

  _, err = stmt.Exec(record.list_index, record.order_id, record.last_updated_at, record.checked_at)
  if err != nil {
    log.Fatal(err)
  }

  err = tx.Commit()
  if err != nil {
    log.Fatal(err)
  }
}
