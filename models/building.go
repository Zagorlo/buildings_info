package models

import (
	"time"
)

//Задание:
//----------
//Сделать сервис с одной API - изменение корпуса по ID (Заполнить можно любыми данными)
//
//Есть таблица building с атрибутами (name, floors_count, parking_available, parking_count)
//При изменении хранить 10 последних изменений только полей, которые изменились.
//
//Никаких требований по используемым библиотек нет.
//
//Результат предоставить в виде ссылки на репозиторий в гитхабе.
//--------- end

type Building struct {
	tableName        struct{} `pg:"buildings"`
	ID               int64    `json:"id" pg:"id, pk"`
	Name             *string  `json:"name" pg:"name"`
	FloorsCount      *int16   `json:"floors_count" pg:"floors_count"`
	ParkingCount     *int16   `json:"parking_count" pg:"parking_count"`
	ParkingAvailable *bool    `json:"parking_available" pg:"parking_available"`
}

type BuildingUpdate struct {
	tableName        struct{}  `pg:"buildings_updates"`
	ID               int64     `json:"id" pg:"id"`
	Name             string    `json:"name" pg:"name"`
	FloorsCount      int16     `json:"floors_count" pg:"floors_count"`
	ParkingCount     int16     `json:"parking_count" pg:"parking_count"`
	ParkingAvailable bool      `json:"parking_available" pg:"parking_available"`
	RemovedAt        time.Time `json:"removed_at" pf:"removed_at"`
}

func (b *Building) PrepareValues() {
	b.ID = 0
}

func (b *Building) RetrieveChanges(a *Building) []CacheItem {
	var changes []CacheItem
	now := time.Now()

	if a == nil {
		changes = []CacheItem{
			{
				b.ID, "name", nil, *b.Name, now,
			},
			{
				b.ID, "floors_count", nil, *b.Name, now,
			},
			{
				b.ID, "parkings_count", nil, *b.ParkingCount, now,
			},
			{
				b.ID, "parking_available", nil, *b.ParkingAvailable, now,
			},
		}

		return changes
	}

	if *b.Name != *a.Name {
		changes = append(changes, CacheItem{b.ID, "name", *a.Name, *b.Name, now})
	}

	if *b.FloorsCount != *a.FloorsCount {
		changes = append(changes, CacheItem{b.ID, "floors_count", *a.FloorsCount, *b.FloorsCount, now})
	}

	if *b.ParkingCount != *a.ParkingCount {
		changes = append(changes, CacheItem{b.ID, "parking_count", *a.ParkingCount, *b.ParkingCount, now})
	}

	if *b.ParkingAvailable != *a.ParkingAvailable {
		changes = append(changes, CacheItem{b.ID, "parking_available", *a.ParkingAvailable, *b.ParkingAvailable, now})
	}

	return changes
}

func (bu *BuildingUpdate) RetrieveBuildingUpdates() []CacheItem {
	return []CacheItem{
		{
			bu.ID, "name", nil, bu.Name, bu.RemovedAt,
		},
		{
			bu.ID, "floors_count", nil, bu.Name, bu.RemovedAt,
		},
		{
			bu.ID, "parkings_count", nil, bu.ParkingCount, bu.RemovedAt,
		},
		{
			bu.ID, "parking_available", nil, bu.ParkingAvailable, bu.RemovedAt,
		},
	}
}
