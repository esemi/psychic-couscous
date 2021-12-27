package domain

type (
	MovieID string

	Movie struct {
		ID             MovieID
		TitleType      string
		PrimaryTitle   string
		OriginalTitle  string
		Adult          bool
		StartYear      int16
		EndYear        int16
		RuntimeMinutes int16
		Genres         []string
	}
)

type MovieRepository interface {
	GetByID(id MovieID) (Movie, error)
	GetByName(name string) (Movie, error)
	Save(name Movie) error
	LoadFromCSV(filename string) error
}
