package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const PORT = ":8080"

type Orders struct {
	OrderID      int       `json:"order_id"`
	CustomerName string    `json:"customer_name"`
	OrderedAt    time.Time `json:"ordered_at"`
	Item         []Items   `json:"items"`
}

type Items struct {
	ItemID      int    `json:"item_id" example:"1"`
	ItemCode    string `json:"item_code" example:"TEST01"`
	Description string `json:"description" example:"DESC TEST ITEM01"`
	Quantity    int    `json:"quantity" example:"5"`
	OrderId     int    `json:"order_id" example:"1"`
}

type Request struct {
	// OrderID      int     `json:"order_id,omitempty"`
	CustomerName string  `json:"customer_name" example:"Sigit Setiawan"`
	Item         []Items `json:"items"`
}

type Response struct {
	CustomerName string    `json:"customer_name" example:"Sigit Setiawan"`
	OrderID      int       `json:"order_id" example:"1"`
	OrderedAt    time.Time `json:"ordered_at" example:"2022-10-04T10:09:55.6868076+07:00"`
	Item         []Items   `json:"items"`
}

type orderService struct {
	db  []*Orders
	db2 []*Items
}

type OrderIface interface {
	CreateOrder(c *gin.Context)
	GetOrders(c *gin.Context)
	// GetOrderID(c *gin.Context) int
	UpdateOrder(c *gin.Context)
	DeleteOrder(c *gin.Context)
}

func NewOrderService(db []*Orders, db2 []*Items) OrderIface {
	return &orderService{
		db:  db,
		db2: db2,
	}
}

// Create One Order
// @Summary Create New Order
// @Description Create New Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Param data body Request true "Order"
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /orders [post]
func (o *orderService) CreateOrder(c *gin.Context) {
	if c.Request.Method == http.MethodPost {
		var req Request
		var res Response
		var item []Items

		// Decode Request
		err := json.NewDecoder(c.Request.Body).Decode(&req)
		if err != nil {
			fmt.Println(err.Error())
		}

		orderID := o.GetOrderID(c)

		// Insert Order
		o.db = append(o.db, &Orders{
			OrderID:      orderID,
			CustomerName: req.CustomerName,
			OrderedAt:    time.Now(),
		})

		// Create Response
		for i := range req.Item {
			itemID := o.GetItemID(c)

			// Insert Item
			o.db2 = append(o.db2, &Items{
				ItemID:      itemID,
				ItemCode:    req.Item[i].ItemCode,
				Description: req.Item[i].Description,
				Quantity:    req.Item[i].Quantity,
				OrderId:     orderID,
			})

			item = append(item, Items{
				ItemID:      itemID,
				ItemCode:    req.Item[i].ItemCode,
				Description: req.Item[i].Description,
				Quantity:    req.Item[i].Quantity,
				OrderId:     orderID,
			})
		}

		res.CustomerName = req.CustomerName
		res.OrderID = orderID
		res.OrderedAt = o.db[0].OrderedAt
		res.Item = item

		c.JSON(http.StatusCreated, res)
	}
}

// Get Order
// @Summary Get All Order
// @Description Get All Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /orders [get]
func (o *orderService) GetOrders(c *gin.Context) {
	var res []Response

	for _, v := range o.db {
		res = append(res, Response{
			CustomerName: v.CustomerName,
			OrderID:      v.OrderID,
			OrderedAt:    v.OrderedAt,
		})
	}

	for i := 0; i < len(res); i++ {
		for j := 0; j < len(o.db2); j++ {
			if o.db2[j].OrderId == res[i].OrderID {
				res[i].Item = append(res[i].Item, Items{
					ItemID:      o.db2[j].ItemID,
					ItemCode:    o.db2[j].ItemCode,
					Description: o.db2[j].Description,
					Quantity:    o.db2[j].Quantity,
					OrderId:     o.db2[j].OrderId,
				})
			}
		}
	}

	c.JSON(http.StatusOK, res)
}

// Update Order
// @Summary Update Order
// @Description Update Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /order/:id [put]
func (o *orderService) UpdateOrder(c *gin.Context) {
	var req Request
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		panic(err)
	}

	// Get id order
	id := c.Params.ByName("id")
	orderId, _ := strconv.Atoi(id)

	for _, v := range o.db {
		if v.OrderID == orderId {
			v.CustomerName = req.CustomerName
		}
	}

	for i, v := range o.db2 {
		if v.OrderId == orderId {
			if v.ItemID == orderId {
				v.ItemCode = req.Item[i].ItemCode
				v.Description = req.Item[i].Description
				v.Quantity = req.Item[i].Quantity
			}
		}
	}

	c.JSON(http.StatusOK, req)
}

// Delete Order
// @Summary Delete Order
// @Description Delete Order
// @Tags Orders
// @Accept  */*
// @Produce  json
// @Success 200 {object} Response
// @Failure 500 {object} string "error"
// @Router /order/:id [delete]
func (o *orderService) DeleteOrder(c *gin.Context) {
	// Get id order
	id := c.Params.ByName("id")
	orderId, _ := strconv.Atoi(id)

	for i := range o.db {
		if orderId == o.db[i].OrderID {
			o.db[i] = &Orders{}
		}
	}

	o.db = RemoveIndex(o.db, 0)
	c.JSON(http.StatusOK, o.db)
}

func (o *orderService) GetOrderID(c *gin.Context) int {
	max := 0
	for _, v := range o.db {
		if v.OrderID > max {
			max = v.OrderID
		}
	}

	return max + 1
}

func (o *orderService) GetItemID(c *gin.Context) int {
	max := 0
	for _, v := range o.db2 {
		if v.ItemID > max {
			max = v.ItemID
		}
	}
	return max + 1
}

func RemoveIndex(s []*Orders, index int) []*Orders {
	return append(s[:index], s[index+1:]...)
}

// @title Gin Swagger Example API
// @version 2.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {
	var db []*Orders
	var db2 []*Items

	orderSvc := NewOrderService(db, db2)

	router := gin.Default()

	p := router.Group("/ping")
	{
		p.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status": "pong",
			})
		})
	}

	router.POST("/orders", orderSvc.CreateOrder)
	router.GET("/orders", orderSvc.GetOrders)
	router.PUT("/order/:id", orderSvc.UpdateOrder)
	router.DELETE("/order/:id", orderSvc.DeleteOrder)

	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json")
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	router.Run(PORT)
	fmt.Println("Server running at port", PORT)
}
