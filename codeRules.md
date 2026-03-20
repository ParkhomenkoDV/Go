# Правила написания и оформления кода Go

1. **Переменные в request неприкасаемы для перезаписи!** Используй структуру `Event`:
```go
type Event struct {
    Req *Request
    // ...
}
```

1. Массив структур объектов лучше масива атрибутов объектов. Агрегируй значения по смыслу объекта перед основной логикой. 
```go
// плохая практика
type Event struct {
    ParticipantID   []int
    ParticipantName []string
    ParticipantAge  []int
}
```
```go
// хорошая практика
type Event struct {
    Participants []Participant
}

type Participant struct {
    ID   int
    Name string
    Age  int
}
```

1. Комментарии спасают! Пиши их как для своей бабушки
```go
// UserService предоставляет методы для работы с пользователями.
// Не является потокобезопасным - синхронизацию должен обеспечивать вызывающий код.
type UserService struct {
    // ...
}

// Create создает нового пользователя с заданными параметрами.
// Возвращает ErrEmailExists если пользователь с таким email уже существует.
// Внутренние ошибки БД оборачиваются в ErrRepository.
func (s *UserService) Create(ctx context.Context, req UserCreateRequest) (*User, error) {
    // ...
}
```

1. Функция должна возвращать новопосчитанное значение **без изменения атрибутов структуры внутри**!
```go
// плохая практика
func GetPayment(event *Event) float64 {
    // представим тут страшную скрытую логику расчета KeyRate
    event.KeyRate += xmath.Pow(math.log(float64(event.Term))-1, 12) / 1200 // страшная скрытая логика 
    event.Payment := event.KeyRate * event.Amount
    return event.Payment
}
func main(event *Event) {
    paymentBefore := event.Payment // предыдущее значение Payment
    event.Payment = GetPayment(event)
    if !IsOkPayment(event.Payment) {
        event.Payment = paymentBefore // откат Payment назад, но не KeyRate => ошибки
    }
}
```
```go
// хорошая практика
func GetPayment(keyRate float64, amount float64) float64 {
    return keyRate * amount // независимые параменные без скрытой логики
}
keyRateNew = event.KeyRate + 1 // очевидная логика
paymentNew = GetPayment(keyRate, event.Amount)
if IsOkPayment(payment) {
    event.KeyRate = keyRateNew
    event.Payment = paymentNew
}
```

1. Магические числа должны быть в пакете `constants`
```go
package main

import co "constants"

func main() {
    // плохая практика
    if age > 18 { // магическое число
        // ...
    }

    // хорошая практика
    if age > co.LegalAge { // теперь хотя бы понятно что за магия и где о ней почитать 
        // ...
    }
}
```

1. Название переменных, функций и методов должны быть говорящие
```go
// плохая практика
func chkUsr() float64 // chk? usr? что проверяется?
```
```go
// хорошая практика
func IsUserActive() float64 // проверяет активен ли пользователь
```

1. Названия булевых переменных, функций и методов должны начинаться с: `Is`, `Has`, `Can`
```go
var isEnabled bool
func HasPermission() bool
```

1. Функции и методы получения чего-либо БЕЗ расчета должны иметь `Get` префикс
```go
func GetID() (int, error) 
```

1. Функции и методы получения чего-либо С расчетом должны иметь `Calc` префикс
```go
func GetID() (int, error) 
```

Унифицированный нэйминг упрощает поиск по стратегии

## Основной алгоритм стратегии

1. Валидация `request`
2. Создаем событие 
```go
func (event *Event) NewEvent(req *Request) {
    // Никакой бизнес логики. Только техническая составляющая! 
    // Перевод массива атрибутов объектов в массив структур объектов
    // Конвертирование атрибутов в нужный тип => меньше возвращаемых error в бизнес-смысле
    // ...
    return &Event{
        Req: req, // Cырые данные
        // Парсинг остальных полей req в объективные атрибуты
    }
}
```
3. Основная логика (бизнес смысл)
4. Собираем `response`