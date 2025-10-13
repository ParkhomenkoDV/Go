*Полезно знать: если ваш сервер на Go не запускается и при этом ничего не выводится в консоль, проверьте, обрабатываете ли вы ошибку при запуске. Вместо `http.ListenAndServe("127.0.0.1:8080", mux)` пишите `log.Fatal(http.ListenAndServe("127.0.0.1:8080", mux))`. 
Это позволит увидеть сообщение об ошибке. Например, `listen tcp 127.0.0.1:8080: bind: address already in use` может означать, что порт уже занят другой программой.*

# Интроспекция ошибок

Ни одна программа не застрахована от ошибок. Если при неверных входных данных или недоступности ресурсов можно просто вывести сообщение для пользователя, то причины аварийных ситуаций, приводящих к краху программы, бывает сложно установить. Часто игнорирование простых ошибок приводит к серьёзным последствиям. Исправить их стоит гораздо дороже, чем изначально уделить больше внимания потенциальным проблемам.

Поскольку избежать ошибок невозможно, языки программирования предоставляют инструментарий для их обработки. В некоторых языках есть механизм исключений, при котором перехватывать и обрабатывать ошибки можно в любом месте программы. Язык Go предлагает отслеживать ошибки и реагировать на них непосредственно в момент возникновения. Код получается менее компактным, зато разработчик лучше контролирует выполнение программы.

Так как Go позволяет возвращать из функций несколько значений, принято последним значением возвращать ошибку. Если возвращаемое значение ошибки равно `nil`, функция завершилась корректно. В противном случае нужно обработать ошибку и/или вернуть её выше по стеку. Если функция завершилась с ошибкой, не стоит использовать остальные возвращаемые значения: они могут быть не определены, функция может выполниться не полностью и не успеть вычислить значения.

В этой теме вы:
- познакомитесь c типом `error`;
- научитесь обёртывать и разворачивать ошибки;
- изучите функции для интроспекции ошибок.

## Тип error

Ошибкой в языке Go может быть значение любого типа, который совместим с интерфейсным типом `error`:
```go
type error interface {
    Error() string
}
```

Для создания ошибки чаще всего применяют функцию `fmt.Errorf(format string, a ...interface{})`, где указывают шаблон форматирования и дополнительные параметры. Ещё можно использовать функцию `errors.New(text string)`, принимающую в параметре строку:
```go
func GetUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user's id")
    }
    // FindUser ищет в БД пользователя с указанным id
    // если пользователь не найден, то user равен nil
    // также может возвращаться ошибка, возникшая в процессе поиска
    user, err := FindUser(id)
    if err != nil {
        return nil, err
    }
    if user == nil {
        return nil, fmt.Errorf("can't find user (id: %d)", id)
    }
    return user, nil
}
```

По умолчанию ошибки имеют тип, который может хранить только строку.

Чтобы сохранить и передать параметры окружения и настройки, при которых произошла ошибка, можно создать свой тип ошибки с нужными полями и определить для него метод `Error() string`. Это позволит включать в ошибку всю необходимую дополнительную информацию.

Например, чтобы сохранить время возникновения ошибки, можно определить такой тип:
```go
// TimeError предназначен для ошибок с фиксацией времени возникновения.
type TimeError struct {
    Time time.Time
    Err  error
}

// Error добавляет поддержку интерфейса error для типа TimeError.
func (te *TimeError) Error() string {
    return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

// NewTimeError записывает ошибку err в тип TimeError c текущим временем.
func NewTimeError(err error) error {
    return &TimeError{
        Time: time.Now(),
        Err:  err,
    }
}
```

**Важно**
При создании ошибок лучше возвращать не структуру, а указатель на структуру. Если вместо указателя возвращать структуру, то ошибки из разных пакетов, но с одинаковым текстом будут равны. А если возвращать указатель и создавать ошибки с одинаковыми данными, они не будут равны друг другу, так как будут ссылаться на разные области памяти.

Для примера рассмотрим, как по умолчанию определяется ошибка в пакете `errors`:
```go
type errorString struct {
    s string
}

func New(text string) error {
    return &errorString{text}
}

func (e *errorString) Error() string {
    return e.s
}
```

Можно было бы реализовать функцию New() и метод Error() таким образом:
```go
func New(text string) error {
    return errorString{text}
}

func (e errorString) Error() string {
    return e.s
}
```

