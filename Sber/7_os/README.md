# Чтение и запись в файл

Как правило, в любом приложении есть логика, организованная через работу с файлами. Это может быть хранение настроек или конфигураций, логирование событий, использование файлов для экспорта или импорта данных. С подобными задачами вам придётся сталкиваться регулярно, поэтому расскажем, как управлять файлами и директориями. 

В этой теме вы познакомитесь с функциями и методами пакета `os`, которые отвечают за работу с файлами. Вы научитесь создавать и открывать файлы, читать их и записывать в них данные. Узнаете, как работать с временными файлами и директориями.

Представим, что есть два сервиса: поставщик и потребитель событий — `Producer` и `Consumer`. Поставщик сохраняет в файл историю покупок автомобилей в виде JSON-строк, а потребитель — читает историю событий из этого файла. Событие представлено в виде структуры:
```go
type Event struct {
    ID       uint    `json:"id"`
    CarModel string  `json:"car_model"`
    Price    float64 `json:"price"`
}
```

## Открытие и закрытие файла

Первое, что нужно сделать, — открыть файл для записи. Если файла нет, то его нужно создать. Для этого воспользуемся функцией `OpenFile()`, которая позволяет открыть файл в разных режимах:
```go
func OpenFile(name string, flag int, perm FileMode) (*File, error)
```

Здесь параметр `flag` — это битовая маска нужных констант. Список констант:
```go
const (
    // можно выбрать только один режим среди O_RDONLY, O_WRONLY и O_RDWR
    O_RDONLY int = syscall.O_RDONLY // открыть файл в режиме чтения
    O_WRONLY int = syscall.O_WRONLY // открыть файл в режиме записи
    O_RDWR   int = syscall.O_RDWR   // открыть файл в режиме чтения и записи

    // значения для управления поведением файла
    O_APPEND int = syscall.O_APPEND // добавлять новые данные в файл при записи
    O_CREATE int = syscall.O_CREAT  // создать новый файл, если файла не существует
    O_EXCL   int = syscall.O_EXCL   // используется вместе с O_CREATE и возвращает
                                    // ошибку, если файл уже существует
    O_SYNC  int = syscall.O_SYNC  // открыть в режиме синхронного ввода/вывода
    O_TRUNC int = syscall.O_TRUNC // очистить файл при открытии
)
```

