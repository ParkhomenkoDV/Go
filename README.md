# Туториал по языку программирования Go

![Гофер](/assets/images/Golang_LOGO.png)

# Инициализация проекта

Проект должен иметь:
- Файл `go.mod` = хранит версии зависимостей
- Файл `main.go`: = название проекта. **Название `main` обязательно, если это основной проект!** В остальных случаях название не регламентируется.
    - `package main` = название пакета. **Название `main` обязательно, если это основной пакет!** В остальных случаях пакет называется как и директория, в которой он находится.
    - `func main()` - функция, с которой стартует любой скрипт

Файл `go.mod` cоздается командой:

```bash
go mod init project_name
```

# Компиляция и запуск

Если Golang установлен из `SberUserSoft`, то требуется прописать
```bash
export PATH="$PATH:/usr/local/go/bin"
```

Команда компилирует скрипт, запускает и его и удаляет скомпилированный файл сразу после исполнения:

```bash
go run main.go
```

Компиляция скрипта в исполняемый файл:

```bash
go build main.go
```

Запуск файла:

```bash
./main
```

Компиляция в исполняемый файл с конкретным названием:

```bash
go build -o compiled main.go
```

Запуск файла:

```bash
./compiled
```

# Базовые типы данных

![](/assets/images/base_types.png)

Комментарии в Go:
- Однострочные: `// oneline comment`
- Многострочные: 
```go
/*
manyline 
comment
*/
```

# Переменные

Все переменные в Go должны подчиняться **CamelCase**, ~~snake_case~~ недопустим!

Неиспользуемые переменные недопустимы!

```go
// Наименование пакета.
// В любом проекте обзательно присутствие файла с "package main" и func main()
package main

// Импорты пакетов.
// Служебные методы импортируются из пакетов
import (
    "fmt"
    "reflect"
)

// Объявление глобальной переменной. Единственный способ
var temp = "temp"

// Объявление глобальных констант.
const (
    minRate = 1
    maxRate = 100
)

// Обязательная функция, вход и старт скрипта, зарезирвированное слово.
func main() {
    fmt.Println("hello world")

    // Вариант 1
    a := 10
    fmt.Println(a)
    fmt.Println(reflect.TypeOf(a)) // проверка типа

    // Вариант 2
    var b = 10
    fmt.Println(b)

    // Вариант 3
    var c int // дефолтное значение 0
    fmt.Println(c)
    c = 10
    fmt.Println(c)

    // Нулевые (дефолтные) значения переменных по типам.
    var a1 int  // 0
    var b1 float64  // 0.0
    var c1 string  // ""
    var d1 bool  // false
    var m map[string]int  // nil

    fmt.Println(a1, b1, c1, d1)

    // Множественное присвоение. Вариант 1. Предпочтительный.
    var (
        a2 int     = 10
        b2 string  = "10"
        c2 float64 = 10.0
        d2 bool    = true
        // остальные типы обсудим позже
    )
    fmt.Println(a2, b2, c2, d2)

    // Множественное присвоение. Вариант 2.
    a3, b3, c3, d3 := 10, "10", 10.0, true
    fmt.Println(a3, b3, c3, d3)

    // Объявление одной константы.
    const myConst = 100
    fmt.Println(myConst)

    // Множественное объявление констант.
    const (
        minLimit int = 100
        maxLimit     = 1000.5
    )
    fmt.Println(minLimit, maxLimit)

    // Печать глобальной переменной.
    fmt.Println(temp)

    // Печать глобальных констант.
    fmt.Println(minRate, maxRate)
}
```

# Область видимости

![](/assets/images/view.png)

# Сигнатура

![](/assets/images/func.png)

# Функции

Именованных аргументов нет!

# Функция без аргументов

```go
package main

import "fmt"

func myFunc() {
    a := 10
    fmt.Println(a)
}

func main() {
    myFunc()
}
```

# Функция с аргументами одинакового типа

```go
package main

import "fmt"

func myFunc(arg1, arg2 int) int {
    return arg1 + arg2
}

func main() {
    result := myFunc(1, 2)
    fmt.Println(result)
}
```

# Функция с аргументами разного типа

```go
package main

import "fmt"

func myFunc(arg1 int, arg2 float64) int {
    return arg1 + int(arg2)  // со строками такой перевод не работает!!!
}

func main() {
    result := myFunc(1, 2)
    fmt.Println(result)
}
```

# Функция с несколькими возвращаемыми значениями

```go
package main

import "fmt"

func myFunc(arg1, arg2 int) (int, bool) {
    if arg1 > 0 {
        return arg1 + arg2, true
    }
    
    return 0, false
}

func main() {
    result, flag := myFunc(1, 2)
    fmt.Println(result, flag)
}
```

# Функция с предварительно известными возвращаемыми переменными

```go
package main

import "fmt"

func myFunc(arg1, arg2 int) (sum int, isPositive bool) {
    sum = arg1 + arg2
    isPositive = sum > 0
    
    return // указывать ничего не нужно: автоматом выйдут sum и isPositive
}

func main() {
    result, flag := myFunc(1, 2)
    fmt.Println(result, flag)
}
```

# Неограниченное количество аргументов

```go
package main

import (
    "fmt"
    "reflect"
)

func myFunc(args ...int) int {
    fmt.Println(args)
    fmt.Println(reflect.TypeOf(args))
    
    var sum int
    
    for _, val := range args {
        sum += val
    }
    
    return sum
}

func main() {
    result := myFunc(1, 2)
    fmt.Println(result)
}
```

# Безымянная функция

```go
package main

import (
    "fmt"
)

// lambda-функция. Вариант 1
var sum = func(foo, bar int) int {
    return foo + bar
}

func main() {
    // lambda-функция. Вариант 1
    result := sum(10, 20)
    fmt.Println(result)
    
    // lambda-функция. Вариант 2
    var lambdaRes int
    
    func(foo, bar int) {
        lambdaRes = foo + bar
    }(10, 20)
    
    fmt.Println("Res: ", lambdaRes)
}
```

# Функция, как аргумент другой функции

```go
package main

import (
    "fmt"
)

func squareOfSums(foo, bar int, sumFunc func(a, b int) int) int {
    return sumFunc(foo, bar) * sumFunc(foo, bar) // 30 * 30 = 900
}

func main() {
    sum := func(foo, bar int) int {
        return foo + bar
    }
    
    fmt.Println(squareOfSums(10, 20, sum))
}
```

# Замыкание = принятие функции и сохранение своего состояния

```go
package main

import (
    "fmt"
)

func closure() func() int {
    counter := 0
    return func() int {
        counter++
        return counter
    }
}

func main() {
    count := closure()  // 0
    
    fmt.Println(count())  // 1
    fmt.Println(count())  // 2
    fmt.Println(count())  // 3
    fmt.Println(count())  // 4
    fmt.Println(count())  // 5
}
```

# Операторы сравнения
```go
package main

import (
    "fmt"
)

func main() {
    var (
        foo = 10
        bar = 20
    )
    
    fmt.Println(foo == bar)
    fmt.Println(foo != bar)
    fmt.Println(foo > bar)
    fmt.Println(foo < bar)
    fmt.Println(foo >= bar)
    fmt.Println(foo <= bar)
    
    fmt.Println(foo > 10 || bar < 20)  // or
    fmt.Println(foo > 10 && bar < 20)  // and
}
```

# Ветвление
```go
package main

import (
    "fmt"
)

func main() {
    num := 10
    
    if num%2 == 0 {
        fmt.Println("divisible by 2")
    } else if num%3 == 0 {  // elif
        fmt.Println("divisible by 3")
    } else if num%5 == 0 {  // elif
        fmt.Println("divisible by 5")
    } else {
        fmt.Println("else branch")
    }
}
```

# Область видимости

Переменная, **объявленная выше**, видна **во вложенных ниже** фигурных скобках.

```go
package main

import (
	"fmt"
)

func main() {
    if true {
        // переменная number объявлена в блоке 'if'
        number := 5
        if number == 5 {
            // переменная видна внутри блока
            fmt.Println(number)
        }
    }
    // Переменная выделена красным цветом, поскольку переменная number существует только в блоке условия.
    // За пределами блока переменная не существует
    // раскомментируйте ниже
    //fmt.Println(number)
}
```

