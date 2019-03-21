package database

type CountryDB interface {
	Get(name string) (Country, error)
	Set(name string, data Country) error
	Del(name string) (Country, error)
}
