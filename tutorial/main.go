// https://golangify.com

// Объявляет пакет, которому принадлежит код
package main

import (
	"fmt" // Делает пакет fmt (формат) доступным для использования
	"math/rand"
	"strings"
	"time"
)

func f1() {
	// Выводит текст Hello, playground на экран
	fmt.Println("Hello, playground")
}

/*
func f2() // отсутствует тело функции
{ // ошибка синтаксиса: лишняя точка с запятой или
} // новая строка перед {
*/

func f3() {
	fmt.Println("Hello, Nathan")
	fmt.Println("こんにちは Здравствуйте Hola")
}

func f4() {
	fmt.Print("Мой вес на поверхности Марса равен ")
	fmt.Print(55.0 * 0.3783) // В результате 20.8065
	fmt.Println(" килограммам, а мой возраст равен", 41*365/687, "годам.")
}

func f5() {
	// Выводит: Мой вес на поверхности Марса равен 20.8065 килограммам,
	fmt.Printf("Мой вес на поверхности Марса равен %v килограммам, ", 55.0*0.3783)
	// Выводит: а мой возраст равен 21 годам.
	fmt.Printf("а мой возраст равен %v годам.\n", 41*365/687)
}

func f6() {
	// Выводит: Мой вес на поверхности Земли равен 55 килограммам.
	fmt.Printf("Мой вес на поверхности %v равен %v килограммам.\n", "Земли", 55)
}

func f7() {
	fmt.Printf("%-15v $%4v\n", "SpaceX", 94)
	fmt.Printf("%-15v $%4v\n", "Virgin Galactic", 100)
}

func f8() {
	const lightSpeed = 299792 // км/с
	var distance = 56000000   // км

	fmt.Println(distance/lightSpeed, "секунд") // В результате 186 секунд

	distance = 401000000
	fmt.Println(distance/lightSpeed, "секунд") // В результате 1337 секунд
}

func f9() {
	var distance1 = "1"
	var speed1 = 1

	fmt.Printf("%v %v", distance1, speed1)

	var (
		distance2 = "2"
		speed2    = 2
	)

	fmt.Printf("%v %v", distance2, speed2)

	var distance3, speed3 = "3", 3

	fmt.Printf("%v %v", distance3, speed3)

}

func f10() {
	var weight = 149.0
	fmt.Println(weight)
	weight = weight * 0.3783
	fmt.Println(weight)
	weight *= 0.3783
	fmt.Println(weight)

	var age = 41
	fmt.Println(age)
	age = age + 1 // С днем рождения!
	fmt.Println(age)
	age += 1
	fmt.Println(age)
	age++
	fmt.Println(age)
}

func f11() {
	var num = rand.Intn(10) + 1
	fmt.Println(num)

	num = rand.Intn(10) + 1
	fmt.Println(num)
}

/*
Расстояние между Землей и Марсом в разное время отличается и зависит от того,
где планеты в данный конкретный момент времени находятся на орбите Солнца.
Напишите программу для генерации случайного расстояния в промежутке от 56 000 000 до 401 000 000 км.
*/
func f12() {
	var distance = rand.Intn(401_000_000-56_000_000) + 56_000_000
	fmt.Println(distance)
}

/*
Напишите программу, которая посчитает, как быстро должна передвигаться ракета (км/ч),
чтобы добраться до Марса за 28 дней.
Предположим, что расстояние от Земли до Марса равно 56 000 000 км.
*/
func f13() {
	const hoursPerDay = 24

	var days = 28
	var distance = 56_000_000 // km

	fmt.Println(distance/(days*hoursPerDay), "км/ч")
}

func f14() {
	var (
		walkOutside     = true
		takeTheBluePill = false
	)
	fmt.Printf("%v %v", walkOutside, takeTheBluePill)
}

func f15() {
	fmt.Println("Вы находитесь в темной пещере.")

	var command = "выйти наружу"
	var exit = strings.Contains(command, "наружу")

	fmt.Println("Вы покидаете пещеру:", exit) // Выводит: Вы покидаете пещеру: true
}

/*
== равно
!= не равно
< меньше
> больше
<= меньше или равно
>= больше или равно
*/

