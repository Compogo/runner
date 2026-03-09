# Compogo Runner 🏃

**Compogo Runner** — это компонент для управления конкурентными задачами (воркерами, серверами, фоновыми процессами) с автоматическим graceful shutdown. Интегрируется с Compogo одной строкой и берёт на себя всю работу по запуску, мониторингу и остановке задач.

## 🚀 Установка

```bash
go get github.com/Compogo/runner
```

## 📦 Быстрый старт

```go
// Пример компонента, использующего раннер
var myWorkerComponent = &component.Component{
    Dependencies: component.Components{runner.Component},
    Run: component.StepFunc(func(c container.Container) error {
        return c.Invoke(func(r runner.Runner) {
            // Создаём задачу из обычной функции
            task := runner.NewTask("worker", runner.ProcessFunc(func(ctx context.Context) error {
                ticker := time.NewTicker(1 * time.Second)
                defer ticker.Stop()
                
                for {
                    select {
                    case <-ctx.Done():
                        return nil  // graceful shutdown
                    case <-ticker.C:
                        doWork()
                    }
                }
            }))
            
            r.RunTask(task)  // ← раннер сам управляет жизнью задачи
        })
    }),
}
```

## ✨ Возможности

### 🎯 Простой интерфейс задачи

Любая задача должна реализовывать `Process` — всего один метод:

```go
type Process interface {
    Process(ctx context.Context) error
}
```

А для функций есть адаптер `ProcessFunc`:

```go
task := runner.NewTask("worker", runner.ProcessFunc(func(ctx context.Context) error {
    // ваша логика
    return nil
}))
```

### 🛡️ Безопасность из коробки

* Panic-защита — если задача паникует, раннер ловит панику, логирует её и продолжает работу
* Детект дубликатов — нельзя запустить одну и ту же задачу дважды, без предварительной остановки.
* Потокобезопасность — все операции защищены мьютексами

### 🔌 Интеграция с жизненным циклом Compogo

* Контекст задачи привязан к closer — при получении сигнала все задачи получают ctx.Done()
* В Stop-фазе раннер автоматически останавливает все задачи и ждёт их завершения

### 📋 Управление задачами

```go
r.RunTask(task)              // запустить одну задачу
r.RunTasks(task1, task2)     // запустить несколько
r.StopTask(task)             // остановить по ссылке
r.StopTaskByName("worker")    // остановить по имени
r.Close()                     // остановить всё и подождать
```

## 🎨 Примеры использования

### HTTP-сервер как задача

```go
task := runner.NewTask("http-server", runner.ProcessFunc(func(ctx context.Context) error {
    server := &http.Server{Addr: ":8080", Handler: router}
    
    go func() {
        <-ctx.Done()
        server.Shutdown(context.Background())
    }()
    
    return server.ListenAndServe()
}))
```

### Пул воркеров

```go
for i := 0; i < 10; i++ {
    workerID := i
    task := runner.NewTask(fmt.Sprintf("worker-%d", i), runner.ProcessFunc(func(ctx context.Context) error {
        for {
            select {
            case <-ctx.Done():
                return nil
            case job := <-jobQueue:
                processJob(job)
            }
        }
    }))
    r.RunTask(task)
}
```

### Периодическая задача

```go
task := runner.NewTask("ticker", runner.ProcessFunc(func(ctx context.Context) error {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return nil
        case <-ticker.C:
            cleanupOldData()
        }
    }
}))
```
