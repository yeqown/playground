package gomock

// IAnimal .
type IAnimal interface {
	Quack(times int) error
}

// AnimalQuack .
func AnimalQuack(animal IAnimal, times int) error {
	return animal.Quack(times)
}

// 注意开头没有空格
//go:generate mockgen -source ./mock.go -destination mock_generated_mock.go -package gomock