Но тогда ошибки с одинаковым текстом, созданные в разных пакетах, были бы равны. Например, в этом случае ошибка `errors.New("EOF")` была бы равна ошибке `io.EOF`, которая определяется точно так же в пакете `io`, что могло бы привести к неверной обработке ошибок.

Посмотрите, как на практике используется тип `TimeError`. Допустим, при старте программы нужно прочитать файл конфигурации, а если не удалось его прочитать, то следует вывести ошибку и завершить работу:
```go
package main

import (
    "fmt"
    "os"
    "time"
)

// TimeError предназначен для ошибок с фиксацией времени возникновения.
type TimeError struct {
    Time time.Time
    Err  error
}

// Error добавляет поддержку интерфейса error для типа TimeError.
func (te *TimeError) Error() string {
    return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

// NewTimeError записывает ошибку err в тип TimeError c текущим временем.
func NewTimeError(err error) error {
    return &TimeError{
        Time: time.Now(),
        Err:  err,
    }
}

func ReadTextFile(filename string) (string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        return ``, NewTimeError(err)
    }
    return string(data), nil
}

func main() {
    data, err := ReadTextFile("myconfig.yaml")
    if err != nil {
        fmt.Println(err)
        os.Exit(0)
    }
    // ...
}
```

При отсутствии файла программа будет выводить время возникновения ошибки и текст:
```
2023/01/24 23:00:00 open myconfig.yaml: no such file or directory
```

Создание своих типов ошибок позволяет собирать дополнительную информацию о месте, времени и причинах возникновения ошибки.

## Упаковка ошибок

Часто ошибка возвращается вверх по стеку с дополнительным комментарием. В момент, когда функция хочет её обработать, уже сложно понять, какой была исходная ошибка.

![](./assets/images/errors.png)

Чтобы решить эту проблему, Go даёт возможность упаковывать ошибки (процесс называется e**rror wrapping**). Для этого достаточно в функции `Errorf()` добавить к ошибке спецификатор `%w`. 

Функция `errors.Unwrap()` снимает один уровень обёртки. Если ошибка была обёрнута несколько раз, то для получения исходной ошибки нужно вызывать `errors.Unwrap()` до тех пор, пока она не начнёт возвращать `nil`.

```go
package main

import (
    "errors"
    "fmt"
    "os"
    "time"
)

func ReadTextFile(filename string) (string, error) {
    data, err := os.ReadFile(filename)
    if err != nil {
        // добавляем время и обёртываем ошибку
        now := time.Now().Format("2006/01/02 15:04:05")
        return "", fmt.Errorf("%s %w", now, err)
    }
    return string(data), nil
}

func main() {
    data, err := ReadTextFile("myconfig.yaml")
    if err != nil {
        fmt.Println(err)
        // можем узнать оригинальную ошибку
        fmt.Println("Original error:", errors.Unwrap(err))
        os.Exit(0)
    }
    fmt.Println(data)
    // ...
}
```

Результат:
```
2025/10/05 10:04:57 open myconfig.yaml: no such file or directory
Original error: open myconfig.yaml: no such file or directory
```

В этом примере исходная ошибка распаковывается функцией `Unwrap()` и выводится в консоль. 

Если тип ошибки не имеет метода `Unwrap()` или упакованная ошибка отсутствует, то `errors.Unwrap()` вернёт `nil`. Поэтому метод `Unwrap()` следует определять для типов ошибок, которые упаковывают исходные ошибки.

Вернёмся к примеру с `TimeError`. Чтобы иметь возможность получать исходную ошибку, добавим для `TimeError` метод `Unwrap()`:
```go
package main

import (
	"errors"
	"fmt"
	"os"
	"time"
)

type TimeError struct {
	Time time.Time
	Err  error
}

func (te *TimeError) Error() string {
	return fmt.Sprintf("%v %v", te.Time.Format("2006/01/02 15:04:05"), te.Err)
}

func NewTimeError(err error) error {
	return &TimeError{
		Time: time.Now(),
		Err:  err,
	}
}

func (te *TimeError) Unwrap() error {
	return te.Err
}

func ReadTextFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return ``, NewTimeError(err)
	}
	return string(data), nil
}

func main() {
	data, err := ReadTextFile("myconfig.yaml")
	if err != nil {
		fmt.Println(err)
		// можем узнать оригинальную ошибку для TimeError
		fmt.Println("Original error:", errors.Unwrap(err))
		os.Exit(0)
	}
	fmt.Println(data)
}
```

