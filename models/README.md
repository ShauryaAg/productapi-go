#### Database

There are 3 models in the database:

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

#### Indexing

Indexing is done by defining the custom Tag `mongo` against the fields that need to be indexed/unique.

The `mongo` is parsed under `utils/db.go` via **`Migrate`** function, which called under `models/db/db.go` by **`InitDatabase`** function whilst initializing the database.

The `Migrate` function checks if the collection/index exists, if not it creates it, and returns all the collections in a map.