# Switch

## Switch со значением

```go
package main

import (
    "fmt"
)

func main() {
    num := 2
    
	// if. else if. else.
    if num == 0 {
        fmt.Println("0")
    } else if num == 1 {
        fmt.Println("One")
    } else if num == 2 {
        fmt.Println("Two")
    } else {
        fmt.Println("Unknown Number")
    }
	
	// switch
    switch num {
    case 0:
        fmt.Println("Zero")
    case 1:
        fmt.Println("One")
    case 2:
        fmt.Println("Two")
    // case может содержать множество значений для сравнения
    case 3, 4, 5:
        fmt.Println("3, 4, 5")
    // выполняется, если остальные проверки не прошли (else)
    default:
        fmt.Println("Unknown Number")
    }
}
```

```go
package main

import (
    "fmt"
)

func main() {
    var (
        foo = 10
        bar = 20
    )
    
    // switch и каждый его case создаёт новую область видимости
    switch res := 9; {
    case foo-res == bar:
        fmt.Println("foo minus res equals bar")
    case bar%foo == 0:
        fmt.Println("bar is divisible by foo")
    case foo%bar == 2:
        // прерывает switch
        break
    default:
        fmt.Println("default branch")
    }
    
    // переменная res не видна за пределами switch
    //fmt.Println(res)
}
```

## Switch без значения

```go
package main

import "fmt"

func main() {
    var foo = 10
    
    switch {
    case foo > 1:
        fmt.Println("foo > 1")
    case foo < 1:
        fmt.Println("foo < 1")
    default:
    }
}
```

# Bool сравнения
```go
package main

import "fmt"

// Bool сравнения.
func ex8(num1, num2 int) bool {
    if num1+num2 > 20 {
        return true
    }
    
    return false
}

func main() {
    // Bool сравнения.
    if ex8(1, 20) {
        fmt.Println("> 20")
    }
    if !ex8(1, 2) {
        fmt.Println("< 20")
    }
}
```

# Ошибки

Ошибки = отдельный интерфейс

# Инициализация кастомной ошибки

```go
package main

import (
    "errors"
    "fmt"
    "reflect"
)

func main() {
    // Вариант 1. Используется если ошибка встречается часто
    customErr := errors.New("наша ошибка")
    fmt.Println(reflect.TypeOf(customErr))
    fmt.Println(customErr.Error()) // преобразование из ошибки в строку
    fmt.Println(customErr)
    
    // Вариант 2. Использвется если ошибка нужна здесь и сейчас
    customErr1 := fmt.Errorf("наша ошибка 2")
    fmt.Println(reflect.TypeOf(customErr1))
    fmt.Println(customErr1)
    
    // Ошибка с нулевым дефолтным значением типа (nil).
    var customErr2 error
    fmt.Println(reflect.TypeOf(customErr2))
    fmt.Println(customErr2)
    customErr2 = fmt.Errorf("наша ошибка 2")
    fmt.Println(customErr2)
    fmt.Println(reflect.TypeOf(customErr2))
}
```

## Возврат числа и ошибки

```go
package main

import (
    "fmt"
    "log"
)

const (
    maxSumForStudent = 1000
    maxSumForPension = 2000
)

var AgeErr = errors.New("неподходящий возраст")

// Возврат числа и ошибки.
func errReturn(age int) (int, error) {
    switch {
    case age < 18:
        return 0, AgeErr
    case age < 25:
        return maxSumForStudent, nil
    case age < 80 && age > 55:
        return maxSumForPension, nil
    default:
        return 0, fmt.Errorf("возраст не подходит")
    }
}

func main() {
    age := 11
    result, err := errReturn(age)
    // обработка ошибки
    if err != nil {
        log.Fatalln(err)
    }
    
    fmt.Println(result)
}
```

# Оборачивание ошибки

Если функция, которую вызывает другая функция может вернуть ошибку, то считается хорошим тоном ее обработать.

Функция `wrapError()` вызывает функцию `createAndReturnError()`, которая может вернуть ошибку. 

Что бы "прокинуть" эту ошибку дальше по стеку, в нашем случае в функцию `main()`, нужно в функции `wrapError()` создать новую ошибку и добавить к ней текст ошибки из `createAndReturnError()`, упаковать одну ошибку в другую.

Это позволит более разернуто видеть причину, текст ошибки и проще найти место ее возникновения. 

```go
package main

import "fmt"

// создает и возвращает ошибку
func createAndReturnError(num int) (int, error) {
    if num == 0 {
        return 0, fmt.Errorf("my custom error %d", num)
    }
    
    return num, nil
}

// оборачивает ошибку и передает ее дальше
func wrapError(num int) (int, error) {
    res, err := createAndReturnError(num)
    if err != nil {
        return 0, fmt.Errorf("wrap error: %w", err)
    }
    
    return res, nil
}

func main() {
    res, err := wrapError(0)
    fmt.Println(res, err)
    
    res, err = wrapError(1)
    fmt.Println(res, err)
}
```

## Приведение к типу

```go
package main

import (
    "fmt"
    "log"
    "reflect"
    "strconv"
)

// Приведение int к string.
func intToString() string {
    number := 10
    stringResult := strconv.Itoa(number)  // string(number) не работает!!!
    
    return stringResult
}

// Приведение string к int.
func stringToInt(str string) (int, error) {
    intResult, err := strconv.Atoi(str)
    if err != nil {
        return 0, fmt.Errorf("func stringToInt %w", err)
    }
    
    return intResult, nil
}

func main() {
    // Приведение int к string.
    resStr := intToString()
    fmt.Println(resStr)
    fmt.Println(reflect.TypeOf(resStr))
    
    // Приведение string к int.
    resInt, err := stringToInt("1235")
    if err != nil {
        log.Fatalln("ошибка преобразования в строку: %w", err.Error()) // Error() - преобразование к строчке
    }
    
    fmt.Println(resInt)
    fmt.Println(reflect.TypeOf(resInt))

    resInt, err := stringToInt("aa")
    if err != nil {
        log.Fatalln("ошибка преобразования в строку: %w", err.Error()) // Error() - преобразование к строчке
    }
    
    fmt.Println(resInt)
    fmt.Println(reflect.TypeOf(resInt))
}
```

## Обработка ошибки в 1 строку

```go
package main

import (
    "errors"
    "log"
)

func oneLineErr() error {
    return errors.New("error")
}

func main() {
    var err error

    // 2 строчки
    err = oneLineErr()
    if err != nil {
        log.Fatalln(err)
    }

    // 1 строчка
    if err = oneLineErr(); err != nil {
        log.Fatalln(err)
    }
}
```

## Defer = Отложенный вызов

```go
package main

import (
    "fmt"
)

func deferFunc(foo, bar int) int {
    // выполнится после return
    defer fmt.Println("3. Defer call")
    
    fmt.Println("1. Start func")
    res := foo + bar
    fmt.Printf("2. Print res in func %d\n", res)
    
    return foo + bar
    // вот тут случиться defer
}

func main() {
    res := deferFunc(10, 20)
    fmt.Printf("4. Print res in main %d", res)
}
```

## errors.Is()
`errors.Is` позволяет выяснить, содержит ли ошибка `err` текст ошибки `timeOutErr - time out`.

`errors.Is` пройдет по цепочке оберток и вернет bool, если найдет/не найдет текст ошибки. 

```go
package main

import (
    "errors"
    "fmt"
)

func main() {
    timeOutErr := errors.New("time out")
    
    // обертка для timeOutErr
    err := fmt.Errorf("wrap: %w", timeOutErr)
    
    if errors.Is(err, timeOutErr) {
        fmt.Println("my error")
    } else {
        fmt.Println("not my error")
    }
}
```

## errors.As()
`errors.As` проверяет совпадение текущей ошибки `err` с переданным типом `pathError`.

