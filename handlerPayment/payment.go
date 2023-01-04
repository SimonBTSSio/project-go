package handlerPayment

import (
	"go-project/broadcast"
	"go-project/payment"
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Message struct {
	Text  string
	Name  string
	Price float64
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type paymentHandler struct {
	paymentService payment.Service
	bc             broadcast.Broadcaster
}

func NewPaymentHandler(paymentService payment.Service, bc broadcast.Broadcaster) *paymentHandler {
	return &paymentHandler{paymentService, bc}
}

func (th *paymentHandler) Create(c *gin.Context) {
	var input payment.InputPayment
	err := c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	newPayment, err := th.paymentService.Create(input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	th.bc.Submit(Message{
		Text:  "Nouveau paiement vient d'etre fait",
		Name:  newPayment.Product.Name,
		Price: newPayment.Product.Price,
	})

	response := &Response{
		Success: true,
		Message: "New payment created",
		Data:    newPayment,
	}
	c.JSON(http.StatusCreated, response)
}

func (th *paymentHandler) GetAll(c *gin.Context) {
	payments, err := th.paymentService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payments,
	})
}

func (th *paymentHandler) GetById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	payment, err := th.paymentService.GetById(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Data:    payment,
	})
}

func (th *paymentHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	var input payment.InputPayment
	err = c.ShouldBindJSON(&input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Cannot extract JSON body",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	payment, err := th.paymentService.Update(id, input)
	if err != nil {
		response := &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		}
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := &Response{
		Success: true,
		Message: "New payment created",
		Data:    payment,
	}
	c.JSON(http.StatusCreated, response)
}

func (th *paymentHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Wrong id parameter",
			Data:    err.Error(),
		})
		return
	}

	err = th.paymentService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, &Response{
			Success: false,
			Message: "Something went wrong",
			Data:    err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, &Response{
		Success: true,
		Message: "Payment successfully deleted",
	})
}

func (th *paymentHandler) Stream(c *gin.Context) {
	listener := make(chan interface{})

	th.bc.Register(listener)
	defer th.bc.Unregister(listener)

	clientGone := c.Request.Context().Done()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			serviceMsg, ok := message.(Message)
			if !ok {
				c.SSEvent("message", message)
				return false
			}
			c.SSEvent("message", serviceMsg.Text+", "+"Produit : "+serviceMsg.Name+", "+"Prix : "+strconv.FormatFloat(serviceMsg.Price, 'E', -1, 64))
			return true
		}
	})
}
