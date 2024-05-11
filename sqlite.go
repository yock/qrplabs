package main

import (
  "time"
  "log"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func openDB() (*sql.DB, error) {
  log.Println("Connecting to sqlite")
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
  log.Println("Creating table if needed")
  _, err = db.Exec(sql)

  return db, err
}

func insert(record QRPLabsOrder) {
  log.Println("Insertting record")
  db, err := openDB()
  if err != nil {
    log.Fatal(err)
  }
  defer db.Close()

  var last_updated_at time.Time
  sql := "select last_updated_at from updates order by checked_at desc limit 1"
  db.QueryRow(sql).Scan(&last_updated_at)
  
  if last_updated_at.Compare(record.last_updated_at) < 1 {
    log.Println("No update since last checked")
    return
  }

  log.Println("Beginning transaction")
  tx, err := db.Begin()
  if err != nil {
    log.Fatal(err)
  }
  
  sql = "INSERT INTO updates (list_index, order_id, last_updated_at, checked_at) values(?, ?, ?, ?)"
  log.Println("Preparing statement")
  stmt, err := tx.Prepare(sql)
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()

  log.Println("executing insert")
  _, err = stmt.Exec(record.list_index, record.order_id, record.last_updated_at, record.checked_at)
  if err != nil {
    log.Fatal(err)
  }

  log.Println("Committing transaction")
  err = tx.Commit()
  if err != nil {
    log.Fatal(err)
  }
  log.Println("done")
}