```go
package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
)

func main() {
    if _, err := os.Open("non-existing"); err != nil {
        var pathError *fs.PathError
        if errors.As(err, &pathError) {
            fmt.Println("Failed at path:", pathError.Path)
        } else {
            fmt.Println(err)
        }
    }

	//Output: Failed at path: non - existing
}
```

# Паника = авариное завершение программы, если нет recover

```go
package main

import (
    "fmt"
    "log"
)

// Паникующая функция.
func makePanic() {
    panic("my panic")  // вынужденная ошибка
}

// Обращение к несуществующему индексу
func indexOutOfRange() {
    arr := []string{"index0", "index1", "index2"}
    fmt.Println(arr[3])
}

// Разыменование пустого указателя (nil)
func nilPointerException() {
    // nil pointer dereference
    var foo *int
    fmt.Println(*foo)
}

// Обработка паники.
func recoverPanic() {
    defer func() {
        if err := recover(); err != nil {
            log.Println("panic name:", err)
        }
    }()
    
    fmt.Println("do smth dangerous!")
    makePanic()
    //indexOutOfRange()
    //nilPointerException()
    fmt.Println("Ooops!")
}

func main() {
    recoverPanic()
    fmt.Println("Program continues")
}
```

# Указатели

![](/assets/images/ptr1.png)


- foo - переменная
- int - тип переменной(целое число)
- значение переменной (20)
- &foo - указатель на переменную foo
- *int - тип указателя(указатель на целое число)
- значение. Условные номер ячейки в памяти.

Вызов `&foo` вернет значение `0xv00000780`. По этому адресу находится значение переменной `foo` 20.

- bar - переменная, указатель на foo
- *int - тип указателя(указатель на целое число)
- значение переменной (0xv00000780)

Вызов `*bar` вернет значение 20, так как мы обратильсь к `0xv00000780` и разыменовали этот указатель, т.е получили
значение.

## Пустой указатель. Дефолтное значение

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    // пустой указатель
    var intPointer *int
    
    fmt.Println(intPointer)                 // nil
    fmt.Println(reflect.TypeOf(intPointer)) // *int
}
```

## Указатель со значением
```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    a := 10
    
    pointerOnA := &a
    fmt.Println(pointerOnA)                 // 0xc00001c088 - услоный номер ячейки в RAM
    fmt.Println(reflect.TypeOf(pointerOnA)) // *int
    
    // разименование
    fmt.Println(*pointerOnA) // 10
}
```

## Передача аргумента по значению (по копии)
```go
package main

import "fmt"

func argByValue(num int) int {
    return num + 100
}

func main() {
    numberValue := 10
    fmt.Println(numberValue) // 10
    
    res := argByValue(numberValue)
    
    fmt.Println(numberValue) // 10
    fmt.Println(res)         // 110
}
```

## Передача аргумента по указателю
```go
package main

import "fmt"

func argByPointer(num *int) {  // возвращать не надо, т.к. модификация происходит в ячеке памяти
    *num += 100  // * = разыменование
}

func main() {
    numberPointer := 10
    
    fmt.Println(numberPointer) // 10
    
    argByPointer(&numberPointer)  // переприсваивать ничего не надо. Все происходит в ячеке памяти
    
    fmt.Println(numberPointer) // 110
}
```

## Ошибка разыменования пустого указателя
```go
package main

import "fmt"

func nilPointerException() {
    a := 10
    // не пустой указатель
    pointerOnA := &a
    // разименование без ошибок
    fmt.Println(*pointerOnA) // 10
    
    // пустой указатель (значение nil)
    var intPointer *int
    // ошибока разименования
    fmt.Println(*intPointer) // panic: runtime error: invalid memory address or nil pointer dereference
}

func main() {
    nilPointerException()
}
```

## Пример. Проверка активности счета

```go
package main

import "fmt"

// Проверка активности счета.
func isActiveAccount(num *int) bool {
    if num != nil {
        return true
    }
    // если значение баланса счета равно nil
    return false
}

func main() {
    // счет активный и его баланс 10
    accountBalance1 := 10
    fmt.Println(isActiveAccount(&accountBalance1))  // true
    
    // счет активный и его баланс 0, явно указан
    accountBalance2 := 0
    fmt.Println(isActiveAccount(&accountBalance2))  // true
    
    // счет активный и его баланс 0, дефолтное значение типа int
    var accountBalance3 int
    fmt.Println(isActiveAccount(&accountBalance3))  // true
    
    // счет не активный, дефолтное значение типа *int (указатель на int) равен nil
    var accountBalance4 *int
    fmt.Println(isActiveAccount(accountBalance4))  // false
}
```

**С простыми (несоставными) типами: int, float, string через указатели работать не надо!** Это итак происходит быстро. Разыменование и взятие адреса лишь дополнительно создаст нагрзку. Указатели показывают себя во всей красе со структурами

```go
type Customer struct {
    Name string
    Surname string
    Age int
}

// хорошая практика: экономия памяти и быстродействие
func DoSmthWithCustomer(c *Customer) {
    fmt.Println(c.Age)
    //Разыменование * для структуры не нужна. Go понимает сам
    c.Age = 0  // я родился!
    fmt.Println(c.Age)
}

func main() {
    customer := Customer{
        Name "Daniil"
        Surname "Andryushin"
        Age 24
    }

    fmt.Println(customer)
    DoSmthWithCustomer(&customer)
    fmt.Println(customer)
}
```

# Массивы

Массив = структра **однотипных** данных, которая выделяет *непрерывню цепь в памяти*.

Массив в Go = **неизменяемая** структра. Метод `append` реализован только для slice!

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    // пустой массив смысла не имеет!
    var array0 [0]int
    fmt.Println(array0)
    fmt.Println(reflect.TypeOf(array0))
    
    // пустой массив с дефолтными значениями типа
    var array1 [3]int
    fmt.Println(array1)
    fmt.Println(reflect.TypeOf(array1))
    
    // пустой массив с дефолтными значениями типа
    array2 := [3]int{}
    fmt.Println(array2)
    fmt.Println(reflect.TypeOf(array2))
    
    // массив с заполненными значениями
    array3 := [3]int{1, 2, 34}
    fmt.Println(array3)
    fmt.Println(reflect.TypeOf(array3))
}
```

## Массив, как аргумент функции

```go
package main

import (
    "fmt"
)

func arrArg(arr [3]int) [3]int {
    fmt.Println(arr)
    return arr
}

func main() {
    arr := [3]int{1, 2, 3}
    res := arrArg(arr) // отправляется копия!
    fmt.Println(res)
}
```

## Индексация

```go
package main

import (
    "fmt"
)

func main() {
    arr := [5]int{2, 33, 343, 88, 99}
    fmt.Println(arr[0])
    
    // получить значение
    a := arr[2]
    fmt.Println(a)
    
    // изменить значение
    arr[0] = 1000
    fmt.Println(arr)
}
```

## Длина и вместимость

Длина = количество элементов массива
Вместимость = количество элементов массива
Для slice длина != вместимость

```go
package main

import (
    "fmt"
)

func main() {
    arr := [5]int{2, 33, 343, 88, 99}
    
    // для массивов одинаковые значения
    fmt.Println(len(arr))
    fmt.Println(cap(arr))
}
```

## Передача массива в функцию по указателю
```go
package main

import (
    "fmt"
)

// Аргумент передан по значению (создается копия arr)
func arrOnValue(arr [3]int) [3]int {
    arr[2] = 100
    return arr
}

// Аргумент передан по указателю. Работаем с оригиналом arr из main()
func arrOnPointer(arr *[3]int) {
    arr[2] = 100  // разыменование не нужно! Go все понимает
}

func main() {
    arr := [3]int{1, 2, 3}
    
    fmt.Println(arr) // [1 2 3]
    res := arrOnValue(arr)
    fmt.Println(arr) // [1 2 3]
    fmt.Println(res) // [1 2 100]
    
    fmt.Println(arr) // [1 2 3]
    arrOnPointer(&arr)
    fmt.Println(arr) // [1 2 100]
}
```

## Работа с последним элементом

**Отрицательных индексов в Go нет!**

