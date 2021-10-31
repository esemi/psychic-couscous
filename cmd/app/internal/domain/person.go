package domain

type (
	PersonID string

	Person struct {
		ID                 PersonID
		PrimaryName        string
		BirthYear          int16
		DeathYear          int16
		PrimaryProfessions []string
		KnownForTitles     []TitleID
	}
)

type PersonRepository interface {
	GetByID(id PersonID) (Person, error)
	GetByName(name string) (Person, error)
	Save(name Person) error
	LoadFromCSV(filename string) error
}
