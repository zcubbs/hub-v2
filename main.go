package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"gopkg.in/yaml.v2"
	"hub-v2/models"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"strconv"
)

func main() {
	groups := loadGroupsFromYaml()
	title := getTitle()
	engine := html.New("./views", ".html")

	app := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/main",
	})
	app.Static("/css", "./public/css")
	app.Static("/assets", "./public/assets")

	engine.Reload(isDevMode())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":  title,
			"Groups": groups,
		})
	})

	app.Get("/tag/:group/:caption", func(c *fiber.Ctx) error {
		groupCaption, err := url.QueryUnescape(c.Params("group"))
		if err != nil {
			log.Println(err)
		}
		tagCaption, err := url.QueryUnescape(c.Params("caption"))
		if err != nil {
			log.Println(err)
		}
		for _, group := range *groups {
			if group.Caption == groupCaption {
				for _, tag := range *group.Tags {
					if tag.Caption == tagCaption {
						return c.Render("tag", fiber.Map{
							"Title": title,
							"Tag":   tag,
						})
					}
				}
			}
		}

		return c.Render("index", fiber.Map{
			"Title": title,
			"Tags":  groups,
		})

	})

	log.Fatal(app.Listen(":8000"))
}

func loadGroupsFromYaml() *[]models.Group {
	yamlFile, err := ioutil.ReadFile(getConfigPath())
	if err != nil {
		log.Fatalf("yamlFile.Get err   #%v ", err)
	}

	yamlGroups := &models.Groups{}

	err = yaml.Unmarshal(yamlFile, yamlGroups)
	if err != nil {
		log.Fatal(err)
	}

	return yamlGroups.Groups
}

func getConfigPath() string {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}
	return configPath
}

func isDevMode() bool {
	devMode := os.Getenv("DEV_MODE")
	if devMode == "" {
		devMode = "false"
	}
	boolVal, _ := strconv.ParseBool(devMode)
	return boolVal
}

func getTitle() string {
	appTitle := os.Getenv("APP_TITLE")
	if appTitle == "" {
		appTitle = "z/HuB"
	}
	return appTitle
}
