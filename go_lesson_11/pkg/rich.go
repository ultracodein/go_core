package iface

// RichEmployee - сотрудник
type RichEmployee struct {
	age int
}

// Age возвращает возраст сотрудника
func (e *RichEmployee) Age() int {
	return e.age
}

// NewRichEmployee возвращает ссылку на созданный объект сотрудника
func NewRichEmployee(age int) *RichEmployee {
	return &RichEmployee{
		age: age,
	}
}

// RichCustomer - клиент
type RichCustomer struct {
	age int
}

// NewRichCustomer возвращает ссылку на созданный объект клиента
func NewRichCustomer(age int) *RichCustomer {
	return &RichCustomer{
		age: age,
	}
}

// Age возвращает возраст клиента
func (c *RichCustomer) Age() int {
	return c.age
}

// AgeReporter - абстракция объекта, который может сообщить свой возраст
type AgeReporter interface {
	Age() int
}

// MaxRichPersonAge возвращает возраст старшего человека
func MaxRichPersonAge(persons ...AgeReporter) int {
	maxAge := 0
	for _, person := range persons {
		age := person.Age()
		if age > maxAge {
			maxAge = age
		}
	}
	return maxAge
}
