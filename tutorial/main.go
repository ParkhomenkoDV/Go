// Объявляет пакет, которому принадлежит код
package main

import (
	// Делает пакет fmt (формат) доступным для использования
	"fmt"
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

// main является функцией, с которой все начинается
func main() {
	f10()
}
