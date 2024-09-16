package controllers

import (
	"ankasa-be/src/models"
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func GenerateMerchantSeed(c *fiber.Ctx) error {
	merchantExists := models.SelectAllMerchants()
	if len(merchantExists) > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Cannot seed while already data in here",
		})
	}

	merchants := []struct {
		Email       string
		Name        string
		Image       string
		Description string
		Classes     []models.Class
	}{
		{
			Email:       "ankasa-merchant@garuda-indonesia.com",
			Name:        "Garuda Indonesia",
			Image:       "https://s3-alpha-sig.figma.com/img/d670/ed6a/9d205fa306085ffa6cc1365eef78958f?Expires=1725840000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=Ia-KU3qxHdoMyUayAVXixT5wIecYVtkYn0~XcMJ3joud3TWL5U52yafo2eFa2ML1DIil6X0QGtGhEE3NQEAQkufvx6WR3hVuOwZSg~fza52kl0YJ9AVih6E21ZnITEOPk-oqDCHtsKTGL60cuvOpJ6upY8FOJZT2LFOucgNFUdSJjJmHyt4z~i9rdlSBdXR8JBn5TR2jRmWXuqvn0xNHa5EDmCx4FwDEYgZDotUPvzH3VpxYgQjKJ~Vz-EF1wS6jyRxh28J~I6V3gmNc20q5aPp6z4YNCDekWZ4Ev~cvQ5HBnnqKNMOzDWXzhihDFKL~PPDLUWByQSirQh7NNWEOog__",
			Description: "Garuda Indonesia is the national airline of Indonesia, offering a wide range of domestic and international flights.",
			Classes: []models.Class{
				{
					Name:           "First Class",
					Price:          6000000,
					Seats:          8,
					RowSeats:       4,
					IsRefund:       true,
					IsReschedule:   true,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         true,
				},
				{
					Name:           "Business",
					Price:          3500000,
					Seats:          38,
					RowSeats:       7,
					IsRefund:       true,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         true,
				},
				{
					Name:           "Economy",
					Price:          1500000,
					Seats:          268,
					RowSeats:       10,
					IsRefund:       false,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         false,
				},
			},
		},
		{
			Email:       "merchant_ankasa@airasia.com",
			Name:        "Air Asia",
			Image:       "https://s3-alpha-sig.figma.com/img/d2d5/5860/4c04453948ea266977fb0d4595f0fb73?Expires=1725840000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=lHJBLOixZ-LBc5rb4kSCMuJ6FfIGo~H2gVNEiXKYiRKrGHGcFot-yn~lkZ~t8fFeRkfStTQMkeVRje8LHLY5ddqRedLCLzvj4ejoUbDoVRjPOPB16-4PLhovK1wvSQEaTKFZKw6uOnGAqkjH8CNDRbJkRvq2GnGsGMui1tv5e3mvqjfxWZqIlimpMcI6wnTJIAHgamyc01jt~8dddVLtnvVFFE4FTbGAtPylmemhfdwni~PCDZFRZeHSn5bEwQF9Gea5IRaLQds3Gk~gv6q7JcDXtUl9pqXt2Iaj7XCrwkfpa4QgsIczlig8dhHBWlDIyDBFCfkYK8le-0dwNfxI3w__",
			Description: "Air Asia is a low-cost airline headquartered in Malaysia, known for its affordable fares and extensive route network.",
			Classes: []models.Class{
				{
					Name:           "Premium Economy",
					Price:          2200000,
					Seats:          110,
					RowSeats:       9,
					IsRefund:       true,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         true,
				},
				{
					Name:           "Economy",
					Price:          1000000,
					Seats:          270,
					RowSeats:       9,
					IsRefund:       false,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: false,
					IsWifi:         false,
				},
			},
		},
		{
			Email:       "ankasamerchantteam@lionair.com",
			Name:        "Lion Air",
			Image:       "https://s3-alpha-sig.figma.com/img/3237/1e89/565a85aff8320c77024052633bf9ac42?Expires=1725840000&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=pc0U~SMtYFz8WzPGoidVrJueHhAmWWFw86TPcOE69u9y8XJZXkgmyYCyVpkBDc-PDaDJG2HPGIlyC-TxoZcOiniMOpPWUTtHd7H8-jbFFTgxE0RvGo2NBYxVCgYwYhFRW25311LcHLNYfPA4Sx132RsVycPyK30br09HPgV40UpQ2R1iTbrCO88SdvEEI6dS9QffSvpF6QXp7DXEY4eT8oz7gCMH9hQByqWAUnHeZqYt16guQj3INXP0ajXE572vSOkPe3ssnonVFmKtYwc2cDiQgWxU-SFgcPbNPWr09IkrFXn4~XXATyngaIHP45H0a-BX4nUHuYl-80yAxT1rOQ__",
			Description: "Lion Air is an Indonesian low-cost airline, providing domestic and international flights with a focus on Southeast Asia.",
			Classes: []models.Class{
				{
					Name:           "Business",
					Price:          2500000,
					Seats:          10,
					RowSeats:       4,
					IsRefund:       true,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         true,
				},
				{
					Name:           "Economy",
					Price:          900000,
					Seats:          198,
					RowSeats:       6,
					IsRefund:       false,
					IsReschedule:   false,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         false,
				},
			},
		},
	}

	var createdMerchants []fiber.Map

	for _, merchantData := range merchants {
		tokenBytes := make([]byte, 14)
		if _, err := rand.Read(tokenBytes); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Password error",
			})
		}

		token := hex.EncodeToString(tokenBytes)

		hashPassword, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.DefaultCost)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Password error",
			})
		}

		user := models.User{
			Email:    merchantData.Email,
			Password: string(hashPassword),
			Role:     "merchant",
			IsVerify: "true",
		}

		userID, err := models.CreateUser(&user)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Failed to create user",
			})
		}

		merchant := models.Merchant{
			UserID:      userID,
			Name:        merchantData.Name,
			Image:       merchantData.Image,
			Description: merchantData.Description,
		}

		merchantID, err := models.CreateMerchant(&merchant)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Failed to create merchant",
			})
		}

		for _, classData := range merchantData.Classes {
			class := models.Class{
				MerchantID:     merchantID,
				Name:           classData.Name,
				Price:          classData.Price,
				Seats:          classData.Seats,
				RowSeats:       classData.RowSeats,
				IsRefund:       classData.IsRefund,
				IsReschedule:   classData.IsReschedule,
				IsLuggage:      classData.IsLuggage,
				IsInflightMeal: classData.IsInflightMeal,
				IsWifi:         classData.IsWifi,
			}

			if err := models.CreateClass(&class); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "server error",
					"statusCode": 500,
					"message":    "Failed to create class",
				})
			}
		}

		createdMerchants = append(createdMerchants, fiber.Map{
			"email":    merchantData.Email,
			"password": token,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":      "success",
		"statusCode":  200,
		"message":     "Merchants created successfully.",
		"dataCreated": createdMerchants,
	})
}

