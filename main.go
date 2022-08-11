package main

import (
	"blogapi/blog"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.New()

	//setup the database
	blog.SetupDB()
	//Setup the cache
	blog.ConnectRedisCache()

	//routes and endpoints

	r.POST("/user", blog.UserRegistration)
	r.POST("/login", blog.UserLogin)

	//authorization middleware
	auth := r.Group("/api", blog.AuthMiddleware)

	//Create a blog in database and cache
	auth.POST("/create", blog.CreateBlogPost)

	//updateBlogposts
	auth.POST("/update", blog.UpdateBlogPosts)

	//delete
	auth.DELETE("/delete/:title", blog.DeleteBlogPosts)

	//Get all blogs from database
	r.GET("/fetchdb", blog.FetchFromDatabase)

	//Get a  blog
	r.GET("/:title", blog.FetchByTitle)

	r.Run(":3000")
}
