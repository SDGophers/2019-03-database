package database

import "os"

func NewCountryDBImpl(file string) (CountryDB, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	return &CountryDBImpl{file: f}, nil
}

// A list implementation of the country DB
type CountryDBImpl struct {
	file *os.File
}

func (db *CountryDBImpl) Get(name string) (Country, error) {
	return Country{}, nil
}

func (db *CountryDBImpl) Set(name string, data Country) error {
	return nil
}

func (db *CountryDBImpl) Del(name string) (Country, error) {
	return Country{}, nil
}
