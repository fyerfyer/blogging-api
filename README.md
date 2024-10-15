### 1.使用方法

&emsp;&emsp;使用如下语句运行go文件：

```bash
$ go run . # -addr: YourAddr -dsn: YourDSN
```

&emsp;&emsp;使用`curl`执行`CREATE`，`UPDATE`，`DELETE`，`GET`，`GETALL`请求：

* **`POST`**

```bash
curl -X POST http://localhost:4000/posts \
-H 'Content-Type: application/json' \
-d '{"title": "My Second Blog Post", "content": "Content for my second post.", "category": "Technology", "tags": ["Tech", "Go"]}'
```

* **`GETALL`**

```bash
curl -X GET http://localhost:4000/posts
```

&emsp;&emsp;如果要获取特定类别的文章，则如下指定类别：

```bash
curl -X GET "http://localhost:4000/posts?category=Technology"
```

* **`GET`**

```bash
curl -X GET http://localhost:4000/posts/1
```

* **`UPDATE`**

```bash
curl -X PUT http://localhost:4000/posts/1 \
-H 'Content-Type: application/json' \
-d '{"title": "Updated Blog Post Title", "content": "Updated content.", "category": "Lifestyle", "tags": ["Update", "Go"]}'
```

* **`DELETE`**

```bash
curl -X DELETE http://localhost:4000/posts/1
```

### 2.API实现概述

&emsp;&emsp;项目实现的文件结构如下：

```go
/blogging-api
  ├── main.go     // setting up database connecting pool & read client parameters & running server
  ├── models.go   // the database table model
  ├── router.go   // the routers creation
  ├── handler.go  // the handlers for each router
```

&emsp;&emsp;一些使用到的方法如下：

* **依赖注入**：

```go
// main.go
type Application struct {
	Post *PostModel
}

...

app := &Application{
	Post: &PostModel{DB: db},
}

srv := &http.Server{
	Addr:    *addr,
	Handler: app.Routers(),
}
```

* **JSON编解码**

```go
err := json.NewDecoder(r.Body).Decode(&post)
if err != nil {
	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	return
}

...

if err := json.Unmarshal(tagsJSON, &post.Tags); err != nil {
	http.Error(w, "Failed to unmarshal tags", http.StatusInternalServerError)
	return
}

...

if err := json.NewEncoder(w).Encode(post); err != nil {
	http.Error(w, "Failed to encode post", http.StatusInternalServerError)
	return
}
```