func f16() {
	fmt.Println("На знаке снаружи написано 'Несовершеннолетним вход запрещен'.")

	var age = 41
	var adult = age >= 18

	fmt.Printf("В возрасте %v, я совершеннолетний? %v\n", age, adult)
}

func f17() {
	fmt.Println("яблоко" > "банан")
}

func f18() {
	var room = "пещера"

	if room == "пещера" {
		fmt.Println("Вы находитесь в тускло освещенной пещере.")
	} else if room == "вход" {
		fmt.Println("Здесь есть вход в пещеру и путь на восток.")
	} else if room == "гора" {
		fmt.Println("Здесь крутой утес. Тропа ведет к подножью горы.")
	} else {
		fmt.Println("Здесь ничего нет.")
	}
}

/*
else if и else являются опциональными.
Когда рассматривается несколько вариантов, можно повторять else if столько раз, сколько требуется.
*/

/*
Напишем код, что должен определить, будет ли 2100 год високосным.
Правила определения високосного года таковы:
- Любой год, что делится без остатка на четыре, но не делится без остатка на 100;
- Или любой год, что делится без остатка на 400.
*/
func f19() {
	fmt.Println("На дворе 2100 год. Он високосный?")

	var year = 2100
	var leap = year%400 == 0 || (year%4 == 0 && year%100 != 0)

	if leap {
		fmt.Println("Этот год високосный!")
	} else {
		fmt.Println("К сожалению, нет. Этот год не високосный.")
	}
}

func f20() {
	var haveTorch = true
	var litTorch = false

	if !haveTorch || !litTorch {
		fmt.Println("Ничего не видно.") // Вывод: Ничего не видно.
	}
}

func f21() {
	fmt.Println("Здесь вход в пещеру и путь на восток.")
	var command = "зайти внутрь"

	switch command { // Сравнивает case с command
	case "идти на восток":
		fmt.Println("Вы направляетесь к горе.")
	case "зайти в пещеру", "зайти внутрь": // Запятая разделяет список возможных значений
		fmt.Println("Вы находитесь в тускло освещенной пещере.")
	case "прочитать знак":
		fmt.Println("На знаке написано 'Несовершеннолетним вход запрещен'.")
	default:
		fmt.Println("Пока не совсем понятно.")
	}
}

func f22() {
	var room = "озеро"

	switch { // Выражения для каждого случая
	case room == "пещера":
		fmt.Println("Вы находитесь в тускло освещенной пещере.")
	case room == "озеро":
		fmt.Println("Лед кажется достаточно крепким.")
		fallthrough // Переходит на следующий случай бкз сравнения!
	case room == "глубина":
		fmt.Println("Вода такая холодная, что сводит кости.")
	}
}

func f23() {
	var count = 10 // Объявление и инициализация

	for count > 0 { // Условие
		fmt.Println(count)
		time.Sleep(time.Second)
		count-- // Обратный отсчет; в противном случае цикл будет длиться вечно
	}
	fmt.Println("Запуск!")
}

func f24() {
	var degrees = 0

	for {
		fmt.Println(degrees)

		degrees++
		if degrees >= 360 {
			degrees = 0
			if rand.Intn(10) == 0 {
				break
			}
		}
	}
}

/*
Не каждый запуск проходит по плану.
Реализуйте обратный отсчет, где на каждую секунду приходится шанс 1 к 100,
что ввиду определенных обстоятельств запуск прервется, и счетчик остановится.
*/
func f25() {
	var count = 10

	for count > 0 {
		fmt.Println(count)
		time.Sleep(time.Second)
		if rand.Intn(100) == 0 {
			break
		}
		count--

	}
	if count == 0 {
		fmt.Println("Запуск!")
	} else {
		fmt.Println("Запуск отменяется.")
	}
}

func f26() {
	switch time.Now().Weekday() {

	case time.Monday:
		fmt.Println("Сегодня понедельник.")

	case time.Tuesday:
		fmt.Println("Сегодня вторник.")

	case time.Wednesday:
		fmt.Println("Сегодня среда.")

	case time.Thursday:
		fmt.Println("Сегодня четверг.")

	case time.Friday:
		fmt.Println("Сегодня пятница.")

	case time.Saturday:
		fmt.Println("Сегодня суббота.")

	case time.Sunday:
		fmt.Println("Сегодня воскресенье.")
	}
}

