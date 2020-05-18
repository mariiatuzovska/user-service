# user-service

## DB table

```GO
type User struct {
	gorm.Model
	Name    string `gorm:"not_null"` // Required
	Age     int    `gorm:"not_null"` // Required
	Email   string `gorm:"not_null"` // Required
	Address string
}
```

## Configuration file

```JSON
{
    "DBContext": {
        "Shema": "postgres",
        "User": "mari",
        "Password": "postgres",
        "Host": "127.0.0.1",
        "Port": "5432"
    },
    "APIContext": {
        "Host": "127.0.0.1",
        "Port": "8080"
    }
}
```

## Request sample 

```JSON
{
	"Name":"Mari",
	"Age":21, 
	"Email":"marriia@gmail.com"
}
```