```go
package main

import (
    "fmt"
    "reflect"
)

func main() {
    // получить последний элемент
    array1 := [5]int{2, 33, 343, 88, 99}
    last1 := array1[len(array1)-1]
    fmt.Println(last1) // 99
    
    // удалить последний элемент. Требуется создание нового слайса
    array2 := [5]int{2, 33, 343, 88, 99}
    newAr := array2[:len(array2)-1]
    fmt.Println(newAr)                 // [2 33 343 88]
    fmt.Println(reflect.TypeOf(newAr)) // []int
}
```

Внимание!
- *[3]int = указатель на массив из 3х элементов
- [3]*int = массив из 3х элементов типа указатель

Срез от массива = slice

```go
import (
    "fmt"
    "reflect"
)
func main() {
    arr := [5]int{1,2,3,4,5}
    fmt.Printf("Type of %v is %T", arr, reflect.TypeOf(arr))
    fmt.Printf("Type of %v is %T", arr[1:3], reflect.TypeOf(arr[1:3]))
}
```

# Slices

## Создание среза (слайса)

```go
package main

import (
    "fmt"
)

func main() {
    var sl1 []int
    fmt.Println(sl1) // []
    
    sl2 := []int{1, 2, 3}
    fmt.Println(sl2) // [1 2 3]
    
    var sl3 = []int{1, 2, 3}
    fmt.Println(sl3) // [1 2 3]

    fmt.Println(sl3[3]) // panic = выход за границу массива
}
```

## Добавление элемента в слайс

```go
package main

import (
    "fmt"
)

func main() {
    sl := []int{1, 2, 3}
    fmt.Println(sl) // [1 2 3]
    
    sl = append(sl, 4)
    fmt.Println(sl) // [1 2 3 4]
    sl = append(sl, 5, 6, 7)  // добавление в конец
    fmt.Println(sl) // [1 2 3 4 5]
}
```

Slice = вид на массив, который имеет длину.

```go
type Slice struct {
    Arr *arr  // настоящий массив на который ссылается Slice
    len int
    cap int
}
```

- Длина = количество элементов в Slice
- Вместимость = количество элементов в массиве, на который смотри Slice

При нехватке места для добавления элементов выделяется **новый** массив двойной вместимости от первоначальной. В обратную сторону вместимость не уменьшается!

## Длина `len()` и вместимость `cap()`
```go
package main

import "fmt"

func main() {
    // cap = 3
    sl := []int{1, 2, 3}
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl, len(sl), cap(sl)) // срез:[1 2 3], len:3, cap:3
    
    // cap увеличился в 2 раза = 6
    sl = append(sl, 4)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl, len(sl), cap(sl)) // срез:[1 2 3 4], len:4, cap:6
    
    // cap = 6
    sl = append(sl, 5)
    sl = append(sl, 6)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl, len(sl), cap(sl)) // срез:[1 2 3 4 5 6], len:6, cap:6
    
    // cap увеличился в 2 раза = 12
    sl = append(sl, 7)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl, len(sl), cap(sl)) // срез:[1 2 3 4 5 6 7], len:7, cap:12
}
```

## Создание слайса через `make` с заданными len и cap

```go
package main

import (
    "fmt"
)

func main() {
    //  len - 0 cap - 5
    sl1 := make([]int, 0, 5)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl1, len(sl1), cap(sl1)) // срез:[], len:0, cap:5
    sl1 = append(sl1, 1)
    sl1 = append(sl1, 2)
    sl1 = append(sl1, 3)
    sl1 = append(sl1, 4)
    sl1 = append(sl1, 5)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl1, len(sl1), cap(sl1)) // срез:[1 2 3 4 5], len:5, cap:5
    sl1 = append(sl1, 6)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl1, len(sl1), cap(sl1)) // срез:[1 2 3 4 5 6], len:6, cap:10
    
    // 0 0 0
    sl2 := make([]int, 3, 5)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl2, len(sl2), cap(sl2)) // срез:[0 0 0], len:3, cap:5
    sl2 = append(sl2, 1)
    sl2 = append(sl2, 2)
    sl2 = append(sl2, 3)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl2, len(sl2), cap(sl2)) // срез:[0 0 0 1 2 3], len:6, cap:10
    
    // если задавать только len, то cap будет того же размера
    sl3 := make([]int, 3)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl3, len(sl3), cap(sl3)) // срез:[0 0 0], len:3, cap:3
}
```

## Передача слайса в функцию
Слайс передается по указателю, т.к. является указательным типом. 
**Нет необходимости передавать его по указателю!**

```go
package main

import "fmt"

func sliceChange(sl []int) {
    sl[0] = 100
}

func main() {
    sl := []int{1, 2, 3}
    
    fmt.Println(sl) // [1 2 3]
    sliceChange(sl)
    fmt.Println(sl) // [100 2 3]
}
```

## Возврат слайса из функции
Если в функции изменяется len и cap переданного слайса, то требуется возвращать новый.
**При работе со Slice всегда должен быть `return`!** 

```go
package main

import "fmt"

// Если в функции изменяется len и cap, то требуется возвращать новый слайс.
func sliceNewArray(sl []int) []int {
    sl = append(sl, 4)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl, len(sl), cap(sl)) // срез:[1 2 3 4], len:4, cap:6
    
    return sl
}

func main() {
    // Return слайса из функции.
    sl1 := []int{1, 2, 3}
    fmt.Printf("срез:%v, len:%v, cap:%v\n", sl1, len(sl1), cap(sl1)) // срез:[1 2 3], len:3, cap:3
    
    // перезаписался
    newSl := sliceNewArray(sl1)
    fmt.Printf("срез:%v, len:%v, cap:%v\n", newSl, len(newSl), cap(newSl)) // срез:[1 2 3 4], len:4, cap:6
}
```

## Получение среза от слайса
```go
package main

import "fmt"

func main() {
    arr := []int{1, 2, 3, 4, 5, 6}
	
    // cap высчитывается от индекса первого взятого элемента до конца массива
    newArr := arr[2:4]
    fmt.Printf("срез:%v, len:%v, cap:%v\n", newArr, len(newArr), cap(newArr)) // срез:[3 4], len:2, cap:4
    
    sl := []int{1, 2, 3, 4, 5, 6}
    newSl := sl[:]
    fmt.Printf("срез:%v, len:%v, cap:%v\n", newSl, len(newSl), cap(newSl)) // срез:[1 2 3 4 5 6], len:6, cap:6
}
```

## Удаление элемента слайса по индексу
```go
package main

import "fmt"

func main() {
    sl := []int{1, 22, 334, 24, 35, 46}
    fmt.Println(sl) // [1 22 334 24 35 46]
    
    index := 2
    // создается новый слайс
    newSl := append(sl[:index], sl[index+1:]...)  // многоточие обязательно
    fmt.Println(newSl) // [1 22 24 35 46]
    fmt.Println(sl) // [1 22 24 35 46 46]
}
```

## Сортировка слайса с типом int

```go
package main

import (
	"fmt"
	"sort"
)

func main() {
	sl := []int{24, 1, 22, 334, 35, 46}

	fmt.Println(sl) // [24 1 22 334 35 46]
	sort.Ints(sl)  // работает по указателям => переприсваивание не нужно
	fmt.Println(sl) // [1 22 24 35 46 334]
}
```

# Циклы

## Цикл for c использованием `range`
```go
package main

import "fmt"

func main() {
    sl := []int{12, 25, 34, 46, 59}
    
    // индекс и значение
    for idx, val := range sl {
        fmt.Printf("index:%v, value:%v\n", idx, val)
    }
    
    // значение
    for _, val := range sl {  // первый всегда индекс!
        fmt.Printf("value:%v\n", val)
    }
    
    // индекс
    for idx := range sl {
        fmt.Printf("index:%v\n", idx)
    }

    // Цикл индексов через range. Доступно с 1.22
    for idx := range 10 {
        fmt.Println(idx)
    }
}
```

## Цикл "while"
```go
package main

import "fmt"

func main() {
    sl := []int{12, 25, 34, 46, 59}
    
    // варинат 1
    j := 0
    for j < len(sl) {
        fmt.Println(sl[j])
        j += 1
    }
    
    // вариант 2
    for i := 0; i < len(sl); i++ {
        fmt.Printf("counter:%v, value:%v\n", i, sl[i])
    }
}
```

