package controllers

import "github.com/gin-gonic/gin"

type GroupControllers struct{}

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

	}
}

func (g *GroupControllers) AddGroupMembers() gin.HandlerFunc {
	return func(c *gin.Context) {

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
