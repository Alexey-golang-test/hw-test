package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in nn goroutines and stops its work when receiving mm errors from tasks.
func Run(tasks []Task, nn, mm int) error {
	if nn < 1 {
		// Защита входных данных. Если кол-во горутин менее 1, то будет работать в 1 поток
		nn = 1
	}

	if mm <= 0 {
		// Если кол-во ошибок меньше либо равно 0, то ошибок быть не должно
		mm = 0
	}

	ch := make(chan Task, nn)
	var wg sync.WaitGroup
	var errAtomic atomic.Int64

	// Запуск горутин на выполнение: слушают канал, выполняют таск, фиксируют ошибки (если есть)
	for ii := 0; ii < nn; ii++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				// Читаем канал
				tt, ok := <-ch

				// Канал закрыт -> выход
				if !ok {
					return
				}

				// Лучше сразу проверить кол-во ошибок, чтобы не делать лишнюю таску
				if ee := errAtomic.Load(); ee >= int64(mm) && ee > 0 {
					return
				}

				// Выполнение таски, нужно увеличивать счетчик ошибок
				if err := tt(); err != nil {
					errAtomic.Add(1)
				}
			}
		}()
	}

	// Отправка тасок в канал
	for _, tt := range tasks {
		ch <- tt
	}
	close(ch)

	wg.Wait()

	if ee := errAtomic.Load(); ee >= int64(mm) && ee > 0 {
		return ErrErrorsLimitExceeded
	}

	return nil
}
