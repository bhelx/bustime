package lib

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type sqliteStorage struct {
	database *sql.DB
}

type Storage interface {
	Store(*Vehicle) error
	Close()
}

func NewStorage() (Storage, error) {
	db, sqlerr := sql.Open("sqlite3", "./vehicle_readings.db")
	if sqlerr != nil {
		return nil, sqlerr
	}
	_, err := os.Stat("./vehicle_readings.db")
	if err != nil {
		if os.IsNotExist(err) {
			_, err = db.Exec(CreateVehiclesSQL)
			if err != nil {
				fmt.Printf("%q: %s\n", err, CreateVehiclesSQL)
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	s := &sqliteStorage{database: db}
	return s, nil
}

func (s *sqliteStorage) Store(v *Vehicle) error {
	// first look to see if this is a new reading
	rows := s.database.QueryRow(fmt.Sprintf("select count(id) from vehicle_readings where tmstmp >= %d and vid = '%s'", v.Tmstmp.Time.Unix(), v.Vid))
	var count int
	if err := rows.Scan(&count); err != nil {
		return err
	}

	if count <= 0 {
		tx, err := s.database.Begin()
		if err != nil {
			return err
		}
		cols, vals := v.ToSql()
		sql := fmt.Sprintf("insert into vehicle_readings(%s) values(%s)", strings.Join(cols, ","), strings.Join(vals, ","))
		//	log.Println(sql)
		_, err = tx.Exec(sql)
		if err != nil {
			return err
		}
		tx.Commit()
	}

	return nil
}

func (s *sqliteStorage) Close() {
	s.database.Close()
}
