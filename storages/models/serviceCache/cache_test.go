package serviceCache

import (
	"buildings_info/consts"
	"buildings_info/logging"
	"buildings_info/models"
	"context"
	uuid2 "github.com/satori/go.uuid"
	"strconv"
	"time"

	"github.com/stretchr/testify/require"
	"testing"
)

func TestPostgresContainer(t *testing.T) {
	logging.InitLogger()
	ctx := context.Background()
	ctx = context.WithValue(ctx, consts.ContextUUIDKey, uuid2.NewV4())

	cache := NewJsonOrderedCache(models.BuildingsCache{10})

	bytes := cache.RetrieveItems()
	require.Equal(t, string(bytes), `[]`)

	now1 := time.Time{}.Add(5 * time.Hour)
	now2 := time.Time{}.Add(777 * time.Millisecond)
	now3 := time.Time{}.Add(1000000 * time.Hour)

	cache.PrependItems(ctx, []models.CacheItem{
		{
			1,
			"a",
			nil,
			76,
			now1,
		},
		{
			2,
			"a",
			nil,
			76,
			now1,
		},
		{
			3,
			"c",
			nil,
			76,
			now1,
		},
	})

	bytes = cache.RetrieveItems()
	require.Equal(t, string(bytes), `[{"id":1,"field_name":"a","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"},{"id":2,"field_name":"a","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"},{"id":3,"field_name":"c","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"}]`)
	cache.PrependItems(ctx, []models.CacheItem{
		{
			4,
			"a",
			nil,
			77,
			now2,
		},
		{
			5,
			"a",
			nil,
			5,
			now2,
		},
		{
			0,
			"c",
			nil,
			"abcde",
			now2,
		},
	})

	bytes = cache.RetrieveItems()
	require.Equal(t, string(bytes), `[{"id":4,"field_name":"a","old_value":null,"new_value":77,"updated_at":"0001-01-01T00:00:00.777Z"},{"id":5,"field_name":"a","old_value":null,"new_value":5,"updated_at":"0001-01-01T00:00:00.777Z"},{"id":0,"field_name":"c","old_value":null,"new_value":"abcde","updated_at":"0001-01-01T00:00:00.777Z"},{"id":1,"field_name":"a","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"},{"id":2,"field_name":"a","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"},{"id":3,"field_name":"c","old_value":null,"new_value":76,"updated_at":"0001-01-01T05:00:00Z"}]`)
	var items = make([]models.CacheItem, 0, 100)
	for i := int64(0); i < 100; i++ {
		items = append(items, models.CacheItem{
			i,
			"field" + strconv.FormatInt(i, 10),
			-i,
			2 * i,
			now3,
		})
	}

	cache.PrependItems(ctx, items)
	bytes = cache.RetrieveItems()
	require.Equal(t, string(bytes), `[{"id":0,"field_name":"field0","old_value":0,"new_value":0,"updated_at":"0115-01-30T16:00:00Z"},{"id":1,"field_name":"field1","old_value":-1,"new_value":2,"updated_at":"0115-01-30T16:00:00Z"},{"id":2,"field_name":"field2","old_value":-2,"new_value":4,"updated_at":"0115-01-30T16:00:00Z"},{"id":3,"field_name":"field3","old_value":-3,"new_value":6,"updated_at":"0115-01-30T16:00:00Z"},{"id":4,"field_name":"field4","old_value":-4,"new_value":8,"updated_at":"0115-01-30T16:00:00Z"},{"id":5,"field_name":"field5","old_value":-5,"new_value":10,"updated_at":"0115-01-30T16:00:00Z"},{"id":6,"field_name":"field6","old_value":-6,"new_value":12,"updated_at":"0115-01-30T16:00:00Z"},{"id":7,"field_name":"field7","old_value":-7,"new_value":14,"updated_at":"0115-01-30T16:00:00Z"},{"id":8,"field_name":"field8","old_value":-8,"new_value":16,"updated_at":"0115-01-30T16:00:00Z"},{"id":9,"field_name":"field9","old_value":-9,"new_value":18,"updated_at":"0115-01-30T16:00:00Z"}]`)
}