Результат:
```
2025/10/05 10:07:23 open myconfig.yaml: no such file or directory
Original error: open myconfig.yaml: no such file or directory
```

___
Представьте, что на стажировке вы работаете с программой, в которой определён тип `LabelError`. Он добавляет слева к ошибке строковую метку `[МЕТКА] ошибка`. Например: `[FILE] текст ошибки`.

Опишите для этого типа методы `Error()` и `NewLabelError(label string, err error)`, чтобы программа выводила `[FILE] open mytest.txt: no such file or directory`.

```go
package main

import (
    "fmt"
    "os"
    "strings"
)

// LabelError описывает ошибку с дополнительной меткой.
type LabelError struct {
    Label string // метка должна быть в верхнем регистре
    Err   error
}

// добавьте методы Error() и NewLabelError(label string, err error)

// ...

func main() {
    _, err := os.ReadFile("mytest.txt")
    if err != nil {
        err = NewLabelError("file", err)
    }
    fmt.Println(err)
    // должна выводить
    // [FILE] open mytest.txt: no such file or directory
}
```

Решение:
```go
// Error добавляет поддержку интерфейса error для типа LabelError.
func (le *LabelError) Error() string {
    return fmt.Sprintf("[%s] %v", le.Label, le.Err)
}

// NewLabelError упаковывает ошибку err в тип LabelError.
func NewLabelError(label string, err error) error {
    return &LabelError{
        Label: strings.ToUpper(label),
        Err:   err,
    }
}
```
___

## Интроспекция ошибок

Используя спецификатор `%w` или упаковывая ошибки в свои типы, можно создавать длинную последовательность вложенных ошибок. Тогда применять функцию `Unwrap()`, чтобы найти изначальную ошибку, становится неудобно. В пакете `errors` есть ещё две функции для работы с упакованными ошибками: `Is()` и `As()`.

Сначала рассмотрим функцию `Is()`:
```go
func Is(err, target error) bool
```

Эта функция предназначена для сравнения ошибок. Допустим, вы получили ошибку `err` и вам нужно сравнить её с ошибкой `target`. Многие пакеты содержат предопределённые ошибки. 

Например, нужно узнать, равна ли полученная ошибка `err` ошибке `io.EOF`. Использовать сравнение `err == io.EOF` нельзя, так как `err` может содержать обёрнутые ошибки, включая `io.EOF`. Функция `errors.Is()` сравнивает ошибку `err` с `target`, причём ошибка `target` ищется среди всех вложенных ошибок. Функция возвращает `true`, если текущая ошибка `err` равна `target` или содержит ошибку `target`. 

Предположим, потребовалось расширить функциональность: если отсутствует файл настроек, то нужно его создать. Если файла не существует, то `os.ReadFile` будет возвращать ошибку `os.ErrNotExist`, которая определена в пакете `os`. Добавим проверку этой ошибки в предыдущий пример с `TimeError` — в функции `main` потребуются совсем небольшие изменения:
```go
func main() {
    data, err := ReadTextFile("myconfig.yaml")
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            // создаём файл
            // ...
        } else {
            fmt.Println(err)
            os.Exit(0)
        }
    }
    // ...
}
```

Несмотря на то что `ReadTestFile` возвращает ошибку типа `*TimeError`, функция `Is()` вернёт `true`, если начальная ошибка была `os.ErrNotExist`. Без функции `Is()` для правильной работы программы пришлось бы вызывать `Unwrap()` и сравнивать возвращаемую ей ошибку с `os.ErrNotExist`.

Теперь перейдём ко второй функции — функции `As()`. Рассмотрим ситуацию, при которой нужно привести ошибку к определённому типу:
```go
var myErr *TimeError
if myErr, ok = err.(*TimeError); ok {
    // ... 
}
```

Этот код будет работать при отсутствии вложенных ошибок, но если ошибка с нужным типом была обёрнута, то утверждение типа `err.(*TimeError)` закончится неудачей и значение `ok` будет `false`. Функция `errors.As()` проверяет ошибку на соответствие определённому типу и приводит её к таковому, но при этом учитываются все обёрнутые ошибки:
```go
func As(err error, target any) bool
```

```go
var myErr *TimeError
if errors.As(err, &myErr) {
    // ... 
}
```