func f27() {
	switch time.Now().Weekday() {

	case time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday:
		fmt.Println("будний день")
	case time.Saturday, time.Sunday:
		fmt.Println("выходные дни")
	}
}

func f28() {
	size := "XXXL"

	switch size {

	case "XXS":
		fmt.Println("очень очень маленький")

	case "XS":
		fmt.Println("очень маленький")

	case "S":
		fmt.Println("маленький")

	case "M":
		fmt.Println("средний")

	case "L":
		fmt.Println("большой")

	case "XL":
		fmt.Println("очень большой")

	case "XXL":
		fmt.Println("очень очень большой")

	default:
		fmt.Println("неизвестно")
	}
}

func f29() {
	switch num := 6; num%2 == 0 {

	case true:
		fmt.Println("even value")

	case false:
		fmt.Println("odd value")
	}
}

func f30() {
	w := "a b c\td\nefg hi"

	for _, e := range w {

		switch e {
		case ' ', '\t', '\n':
			break // заканчивает switch
		default:
			fmt.Printf("%c\n", e)
		}
	}
}

func f31() {

	now := time.Now()

	switch {
	case now.Hour() < 12:
		fmt.Println("AM")

	default:
		fmt.Println("PM")
	}
}

// A -> B -> C -> D -> E

func f32() {

	nextstop := "B"

	fmt.Println("Stops ahead of us:")

	switch nextstop {

	case "A":
		fmt.Println("A")
		fallthrough

	case "B":
		fmt.Println("B")
		fallthrough

	case "C":
		fmt.Println("C")
		fallthrough

	case "D":
		fmt.Println("D")
		fallthrough

	case "E":
		fmt.Println("E")
	}
}

func f33() {

	var data interface{}

	data = 112523652346.23463246345

	switch mytype := data.(type) {

	case string:
		fmt.Println("string")

	case bool:
		fmt.Println("boolean")

	case float64:
		fmt.Println("float64 type")

	case float32:
		fmt.Println("float32 type")

	case int:
		fmt.Println("int")

	default:
		fmt.Printf("%T", mytype)
	}
}

/*
В Go область видимости начинает и заканчивается фигурными скобками {}.
В следующей программе функция main начинает область видимости,
а вместе с циклом for стартует вложенная область.
*/
func f34() {
	var count = 0

	for count < 10 { // Начало области видимости
		var num = rand.Intn(10) + 1
		fmt.Println(num)

		count++
	} // Конец области видимости
}

/*
Переменная count объявляется внутри области видимости функции,
она видима до конца функции main,
в то время как переменная num объявляется внутри области видимости цикла for.
По завершении цикла переменная num выходит из области видимости.

При попытке получить доступ к переменной num после цикла, компилятор Go выведет ошибку.
Однако получить доступ к переменной count по завершении цикла for все еще можно,
ведь ее объявили за пределами цикла, хотя особой причины для этого не было.
Для заключения переменной count в области видимости цикла понадобится использовать другой способ объявления переменных Go.
*/

// краткое объявление переменных в go
func f35() {
	var count1 = 10
	fmt.Println(count1)
	// аналогичная запись
	count2 := 10
	fmt.Println(count2)
}

/*
Поначалу может показаться, что разница невелика,
однако разница в три символа делает сокращенный вариант намного популярнее способа с var.
Кроме того, краткое объявление может использоваться в некоторых местах, где недопустимо ключевое слово var.

В следующей программе показан пример цикла for, что совмещает инициализацию, условие и последующий оператор,
что уменьшает значение count.
При использовании данной формы цикла for очень важен порядок: инициализация, условие, операция.
*/

func f36() {
	var count = 0

	for count = 10; count > 0; count-- {
		fmt.Println(count)
	}

	fmt.Println(count) // count остается в области видимости
}

