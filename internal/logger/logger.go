package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"time"
)

func InitLogger() {
	zerolog.TimeFieldFormat = time.RFC3339
}

func Info(msg string) {
	log.Info().Msg(msg)
}

func Error(msg string) {
	log.Error().Msg(msg)
}
func Fatal(msg string) {
	log.Fatal().Msg(msg)
}
