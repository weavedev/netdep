package servicecalls

type TESTService interface {
	FirstMethod(a int, b int, c int) error
	SecondMethod(a int, b int, c int) error
	ThirdMethod(a int, b int, c int) error
}