Если параметр `err` имеет такой же тип, на какой указывает `target`, или содержит обёрнутую ошибку такого же типа, то `As()` присваивает значению `target` эту ошибку и возвращает `true`. В противном случае возвращается `false`. Параметр `target` должен быть ненулевым указателем.

Посмотрите, как можно получить все поля ошибки для типа `TimeError`:
```go
type TimeError struct {
    Time time.Time
    Err  error
}

func main() {
    data, err := ReadTextFile("myconfig.yaml")
    if err != nil {
        var te *TimeError
        if errors.As(err, &te) {
            fmt.Printf("ошибка %v возникла %s", te.Err,
                te.Time.Format("02.01.06 15:04:05"))
        }
        // ...
    }
    // ...
}
```

Так как при создании ошибки функция `NewTimeError()` возвращает указатель, а не структуру, в примере выше `te` — это указатель на `TimeError`.

Если в пакете есть экспортируемая ошибка и в дальнейшем возможны её изменения, то лучше возвращать её в упакованном виде.

```go
// экспортируемая ошибка, доступна из других пакетов
var ErrAccessDenied = errors.New("access denied")

func LoadSettings(userID int) error {
    if !AdminUser(userID) {
        return fmt.Errorf("%w", ErrAccessDenied)
    }
    // загружаем настройки
}
```

Нет смысла упаковывать ошибки в пакете Go, когда:
- Функция выполняет простую операцию и возвращает ошибку из стандартной библиотеки.
- Не предполагается обработка ошибки. Например, если функция используется только в рамках конкретного приложения и возвращаемая ошибка записывается в лог.
- Упаковка ошибок усложняет код. Если упаковка ошибок не облегчает обработку и кажется излишней, то её можно не использовать.

![](./assets/images/packing_error.png)

Закрепим материал примерами правильной и неправильной работы с ошибками:
```go
// плохо
if err == ErrAccessDenied {
}
// хорошо
if errors.Is(err, ErrAccessDenied) {
}

var myErr *MyError
// плохо
if myErr, ok = err.(*MyError); ok {
}
// хорошо
if errors.As(err, &myErr) {
}
```

Разберём небольшой пример классификации ошибок на retryable/non-retryable на примере PostgreSQL.

Retryable-ошибки — это ошибки, которые могут быть исправлены повторной попыткой выполнения операции. Они могут быть вызваны различными причинами, такими как: перегрузка сервера, недоступность сети или баги в коде программы.

Повторную отправку запросов обычно выполняют после некоторой задержки, используя различные политики реализации повторений: fixed delay, exponential backoff, linear backoff, random delay, adaptive backoff.

## Примеры retryable-ошибок и стратегии повтора

Тип ошибки | Краткое описание | Возможная причина | Рекомендуемая стратегия повтора
--- | --- | --- | ---
Кратковременная потеря соединения | Нет соединения, запрос не ушёл или не вернулся ответ | Временная недоступность сети или сервера | Fixed delay, Linear backoff
База данных / сервис не отвечает | БД/сервис отвечает с задержкой или не отвечает | Высокая нагрузка, сбой в передаче | Exponential backoff, Adaptive backoff
Ошибка доступа к заблокированному файлу | Файл занят другим процессом | Конкурентный доступ | Linear backoff, Random delay

## Описание стратегий повтора

Стратегия | Описание | Применение
--- | --- | ---
Fixed delay | Повтор с фиксированной задержкой | Простой и предсказуемый подход при кратковременных сбоях
Exponential backoff | Экспоненциальное увеличение задержки | Эффективен при высокой нагрузке или DDoS
Linear backoff | Линейное увеличение задержки | Более мягкий рост времени ожидания, чем при Exponential
Random delay | Случайная задержка в пределах диапазона | Предотвращает одновременные повторы от множества клиентов
Adaptive backoff | Подстройка задержки на основе типа ошибки и состояния системы | Для оптимизации производительности

