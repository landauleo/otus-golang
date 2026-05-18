package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	proxyIn := make(chan interface{})

	// 1. Менеджер ВХОДА
	go func() {
		defer close(proxyIn)
		for {
			select {
			case <-done:
				return
			case val, ok := <-in:
				if !ok {
					return
				}
				// Передаем дальше, но если пришел done — выходим
				select {
				case <-done:
					return
				case proxyIn <- val:
				}
			}
		}
	}()

	// 2. Строим цепочку стейджей
	var currentIn In = proxyIn
	for _, stage := range stages {
		currentIn = stage(currentIn)
	}

	// 3. Менеджер ВЫХОДА
	proxyOut := make(chan interface{})
	go func() {
		defer close(proxyOut)
		for {
			select {
			case <-done:
				// ВАЖНО: Если пришел done, мы не можем просто сделать return!
				// Нам нужно "очистить" канал currentIn, иначе стейджи зависнут.
				// Запускаем горутину-пылесос (drain)
				go func() {
					for range currentIn {
						// Просто читаем всё подряд и выбрасываем в никуда,
						// пока последний стейдж не закроет свой канал.
					}
				}()
				return // Теперь спокойно выходим, proxyOut закроется, а стейджи не зависнут!

			case val, ok := <-currentIn:
				if !ok {
					return
				}

				select {
				case <-done:
					// Если done пришел прямо сейчас, запускаем пылесос для остатков
					go func() {
						for range currentIn {
						}
					}()
					return
				case proxyOut <- val:
				}
			}
		}
	}()

	return proxyOut
}