/*
Не используя краткое объявление, переменную count нужно было бы объявить за пределами цикла,
следовательно, переменная в таком случае после завершения цикла остается в области видимости.

В следующей программе при задействовании краткого объявления,
переменная count объявляется и инициализируется как часть цикла for,
по завершении цикла выходит из области видимости.
Если бы к переменной count доступ был получен за пределами цикла,
тогда компилятор Go выдал бы ошибку undefined: count.
*/

func f37() {
	for count := 10; count > 0; count-- {
		fmt.Println(count)
	} // count больше не в области видимости
}

/*
Краткое объявление дает возможность объявить новую переменную в операторе if.
В следующем коде переменная num может использовать в любом ответвлении оператора if.
*/

func f38() {
	if num := rand.Intn(3); num == 0 {
		fmt.Println("Space Adventures")
	} else if num == 1 {
		fmt.Println("SpaceX")
	} else {
		fmt.Println("Virgin Galactic")
	} //num больше не в области видимости
}

// Краткое объявление может использоваться с оператором switch, как показано в следующей программе:

func f39() {
	switch num := rand.Intn(10); num {
	case 0:
		fmt.Println("Space Adventures")
	case 1:
		fmt.Println("SpaceX")
	case 2:
		fmt.Println("Virgin Galactic")
	default:
		fmt.Println("Random spaceline #", num)
	}
}

//Локальная и глобальная область видимости

/*
Следующий код генерирует и отображает случайную дату — к примеру, дату вылета на Марс.
В нем также представлено несколько областей видимости и показано,
почему особенно важно задуматься об области видимости во время объявления переменной.
*/

var era = "AD" // переменная era доступна через пакет

/*
На заметку:
Краткое объявление недоступно для переменных, объявленных в области видимости пакета,
поэтому переменную era нельзя объявить через era := "AD" в ее текущей позиции.
*/

func f40() {
	year := 2018 // переменные era и year находятся в области видимости

	switch month := rand.Intn(12) + 1; month { // переменные era, year и month в области видимости
	case 2:
		day := rand.Intn(28) + 1 // новый день
		fmt.Println(era, year, month, day)
	case 4, 6, 9, 11:
		day := rand.Intn(30) + 1
		fmt.Println(era, year, month, day)
	default:
		day := rand.Intn(31) + 1
		fmt.Println(era, year, month, day)
	} // month и day за пределами области видимости
} // year за пределами области видимости

/*
Переменная year видна только внутри функции main.
Другие функции видят era, но не year.
Область видимости функции уже, чем область видимости пакета.
Она начинается c ключевого слова func и заканчивается закрывающей скобкой.
*/

/*
Переменная month доступна внутри оператора switch,
но как только оператор switch заканчивается, month выводится из области видимости.
Область видимости начинается с ключевого слова switch и заканчивается закрывающей скобкой switch.
*/

/*
У каждого case есть собственная область видимости и три независимые переменные day.
Как только каждый случай заканчивается, объявленная внутри case переменная day выходит за пределы области видимости.
Это единственная ситуация, когда для обозначения области видимости не используются скобки.
*/

/*
Код примера не идеален.
Узкая область видимости переменных month и day приводит к дубликату кода (Println, Println, Println).
Когда код дублируется, кто-то может пересмотреть код в одной области,
но не в другой (при решений не выводить era, но забыв изменить один case).
Иногда дублированный код имеет смысл,
однако чаше всего он рассматривается как код с запашком и указывается на возможные проблемы в программе.
*/

/*
Для удаления дубликатов и упрощения кода переменные нужно было объявлять в более широкой области видимости функции, делая их доступными после оператора switch для дальнейшей работы.
*/

/*
Время рефакторинга!
Рефакторинг предполагает модификацию кода без изменения его поведения.
Следующая программа по-прежнему выводит случайную дату.
*/

func f41() {
	year := 2018
	month := rand.Intn(12) + 1
	daysInMonth := 31

	switch month {
	case 2:
		daysInMonth = 28
	case 4, 6, 9, 11:
		daysInMonth = 30
	}

	day := rand.Intn(daysInMonth) + 1
	fmt.Println(era, year, month, day)

}