## Break
```go
package main

import "fmt"

func main() {
    sl := []int{12, 25, 34, 46, 59}
  
    for _, v := range sl {
        if v == 34 {
            break
        }
        fmt.Println(v)
    }
}
```

## Continue
```go
package main

import "fmt"

func main() {
    sl := []int{12, 25, 34, 46, 59}
    
    for _, v := range sl {
        if v == 34 {
            continue
        }
        fmt.Println(v)
    }
}
```

## Именование цикла
```go
package main

import "fmt"

func main() {
    sl1 := []int{10, 20, 30, 40, 50}
    sl2 := []int{1, 2, 3, 4, 5}
    
    var result []int

mainLoop:
    for _, v := range sl1 {
        for i := 0; i < len(sl2); i++ {
            if v*sl2[i] == 200 {
                break mainLoop  // разрыв цикла mainloop, а не ближайшего цикла for
            }
            result = append(result, v*sl2[i])
        }
    }
    
    fmt.Println(result)
}
```

## ДЗ
Задача-1. 
- Написать функцию task1, которая приниает на вход аргумент типа float64 и его печатает. 

Задача-2.
- Написать функцию task2, которая приниает на вход n int. Вывести все числа от 0 до n.

Задача-3.
- Написать функцию task3, которая принимает на вход два положительных числа k и n (k < n) int. Вывести все числа от k до n

Задача-4.
- Написать функцию task4, которая принимает на вход число n. Найти 2 суммы всех четных и всех нечетных чисел от 0 до n и вернуть из функции 2 значения.

Задача-5
- Написать функции task5ForRange и task5While которые принимают на вход строку и слайс строк.
- Проверить, входит ли строка в список строк. Вернуть true false.
- sl := []string{"abc", "bcd", "xvm", "abc", "abd", "bcd", "abc"}
- Решение в двух вариантах: используя конструкцию for range и for как "while". Питоновского "in" в golang нет :)

Задача-6
- Подсчет количества вхождений заданной строки "abc" в список строк.
- Написать функции task6ForRange и task6While которые принимают на вход строку и слайс строк.
- Посчитать количество вхождений строки в слайс. Если строка не входит в слайс, создать ошибку, вернуть ее и обработать.

Задача-7
- Дан слайс sl := []int{10, 20, 30, 40, 50}, требуетя создать новый перемножив каждое число на множитель.
- Написать функции task7ForRange и task7While которые принимают на вход слайс int и можитель.
- Проверить длинну слайса, если == 0, то создать ошибку, вернуть ее и обработать.
- Требуется избежать ненужных аллокаций памяти(RAM).

Задача-8
- Дан список, заполненный произвольными целыми числами от -50 до 50.
- Написать функцию task8Example1. Получите новый список, элементами которого будут неповторяющиеся(уникальные) элементы
  исходного списка: например, [1, 2, 4, 5, 6, 2, 5, 2], нужно получить [1, 2, 4, 5, 6].
- Написать функцию task8Example2. Вернуть количество элементов списка, значение которых не превышает 3.
- Написать функцию task8Example3. Вернуть сумму всех положительных четных элементов списка.
- Написать функцию task8Example4. Вернуть среднее арифметическое всех положительных нечетных элементов.

Задача-9
- Дан список заполненный произвольными числами от -50 до 50.
- Написать функцию task9, которая принимает на вход список и ищет разницу между самым большим и самым маленьким элементами списка.

# Структуры

## (u User) - value receiver, (u *User) - pointer receiver
```go
package main

import "fmt"

// User структура пользователя.
type User struct {
    ID      int
    Name    string
    Surname string
}
// если в полях не указывать значения, то они будут дефолтными

// New - конструктор структуры User. Нужно писать ручками
func New(id int, name, surname string) *User {
    return &User{  // экономия памяти при передаче через указатель
        ID:      id,
        Name:    name,
        Surname: surname,
    }
}

// SetName - установка пользователю нового имени.
func (u User) SetName1(name string) {  // u = ссылка на структуру. Принято называть первой маленькой буквой 
    u.Name = name  // передача по копии = плохая практика!
}

func (u *User) SetName2(name string) {
    u.Name = name // хорошая практика. Во всех методах надо писать *
}

// GetName - возвращение имени пользователя.
func (u User) GetName1() string {
    return u.Name  
}

func (u *User) GetName2() string {
    return u.Name
}

func main() {
    // Создание пользователя с конструктором.
    user1 := New(1, "Ivan", "Ivanov")
    user2 := User{
        ID: 2
        Name: "Daniil"
        Surname: "Andryushin"
    }

    fmt.Printf("%+v\n", user1) // &{ID:1 Name:Ivan Surname:Ivanov}
    fmt.Printf("%+v\n", user2) // &{ID:2 Name:Daniil Surname:Andryushin}

    // установка нового имени
    user1.SetName1("Andrey")
    user2.SetName1("Andrey")
    fmt.Printf("%+v\n", user1) // &{ID:1 Name:Andrey Surname:Ivanov}
    fmt.Printf("%+v\n", user2) // &{ID:2 Name:Andrey Surname:Andryushin}
    user1.SetName2("Dmitriy")
    user2.SetName2("Dmitriy")
    fmt.Printf("%+v\n", user1) // &{ID:1 Name:Dmitriy Surname:Ivanov}
    fmt.Printf("%+v\n", user2) // &{ID:2 Name:Dmitriy Surname:Andryushin}

    // получение нового имени
    fmt.Println(user1.GetName1()) // Andrey
    fmt.Println(user2.GetName1()) // Andrey
    fmt.Println(user1.GetName2()) // Andrey
    fmt.Println(user2.GetName2()) // Andrey
}
```

# Импорт пакетов

иерархия:
```
temp
|   user
|   |___ user.go
| go.mod
| main.go
```

файл `go.mod`:
```go
module <root> // корневая папка
```

файл main.go
```go
import "temp/user"

u := user.User{
    ...
}
```

# Инкапсуляция

Если обозвать структуру с *маленькой* буквы, то ее импортировать не полчится.
Аналогично для полей структуры.

# Встраивание

```go
package customer

// Customer - структура клента.
type Customer struct {
	ID            int
	Name, Surname string
	Address       Address
	Phones        []Phones
	Scores        Scores
	Account       Account // композиция
}

// Address - структура адреса.
type Address struct {
	Land, City, District string
}

// Phones - структура телефона клиента.
type Phones struct {
	Operator string
	Number   int
}

// Scores - структура скоров клиента.
type Scores struct {
	Score1 []int
	Score2 [5]int
}

// Account - структура активности счета.
type Account struct {
	Balance *int
}

// IsActive - композиция. Проверка активности баланса.
func (a *Account) IsActive() bool {
	return a.Balance != nil
}

// New - конструктор структуры пользователя.
func New(id int, name, surname, land, city, district string) *Customer {
	var balance = 100

	return &Customer{
		ID:      id,
		Name:    name,
		Surname: surname,
		Address: Address{
			Land:     land,
			City:     city,
			District: district,
		},
		Account: Account{Balance: &balance},
		// Phones и Scores оставим пустыми. Они равны дефолтным значениям типа.
	}
}

// AddPhone - добавление полей поля (встроенной структуры) Phones.
func (c *Customer) AddPhone(operator string, phoneNumber int) {
	phone := Phones{
		Operator: operator,
		Number:   phoneNumber,
	}

	c.Phones = append(c.Phones, phone)
}

// AddScores - добавление полей поля (встроенной структуры) Scores.
func (c *Customer) AddScores(scores []int) {
	for i, v := range scores {
		if v%2 == 0 {
			c.Scores.Score1 = append(c.Scores.Score1, v)
		} else {
			c.Scores.Score2[i] = v
		}
	}
}

func main() {
    c := customer.Customer{}
    fmt.Printf("%+v", c)
    c.AddPhone("MTS", 88005553535)
    c.AddPhone("Mgafone", 88009999999)
    fmt.Printf("%+v", c)
}
```

