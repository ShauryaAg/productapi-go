## ProductAPI

ProductAPI is a RESTful API for managing products created as IDT.NET task

The API consists of three models:

**User**

```go
type User struct {
	Id       primitive.ObjectID
	Name     string
	Password string
	Email    string
}
```

**Product**

```go
type Product struct {
	Id                primitive.ObjectID
	Name              string
	Description       string
	ThumbnailImageUrl string
	Reviews           []Review
	RatingSum         float64
	RatingCount       int
}
```

**Review**

```go
type Review struct {
	Id       primitive.ObjectID
	Text     string
	Rating   int
	Reviewer User
}
```

### Setup locally

- Clone the repo using `git clone https://github.com/ShauryaAg/ProductAPI.git`
- Move into the project folder `cd ProductAPI/`
- Create a `.env` file in the project folder

```
PORT=<PORT>
MONGODB_URI=<Your Mongo URI>
SECRET=<Secret>
SEED_DB=<true/false>
```

##### Using Docker <a href="https://www.docker.com/"> <img alt="ProductAPI" src="https://www.docker.com/sites/default/files/d8/styles/role_icon/public/2019-07/vertical-logo-monochromatic.png" width="50" /> </a>

- Run using `docker-compose up --build`

### **OR**

##### Using Golang <a href="https://golang.org/"> <img alt="ProductAPI" src="https://golang.org/lib/godoc/images/go-logo-blue.svg" width="50" /> </a>

- Install the dependecies using `go mod download`
- Run using `go run server.go`
  > Note: you need to have the mongo database running on your local machine for this to work

###### Seeding the Database

> Set the `SEED_DB` environment variable to `true` to seed the database with sample data.

### Endpoint Usage

#### Auth

- `/api/auth/register`

  - _Allowed Methods:_ `POST`
  - _Accepted Fields:_ `{name, email, password}`
  - _Returns:_ `{id, email, token}`

- `/api/auth/login`

  - _Allowed Methods:_ `POST`
  - _Accepted Fields:_ `{email, password}`
  - _Returns:_ `{id, name, email, token}`

- `/api/auth/user`

  - _Allowed Methods:_ `GET`
  - _Authorization:_ `Bearer <Token>`
  - _Returns:_ `User details`

#### Products

- `api/product`

  - _Allowed Methods:_ `GET, POST`

    - `GET`

      - Query params: `{q, page, limit}`
      - _Returns:_ `All products matching the query at the given page`

    - `POST`

      - _Authorization:_ `Bearer <Token>`
      - _Accepted Fields:_ `{name, description, thumbnailImageUrl}`
      - _Returns:_ `Product Details`

#### Reviews

- `api/review/<product_id>`

  - _Allowed Methods:_ `POST`
  - _Authorization:_ `Bearer <Token>`
  - _Accepted Fields:_ `{text, rating}`
  - _Returns:_ `Review Details`