Напишем небольшую реализацию классификатора ошибок с конвертацией из интерфейса `error` стандартной библиотеки. Чтобы определить тип ошибки PostgreSQL, с которой завершился запрос, можно воспользоваться библиотекой `github.com/jackc/pgerrcode`:
```go
package pgerrors

import (
    "errors"

    "github.com/jackc/pgerrcode"
    "github.com/jackc/pgx/v5/pgconn"
)

// ErrorClassification тип для классификации ошибок
type PGErrorClassification int

const (
    // NonRetriable - операцию не следует повторять
    NonRetriable PGErrorClassification = iota

    // Retriable - операцию можно повторить
    Retriable
)

// PostgresErrorClassifier классификатор ошибок PostgreSQL
type PostgresErrorClassifier struct{}

func NewPostgresErrorClassifier() *PostgresErrorClassifier {
    return &PostgresErrorClassifier{}
}

// Classify классифицирует ошибку и возвращает PGErrorClassification
func (c *PostgresErrorClassifier) Classify(err error) PGErrorClassification {
    if err == nil {
        return NonRetriable
    }

    // Проверяем и конвертируем в pgconn.PgError, если это возможно
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        return СlassifyPgError(pgErr)
    }

    // По умолчанию считаем ошибку неповторяемой
    return NonRetriable
}

func СlassifyPgError(pgErr *pgconn.PgError) PGErrorClassification {
    // Коды ошибок PostgreSQL: https://www.postgresql.org/docs/current/errcodes-appendix.html

    switch pgErr.Code {
    // Класс 08 - Ошибки соединения
    case pgerrcode.ConnectionException,
        pgerrcode.ConnectionDoesNotExist,
        pgerrcode.ConnectionFailure:
        return Retriable

    // Класс 40 - Откат транзакции
    case pgerrcode.TransactionRollback, // 40000
        pgerrcode.SerializationFailure, // 40001
        pgerrcode.DeadlockDetected:     // 40P01
        return Retriable

    // Класс 57 - Ошибка оператора
    case pgerrcode.CannotConnectNow: // 57P03
        return Retriable
    }

    // Можно добавить более конкретные проверки с использованием констант pgerrcode
    switch pgErr.Code {
    // Класс 22 - Ошибки данных
    case pgerrcode.DataException,
        pgerrcode.NullValueNotAllowedDataException:
        return NonRetriable

    // Класс 23 - Нарушение ограничений целостности
    case pgerrcode.IntegrityConstraintViolation,
        pgerrcode.RestrictViolation,
        pgerrcode.NotNullViolation,
        pgerrcode.ForeignKeyViolation,
        pgerrcode.UniqueViolation,
        pgerrcode.CheckViolation:
        return NonRetriable

    // Класс 42 - Синтаксические ошибки
    case pgerrcode.SyntaxErrorOrAccessRuleViolation,
        pgerrcode.SyntaxError,
        pgerrcode.UndefinedColumn,
        pgerrcode.UndefinedTable,
        pgerrcode.UndefinedFunction:
        return NonRetriable
    }

    // По умолчанию считаем ошибку неповторяемой
    return NonRetriable
}
```

Пример использования:
```go
// executeWithRetry демонстрируем повторение операции с ошибкой
func executeWithRetry(db *pgx.DB) error {
    const maxRetries = 3
    var lastErr error

    classifier := pgerrors.NewPostgresErrorClassifier()

    for attempt := 0; attempt < maxRetries; attempt++ {
        ctx := context.Background()

        _, err := db.Exec(ctx, "полезный SQL запрос")
        if err == nil {
            return nil
        }

        // Определяем классификацию ошибки
        classification := classifier.Classify(err)

        if classification == pgerrors.NonRetriable {
            // Нет смысла повторять, возвращаем ошибку
            fmt.Printf("Непредвиденная ошибка: %v\n", err)
            return err
        }

        // .... делаем что-то полезное
    }

    return fmt.Errorf("операция прервана после %d попыток: %w", maxRetries, lastErr)
}
```

___
В функцию `main` из задания 1 были внесены изменения. Допишите код, чтобы она работала правильно.

```go
// программа должна выводить правильное значение

// ...

func main() {
    _, err := os.ReadFile("mytest.txt")
    if err != nil {
        err = NewLabelError("file", err)
    }
    fmt.Println(errors.Is(err, os.ErrNotExist), err)
    // должна выводить текст:
    // true [FILE] open mytest.txt: no such file or directory
}
```

Решение:
```go
// Unwrap() возвращает исходную ошибку.
func (le *LabelError) Unwrap() error {
    return le.Err
} 
```
___

## Функция Join()

В версии Go 1.20 появилась ещё одна функция:
```go
func Join(errs ...error) error
```

Ей можно передать несколько ошибок, которые она объединит в одну. Если в параметрах `Join()` будут переданы ошибки, равные `nil`, то они будут пропущены. Таким образом, если все параметры равны `nil`, функция тоже вернёт `nil`.

