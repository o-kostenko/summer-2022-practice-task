package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"strconv"
	"time"
)

var (
	ErrWrongCriteria         = errors.New("unsupported criteria")
	ErrEmptyDepartureStation = errors.New("empty departure station")
	ErrEmptyArrivalStation   = errors.New("empty arrival station")
	ErrBadDepartureStation   = errors.New("bad arrival station input")
	ErrBadArrivalStation     = errors.New("bad departure station input")
)

type Trains []Train

type Train struct {
	TrainID            int       `json:"trainId"`
	DepartureStationID int       `json:"departureStationId"`
	ArrivalStationID   int       `json:"arrivalStationId"`
	Price              float32   `json:"price"`
	ArrivalTime        time.Time `json:"arrivalTime"`
	DepartureTime      time.Time `json:"departureTime"`
}

func main() {
	//	... запит даних від користувача
	var (
		departureStation, arrivalStation, criteria string
	)
	fmt.Println("Введіть станцію відправлення:")
	_, err := fmt.Scanf("%s", &departureStation)
	if err != nil {
		log.Fatalf("error when input %s", err.Error())
	}

	fmt.Println("Введіть станцію прибуття:")
	_, err = fmt.Scanf("%s", &arrivalStation)
	if err != nil {
		log.Fatalf("error when input %s", err.Error())
	}

	fmt.Println("Введіть критерій, по якому сортувати результат (price, arrival-time, departure-time):")
	_, err = fmt.Scanf("%s", &criteria)
	if err != nil {
		log.Fatalf("error when input %s", err.Error())
	}

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	//	... обробка помилки
	if err != nil {
		log.Fatal(err)
	}

	//	... друк result
	PrintTrains(result)
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// ... код
	var data, result Trains

	if departureStation == "" {
		return nil, ErrEmptyDepartureStation
	}

	dSt, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, ErrBadDepartureStation
	}

	if dSt < 1 {
		return nil, ErrBadDepartureStation
	}

	if arrivalStation == "" {
		return nil, ErrEmptyArrivalStation
	}

	aSt, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, ErrBadArrivalStation
	}

	if aSt < 1 {
		return nil, ErrBadArrivalStation
	}

	file, err := ioutil.ReadFile("data.json")
	if err != nil {
		return nil, fmt.Errorf("file read error:%w", err)
	}

	err = json.Unmarshal(file, &data)
	if err != nil {
		return nil, fmt.Errorf("parse json error:%w", err)
	}

	switch criteria {
	case "price", "arrival-time", "departure-time":
	default:
		return nil, ErrWrongCriteria
	}

	for i := range data {
		if data[i].DepartureStationID == dSt && data[i].ArrivalStationID == aSt {
			result = append(result, data[i])
		}
	}

	switch criteria {
	case "price":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].Price < result[j].Price
		})

	case "arrival-time":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].ArrivalTime.Sub(result[j].ArrivalTime) < 0
		})

	case "departure-time":
		sort.SliceStable(result, func(i, j int) bool {
			return result[i].DepartureTime.Sub(result[j].DepartureTime) < 0
		})
	}

	// маєте повернути правильні значення
	return result, nil
}

func PrintTrains(trains Trains) {
	fmt.Printf("trainId depStId arrStId price   departureTime   arrivalTime\n")
	for i := range trains {
		fmt.Printf("%v \t%v  \t%v  \t%v \t%d:%d:%d   \t%d:%d:%d \n",
			trains[i].TrainID,
			trains[i].DepartureStationID,
			trains[i].ArrivalStationID,
			trains[i].Price,
			trains[i].DepartureTime.Hour(),
			trains[i].DepartureTime.Minute(),
			trains[i].DepartureTime.Second(),
			trains[i].ArrivalTime.Hour(),
			trains[i].ArrivalTime.Minute(),
			trains[i].ArrivalTime.Second())

	}
}

func (t *Train) UnmarshalJSON(data []byte) error {
	type Alias Train

	aux := &struct {
		ArrivalTime   string `json:"arrivalTime"`
		DepartureTime string `json:"departureTime"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}

	err := json.Unmarshal(data, aux)
	if err != nil {
		return fmt.Errorf("read time as string: %w", err)
	}

	t.ArrivalTime, err = time.Parse("15:04:05", aux.ArrivalTime)
	if err != nil {
		return fmt.Errorf("error parse arrival time: %w", err)
	}

	t.DepartureTime, err = time.Parse("15:04:05", aux.DepartureTime)
	if err != nil {
		return fmt.Errorf("error parse departure time: %w", err)
	}

	return nil
}
