package scheduler

import (
	"encoding/gob"
	"os"
	"time"
)

// Service реализует простой планировщик сканирований
type Service struct {
	History map[string]time.Time
	ExpDays int
}

// New - конструктор.
func New(sites []string, expdays int) *Service {
	history := make(map[string]time.Time)
	for _, site := range sites {
		history[site] = time.Unix(0, 0)
	}
	s := Service{
		History: history,
		ExpDays: expdays,
	}
	return &s
}

// UpdateHistory обновляет даты сканирования сайтов в истории
func (s *Service) UpdateHistory(sites []string) {
	for _, site := range sites {
		s.History[site] = time.Now().Truncate(time.Hour * 24)
	}
}

// ExpiredSites возвращает устаревшие сайты
func (s *Service) ExpiredSites() []string {
	expired := make([]string, 0)
	for site, scanDate := range s.History {
		days := int(time.Now().Sub(scanDate).Hours() / 24)
		if days >= s.ExpDays {
			expired = append(expired, site)
		}
	}
	return expired
}

// SaveTo сохраняет планировщик в файл
func (s *Service) SaveTo(fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	enc := gob.NewEncoder(file)
	err = enc.Encode(s)
	if err != nil {
		return err
	}
	return nil
}

// LoadFrom создает планировщик из файла
func LoadFrom(fileName string) (*Service, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var s Service
	dec := gob.NewDecoder(file)
	err = dec.Decode(&s)
	if err != nil {
		return nil, err
	}
	return &s, nil
}
