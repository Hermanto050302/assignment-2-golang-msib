package handler

import (
	"assignment-2-golang-msib/database"
	"assignment-2-golang-msib/dto"
	"assignment-2-golang-msib/pkg/helpers"
	"assignment-2-golang-msib/repository/item_repository/item_pg"
	"assignment-2-golang-msib/repository/order_repository/order_pg"
	"assignment-2-golang-msib/service"

	"github.com/gin-gonic/gin"
)

func StartApp() {
	database.InitiliazeDatabase()

	db := database.GetDatabaseInstance()

	itemRepo := item_pg.NewItemPG(db)

	itemService := service.NewItemService(itemRepo)

	orderRepo := order_pg.NewOrderPG(db)

	orderService := service.NewOrderService(orderRepo, itemService)

	r := gin.Default()

	r.POST("/orders", func(ctx *gin.Context) {
		var orderRequest dto.NewOrderRequest

		if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
			ctx.JSON(422, gin.H{
				"errMessage": "err brader",
			})
			return
		}

		result, err := orderService.CreateOrder(orderRequest)

		if err != nil {
			ctx.JSON(500, gin.H{
				"errMessage": "err internal brader",
			})
			return
		}
		ctx.JSON(201, result)

	})

	//get all orders
	r.GET("/orders", func(ctx *gin.Context) {
		orders, err := orderService.GetAllOrders()
		if err != nil {
			ctx.JSON(400, gin.H{
				"errMessage": err.Error(),
			})
			return
		}
		ctx.JSON(200, orders)
	})

	r.PUT("/orders/:orderId", func(ctx *gin.Context) {
		var orderRequest dto.NewOrderRequest

		if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
			ctx.JSON(422, gin.H{
				"errMessage": "unprocessible entity",
			})
			return
		}

		orderId, err := helpers.GetParamId(ctx, "orderId")

		if err != nil {
			ctx.JSON(400, gin.H{
				"errMessage": err.Error(),
			})
			return
		}

		updatedOrder, err := orderService.UpdateOrder(orderId, orderRequest)

		if err != nil {
			ctx.JSON(400, gin.H{
				"errMessage": err.Error(),
			})
			return
		}

		ctx.JSON(updatedOrder.Code, updatedOrder)

	})

	//delete order
	r.DELETE("/orders/:orderId", func(ctx *gin.Context) {
		//delete order
		orderId, err := helpers.GetParamId(ctx, "orderId")

		if err != nil {
			ctx.JSON(400, gin.H{
				"errMessage": err.Error(),
			})
			return
		}

		_, err = orderService.DeleteOrder(orderId)

		if err != nil {
			ctx.JSON(400, gin.H{
				"errMessage": err.Error(),
			})
			return
		}

		ctx.JSON(204, nil)
	})


	r.Run(":8080")
}
