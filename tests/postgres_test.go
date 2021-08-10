package tests

import (
	"buildings_info/consts"
	"buildings_info/logging"
	"buildings_info/models"
	"buildings_info/service"
	"buildings_info/storages/buildings"
	"buildings_info/storages/models/postgresDriver"
	"buildings_info/storages/models/serviceCache"
	"context"
	"encoding/json"
	"fmt"
	uuid2 "github.com/satori/go.uuid"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"testing"
)

func TestPostgresContainer(t *testing.T) {
	logging.InitLogger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.ContextUUIDKey, uuid2.NewV4())

	postgresDB, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:11.6-alpine",
			ExposedPorts: []string{"5432/tcp"},
			SkipReaper:   true,
		},
		Started: true,
	})
	require.Nil(t, err)

	defer postgresDB.Terminate(ctx)

	ip, err := postgresDB.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := postgresDB.MappedPort(ctx, "5432")
	require.Nil(t, err)

	cfg := testCfg(ip, port.Port(), ctx)

	driver := postgresDriver.NewPostgresDriver(&cfg.Postgres, cfg.Context)
	ec := driver.PostgresOpen()
	require.Equal(t, ec.NotNil(), false)
	defer driver.PostgresClose()

	ec = driver.PostgresCreations("../init.sql")
	require.Equal(t, ec.NotNil(), false)

	buildingsModule := &buildings.BuildingsModule{}
	buildingsModule.Init(cfg.Context, driver.GetConn(), serviceCache.NewJsonOrderedCache(cfg.BuildingsCache))
	ser := &service.Service{
		buildingsModule,
	}
	ser.FillBuildingsCache(cfg.Context)

	bytes := ser.GetBuildingsChangesInfo()
	require.Equal(t, string(bytes), `[{"id":2,"field_name":"name","old_value":null,"new_value":"eee","updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"floors_count","old_value":null,"new_value":"eee","updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"parkings_count","old_value":null,"new_value":0,"updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"parking_available","old_value":null,"new_value":false,"updated_at":"2012-12-20T00:00:00Z"},{"id":1,"field_name":"name","old_value":null,"new_value":"ddd","updated_at":"2012-12-18T00:00:00Z"},{"id":1,"field_name":"floors_count","old_value":null,"new_value":"ddd","updated_at":"2012-12-18T00:00:00Z"},{"id":1,"field_name":"parkings_count","old_value":null,"new_value":35,"updated_at":"2012-12-18T00:00:00Z"},{"id":1,"field_name":"parking_available","old_value":null,"new_value":true,"updated_at":"2012-12-18T00:00:00Z"}]`)

	ec = ser.InsertBuildingInfo(cfg.Context, &buildingToInsert)
	require.Equal(t, ec.NotNil(), false)

	ec = ser.UpdateBuildingInfo(cfg.Context, &buildingToUpdateIncorrect)
	require.Equal(t, ec.NotNil(), true)

	ec = ser.UpdateBuildingInfo(cfg.Context, &buildingToUpdateCorrect1)
	require.Equal(t, ec.NotNil(), false)

	ec = ser.UpdateBuildingInfo(cfg.Context, &buildingToUpdateCorrect2)
	require.Equal(t, ec.NotNil(), false)

	//bytes = ser.GetBuildingsChangesInfo()
	//require.Equal(t, string(bytes), `[{"id":4,"field_name":"name","old_value":null,"new_value":"111","updated_at":"2021-08-10T17:30:59.557782+03:00"},{"id":4,"field_name":"floors_count","old_value":null,"new_value":"111","updated_at":"2021-08-10T17:30:59.557782+03:00"},{"id":4,"field_name":"parkings_count","old_value":null,"new_value":0,"updated_at":"2021-08-10T17:30:59.557782+03:00"},{"id":4,"field_name":"parking_available","old_value":null,"new_value":false,"updated_at":"2021-08-10T17:30:59.557782+03:00"},{"id":2,"field_name":"name","old_value":null,"new_value":"eee","updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"floors_count","old_value":null,"new_value":"eee","updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"parkings_count","old_value":null,"new_value":0,"updated_at":"2012-12-20T00:00:00Z"},{"id":2,"field_name":"parking_available","old_value":null,"new_value":false,"updated_at":"2012-12-20T00:00:00Z"},{"id":1,"field_name":"name","old_value":null,"new_value":"ddd","updated_at":"2012-12-18T00:00:00Z"},{"id":1,"field_name":"floors_count","old_value":null,"new_value":"ddd","updated_at":"2012-12-18T00:00:00Z"}]`)

	fmt.Println(string(bytes))
}

const (
	buildingToInsertJSON          = `{"id":1234556789,"name":"111","floors_count":8,"parking_count":0,"parking_available":false}`
	buildingToUpdateIncorrectJSON = `{"id":7,"name":"111","floors_count":8,"parking_count":0,"parking_available":false}`
	buildingToUpdateCorrect1JSON  = `{"id":2,"name":"222","floors_count":2,"parking_count":0,"parking_available":false}`
	buildingToUpdateCorrect2JSON  = `{"id":2,"name":"333","floors_count":2,"parking_count":0,"parking_available":true}`
)

var (
	buildingToInsert, buildingToUpdateIncorrect, buildingToUpdateCorrect1, buildingToUpdateCorrect2 models.Building

	_ = json.Unmarshal([]byte(`{"id":1234556789,"name":"111","floors_count":8,"parking_count":0,"parking_available":false}`), &buildingToInsert)
	_ = json.Unmarshal([]byte(`{"id":7,"name":"111","floors_count":8,"parking_count":0,"parking_available":false}`), &buildingToUpdateIncorrect)
	_ = json.Unmarshal([]byte(`{"id":2,"name":"222","floors_count":2,"parking_count":0,"parking_available":false}`), &buildingToUpdateCorrect1)
	_ = json.Unmarshal([]byte(`{"id":2,"name":"333","floors_count":2,"parking_count":0,"parking_available":true}`), &buildingToUpdateCorrect2)
)
