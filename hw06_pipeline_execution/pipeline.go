package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	// Эта функция слушает два канала. Служит как переходник для Stage
	// С одной стороны надо переложить In в Out
	// С другой стороны надо прерывать работу через done-канал
	stageFunc := func(done In, in In) Out {
		resChan := make(Bi)

		go func() {
			defer close(resChan)

			for ii := range in {
				select {
				case <-done:
					return
				case resChan <- ii:
				}
			}
		}()

		return resChan
	}

	// Последовательный запуск стадий пайплайна.
	for _, ss := range stages {
		// in переопределяется, это не ошибка.
		// На каждой итерации цикла исходящий канал функции сохраняется в переменную,
		// и уже на следующей итерации является входящим каналом в функцию
		// Это аналог рекурсивного вызова ss(stageFunc(done, ss(stageFunc(done, ss(stageFunc(done, ... и т.д.))))))
		in = ss(stageFunc(done, in))
	}

	return in
}