Третий параметр `perm` функции `OpenFile()` — это набор модификаторов доступа, с которыми будет создан новый файл, если его не существует. Набор доступов задаётся в виде числа, как в Unix-подобных системах. Если вы не знакомы с ними, почитать о них можно на [linux.com](https://www.linux.com/training-tutorials/understanding-linux-file-permissions/).

Функция `OpenFile()` возвращает указатель на структуру `File`, у которой есть много методов для работы с содержимым файла, так что она реализует основные интерфейсы пакета `io`: `io.Reader`, `io.ReaderAt`, `io.ReaderFrom`, `io.Writer`, `io.WriterAt`, `io.Seeker`, `io.Closer` и другие.

Вот некоторые из методов `*File`:
```go
// Read читает из файла и записывает в буфер p.
func (f *File) Read(p []byte) (n int, err error)

// ReadAt читает из файла со смещением off и записывает в буфер p.
func (f *File) ReadAt(p []byte, off int64) (n int, err error)

// ReadFrom читает из r и записывает в файл.
func (f *File) ReadFrom(r Reader) (n int64, err error)

// Write записывает из буфера p в файл.
func (f *File) Write(p []byte) (n int, err error)

// WriteAt записывает из буфера p в файл со смещением off.
func (f *File) WriteAt(p []byte, off int64) (n int, err error)

// Seek устанавливает смещение offset для следующего чтения
// или записи в зависимости от whence:
// 0 — относительно начала файла,
// 1 — относительно текущего смещения,
// 2 — относительно конца файла.
func (f *File) Seek(offset int64, whence int) (int64, error)

// Close закрывает файл.
func (f *File) Close() error
```

Чтобы начать работу с файлом, его нужно открыть, получив указатель на структуру `File`. После этого можно будет пользоваться её методами для чтения и изменения его содержимого. Пока файл открыт, можно читать и записывать данные. По умолчанию чтение и запись продолжается с того места, где закончилась предыдущая операция. После выполнения всех необходимых действий с файлом его следует закрыть.

*Не забывайте закрывать файл после окончания работы с ним. Для этого сразу после открытия можно использовать `defer f.Close()`. Несвоевременное закрытие файла может быть чревато утечкой ресурсов и потерей данных.*

Но вернёмся к нашему примеру. Со стороны производителя нужно открыть файл с такими флагами:
- `O_WRONLY` — открыть файл в режиме записи;
- `O_CREATE` — если файла не существует, создать новый;
- `O_APPEND` — добавлять новые данные в файл.

А со стороны потребителя — с такими:
- `O_RDONLY` — открыть файл в режиме чтения;
- `O_CREATE` — если файла не существует, создать новый;

Битовые маски в коде можно посчитать так: 
```go
flag1 := os.O_WRONLY | os.O_CREATE | os.O_APPEND

flag2 := os.O_RDONLY | os.O_CREATE
```

Теперь сведём всё вместе и напишем функции открытия и закрытия файла — как для производителя, так и для потребителя.

```go
type Producer struct {
    file *os.File // файл для записи
}

func NewProducer(filename string) (*Producer, error) {
    // открываем файл для записи в конец
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    return &Producer{file: file}, nil
}

func (p *Producer) Close() error {
    // закрываем файл
    return p.file.Close()
}

type Consumer struct {
    file *os.File // файл для чтения
}

func NewConsumer(filename string) (*Consumer, error) {
    // открываем файл для чтения
    file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
    if err != nil {
        return nil, err
    }

    return &Consumer{file: file}, nil
}

func (c *Consumer) Close() error {
    // закрываем файл
    return c.file.Close()
}
```

Стоит сказать, что в пакете `os` есть две простые функции, которые можно использовать, если файл нужно создать или открыть только для чтения:
```go
func Open(name string) (*File, error)
```

Функция `Open()` попытается открыть файл в режиме чтения. Если файла с указанным именем не существует, `Open` вернёт ошибку.

Создать и открыть новый файл можно функцией `Create()`:
```go
func Create(name string) (*File, error) 
```

Она создаёт и открывает файл в режиме чтения и записи. Обратите внимание: если файл уже существует, его содержимое удалится.

___
```go
fd, err := os.Create("fileName.txt", os.O_WRONLY|os.O_APPEND)
```

Что нужно изменить в приведённом выше коде, чтобы открыть существующий файл в режиме чтение/запись и стереть его содержимое?

Верные:
- Оставить только аргумент `"fileName.txt"` у вызова функции `os.Create` (Этот код откроет существующий файл в режиме `O_RDWR` и сотрёт его содержимое.)
- Заменить вызов `os.Create` на `os.OpenFile`, добавить флаг `os.O_TRUNC` и третий аргумент с желаемыми правами доступа (С помощью `os.OpenFile` можно гибко управлять параметрами взаимойдествия с файлом.)

Неверные:
- Оставить код как есть, всё должно работать корректно
- Такие условия соблюсти невозможно: содержимое существующего файла нельзя стереть при открытии
___
Какой код откроет файл для записи в конец файла?

Верные:
- `os.OpenFile("fileName.txt", os.O_RDWR|os.O_APPEND, 644)`
- `os.OpenFile("fileName.txt", os.O_WRONLY|os.O_APPEND, 644)`

Неверные:
- `os.OpenFile("fileName.txt", os.O_RDWR|os.O_CREATE, 644)`
- `os.OpenFile("fileName.txt", os.O_RDONLY, 644)`
___

## Запись и чтение из файла

Теперь на нашем примере рассмотрим, как реализуется запись и чтение данных. Напомним, `Producer` должен записывать строки в JSON-формате в файл, а `Consumer` — последовательно читать эту информацию из того же самого файла.

Реализуем запись событий в файл с помощью метода `Write` структуры `File`:
```go
func (f *File) Write(p []byte) (n int, err error)
```

Метод `Write` сохраняет данные в файл и возвращает количество записанных байт и ошибку.

Поскольку каждое событие хранится в JSON-формате, преобразуем структуру `Event` в слайс байт функцией `json.Marshal`. События разделим символом переноса строки:
```go
func (p *Producer) WriteEvent(event *Event) error {
    data, err := json.Marshal(&event)
    if err != nil {
        return err
    }
    // добавляем перенос строки
    data = append(data, '\n')

    _, err = p.file.Write(data)
    return err
}
```

Теперь реализуем чтение из файла. Для этого нужно:
1. Создать буфер в виде слайса байт.
1. Прочитать методом `Read` данные, записав их в буфер.
1. Найти в буфере символ переноса.
1. Выделить нужную строку.
1. Перенести каретку в файле на символ переноса.

Чтобы не писать этот алгоритм вручную, воспользуемся пакетом [bufio](https://golang.org/pkg/bufio/). Он содержит три типа для работы с буферизованным вводом и выводом: `Reader`, `Writer`, `ReadWriter`.

Перепишем код производителя и потребителя, используя `Reader` и `Writer` из пакета `bufio`:
```go
type Producer struct {
    file *os.File
    // добавляем Writer в Producer
    writer *bufio.Writer
}

func NewProducer(filename string) (*Producer, error) {
    file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    return &Producer{
        file: file,
        // создаём новый Writer
        writer: bufio.NewWriter(file),
    }, nil
}

func (p *Producer) WriteEvent(event *Event) error {
    data, err := json.Marshal(&event)
    if err != nil {
        return err
    }

    // записываем событие в буфер
    if _, err := p.writer.Write(data); err != nil {
        return err
    }

    // добавляем перенос строки
    if err := p.writer.WriteByte('\n'); err != nil {
        return err
    }

    // записываем буфер в файл
    return p.writer.Flush()
}

type Consumer struct {
    file *os.File
    // добавляем reader в Consumer
    reader *bufio.Reader
}

func NewConsumer(filename string) (*Consumer, error) {
    file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
    if err != nil {
        return nil, err
    }

    return &Consumer{
        file: file,
        // создаём новый Reader
        reader: bufio.NewReader(file),
    }, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
    // читаем данные до символа переноса строки
    data, err := c.reader.ReadBytes('\n')
    if err != nil {
        return nil, err
    }

    // преобразуем данные из JSON-представления в структуру
    event := Event{}
    err = json.Unmarshal(data, &event)
    if err != nil {
        return nil, err
    }

    return &event, nil
}
```

*Не забудьте вызвать метод `Writer.Flush` для сброса буфера, иначе данные так и останутся в нём.*

Также в пакете `bufio` есть структура `Scanner`. Она предоставляет удобный интерфейс для чтения данных (например, файла строк). Структура `Scanner` реализует методы:
```go
// Buffer устанавливает свой буфер buf с максимальным размером max.
func (s *Scanner) Buffer(buf []byte, max int)

// Bytes возвращает данные последнего сканирования.
func (s *Scanner) Bytes() []byte

// Err возвращает ошибку сканирования.
// Если scanner достиг конца файла, в качестве ошибки вернётся nil.
func (s *Scanner) Err() error

// Scan производит сканирование до следующего токена (разделителя).
func (s *Scanner) Scan() bool

// Split задаёт свою функцию сканирования.
// По умолчанию используется построчное сканирование.
func (s *Scanner) Split(split SplitFunc)

// Text возвращает данные последнего сканирования в виде строки.
func (s *Scanner) Text() string
```

Перепишем потребитель, используя `Scanner`. Код стал проще. Теперь он выглядит так: 
```go
type Consumer struct {
    file *os.File
    // заменяем Reader на Scanner
    scanner *bufio.Scanner
}

func NewConsumer(filename string) (*Consumer, error) {
    file, err := os.OpenFile(filename, os.O_RDONLY|os.O_CREATE, 0666)
    if err != nil {
        return nil, err
    }

    return &Consumer{
        file: file,
        // создаём новый scanner
        scanner: bufio.NewScanner(file),
    }, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
    // одиночное сканирование до следующей строки
    if !c.scanner.Scan() {
        return nil, c.scanner.Err()
    }
    // читаем данные из scanner
    data := c.scanner.Bytes()

    event := Event{}
    err := json.Unmarshal(data, &event)
    if err != nil {
        return nil, err
    }

    return &event, nil
}

func (c *Consumer) Close() error {
    return c.file.Close()
}
```

Стоит заметить, что если файл одновременно читается и пишется, то буферизацию лучше не использовать. Не будет гарантии, что прочитаны все записанные данные, так как часть из них может находиться в буфере записи.

___
В предыдущих уроках вы изучили сериализацию данных с помощью пакета `json`. Узнали, как работают `json.Encoder` и `json.Decoder`. Реализуйте поставщика и потребителя с помощью `json.Encoder` и `json.Decoder`: 
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
)

type Event struct {
    ID       uint    `json:"id"`
    CarModel string  `json:"car_model"`
    Price    float64 `json:"price"`
}

type Producer struct {
    file    *os.File
    encoder *json.Encoder
}

func NewProducer(filename string) (*Producer, error) {
    // откройте файл и создайте для него json.Encoder
    // допишите код здесь
    // ...
}

func (p *Producer) WriteEvent(event *Event) error {
    // добавьте вызов Encode для сериализации и записи
    // допишите код здесь
    // ...
}

func (p *Producer) Close() error {
    return p.file.Close()
}

type Consumer struct {
    file    *os.File
    decoder *json.Decoder
}

func NewConsumer(filename string) (*Consumer, error) {
    // откройте файл и создайте для него json.Decoder
    // допишите код здесь
    // ...
}

func (c *Consumer) ReadEvent() (*Event, error) {
    // добавьте вызов Decode для чтения и десериализации
    // допишите код здесь
    // ...
}

func (c *Consumer) Close() error {
    return c.file.Close()
}

var events = []*Event{
    {
        ID:       1,
        CarModel: "Lada",
        Price:    400000,
    },
    {
        ID:       2,
        CarModel: "Mitsubishi",
        Price:    650000,
    },
    {
        ID:       3,
        CarModel: "Toyota",
        Price:    800000,
    },
    {
        ID:       4,
        CarModel: "BMW",
        Price:    875000,
    },
    {
        ID:       5,
        CarModel: "Mercedes",
        Price:    999999,
    },
}

func main() {
    fileName := "events.log"
    defer os.Remove(fileName)

    Producer, err := NewProducer(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer Producer.Close()

    Consumer, err := NewConsumer(fileName)
    if err != nil {
        log.Fatal(err)
    }
    defer Consumer.Close()

    for _, event := range events {
        if err := Producer.WriteEvent(event); err != nil {
            log.Fatal(err)
        }

        readEvent, err := Consumer.ReadEvent()
        if err != nil {
            log.Fatal(err)
        }

        fmt.Println(readEvent)
    }
}
```
Решение:
```go
func NewProducer(fileName string) (*Producer, error) {
    file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        return nil, err
    }

    return &Producer{
        file:    file,
        encoder: json.NewEncoder(file),
    }, nil
}

func (p *Producer) WriteEvent(event *Event) error {
    return p.encoder.Encode(&event)
}

func NewConsumer(fileName string) (*Consumer, error) {
    file, err := os.OpenFile(fileName, os.O_RDONLY|os.O_CREATE, 0666)
    if err != nil {
        return nil, err
    }

    return &Consumer{
        file:    file,
        decoder: json.NewDecoder(file),
    }, nil
}

func (c *Consumer) ReadEvent() (*Event, error) {
    event := &Event{}
    if err := c.decoder.Decode(&event); err != nil {
        return nil, err
    }

    return event, nil
}
```
___

## Временные файлы и каталоги

Иногда в приложениях необходимо сохранять промежуточные результаты, для которых нужны временные файлы. Например, при запуске `go run` компилятор создаёт временный исполняемый файл и выполняет его. Пакет `os` содержит несколько функций для создания временных файлов и директорий:
```go
// TempDir возвращает путь до временного каталога (может быть пустым).
func TempDir() string

// MkdirTemp создаёт временный каталог по пути dir.
func MkdirTemp(dir, pattern string) (string, error)

// CreateTemp создаёт временный файл по пути dir.
func CreateTemp(dir, pattern string) (*File, error)
```

Функции `MkdirTemp` и `CreateTemp` принимают `pattern`, который будет использоваться для генерации имени временного файла или каталога. Если `pattern` содержит символ `*`, вместо него подставится случайная последовательность символов. Если же символа `*` нет, случайная последовательность символов добавится к концу `pattern`. Например, по `pattern my*test` сгенерируется имя `my998518050test`, а по `pattern` `mytest` — `mytest3745872346`.

Стоит заметить, что временные файлы не удаляются автоматически по окончании работы программы. Удалять их нужно самостоятельно — методом `Remove` или `RemoveAll`:
```go
// Remove удаляет файл или пустой каталог.
func Remove(name string) error

// RemoveAll удаляет каталог со всеми вложенными файлами и каталогами.
func RemoveAll(path string) error
```

## Сохранение и чтение конфигураций

Теперь узнаем, какие функции пакета `os` используются для сохранения настроек и конфигураций. Как правило,  настройки в файлах не занимают большого объёма и могут быть записаны или прочитаны без дополнительной буферизации. Для записи данных в файл можно использовать функцию `WriteFile(name string, data []byte, perm FileMode) error`, а для чтения — `ReadFile(name string) ([]byte, error)`. Эти функции не требуют дополнительного открытия и закрытия файла: `WriteFile` записывает в файл указанный слайс байт, а `ReadFile` читает всё содержимое файла и возвращает тоже слайс байт.

Рассмотрим пример, как сохранить настройки в `.json`-файл.

```go
package settings

import (
    "encoding/json"
    "os"
    "testing"
)

// Settings содержит настройки приложения.
type Settings struct {
    Port int    `json:"port"`
    Host string `json:"host"`
}

// Save сохраняет настройки в файле fname.
func (settings Settings) Save(fname string) error {
    // сериализуем структуру в JSON формат
    data, err := json.MarshalIndent(settings, "", "   ")
    if err != nil {
        return err
    }
    // сохраняем данные в файл
    return os.WriteFile(fname, data, 0666)
}

func TestSettings(t *testing.T) {
    fname := `settings.json`
    settings := Settings{
        Port: 3000,
        Host: `localhost`,
    }
    if err := settings.Save(fname); err != nil {
        t.Error(err)
    }
}
```

![](./assets/images/bookmarks.png)

___
Дополните пример выше методом, который читает настройки `Settings` из файла.
```go
// Load читает настройки из файла fname.
func (settings *Settings) Load(fname string) error {
    // прочитайте файл с помощью os.ReadFile
    // десериализуйте данные, используя json.Unmarshal
    // ...
}

func TestSettings(t *testing.T) {
    fname := `settings.json`
    settings := Settings{
        Port: 3000,
        Host: `localhost`,
    }
    if err := settings.Save(fname); err != nil {
        t.Error(err)
    }
    var result Settings
    if err := (&result).Load(fname); err != nil {
        t.Error(err)
    }
    if settings != result {
        t.Errorf(`%+v не равно %+v`, settings, result)
    }
    // удалим файл settings.json
    if err := os.Remove(fname); err != nil {
        t.Error(err)
    }
}
```
Решение:
```go
func (settings *Settings) Load(fname string) error {
    data, err := os.ReadFile(fname)
    if err != nil {
        return err
    }
    if err := json.Unmarshal(data, settings); err != nil {
        return err
    }
    return nil
}
```

## Дополнительные материалы

- [go.dev/os](https://pkg.go.dev/os) — документация пакета `os`.
- [go.dev/bufio](https://pkg.go.dev/bufio) — документация пакета `bufio`.
- [The Linux Foundation | Classic SysAdmin: Understanding Linux File Permissions](https://www.linux.com/training-tutorials/understanding-linux-file-permissions/) — модификаторы доступа в Linux.

# Инкремент 9

## Задание по треку «Сервис сокращения URL»

Сохраните все сокращённые URL на диск в виде файла. При перезапуске сервера все URL должны быть восстановлены.

Сервер должен принимать соответствующие параметры конфигурации через флаги и переменные окружения:
- Флаг `-f`, переменная окружения `FILE_STORAGE_PATH` — путь до файла, куда сохраняются данные в формате JSON. Имя файла для значения по умолчанию придумайте сами.

Пример содержимого файла:
```json
[
  {"uuid":"1","short_url":"4rSPg8ap","original_url":"http://yandex.ru"},
  {"uuid":"2","short_url":"edVPg3ks","original_url":"http://ya.ru"},
  {"uuid":"3","short_url":"dG56Hqxm","original_url":"http://practicum.yandex.ru"},
  ...
]
```
Приоритет параметров сервера должен быть таким:
- Если указана переменная окружения, то используется она.
- Если нет переменной окружения, но есть флаг, то используется он.
- Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

## Задание по треку «Сервис сбора метрик и алертинга»

Доработайте код сервера, чтобы он мог с заданной периодичностью сохранять текущие значения метрик на диск в указанный файл, а на старте — опционально загружать сохранённые ранее значения.

Сервер должен принимать соответствующие параметры конфигурации через флаги и переменные окружения:
- Флаг `-i`, переменная окружения `STORE_INTERVAL` — интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение `0` делает запись синхронной).
- Флаг `-f`, переменная окружения `FILE_STORAGE_PATH` — путь до файла, куда сохраняются текущие значения. Имя файла для значения по умолчанию придумайте сами.
- Флаг `-r`, переменная окружения `RESTORE` — булево значение (`true`/`false`), определяющее, следует ли загружать ранее сохранённые значения из указанного файла при старте сервера.

Пример содержимого файла:
```json
[
  {"id":"LastGC","type":"gauge","value":1257894000000000000},
  {"id":"NumGC","type":"counter","delta":42},
  ...
]
```

Приоритет параметров сервера должен быть таким:
- Если указана переменная окружения, то используется она.
- Если нет переменной окружения, но есть флаг, то используется он.
- Если нет ни переменной окружения, ни флага, то используется значение по умолчанию.

# Спринт 2. Код-ревью инкрементов

Как только инкремент будет готов, нажимайте кнопку «Сдать работу». Вас перебросит на вкладку «Ревью».

*Внимание! Для повторной проверки кода нужно снова отправить ссылку на Pull Request через платформу. Только так ревьюер узнает, что вы закоммитили изменения.*