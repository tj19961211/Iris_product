package controllers

import (
	"Iris_product/datamodels"
	"Iris_product/rabbitmq"
	"Iris_product/services"
	"encoding/json"
	"strconv"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
	"github.com/kataras/iris/sessions"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	RabbitMQ       *rabbitmq.RabbitMQ
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	product, err := p.ProductService.GetProductByID(4)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() []byte {
	productString := p.Ctx.URLParam("productID")
	userString := p.Ctx.GetCookie("uid")
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	userID, err := strconv.ParseInt(userString, 10, 64)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 创建消息体
	message := datamodels.NewMessage(userID, productID)
	//类型转化
	byteMessage, err := json.Marshal(message)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	err = p.RabbitMQ.PublishSimple(string(byteMessage))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}
	return []byte("true")

	// product, err := p.ProductService.GetProductByID(int64(productID))
	// if err != nil {
	// 	p.Ctx.Application().Logger().Debug(err)
	// }
	// var orderID int64
	// showMessage := "抢购失败！"
	// //判断商品数量是否满足需求
	// if product.ProductNum > 0 {
	// 	//扣除商品数量
	// 	product.ProductNum -= 1
	// 	err := p.ProductService.UpdateProduct(product)
	// 	if err != nil {
	// 		p.Ctx.Application().Logger().Debug(err)
	// 	}
	// 	//创建订单
	// 	userID, err := strconv.Atoi(userString)
	// 	if err != nil {
	// 		p.Ctx.Application().Logger().Debug(err)
	// 	}

	// 	order := &datamodels.Order{
	// 		UserId:      int64(userID),
	// 		ProductId:   int64(productID),
	// 		OrderStatus: datamodels.OrderSuccess,
	// 	}
	// 	//新建订单
	// 	orderID, err = p.OrderService.InsertOrder(order)
	// 	if err != nil {
	// 		p.Ctx.Application().Logger().Debug(err)
	// 	} else {
	// 		showMessage = "抢购成功！"
	// 	}
	// }

	// return mvc.View{
	// 	Layout: "shared/productLayout.html",
	// 	Name:   "product/result.html",
	// 	Data: iris.Map{
	// 		"orderID":     orderID,
	// 		"showMessage": showMessage,
	// 	},
	// }
}
