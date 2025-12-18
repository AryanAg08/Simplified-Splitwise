package controllers

import (
	"net/http"

	"github.com/AryanAg08/Simplified-Splitwise/services"
	"github.com/gin-gonic/gin"
)

type GroupControllers struct {
	balanceService  services.BalanceService
	expensesService services.ExpensesSerive
	groupService    services.GroupsService
}

func (g *GroupControllers) InitGroupController(router *gin.Engine) {
	group := router.Group("/group")

	group.POST("/", g.CreateGroup())                     // create new Grp
	group.POST("/:groupId/members", g.AddGroupMembers()) // adding members
	group.GET("/:groupId", g.GetGroupDetails())          // grp details
	group.GET("/", g.GetGroups())                        // All Groups

	// expenses
	group.POST("/:groupId/expenses", g.AddExpenses()) // add expenses
	group.GET("/:groupId/expenses", g.GetExpenses())  // List all expenses

	// balance
	group.GET("/:groupId/balance", g.GetBalance()) // get balances
}

func (g *GroupControllers) CreateGroup() gin.HandlerFunc {
	return func(c *gin.Context) {

		type Body struct {
			Name    string   `json:"name" binding:"required"`
			Members []string `json:"members"`
		}

		var body Body

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

	}
}

func (g *GroupControllers) AddGroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {
		// groupId := c.Param("groupId")

		type Body struct {
			Members []string `json:"members" binding:"required"`
		}

		var body Body

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}

	}
}

func (g *GroupControllers) GetGroupDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GroupControllers) GetGroups() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GroupControllers) AddExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

		// groupId := c.Param("groupId")

		type Body struct {
			Description  string   `json:"description" binding:"required"`
			PaidBy       string   `json:"paidBy" binding:"required"`
			Amount       float64  `json:"amount" binding:"required"`
			SplitBetween []string `json:"splitBetween" binding:"required"`
		}

		var body Body

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	}
}

func (g *GroupControllers) GetExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}

func (g *GroupControllers) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