/*
- Открывающая фигурная скобка { вводит новую область видимости, что оканчивается закрывающей скобкой };
- Ключевые слова case и default также вводят новую область видимости, хотя фигурные скобки здесь уже не используются;
- Место объявления переменной определяется тем, в какой она области видимости;
- Переменные, объявленные на той же строке, что и ключевые слова for, if или switch находятся в области видимости до окончания данного оператора;
- Иногда широкая область видимости лучше, а иногда — узкая, все зависит от ситуации.
*/

/*
Измените следующую программу для обработки високосных годов. Код должен:

- Генерировать случайный год вместо постоянного использования 2018;
- Для февраля присвойте daysInMonth на 29 для високосных годов, и 28 для всех остальных. Можете использовать оператор if вместо блока case;
- Используйте цикл for для генерации и отображения 10 случайных дат.
*/

func f42() {
	for count := 0; count < 10; count++ {
		year := 2018 + rand.Intn(10)
		leap := year%400 == 0 || (year%4 == 0 && year%100 != 0)
		month := rand.Intn(12) + 1

		daysInMonth := 31
		switch month {
		case 2:
			daysInMonth = 28
			if leap {
				daysInMonth = 29
			}
		case 4, 6, 9, 11:
			daysInMonth = 30
		}

		day := rand.Intn(daysInMonth) + 1
		fmt.Println(era, year, month, day)
	}
}

//Создание программы для покупки билетов в Golang

/*
Пришло время проверить свои силы. Напишем в Go Playground программу для покупки билетов для путешествия на Марс. В коде используем переменные, константы, switch, if и for.
Для отображения, выравнивания текста и генерации случайных чисел будут задействованы пакеты fmt и math/rand.
*/

/*
Создание программы для покупки билетов в Golang
Примеры кода программ на Golang
Пришло время проверить свои силы. Напишем в Go Playground программу для покупки билетов для путешествия на Марс. В коде используем переменные, константы, switch, if и for. Для отображения, выравнивания текста и генерации случайных чисел будут задействованы пакеты fmt и math/rand.

При планировании поездки на Марс будет удобно собрать расценки различных космических станций в одном месте.
Есть множество сайтов для авиалиний, но не для космических.
Для нас это не будет проблемой. При умелом руководстве, Go сможет решить проблемы подобного рода.
*/

/*
В таблице четыре столбца:

- Космическая станция (Spaceline), что предоставляет услуги;
- Продолжительность (Duration) в днях поездки на Марс в один конец;
- Покрывает ли цена поездку туда и обратно (Trip type);
- Цена (Price) в миллионах долларов.
*/

/*
Для каждого билета случайным образом выбирается космическая станция: Space Adventures, SpaceX или Virgin Galactic.
*/

/*
Датой отправления на каждом билете значится 13 Октября 2020 года. В этот день Марс будет на расстоянии 62 100 000 км от Земли.

Скорость космического корабля будет выбрана случайным образом из диапазона от 16 до 30 км/ч.
Это определит продолжительность поездки на Марс, а также цену билета.
Более быстрые корабли намного дороже. Цены на билеты варьируются от $36 до $50 миллионов.
Цена для поездки туда-обратно удваивается.
*/

const secondsPerDay = 86400

func f43() {
	distance := 62_100_000
	company := ""
	trip := ""

	fmt.Println("Spaceline        Days Trip type  Price")
	fmt.Println("======================================")

	for count := 0; count < 10; count++ {
		switch rand.Intn(3) {
		case 0:
			company = "Space Adventures"
		case 1:
			company = "SpaceX"
		case 2:
			company = "Virgin Galactic"
		}

		speed := rand.Intn(15) + 16                  // 16-30 km/s
		duration := distance / speed / secondsPerDay // days
		price := 20.0 + speed                        // millions

		if rand.Intn(2) == 1 {
			trip = "Round-trip"
			price = price * 2
		} else {
			trip = "One-way"
		}

		fmt.Printf("%-16v %4v %-10v $%4v\n", company, duration, trip, price)
	}
}

//Вещественные числа в Golang — float64 и float32

// main является функцией, с которой все начинается
func main() {
	f43()
}
