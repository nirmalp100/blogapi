# blogapi
Build the following back-end APIs for a blogging application in Golang using gin-gonic.
An API to create a blog post.
The post will consist of title, content, author name and publish date time.
The blog post must be stored in PostgreSQL
An API to get the blog posts.
The blog posts must be read from PostgreSQL
An API to get a blog post by post ID
The post must be fetched from Redis
If the post is not found in Redis, read from PostgreSQL, store in Redis and then return the post
An API to update blog post
The blog post must be updated both in Redis and PostgreSQL
Question - What order must the updates happen and why?
An API to delete a blog post
The blog post must be deleted in both Redis and PostgreSQL



Once the above is done,
Build following back-end APIs for User Authentication
An API to register user.
Collect and store name, email and password
Save the user in DB with an ID field
Password must be stored securely (hash + salt)
An API to login user
Accept email and password
Authenticate using (password + salt == hash)
After the authentication is complete, update the blogging APIs with the following changes
Use JWT to secure the Create, Update and Delete APIs
If the request does not contain a JWT token, the APIs should respond with 401 error code.
Add an author ID field to the posts table
This should be set to the current user's ID who is creating the post
Update and delete must be done by the user who created the post
Anyone can list and read the posts