# JSON

## Работа с JSON. Парсинг в структуру

Парсить в map = плохая практика в виду разных типов данных! 
Надо парсить в структуру!


Файл `in.json`:
```json
{
    "score_name": "score_1",
    "score_value": 9999
}
```

Файл `main.go`:
```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "os"
)

// Score - структура, в которую парсится json.
type Score struct {
    Name  string `json:"score_name"` // в теге 'json' указаывается поле из json
    Value int    `json:"score_value"`
}

func main() {
    // чтение файла
    file, err := os.ReadFile("in.json")
    if err != nil {
        log.Fatalln(err)
    }
    
    // объект структуры
    s := Score{}
    
    // парсинг json в структуру
    if err := json.Unmarshal(file, &s); err != nil {
        log.Fatalln(err)
    }
    fmt.Printf("%+v\n", s) // {Name:score_1 Value:9999}
    
    // парсинг структуры в json
    payload, err := json.Marshal(s)
    if err != nil {
        log.Fatalln(err)
    }
    
    // payload - []byte, с которым можно работать дальше
    // отправить по сети, записать в файл и тд
    fmt.Println(payload) // [123 34 115 99 111 114 101 95 110 97 109 101 34 58 34 115 99 111 114 101 95 49 34 44 34 115 99 111 114 101 95 118 97 108 117 101 34 58 57 57 57 57 125]
}
```

## Валидация JSON

Указательные поля необходимы только для того, что при условии, что поля не будет => его значение будет дефолтным, (т.е. для int нулевым), а это уже может что-то значить!!! например 0 штрафов.
А дефолтное значение у указателя это `nil`. Тут сразу все ясно.

Указательными типами мы делаем те поля, которые обязательно должны приходить! 

Иерархия:
```go
temp
    customer
        customer.go
        validation.go
    data
        in.json
        out.json
    main.go
```

Файл `customer.go`:
```go
package customer

// Customer - структура клиента.
type Customer struct {
	ID         int      `json:"id"`
	Name       *string  `json:"name"`  // обязательные поля
	Surname    *string  `json:"surname"`  // обязательные поля
	SalaryFlag *int     `json:"salary_flag"`  // обязательные поля
	Address    Address  `json:"address"`
	Phones     []Phones `json:"phones"`
	Scores     Scores   `json:"scores"`
}

type Address struct {
	Land     string `json:"land"`
	City     string `json:"city"`
	District string `json:"district"`
}

type Phones struct {
	Operator string `json:"operator"`
	Number   int    `json:"number"`
}

type Scores struct {
	Score1 []int  `json:"score_1"`
	Score2 [5]int `json:"score_2"`
}

func New() *Customer {
	return &Customer{}
}
```

Файл `validation.go`:
```go
package customer

import (
	"errors"
	"fmt"
)

// ошибки вынесены в глобальные переменные.
var (
	// ErrRequired - обязательное поле.
	ErrRequired = errors.New("field is required")
	// ErrLess1Greater30 - поле меньше 1 или больше 30 символов.
	ErrLess1Greater30 = errors.New("field value is less than 1 or greater than 30")
	// ErrValue - поле должно иметь значение 0 или 1
	ErrValue = errors.New("field value not 0 or 1")
)

// Validate - валидация структуры.
func Validate(c *Customer) error {
	// проверка обязательного поля
	if c.Name == nil {
		return fmt.Errorf("failed check field 'Name': %w", ErrRequired)
	}

	// проверка обязательного поля
	if c.Surname == nil {
		return fmt.Errorf("failed check field 'Surname': %w", ErrRequired)
	}

	// проверка необязательного поля, значение которого может быть 0 или 1
	if c.SalaryFlag != nil {
		if *c.SalaryFlag != 0 && *c.SalaryFlag != 1 {
			return fmt.Errorf("failed check field 'SalaryFlag': %w", ErrValue)
		}
	}

	// проверка необязательного поля на условие значения
	for idx, score := range c.Scores.Score2 {
		if score < 1 || score > 30 {
			return fmt.Errorf("failed check field 'Scores.Score2': %w, err in %dth array index", ErrLess1Greater30, idx)
		}
	}

	return nil
}
```

Файл `main.go`:
```go
package main

func readFile() ([]byte, error) {
    file, err := os.ReadFile("data/in.json")
    if err != nil {
        return nil, fmt.Errorf("file read error: %w", err)
    }
    return file, nil
}

func jsonParse(file []byte) (*customer.Customer, err) {
    c := customer.Customer.New()
    err := json.Unmarshal(file, c)
    if err != nil {
        log.Fatalln(err)
    }
    return c, nil
}

func fileWrite(c *customer.Customer) error {
    payload, err := json.Marshal(c)
    if err != nil {
        return fmt.Errorf("marshal error: %w", err)
    }
    if err = os.WriteFile("out.json", payload, 0644) {
        log.Fatalln(err)
    }
    return nil
}

func main() {
    file, err := readFile()
    if err != nil {
        log.Fatalln(err)
    }
    fmt.Println(file)
    c, err := json.Parse(file)

    fmt.Printf("%+v\n", c)

    if err := fileWrite(c); err != nil {
        log.Fatalln(err)
    }
}
```

`in.json`
```json
{
    "id": 1,
    "name": "Ivan",
    "surname": "Ivanov",
    "salary_flag": 1,
    "address": {
        "land": "Russia",
        "city": "Moscow",
        "district": "Butovo"
    },
    "phones": [
        {
            "operator": "MTS",
            "number": 89850000001
        },
        {
            "operator": "МегаФон",
            "number": 89850000002
        }
    ],
    "scores": {
        "score_1": [
            1,
            2,
            3
        ],
        "score_2": [
            4,
            5,
            30,
            6,
            24
        ]
    }
}
```

`out.json`
```json
{
    "id": 1,
    "name": "Ivan",
    "surname": "Ivanov",
    "salary_flag": 1,
    "address": {
        "land": "Russia",
        "city": "Moscow",
        "district": "Butovo"
    },
    "phones": [
        {
            "operator": "MTS",
            "number": 89850000001
        },
        {
            "operator": "МегаФон",
            "number": 89850000002
        }
    ],
    "scores": {
        "score_1": [
            1,
            2,
            3
        ],
        "score_2": [
            4,
            5,
            30,
            6,
            24
        ]
    }
}
```

# Интерфейсы

**Интерфейс** = контракт, в котором указываются методы, но не указываются поля.

Структуры `Square` и `Circle` имплементируют интерфейс `Shape`, т.к. каждая имеет метод с такой же сигнатурой.

```go
package main

import (
    "fmt"
    "math"
)

// Shape - интерфейс фигуры.
type Shape interface {
    Area() float64
}

//------------------------------

// Square - структура квадрата.
type Square struct {
    side float64
}

// NewSquare - конструктор, возвращающий тип интерфейса.
func NewSquare(side float64) Shape { // возвращвет интерфейс, а не Square!
    return &Square{side: side}
}

// Area - метод вычисления площади квадрата.
func (s *Square) Area() float64 {
    return s.side * s.side
}

//------------------------------

// Circle - структура круга.
type Circle struct {
    radius float64
}

// NewCircle - конструктор, возвращающий тип интерфейса.
func NewCircle(radius float64) Shape { // возвращвет интерфейс, а не Circle!
    return &Circle{radius: radius}
}

// Area - метод вычисления площади круга.
func (c *Circle) Area() float64 {
    return c.radius * c.radius * math.Pi
}

//------------------------------

// ShapeArea - универсальная функция вычисления площади.
func ShapeArea(shape Shape) float64 {
    return shape.Area()
}

func main() {
    // создание объектов интерфейса Shape
    square := NewSquare(123)
    circle := NewCircle(33)
    
    fmt.Println(square.Area()) // 15129
    fmt.Println(circle.Area()) // 3421.194399759285
    
    // вызов методов вычисления площади
    fmt.Println(ShapeArea(square)) // 15129
    fmt.Println(ShapeArea(circle)) // 3421.194399759285
    
    // слайс объектов интерфейса
    slOfShapes := []Shape{square, circle} // = полиморфизм (список с разными типами структур)
    
    for _, shape := range slOfShapes {
        fmt.Println(shape.Area())
    }
    // 15129
    // 3421.194399759285
    
    for _, shape := range slOfShapes {
        fmt.Println(ShapeArea(shape))
    }
    // 15129
    // 3421.194399759285
}
```