func GenerateCityCountrySeed(c *fiber.Ctx) error {
	cityExists, _ := models.SelectAllCities()
	if len(cityExists) > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Cannot seed while already data in here",
		})
	}

	countryExists, _ := models.SelectAllCountries()
	if len(countryExists) > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Cannot seed while already data in here",
		})
	}

	countries := []struct {
		Name   string
		Code   string
		Cities []struct {
			Name  string
			Image string
		}
	}{
		{
			Name: "Indonesia",
			Code: "IDN",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "Jakarta", Image: "https://c4.wallpaperflare.com/wallpaper/702/43/725/jakarta-city-cityscape-wallpaper-preview.jpg"},
				{Name: "Bali", Image: "https://s3-alpha-sig.figma.com/img/0656/5c4d/b656f5710f5e0f053aa347575fbc57a0?Expires=1725235200&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=Rzg75UjdbZ3okMWWhD-fvol15kkYQU6MwIc7hzfW8ei843Tf6lNWA02DJ7dihj9ogfgaP2-eOctXToKeaRiS721qkfYrjqm~bS95Sa9~jT8I6wUPsOlQAQKkeOtbToLw4~MGlEQZYsQyOibUXq05gpEBNRaLVNYa9602AXfmHLm-LCfO8DhSpYDpZM1zNg~akdjKGP-tVZjusqvKQANzc5UQl6oGSsS74uJVsXMgxht094T9GeJvnFaMU7SHeUha5ZwIuGtL3WYQn93G-vz3AXmJE-rQbJrEpZrguwgcBGhv-LsGnM0pnhXhdBe0shnVh3J5P2gCZEt2biwuY26BiA__"},
			},
		},
		{
			Name: "Singapore",
			Code: "SGP",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "Changi Bay", Image: "https://s3-alpha-sig.figma.com/img/91eb/c3b9/e985416af059aab94180bce2220da23c?Expires=1725235200&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=gFso4Fi-j-cDRRLtDLhaWLrpXoIJUIO~-l6zPhzQtMGih44nbJRxedYFGEfLh4DsbH2kIHQU3N~0He7aSdtJBX3eioys2PyMP0UAXZvAlJtJbms3WQsz342eusS~DqUUz6bDYf1aLioV4ZYF0zL77mwjRZh7n6tjLZvhyUyfEzYUmB~CCOwRTBW20ll~73EGvxbXBg~tvgc3lGU7oU6KwTvo4MoBbtTNj9D~RtrluXGVN-0yv9WOt~y4hJdea~yHJlL1nup01RbdSZys6p4PGUJfrA7G9zog8FmlcNWPCg~M1Wq6xF1cjeOy-pKHLNHDKiuXWBYxrPEdU7TSrW3KXQ__"},
			},
		},
		{
			Name: "Malaysia",
			Code: "MYS",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "Kuala Lumpur", Image: "https://media.istockphoto.com/id/955628078/id/foto/singapura-singapura.jpg?s=612x612&w=0&k=20&c=Di2CyvzJQWSYDe0pYPFxjNKiU3I8rUf4SP8hGfbtYDE="},
			},
		},
		{
			Name: "Japan",
			Code: "JPN",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "Tokyo", Image: "https://a.loveholidays.com/media-library/~production/6d7b5475d338ee30647c4b88d27fea204429081c-3840x1408.jpg?auto=avif%2Cwebp&quality=80&dpr=1.5&optimize=high&fit=crop&width=1280&height=380"},
			},
		},
		{
			Name: "France",
			Code: "FRA",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "Paris", Image: "https://s3-alpha-sig.figma.com/img/206e/fc86/f01e0ab7c33981f586273a726b9e138f?Expires=1725235200&Key-Pair-Id=APKAQ4GOSFWCVNEHN3O4&Signature=LPBoB91Apbv21a8-x031xV5yX6ZRFjVR23GG9qllmDTQjsM0jNJ8wFMrv1p0uCoLztWbizaUB-CYc19PtLI9Y-Aauvlj5AW1WziwcstqkAKbRKwfDGLJUDvyiumhPc3HKhdAYJehoXh-vaPz0QJbyeFoULBVMQzi1fGQGpJ~q8rKwVl7sHwBg1uHcnZkRmfKOyrtrWQEnnBzLYHk0qxaiuFiFG4m~B6FA6sr-rnI~NUBsOaU4-vNyFaxn0y8fq5fZjluKX~i2V2T-ASBh8OFKucIy1USzccg8xuTFKuLEEkvJvFokymIRJVTwiENtW1-LVlzkcDZj-269XEYtfwOJg__"},
			},
		},
		{
			Name: "United States",
			Code: "USA",
			Cities: []struct {
				Name  string
				Image string
			}{
				{Name: "New York", Image: "https://images.pexels.com/photos/597909/pexels-photo-597909.jpeg"},
			},
		},
	}

	for _, countryData := range countries {
		country := models.Country{
			Name: countryData.Name,
			Code: countryData.Code,
		}

		countryID, err := models.CreateCountry(&country)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "server error",
				"statusCode": 500,
				"message":    "Failed to create country",
			})
		}

		for _, cityData := range countryData.Cities {
			city := models.City{
				Name:      cityData.Name,
				Image:     cityData.Image,
				CountryID: countryID,
			}
			if err := models.CreateCity(&city); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "server error",
					"statusCode": 500,
					"message":    "Failed to create city",
				})
			}
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Cities and countries created successfully.",
	})
}

