package config

import (
	"flag"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)


type Config struct {
	Addr   		string       	`env:"RUN_ADDRESS"`
	PollInt   	time.Duration 	`env:"POLL_WORKER_INTERVAL"`
	DBAddr 		string		`env:"DATABASE_URI"`
	AccrualAddr     string		`env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() (*Config, error) {
	c := &Config{}

	addrFlag := flag.String("a", "127.0.0.1:8081", "host:port")
	pollIntFlag := flag.Duration("p", 1 * time.Second, "poll interval")
	DBAddressFlag := flag.String("d", "postgres://postgres:postgres@localhost:5432/marketDB", "DataBase address")
	accrualAddrFlag := flag.String("r", "http://127.0.0.1:8082", "host:port")

	flag.Parse()

	if err := env.Parse(c); err != nil {
		return nil, err
	}

	if _, ok := os.LookupEnv("POLL_WORKER_INTERVAL"); !ok {
		c.Addr = *addrFlag
	}

	if _, ok := os.LookupEnv("RESTORE"); !ok {
		c.PollInt = *pollIntFlag
	}

	if _, ok := os.LookupEnv("DATABASE_URI"); !ok {
		c.DBAddr = *DBAddressFlag
	}

	if _, ok := os.LookupEnv("ACCRUAL_SYSTEM_ADDRESS"); !ok {
		c.AccrualAddr = *accrualAddrFlag
	}

	return c, nil
}

/*
### Конфигурирование сервиса накопительной системы лояльности

- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `a`;
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `d`;
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `r`.
*/
