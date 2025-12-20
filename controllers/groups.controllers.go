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

	group.POST("/", g.CreateGroup())                     // create new Grp //done!!
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

		group, err := g.groupService.CreateGroupService(body.Name, body.Members)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}

		c.JSON(http.StatusAccepted, gin.H{
			"data": group,
		})

	}
}

func (g *GroupControllers) AddGroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {

		groupId := c.Param("groupId")

		type Body struct {
			Members []string `json:"members" binding:"required"`
		}

		var body Body

		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		group, err := g.groupService.AddGroupMembers(groupId, body.Members)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusAccepted, gin.H{
			"data": group,
		})
	}
}

func (g *GroupControllers) GetGroupDetails() gin.HandlerFunc {
	return func(c *gin.Context) {

		groupId := c.Param("groupId")

		group, err := g.groupService.GroupDetails(groupId)

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": group,
		})
	}
}

func (g *GroupControllers) GetGroups() gin.HandlerFunc {
	return func(c *gin.Context) {

		groups, err := g.groupService.GetAllGroups()

		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": groups,
		})
	}
}

func (g *GroupControllers) AddExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

		groupId := c.Param("groupId")

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

		expense, err := g.expensesService.AddExpensesService(body.Description, body.PaidBy, body.Amount, body.SplitBetween, groupId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": expense,
		})
	}
}

func (g *GroupControllers) GetExpenses() gin.HandlerFunc {
	return func(c *gin.Context) {

		groupId := c.Param("groupId")

		expenses, err := g.expensesService.GetAllExpensesService(groupId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": expenses,
		})
	}
}

func (g *GroupControllers) GetBalance() gin.HandlerFunc {
	return func(c *gin.Context) {

		groupId := c.Param("groupId")

		balance, err := g.balanceService.GetGroupBalance(groupId)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"balance": balance,
		})
	}
}
