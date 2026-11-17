package helper

import (
	"context"
	"time"

	"github.com/rs/zerolog"
)

// Sleeper интерефейс временных задержек между повторами
type Sleeper interface {
	Sleep(time.Duration) <-chan time.Time
}

// StandardSleeper структура стандартного подхода задержек между повторами
type StandardSleeper struct{}

func (s *StandardSleeper) Sleep(duration time.Duration) <-chan time.Time {
	return time.After(duration)
}

// RetryableHelper структура хелпера повторов
type RetryableHelper struct {
	sleeper             Sleeper
	triesPeriods        []time.Duration
	applicableErrorFunc func(error) bool
	logger              *zerolog.Logger
}

// Run запускает функцию retryableFunc, которую перезапустит через RetryableHelper.triesPeriods в зависимости
// от попытки, в случае успеха вернет nil, иначе ошибку
func (helper *RetryableHelper) Run(ctx context.Context, retryableFunc func() error) error {
	tries := append([]time.Duration{0}, helper.triesPeriods...)
	try := 0

	var err error

	for _, sleepDur := range tries {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-helper.sleeper.Sleep(sleepDur):
		}

		if try > 0 {
			helper.logger.Info().Err(err).Int("try", try).Msg("retry")
		}
		err = retryableFunc()
		if err == nil {
			return nil
		}

		if helper.applicableErrorFunc(err) {
			try++
		} else {
			return err
		}
	}

	return err
}

// NewRetryableHelper создает и возвращает ссылку на RetryableHelper, беря на вход sleeper, функцию которая должна
// вернуть true, если повторяемая функция вернула подходящую ошибку, иначе эта ошибка сразу будет возвращена и повторов
// не будет. Иначе согласно задержкам в triesPeriods будут выполнены повторы функции.
func NewRetryableHelper(
	sleeper Sleeper,
	applicableError func(err error) bool,
	triesPeriods []time.Duration,
	logger *zerolog.Logger,
) *RetryableHelper {
	if triesPeriods == nil {
		triesPeriods = []time.Duration{
			time.Second,
			time.Second * 3,
			time.Second * 5,
		}
	}

	return &RetryableHelper{
		sleeper:             sleeper,
		triesPeriods:        triesPeriods,
		applicableErrorFunc: applicableError,
		logger:              logger,
	}
}
