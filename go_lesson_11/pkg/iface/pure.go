package iface

import (
	"errors"
)

// Employee - сотрудник
type Employee struct {
	age int
}

// NewEmployee возвращает ссылку на созданный объект сотрудника
func NewEmployee(age int) *Employee {
	return &Employee{
		age: age,
	}
}

// Customer - клиент
type Customer struct {
	age int
}

// NewCustomer возвращает ссылку на созданный объект клиента
func NewCustomer(age int) *Customer {
	return &Customer{
		age: age,
	}
}

// MaxPersonAge возвращает объект старшего человека
func MaxPersonAge(persons ...interface{}) (interface{}, error) {
	maxAge := 0
	var oldestPerson interface{}
	for _, person := range persons {
		age := 0
		switch person.(type) {
		case *Employee:
			age = person.(*Employee).age
		case *Customer:
			age = person.(*Customer).age
		default:
			return nil, errors.New("arg type is not supported")
		}
		if age > maxAge {
			maxAge = age
			oldestPerson = person
		}
	}
	return oldestPerson, nil
}