Например, это может быть полезно в случае, когда полученные данные должны пройти валидацию. Если сразу возвращать первую найденную ошибку, пользователь будет злиться, отправляя сотый раз форму с очередной исправленной ошибкой. Лучше сразу найти все ошибки и вернуть их пользователю.

Для ошибки, возвращаемой `Join()`, можно вызывать функции `Is()` и `As()`, но `Uwrap()` будет возвращать `nil`. 

Вот пример использования функции `Join()`:
```go
package main

import (
    "encoding/json"
    "errors"
    "fmt"
)

type Settings struct {
    Host string
    Port int
}

var (
    ErrNoHost = errors.New("Не указан host")
    ErrNoPort = errors.New("Не указан port")
)

func ParseSettings(input string) (*Settings, error) {
    var settings Settings
    err := json.Unmarshal([]byte(input), &settings)
    if err != nil {
        return nil, err
    }
    // находим сразу все ошибки
    var errs []error
    if len(settings.Host) == 0 {
        errs = append(errs, ErrNoHost)
    }
    if settings.Port == 0 {
        errs = append(errs, ErrNoPort)
    }
    return &settings, errors.Join(errs...)
}

func main() {
    settings, err := ParseSettings(`{"host":"localhost", "port": 3000}`)
    fmt.Println(err, settings)

    _, err = ParseSettings("{}")
    fmt.Println(err)
    fmt.Println(errors.Is(err, ErrNoHost), errors.Is(err, ErrNoPort))
}
```

Эта программа успешно разберёт первый JSON и вернёт ошибку во втором случае:
```
<nil> &{localhost 3000}
Не указан host
Не указан port
true true
```

Резюмируя, можно сказать, что программирование на Go невозможно без обработки ошибок. Пакет `errors` содержит небольшое количество функций — Go-разработчику важно понимать назначение каждой из них и уметь использовать эти функции на практике. 

## Дополнительные материалы

