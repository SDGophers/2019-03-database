package database

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"testing"
)

func BenchmarkPopulate(b *testing.B) {
	data := map[string]Country{}

	for _, v := range []string{"pop.json", "geo.json"} {
		byt, err := ioutil.ReadFile(v)
		if err != nil {
			panic(err)
		}

		switch v {
		case "pop.json":
			jdata := []struct {
				Country, Population string
			}{}

			err = json.Unmarshal(byt, &jdata)
			if err != nil {
				panic(err)
			}

			for _, v := range jdata {
				vv := Country{}
				copy(vv.Name[:], []byte(v.Country))
				vv.Population, err = strconv.ParseUint(v.Population, 10, 64)
				if err != nil {
					vv.Population = 0
				}
			}
		case "geo.json":
			jdata := []struct {
				Country                  string
				North, South, East, West string
			}{}

			err = json.Unmarshal(byt, &jdata)
			if err != nil {
				panic(err)
			}

			for _, v := range jdata {
				vv := data[v.Country]
				vv.North, err = strconv.ParseFloat(v.North, 64)
				vv.South, err = strconv.ParseFloat(v.South, 64)
				vv.East, err = strconv.ParseFloat(v.East, 64)
				vv.West, err = strconv.ParseFloat(v.West, 64)

				data[v.Country] = vv
			}
		}
	}

	db, err := NewCountryDBImpl("test.bin")
	if err != nil {
		panic(err)
	}

	for n := 0; n < b.N; n++ {
		for k, v := range data {
			err = db.Set(k, v)
			if err != nil {
				panic(err)
			}
		}

		for k := range data {
			_, err := db.Del(k)
			if err != nil {
				panic(err)
			}
		}
	}

}