func GenerateTicketSeed(c *fiber.Ctx) error {
	ticketExists := models.SelectAllTickets()
	if len(ticketExists) > 0 {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Cannot seed while already data in here",
		})
	}

	merchants := models.SelectAllMerchants()
	if len(merchants) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Merchants not found",
		})
	}

	cities, _ := models.SelectAllCities()
	if len(cities) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Cities not found",
		})
	}

	// rand.Seed(time.Now().UnixNano())

	for _, merchant := range merchants {
		rowStart := 1
		for _, class := range merchant.Classes {
			stock := class.Seats
			if stock == 0 {
				return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"status":     "not found",
					"statusCode": 404,
					"message":    "No stock in here",
				})
			}
			price := class.Price
			rowSeats := class.RowSeats

			departureCity := cities[rand.Intn(len(cities))]
			arrivalCity := cities[rand.Intn(len(cities))]
			for arrivalCity.ID == departureCity.ID {
				arrivalCity = cities[rand.Intn(len(cities))]
			}

			departureTime := generateRandomTime()
			arrivalTime := departureTime.Add(time.Duration(rand.Intn(3)+2) * time.Hour) // 2-4 jam setelah arrival

			departure := models.Departure{
				Schedule: &departureTime,
				CityID:   departureCity.ID,
			}

			arrival := models.Arrival{
				Schedule: &arrivalTime,
				CityID:   arrivalCity.ID,
			}

			ticket := models.Ticket{
				Stock:      stock,
				Price:      price,
				MerchantID: merchant.ID,
				ClassID:    class.ID,
				Gate:       "G-" + strconv.Itoa(rand.Intn(20)+1),
				Arrival:    arrival,
				Departure:  departure,
			}

			ticketID, err := models.CreateTicket(&ticket)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "server error",
					"statusCode": 500,
					"message":    "Failed to create ticket",
				})
			}

			seats := make([]models.Seat, stock)
			rowLetters := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
			for i := 0; i < stock; i++ {
				row := (rowStart) + (i / rowSeats)
				seatIndex := i % rowSeats
				seatCode := fmt.Sprintf("%d-%c", row, rowLetters[seatIndex])

				seats[i] = models.Seat{
					Code:      seatCode,
					IsBooking: false,
					TicketID:  ticketID,
				}

				if err := models.CreateSeat(&seats[i]); err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
						"status":     "server error",
						"statusCode": 500,
						"message":    "Failed to create seat",
					})
				}

			}
			rowStart += (stock / rowSeats) + 1

		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Tickets with arrival, departure, and seats created successfully.",
	})
}

func generateRandomTime() time.Time {
	now := time.Now()
	randomOffset := time.Duration(rand.Intn(48)) * time.Hour
	return now.Add(randomOffset)
}
