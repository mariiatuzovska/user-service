package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/labstack/echo"
	"github.com/urfave/cli"
)

var (
	ServiceName = "user-service"
	Version     = "0.0.1"
)

func main() {
	app := cli.NewApp()
	app.Name = ServiceName
	app.Usage = "command line client"
	app.Description = "Users API service"
	app.Version = Version
	app.Authors = []cli.Author{cli.Author{Name: "Tuzovska Mariia"}}
	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "starting service via http",
			Action: func(c *cli.Context) error {
				config, err := NewConfiguration(c.String("config"))
				if err != nil {
					log.Fatalln(err)
				}
				srv, err := NewService(config)
				if err != nil {
					log.Fatalln(err)
				}
				if c.String("host") != "" {
					config.APIContext.Host = c.String("host")
				}
				if c.String("port") != "" {
					config.APIContext.Port = c.String("port")
				}
				log.Println(fmt.Sprintf("| General | Starting service at %s:%s%s",
					config.APIContext.Host, config.APIContext.Port, getURL()))
				return srv.Start(fmt.Sprintf("%s:%s", config.APIContext.Host, config.APIContext.Port))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "host",
					Usage: "Runs service with this host",
				},
				&cli.StringFlag{
					Name:  "port",
					Usage: "Runs service with this port",
				},
				&cli.StringFlag{
					Name:  "config",
					Usage: "Path to configuration file",
					Value: "./user-configuration.json",
				},
			},
		},
	}
	app.Run(os.Args)
}

func getURL() string {
	s := strings.Split(Version, ".")
	return fmt.Sprintf("/api/v%s%s%s", s[0], s[1], s[2])
}

type Service struct {
	*gorm.DB
	*echo.Echo
}

type Configuration struct {
	DBContext struct {
		Shema, User, Password, Host, Port string
	}
	APIContext struct {
		Host string
		Port string
	}
}

type User struct {
	gorm.Model
	Name    string `gorm:"not_null"`
	Age     int    `gorm:"not_null"`
	Email   string `gorm:"not_null"`
	Address string
}

func NewConfiguration(configPath string) (*Configuration, error) {
	config := new(Configuration)
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, config)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func NewService(config *Configuration) (*Service, error) {
	db, err := gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.DBContext.User, config.DBContext.Password, config.DBContext.Host, config.DBContext.Port, config.DBContext.Shema))
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&User{})
	srv := &Service{db, echo.New()}
	srv.HidePort = true
	srv.HideBanner = true
	srv.GET(getURL(), srv.GetUsers)
	srv.POST(getURL(), srv.CreateUser)
	srv.PUT(getURL(), srv.UpdateUser)
	srv.DELETE(getURL(), srv.DeleteUser)
	return srv, nil
}

func (srv *Service) GetUsers(c echo.Context) error {
	log.Println(fmt.Sprintf("| Info | %s - GetUsers has been called", getURL()))
	query := new(User)
	if err := c.Bind(query); err != nil {
		log.Println(fmt.Sprintf("| Error | %s - can't parse query", getURL()))
		return c.NoContent(http.StatusBadRequest)
	}
	users := []User{}
	if srv.Find(&users, query).RecordNotFound() {
		log.Println(fmt.Sprintf("| Error | %s - user(s) not found", getURL()))
		c.NoContent(http.StatusNotFound)
	}
	return c.JSON(http.StatusOK, users)
}

func (srv *Service) CreateUser(c echo.Context) error {
	log.Println(fmt.Sprintf("| Info | %s - CreateUser has been called", getURL()))
	query := new(User)
	if err := c.Bind(query); err != nil {
		log.Println(fmt.Sprintf("| Error | %s - can't parse query", getURL()))
		return c.NoContent(http.StatusBadRequest)
	}
	if query.Validate() {
		user := new(User)
		query = &User{Name: query.Name, Age: query.Age, Email: query.Email, Address: query.Address}
		srv.Model(User{}).Create(query).Last(user, query)
		log.Println(fmt.Sprintf("| Info | %s - OK", getURL()))
		return c.JSON(http.StatusOK, user)
	}
	log.Println(fmt.Sprintf("| Error | %s - query is not valid", getURL()))
	return c.NoContent(http.StatusBadRequest)
}

func (srv *Service) UpdateUser(c echo.Context) error {
	log.Println(fmt.Sprintf("| Info | %s - UpdateUser has been called", getURL()))
	query := new(User)
	if err := c.Bind(query); err != nil {
		log.Println(fmt.Sprintf("| Error | %s - can't parse query", getURL()))
		return c.NoContent(http.StatusBadRequest)
	}
	if !query.Validate() {
		log.Println(fmt.Sprintf("| Error | %s - query is not valid", getURL()))
		c.NoContent(http.StatusBadRequest)
	}
	user := new(User)
	if srv.Model(User{}).First(user, query.ID).RecordNotFound() {
		log.Println(fmt.Sprintf("| Error | %s - user not found", getURL()))
		c.NoContent(http.StatusNotFound)
	}
	query = &User{Name: query.Name, Age: query.Age, Email: query.Email, Address: query.Address}
	srv.Model(user).Update(query)
	srv.Model(User{}).First(user, query.ID)
	log.Println(fmt.Sprintf("| Info | %s - OK", getURL()))
	return c.JSON(http.StatusOK, user)
}

func (srv *Service) DeleteUser(c echo.Context) error {
	log.Println(fmt.Sprintf("| Info | %s - DeleteUser has been called", getURL()))
	query := new(User)
	if err := c.Bind(query); err != nil {
		log.Println(fmt.Sprintf("| Error | %s - can't parse query", getURL()))
		return c.NoContent(http.StatusBadRequest)
	}
	user := new(User)
	if srv.Model(User{}).First(user, query.ID).RecordNotFound() {
		log.Println(fmt.Sprintf("| Error | %s - user not found", getURL()))
		c.NoContent(http.StatusNotFound)
	}
	srv.Delete(user)
	log.Println(fmt.Sprintf("| Info | %s - OK", getURL()))
	return c.NoContent(http.StatusNoContent)
}

func (user *User) Validate() bool {
	regex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !regex.MatchString(user.Email) || user.Name == "" || user.Age < 1 {
		return false
	}
	return true
}
