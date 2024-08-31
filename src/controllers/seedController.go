package controllers

import (
	"ankasa-be/src/models"
	"crypto/rand"
	"encoding/hex"

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
					Name:           "Economy",
					Quantity:       200,
					Price:          1500000,
					IsRefund:       true,
					IsReschedule:   true,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         false,
				},
				{
					Name:           "Business",
					Quantity:       50,
					Price:          3000000,
					IsRefund:       true,
					IsReschedule:   true,
					IsLuggage:      true,
					IsInflightMeal: true,
					IsWifi:         true,
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
					Name:           "Economy",
					Quantity:       180,
					Price:          1000000,
					IsRefund:       false,
					IsReschedule:   true,
					IsLuggage:      false,
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
					Name:           "Economy",
					Quantity:       200,
					Price:          1200000,
					IsRefund:       false,
					IsReschedule:   true,
					IsLuggage:      true,
					IsInflightMeal: false,
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
				Quantity:       classData.Quantity,
				Price:          classData.Price,
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
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Cities and countries created successfully.",
		// "dataCreated": ticketMerchants,
	})
}

func GenerateTicketSeed(c *fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "success",
		"statusCode": 200,
		"message":    "Tickets created successfully.",
		// "dataCreated": ticketMerchants,
	})
}
