package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

const (
	getLastInsertedRowTime = `SELECT created_at FROM default.info ORDER BY created_at DESC LIMIT 1`
	getCountOfNewRows      = `SELECT count(*) FROM default.info_queue`
	getNewRows             = `SELECT name, age, created_at FROM default.info_queue`
)

type DB struct {
	Conn driver.Conn
}

func NewDB(conn *driver.Conn) *DB {
	return &DB{*conn}
}

type ClickhouseConfig struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
}

func NewClickhouseDB(cfg ClickhouseConfig) (*driver.Conn, error) {
	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)},
		Auth: clickhouse.Auth{
			Database: cfg.DBName,
			Username: cfg.Username,
			Password: cfg.Password,
		},
	})

	if err != nil {
		return nil, err
	}

	if err := conn.Ping(context.Background()); err != nil {
		return nil, err
	}
	return &conn, nil
}

type RowsStruct struct {
	Name      string
	Age       int32
	CreatedAt time.Time
}

func (d *DB) BatchInsert(ctx context.Context, batchSize uint64) error {
	// row := d.Conn.QueryRow(ctx, getLastInsertedRowTime)
	// var lastEntryTime time.Time
	// err := row.Scan(&lastEntryTime)
	// if err != sql.ErrNoRows {
	// 	if err != nil || lastEntryTime.IsZero() {
	// 		return err
	// 	}
	// }

	var countOfNewRows uint64
	// countRow := d.Conn.QueryRow(ctx, getCountOfNewRows)
	// err := countRow.Scan(&countOfNewRows)
	// if err != nil {
	// 	return err
	// }
	countOfNewRows = 25
	// fmt.Println("count : ", countOfNewRows)
	if countOfNewRows >= batchSize {
		newRows, err := d.Conn.Query(ctx, getNewRows)
		if err != nil {
			return err
		}
		fmt.Println("Got new rows")
		infos := make([]RowsStruct, 0)
		for newRows.Next() {
			var info RowsStruct
			err := newRows.Scan(&info.Name, &info.Age, &info.CreatedAt)
			if err != nil {
				return err
			}
			infos = append(infos, info)
		}

		batch, err := d.Conn.PrepareBatch(ctx, "INSERT INTO default.info")
		if err != nil {
			return err
		}
		fmt.Println("Prepared batch")
		for _, info := range infos {
			err := batch.Append(info.Name, info.Age, info.CreatedAt)
			if err != nil {
				return err
			}
		}
		fmt.Println("Send batch")
		return batch.Send()
	}

	return nil
}
