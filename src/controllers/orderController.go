package controllers

import (
	"ankasa-be/src/helpers"
	"ankasa-be/src/middlewares"
	"ankasa-be/src/models"
	"ankasa-be/src/services"
	"context"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/xendit/xendit-go/v6/payment_method"
	"github.com/xendit/xendit-go/v6/payment_request"
)

func GetAllOrders(c *fiber.Ctx) error {
	user_id, err := middlewares.JWTAuthorize(c, "customer")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	user := models.SelectUserfromID(int(user_id))
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "User not found",
		})
	}

	customer := models.SelectCustomerfromUserID(int(user_id))
	if customer.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Customer not found",
		})
	}

	orders := models.SelectOrdersbyCustomerID(int(customer.ID))
	if len(orders) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "success",
			"statusCode": fiber.StatusOK,
			"message":    "orders unavailable",
		})
	}

	resultOrders := make([]map[string]interface{}, len(orders))
	for i, order := range orders {
		resultOrders[i] = map[string]interface{}{
			"id":                     order.ID,
			"created_at":             order.CreatedAt,
			"updated_at":             order.UpdatedAt,
			"merchant_name":          order.Ticket.Merchant.Name,
			"merchant_image":         order.Ticket.Merchant.Image,
			"departure_schedule":     order.Ticket.Departure.Schedule,
			"departure_country_code": order.Ticket.Departure.City.Country.Code,
			"arrival_country_code":   order.Ticket.Arrival.City.Country.Code,
			"payment_method":         order.PaymentMethod,
			"payment_status":         order.PaymentStatus,
			"payment_url":            order.PaymentURL,
			"is_active":              order.IsActive,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       resultOrders,
	})
}

func GetBookingPass(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid ID format",
		})
	}

	order := models.SelectOrderbyID(&id)
	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Order not found",
		})
	}

	resultBookings := make([]map[string]interface{}, len(order.Passengers))
	for i, passenger := range order.Passengers {
		resultBookings[i] = map[string]interface{}{
			"id":                     order.ID,
			"created_at":             order.CreatedAt,
			"merchant_name":          order.Ticket.Merchant.Name,
			"merchant_image":         order.Ticket.Merchant.Image,
			"departure_schedule":     order.Ticket.Departure.Schedule,
			"departure_country_code": order.Ticket.Departure.City.Country.Code,
			"arrival_country_code":   order.Ticket.Arrival.City.Country.Code,
			"gate":                   order.Ticket.Gate,
			"class_name":             order.Ticket.Class.Name,
			"passenger_name":         passenger.Name,
			"passenger_seat_code":    passenger.Seat.Code,
			"passenger_category":     passenger.Category,
		}
	}

	if len(resultBookings) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Booking unavailable",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":     "success",
		"statusCode": fiber.StatusOK,
		"data":       resultBookings,
	})
}

