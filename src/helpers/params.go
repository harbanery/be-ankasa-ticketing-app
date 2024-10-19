package helpers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetSortParams(sort, orderBy string) string {
	if sort != "ASC" && sort != "DESC" {
		sort = "DESC"
	}

	if orderBy == "" {
		orderBy = "updated_at"
	}

	return orderBy + " " + sort
}

func GetPaginationParams(oldLimit, oldPage string) (int, int, int) {
	page, _ := strconv.Atoi(oldPage)
	if page == 0 {
		page = 1
	}

	limit, _ := strconv.Atoi(oldLimit)
	if limit == 0 {
		limit = 10
	}

	offset := (page - 1) * limit

	return page, limit, offset
}

func GetFilterParams(c *fiber.Ctx) (map[string]interface{}, error)  {
	filter := make(map[string]interface{})
	params := c.Queries()

	merchant := params["merchant"]
	if len(merchant) > 0 {
        filter["merchantValues"] = merchant
    }

	class := params["class"]
	if class != "" {
		classSplit := strings.Split(class, ",")
		filter["classValues"] = classSplit
	}

	arrival := params["arrival"]
	if arrival != "" {
		arrivalSplit := strings.Split(arrival, ",")
		arrivalTimes := make([]time.Time, len(arrivalSplit))
		for i, a := range arrivalSplit {
			arrivalTime, err := time.Parse(time.RFC3339, a)
			if err != nil {
				return nil, fmt.Errorf("invalid arrival time format: %v", err)
			}
			arrivalTimes[i] = arrivalTime
		}
		filter["arrivalValues"] = arrivalTimes
	}

	// Memproses departure time
	departure := params["departure"]
	if departure != "" {
		departureSplit := strings.Split(departure, ",")
		departureTimes := make([]time.Time, len(departureSplit))
		for i, d := range departureSplit {
			departureTime, err := time.Parse(time.RFC3339, d)
			if err != nil {
				return nil, fmt.Errorf("invalid departure time format: %v", err)
			}
			departureTimes[i] = departureTime
		}
		filter["departureValues"] = departureTimes
	}

	return filter, nil

}