- [go.dev/errors](https://pkg.go.dev/errors) — документация пакета `errors`.
- [bitfieldconsulting | Error wrapping in Go](https://bitfieldconsulting.com/golang/wrapping-errors) — упаковка ошибок в Go.

# Обучение Алисы 12

Наш навык для Алисы практически обрёл свою финальную форму — осталась всего пара штрихов. Сегодня добавим возможность регистрировать нового пользователя, ведь без пользователей нам некому будет отсылать и доставлять сообщения.

Изменим файл `internal/store/store.go`. Добавим в интерфейс хранилища метод `RegisterUser()` и особую ошибку:
```go
package store

import (
    "context"
    "errors"
    "time"
)

// ErrConflict указывает на конфликт данных в хранилище.
var ErrConflict = errors.New("data conflict")

type MessageStore interface {
    FindRecepient(ctx context.Context, username string) (userID string, err error)
    ListMessages(ctx context.Context, userID string) ([]Message, error)
    GetMessage(ctx context.Context, id int64) (*Message, error)
    SaveMessage(ctx context.Context, userID string, msg Message) error
    // RegisterUser регистрирует нового пользователя
    RegisterUser(ctx context.Context, userID, username string) error
}

...
```

Ошибку `ErrConflict` будем использовать во всех реализациях хранилища для оповещения о возможном нарушении целостности данных. В частности, здесь используем её для того, чтобы указать на попытку зарегистрировать нового пользователя с уже имеющимся именем.

Дополним файл `internal/store/pg/store.go` новым методом `RegisterUser()`:
```go
package pg

import (
    "context"
    "database/sql"
    "errors"
    "time"

    "github.com/bluegopher/alice-skill/internal/store"
    "github.com/jackc/pgerrcode"
    "github.com/jackc/pgx/v5/pgconn"
)

// Store реализует интерфейс store.MessageStore и позволяет взаимодействовать с СУБД PostgreSQL.
type Store struct {
    // Поле conn содержит объект соединения с СУБД.
    conn *sql.DB
}

...

func (s *Store) RegisterUser(ctx context.Context, userID, username string) error {
    // добавляем новую запись пользователя
    _, err := s.conn.ExecContext(ctx, `
        INSERT INTO users
        (id, username)
        VALUES
        ($1, $2);
    `, userID, username)

    if err != nil {
        // проверяем, что ошибка сигнализирует о потенциальном нарушении целостности данных
        var pgErr *pgconn.PgError
        if errors.As(err, &pgErr) && pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
            err = store.ErrConflict
        }
    }

    return err
}
```

Не забудем также выполнить команду `$ go mod tidy`, чтобы внести в список зависимостей пакет `github.com/jackc/pgerrcode`, который отвечает за интроспекцию специфичных ошибок СУБД PostgreSQL.

Теперь осталось добавить поддержку команды регистрации в хендлер. Обновим код в файле `cmd/skill/app.go`:
```go
func (a *app) webhook(w http.ResponseWriter, r *http.Request) {
    ...

    // текст ответа навыка
    var text string

    switch true {
    case strings.HasPrefix(req.Request.Command, "Отправь"):
        ...
    case strings.HasPrefix(req.Request.Command, "Прочитай"):
        ...

    // пользователь хочет зарегистрироваться
    case strings.HasPrefix(req.Request.Command, "Зарегистрируй"):
        // гипотетическая функция parseRegisterCommand вычленит из запроса 
        // желаемое имя нового пользователя
        username := parseRegisterCommand(req.Request.Command)

        // регистрируем пользователя
        err := a.store.RegisterUser(ctx, req.Session.User.UserID, username)
        // наличие неспецифичной ошибки
        if err != nil && !errors.Is(err, store.ErrConflict) {
            logger.Log.Debug("cannot register user", zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // определяем правильное ответное сообщение пользователю
        text = fmt.Sprintf("Вы успешно зарегистрированы под именем %s", username)
        if errors.Is(err, store.ErrConflict) {
            // ошибка специфична для случая конфликта имён пользователей
            text = "Извините, такое имя уже занято. Попробуйте другое."
        }

    default:
        ...
    }

    // заполняем модель ответа
    resp := models.Response{
        Response: models.ResponsePayload{
            Text: text, // Алиса проговорит текст
        },
        Version: "1.0",
    }

    w.Header().Set("Content-Type", "application/json")

    // сериализуем ответ сервера
    enc := json.NewEncoder(w)
    if err := enc.Encode(resp); err != nil {
        logger.Log.Debug("error encoding response", zap.Error(err))
        return
    }
    logger.Log.Debug("sending HTTP 200 response")
}
```

Вуаля, больше ничего не мешает нашим пользователям обмениваться сообщениями!

В данном инкременте мы научили наш сервис более понятно и дружелюбно сообщать пользователю о возникших ошибках.

В следующем инкременте мы повысим производительность нашего сервиса за счет многопоточного сбора и множествнной записи новых сообщений в БД.

# Инкремент 13

## Задание для трека «Сервис сокращения URL»

Сделайте в таблице базы данных с сокращёнными URL уникальный индекс для поля с исходным URL. Это позволит избавиться от дублирующих записей в БД. 
Используйте инструмент миграций для изменения схемы базы данных.

При попытке пользователя сократить уже имеющийся в базе URL через хендлеры `POST /` и `POST /api/shorten` сервис должен вернуть HTTP-статус `409 Conflict`, а в теле ответа — уже имеющийся сокращённый URL в правильном для хендлера формате.

Стратегии реализации:
1. Чтобы не проверять наличие оригинального URL в базе данных отдельным запросом, можно воспользоваться конструкцией `INSERT ... ON CONFLICT` в PostgreSQL. Однако в таком случае придётся самостоятельно возвращать и проверять собственную ошибку.
1. Чтобы определить тип ошибки PostgreSQL, с которой завершился запрос, можно воспользоваться библиотекой `github.com/jackc/pgerrcode`, в частности `pgerrcode.UniqueViolation`. В таком случае придётся делать дополнительный запрос к хранилищу, чтобы определить сокращённый вариант URL.

## Задание для трека «Сервис сбора метрик и алертинга»

Измените весь свой код в соответствии со знаниями, полученными в этой теме. Добавьте обработку retriable-ошибок.

Сценарии возможных ошибок:
- Агент не сумел с первой попытки выгрузить данные на сервер из-за временной невозможности установить с ним соединение.
- При обращении к PostgreSQL cервер получил ошибку транспорта (из категории Class 08 — Connection Exception).

Стратегия реализации:
- Количество повторов должно быть ограничено тремя дополнительными попытками.
- Интервалы между повторами должны увеличиваться: 1s, 3s, 5s.
- Чтобы определить тип ошибки PostgreSQL, с которой завершился запрос, можно воспользоваться библиотекой `github.com/jackc/pgerrcode`, в частности `pgerrcode.UniqueViolation`.