func CreatePaymentOrder(c *fiber.Ctx) error {
	user_id, err := middlewares.JWTAuthorize(c, "customer")
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			return c.Status(fiberErr.Code).JSON(fiber.Map{
				"status":     fiberErr.Message,
				"statusCode": fiberErr.Code,
				"message":    fiberErr.Message,
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	user := models.SelectUserfromID(int(user_id))
	if user.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "User not found",
		})
	}

	customer := models.SelectCustomerfromUserID(int(user_id))
	if customer.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Customer not found",
		})
	}

	var requestOrder models.OrderRequest
	if err := c.BodyParser(&requestOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	orderReq := middlewares.XSSMiddleware(&requestOrder).(*models.OrderRequest)
	order := models.Order{
		TicketID:      orderReq.TicketID,
		Passengers:    orderReq.Passengers,
		TotalPrice:    orderReq.TotalPrice,
		PaymentMethod: orderReq.PaymentMethodType,
		PaymentStatus: payment_request.PAYMENTREQUESTSTATUS_PENDING.String(),
	}

	ticket, _ := models.SelectTicketById(int(order.TicketID))
	if ticket.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Ticket not found",
		})
	}

	if ticket.Stock == 0 || ticket.Stock < len(order.Passengers) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Ticket cannot booking",
		})
	}

	order.CustomerID = customer.ID
	order.TotalPassengers = len(order.Passengers)
	for _, passenger := range order.Passengers {
		seat := models.SelectSeatBooking(int(passenger.SeatID))
		if seat.ID != 0 {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":     "forbidden",
				"statusCode": 403,
				"message":    "Seat already booking",
			})
		}
	}

	var referenceID *string
	if orderReq.PaymentMethodType != "INTERNAL" {
		var err error
		if orderReq.PaymentMethodType == "EWALLET" {
			referenceID, err = services.EWalletPaymentMethod(orderReq.PaymentMethodEWallet)
		} else {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":     "forbidden",
				"statusCode": 403,
				"message":    "Another payment method still in development",
			})
		}

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}
	} else {
		order.PaymentStatus = payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED.String()
		now := time.Now()
		order.PaidAt = &now
		order.IsActive = true

		for _, passenger := range order.Passengers {
			if err := models.UpdateSeatIsBooking(int(passenger.SeatID), true); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		ticket.Stock -= len(order.Passengers)
		if err := models.UpdateTicketSingle(int(ticket.ID), "stock", ticket.Stock); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		err := models.CreateOrder(&order)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"status":     "created",
			"statusCode": fiber.StatusCreated,
			"message":    "Order completed and payment succeeded but payment ankasapay (wallet integration) still in development.",
		})
	}

	if referenceID == nil && orderReq.PaymentMethodType != "INTERNAL" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Payment can't completed",
		})
	}
	order.PaymentID = *referenceID

	for _, passenger := range order.Passengers {
		if err := models.UpdateSeatIsBooking(int(passenger.SeatID), true); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}
	}

	ticket.Stock -= len(order.Passengers)
	if err := models.UpdateTicketSingle(int(ticket.ID), "stock", ticket.Stock); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	err = models.CreateOrder(&order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "internal server error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":     "created",
		"statusCode": fiber.StatusCreated,
		"message":    "Order has been created successfully",
	})
}

func HandlePaymentMethodCallback(c *fiber.Ctx) error {
	var requestPayload payment_method.PaymentMethodCallback
	if err := c.BodyParser(&requestPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	paymentMethod := middlewares.XSSMiddleware(&requestPayload).(*payment_method.PaymentMethodCallback)
	data, ok := paymentMethod.GetDataOk()
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Data payment method not found",
		})
	}

	order := models.SelectOrderSingle("payment_id", data.ReferenceId)
	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Order not found",
		})
	}

	status := *data.Status
	if status == payment_method.PAYMENTMETHODSTATUS_ACTIVE {

		referenceId := "ankasa-" + helpers.GenerateString(16) + "-" + time.Now().Format("20060102150405")
		paymentRequestParameters := payment_request.PaymentRequestParameters{
			ReferenceId:     &referenceId,
			Amount:          &order.TotalPrice,
			Currency:        "IDR",
			PaymentMethodId: &data.Id,
		}

		paymentRequest, _, err := services.Client.PaymentRequestApi.
			CreatePaymentRequest(context.Background()).
			PaymentRequestParameters(paymentRequestParameters).
			Execute()

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "Internal Server Error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		if err := models.UpdateOrderSingle(int(order.ID), "external_id", &referenceId); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		if paymentRequest.Status == payment_request.PAYMENTREQUESTSTATUS_REQUIRES_ACTION {
			for _, action := range paymentRequest.Actions {
				if action.UrlType == "WEB" {
					url, ok := action.GetUrlOk()
					if ok {
						if err := models.UpdateOrderSingle(int(order.ID), "payment_url", &url); err != nil {
							return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
								"status":     "internal server error",
								"statusCode": fiber.StatusInternalServerError,
								"message":    err.Error(),
							})
						}
					}
				}
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment method activated",
			"data":       paymentRequest,
		})

	} else if status == payment_method.PAYMENTMETHODSTATUS_EXPIRED {

		if !order.IsActive && order.PaymentStatus == payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED.String() {
			if err := models.UpdateOrderSingle(int(order.ID), "is_active", true); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment method expired",
		})

	} else if status == payment_method.PAYMENTMETHODSTATUS_INACTIVE || status == payment_method.PAYMENTMETHODSTATUS_PENDING {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment method could be pending or inactive",
		})

	} else if status == payment_method.PAYMENTMETHODSTATUS_REQUIRES_ACTION {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment method require action still in development",
		})

	} else if status == payment_method.PAYMENTMETHODSTATUS_FAILED {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    "Payment method failed",
		})

	} else {

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Payment method not valid",
		})

	}
}

