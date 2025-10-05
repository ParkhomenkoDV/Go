# Пакет gomock: имитация баз данных

Обычно разработчики спрашивают, что означает символ `_` перед импортом пакета для работы с базой данных. Или как подключиться к базе данных на Go. Одних интересуют транзакции в SQL, других — конвертация записи из базы данных в структуру. 

![](./assets/images/questions1.png)

Следующая тема посвящена взаимодействию с SQL-базами данных. В рамках её изучения вы: 
- узнаете, как тестировать код, когда база данных еще не готова, но уже почти написана бизнес-логика приложения;
- создадите базы данных и наполните их;
- научитесь выполнять запросы к БД и получать информацию из неё;
- узнаете, как изменять, добавлять и удалять записи из базы данных;
- изучите транзакции и научитесь их применять.

Начнём с инструмента `gomock`, который позволяет имитировать работу базы данных при тестировании.

Как вы знаете, код нужно тестировать, чтобы он был надёжным и живучим. Для юнит-тестирования предназначен пакет `testing`, который обеспечивает поддержку автоматического тестирования пакетов Go. 

Хотя пакет `testing` справляется со своими задачами, есть ситуации, когда нужно прибегнуть к дополнительным инструментам тестирования. Например, когда тестируемый код активно работает с внешним окружением. Не всегда тесты могут отправлять множество запросов к базе данных или сервисам. Более того, запускать тесты на рабочей БД нельзя. Иначе можно поломать производственные процессы. 

Приходится создавать тестовое окружение — песочницу. Чтобы протестировать код, работающий с БД, нужно установить базу данных, создать копию или срез рабочей базы или сгенерировать базу со случайными данными. Эти действия выходят за рамки юнит-тестов: получается интеграционное тестирование. 

![](./assets/images/gomok.png)

Бывает, что в базе данных или сервисе ещё не реализованы нужные возможности, а код, требующий тестирования, уже написан. Во всех этих ситуациях помогут заглушки — **моки** (от англ. mock — «пародия», «имитация»).

Использование моков — это стандартный подход в юнит-тестировании. Мок заменяет реальный объект — зависимость тестируемого класса. 

Сервис описывают интерфейсом. В тестовой среде этот сервис подменивают моком. Заглушка реализует интерфейс, но фактических вычислений не выполняет. 

Писать заглушки вручную — занятие трудоёмкое. Для экономии времени используют пакеты, которые создают моки. Код заглушек автоматически генерируется из объявленных интерфейсов. Один из таких пакетов рассмотрим подробнее. 

## Создаём моки с помощью пакета gomock

