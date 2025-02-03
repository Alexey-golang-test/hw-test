package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func chanFunc(in In, doneAll In) Out {
	newIn := make(Bi, 10)

	go func() {
		for {
			select {
			case <-doneAll:
				close(newIn)
				go func() {
					// Пропускаем сообщения
					for range in {
						_ = newIn
					}
				}()
				go func() {
					// Пропускаем сообщения
					for range newIn {
						_ = newIn
					}
				}()
				return
			case value, exist := <-in:
				if !exist {
					close(newIn)
					return
				}
				newIn <- value
			}
		}
	}()
	return newIn
}

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Горутина (+ канал) для защиты от подачи в done канал данных, а не его закрытия
	// Например сообщение "stop"
	// Если бы в описании задачи/контракте был бы жестко прописан способ работы done-канала,
	// то этот блок можно было удалить и работать через входящий канал done
	doneAll := make(Bi)
	go func() {
		defer close(doneAll)
		<-done
	}()

	// Последовательный запуск стадий пайплайна.
	for _, ss := range stages {
		// Встраиваемся и перекладываем сообщения в in-канале, чтобы можно было встроить логику работы done-канала
		// in переначитывается, это не ошибка.
		in = ss(chanFunc(in, doneAll))
	}

	resChan := make(Bi, 10)

	// Наличие этого блока кода - скорее защита, т.к. тесты требуют, чтобы в случае работы done-канала
	// результатов не было вообще. Чаще всего, если объект прошел по всем стадиям, то этот результат имеет ценность.
	for {
		select {
		case <-doneAll:
			emptyChan := make(Bi)
			close(emptyChan)
			return emptyChan
		case value, exist := <-in:
			if !exist {
				close(resChan)
				return resChan
			}
			resChan <- value
		}
	}
}
