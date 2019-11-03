# Go GraphQL MongoDB

Get Authors
[query={authors{id,firstname,lastname}}](http://localhost:8088/graphql?query={authors{id,firstname,lastname}})

Create Author
[query=mutation+_{createAuthor(firstname:"Bruno",lastname:"Ramalho"){id,firstname,lastname}}](http://localhost:8088/graphql?query=mutation+_{createAuthor(firstname:"Bruno",lastname:"Ramalho"){id,firstname,lastname}})

Get Author
[query={author(id:"123"){id,firstname,lastname}}](http://localhost:8088/graphql?query={author(id:"123"){id,firstname,lastname}})

Get Blogs
[](http://localhost:8088/graphql?query={blogs{id,title,body,author}})

Create Blog
[query=mutation+_{createBlog(authorID:"5dbf012c3778b9648b8609a1",title:"Hello_World",body:"It_works!"){id,title,body}}](http://localhost:8088/graphql?query=mutation+_{createBlog(authorID:"5dbf012c3778b9648b8609a1",title:"Hello_World",body:"It_works!"){id,title,body}})

Get Blog
[query={blog(id:"123"){id,title,body,author}}](http://localhost:8088/graphql?query={blog(id:"123"){id,title,body,author}})
