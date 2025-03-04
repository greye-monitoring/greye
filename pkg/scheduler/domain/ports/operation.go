package ports

type Operation interface {
	Add()
	Delete()
	Update()

	GetById()
	GetAll()
}
