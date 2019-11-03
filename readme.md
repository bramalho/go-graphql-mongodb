# Go GraphQL DB

Get Authors
<http://localhost:8088/graphql?query={authors{id,firstname,lastname}}>

Create Author
<http://localhost:8088/graphql?query=mutation+_{createAuthor(firstname:"Bruno",lastname:"Ramalho"){id,firstname,lastname}}>

Get Author
<http://localhost:8088/graphql?query={author(id:"123"){id,firstname,lastname}}>

Get Blogs
<http://localhost:8088/graphql?query={blogs{id,title,body,author}}>

Create Blog
<http://localhost:8088/graphql?query=mutation+_{createBlog(authorID:"5dbf012c3778b9648b8609a1",title:"Hello World",body:"It works!"){id,title,body}}>

Get Blog
<http://localhost:8088/graphql?query={blog(id:"123"){id,title,body,author}}>
