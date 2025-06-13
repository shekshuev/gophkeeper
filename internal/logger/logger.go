package logger

import (
	"log"
	"sync"

	"go.uber.org/zap"
)

// Logger — обёртка над zap.Logger с синглтон-инициализацией.
// Используется для централизованного логирования в приложении.
type Logger struct {
	Log *zap.Logger // Встроенный zap-логгер
}

var (
	instance *Logger   // глобальный синглтон
	once     sync.Once // обеспечивает однократную инициализацию
)

// NewLogger возвращает синглтон-инстанс логгера.
// По умолчанию инициализирует zap в режиме Production с уровнем "info".
func NewLogger() *Logger {
	once.Do(func() {
		instance = &Logger{Log: zap.NewNop()} // временный "пустой" логгер, пока инициализация не завершится
		if err := instance.initialize("info"); err != nil {
			log.Fatalf("Error initializing zap logger: %v", err)
		}
	})
	return instance
}

// initialize конфигурирует zap.Logger с заданным уровнем логирования.
// Например: "debug", "info", "warn", "error".
func (l *Logger) initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl
	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	l.Log = zl
	return nil
}