## Приведение типов. Пустой интерфейс

```go
package main

import (
    "fmt"
)

// Runner - интерфейс.
type Runner interface {
    Run()
}

type Dog struct {
    Name string
}

func (d *Dog) Run() {
    fmt.Println("Dog runs")
}

type Cat struct {
    Name string
}

func (c *Cat) Run() {
    fmt.Println("Cat runs")
}

// TypeAssertion - функция приведения типов.
// Можно использовать тип runner - Runner, но тогда не получится проверить типы int и string
// так как они не реализуют интерфейс Runner.
// Используем универсальное решение - пустой интерфейс, который реализуют ВСË, но это черевато ошибками и ресурсом
func TypeAssertion(runner interface{}) {
    switch v := runner.(type) {
    case *Dog:
        v.Run()
    case *Cat:
        v.Run()
    case int:
        fmt.Println("это int")
    case string:
        fmt.Println("это string")
    default:
        fmt.Println("незнакомый тип")
    }
}

func main() {
    // создание объекта интерфейса Runner
    var r Runner
    fmt.Printf("Type: %T Value: %#v\n", r, r) // Type: <nil> Value: <nil>
    
    // инициализация структуры Dog через интерфейс, который она реализует
    r = &Dog{Name: "dog"}
    fmt.Printf("Type: %T Value: %#v\n", r, r) // Type: *t.Dog Value: &t.Dog{Name:"dog"}
    TypeAssertion(r)
    
    // инициализация структуры Cat через интерфейс, который она реализует
    r = &Cat{Name: "cat"}
    fmt.Printf("Type: %T Value: %#v\n", r, r) // Type: *t.Cat Value: &t.Cat{Name:"cat"}
    TypeAssertion(r)
    
    // проверка остальных типов
    TypeAssertion(111)   // это int
    TypeAssertion("111") // это string
    TypeAssertion(true)  // незнакомый тип
}
```

# Наследование = Встраивание (Embedding). 

## Пример 1

```go
package main

import "fmt"

// Parent - родительская структура.
type Parent struct {
    ParentField int
}

// Hello - Родительский метод.
func (p *Parent) Hello() {
    fmt.Println("Parent func 'hello'")
}

// Child - Структура наследника.
type Child struct {
    Parent     // встраиваем
    ChildField int
}

// Hello - Метод наследника.
func (c *Child) Hello() {
    fmt.Println("Child func 'hello'")
}

func main() {
    ch := Child{
        Parent:     Parent{ParentField: 1111},
        ChildField: 2222,
    }
    
    fmt.Println(ch.Parent.ParentField) // 1111
    ch.Hello()                         // Child func 'hello'
    ch.Parent.Hello()                  // Parent func 'hello'
}
```

## Пример 2. Pattern Builder

```go
package main

import (
    "fmt"
)

// Decider - интерфейс выбора пути рассчета.
type Decider interface {
    LoanCalc()
}

// Customer - клиент.
type Customer struct {
    Name string
    Age  int
}

// WalkIn - клиент улица.
type WalkIn struct {
    Customer
}

// LoanCalc - расчет кредита для улицы.
func (w WalkIn) LoanCalc() {
    fmt.Println("Расчет для улицы")
}

// Salary - клиент зп.
type Salary struct {
    Customer
}

// LoanCalc - расчет кредита для зарплатника.
func (s Salary) LoanCalc() {
    fmt.Println("Расчет для зарплатника")
}

// Pension - клиент пенсионер.
type Pension struct {
    Customer
}

// LoanCalc - расчет кредита для пенсионера.
func (p Pension) LoanCalc() {
    fmt.Println("Расчет для пенсионера")
}

// Decision - структура со встроенным интерфейсом.
type Decision struct {
    Strategy string
    Decider
}

// New - универсальный конструктор для создания всех типов клиентов.
func New(strategy string, customer Decider) *Decision {
    return &Decision{
        Strategy: strategy,
        Decider:  customer,
    }
}

func main() {
    // Встраивание. Клиент Улица.
    walkInCustomer := New(
        "Улица",
        WalkIn{
            Customer: Customer{
                Name: "walkin customer",
                Age:  30,
            },
        },
    )
    walkInCustomer.LoanCalc()
    
    // Встраивание. Клиент Зарплатник
    salaryCustomer := New(
        "Зарплата",
        Salary{
            Customer: Customer{
                Name: "salary customer",
                Age:  40,
            },
        },
    )
    salaryCustomer.LoanCalc()
    
    // Встраивание. Клиент пенсионер
    pensCustomer := New(
        "Пенсия",
        Pension{
            Customer: Customer{
                Name: "pension customer",
                Age:  70,
            },
        },
    )
    pensCustomer.LoanCalc()
}
```

## Пример 2. Pattern Builder. Правильная иерархия

`choicer.go`
```go
package usecase4

type Strategy struct {
	WalkIn  WalkIn
	Salary  Salary
	Pension Pension
}

func NewStrategy() *Strategy {
	return &Strategy{
		WalkIn:  NewWalkIn(),
		Salary:  NewSalary(),
		Pension: NewPension(),
	}
}
```

`customer.go`
```go
package usecase4

// Customer - клиент.
type Customer struct {
	Name string
	Age  int
}
```

`pension.go`
```go
package usecase4

import "fmt"

type Pension interface {
	PrintPension()
	GetPensionParams() (string, int)
}

type PensionStruct struct {
	Customer
}

func NewPension() Pension {
	return &PensionStruct{
		Customer{
			Name: "Pension",
			Age:  70,
		},
	}
}

func (s *PensionStruct) PrintPension() {
	fmt.Println("This is PensionStruct method")
}

func (s *PensionStruct) GetPensionParams() (string, int) {
	fmt.Println("Pension params")
	return s.Name, s.Age
}
```

`salary.go`
```go
package usecase4

import "fmt"

type Salary interface {
	PrintSalary()
	GetSalaryParams() (string, int)
}

type SalaryStruct struct {
	Customer
}

func NewSalary() Salary {
	return &SalaryStruct{
		Customer{
			Name: "Salary",
			Age:  40,
		},
	}
}

func (s *SalaryStruct) PrintSalary() {
	fmt.Println("This is SalaryStruct method")
}

func (s *SalaryStruct) GetSalaryParams() (string, int) {
	fmt.Println("Salary params")
	return s.Name, s.Age
}
```

`walkin.go`
```go
package usecase4

import "fmt"

type WalkIn interface {
	PrintWalkIn()
	GetWalkInParams() (string, int)
}

type WalkInStruct struct {
	Customer
}

func NewWalkIn() WalkIn {
	return &WalkInStruct{
		Customer{
			Name: "WalkIn",
			Age:  20,
		},
	}
}

func (w *WalkInStruct) PrintWalkIn() {
	fmt.Println("This is WalkInStruct method")
}

func (w *WalkInStruct) GetWalkInParams() (string, int) {
	fmt.Println("WalkIn params")
	return w.Name, w.Age
}
```

# Подробнее про ошибки

`error` = интерфейс с 1 методом:
```go
type error interface {
    Error() string
}
```

```go
// CustomError - структура кастомной ошибки.
type CustomError struct {
	Err string
}

// Error - метод, реализующий интерфейс Error.
func (c *CustomErr) Error() string {
	return c.Err
}

func main() {
    err1 := CustomError{
        Err: "error1",
    }
    res := err1.Error()
    fmt.Println(res) // error1

    err2 := errors.New("error2")
    res := err2.Error()
    fmt.Println(res) // error2
}
```

# Словари = Карты = map

map = хэш-таблицв, где есть ключ-значение.
1. берется хэш от ключа
2. считается отстаток от деления хэша пункта 1. на количество бакетов
3. по полученному значени. пункта 2. как по индексу массива кладется значение