func HandlePaymentRequestCallback(c *fiber.Ctx) error {
	var requestPayload payment_request.PaymentCallback
	if err := c.BodyParser(&requestPayload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":     "bad request",
			"statusCode": 400,
			"message":    "Invalid request body",
		})
	}

	payment := middlewares.XSSMiddleware(&requestPayload).(*payment_request.PaymentCallback)
	data, ok := payment.GetDataOk()
	if !ok {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Data payment not found",
		})
	}

	order := models.SelectOrderSingle("external_id", data.ReferenceId)
	if order.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Order not found",
		})
	}

	ticket, _ := models.SelectTicketById(int(order.TicketID))
	if ticket.ID == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":     "not found",
			"statusCode": 404,
			"message":    "Ticket not found",
		})
	}

	status, _ := payment_request.NewPaymentRequestStatusFromValue(data.Status)
	if *status == payment_request.PAYMENTREQUESTSTATUS_SUCCEEDED {

		order.PaymentStatus = string(*status)
		now := time.Now()
		order.PaidAt = &now
		order.IsActive = true

		if _, err := models.UpdateOrderById(int(order.ID), *order); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment succeeded",
		})

	} else if *status == payment_request.PAYMENTREQUESTSTATUS_PENDING {

		if order.PaymentStatus != "PENDING" {
			if err := models.UpdateOrderSingle(int(order.ID), "status", string(*status)); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		if order.IsActive {
			if err := models.UpdateOrderSingle(int(order.ID), "is_active", false); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment pending",
		})

	} else if *status == payment_request.PAYMENTREQUESTSTATUS_VOIDED || *status == payment_request.PAYMENTREQUESTSTATUS_CANCELED {

		order.PaymentStatus = payment_request.PAYMENTREQUESTSTATUS_CANCELED.String()
		if err := models.UpdateOrderSingle(int(order.ID), "status", order.PaymentStatus); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		if order.IsActive {
			if err := models.UpdateOrderSingle(int(order.ID), "is_active", false); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		for _, passenger := range order.Passengers {
			if err := models.UpdateSeatIsBooking(int(passenger.SeatID), false); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		ticket.Stock += len(order.Passengers)
		if err := models.UpdateTicketSingle(int(ticket.ID), "stock", ticket.Stock); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment cancelled",
		})

	} else if *status == payment_request.PAYMENTREQUESTSTATUS_EXPIRED {

		if err := models.UpdateOrderSingle(int(order.ID), "status", string(*status)); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		if order.IsActive {
			if err := models.UpdateOrderSingle(int(order.ID), "is_active", false); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		for _, passenger := range order.Passengers {
			if err := models.UpdateSeatIsBooking(int(passenger.SeatID), false); err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":     "internal server error",
					"statusCode": fiber.StatusInternalServerError,
					"message":    err.Error(),
				})
			}
		}

		ticket.Stock += len(order.Passengers)
		if err := models.UpdateTicketSingle(int(ticket.ID), "stock", ticket.Stock); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":     "internal server error",
				"statusCode": fiber.StatusInternalServerError,
				"message":    err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment expired",
		})

	} else if *status == payment_request.PAYMENTREQUESTSTATUS_REQUIRES_ACTION || *status == payment_request.PAYMENTREQUESTSTATUS_AWAITING_CAPTURE {

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":     "ok",
			"statusCode": fiber.StatusOK,
			"message":    "Payment require action or awaiting capture still in development",
		})

	} else if *status == payment_request.PAYMENTREQUESTSTATUS_FAILED {

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":     "Internal Server Error",
			"statusCode": fiber.StatusInternalServerError,
			"message":    "Payment failed",
		})

	} else {

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":     "forbidden",
			"statusCode": 403,
			"message":    "Payment not valid",
		})

	}
}