Вы уже знакомы с библиотекой [testify](https://github.com/stretchr/testify), которая содержит различные пакеты для тестирования, в том числе пакет `mock`. Рассмотрим другой популярный пакет для работы с моками — [gomock](https://github.com/golang/mock).

Пакет устанавливается так же, как и другие инструменты в Go: 
`go install github.com/golang/mock/mockgen@latest`

Получаем исполняемую команду `mockgen`, которая и генерирует заглушки. `mockgen` работает в двух основных режимах:
1. Режим исходного кода. Мок-интерфейсы генерируются из файла-исходника, указанного в параметре `-source`. Этот режим может работать с неэкспортируемыми интерфейсами.
`mockgen -source=src_file.go`
1. Режим рефлексии (англ. reflection). Мок-интерфейсы генерируются через анализ интерфейсов. Код сначала компилируется, потом разбирается при помощи пакета `reflect`. Вы уже сталкивались с этим пакетом в уроках про структурные теги и encoding. Для запуска в режиме рефлексии нужно указать два позиционных аргумента: путь тестируемого пакета (такой же, как в инструкции `import`) и список разделённых запятой интерфейсов, для которых требуется сгенерировать заглушки. Режим рефлексии можно использовать с директивами `go:generate`.
`mockgen github.com/user/project Store,Driver`

Чтобы сохранить сгенерированные заглушки в файл, укажите его имя в параметре `-destination`.

У `mockgen` есть и другие флаги, детализирующие поведение утилиты. Их описание можно открыть на сайте разработчика или вызвать командой `mockgen -help`.

Попробуем `gomock` в деле. Будем использовать режим рефлексии. Он даёт больше контроля над тем, что, где и когда генерируется.

Рассмотрим хранилище, описанное интерфейсом:
```go
// файл store/store.go

package store

type Store interface {
   Set(key string, value []byte) error
   Get(key string) ([]byte, error)
   Delete(key string) error
}
```

А вот использующая его функция, которую нужно протестировать:
```go
// файл persistent/persistent.go

package persistent

import (
   "project/store"
)

func Lookup(s store.Store, key string) ([]byte, error) {
   // ...
   return s.Get(key)
}
```

Предположим, что пакеты проекта организованы следующим образом:
```
project
|___ store
|    |___ store.go
|
|___ mocks
|    |___ mock_store.go
|
|___ persistent
     |___ persistent.go
     |___ persistent_test.go
```

Папку `mocks` для заглушек нужно создать вручную, а файл `mock_store.go` создаст `mockgen`.

Теперь, чтобы запустить `mockgen` в режиме рефлексии, выполним в корневой директории проекта команду:
`mockgen -destination=mocks/mock_store.go -package=mocks project/store Store`

Разберём параметры этой команды:
- `-destination` указывает имя файла, куда нужно записать сгенерированный код. Если этот флаг не указать, результаты будут просто выведены в консоль.
- `-package` — путь пакета сгенерированных заглушек. По умолчанию будет `mock_store`.
- `project/store` — имя пакета, для интерфейсов которого делаем заглушки.
- `Store` — имя интерфейса, для которого будут созданы моки.

В результате должен сгенерироваться файл `mocks/mock_store.go` с примерно таким кодом:
```go
// Code generated by MockGen. DO NOT EDIT.
// Source: project/store (interfaces: Store)

// Package mocks is a generated GoMock package.
package mocks

import (
   reflect "reflect"

   gomock "github.com/golang/mock/gomock"
)

// MockStore is a mock of Store interface.
type MockStore struct {
   ctrl     *gomock.Controller
   recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
   mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
   mock := &MockStore{ctrl: ctrl}
   mock.recorder = &MockStoreMockRecorder{mock}
   return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
   return m.recorder
}

// ...

// Get mocks base method.
func (m *MockStore) Get(arg0 string) ([]byte, error) {
   m.ctrl.T.Helper()
   ret := m.ctrl.Call(m, "Get", arg0)
   ret0, _ := ret[0].([]byte)
   ret1, _ := ret[1].(error)
   return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStoreMockRecorder) Get(arg0 interface{}) *gomock.Call {
   mr.mock.ctrl.T.Helper()
   return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), arg0)
}

// ...
```

`MockStore` — это тип-заглушка, реализующий интерфейс `Store`. Объекты этого типа создаются функцией `NewMockStore(...)`.  У каждого объекта есть метод `EXPECT()`, который возвращает экземпляр объекта для записи действий — в нашем случае `MockStoreMockRecorder`. `MockStoreMockRecorder` также имеет все методы интерфейса `Store`, но типы параметров у них — `interface{}`.

Чтобы подключить мок к тесту, нужно:
1. Создать контроллер с помощью функции `gomock.NewController(t)`.
1. Создать фиктивный объект. В нашем случае — вызвать `NewMockStore(ctrl)`.
1. Вызвать для объекта метод `EXPECT()`, затем метод-заглушку с нужными параметрами и метод `Return()` c требуемыми возвращаемыми значениями. Например: `m.EXPECT().Get("key").Return("value", nil)`. Это гарантирует, что при вызове конкретного метода с указанными параметрами будут возвращены именно такие значения.

Теперь можно передать объект-заглушку в тестируемую функцию. Реализуем это в тесте:
```go
// файл persistent/persistent_test.go

package persistent

import (
   "project/mocks"
   "testing"

   "github.com/golang/mock/gomock"
   "github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
   // создаём контроллер
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()

   // создаём объект-заглушку
   m := mocks.NewMockStore(ctrl)

   // гарантируем, что заглушка
   // при вызове с аргументом "Key" вернёт "Value"
   value := []byte("Value")
   m.EXPECT().Get("Key").Return(value, nil)

   // тестируем функцию Lookup, передав в неё объект-заглушку
   val, err := Lookup(m, "Key")
   // и проверяем возвращаемые значения
   require.NoError(t, err)
   require.Equal(t, val, value)
}
```

У пакета `gomock` есть функции и методы для создания разнообразных сценариев тестирования. Он позволяет: 
- установить максимальное и минимальное количество вызовов метода,
```go
   m.EXPECT().Get("key").Return("value", nil).MaxTimes(5)
   m.EXPECT().Get("key").MinTimes(1).MaxTimes(10)
```
- принимать любые параметры,
```go
   m.EXPECT().Get(gomock.Any())
```
- вызывать методы в определённом порядке.
```go
gomock.InOrder(
   m.EXPECT().Get("1"),
   m.EXPECT().Get("2"),
   m.EXPECT().Get("3"),
   m.EXPECT().Get("4"),
)
```

Вот так будет выглядеть пример теста с моками:
```go
func TestGet(t *testing.T) {
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()

   m := mocks.NewMockStore(ctrl)

   value := []byte("Value")
   m.EXPECT().
      Get(gomock.Any()).
      Return(value, nil).
      MaxTimes(5)

   for _, s := range []string{"Валерия", "Иван", "Екатерина"} {
      val, err := Lookup(m, s)
      require.NoError(t, err)
      require.Equal(t, val, value)
   }
}
```

___
Реализуйте сценарий тестирования, во время которого заглушка вернёт ошибку при пустом ключе.
```go
package persistent

import (
   "errors"
   "project/mocks"
   "testing"

   "github.com/golang/mock/gomock"
   "github.com/stretchr/testify/require"
)

func TestGet(t *testing.T) {
   ctrl := gomock.NewController(t)
   defer ctrl.Finish()

   m := mocks.NewMockStore(ctrl)

   // возвращаемая ошибка
   errEmptyKey := errors.New("Указан пустой ключ")

   // допишите код
   // ...
}
```

Решение:
```go
   m.EXPECT().Get("").Return([]byte(""), errEmptyKey)
   _, err := Lookup(m, "")
   require.ErrorIs(t, err, errEmptyKey)
```

## Дополнительные материалы

- [go.dev/github.com/golang/mock/gomock](https://pkg.go.dev/github.com/golang/mock/gomock) — документация пакета `gomock`.
- [GitHub | Testing with GoMock: A Tutorial](https://gist.github.com/thiagozs/4276432d12c2e5b152ea15b3f8b0012e) — руководство по тестированию с пакетами `testing` и `gomock`.
- [Medium | A GoMock Quick Start Guide](https://readmedium.com/a-gomock-quick-start-guide-71bee4b3a6f1) — гайд по тестированию с пакетом `gomock`.

# Обучение Алисы 10

Навык Алисы принимает голосовые сообщения и передаёт их между пользователями. Поэтому нам надо позаботиться о том, куда будут сохраняться сообщения и откуда их при необходимости предстоит извлекать. Создадим постоянное хранилище.  

Спроектируем абстрактное хранилище, используя интерфейсы Go. Также добавим новые сценарии в автотесты.

Создадим директорию `internal/store` и поместим туда файл `store.go`:
```
> ~/dev/alice-skill
       |
       |--- cmd
       |     |--- skill
       |            |--- flags.go
       |            |--- gzip.go
       |            |--- main.go
       |            |--- main_test.go
       |--- internal
       |       |--- logger
       |       |      |--- logger.go
       |       |--- models
       |       |      |--- models.go
       |       |--- store
       |              |--- store.go
       |--- go.mod
       |--- go.sum
```

Добавим в новый файл следующий код:
```go
package store

import (
    "context"
    "time"
)

// MessageStore описывает абстрактное хранилище сообщений пользователей
type MessageStore interface {
    // FindRecipient возвращает внутренний идентификатор пользователя по человекопонятному имени
    FindRecipient(ctx context.Context, username string) (userID string, err error)
    // ListMessages возвращает список всех сообщений для определённого получателя
    ListMessages(ctx context.Context, userID string) ([]Message, error)
    // GetMessage возвращает сообщение с определённым ID
    GetMessage(ctx context.Context, id int64) (*Message, error)
    // SaveMessage сохраняет новое сообщение
    SaveMessage(ctx context.Context, userID string, msg Message) error
}

// Message описывает объект сообщения
type Message struct {
    ID      int64     // внутренний идентификатор сообщения
    Sender  string    // отправитель
    Time    time.Time // время отправления
    Payload string    // текст сообщения
}
```

Также, для сохранения данных пользователя в наше хранилище, нам потребуется принимать идентификатор пользователя в запросе к нашему навыку. Для это обновим файл `internal/models/models.go`:
```go
type Session struct {
    New  bool        `json:"new"`
    User RequestUser `json:"user"`
}

// RequestUser содержит данные об авторизованном пользователе навыка
type RequestUser struct {
    UserID string `json:"user_id"`
}
```

Теперь можем добавить интерфейс `store.MessageStore` как зависимость навыка и обновить код для работы с хранилищем сообщений.

Лучше всего при проектировании приложения с внешними зависимостями инкапсулировать их внутрь объекта, реализующего бизнес-логику приложения.
Для этого вынесем `webhook` из файла `cmd/skill/main.go` в файл `cmd/skill/app.go` и превратим его в метод структуры вашего приложения:
```go
package main

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "github.com/bluegopher/alice-skill/internal/logger"
    "github.com/bluegopher/alice-skill/internal/models"
    "github.com/bluegopher/alice-skill/internal/store"
    "go.uber.org/zap"
)

// app инкапсулирует в себя все зависимости и логику приложения
type app struct {
    store store.MessageStore
}

// newApp принимает на вход внешние зависимости приложения и возвращает новый объект app
func newApp(s store.MessageStore) *app {
    return &app{store: s}
}

func (a *app) webhook(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if r.Method != http.MethodPost {
        logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    // десериализуем запрос в структуру модели
    logger.Log.Debug("decoding request")
    var req models.Request
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&req); err != nil {
        logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // проверим, что пришёл запрос понятного типа
    if req.Request.Type != models.TypeSimpleUtterance {
        logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    // получаем список сообщений для текущего пользователя
    messages, err := a.store.ListMessages(ctx, req.Session.User.UserID)
    if err != nil {
        logger.Log.Debug("cannot load messages for user", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    // формируем текст с количеством сообщений
    text := "Для вас нет новых сообщений."
    if len(messages) > 0 {
        text = fmt.Sprintf("Для вас %d новых сообщений.", len(messages))
    }

    // первый запрос новой сессии
    if req.Session.New {
        // обработаем поле Timezone запроса
        tz, err := time.LoadLocation(req.Timezone)
        if err != nil {
            logger.Log.Debug("cannot parse timezone")
            w.WriteHeader(http.StatusBadRequest)
            return
        }

        // получим текущее время в часовом поясе пользователя
        now := time.Now().In(tz)
        hour, minute, _ := now.Clock()

        // формируем новый текст приветствия
        text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, text)
    }

    // заполним модель ответа
    resp := models.Response{
        Response: models.ResponsePayload{
            Text: text, // Алиса проговорит наш новый текст
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

Теперь структура приложения выглядит так:
```
> ~/dev/alice-skill
       |
       |--- cmd
       |     |--- skill
       |            |--- app.go
       |            |--- flags.go
       |            |--- gzip.go
       |            |--- main.go
       |            |--- main_test.go
       |--- internal
       |       |--- logger
       |       |      |--- logger.go
       |       |--- models
       |       |      |--- models.go
       |       |--- store
       |              |--- store.go
       |--- go.mod
       |--- go.sum
```

Изменим код в файле `cmd/skill/main.go`, удалив функцию `webhook`, которая теперь является методом структуры `app`. Добавим следующие строчки:
```go
func run() error {
    if err := logger.Initialize(flagLogLevel); err != nil {
        return err
    }
    
    // создаём экземпляр приложения, пока без внешней зависимости хранилища сообщений
    appInstance := newApp(nil)

    logger.Log.Info("Running server", zap.String("address", flagRunAddr))
    // обернём хендлер webhook в middleware с логированием и поддержкой gzip
    return http.ListenAndServe(flagRunAddr, logger.RequestLogger(gzipMiddleware(appInstance.webhook)))
}
```

Последним шагом нужно удостовериться, что код хендлера корректно работает с хранилищем. Для этого сгенерируем мок для интерфейса хранилища и обновим зависимости проекта:
```bash
$ mockgen -source=internal/store/store.go -destination=internal/store/mock/store.go -package=mock MessageStore
$ go mod tidy
```

Если в дереве проекта появился файл `internal/store/mock/store.go`, значит, операция прошла успешно:
```
> ~/dev/alice-skill
       |
       |--- cmd
       |     |--- skill
       |            |--- app.go
       |            |--- flags.go
       |            |--- gzip.go
       |            |--- main.go
       |            |--- main_test.go
       |--- internal
       |       |--- logger
       |       |      |--- logger.go
       |       |--- models
       |       |      |--- models.go
       |       |--- store
       |              |--- mock
       |                    |--- store.go
       |              |--- store.go
       |--- go.mod
       |--- go.sum
```

Теперь осталось выполнить последний шаг. Дополним автотесты, чтобы они учитывали последние изменения в коде.

Для этого исправим файл `cmd/skill/main_test.go`:
```go
package main

import (
    "bytes"
    "compress/gzip"
    "io"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"

    "github.com/bluegopher/alice-skill/internal/store"
    "github.com/bluegopher/alice-skill/internal/store/mock"
    "github.com/go-resty/resty/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "go.uber.org/mock/gomock"
)

func TestWebhook(t *testing.T) {
    // создадим конроллер моков и экземпляр мок-хранилища
    ctrl := gomock.NewController(t)
    s := mock.NewMockMessageStore(ctrl)

    // определим, какой результат будем получать от «хранилища»
    messages := []store.Message{
        {
            Sender:  "411419e5-f5be-4cdb-83aa-2ca2b6648353",
            Time:    time.Now(),
            Payload: "Hello!",
        },
    }

    // установим условие: при любом вызове метода ListMessages возвращать массив messages без ошибки
    s.EXPECT().
        ListMessages(gomock.Any(), gomock.Any()).
        Return(messages, nil)

    // создадим экземпляр приложения и передадим ему «хранилище»
    appInstance := newApp(s)

    handler := http.HandlerFunc(appInstance.webhook)
    srv := httptest.NewServer(handler)
    defer srv.Close()

    testCases := []struct {
        name         string // добавим название тестов
        method       string
        body         string // добавим тело запроса в табличные тесты
        expectedCode int
        expectedBody string
    }{
        {
            name:         "method_get",
            method:       http.MethodGet,
            expectedCode: http.StatusMethodNotAllowed,
            expectedBody: "",
        },
        {
            name:         "method_put",
            method:       http.MethodPut,
            expectedCode: http.StatusMethodNotAllowed,
            expectedBody: "",
        },
        {
            name:         "method_delete",
            method:       http.MethodDelete,
            expectedCode: http.StatusMethodNotAllowed,
            expectedBody: "",
        },
        {
            name:         "method_post_without_body",
            method:       http.MethodPost,
            expectedCode: http.StatusInternalServerError,
            expectedBody: "",
        },
        {
            name:         "method_post_unsupported_type",
            method:       http.MethodPost,
            body:         `{"request": {"type": "idunno", "command": "do something"}, "version": "1.0"}`,
            expectedCode: http.StatusUnprocessableEntity,
            expectedBody: "",
        },
        {
            name:         "method_post_success",
            method:       http.MethodPost,
            body:         `{"request": {"type": "SimpleUtterance", "command": "sudo do something"}, "session": {"new": true}, "version": "1.0"}`,
            expectedCode: http.StatusOK,
            expectedBody: `Точное время .* часов, .* минут. Для вас 1 новых сообщений.`,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.method, func(t *testing.T) {
            req := resty.New().R()
            req.Method = tc.method
            req.URL = srv.URL

            if len(tc.body) > 0 {
                req.SetHeader("Content-Type", "application/json")
                req.SetBody(tc.body)
            }

            resp, err := req.Send()
            assert.NoError(t, err, "error making HTTP request")

            assert.Equal(t, tc.expectedCode, resp.StatusCode(), "Response code didn't match expected")
            if tc.expectedBody != "" {
                assert.Regexp(t, tc.expectedBody, string(resp.Body()))
            }
        })
    }
}

func TestGzipCompression(t *testing.T) {
    ctrl := gomock.NewController(t)
    s := mock.NewMockMessageStore(ctrl)

    messages := []store.Message{
        {
            Sender:  "411419e5-f5be-4cdb-83aa-2ca2b6648353",
            Time:    time.Now(),
            Payload: "Hello!",
        },
    }

    s.EXPECT().
        ListMessages(gomock.Any(), gomock.Any()).
        Return(messages, nil).
        // удостоверимся, что мок вызывается точно 2 раза
        Times(2)

    appInstance := newApp(s)

    handler := http.HandlerFunc(gzipMiddleware(appInstance.webhook))

    srv := httptest.NewServer(handler)
    defer srv.Close()

    requestBody := `{
        "request": {
            "type": "SimpleUtterance",
            "command": "sudo do something"
        },
        "version": "1.0"
    }`

    // ожидаемое содержимое тела ответа при успешном запросе
    successBody := `{
        "response": {
            "text": "Для вас 1 новых сообщений."
        },
        "version": "1.0"
    }`

    t.Run("sends_gzip", func(t *testing.T) {
        buf := bytes.NewBuffer(nil)
        zb := gzip.NewWriter(buf)
        _, err := zb.Write([]byte(requestBody))
        require.NoError(t, err)
        err = zb.Close()
        require.NoError(t, err)

        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Content-Encoding", "gzip")
        r.Header.Set("Accept-Encoding", "")

        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)

        defer resp.Body.Close()

        b, err := io.ReadAll(resp.Body)
        require.NoError(t, err)
        require.JSONEq(t, successBody, string(b))
    })

    t.Run("accepts_gzip", func(t *testing.T) {
        buf := bytes.NewBufferString(requestBody)
        r := httptest.NewRequest("POST", srv.URL, buf)
        r.RequestURI = ""
        r.Header.Set("Accept-Encoding", "gzip")

        resp, err := http.DefaultClient.Do(r)
        require.NoError(t, err)
        require.Equal(t, http.StatusOK, resp.StatusCode)

        defer resp.Body.Close()

        zr, err := gzip.NewReader(resp.Body)
        require.NoError(t, err)

        b, err := io.ReadAll(zr)
        require.NoError(t, err)

        require.JSONEq(t, successBody, string(b))
    })
}
```

На данном этапе у нас готова вся вспомогательная функциональность: обработка входящих запросов, конфигурация, логирование и работа со сжатыми данными.

В следующий раз подключим навык к настоящей базе данных и допишем автотесты, чтобы проверить корректность работы навыка с хранилищем. Это будет решающий шаг на пути к реализации основной функции нашего сервиса — получения, сохранения и передачи голосовых сообщений между пользователями.

# Абстрактный интерфейс и SQL-драйверы

В Go для работы с данными разработчики используют БД, СУБД (систему управления базами данных) и драйвер, который отправляет SQL-запросы серверу. 

Драйвер — это прослойка между приложением и конкретной СУБД. С помощью соответствующего драйвера можно работать со многими СУБД, в том числе SQLite, MySQL, MariaDB, MS SQL, PostgreSQL, Oracle. Для одной базы данных могут использоваться несколько альтернативных драйверов.

Связать приложение с драйвером БД помогает пакет [database/sql](https://pkg.go.dev/database/sql). Его функции и методы выступают в роли промежуточного интерфейса между приложением, подключённым драйвером и сервером базы данных. Если приложение использует `database/sql`, то для переключения с одной БД на другую достаточно указать новый драйвер и поправить синтаксис SQL-запросов, так как в разных СУБД он может различаться.

Такой подход в работе с SQL-базами данных позволяет писать код, который будет по минимуму привязан к конкретной БД и при необходимости мигрировать на другую СУБД с наименьшими затратами.

![](./assets/images/database.png)

В этом уроке мы будем работать в основном с SQLite — встраиваемой кроссплатформенной СУБД, не использующей парадигму клиент-сервер. Она представляет собой библиотеку, которая становится частью программы. 
Также затронем работу с двумя другими популярными СУБД — MySQL и PostgreSQL. В заданиях вы можете использовать MySQL вместо SQLite.

*PostgreSQL пригодится вам для работы над проектом — обратите на неё особое внимание.*

## Установка MySQL и PostgreSQL

Для начала установим систему управления базами данных. 

Чтобы установить MySQL, [скачайте дистрибутив](https://dev.mysql.com/downloads/mysql/) для вашей операционной системы и запустите его. В процессе установки отметьте варианты по умолчанию, если затрудняетесь с выбором. 

Для Linux можно установить MySQL из официального репозитория. Вот пример для Linux Ubuntu:
```bash
$ sudo apt install mysql-server mysql-client

# далее нужно установить настройки безопасности
$ mysql_secure_installation
```

[На этом сайте](https://metanit.com/sql/mysql/1.1.php) вы найдёте подробную инструкцию по установке MySQL для Windows, [на этом](https://losst.pro/ustanovka-mysql-ubuntu-16-04) — для Linux. А [здесь](https://www.lifewire.com/how-to-install-mysql-on-macos-4586389) вы увидите инструкцию по установке для macOS. 

*Можно также использовать [MariaDB](https://mariadb.org/) — это форк базы кода MySQL, который поддерживается и развивается сообществом.*

Аналогично происходит установка PostgreSQL. Нужно [скачать](https://www.postgresql.org/download/) дистрибутив и запустить его на вашей ОС. Установить PostgreSQL на Linux можно из официального репозитория:
```bash
$ sudo apt -y install postgresql
```

[Здесь](https://winitpro.ru/index.php/2019/10/25/ustanovka-nastrojka-postgresql-v-windows/) вы найдёте подробную инструкцию по установке PostgreSQL на Windows, а [здесь](https://ruvds.com/ru/helpcenter/postgresql-pgadmin-ubuntu/) — на Linux Ubuntu. [Здесь](https://wiki.postgresql.org/wiki/Russian/PostgreSQL-One-click-Installer-Guide) — инструкция по установке на macOS.

Итак, вы установили СУБД. Чтобы начать с ней работать, нужно создать базу данных. Воспользуемся консольным клиентом `mysql` для MySQL и MariaDB, а также `psql` — для PostgreSQL.

Вот пример создания БД PostgreSQL на Linux:
```bash
$ sudo -i -u postgres
$ psql -U postgres
postgres=# create database dbname;
postgres=# create user username with encrypted password 'userpassword';
postgres=# grant all privileges on database dbname to username;
```

Здесь `dbname` — это имя БД, `username` — имя пользователя, а `userpassword` — пароль пользователя.

Чтобы не запутаться в дальнейшем, лучше определить одинаковые имена для базы данных и пользователя. Например, в уроке мы будем работать с базой данных видео, поэтому определим такие имена:
```bash
postgres=# create database videos;
postgres=# create user videos with encrypted password 'userpassword';
postgres=# grant all privileges on database videos to videos;
```

## Установка SQLite

SQLite встраивается в приложение и не требует установки сервера базы данных или дополнительных модулей на ваш компьютер. Эта СУБД отлично подходит для небольших баз данных, которые можно легко скопировать и на другой компьютер.

[Скачайте](https://www.sqlite.org/download.html) комплект консольных утилит для вашей ОС. Вам понадобится программа `sqlite3`, чтобы импортировать `csv`-файл в БД SQLite.

Если вы не привыкли к консоли, можно использовать любой клиент с графическим интерфейсом, который работает с нужной вам СУБД. Например, [DBeaver](https://dbeaver.io/).

## Выбор драйвера

Теперь нужно выбрать драйвер, который обеспечит связь между кодом на Go и СУБД. 

Посмотрите [список доступных драйверов](https://go.dev/wiki/SQLDrivers) на GitHub. Их очень много — больше 50, и для каждой СУБД можно выделить два-три популярных. Поэтому самого лучшего, универсального драйвера не существует.

В таблице приведены наиболее часто используемые драйверы. Обратите внимание: некоторые из них требуют [cgo](https://pkg.go.dev/cmd/cgo), некоторые — нет. Если режим `cgo` включён, можно напрямую исполнять C-код, а результат использовать в Go-коде. Для применения `cgo` нужно дополнительно устанавливать C-компилятор.

База данных | Драйвер | Описание
---|---|---
SQLite | https://github.com/mattn/go-sqlite3 | Один из самых популярных пакетов для SQLite. Хорошая документация. Требует cgo.
SQLite | https://pkg.go.dev/modernc.org/sqlite | Не требует cgo. Написан на чистом Go.
MySQL | https://github.com/go-sql-driver/mysql/ | Не требует cgo. Быстрый и стабильный драйвер, поддерживает запросы более 16 Мб, соединение через IPv4, IPv6 и Unix domain sockets. Подробная документация.
MySQL | https://github.com/go-mysql-org/go-mysql | Это не только драйвер, но и набор утилит. Включает в себя инструменты репликаций, синхронизаций с другими БД ― например, Redis.
PostgreSQL | https://github.com/lib/pq | Популярная раньше библиотека, сейчас практически не поддерживается. Авторы рекомендуют использовать pgx.
PostgreSQL | https://github.com/jackc/pgx | Самый популярный драйвер для проекта с БД PostgreSQL. Быстрый, написан на чистом Go. Можно использовать как через драйвер, так и без пакета `database/sql`. Использование без пакета немного увеличит производительность, а также обеспечит поддержку дополнительных типов атрибутов.

*Помните, что список драйверов периодически обновляется. Могут появиться новые драйверы, которые будут поддерживать больше возможностей и работать быстрее. При написании проекта обязательно проверяйте полный список доступных драйверов.*

## Пакет database/sql

Перейдём к работе с пакетом `database/sql`. Его особенность в том, что он написан обобщённо. Это означает, что все привязки к конкретному движку вынесены в драйвер.

Чтобы получить доступ к базе данных, будем использовать структуру `sql.DB`.

*`sql.DB` — это интерфейс к драйверу базы данных. А драйвер — реализация интерфейса, который содержит логику работы с конкретной СУБД. Драйвер спроектирован так, что в коде не нужно вручную работать с базой данных и достаточно пользоваться только методами `sql.DB`. Код не зависит от драйвера, поэтому драйвер можно легко сменить.*

Далее следует установить соединение с сервером, по которому клиент будет отправлять запросы и получать ответы. На этом этапе важно учесть ещё одну особенность. Сервер выполняет запросы синхронно: клиент должен дождаться ответа, прежде чем отправлять следующий запрос. Если установить одно общее соединение, при большом количестве запросов возникнет очередь и увеличится время обработки.

Если на каждый запрос устанавливать отдельное соединение, можно столкнуться с ограничением количества соединений. Да и такой подход не самый удобный: нужно тратить больше времени и дополнительных ресурсов.

Проблему решает **пул соединений** из пакета `database/sql`. Пул распределяет клиентов по разным соединениям и контролирует количество этих соединений, не тратя время на их повторное открытие.

Пулом соединений управляет структура `sql.DB`. Благодаря ей можно сконфигурировать разные политики: установить фиксированное количество соединений или настроить динамическое открытие новых соединений (при повышенной нагрузке) и закрытие неиспользуемых по тайм-ауту.

Обращение к методам переменной `*sql.DB` потокобезопасно. Можно обращаться к ней из разных горутин.

Функция `sql.Open(driverName, dataSourceName string) (*DB, error)` инициализирует переменную `*sql.DB` и создаёт пул соединений. Первым параметром функции `Open()` указывают драйвер, который будет использоваться для общения с базой данных. Второй аргумент — это строка для драйвера (её формат определяется самим драйвером), где указана необходимая информация для подключения к БД.

Перед открытием нужно импортировать пакет с используемым драйвером:
```go
import (
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)
```

*Обратите внимание: пакет с драйвером импортирован анонимно, поскольку нет необходимости обращаться напрямую к его функциям и методам. При импорте будет вызвана `init()`-функция, в которой драйвер зарегистрирует себя самостоятельно и станет доступен для использования через `sql.DB`.*

После использования база данных закрывается методом `(db *DB) Close() error`. Открывать и закрывать соединение в каждой функции не нужно: это снижает скорость работы программы. Пакет `database/sql` написан с расчётом на долговременное использование структуры `sql.DB`, что позволяет эффективно работать с пулом подключений.

Обычно функция `Open()` вызывается один раз при запуске программы, а метод `db.Close()` — при завершении работы.

Теперь рассмотрим, как подключаться к разным БД.

## MySQL

Для установки MySQL необходим драйвер `go-sql-driver/mysql`.

При использовании `go-sql-driver/mysql` нужно указать имя драйвера `mysql`. Информация для подключения передаётся в формате `[username[:password]@][protocol[(address)]]/dbname`. В квадратных скобках — опциональные параметры. Иногда может быть указано только имя базы данных — `/dbname`.

```go
package main

import (
    "database/sql"
    _ "github.com/go-sql-driver/mysql"
)

func main() {
    db, err := sql.Open("mysql", "user:password@/dbname")
    // вот полный вариант
    // db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/dbname")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    // ...
}
```

## PostgreSQL

Для подключения к PostgreSQL возьмём драйвер `pgx`. Формат строки для подключения будет таким: `host=host port=port user=myuser password=xxxx dbname=mydb sslmode=disable`. Если используется стандартный для PostgreSQL порт `5432`, его можно не указывать.

```go
package main

import (
    "database/sql"
    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    ps := fmt.Sprintf("host=%s user=%s password=%s dbname=%s sslmode=disable",
                      `localhost`, `video`, `XXXXXXXX`, `video`)

    db, err := sql.Open("pgx", ps)
    if err != nil {
        panic(err)
    }
    defer db.Close()
    // ...
}
```

## SQLite

Чтобы подключить SQLite, воспользуемся драйвером `modernc.org/sqlite`, который реализован на чистом Go. В параметрах подключения укажем имя файла базы данных. Если файла не существует, то он будет создан.

*Если вы работаете на Linux или Windows и у вас уже установлен C-компилятор, то для работы с SQLite вы можете использовать драйвер `github.com/mattn/go-sqlite3`.*

```go
package main

import (
    "database/sql"

    _ "modernc.org/sqlite"
)

func main() {
    // при использовании пакета go-sqlite3 имя драйвера — sqlite3
    db, err := sql.Open("sqlite", "video.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()
    // ...
}
```

Напомним, что в примерах мы будем использовать SQLite. Вы можете выбрать любую другую SQL-базу данных для самостоятельного разбора примеров. 

## Проверка подключения к БД

В коде важно указать верные параметры доступа к базе данных: имя БД, пароль, логин. Иначе функция `sql.Open()` не вернёт ошибку сразу. Вы получите её в тот момент, когда попытаетесь отправить SQL-запрос. Так происходит потому, что при использовании пакета `database/sql` соединение с БД устанавливается только после выполнения первого запроса. 

Можно использовать метод `(db *DB) PingContext(ctx context.Context) error` — он проверяет соединение и при необходимости восстанавливает его. Этот метод принимает контекст и возвращает ошибку, если соединение отсутствует и восстановить его невозможно.

```go
package main

import (
    "context"
    "database/sql"
    "time"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "video.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    if err = db.PingContext(ctx); err != nil {
        panic(err)
    }
    // ...
}
```

___
Что происходит при импорте драйвера базы данных?
~~Отрабатывают `init()`-функции пакета драйвера.~~
___
Можно ли использовать `sql.DB` из нескольких потоков исполнения горутин?
~~Да. Структура `sql.DB` потокобезопасна и может быть использована параллельно из нескольких горутин.~~
___

В следующих уроках вы научитесь получать данные из таблицы, применять подготовленные SQL-запросы и транзакции. А пока подготовим базу данных. Использовать будем SQLite.

Для работы возьмём готовый набор данных [о популярных в США видеороликах на YouTube](https://www.kaggle.com/datasnaek/youtube-new?select=USvideos.csv). Скачайте `csv`-файл [здесь](https://code.s3.yandex.net/go/USvideos.csv).

Вы можете открыть файл в текстовом редакторе и посмотреть его содержимое. Файл состоит из множества полей, но нам нужны лишь несколько из них:

атрибут | описание
---|---
`video_id` | `id` видео на YouTube, например: `2kyS6SvSYSE`
`trending_date` | дата, когда видео было в топе: `17.14.11`
`title` | название: `WE WANT TO TALK ABOUT OUR MARRIAGE`
`channel_title` | название канала: `CaseyNeistat`
`category_id` | `id` категории видео: `2`
`publish_time` | время публикации: `2017-11-13T17:13:01.000Z`
`tags` | теги через `|`: `"Shantell"|"martin"|`
`views` | количество просмотров: `4949484`
`likes` | количество лайков: `488499`
`dislikes` | количество дизлайков: `8484`
`comment_count` | количество комментариев: `4884`
`thumbnail_link` | ссылка на заставку: `https://i.ytimg.com/vi/2kyS6SvSYSE/default.jpg`
`comments_disabled` | отключены ли комментарии
`ratings_disabled` | отключён ли рейтинг видео
`video_eror_or_removed` | было ли видео удалено
`description` | описание видео

*Для дальнейшего импорта из `csv`-файла в БД вам понадобится консольный клиент `sqlite3`. [Скачайте](https://sqlite.org/download.html) `zip`-архив с готовыми исполняемыми файлами для вашей ОС.*

Воспользуйтесь `sqlite3` и загрузите данные напрямую из `csv` в новую таблицу:
1. Откройте базу данных. Если укажете несуществующий файл, `sqlite3` просто создаст новую базу.
```bash
$ ./sqlite3 video.db
```
2. Создайте таблицу `videos`. Можно и не создавать, но тогда при импорте `csv`-файла все поля по умолчанию получат тип `TEXT`. С такой таблицей в Go работать неудобно, её трудно конвертировать в структуру языка — потеряется типизация.
```sql
CREATE TABLE videos(
"video_id" TEXT,
"trending_date" TEXT,
"title" TEXT,
"channel_title" TEXT,
"category_id" INTEGER,
"publish_time" TEXT,
"tags" TEXT,
"views" INTEGER,
"likes" INTEGER,
"dislikes" INTEGER,
"comment_count" INTEGER,
"thumbnail_link" TEXT,
"comments_disabled" BOOLEAN,
"ratings_disabled" BOOLEAN,
"video_error_or_removed" BOOLEAN,
"description" TEXT
);
```
У SQLite — динамическая типизация данных с четырьмя типами: `INTEGER`, `REAL`, `TEXT` и `BLOB`. Значение колонки может быть любым из этих типов, независимо от типа, указанного при создании таблицы.

3. Положите файл `USvideos.csv` в ту же директорию, куда распакованы исполняемые файлы `sqlite`. Переключитесь в режим `csv` и импортируйте данные из `csv`-файла в таблицу `videos`:
```bash
sqlite> .mode csv videos
sqlite> .import USvideos.csv videos
```

Теперь проверьте количество записей в базе, чтобы понять, с каким объёмом данных предстоит работать: 
```sql
sqlite> SELECT COUNT(*) From videos;
40950
```

Если откроете `csv`-файл, то увидите, что количество строк будет больше 48 000. Дело в том, что некоторые записи содержат текст с переводом строк в двойных кавычках.

Проверьте, корректно ли импортировались данные о видеороликах. Например, выведите 10 самых непопулярных роликов. Перед этим не забудьте переключиться в режим таблицы, так как при импорте использовался режим `csv`:
```sql
sqlite> .mode table
sqlite> SELECT video_id, title, views from videos ORDER BY views LIMIT 10;
```

Если всё выполнено верно, вывод SQL-команды будет примерно таким: 
video_id | title | views
--- | --- | ---
y6KYFcta4SE | 1 dead, others injured after Ky. school shooting | 554
y6KYFcta4SE | 1 dead, others injured after Ky. school shooting | 559
dQMZLXaa1L8 | Artwork Forge | 658
-JVITToppE0 | Coach Taggart Monday Presser Ahead of Arizona | 687
-JVITToppE0 | Coach Taggart Monday Presser Ahead of Arizona | 704
dQMZLXaa1L8 | Artwork Forge | 713
dQMZLXaa1L8 | Artwork Forge | 745
zeQaJGkFyqQ | Raw: 3 South Carolina Deputies, 1 Officer Shot | 748
qg0GdM60syI | Huffy Metaloid Bicycle Commercial 1997 | 773
JNv4w6DFoYs | OjO Commuter Scooter®- Light Electric Vehicle for the Bike Lane & Beyond! | 789

Данные из `csv`-файла импортированы корректно, структура данных передана верно.

*Обратите внимание: некоторые сущности в файле повторяются. Учтите это при дальнейшей работе.*

Теперь найдите самый просматриваемый ролик:
```sql
sqlite> SELECT title, channel_title, views from videos 
        WHERE views = (SELECT MAX(views) from videos);
```

Вместо реальной записи выводится `title`, `channel_title` и `views`. Всё просто. Так как вы предварительно создали таблицу и уже потом импортировали данные, первая строка из `csv`-файла с заголовками столбцов попала в таблицу. Удалите эту запись:
```sql
sqlite> DELETE from videos where video_id = 'video_id';
```

Запрос выполнился — в базе больше не существует такой записи. Снова запросите видео с максимальным количеством просмотров:
```sql
sqlite> SELECT title, channel_title, views from videos 
        WHERE views = (SELECT MAX(views) from videos);
```

Вывод будет таким:
title | channel_title | views
--- | --- | ---
Childish Gambino — This Is America (Official Video) | ChildishGambinoVEVO | 225211923

Вот и всё! C помощью консольной программы `sqlite3` вы создали таблицу в базе данных SQLite и импортировали туда информацию из `csv`-файла. Эта таблица пригодится вам в следующем уроке, где мы расскажем, как на Go выполнять `SELECT`-запросы к базе данных.

___
При работе с БД нужно использовать консольный клиент. Почему все SQL-запросы не пишут сразу на Go? Выберите верные утверждения. 
- ~~В Go могут быть недоступны некоторые удобные возможности консольного клиента. Например, `sqlite> .mode csv` и `sqlite> .import file.csv tablename`.~~
- ~~Чтобы проверить корректность SQL-запроса, быстрее и удобнее вывести из консоли, чем переписать на Go и перекомпилировать.~~

## Дополнительные материалы

- [go.dev/database/sql](https://pkg.go.dev/database/sql) — документация пакета `database/sql`.
- [SQLDrivers](https://go.dev/wiki/SQLDrivers) — список доступных драйверов.
- [Хабр | Go и MySQL: настраиваем пул соединений](https://habr.com/ru/company/oleg-bunin/blog/583558/) — статья о MySQL.
- [Golang Bot | MySQL Tutorial: Connecting to MySQL and Creating a DB using Go](https://golangbot.com/connect-create-db-mysql/) — руководство по работе с MySQL.
- [Go database/sql tutorial](http://go-database-sql.org/) — популярный самоучитель по работе с БД в Go.
- [Soham Kamani | Using an SQL Database in Go](https://www.sohamkamani.com/golang/sql-database/) — о работе с БД SQL в Go.
- [Go Linux Cloud | Golang SQLite3 Tutorial](https://www.golinuxcloud.com/golang-sqlite3/) — руководство по работе с SQLite3 с примерами.
- [Ardan Labs | Working with SQLite using Go and Python](https://www.ardanlabs.com/blog/2020/11/working-with-sqlite-using-go-python.html) — о работе с SQLite3.
- [YouTube | Go Northwest | David Crawshaw SQLite and Go](https://www.youtube.com/watch?v=RqubKSF3wig&ab_channel=GoNorthwest) — видео о Go и SQLite3.

# Инкремент 10

## Задание по треку «Сервис сокращения URL»

- Добавьте в сервис функциональность подключения к базе данных. В качестве СУБД используйте PostgreSQL не ниже 10 версии.
- Добавьте в сервис хендлер `GET /ping`, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус `200 OK`, при неуспешной — `500 Internal Server Error`.

Строка с адресом подключения к БД должна получаться из переменной окружения `DATABASE_DSN` или флага командной строки `-d`.

Для работы с БД используйте один из следующих пакетов:
- `database/sql`,
- `github.com/jackc/pgx`,
- `github.com/lib/pq`,
- `github.com/jmoiron/sqlx`.

## Задание по треку «Сервис сбора метрик и алертинга»

Сервер:
- Добавьте функциональность подключения к базе данных. В качестве СУБД используйте PostgreSQL не ниже 10 версии.
- Добавьте в сервер хендлер `GET /ping`, который при запросе проверяет соединение с базой данных. При успешной проверке хендлер должен вернуть HTTP-статус `200 OK`, при неуспешной — `500 Internal Server Error`.

Строка с адресом подключения к БД должна получаться из переменной окружения `DATABASE_DSN` или флага командной строки `-d`.

Для работы с БД используйте один из следующих пакетов:
- `database/sql`,
- `github.com/jackc/pgx`,
- `github.com/lib/pq`,
- `github.com/jmoiron/sqlx`.

# Запросы к базе данных

Итак, вы уже установили базу данных и подключились к ней. И даже создали SQLite БД с таблицей `videos`. Теперь вы начнёте работать с этой базой видеороликов на Go. Но для этого нужно обучиться языку SQL-запросов. Их условно делят на две группы:
- запросы, которые возвращают данные: `SELECT`-запросы;
- запросы без возвращения данных: `CREATE`, `UPDATE`, `DROP`, `DELETE` и так далее.

В этом уроке рассмотрим, как выполнять `SELECT`-запросы (англ. select — выбирать). Вы научитесь отправлять запросы и обрабатывать ответы на Go. Также разберём основные SQL-запросы, которые понадобятся для выполнения следующего инкремента спринта. Без этих знаний невозможно работать с базами данных.

## SELECT-запрос

`SELECT`-запрос возвращает данные из таблицы в соответствии с указанными в нём параметрами. Рассмотрим базовый синтаксис `SELECT`-запросов. Простейший запрос для получения всех данных из таблицы выглядит так:
```sql
-- получаем все записи со всеми полями из таблицы videos
SELECT * FROM videos;
```

Допустим, нам не нужны все колонки из таблицы `videos`. В таком случае указываем только нужные поля:
```sql
-- все записи со значениями полей video_id,title,views из таблицы videos
SELECT video_id,title,views FROM videos;
```

Если нужно отфильтровать записи по определённым условиям, то нужно указать `WHERE` с требуемыми условиями. Вы можете использовать:
- операции сравнения: `=`, `<`, `>`, `<=`, `>=`, `!=`;
- операторы `AND` и `OR`;
- операторы `IN` и `NOT IN` для фильтрации по нескольким значениям;
- оператор `BETWEEN ... AND` для указания диапазона.

Вот так выглядит запрос, если задать определённые условия:  
```sql
-- записи с количеством просмотров от 1000 до 2000
SELECT * FROM videos WHERE views >= 1000 AND views <= 2000;

-- эквивалент
SELECT * FROM videos WHERE views BETWEEN 1000 AND 2000;

-- записи с количеством лайков - 0, 5 или 10
SELECT title, likes FROM videos WHERE likes IN (0,5,10);
```

Для агрегации данных, то есть преобразования набора данных в единственное результирующее значение, следует использовать `GROUP BY` и агрегирующие функции:
- `COUNT()` — количество записей;
- `SUM()` — сумма числовых значений;
- `MAX()`, `MIN()` — максимальное и минимальное значение;
- `AVG()` — среднее значение.

Если в `GROUP BY` нужна фильтрация по агрегированным данным, то вместо `WHERE` используется `HAVING`.

```sql
-- группируем записи по лайкам от 0 до 9 и получаем для каждой группы количество
-- записей и среднее количество просмотров в группе
SELECT COUNT(*) AS count_likes, AVG(views) AS avg_views, likes 
FROM videos GROUP BY likes HAVING likes < 10;
```

Ключевое слово `AS` определяет алиас для столбца данных. В примере выше у колонок будут имена `count_likes` и `avg_views`, а не `COUNT(*)` и `AVG(views)`.

По умолчанию порядок строк в результатах не гарантирован и может быть каждый раз разным. Чтобы упорядочить строки, можно воспользоваться командой `ORDER BY column`. Она сортирует возвращаемые записи по возрастанию значений указанного столбца. Если требуется сортировка по убыванию, то после имени поля нужно добавить слово `DESC`.

```sql
-- записи с просмотрами меньше 1000 отсортированы по возрастанию
SELECT * FROM videos WHERE views < 1000 ORDER BY views;

-- записи отсортированы по убыванию лайков
SELECT * FROM videos ORDER BY likes DESC;
```

Если нужно ограничить количество получаемых записей, то следует воспользоваться LIMIT:
```sql
-- первые десять роликов с самым большим количеством просмотров
SELECT * FROM videos ORDER BY views DESC LIMIT 10;
```

## Памятка по SQL-командам

### Создание и удаление таблиц

С возможностями команды `SELECT` вы уже знакомы. Теперь рассмотрим команды, которые, в отличие от `SELECT`, могут изменять данные.

`CREATE TABLE` используется для создания таблицы. Эта команда вам встречалась при создании таблицы `videos` из `sqlite3`.

```sql
-- создаём таблицу movies
CREATE TABLE movies (
    "id" INTEGER PRIMARY KEY,
    "title" VARCHAR(250) NOT NULL DEFAULT '',
    "created" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    "tags" TEXT,
    "views" INTEGER NOT NULL DEFAULT 0
)
```

Объясним значения некоторых ключевых слов в этом выражении:
- `PRIMARY KEY` определяет первичный ключ, для которого автоматически создаётся индекс;
- `VARCHAR(250)` означает строковый тип переменного размера (до 250 байт);
- `NOT NULL` указывает на то, что значение колонки не может принимать `NULL`;
- `DEFAULT` определяет значение по умолчанию, если при добавлении записи значение колонки не указано;
- `CURRENT_TIMESTAMP` означает подстановку текущего времени.

Команда `DROP TABLE` удаляет таблицу.

```sql
-- удаляем таблицу videos
DROP TABLE videos
```

Заметьте, что `DROP TABLE` полностью удаляет таблицу из БД. Если нужно удалить все записи, но оставить структуру таблицы без изменений, то нужно использовать `TRUNCATE TABLE`.

```sql
-- удаляем все записи в таблице videos
TRUNCATE TABLE videos
```

### Добавление новых записей

Для вставки записей используется команда `INSERT INTO`. 

```sql
-- вставляем запись, у которой указаны значения только четырёх колонок
INSERT INTO videos (video_id, title, views, likes)
    VALUES ('1', 'Something about', 100, 20)
```

Можно вставлять несколько записей одновременно:
```sql
-- вставляем три записи
INSERT INTO videos (video_id, title, views, likes)
    VALUES ('1', 'Something about', 100, 20),
           ('2', 'Happy birthday', 150, 31),
           ('3', 'To be or not to be', 97, 12)
```

### Изменение записей

Если у записей необходимо изменить отдельные поля, то используется команда `UPDATE ... SET`. Как правило, она применяется с условием `WHERE`. Условие указывается так же, как в команде `SELECT`.

```sql
-- изменяем заголовок, количество просмотров и лайков у роликов 
-- с video-id, равным 2kyS6SvSYSE
UPDATE videos SET title = 'The best movie', views = 10001, likes = 578
    WHERE video_id = '2kyS6SvSYSE'

-- обнуляем количество лайков, если их меньше 100
UPDATE videos SET likes = 0 WHERE likes < 100
```

### Удаление записей

Команда `DELETE FROM` служит для удаления записей. Её нужно применять с условием `WHERE`. Иначе все записи в таблице будут удалены.

```sql
-- удаляем записи, у которых количество просмотров меньше 100000
DELETE FROM videos WHERE views < 100000
```

Переходим к следующему этапу урока. Теперь вам предстоит понять логику SQL-выражений.

## Методы Go для отправки запросов

Сосредоточимся на языке Go и вернёмся к пакету `database/sql`. Для отправки запросов в нём есть три основных метода:
```go
(db *DB) ExecContext(ctx context.Context, query string, args ...any) (Result, error)
(db *DB) QueryContext(ctx context.Context, query string, args ...any) (*Rows, error)
(db *DB) QueryRowContext(ctx context.Context, query string, args ...any) *Row
```

Есть аналогичные методы без контекста: `Exec(...)`, `Query(...)` и `QueryRow(...)`. Какие методы выбрать, решает сам разработчик. Но в код лучше сразу заложить использование контекста.

Методы, начинающиеся с `Query`, используются для отправки `SELECT`-запросов:
- `QueryContext()` и `Query()` возвращают все полученные записи;
- `QueryRowContext()` и `QueryRow()` возвращают только одну запись.

Методы `ExecContext()` и `Exec()` используются для запросов без возвращения данных: `INSERT`, `DELETE`, `UPDATE` и так далее. Как с ними работать, подробно расскажем в следующем уроке. 

*Важно: контекст в операциях с БД присутствует всегда. Методы `Query()` или `Exec()` под капотом используют `context.Background()`. Если вы пользуетесь инструментарием `context.Context`, нужно применять методы `QueryContext()` и `ExecContext()`, позволяющие установить контекст явно.*

Все `Query()`-функции возвращают значения типа `*Rows` или `*Row`. Чтобы перевести полученные данные в Go-типы, нужно использовать методы `Scan()`. Они последовательно копируют колонки текущей записи в переменные или области, на которые указывают передаваемые параметры.

```go
func (r *Row) Scan(dest ...any) error
func (rs *Rows) Scan(dest ...any) error
```

Теперь нам понадобится база данных с видеороликами. Попробуем отправить простейший запрос на Go. Посчитаем количество всех записей в базе данных:
```go
package main

import (
    "context"
    "database/sql"
    "fmt"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "video.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // делаем запрос
    row := db.QueryRowContext(context.Background(),
        "SELECT COUNT(*) as count FROM videos")
    // готовим переменную для чтения результата
    var id int64
    err = row.Scan(&id)  // разбираем результат
    if err != nil {
        panic(err)
    }
    fmt.Println(id)
}
```

Порядок действий такой:
1. Выполняем запрос к базе данных. Так как результатом всегда будет одна запись, используем не `QueryContext()`, а `QueryRowContext()`. Это позволит сократить объём кода.
1. Вызываем метод `Scan()`, который записывает результат в переменную. В случае нескольких значений указываем их последовательно. Так же, как выводим в SQL.
1. В случае ошибки выводим сообщение об ошибке.

В качестве контекста передали `context.Background()`. Контекст позволяет ограничить по времени или прервать слишком долгие или уже не нужные операции с базой данных. Назначить для них дедлайн или тайм-аут. Вот так выглядит установка трёхсекундного тайм-аута, после которого затянувшаяся операция с БД будет прервана:
```go
ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
// не забываем освободить ресурс
defer cancel()
// обращаемся к БД
row := db.QueryRowContext(ctx, "SELECT COUNT(*) as count FROM videos")
```

Если вы делаете SQL-запрос к БД из функций-обработчиков (`http.Handler`) HTTP-сервера, то наследуйте контекст HTTP-запроса:
```go
func MyHandler(w http.ResponseWriter, r *http.Request) {
    // наследуем контекст запроса r *http.Request, оснащая его Timeout
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()
    // делаем обращение к db в рамках полученного контекста
    row := db.QueryRowContext(ctx, "SELECT COUNT(*) as count FROM videos")
    // ...
}
```

___
Выберите верные утверждения.

Верные:
- Для отправки `SELECT`-запроса в БД нужно использовать методы `Query...()`. (Методы, начинающиеся с `Query`, используются для отправки `SELECT`-запросов.)
- Метод `Scan()` преобразует полученные данные в нужные Go-типы. (`Scan()` последовательно копирует столбцы записи в переменные или области, на которые указывают передаваемые параметры.)

Неверные:
- Запрос `SELECT * FROM users LIMIT 100` выдаст ошибку, если количество записей в таблице `users` меньше 100. (`LIMIT` только устанавливает максимальное количество возвращаемых записей.)
- Метод `ExecContext()` можно использовать для отправки любых SQL-запросов. (`ExecContext()` не подходит для SELECT-запросов.)
___

Разберём подробнее работу метода `Scan()`. Все значения, полученные в результате SQL-запроса в `*sql.Rows` и `*sql.Row`, имеют тип `interface{}`. Работать с этим типом не совсем удобно, поэтому и используется метод `Scan()`. Он преобразует значения в тот тип переменных, который нужен для дальнейших действий.

![](./assets/images/scan.png)

Метод `Scan()` поддерживает: 
- Указатели на многие встроенные типы языка Go ― например, `int64`, `float64`, `bool`, `[]byte`, `string`, `time.Time`, `nil`.
- Специфические типы пакета `database/sql` ― например, `*sql.NullString`.
- Любые пользовательские типы данных, реализующие интерфейсы `sql.Scanner` и `driver.Valuer` (пакет `database/sql/driver`).

Рассмотрим пример. В таблице `videos` есть колонка `comments_disabled` типа `BOOLEAN` со значениями `TRUE` или `FALSE`. Но это могла бы быть и колонка типа `TEXT` со значениями: `true`, `T`, `TRUE`, `1` или `0`, `false`, `F`. Если в `Scan()` для этого поля передать ссылку на строку, то после, чтобы получить логическое значение, нужно использовать `strconv.ParseBool()`. Если же сразу передать ссылку на переменную типа `bool`, `Scan()` сам сделает преобразования и запишет логическое значение по указанному адресу или вернёт ошибку.

Вот как получить и разобрать три значения разных типов из таблицы `videos`. Попробуйте запустить эту программу и посмотрите на результат.

```go
package main

import (
    "context"
    "database/sql"
    "fmt"

    _ "modernc.org/sqlite"
)

func main() {
    db, err := sql.Open("sqlite", "video.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    row := db.QueryRowContext(context.Background(),
        "SELECT title, likes, comments_disabled "+
            "FROM videos ORDER BY likes DESC LIMIT 1")
    var (
        title  string
        likes  int
        comdis bool
    )
    // порядок переменных должен соответствовать порядку колонок в запросе
    err = row.Scan(&title, &likes, &comdis)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%s | %d | %t \r\n", title, likes, comdis)
}
```

## Передача параметров

До этого все рассмотренные нами SQL-запросы не содержали переменных величин. Никогда не используйте для этого функцию `fmt.Sprintf` или объединение строк. 

```go
// если в id окажется "ooops; DROP TABLE video", 
// то вам придётся создавать таблицу video заново
query := fmt.Sprintf(`SELECT title FROM video WHERE video_id = %s`, id)
```

Вот пример правильной подстановки значения:
```go
row := db.QueryRowContext(ctx, "SELECT title FROM videos WHERE video_id = ?", id)
```

Для подстановки значений в запрос используются специальные параметры `?`, `$N`, `:param`, которые зависят от базы данных.

```go
// запрос для SQLite и MySQL
db.Query("SELECT * FROM video WHERE likes BETWEEN ? AND ?", minVal, maxVal)
// запрос для PostgreSQL
db.Query("SELECT * FROM video WHERE likes BETWEEN $1 AND $2", minVal, maxVal)
```

*Параметр `?` будет разным для разных типов баз данных. Для MySQL, SQLite3 и MS SQL это будет `?`. Для PostgreSQL — `$N`, где N — число. Для Oracle — `:param`.*

## Работа с NULL-значениями

В SQL есть специальное значение `NULL`, которое может быть записано в поле любого типа. Оно означает, что данное поле не содержит значения. Например, `NULL` нужен, чтобы отличить пустую строку в текстовом поле от отсутствия значения (значение не определено). Чтобы отслеживать такие ситуации, пакет `database/sql` содержит типы, начинающиеся с `Null`: `sql.NullString`, `sql.NullInt64`, `sql.NullFloat64`, `sql.NullByte` и так далее. 

Для примера возьмём тип `sql.NullString`. Он представляет строку, которая в SQL-таблице может быть не установлена, то есть иметь значение `NULL`. Этот тип, как и другие подобные, описывается структурой с полем `Valid bool`. Оно равно `true`, если у поля в SQL-таблице есть значение, и `false` — если значение — `NULL`.

```go
type NullString struct {
    String string
    Valid  bool 
}
```

Предположим, нужно написать функцию, которая по указанному `video_id` ищет ролик и возвращает его описание в поле `description`. Если значение поля не определено, то нужно возвратить `-----`. В этом случае пригодится тип `NullString`. В колонке `description` таблицы `videos` нет значений `NULL`, поэтому не получится проверить работу функции. С другой стороны, тут есть записи с пустыми строками в `description`. В этом случае будет возвращена именно пустая строка, а не `-----`.

Попробуем реализовать такую функцию.
```go
func getDesc(ctx context.Context, db *sql.DB, id string) (string, error) {
    row := db.QueryRowContext(ctx,
        "SELECT description FROM videos WHERE video_id = ?", id)
    var desc sql.NullString
    
    err = row.Scan(&desc)
    if err != nil {
        return "", err
    }
    if desc.Valid {
        return desc.String, nil
    }
    return "-----", nil
}
```

Для эксперимента проверьте, что возвратит эта функция для идентификатора `0EbFotkXOiA` (должна возвратить пустую строку). После этого выполните `UPDATE videos SET description = NULL WHERE video_id = '0EbFotkXOiA';` в программе `sqlite3` или другом клиенте БД и снова вызовите `getDesc` для этого же идентификатора. Если вы всё сделали правильно, то функция возвратит `-----`.

___
Напишите запрос, который выберет из базы данных `title` (название), `views` (количество просмотров) и `channel_title` (название канала) самого популярного ролика на YouTube.

Составьте SQL-запрос самостоятельно и затем сравните его с предложенным решением.

Решение:
```sql
SELECT title, views, channel_title from videos ORDER BY views DESC LIMIT 1;

--- альтернативный вариант
SELECT title, views, channel_title from videos 
    WHERE views = (SELECT MAX(views) from videos) LIMIT 1;
```
___

## Выбор нескольких строк

Многие запросы к БД возвращают не одну, а несколько записей SQL-таблицы. Предположим, в таблице `videos` надо выбрать первые `X` самых популярных роликов, где `X` — количество.

```sql
SELECT video_id, title, views from videos ORDER BY views LIMIT X;
```

Определим структуру `Video`, в которую функция `Scan` будет записывать значения.

```go
// Video — структура видео.
type Video struct {
    Id    string
    Title string
    Views int64
}

// limit — максимальное количество записей.
const limit = 20
```

Функция для подобного запроса к базе данных может быть такой:
```go
func QueryVideos(ctx context.Context, db *sql.DB, limit int) ([]Video, error) {
    videos := make([]Video, 0, limit)

    rows, err := db.QueryContext(ctx, "SELECT video_id, title, views from videos ORDER BY views LIMIT ?", limit)
    if err != nil {
        return nil, err
    }

    // обязательно закрываем перед возвратом функции
    defer rows.Close()

    // пробегаем по всем записям
    for rows.Next() {
        var v Video
        err = rows.Scan(&v.Id, &v.Title, &v.Views)
        if err != nil {
            return nil, err
        }

        videos = append(videos, v)
    }

    // проверяем на ошибки
    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return videos, nil
}
```

Теперь сделаем следующие шаги: 
- Используем `defer rows.Close()` — это обязательно. Пока доступны результаты, соединение активно и не закроется ещё долгое время. Если не закрыть, соединение будет «грязным». За ним будет закреплён курсор с данными, и драйверу придётся закрыть соединение и создать новое, что требует ресурсов.
- Проходим по всем записям методом `rows.Next()` до тех пор, пока не пройдём все доступные результаты.
- После цикла проверяем записи на потенциальные ошибки методом `rows.Err()`. Метод возвращает ошибку, которая могла произойти в момент получения новой порции данных (новой итерации цикла `rows.Next()`). Ошибка заставит итерацию прекратиться, и алгоритм не сможет проверить её другим способом. Пример такой ошибки — разрыв сетевого соединения с сервером базы данных в процессе получения результатов запроса.
- `db.QueryContext(ctx, "SELECT video_id, title, views from videos ORDER BY views LIMIT ?", 20)` — запрос выполняется с помощью специального встроенного форматирования, с помощью символа `?`. Никогда не используйте в запросе `fmt.Sprintf`. Этот способ небезопасен — он приводит к SQL-инъекциям.

![](./assets/images/db_sprintf.png)

___
Составьте SQL запрос, который выводит 30 самых популярных каналов, отсортированных по убыванию количества просмотров. Используйте `GROUP BY` и агрегирующие функции `SUM()`, `COUNT()`, `AVG()`.
Пример результата:
Название канала | Количество просмотров всех видео | Количество видео на канале | Среднее количество просмотров видео
--- | --- | --- | ---
channel_title | 10000 | 3 | 3333.33
channel_title2 | 1000 | 2 | 500

Решение:
```sql
SELECT channel_title, SUM(views), COUNT(video_id), AVG(views)
from videos
GROUP BY channel_title
ORDER BY SUM(views)
DESC LIMIT 30;
```
___

## Расширение поддерживаемых типов

Как уже говорилось выше, `Scan` поддерживает ограниченный набор типов переменных: строки, числа, булевы значения и байты. А что, если нужно использовать массивы или собственные структуры?

Допустим, в базе данных есть теги, представленные одной строкой. В таком формате использовать теги невозможно. Было бы здорово хранить их в массиве строк. Например, в таблице `videos` есть колонка `tags`, которая содержит теги в таком формате: `Robots|"Boston Dynamics"|"SpotMini"|"Legged Locomotion"|"Dynamic robot"`. Попробуем получить их в виде слайса строк.

Для начала определим типы:
```go
type Tags []string

type Video struct {
    Id    string
    Title string
    Tags  Tags
}
```

Для поддержки тегов реализуем интерфейсы `database/driver.Valuer` и `database/sql.Scanner`:
```go
// приводит сложные типы и структуры к простому типу
type Valuer interface {
    Value() (Value, error)
}

// приводит простой тип к сложным типам и структурам Go
type Scanner interface {
    Scan(src any) error
}
```

Определим метод `Value()` для приведения `Tags` к строке:
```go
// Value — функция, реализующая интерфейс driver.Valuer
func (tags Tags) Value() (driver.Value, error) {
    // преобразуем []string в string
    if len(tags) == 0 {
        return "", nil
    }
    return strings.Join(tags, "|"), nil
}
```

Важно преобразовать значения в один из поддерживаемых типов: `int64`, `float64`, `bool`, `[]byte`, `string`, `time.Time`, `nil`.

Теперь можно записывать значения в базу. Но нужно и считывать данные — напишем функцию `Scan()` для реализации интерфейса `sql.Scanner`.

```go
func (tags *Tags) Scan(value interface{}) error {
    // если `value` равен `nil`, будет возвращён пустой массив
    if value == nil {
        *tags = Tags{}
        return nil
    }

    sv, err := driver.String.ConvertValue(value)
    if err != nil {
        return fmt.Errorf("cannot scan value. %w", err)
    }

    v, ok := sv.(string)
    if !ok {
        return errors.New("cannot scan value. cannot convert value to string")
    }
    *tags = strings.Split(v, "|")

    // удаляем кавычки у тегов
    for i, v := range *tags {
        (*tags)[i] = strings.Trim(v, `"`)
    }
    return nil
}
```

Функцией `driver.String.ConvertValue(value)` пытаемся конвертировать значение в строку. После этого преобразуем значение в слайс строк. Напишем функцию, которая будет получать идентификатор, наименование и теги у роликов.

```go
func QueryTagVideos(ctx context.Context, db *sql.DB, limit int) ([]Video, error) {
    videos := make([]Video, 0, limit)

    rows, err := db.QueryContext(ctx, "SELECT video_id, title, tags from videos "+
        "GROUP BY video_id ORDER BY views LIMIT ?", limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var v Video
        // все теги должны автоматически преобразоваться в слайс v.Tags
        err = rows.Scan(&v.Id, &v.Title, &v.Tags)
        if err != nil {
            return nil, err
        }
        videos = append(videos, v)
    }

    err = rows.Err()
    if err != nil {
        return nil, err
    }
    return videos, nil
}

func main() {
    db, err := sql.Open("sqlite", "video.db")
    if err != nil {
        panic(err)
    }
    defer db.Close()

    list, err := QueryTagVideos(context.Background(), db, 5)
    if err != nil {
        panic(err)
    }
    // для теста проверим, какие строки содержит v.Tags
    // выведем по 4 первых тега
    for _, v := range list {
        length := 4
        if len(v.Tags) < length {
            length = len(v.Tags)
        }
        fmt.Println(strings.Join(v.Tags[:length], " # "))
    }
}
```

Программа должна вывести:
```
shooting # kentucky # breaking news # u.s. news
[none]
Oregon # Ducks # college athletics # college football
news # associated press # ap # ap online
90s commercials # Huffy
```

___
В таблице `videos` есть колонка `trending_date`. Её значения хранятся в формате `год.день.месяц`. Например, `18.31.05`, `17.29.12`. Преобразуйте дату в формат `time.Time`. Напишите функцию, которая возвращает количество роликов по каждому дню `trending_date`.

```go
type Trend struct {
    T     time.Time
    Count int
}

func TrendingCount(db *sql.DB) ([]Trend, error) {
    // вставьте код
    // ...
}
```

Подсказка:
Вывести количество роликов по каждому дню можно SQL-выражением 
`SELECT trending_date, COUNT(trending_date) FROM videos GROUP BY trending_date`.

Преобразовать дату в формат `time.Time` по шаблону можно функцией 
`GoDate, err := time.Parse("06.02.01", date)`. `"06.02.01"` — шаблон год.день.месяц.

Решение:
```go
func TrendingCount(db *sql.DB) ([]Trend, error) {
    rows, err := db.Query("SELECT trending_date, COUNT(trending_date) FROM videos" +
        " GROUP BY trending_date")
    if err != nil {
        return nil, err
    }
    trends := make([]Trend, 0)
    date := new(string)
    for rows.Next() {
        trend := Trend{}
        err = rows.Scan(date, &trend.Count)
        if err != nil {
            return nil, err
        }
        var t time.Time
        if t, err = time.Parse("06.02.01", *date); err != nil {
            return nil, err
        }
        trend.T = t
        trends = append(trends, trend)
    }
    return trends, nil
}
```
___

## Запись и изменение данных в таблице

Для вставки записей в таблицу или изменения данных используется методы `ExecContext()` или `Exec()`.

```go
(db *DB) ExecContext(ctx context.Context, query string, args ...any) (Result, error)
(db *DB) Exec(query string, args ...any) (Result, error)
```

Эти методы возвращают переменную интерфейсного типа `Result`:
```go
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
```

*`LastInsertId()` возвращает целое число, сгенерированное базой данных. Обычно это происходит у атрибутов со свойством «автоинкремент» (например, `id` записи). 
`RowsAffected()` возвращает количество строк, на которые повлияли вставка, обновление или удаление.
Важно: не все базы данных и драйверы поддерживают эти функции.*

Пример вставки записей в таблицу с помощью `INSERT INTO`:
```go
func InsertNewVideo(db *sql.DB, title, channelTitle string) (int64, error) {
    result, err := db.Exec("INSERT INTO videos (title, channel_title)
    VALUES ($1, $2)", title, channelTitle)
    if err != nil {
        return 0, err
    }
    return result.LastInsertId()
}
```

Изменение записей в таблице с помощью `UPDATE`:
```go
func IncrementViews(db *sql.DB, videoId int64) error {
    result, err := db.Exec("UPDATE videos 
    SET views = views + 1
    WHERE id = $1", videoId)
    if err != nil {
        return err
    }
    affected, err := result.RowsAffected()
    if err != nil {
    return err
    }
    if affected == 0 {
    return errors.New("Video not found")
    }
    return nil
}
```

Удаление записей из таблицы с помощью `DELETE`:
```go
func DeleteVideo(db *sql.DB, videoId int64) error {
    result, err := db.Exec("DELETE FROM videos
    WHERE id = $1", videoId)
    if err != nil {
        return err
    }
    affected, err := result.RowsAffected()
    if err != nil {
    return err
    }
    if affected == 0 {
    return errors.New("Video not found")
    }
    return nil
}
```

Теперь вы можете делать запросы к БД в Go-программе и знаете, что декларация запроса пишется на стандартном языке SQL и передаётся методам `Query()` в виде строки. Вы изучили, что метод `Rows.Next()` перебирает полученные в результате запроса строки, а разбор строки в целевую структуру Go делает метод `Row.Scan()`. 

Сейчас вы знаете, что конверсия между SQL и Go для базовых типов уже реализована в пакете `database/sql`. А ещё — что пакет предоставляет интерфейс для конвертации пользовательских типов.

Также вы узнали, как отправлять `SELECT`-запросы к SQL-базам и как обрабатывать полученные в ответе данные. Этот полезный навык пригодится, например, при разработке веб-сервисов.

## Дополнительные материалы

- [Хабр | SQL-запросы быстро](https://habr.com/ru/post/480838/) — статья, посвящённая базовому синтаксису SQL.
- [Querying for data - The Go Programming Language](https://go.dev/doc/database/querying) — руководство по работе с SQL-запросами в Go с использованием пакета `database/sql`.
- [Scanners and Valuers with Go](https://husobee.github.io/golang/database/2015/06/12/scanner-valuer.html) — статья о сериализации пользовательских типов в Go (`Scanner` / `Valuer`).

# Миграции

Любой программный продукт со временем изменяется: добавляются новые функциональные возможности, а также изменяется и схема база данных.

Схема базы данных определяет, как данные организованы, хранятся и связаны между собой. Обычно она описывает таблицы и их поля, индексы, ограничения (constraints), связи между таблицами, триггеры и хранимые процедуры. 

Для управления версионированием схемы базы данных используют инструмент миграций. Он позволяет создавать, применять и отслеживать изменения в схеме. Можно сказать, что с его помощью база данных «мигрирует» от одного состояния к другому.

В этом уроке мы разберём, что такое миграции, для чего они нужны и как работают. После изучения материала вы: 
- сможете понимать структуру файлов миграций;
- научитесь создавать, применять и запускать миграции — как для создания таблиц, так и для модификации схемы.

## Преимущества миграций

Миграции можно сравнить с GIT: коммиты — это файлы миграций, история изменений — последовательность миграций и их применение, ревёрт — откат миграции. В них могут даже возникать конфликты, когда несколько разработчиков одновременно изменяют схему базы данных.

Миграции особенно полезны, когда проект разрабатывает команда, а не один программист. Как только один разработчик внёс изменения в структуру БД, эти изменения должны быть доступны и всем остальным членам команды.

Теоретически можно было бы передавать дамп базы данных или же SQL-инструкции для перевода БД из одного состояния в другое, но миграции — более удобный и функциональный инструмент. Вот почему:
- могут не только изменять БД, но и «откатывать» изменения назад, если возникли проблемы;
- создавать их проще, чем писать SQL-запросы;
- можно настроить на автоматическое применение;
- поэтапно разложены по файлам — они самостоятельно документируются и дают возможность проследить, как развивалась БД.

С помощью миграций можно создавать или удалять таблицы, столбцы и другие объекты БД, менять установленные ограничения или типы данных.

## Файлы миграций

Обычно файлы миграций хранят вместе с кодом приложения. Каждый разработчик создаёт такие файлы, когда ему необходимо изменить БД. В них хранятся инструкции для миграций — как именно нужно поменять базу данных.

Файл миграций содержит инструкции в виде классов или функций на языке программирования, либо в каком-то другом формате — например, в виде sql-файлов. В этих инструкциях написано, какие действия надо выполнить, чтобы перевести базу из одного состояния в другое — например, добавить таблицу с колонками, удалить неиспользуемую таблицу, переименовать ряд столбцов и т. д.

Как правило миграция состоит из двух секций:
- `UP` — применение изменений (что нужно сделать).
- `DOWN` — откат изменений (как вернуть обратно).

Инструкции `UP` и `DOWN` могут храниться как в одном файле, так и в двух разных — это зависит от реализации инструмента для управления миграциями. Например, библиотека `goose` предлагает объединять эти инструкции в одном файле, используя комментарии SQL. А библиотека `golang-migrate`, наоборот, предлагает разделять миграцию на два файла.

## Версионирование миграций

Версионирование миграций — это система присвоения уникальных идентификаторов каждой миграции для отслеживания порядка их применения и текущего состояния схемы базы данных.

Каждая миграция получает уникальный номер версии, который определяет:
- Порядок применения миграций
- Текущее состояние схемы БД
- Возможность отката к предыдущим версиям
- Синхронизацию между разными окружениями

Для реализации версионирования инструмент миграций в процессе работы с базой данных создает и изменяет служебную таблицу, в которой хранится информация о применённых миграциях или имя последней применённой миграции. Это позволяет инструменту определить, какие миграции он может откатить, а какие — применить.

Для именования файлов миграций можно использовать разные подходы. Рассмотрим некоторые способы присвоения уникальных идентификаторов.

### Последовательная нумерация

При последовательной нумерации к названиям файлов добавляется префикс в виде порядкового номера миграции. Например:
```
001_create_movies_table.up.sql
002_add_country_column.up.sql
003_add_country_index.up.sql
004_add_rating_column.up.sql
```

Преимущества этого способа:
- Очевидный порядок выполнения
- Краткие имена файлов
- Подходит для небольших команд и проектов

Недостатки:
- Конфликты номеров при параллельной работе
- Сложность вставки миграции между существующими

### Использование временных меток

При таком способе именования в название файла добавляется временная метка — `timestamp`. Например:
```
20060102120000_create_movies_table.up.sql
20060102120000_create_movies_table.down.sql
20060102130000_add_country_column.up.sql
20060102130000_add_country_column.down.sql
```

Преимущества этого способа:
- Практически невозможны конфликты номеров
- Видно, когда создана миграция
- Файлы естественно сортируются по времени
- Подходит для больших команд
- Разработчики могут создавать миграции независимо

Недостатки:
- Длинные имена файлов
- Временные конфликты возможны, но довольно редки
- Сложно понять последовательность без сортировки

Также в названиях миграций можно применять `uuid`-идентификаторы, названия веток или использовать комбинированные способы присвоения идентификаторов.

## Порядок работы с миграциями

Работа с миграциями сводится к следующим действиям:
- Создать файл миграций — в режиме автогенерации с использованием библиотеки или вручную.
- Применить или отменить файлы миграций.

Эти два шага выполняются по очереди, потому что применение миграций может нарушить сохранность данных. Кроме того, перед применением файла миграции важно убедиться, что в нём не написано ничего лишнего.

При командной разработке иногда могут возникать **конфликты в миграциях**. Например:
- **Конфликт имён** — два разработчика создали миграцию с одним именем
- **Конфликт версионирования** — разработчик А успел применить миграцию быстрее, чем была создана миграция разработчиком Б. Таким образом, миграция разработчика Б оказалась более старой, и версионирование сломалось.
- **Конфликты зависимостей и данных** — изменяется один и тот же ресурс схемы БД.

Подобные конфликты можно разрешать следующими способами:
- Переименовать миграцию
- Создать объединяющую миграцию
- Откатить миграцию в БД и пересоздать новую

Как правило, миграции принято отделять от запуска и деплоя приложения. Для этого их применение выносят в отдельный шаг пайплайна деплоя. Такой подход позволяет применить миграции до развёртывания новой версии приложения, заранее подготовив базу данных. В случае ошибки можно откатить миграцию и прервать процесс деплоя.

___
Что из перечисленного лучше всего описывает назначение миграций в разработке?
~~Миграции позволяют управлять изменениями схемы базы данных, включая создание, изменение и откат таблиц, столбцов и других объектов.~~
___
Распределите следующие действия по этапам:
- Создание:
    - Написать SQL-инструкции для изменения схемы БД
    - Сгенерировать файлы миграции UP и DOWN
- Выполнить команду для обновления БД до новой версии
- выполнить команду для возврата БД к предыдущей версии
___
Какие из следующих утверждений о миграциях соответствуют действительности?
~~Миграции создают таблицу в БД для хранения служебной информации.~~
___

## Практическое использование миграций

Рассмотрим пример использования инструмента миграций для создания и изменения простой базы данных фильмов. В качестве инструмента управления миграциями будет использоваться библиотека `golang-migrate`.

Сначала установим библиотеку
```bash
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```
```
go get -u github.com/golang-migrate/migrate/v4 
```

Создадим первоначальную схему базы данных с одной таблицей:
```bash
migrate create -ext sql -dir ./migrations -seq create_movies_table
```

Запуск `CLI`-команды создал два файла для миграций. Теперь добавим sql-инструкции для применения:
```sql
-- migrations/000001_create_movies_table.up.sql
-- Создание таблицы фильмов
CREATE TABLE movies (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    year INTEGER NOT NULL
);

-- Базовый индекс для поиска по названию
CREATE INDEX idx_movies_name ON movies(name);

-- Индекс для поиска по году
CREATE INDEX idx_movies_year ON movies(year);
```

А также добавим sql-инструкции для отката изменений:
```sql
-- migrations/000001_create_movies_table.down.sql
-- Откат создания таблицы фильмов
DROP INDEX IF EXISTS idx_movies_year;
DROP INDEX IF EXISTS idx_movies_name;
DROP TABLE IF EXISTS movies;
```

Применить миграции можно следующей командой:
```bash
migrate -database "postgres://user:password@localhost:5432/movies_db?sslmode=disable" -path ./migrations up
```

Через некоторое время потребовалось добавить информацию о стране производства фильма.
```bash
migrate create -ext sql -dir ./migrations -seq add_country_column
```

Аналогично предыдущей миграции, напишем инструкции для применения и отката миграций:
```sql
-- migrations/000002_add_country_column.up.sql
-- Добавляем столбец country в таблицу movies
ALTER TABLE movies ADD COLUMN country VARCHAR(100);

-- Делаем поле обязательным после заполнения данных
ALTER TABLE movies ALTER COLUMN country SET NOT NULL;

-- Создаем индекс на столбец country для оптимизации поиска по стране
CREATE INDEX idx_movies_country ON movies(country);

-- Создаем составной индекс для поиска по стране и году
CREATE INDEX idx_movies_country_year ON movies(country, year);
```

```sql
-- migrations/000002_add_country_column.down.sql
-- Откат создания индексов на столбец country
DROP INDEX IF EXISTS idx_movies_country_year;
DROP INDEX IF EXISTS idx_movies_country;
-- Откат добавления столбца country
ALTER TABLE movies DROP COLUMN IF EXISTS country;
```

Использование миграций БД позволяет упростить управление изменениями схемы базы данных, синхронизировать передачу этих изменений между разработчиками, а также внедрить версионирование. При выполнении следующих инкрементов рекомендуем использовать инструмент миграций для создания и модификации схемы.

## Дополнительные материалы

- [github.com/golang-migrate/migrate](https://github.com/golang-migrate/migrate) — библиотека `golang-migrate`.
- [github.com/pressly/goose](https://github.com/pressly/goose) — библиотека `goose`.
- [martinfowler.com/articles/evodb.html](https://martinfowler.com/articles/evodb.html) — Evolutionay Database Design, Martin Fowler

# Инкремент 11

## Задание по треку «Сервис сокращения URL»

Доработайте сервис и добавьте поддержку СУБД PostgreSQL для хранения сокращённых URL.

Сервису нужно самостоятельно создать все необходимые таблицы в базе данных. Схема и формат хранения остаются на ваше усмотрение.
Используйте инструмент миграций для создания и изменения схемы базы данных.

При отсутствии переменной окружения `DATABASE_DSN` или флага командной строки `-d` либо при их пустых значениях, вернитесь последовательно к:
- хранению сокращённых URL в файле — при наличии соответствующей переменной окружения или флага командной строки;
- хранению сокращённых URL в памяти.

## Задание по треку «Сервис сбора метрик и алертинга»

Доработайте сервис и добавьте поддержку СУБД PostgreSQL для хранения метрик.

Сервису нужно самостоятельно создать все необходимые таблицы в базе данных. Схема и формат хранения остаются на ваше усмотрение.
Используйте инструмент миграций для создания и изменения схемы базы данных.

Для хранения значений `gauge` рекомендуется использовать тип `double precision`.

При отсутствии переменной окружения `DATABASE_DSN` или флага командной строки `-d` либо при их пустых значениях, вернитесь последовательно к:
- хранению метрик в файле — при наличии соответствующей переменной окружения или флага командной строки;
- хранению метрик в памяти.

# Запись в базу данных

В предыдущем уроке вы научились извлекать данные из таблиц с помощью `SELECT`-запросов. Теперь перейдём к следующему этапу — загрузке данных в базу. На практике мы разберём, как с помощью языка Go сконвертировать `csv`-файл с видеороликами в базу данных.

Вы познакомитесь с транзакциями и подготовленными SQL-запросами — инструментами для надёжной и эффективной работы с данными. Освоив их, вы сможете создавать таблицы, а также добавлять, изменять и удалять записи — всё это необходимо для разработки полноценных приложений на Go.

## Приступаем к практике

Вы уже копировали записи из `USvideos.csv` файла в SQLite базу данных с помощью утилиты `sqlite3` и представляете, как можно импортировать данные. В реальной практике часто возникает потребность в импорте и экспорте данных. Представьте себе такие ситуации:
- заказчик хочет самостоятельно сконвертировать данные, не вникая в тонкости встроенных программ СУБД;
- требуется провести выборочный импорт/экспорт как по колонкам, так и по записям;
- нужен импорт данных из нескольких источников;
- необходимо сделать дополнительные преобразования данных.

Во всех этих случаях не обойтись без написания нужной программы, пакета или функции.

Предположим, что вас попросили написать программу, которая будет импортировать в SQLite таблицу только пять определённых полей из уже знакомого вам файла с видеороликами. Начать стоит с предварительного этапа — чтения `csv`-файла. Так как количество записей небольшое, можно написать функцию, которая читает данные и возвращает слайс структур со всеми записями. После — добавим данные в таблицу.

Опишем структуру для хранения данных:
```go
type Video struct {
    Id          string    // video_id
    Title       string    // title
    PublishTime time.Time // publish_time
    Tags        []string  // tags
    Views       int       // views
}
```

Для чтения `csv`-файла будем использовать пакет `encoding/csv`. Приведём готовую функцию, которая будет использоваться в программе конвертации. Если вы чувствуете, что способны написать функцию `func readVideoCSV(csvFile string) ([]Video, error)` самостоятельно, то попробуйте. Возможно, вы реализуете её элегантнее предложенного варианта. 

```go
func readVideoCSV(csvFile string) ([]Video, error) {
    // открываем csv-файл
    file, err := os.Open(csvFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()
    var videos []Video

    // определим индексы нужных полей
    const (
        Id          = 0 // video_id
        Title       = 2 // title
        PublishTime = 5 // publish_time
        Tags        = 6 // tags
        Views       = 7 // views
    )

    // конструируем Reader из пакета encoding/csv
    // он умеет читать строки csv-файла
    r := csv.NewReader(file)
    // пропустим первую строку с именами полей
    if _, err := r.Read(); err != nil {
        return nil, err
    }

    for {
        // csv.Reader за одну операцию Read() считывает одну csv-запись
        l, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return nil, err
        }
        // инициализируем целевую структуру,
        // в которую будем делать разбор csv-записи
        v := Video{
            Id:    l[Id],
            Title: l[Title],
        }
        // парсинг строковых записей в типизированные поля структуры
        if v.PublishTime, err = time.Parse(time.RFC3339, l[PublishTime]); err != nil {
            return nil, err
        }
        tags := strings.Split(l[Tags], "|")
        for i, v := range tags {
            tags[i] = strings.Trim(v, `"`)
        }
        v.Tags = tags
        if v.Views, err = strconv.Atoi(l[Views]); err != nil {
            return nil, err
        }
        // добавляем полученную структуру в слайс
        videos = append(videos, v)
    }
    return videos, nil
}
```

Теперь можно создать макет программы импорта.

```go
func insertVideos(ctx context.Context, db *sql.DB, videos []Video) error {
    // здесь будем вставлять записи в базу данных
    // ...
    return nil
}

func main() {
    // открываем соединение с БД
    db, err := sql.Open("sqlite", "newvideo.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // читаем записи из файла в слайс []Video вспомогательной функцией
    videos, err := readVideoCSV("USvideos.csv")
    if err != nil {
        log.Fatal(err)
    }
    // записываем []Video в базу данных
    // тоже вспомогательной функцией
    err = insertVideos(context.Background(), db, videos)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Всего csv-записей %v\n", len(videos))
}
```

Запустите код и посмотрите результат:
```bash
$ go run main.go
Всего csv-записей 40949
```

В прошлом уроке мы рассматривали два типа методов: `Query...()` и `Exec...()`. Первые методы используются для отправки `SELECT`-запросов, а вторые — для остальных SQL-команд. Используем метод `ExecContext()` для вставки элементов в таблицу с помощью `INSERT INTO`.

```go
(db *DB) ExecContext(ctx context.Context, query string, args ...any) (Result, error)
```

Этот метод возвращает переменную интерфейсного типа `Result`:
```go
type Result interface {
    LastInsertId() (int64, error)
    RowsAffected() (int64, error)
}
```

*`LastInsertId()` возвращает целое число, сгенерированное базой данных. Обычно это происходит у атрибутов со свойством «автоинкремент» (например, `id` записи). 
`RowsAffected()` возвращает количество строк, на которые повлияли вставка, обновление или удаление.
Важно: не все базы данных и драйверы поддерживают эти функции.*

Добавление всех элементов из слайса `[]Video` можно реализовать таким образом:
```go
for _, v := range videos {
    // в этом случае возвращаемое значение не несёт полезной информации,
    // поэтому его можно игнорировать
    _, err := db.ExecContext(ctx,
        "INSERT INTO videos (video_id, title, publish_time, tags, views)"+
            " VALUES(?,?,?,?,?)", v.Id, v.Title, v.PublishTime, 
            strings.Join(v.Tags, `|`), v.Views)
    if err != nil {
        return err
    }
}
```

___
Реализуйте функцию `insertVideos()` с указанным выше кодом и добавьте в `main()` создание таблицы `videos`. В таблице достаточно определить только импортируемые поля.

Решение:
```go
func insertVideos(ctx context.Context, db *sql.DB, videos []Video) error {
    for _, v := range videos {
        _, err := db.ExecContext(ctx,
            "INSERT INTO videos (video_id, title, publish_time, tags, views)"+
                " VALUES(?,?,?,?,?)", v.Id, v.Title, v.PublishTime, 
                strings.Join(v.Tags, `|`), v.Views)
        if err != nil {
            return err
        }
    }
    return nil
}

func main() {
    db, err := sql.Open("sqlite", "newvideo.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()
    _, err = db.ExecContext(ctx, `CREATE TABLE videos (
        "video_id" TEXT,
        "title" TEXT,
        "publish_time" TEXT,
        "tags" TEXT,
        "views" INTEGER
      )`)
    if err != nil {
        log.Fatal(err)
    }

    videos, err := readVideoCSV("USvideos.csv")
    if err != nil {
        log.Fatal(err)
    }

    // записываем []Video в базу данных
    err = insertVideos(ctx, db, videos)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Всего csv-записей %v\n", len(videos))
}
```
___

Если вы запустите эту программу, то увидите, что она выполняется довольно медленно. На среднем компьютере процесс добавления записей может занимать больше минуты. Стоит отметить, что драйвер `github.com/mattn/go-sqlite3` будет выполнять эту работу на 20–25% быстрее, чем `modernc.org/sqlite`.

Оценить время, затраченное на вставку записей в таблицу, можно так:
```go
start := time.Now()
err = insertVideos(context.Background(), db, videos)
if err != nil {
    log.Fatal(err)
}
fmt.Println(time.Since(start))
```

*Чтобы ускорить вставку записей в базу данных, используйте транзакции — они позволяют объединить множество операций в одну и значительно повысить производительность.*

**Транзакция** — это набор последовательных запросов к БД, который представляет собой одну атомарную операцию с данными. Транзакция может быть выполнена либо полностью, либо не выполнена вообще — тогда она не внесёт изменений в базу данных. 

Рассмотрим пример. Допустим, существует таблица с банковскими счетами клиентов. Когда один клиент переводит деньги другому, необходимо выполнить две операции: списать сумму со счёта плательщика, зачислить сумму на счёт получателя. Если производить операции отдельно, то после списания денег может произойти системный сбой и нужная сумма не зачислится на счёт получателя. Использование транзакций гарантирует, что выполнятся две операции или, в случае ошибки, платёж не пройдёт. В этом и есть суть транзакции: она нужна для сохранения целостности данных.

Транзакциями можно управлять с помощью трёх основных команд: `Begin`, `Commit` и `Rollback`. Рассмотрим каждую из них: 
- `Begin` указывает на начало транзакции.
- `Commit` сохраняет изменения, которые внесли в транзакцию.
- `Rollback` откатывает изменения. Например, `Rollback` поможет вернуть данные пользователя, если их удаление вызывало системную ошибку в работе сервера.

Для начала транзакции нужно вызвать один из методов, которые возвращают переменную типа `*sql.Tx`:
```go
func (db *DB) Begin() (*Tx, error)
func (db *DB) BeginTx(ctx context.Context, opts *TxOptions) (*Tx, error)
```

В методе `BeginTx` передаётся контекст — транзакция будет отменена при отмене контекста. Параметр `opts` позволяет настраивать транзакцию — нам это не требуется, так что будем передавать `nil` в этом параметре.

Метод `(tx *Tx) Commit() error` сохраняет все изменения, сделанные в рамках транзакции, а `(tx *Tx) Rollback() error` отменяет все изменения и откатывает транзакцию.

Для выполнения запросов у типа `*sql.Tx` есть методы `tx.Query...()` и `tx.Exec...()`, аналогичные тем, которые мы использовали.

```go
func insertVideos(ctx context.Context, db *sql.DB, videos []Video) error {
    // начинаем транзакцию
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    for _, v := range videos {
        // все изменения записываются в транзакцию
        _, err := tx.ExecContext(ctx, 
            "INSERT INTO videos (video_id, title, publish_time, tags, views)"+
            " VALUES(?,?,?,?,?)", v.Id, v.Title, v.PublishTime,
            strings.Join(v.Tags, `|`), v.Views)
        if err != nil {
            // если ошибка, то откатываем изменения
            tx.Rollback()
            return err
        }
    }
    // завершаем транзакцию
    return tx.Commit()
}
```

Попробуйте удалить файл базы и запустить этот вариант программы. Процесс создания `newvideo.db` должен быть быстрее в сотню раз.

![](./assets/images/transactions.png)

## Подготовленные SQL-запросы

Каждый раз при отправке запроса СУБД приходится разбирать полученное SQL-выражение. Если выражения сложные и они регулярно выполняются, можно использовать подготовленные — скомпилированные — SQL-запросы. В некоторых ситуациях они могут заметно увеличить скорость работы с БД. Сервер базы данных может запомнить запрос во внутреннем представлении и не тратить каждый раз время на его разбор из текстовой формы. Это существенно ощущается при частом выполнении одних и тех же запросов (возможно, с разными параметрами). 

Для создания скомпилированного запроса используется методы `PrepareContext(...)` или `Prepare(...)`.

```go
func (db *DB) PrepareContext(ctx context.Context, query string) (*Stmt, error)
func (db *DB) Prepare(query string) (*Stmt, error)
```

Они возвращают переменную типа `*sql.Stmt`. У этого типа, как и у `*sql.Tx`, есть методы, аналогичные `Query...()` и `Exec...()`. Например, ниже показан код с `QueryRowContext` из прошлого урока и подобный код с использованием подготовленной SQL-инструкции:
```go
row := db.QueryRowContext(ctx, 
            "SELECT description FROM videos WHERE video_id = ?", id)

// пример с использованием PrepareContext
stmt, err := db.PrepareContext(ctx, "SELECT description FROM videos WHERE video_id = ?")
// не забывайте закрывать запрос
defer stmt.Close()
// сам запрос не указывается, передаются только параметры
row := stmt.QueryRowContext(ctx, id)
```

Тип `*sql.Tx` также имеет методы `PrepareContext(...)` или `Prepare(...)`, которые создают скомпилированные запросы для использования в рамках конкретной транзакции.

___
Создайте в функции `insertVideos()` подготовленный запрос и используйте его для вставки видеороликов. Можно использовать версию как с транзакциями, так и без.

Решение:
```go
func insertVideos(ctx context.Context, db *sql.DB, videos []Video) error {
    tx, err := db.Begin()
    if err != nil {
        return err
    }
    // можно вызвать Rollback в defer,
    // если Commit будет раньше, то откат проигнорируется
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO videos (video_id, title, publish_time, tags, views)"+
        " VALUES(?,?,?,?,?)")
    if err != nil {
        return err
    }
    defer stmt.Close()

    for _, v := range videos {
        _, err := stmt.ExecContext(ctx, v.Id, v.Title, v.PublishTime,
            strings.Join(v.Tags, `|`), v.Views)
        if err != nil {
            return err
        }
    }
    return tx.Commit()
}
```
___


## Множественная вставка

Перенос данных из файла в базу требует больших ресурсов памяти. Сейчас в `csv`-файле — 40 000 записей. А что, если их будет миллион? И размер файла будет 10, 200, 500 гигабайт? Тогда нужно будет записывать строки в базу данных **пачками**. Как это сделать:
- читать `csv`-файл построчно,
- добавлять строки во временный буфер,
- при достижении лимита в 1000 строк делать запись в базу данных,
- повторять, пока не запишем все значения из `csv`-файла.

*Важно: вы можете возразить, что так лучше не делать. Множественную вставку для `INSERT` можно реализовать с помощью SQL. Например, `INSERT (id, name) INTO table VALUES (1, "video1"), (2, "video2");`. А если надо вставлять по 1000 записей? Удобно ли будет перечислять их в одном огромном запросе? Чтобы научиться использовать транзакции, рассмотрим вариант с буферизацией. Хотя для этого примера он избыточен.*

Попробуем реализовать буферизированное добавление данных с минимальными изменениями текущего кода. Функцию `insertVideos()` можно оставить без изменений, она будет получать очередной слайс по мере заполнения буфера. Изменения придётся внести в функцию чтения записей `readVideoCSV()`. 

___
Переработайте функцию `readVideoCSV()` для добавления данных порциями по 1000 записей.
```go
func readVideoCSV(ctx context.Context, db *sql.DB, csvFile string) error {
    // внесите изменения в существующий код,
    // чтобы ролики добавлялись порциями по 1000 записей
    // ...
}

func main() {
    // открываем соединение с БД и создаём таблицу
    // эта часть кода остаётся без изменений
    // ...
    err = readVideoCSV(ctx, db, "USvideos.csv")
    if err != nil {
        log.Fatal(err)
    }
}
```

Решение:
```go
func readVideoCSV(ctx context.Context, db *sql.DB, csvFile string) error {
    // открываем csv файл
    file, err := os.Open(csvFile)
    if err != nil {
        return err
    }
    defer file.Close()
    videos := make([]Video, 0, 1000)

    // определим индексы нужных полей
    const (
        Id          = 0 // video_id
        Title       = 2 // title
        PublishTime = 5 // publish_time
        Tags        = 6 // tags
        Views       = 7 // views
    )

    // конструируем Reader из пакета encoding/csv
    // он умеет читать строки csv-файла
    r := csv.NewReader(file)
    // пропустим первую строку с именами полей
    if _, err := r.Read(); err != nil {
        return err
    }

    for {
        // csv.Reader за одну операцию Read() считывает одну csv-запись
        // в виде []string
        l, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return err
        }
        // инициализируем целевую структуру,
        // в которую будем делать разбор csv-записи
        v := Video{
            Id:    l[Id],
            Title: l[Title],
        }
        // парсинг строковых записей в типизированные поля структуры
        if v.PublishTime, err = time.Parse(time.RFC3339, l[PublishTime]); err != nil {
            return err
        }
        tags := strings.Split(l[Tags], "|")
        for i, v := range tags {
            tags[i] = strings.Trim(v, `"`)
        }
        v.Tags = tags
        if v.Views, err = strconv.Atoi(l[Views]); err != nil {
            return err
        }
        // добавляем полученную структуру в слайс
        videos = append(videos, v)
        if len(videos) == 1000 {
            if err = insertVideos(ctx, db, videos); err != nil {
                return err
            }
            videos = videos[:0]
        }
    }
    // добавляем оставшиеся записи
    return insertVideos(ctx, db, videos)
}
```

## Индексы

Предположим, запросы пользователей требуют поиска записей по определённому полю. Построим простейший стенд для оценки скорости работы поиска по колонке `video_id`. В примере для краткости не используются скомпилированные запросы, но вы можете их добавить.

```go
// query_test.go
package main

import (
    "database/sql"
    "log"
    "testing"

    _ "modernc.org/sqlite"
)

func BenchmarkQuery(b *testing.B) {
    db, err := sql.Open("sqlite", "newvideo.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        db.QueryRow("SELECT title, views FROM videos WHERE video_id = ?", "R39-E3uG5J0")
    }
}
```

Если запустить `go test -bench=. query_test.go`, то видно, что скорость оставляет желать лучшего, так как движку базы данных приходится перебирать все записи в таблице `videos`. И чем больше количество записей, тем медленнее будет проходить выборка данных. 

```
cpu: Intel(R) Core(TM) i5-6260U CPU @ 1.80GHz
BenchmarkQuery-4             216           5551661 ns/op
```

Если предварительно отсортировать записи по `video_id`, то поиск будет значительно быстрее. Но всего должно быть в меру. Не нужно увлекаться индексами: каждый индекс занимает дополнительное место на диске или в памяти, а ещё замедляет операции по изменению таблицы. Например, при добавлении записи будут обновлены все индексы таблицы.

Для создания индекса используется команда `CREATE INDEX name_idx ON table (field1, field2, ...)`, где `name_idx` — имя индекса, `table` — имя таблицы, а `fieldN` — имена полей, по которым создаётся индекс.

```sql
-- создание индекса для таблицы videos по полю video_id
CREATE INDEX video_id ON videos (video_id)
```

Добавим в конец программы добавления роликов создание индекса по колонке `video_id` и заново сконвертируем данные.

```go
_, err = db.ExecContext(ctx, "CREATE INDEX video_id ON videos (video_id)")
if err != nil {
    log.Fatal(err)
}
```

После этого можно снова запустить бенчмарк и сравнить скорость работы с предыдущим запуском. Скорость выполнения запроса должна увеличиться примерно на порядок.

```
cpu: Intel(R) Core(TM) i5-6260U CPU @ 1.80GHz
BenchmarkQuery-4            3272            329454 ns/op
```

Обычно индексы планируются на этапе разработки архитектуры программы вместе со схемой базы, чтобы обеспечить высокую скорость работы программы.

## Изменение записей

Рассмотрим ещё одну SQL-команду — `UPDATE`. Она изменяет значения указанных в ней полей (колонок). С точки зрения Go-разработчика, её использование ничем не отличается от `INSERT`, `CREATE` и прочих команд, которые вызываются методами `Exec...()`. Нужно обращать внимание на корректность самого SQL-выражения. Если, например, выполнить такой запрос: `db.Exec("UPDATE videos SET title = ?", "Marshmello - Blocks")`, новое название будет установлено для всех роликов и таблица потеряет свою ценность.

Чтобы немного попрактиковаться, представим такую ситуацию. Вы создали из `csv`-файла базу данных `newvideo.db` и отдали её руководителю. Через некоторое время было решено избавиться от всех тегов, содержащих слово `best`. Эта задача попадает снова к вам. 

Если бы требовалось просто очистить поле тегов у всех роликов, где встречается слово best, это можно было бы сделать одним запросом типа `UPDATE videos SET tags='' WHERE tags LIKE '%best%'`. Но перед вами задача — убрать конкретные теги. Для этого нужно получить все ролики, содержащие нужный тег, а затем удалить его из каждой записи.

Вы, наверно, заметили, что значения колонки `video_id` не уникальны, но при этом у записей с одинаковым `video_id` нельзя гарантировать совпадание тегов. Чтобы не проверять каждую запись вручную, воспользуемся `GROUP BY`. Сгруппируем записи с одинаковым значением указанного поля. `GROUP BY tags` позволит нам получить только набор различающихся значений `tags`, без дублирующихся записей.

```go
package main

import (
    "context"
    "database/sql"
    "fmt"
    "log"
    "strings"

    _ "modernc.org/sqlite"
)

type TagVideo struct {
    Tags string
}

func getList(ctx context.Context, db *sql.DB) (videos []TagVideo, err error) {
    /* 1. Выберите список роликов, у которых в тегах есть слово best
       2. Получите теги в структуру TagVideo
       3. Добавьте TagVideo в videos
       4. Повторите пункты 2-3 для каждого найденного элемента
    */
}

func main() {
    // открываем соединение с БД
    db, err := sql.Open("sqlite", "newvideo.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()

    videos, err := getList(ctx, db)
    if err != nil {
        log.Fatal(err)
    }

    var updates int64 // подсчитаем количество изменённых записей
    for _, v := range videos {
        var tags []string
        // удаляем лишние теги
        for _, tag := range strings.Split(v.Tags, `|`) {
            if !strings.Contains(strings.ToLower(tag), "best") {
                tags = append(tags, tag)
            }
        }
        res, err := db.ExecContext(ctx,
            "UPDATE videos SET tags = ? WHERE tags = ?",
            strings.Join(tags, "|"), v.Tags)
        if err != nil {
            log.Fatal(err)
        }
        // посмотрим, сколько записей было обновлено
        if upd, err := res.RowsAffected(); err == nil {
            updates += upd
        }
    }
    fmt.Println(updates)
}
```

___
Реализуйте функцию `getList()` в приведённом выше примере. Если всё сделано правильно, то первом запуске программа выведет `2576`, а при повторном - `0`.

Решение:
```go
func getList(ctx context.Context, db *sql.DB) (videos []TagVideo, err error) {
    var rows *sql.Rows

    rows, err = db.QueryContext(ctx,
        "SELECT tags FROM videos"+
            " WHERE tags LIKE '%best%' GROUP BY tags")
    if err != nil {
        return
    }
    defer rows.Close()

    for rows.Next() {
        var video TagVideo

        if err = rows.Scan(&video.Tags); err != nil {
            return
        }
        videos = append(videos, video)
    }
    return
}
```

## Использование именованных параметров в запросе

Представьте, что нужно выполнить `DELETE`-запрос: удалить опубликованные в определённый день видео, у которых меньше 1000 просмотров.

Пример такой операции на языке SQL:
```sql
DELETE FROM videos 
    WHERE publish_time BETWEEN '2018-01-01 00:00:00' AND '2018-01-31 23:59:59'
        AND views < 1000
```

На Go этот запрос можно написать так:
```go
_, err := db.ExecContext(сtx, "DELETE FROM videos WHERE publish_time" +
                " BETWEEN ? AND ? AND views < ?", fromDate,
                toDate, viewsLimit) 
```

Запрос с `?` плохо читается, особенно если в нём больше трёх аргументов. Чтобы запрос читался легче, можно использовать именованные аргументы, которые определяются функцией `Named()`:
```go
func Named(name string, value interface{}) NamedArg
```

Тогда получится переделать запрос:
```go
_, err := db.ExecContext(сtx, "DELETE FROM videos WHERE publish_time" +
                " BETWEEN @start AND @end AND views < @views",
                sql.Named("end", toDate),
                sql.Named("start", fromDate),
                sql.Named("views", viewsLimit),
)
```

Такой подход позволит совершать меньше ошибок при выполнении запросов.

## Настройка пула соединений

Пакет `database/sql` поддерживает пул соединений. Часто у СУБД есть ограничение на количество подключений, которые она может обслуживать. Чтобы поддерживать баланс между скоростью обработки запросов и количеством одновременно работающих клиентов, пакет `database/sql` позволяет настраивать параметры пула подключений. Это можно сделать с помощью методов:
```go
func (db *DB) SetMaxOpenConns(n int)
func (db *DB) SetMaxIdleConns(n int)
func (db *DB) SetConnMaxIdleTime(d time.Duration)
func (db *DB) SetConnMaxLifetime(d time.Duration)
```

Коротко опишем каждый из них: 
- `db.SetMaxOpenConns(n int)` служит для ограничения количества соединений, используемых приложением. Рекомендуемого числа нет: оно зависит от задач и особенностей работы. Можете поиграться с настройками.
- `db.SetMaxIdleConns(n int)` служит для ограничения количества простаивающих соединений. Устанавливайте его с тем же значением, что и `db.SetMaxOpenConns(...)`. Если значение меньше, чем `SetMaxOpenConns()`, соединения могут открываться и закрываться гораздо чаще.
- `db.SetConnMaxIdleTime(d time.Duration)` задаёт максимальное время простаивания одного соединения. Поскольку некоторые промежуточные инфраструктурные приложения закрывают простаивающие соединения после пяти минут, рекомендуется устанавливать значение меньше.
- `db.SetConnMaxLifetime(d time.Duration)` задаёт максимальное время работы одного соединения с момента его создания.

## SQLx

Преимущество пакета `database/sql` — простота. Но есть у него и свои неудобства. В рабочих проектах, в больших таблицах использование функции `Scan()` может быть утомительным. Нужно следить за порядком столбцов, обязательно делать итерацию по результатам.

Пакет [SQLx](https://github.com/jmoiron/sqlx) помогает облегчить жизнь разработчика на Go. Он содержит множество методов и функций. А ещё при работе с ним можно применять методы пакета `database/sql`. Вернёмся к примеру и удалим теги, содержащие слово `worst`. Используем `sqlx`.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "strings"

    "github.com/jmoiron/sqlx"
    _ "modernc.org/sqlite"
)

type TagVideo struct {
    Tags string
}

func getList(ctx context.Context, db *sqlx.DB) (videos []TagVideo, err error) {
    sqlSelect := `SELECT tags FROM videos 
                    WHERE tags LIKE '%worst%' GROUP BY tags`
    err = db.SelectContext(ctx, &videos, sqlSelect)
    return
}

func main() {
    db := sqlx.MustOpen("sqlite", "newvideo.db")
    defer db.Close()

    ctx := context.Background()
    videos, err := getList(ctx, db)
    if err != nil {
        log.Fatal(err)
    }

    var updates int64
    sqlUpdate := "UPDATE videos SET tags = ? WHERE tags = ?"

    for _, v := range videos {
        var tags []string
        // удаляем лишние теги
        for _, tag := range strings.Split(v.Tags, `|`) {
            if !strings.Contains(strings.ToLower(tag), "worst") {
                tags = append(tags, tag)
            }
        }

        res := db.MustExecContext(ctx, sqlUpdate, strings.Join(tags, `|`), v.Tags)
        if upd, err := res.RowsAffected(); err == nil {
            updates += upd
        }
    }
    fmt.Println(updates)
}
```

Здесь используются только три метода: `MustOpen`, `SelectContext` и `MustExecContext`, но при этом количество строк кода уменьшилось на 20%. Уже этот пример показывает, насколько эффективным может быть применение пакета `sqlx`. 

## Почему в Go не используют ORM

**ORM (от англ. Object-Relational Mapping — «объектно-реляционное отображение»)** — это технология, которая связывает таблицы в базах данных и объекты (классы, структуры) в языке программирования. Основное преимущество ORM — абстракция от разных типов базы данных.

В Go не принято использовать ORM. И причин несколько. 
- Каждый раз при использовании ORM нужно изучать его синтаксис. Если у вас есть специфичный запрос, гораздо проще найти решение на чистом SQL, чем для ORM. Использование «костылей» при сложных запросах (или запросах на чистом SQL) сильно замедляет работу.
- Разные компании используют и будут использовать разные ORM — переключение из одного проекта в другой будет требовать от вас очередного погружения в детали работы ORM. Выгоднее использовать чистый SQL.

Но есть и ситуации-исключения, когда ORM всё же можно использовать: 
- Когда в проекте множество однотипных и простых CRUD-операций (`Create`, `Read`, `Update` и `Delete`). В этом случае ORM сэкономит много времени на разработку.
- В собственных небольших проектах, где важно быстро разработать решение «на коленке», ORM могут помочь сконцентрироваться на идее, а не на коде.
- Когда требуется поддержка различных СУБД.

Один из самых популярных ORM-пакетов в Go — [GORM](https://gorm.io/docs/). Миграции, CRUD-операции, хуки, работа со связями — всё это есть в GORM.

___
Метод `Exec()` определён для типов `sql.DB`, `sql.Tx`, `sql.Stmt`. Чем они отличаются?

Верные:
- `Stmt.Exec()` не требует SQL-декларации в качестве аргумента — только параметры запроса. (Декларация `Stmt` уже подготовлена и скомпилирована. Осталось только передать параметры `stmt.Exec()`.)
- Запрос `Tx.Exec()` актуален только в рамках транзакции `Tx`. (Транзакция может быть отменена методом `Rollback()`, и результаты `tx.Exec()` потеряют значимость.)
- `Stmt.Exec()` работает быстрее, чем `DB.Exec()`. (Да. Запрос уже подготовлен и скомпилирован. Осталось только установить параметры.)

Неверные:
- `DB.Exec()` даёт больше гарантий успешного выполнения, потому что это метод самой БД.
(`Stmt.Exec()` и `Tx.Exec()` работают с теми же данными и имеют такой же приоритет.)
- При конкурентном обращении к БД `Stmt.Exec()` получит преимущество перед `DB.Exec()`. (Для движка или драйвера разницы в приоритете обращений нет.)
___

## Дополнительные материалы

- [Хабр | Памятка/шпаргалка по SQL](https://habr.com/ru/post/564390/) — удобная шпаргалка по SQL.
- [GitHub | Пакет SQLx](https://github.com/jmoiron/sqlx) — репозиторий Go-библиотеки SQLx.
- [GORM](https://gorm.io/docs/) — официальная документация по GORM.

# Обучение Алисы 11

Теперь пришло время подключить навык Алисы к реальной СУБД. Первым делом напишем реализацию интерфейса `store.MessageStore` с запросами к СУБД PostgreSQL.

Создадим директорию `internal/store/pg` и поместим туда файл `store.go`:
```
 > ~/dev/alice-skill
        |
        |--- cmd
        |     |--- skill
        |            |--- flags.go
        |            |--- gzip.go
        |            |--- main.go
        |            |--- main_test.go
        |--- internal
        |       |--- logger
        |       |      |--- logger.go
        |       |--- models
        |       |      |--- models.go
        |       |--- store
        |              |--- mock
        |              |     |--- store.go
        |              |--- pg
        |              |     |--- store.go
        |              |--- store.go
        |--- go.mod
        |--- go.sum
```

Добавим в новый файл следующий код:
```go
package pg

import (
    "context"
    "database/sql"
    "time"

    "github.com/bluegopher/alice-skill/internal/store"
)

// Store реализует интерфейс store.MessageStore и позволяет взаимодействовать с СУБД PostgreSQL
type Store struct {
    // Поле conn содержит объект соединения с СУБД
    conn *sql.DB
}

// NewStore возвращает новый экземпляр PostgreSQL-хранилища
func NewStore(conn *sql.DB) *Store {
    return &Store{conn: conn}
}

// Bootstrap подготавливает БД к работе, создавая необходимые таблицы и индексы
func (s *Store) Bootstrap(ctx context.Context) error {
    // запускаем транзакцию
    tx, err := s.conn.BeginTx(ctx, nil)
    if err != nil {
        return err
    }

    // в случае неуспешного коммита все изменения транзакции будут отменены
    defer tx.Rollback()

    // создаём таблицу пользователей и необходимые индексы
    tx.ExecContext(ctx, `
        CREATE TABLE users (
            id varchar(128) PRIMARY KEY,
            username varchar(128)
        )
    `)
    tx.ExecContext(ctx, `CREATE UNIQUE INDEX sender_idx ON users (username)`)

    // создаём таблицу сообщений и необходимые индексы
    tx.ExecContext(ctx, `
        CREATE TABLE messages (
            id serial PRIMARY KEY,
            sender varchar(128),
            recipient varchar(128),
            payload text,
            sent_at timestamp with time zone,
            read_at timestamp with time zone DEFAULT NULL
        )
    `)
    tx.ExecContext(ctx, `CREATE INDEX recipient_idx ON messages (recipient)`)

    // коммитим транзакцию
    return tx.Commit()
}

func (s *Store) FindRecipient(ctx context.Context, username string) (userID string, err error) {
    // запрашиваем внутренний идентификатор пользователя по его имени
    row := s.conn.QueryRowContext(ctx, `SELECT id FROM users WHERE username = $1`, username)
    err = row.Scan(&userID)
    return
}

func (s *Store) ListMessages(ctx context.Context, userID string) ([]store.Message, error) {
    // запрашиваем данные обо всех сообщениях пользователя, без самого текста
    rows, err := s.conn.QueryContext(ctx, `
        SELECT
            m.id,
            u.username AS sender,
            m.sent_at
        FROM messages m
        JOIN users u ON m.sender = u.id
        WHERE
            m.recipient = $1
    `, userID)

    if err != nil {
        return nil, err
    }
    // не забываем закрыть курсор после завершения работы с данными
    defer rows.Close()

    // считываем записи в слайс сообщений
    var messages []store.Message
    for rows.Next() {
        var m store.Message
        if err := rows.Scan(&m.ID, &m.Sender, &m.Time); err != nil {
            return nil, err
        }
        messages = append(messages, m)
    }

    // необходимо проверить ошибки уровня курсора
    if err := rows.Err(); err != nil {
        return nil, err
    }

    return messages, nil
}

func (s *Store) GetMessage(ctx context.Context, id int64) (*store.Message, error) {
    // запрашиваем сообщение по внутреннему идентификатору
    row := s.conn.QueryRowContext(ctx, `
        SELECT
            m.id,
            u.username AS sender,
            m.payload,
            m.sent_at
        FROM messages m
        JOIN users u ON m.sender = u.id
        WHERE
            m.id = $1
    `,
        id,
    )

    // считываем значения из записи БД в соответствующие поля структуры
    var msg store.Message
    err := row.Scan(&msg.ID, &msg.Sender, &msg.Payload, &msg.Time)
    if err != nil {
        return nil, err
    }
    return &msg, nil
}

func (s *Store) SaveMessage(ctx context.Context, userID string, msg store.Message) error {
    // добавляем новое сообщение в БД
    _, err := s.conn.ExecContext(ctx, `
        INSERT INTO messages
        (sender, recipient, payload, sent_at)
        VALUES
        ($1, $2, $3, $4);
    `, msg.Sender, userID, msg.Payload, time.Now())

    return err
}
```

Теперь нам необходимо добавить новый аргумент командной строки, через который мы будем получать строку с параметрами соединения с СУБД.
Обновим файл `cmd/skill/flags.go`:
```go
package main

import (
    "flag"
    "os"
)

var (
    flagRunAddr  string
    flagLogLevel string
    // переменная будет содержать параметры соединения с СУБД
    flagDatabaseURI string
)

func parseFlags() {
    flag.StringVar(&flagRunAddr, "a", ":8080", "address and port to run server")
    flag.StringVar(&flagLogLevel, "l", "info", "log level")
    // обрабатываем аргумент -d
    flag.StringVar(&flagDatabaseURI, "d", "", "database URI")
    flag.Parse()

    if envRunAddr := os.Getenv("RUN_ADDR"); envRunAddr != "" {
        flagRunAddr = envRunAddr
    }
    if envLogLevel := os.Getenv("LOG_LEVEL"); envLogLevel != "" {
        flagLogLevel = envLogLevel
    }
    // обрабатываем переменную окружения DATABASE_URI
    if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
        flagDatabaseURI = envDatabaseURI
    }
}
```

После этого мы можем использовать переменную `flagDatabaseURI` для подключения к СУБД PostgreSQL с помощью драйвера `github.com/jackc/pgx/v5`.
Нам необходимо добавить код в файл `cmd/skill/main.go`:
```go
package main

import (
    "database/sql"
    "net/http"
    "strings"

    "github.com/bluegopher/alice-skill/internal/logger"
    "github.com/bluegopher/alice-skill/internal/store/pg"
    _ "github.com/jackc/pgx/v5/stdlib"
    "go.uber.org/zap"
)

func main() {
    parseFlags()

    if err := run(); err != nil {
        panic(err)
    }
}

func run() error {
    if err := logger.Initialize(flagLogLevel); err != nil {
        return err
    }

    // создаём соединение с СУБД PostgreSQL с помощью аргумента командной строки
    conn, err := sql.Open("pgx", flagDatabaseURI)
    if err != nil {
        return err
    }

    // создаём экземпляр приложения, передавая реализацию хранилища pg в качестве внешней зависимости
    appInstance := newApp(pg.NewStore(conn))

    logger.Log.Info("Running server", zap.String("address", flagRunAddr))
    // обернём хендлер webhook в middleware с логированием и поддержкой gzip
    return http.ListenAndServe(flagRunAddr, logger.RequestLogger(gzipMiddleware(appInstance.webhook)))
}

...
```

Теперь навык Алисы может работать с реальной СУБД. Осталось только добавить в хендлер немного бизнес-логики.
Исправим файл `cmd/skill/app.go`:
```go
func (a *app) webhook(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()

    if r.Method != http.MethodPost {
        logger.Log.Debug("got request with bad method", zap.String("method", r.Method))
        w.WriteHeader(http.StatusMethodNotAllowed)
        return
    }

    logger.Log.Debug("decoding request")
    var req models.Request
    dec := json.NewDecoder(r.Body)
    if err := dec.Decode(&req); err != nil {
        logger.Log.Debug("cannot decode request JSON body", zap.Error(err))
        w.WriteHeader(http.StatusInternalServerError)
        return
    }

    if req.Request.Type != models.TypeSimpleUtterance {
        logger.Log.Debug("unsupported request type", zap.String("type", req.Request.Type))
        w.WriteHeader(http.StatusUnprocessableEntity)
        return
    }

    // текст ответа навыка
    var text string

    switch true {
    // пользователь попросил отправить сообщение
    case strings.HasPrefix(req.Request.Command, "Отправь"):
        // гипотетическая функция parseSendCommand вычленит из запроса логин адресата и текст сообщения
        username, message := parseSendCommand(req.Request.Command)

        // найдём внутренний идентификатор адресата по его логину
        recipientID, err := a.store.FindRecipient(ctx, username)
        if err != nil {
            logger.Log.Debug("cannot find recipient by username", zap.String("username", username), zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // сохраняем новое сообщение в СУБД, после успешного сохранения оно станет доступно для прослушивания получателем
        err = a.store.SaveMessage(ctx, recipientID, store.Message{
            Sender:  req.Session.User.UserID,
            Time:    time.Now(),
            Payload: message,
        })
        if err != nil {
            logger.Log.Debug("cannot save message", zap.String("recipient", recipientID), zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        // Оповестим отправителя об успешности операции
        text = "Сообщение успешно отправлено"

    // пользователь попросил прочитать сообщение
    case strings.HasPrefix(req.Request.Command, "Прочитай"):
        // гипотетическая функция parseReadCommand вычленит из запроса порядковый номер сообщения в списке доступных
        messageIndex := parseReadCommand(req.Request.Command)

        // получим список непрослушанных сообщений пользователя
        messages, err := a.store.ListMessages(ctx, req.Session.User.UserID)
        if err != nil {
            logger.Log.Debug("cannot load messages for user", zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        text = "Для вас нет новых сообщений."
        if len(messages) < messageIndex {
            // пользователь попросил прочитать сообщение, которого нет
            text = "Такого сообщения не существует."
        } else {
            // получим сообщение по идентификатору
            messageID := messages[messageIndex].ID
            message, err := a.store.GetMessage(ctx, messageID)
            if err != nil {
                logger.Log.Debug("cannot load message", zap.Int64("id", messageID), zap.Error(err))
                w.WriteHeader(http.StatusInternalServerError)
                return
            }

            // передадим текст сообщения в ответе
            text = fmt.Sprintf("Сообщение от %s, отправлено %s: %s", message.Sender, message.Time, message.Payload)
        }

    // если не поняли команду, просто скажем пользователю, сколько у него новых сообщений
    default:
        messages, err := a.store.ListMessages(ctx, req.Session.User.UserID)
        if err != nil {
            logger.Log.Debug("cannot load messages for user", zap.Error(err))
            w.WriteHeader(http.StatusInternalServerError)
            return
        }

        text = "Для вас нет новых сообщений."
        if len(messages) > 0 {
            text = fmt.Sprintf("Для вас %d новых сообщений.", len(messages))
        }

        // первый запрос новой сессии
        if req.Session.New {
            // обработаем поле Timezone запроса
            tz, err := time.LoadLocation(req.Timezone)
            if err != nil {
                logger.Log.Debug("cannot parse timezone")
                w.WriteHeader(http.StatusBadRequest)
                return
            }

            // получим текущее время в часовом поясе пользователя
            now := time.Now().In(tz)
            hour, minute, _ := now.Clock()

            // формируем новый текст приветствия
            text = fmt.Sprintf("Точное время %d часов, %d минут. %s", hour, minute, text)
        }
    }

    // заполним модель ответа
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

У навыка Алисы наконец-то появилась реализация основной функция! Теперь можно прослушивать приветствие, полученные сообщения и отправлять новые другим пользователям.

В следующем инкременте мы добавим в наш сервис интроспекцию ошибок, чтобы коммуницировать с пользователем о проблемах в более понятном виде.

# Инкремент 12

## Задание по треку «Сервис сокращения URL»

Добавьте новый хендлер `POST /api/shorten/batch`, принимающий в теле запроса множество URL для сокращения в формате:
```js
[
    {
        "correlation_id": "<строковый идентификатор>",
        "original_url": "<URL для сокращения>"
    },
    ...
]
```

В качестве ответа хендлер должен возвращать данные в формате:
```js
[
    {
        "correlation_id": "<строковый идентификатор из объекта запроса>",
        "short_url": "<результирующий сокращённый URL>"
    },
    ...
]
```

Все записи о коротких URL сохраняйте в базе данных. Не забудьте добавить реализацию для сохранения в файл и в память.

Стоит помнить, что:
- нужно соблюдать обратную совместимость;
- отправлять пустые батчи не нужно;
- вы умеете сжимать контент по алгоритму gzip;
- изменение в базе можно выполнять в рамках одной транзакции или одного запроса;
- необходимо избегать формирования условий для возникновения состояния гонки (race condition).

## Задание по треку «Сервис сбора метрик и алертинга»

Сервер:
- Добавьте новый хендлер `POST /updates/`, принимающий в теле запроса множество метрик в формате: `[]Metrics` (списка метрик).

Агент:
- Научите агент работать с использованием нового API (отправлять метрики батчами).

Стоит помнить, что:
- нужно соблюдать обратную совместимость;
- отправлять пустые батчи не нужно;
- вы умеете сжимать контент по алгоритму gzip;
- изменение в базе можно выполнять в рамках одной транзакции или одного запроса;
- необходимо избегать формирования условий для возникновения состояния гонки (race condition).

# Что вы узнали

Поздравляем! Вы прошли достаточно сложную тему. Теперь ваш список достижений стал ещё больше. Тем не менее, не забывайте, что SQL требует отдельного изучения, а в уроках вы рассмотрели только базовые SQL-конструкции и запросы.

В этой теме вы:
- узнали, как в Go устроены базы данных SQL, драйверы, пакет `database/sql` и дополнительный пакет `sqlx`;
- научились выполнять запросы `SELECT`, `INSERT`, `DELETE`, `UPDATE`, писать для метода `Scan` собственный обрабатываемый тип данных и вставлять записи пачками;
- поняли, почему не стоит использовать ORM в Go.

![](./assets/images/evolution_map_9.png)