Колизия = ситуация, при которой хэш от разных значений совпадает => проблема => 
- метод цепочек (бакет = связный список, т.е. массив, который не последователен в ячейках памяти)
- метод открытой адресации (кладется колизия в пустой бакет справа)

В Go используется метод цепочек. В каждом бакете должно быть около 6.5 значений. При увеличении данного значения каждый раз при обращении в хэш-таблицу автоматически аллоцируется новая память и таким образом потихоньку переписыватеся в новую хэш-таблицу.

## Инициализация Map (словаря)
```go
package main

import "fmt"

func main() {
    // Вариант с make
    dct1 := make(map[string]int)
    fmt.Printf("dct1:%v, len:%v\n", dct1, len(dct1)) // dct1:map[], len:0
    
    // 5 = вместимость (cap), как у slice, но такого параметра нет! len = 0
    dct2 := make(map[string]int, 5)
    fmt.Printf("dct1:%v, len:%v\n", dct2, len(dct2)) // dct1:map[], len:0
    
    // Инициализация со значениями
    dct3 := map[string]int{
        "Ivanov": 1,
        "Petrov": 2,
    }
    fmt.Printf("dct3:%v, len:%v\n", dct3, len(dct3)) // dct3:map[Ivanov:1 Petrov:2], len:2

    dct4 := map[string]float64{} // фигурные скобки обязательны!
    fmt.Printf("dct4:%v, len:%v\n", dct4, len(dct4)) // dct4:map[], len:0
    dct4["a"] = 1.0
    fmt.Printf("dct4:%v, len:%v\n", dct4, len(dct4)) // dct4:map[a:1.0], len:1

    var dct5 map[string]float64 // нет фигурных скобок => бестолковый способ! т.к. см. ниже
    fmt.Printf("dct5:%v, len:%v\n", dct5, len(dct5)) // dct5:map[], len:0
    dct5["a"] = 1.0 // panic    
}
```

## Получение, добавление, удаление 

```go
package main

import (
    "fmt"
    "log"
)

func main() {
    dct := make(map[string]int)
    fmt.Printf("dct:%v\n", dct) // dct:map[]
    
    // вставка значения по ключу
    dct["Ivanov"] = 1
    dct["Petrov"] = 2
    fmt.Printf("dct:%v\n", dct) // dct:map[Ivanov:1 Petrov:2]
    
    // НЕБЕЗОПАСНОЕ получение значения по ключу
    getFromMap := dct["Petrov"]
    fmt.Println(getFromMap) // 2
    // Если в мапе нет значения по ключу, то вернется дефолтное значение типа
    wtf := dct["Sidorov"]
    fmt.Println(wtf) // вернется дефольное значение типа int = 0 !!!

    // Безопасное получение значения по ключу
    value1, ok := dct["Petrov"] // в bool переменную 'ok' будет записан true, в случае находжения ключа
    if !ok {
        log.Fatalln("пользователя нет в словаре")
    }
    fmt.Println(value1) // 2
    value2, ok := dct["Who are U?"]
    if !ok {
        log.Fatalln("пользователя нет в словаре")
    }
    fmt.Println(value2) // пользователя нет в словаре
    
    // изменение значения по ключу
    dct["Ivanov"] = 1_000
    fmt.Printf("dct:%v\n", dct) // dct:map[Ivanov:1000 Petrov:2]
    
    // удаление значения по ключу
    delete(dct, "Petrov")
    fmt.Printf("dct:%v\n", dct) // dct:map[Ivanov:1000]
}
```

## Итерирование
```go
package main

import (
    "fmt"   
)

func main() {
    dct := map[string]int{
        "Ivanov": 1,
        "Petrov": 2,
    }
    
    for k, v := range dct {
        fmt.Printf("key: %v, value: %v\n", k, v)
    }
    // key: Ivanov, value: 1
    // key: Petrov, value: 2
    // порядок НЕ гаратнирован!!!
}
```

## Использование в map разных типов данных 
Тип пустой интерфейс (`interface{}`) дает возможность записывать в map данные разных типов.
Но, появляется необходимость приведения типа значения при работе с конкретной парой, что создает неудобство.
Для разных типов данных принято использовать **структуры**.

```go
package main

import (
    "fmt"
)

func main() {
    // Тип пустой интерфейс дает возможность записывать в мапу данные разных типов
    dct := map[string]interface{}{
        "Name":    "Ivan",
        "Surname": "Ivanov",
        "Age":     30,
        "Married": true,
    }
    
    for k, v := range dct {
        fmt.Println(k, v)
    }
    
    // неудобство приведения типа
    if dct["Age"].(int) == 30 {
        fmt.Println("ok")
    }
}
```

```go
package main

import (
	"fmt"
)

func main() {
    // мапа с разными типами значений
    dct := map[string]interface{}{
        "Name":    "Ivan",
        "Surname": "Ivanov",
        "Age":     30,
        "Married": true,
    }
    
    // приведение типа
    for _, v := range dct {
        switch v.(type) {
        case int:
            fmt.Printf("type: int, val: %d\n", v)
        case float64:
            fmt.Printf("type: float, val: %f\n", v)
        case string:
            fmt.Printf("type: string, val: %s\n", v)
        default:
            fmt.Printf("незнакомый тип, val: %v\n", v)
        }
    }
}
```

## usecase 1. Unique

```go
package usecase1

type User struct {
	Id   int
	Name string
}

// UniqueUsers - Фильтрации по уникальным именам.
func UniqueUsers() map[int]User {
	users := []User{
		{Id: 1, Name: "Sasha"},
		{Id: 222, Name: "Petya"},
		{Id: 3, Name: "Andrey"},
		{Id: 222, Name: "Petya"},
		{Id: 4, Name: "Sergei"},
	}

	uniqueUsers := make(map[int]User)

	for _, user := range users {
		if _, ok := uniqueUsers[user.Id]; !ok {
			uniqueUsers[user.Id] = user
		}
	}

	// map[1:{1 Sasha} 3:{3 Andrey} 4:{4 Sergei} 222:{222 Petya}]
	return uniqueUsers
}
```

## usecase 2.

```go
package usecase2

// Справочник
var pilotCategories = []string{"walkIn", "salary"}

type Customer struct {
	Name  string
	Score int
}

// MostScore - Поиск клиента с максимальным скором.
func MostScore(data map[string][]Customer) Customer {
	// клиент с максимальным скором
	var mostScoreCustomer Customer

	for _, category := range pilotCategories {
		// проверка на категорию
		customers, ok := data[category]
		if !ok { // перебор только интересующих категорий
			continue
		}
		// поиск клиента с максимальным скором
		for _, customer := range customers {
			if customer.Score > mostScoreCustomer.Score {
				mostScoreCustomer = customer
			}
		}
	}

	return mostScoreCustomer
}
```

## usecase 3.

```json
{
    "id": 1,
    "name": "Ivan",
    "surname": "Ivanov",
    "address": {
        "land": "Russia",
        "city": "Moscow",
        "district": "Butovo"
    },
    "phones": [
        {
            "operator": "MTS",
            "number": 89850000001
        },
        {
            "operator": "МегаФон",
            "number": 89850000001
        }
    ],
    "scores": {
        "score_1": [
            1,
            2,
            3
        ],
        "score_2": [
            4,
            5,
            0,
            6,
            0
        ]
    }
}
```

```go
package usecase3

import (
	"encoding/json"
	"fmt"
	"os"
)

// JSONParse - парсинг json в map.
func JSONParse() {
	// чтение файла
	payload, err := os.ReadFile("usecase3/data/in.json")
	if err != nil {
		fmt.Println(err)
	}
	dct := make(map[string]interface{}) // инициализация мапы
	if err := json.Unmarshal(payload, &dct); err != nil { // парсинг
		fmt.Println(err)
	}
	fmt.Println(dct)
	fmt.Println(dct["name"].(string))
    fmt.Println(dct["phones"].([]interface{})[0]) // это костыль!!
	fmt.Println(dct["phones"][0]) // ошибка. Пустой интерфейс!
    // поэтому парсим ТОЛЬКО в структуры!
}
```


