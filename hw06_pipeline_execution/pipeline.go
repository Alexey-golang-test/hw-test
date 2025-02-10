package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func chanDataTransit(in In, out Bi, done In) {
	for {
		select {
		case <-done:
			close(out)
			go func() {
				// Сброс данных из канала
				for range in {
					_ = in
				}
			}()
			go func() {
				// Сброс данных из канала
				for range out {
					_ = out
				}
			}()
			return
		case value, exist := <-in:
			if !exist {
				close(out)
				return
			}
			out <- value
		}
	}
}

func stageRun(in In, done In, stage Stage) Bi {
	outChan := make(Bi)

	go func() {
		stageOut := stage(in)
		chanDataTransit(stageOut, outChan, done)
	}()

	return outChan
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

	// Промежуточный канал
	out := make(Bi)
	go chanDataTransit(in, out, doneAll)

	// Последовательный запуск стадий пайплайна.
	for _, ss := range stages {
		out = stageRun(out, doneAll, ss)
	}

	return out
}
