package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)


type Config struct {
	Addr   		string       	`env:"RUN_ADDRESS"`
	PollInt   	time.Duration 	`env:"POLL_WORKER_INTERVAL"`
	DBAddr 		string			`env:"DATABASE_URI"`
	AccrualAddr string			`env:"ACCRUAL_SYSTEM_ADDRESS"`
}

func New() (*Config, error) {
	c := &Config{}

	addrFlag := flag.String("a", "127.0.0.1:8080", "host:port")
	pollIntFlag := flag.Duration("p", time.Minute, "poll interval")
	DBAddressFlag := flag.String("d", "postgres://postgres:postgres@localhost:5432/marketDB", "DataBase address")
	accrualAddrFlag := flag.String("r", "127.0.0.1:8088", "host:port")

	flag.Parse()

	err := env.Parse(c)
	if err != nil {
		fmt.Println("Ошибка парсинга env", err)
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


	fmt.Printf("Config: %+v\n", c)

	return c, nil
}

/*
### Конфигурирование сервиса накопительной системы лояльности

Сервис должен поддерживать конфигурирование следующими методами:

- адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `a`;
- адрес подключения к базе данных: переменная окружения ОС `DATABASE_URI` или флаг `d`;
- адрес системы расчёта начислений: переменная окружения ОС `ACCRUAL_SYSTEM_ADDRESS` или флаг `r`.
*/