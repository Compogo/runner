# Compogo Runner 

[![Go Reference](https://pkg.go.dev/badge/github.com/Compogo/runner.svg)](https://pkg.go.dev/github.com/Compogo/runner)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Фреймворк для управления фоновыми процессами в приложениях на [Compogo](https://github.com/Compogo/compogo).

Предоставляет:

* Запуск и остановку фоновых задач (процессов)
* Graceful shutdown через io.Closer
* Middleware для логирования, метрик, восстановления после паник и автоматического перезапуска
* Потокобезопасное управление процессами

## Установка

```shell
go get github.com/Compogo/runner
```

## Быстрый старт

```go
package main

import (
    "context"
    "time"

    "github.com/Compogo/compogo"
    "github.com/Compogo/runner"
    runnerlogger "github.com/Compogo/runner/middleware/logger"
)

func main() {
    app := compogo.NewApp("myapp",
        // Подключаем Runner
        compogo.WithComponents(&runner.Component),
        // Подключаем middleware для логирования
        compogo.WithComponents(&runnerlogger.Component),
        // Другие опции
        compogo.WithContainer(container, containerCmp),
        compogo.WithLogger(logger, loggerCmp),
        compogo.WithOsSignalCloser(),
    )

    if err := app.Serve(); err != nil {
        panic(err)
    }
}
```

## Использование

### Создание процесса

```go
// Создание задачи через Task
task := runner.NewTask("worker", func(ctx context.Context) error {
    // Фоновая работа
    select {
    case <-ctx.Done():
        return nil
    case <-time.After(time.Second):
        // делаем работу
    }
    return nil
})

// Или через структуру, реализующую интерфейс Process
type MyProcess struct {
    name string
}

func (p *MyProcess) Process(ctx context.Context) error {
    // работа
    return nil
}

func (p *MyProcess) Name() string {
    return p.name
}

func (p *MyProcess) Close() error {
    // очистка
    return nil
}
```

### Запуск процессов

```go
// Получение Runner из контейнера
var r runner.Runner
container.Invoke(func(runner runner.Runner) { r = runner })

// Запуск одного процесса
r.RunProcess(task)

// Запуск нескольких процессов
r.RunProcesses(task1, task2, task3)
```

### Остановка процессов

```go
// Остановка по объекту
r.StopProcess(task)

// Остановка по имени
r.StopProcessByName("worker")

// Остановка всех процессов (через Close)
r.Close()
```

### Проверка наличия процессов

```go
if r.HasProcess(task) {
    // процесс уже запущен
}

if r.HasProcessByName("worker") {
    // процесс с именем "worker" запущен
}
```

## Middleware

### Logger — Логирование выполнения

Логирует начало, завершение и ошибки выполнения процессов.

```go
import runnerlogger "github.com/Compogo/runner/middleware/logger"

// Подключение через компонент
app.AddComponents(&runnerlogger.Component)

// Или программно
var r runner.Runner
var m *runnerlogger.Logger
container.Invoke(func(runner runner.Runner, logger *runnerlogger.Logger) {
    r.Use(logger)
})
```

### Metric — Сбор метрик Prometheus

Metric — Сбор метрик Prometheus

```go
import runnermetric "github.com/Compogo/runner/middleware/metric"

// Подключение через компонент
app.AddComponents(&runnermetric.Component)

// Метрика: compogo_runner_task{app="myapp"}
```

### Recover — Восстановление после паник

Перехватывает паники в процессах, логирует стек-трейс и преобразует в ошибку.

```go
import runnerrecover "github.com/Compogo/runner/middleware/recover"

// Подключение через компонент
app.AddComponents(&runnerrecover.Component)
```

### Restore — Автоматический перезапуск

При ошибке выполнения процесса перезапускает его с заданной задержкой.

```go
import runnerrestore "github.com/Compogo/runner/middleware/restore"

// Подключение через компонент
app.AddComponents(&runnerrestore.Component)

// Настройка задержки через флаг
// --runner.middleware.restore.delay=5s
```

## Зависимости

* [Compogo](https://github.com/Compogo/compogo) — основной фреймворк
* [Prometheus](https://github.com/prometheus/client_golang) — метрики (опционально)

## Лицензия

```plantuml
MIT License

Copyright (c) 2026 Compogo

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.

```
