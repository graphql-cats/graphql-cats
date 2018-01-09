package main

import (
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
)

type Universe struct {
	Films     Films
	People    People
	Planets   Planets
	Species   Species
	Starships Starships
	Vehicles  Vehicles
}

type JsonUniverse struct {
	Films     map[string]Film
	People    map[string]Person
	Planets   map[string]Planet
	Species   map[string]Specie
	Starships map[string]Starship
	Vehicles  map[string]Vehicle
}

func main() {
	client := &http.Client{}
	root, err := fetchRoot(client)
	if err != nil {
		fmt.Errorf("%v\n", err)
		return
	}

	var universe Universe

	// could make this parallel but feels a little douchy to hammer
	// their API. Could probably clean this up some too
	entities := []struct {
		basename string
		data     Appendable
		result   Resultable
	}{
		{"films", &universe.Films, &FilmResult{}},
		{"people", &universe.People, &PersonResult{}},
		{"planets", &universe.Planets, &PlanetResult{}},
		{"species", &universe.Species, &SpeciesResult{}},
		{"starships", &universe.Starships, &StarshipResult{}},
		{"vehicles", &universe.Vehicles, &VehicleResult{}},
	}

	for _, v := range entities {
		filename := v.basename + ".gob"
		err = fetch(client, root[v.basename], filename, v.data, v.result)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	f, err := os.Create("swapi.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()

	m := s2m(&universe)

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err = enc.Encode(&m); err != nil {
		fmt.Println(err)
		return
	}
}

func s2m(u *Universe) *JsonUniverse {
	var ju JsonUniverse
	// should probably handle dup keys but can't be faffed
	ju.Films = make(map[string]Film)
	for _, v := range u.Films {
		id := string(v.Url)
		v.Id = id
		v.Url = "" // drop the URL in the JSON
		ju.Films[id] = v
	}

	ju.People = make(map[string]Person)
	for _, v := range u.People {
		id := string(v.Url)
		v.Id = id
		v.Url = ""
		ju.People[id] = v
	}

	ju.Planets = make(map[string]Planet)
	for _, v := range u.Planets {
		id := string(v.Url)
		v.Id = id
		v.Url = ""
		ju.Planets[id] = v
	}

	ju.Species = make(map[string]Specie)
	for _, v := range u.Species {
		id := string(v.Url)
		v.Id = id
		v.Url = ""
		ju.Species[id] = v
	}

	ju.Starships = make(map[string]Starship)
	for _, v := range u.Starships {
		id := string(v.Url)
		v.Id = id
		v.Url = ""
		ju.Starships[id] = v
	}

	ju.Vehicles = make(map[string]Vehicle)
	for _, v := range u.Vehicles {
		id := string(v.Url)
		v.Id = id
		v.Url = ""
		ju.Vehicles[id] = v
	}

	return &ju
}

func url2id(url string) string {
	last := len(url) - 1
	pos := strings.LastIndex(url[:last], "/")
	return url[pos+1 : last]
}

func open(filename string) (*os.File, error) {
	return os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_SYNC, 0666)
}

var ErrEmpty = errors.New("Empty file contents")

func decode(f *os.File, e interface{}) error {
	stat, err := f.Stat()
	if err != nil {
		return err
	}

	if stat.Size() > 0 {
		if err = gob.NewDecoder(f).Decode(e); err != nil {
			return err
		}

		return nil
	}

	return ErrEmpty
}

// fetch the API root or if cached use that.
func fetchRoot(c *http.Client) (map[string]string, error) {
	f, err := open("root.gob")
	if err != nil {
		return nil, err
	}
	defer f.Close()

	m := make(map[string]string)
	err = decode(f, &m)
	if err != nil && err != ErrEmpty {
		return nil, err
	} else if err == nil {
		fmt.Printf("read %v entries from root.gob\n", len(m))
		return m, nil
	}

	// fetch from URL
	resp, err := c.Get("https://swapi.co/api/?format=json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return nil, err
	}

	if gob.NewEncoder(f).Encode(m); err != nil {
		return nil, err
	}

	return m, nil
}

// Entity Collections

type Appendable interface {
	Append(result interface{})
	Len() int
}

type People []Person

func (p *People) Append(result interface{}) {
	res := result.(*PersonResult)
	*p = append(*p, res.Results...)
}

func (p *People) Len() int {
	return len(*p)
}

type Planets []Planet

func (p *Planets) Append(result interface{}) {
	res := result.(*PlanetResult)
	*p = append(*p, res.Results...)
}

func (p *Planets) Len() int {
	return len(*p)
}

type Films []Film

func (f *Films) Append(result interface{}) {
	res := result.(*FilmResult)
	*f = append(*f, res.Results...)
}

func (f *Films) Len() int {
	return len(*f)
}

type Species []Specie

func (f *Species) Append(result interface{}) {
	res := result.(*SpeciesResult)
	*f = append(*f, res.Results...)
}

func (f *Species) Len() int {
	return len(*f)
}

type Starships []Starship

func (f *Starships) Append(result interface{}) {
	res := result.(*StarshipResult)
	*f = append(*f, res.Results...)
}

func (f *Starships) Len() int {
	return len(*f)
}

type Vehicles []Vehicle

func (f *Vehicles) Append(result interface{}) {
	res := result.(*VehicleResult)
	*f = append(*f, res.Results...)
}

func (f *Vehicles) Len() int {
	return len(*f)
}

// Entities
type Ids []string

func (fi *Ids) UnmarshalJSON(b []byte) error {
	var urls []string
	err := json.Unmarshal(b, &urls)
	if err != nil {
		return err
	}

	for i, v := range urls {
		urls[i] = url2id(v)
	}

	*fi = urls

	return nil
}

type Url string

func (ptr *Url) UnmarshalJSON(b []byte) error {
	var u string
	err := json.Unmarshal(b, &u)
	if err != nil {
		return err
	}

	*ptr = Url(url2id(u))

	return nil
}

type RestFields struct {
	Id      string
	Created string
	Edited  string
	Url     Url `json:",omitempty"`
}

type Person struct {
	Name      string
	Height    string
	Mass      string
	SkinColor string `json:"skin_color"`
	EyeColor  string `json:"eye_color"`
	BirthYear string `json:"birth_year"`
	Gender    string
	Homeworld string
	Films     Ids
	Species   Ids
	Vehicles  Ids
	Starships Ids
	RestFields
}

type Planet struct {
	Name           string
	Diameter       string
	RotationPeriod string `json:"rotation_period"`
	OrbitalPeriod  string `json:"orbital_period"`
	Gravity        string
	Population     string
	Climate        string
	Terrain        string
	SurfaceWater   string `json:"surface_water"`
	Residents      Ids
	Films          Ids
	RestFields
}

type Film struct {
	Title        string
	EpisodeID    int    `json:"episode_id"`
	OpeningCrawl string `json:"opening_crawl"`
	Director     string
	Producer     string
	ReleaseDate  string `json:"release_date"`
	Species      Ids
	Starships    Ids
	Vehicles     Ids
	Characters   Ids
	Planets      Ids
	RestFields
}

type Specie struct {
	Name            string
	Classification  string
	Designation     string
	AverageHeight   string `json:"average_height"`
	AverageLifespan string `json:"average_lifespan"`
	EyeColors       string `json:"eye_colors"`
	HairColors      string `json:"hair_colors"`
	SkinColors      string `json:"skin_colors"`
	Language        string
	Homeworld       string
	People          Ids
	Films           Ids
	RestFields
}

type Transport struct {
	CargoCapacity        string `json:"cargo_capacity"`
	Consumables          string
	CostInCredits        string `json:"cost_in_credits"`
	Crew                 string
	Films                Ids
	Length               string
	Manufacturer         string
	MaxAtmospheringSpeed string `json:"max_atmosphering_speed"`
	Model                string
	Name                 string
	Passengers           string
	Pilots               Ids

	RestFields
}

type Starship struct {
	Transport

	Mglt             string
	HyperdriveRating string `json:"hyperdrive_rating"`
	StarshipClass    string `json:"starship_class"`
}

type Vehicle struct {
	Transport

	VehicleClass string `json:"vehicle_class"`
}

// API results

type Resultable interface {
	NextPage() string
	Reset()
}

type Result struct {
	Count    int
	Next     string
	Previous string
}

func (r *Result) Reset() {
	r.Count = 0
	r.Next = ""
	r.Previous = ""
}

func (r *Result) NextPage() string {
	return r.Next
}

type PersonResult struct {
	Result
	Results People
}

func (r *PersonResult) Reset() {
	r.Result.Reset()
	r.Results = make(People, 0)
}

type PlanetResult struct {
	Result
	Results Planets
}

func (r *PlanetResult) Reset() {
	r.Result.Reset()
	r.Results = make(Planets, 0)
}

type FilmResult struct {
	Result
	Results Films
}

func (r *FilmResult) Reset() {
	r.Result.Reset()
	r.Results = make(Films, 0)
}

type SpeciesResult struct {
	Result
	Results Species
}

func (r *SpeciesResult) Reset() {
	r.Result.Reset()
	r.Results = make(Species, 0)
}

type StarshipResult struct {
	Result
	Results Starships
}

func (r *StarshipResult) Reset() {
	r.Result.Reset()
	r.Results = make(Starships, 0)
}

type VehicleResult struct {
	Result
	Results Vehicles
}

func (r *VehicleResult) Reset() {
	r.Result.Reset()
	r.Results = make(Vehicles, 0)
}

func fetch(c *http.Client, next, filename string, coll Appendable, result Resultable) error {
	if next == "" {
		return errors.New("no start URL provided")
	}

	f, err := open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	// decode from the GOB file if possible
	err = decode(f, coll)
	if err != nil && err != ErrEmpty {
		return err
	} else if err == nil && coll.Len() > 0 {
		fmt.Printf("read %v entries from %s\n", coll.Len(), filename)
		return nil
	}

	// no entries found, download from the Star Wars API.
	for next != "" {
		fmt.Println(next)
		resp, err := c.Get(next)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
			return err
		}

		coll.Append(result)

		next = result.NextPage()
		result.Reset()
	}

	// write the results out to the GOB file
	if err = gob.NewEncoder(f).Encode(coll); err != nil {
		return err
	}

	return nil
}
