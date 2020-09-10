package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

var log = logrus.New()
var googleCivic civicResponse

// NormalizedInput ... office of a representative
type NormalizedInput struct {
	Line1 string `json:"line1"`
	City  string `json:"city"`
	State string `json:"state"`
	Zip   string `json:"zip"`
}

// Office ... office of a representative
type Office struct {
	Name            string   `json:"name"`
	DivisionID      string   `json:"divisionId"`
	Levels          []string `json:"levels"`
	Roles           []string `json:"roles"`
	OfficialIndices []int    `json:"officialIndices"`
}

// Official ... official in the matching position
type Official struct {
	Name     string `json:"name"`
	Address  string `json:"address"`
	Party    string `json:"party"`
	Phone    string `json:"phone"`
	Urls     string `json:"urls"`
	PhotoURL string `json:"photoUrl"`
	Channels string `json:"channels"`
}

type civicResponse struct {
	NormalizedInput NormalizedInput        `json:"normalizedInput"`
	Kind            string                 `json:"kind"`
	Divisions       map[string]interface{} `json:"divisions"`
	Offices         []Office               `json:"offices"`
	Officials       []Official             `json:"officials"`
}

type localRepResponse struct {
	Index    int
	Office   string
	Name     string
	Location string
}

// Configuration ... configuration data
type Configuration struct {
	KeyName  string
	KeyValue string
}

func setupRouter() *gin.Engine {
	gin.DisableConsoleColor()
	r := gin.New()
	return r
}

func getStatus(c *gin.Context) {
	msg := map[string]interface{}{"Status": "Ok", "msg": "ready", "version": "v1.0.1"}
	c.JSON(http.StatusOK, msg)
}

func localRepresentatives(c *gin.Context) {
	address, _ := c.GetQuery("address")
	configuration := Configuration{}
	err := gonfig.GetConf("./data/config.json", &configuration)

	// read user input of address
	// reader := bufio.NewReader(os.Stdin)
	// fmt.Println("Input Your Address: ")
	// address, _ := reader.ReadString('\n')
	// address = strings.Replace(address, " ", "%20", -1)
	// address = strings.Replace(address, "\n", "", -1)
	// address = "80204"

	resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=" + address + "&includeOffices=true&key=" + configuration.KeyValue)

	// resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=37%20ibis%20dr%20akron%20ohio&includeOffices=true&key=" + civicKey)
	// resp, err := http.Get("https://civicinfo.googleapis.com/civicinfo/v2/representatives?address=80204&includeOffices=true&key=" + configuration.KeyValue)
	if err != nil {
		print(err)
	}

	defer resp.Body.Close()
	byteValue, err := ioutil.ReadAll(resp.Body)

	//fmt.Print(string(body))
	// err = ioutil.WriteFile("denver.json", body, 0644)

	// jsonInputFile := "denver.json"

	// Open our jsonFile
	//jsonFile, err := os.Open(jsonInputFile)
	// if we os.Open returns an error then handle it
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Println("Successfully Opened ibis.json")
	// defer the closing of our jsonFile so that we can parse it later on
	// defer jsonFile.Close()

	// byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, &googleCivic)

	officeMap := make(map[int]string) // map for representatives
	officeDivisionMap := make(map[string]string)
	divisionMap := make(map[string]string)

	// map office names to office indices
	for i := range googleCivic.Offices {
		// fmt.Println(googleCivic.Offices[i].Name)
		for j := range googleCivic.Offices[i].OfficialIndices {
			// fmt.Println(strconv.Itoa(googleCivic.Offices[i].OfficialIndices[j]))
			officeMap[googleCivic.Offices[i].OfficialIndices[j]] = googleCivic.Offices[i].Name
		}
		officeDivisionMap[googleCivic.Offices[i].Name] = googleCivic.Offices[i].DivisionID
	}
	// map division tags to division names (to join to office map)
	for key, value := range googleCivic.Divisions {
		// fmt.Println("key: ", key)
		// fmt.Println("RAW: ", value)
		for key1, value1 := range value.(map[string]interface{}) {
			// fmt.Println("key1: ", key1)
			// fmt.Println("value1: ", value1)
			if key1 == "name" {
				divisionName := fmt.Sprintf("%v", value1)
				divisionMap[key] = divisionName
			}
		}
	}

	var response []localRepResponse
	// fmt.Println("")
	// fmt.Println("")
	for i := 0; i < len(googleCivic.Officials); i++ {
		// fmt.Println(i)
		tempResponse := localRepResponse{Index: i, Office: officeMap[i], Name: googleCivic.Officials[i].Name, Location: divisionMap[officeDivisionMap[officeMap[i]]]}
		response = append(response, tempResponse)
		// fmt.Println(strconv.Itoa(i) + " - " + officeMap[i] + " - " + googleCivic.Officials[i].Name + " - " + divisionMap[officeDivisionMap[officeMap[i]]])
	}

	msg := map[string]interface{}{"Status": "Ok", "address": address, "representatives": response}
	// fmt.Println(response)
	c.JSON(http.StatusOK, msg)
}

// func getCount(c *gin.Context) {
// 	playerstatsCache.ItemCount()
// 	msg := map[string]interface{}{"Status": "Ok", "guids": playerstatsCache.ItemCount()}
// 	c.JSON(http.StatusOK, msg)
// }

func main() {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	r := setupRouter()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	// config.AllowOrigins == []string{"http://google.com", "http://facebook.com"}

	r.Use(cors.New(config))
	log.Info("Starting Application")
	r.GET("/ready", getStatus)
	r.GET("localReps", localRepresentatives)

	r.Run()
}
