package infra

import "fmt"

type FakeDB struct{}

func (f *FakeDB) Create(data any) {
	fmt.Println("DB CREATE:", data)
}
