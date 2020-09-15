package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/tkanos/gonfig"
)

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

type userRepList struct {
	UserGUID string
	RepIndex []int
}

// Configuration ... configuration data
type Configuration struct {
	KeyName  string
	KeyValue string
}

// Jwks stores a slice of JSON Web Keys
type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
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

// create a local example DB of reps
var localRepsDB = []localRepResponse{
	localRepResponse{0, "President of the United States", "Donald J. Trump", "United States"},
	localRepResponse{1, "Vice President of the United States", "Mike Pence", "United States"},
	localRepResponse{2, "U.S. Senator", "Cory Gardner", "Colorado"},
	localRepResponse{3, "U.S. Senator", "Michael F. Bennet", "Colorado"},
	localRepResponse{4, "U.S. Representative", "Diana DeGette", "Colorado's 1st congressional district"},
	localRepResponse{5, "Governor of Colorado", "Jared Polis", "Colorado"},
	localRepResponse{6, "Lieutenant Governor of Colorado", "Dianne Primavera", "Colorado"},
	localRepResponse{7, "CO Secretary of State", "Jena Griswold", "Colorado"},
	localRepResponse{8, "CO Attorney General", "Phil Weiser", "Colorado"},
	localRepResponse{9, "CO State Treasurer", "Dave Young", "Colorado"},
	localRepResponse{10, "CO Supreme Court Justice", "Carlos A. Samour, Jr.", "Colorado"},
	localRepResponse{11, "CO Supreme Court Justice", "Monica M. Márquez", "Colorado"},
	localRepResponse{12, "CO Supreme Court Justice", "Richard L. Gabriel", "Colorado"},
	localRepResponse{13, "CO Supreme Court Justice", "Brian D. Boatright", "Colorado"},
	localRepResponse{14, "CO Supreme Court Justice", "Nathan B. Coats", "Colorado"},
	localRepResponse{15, "CO Supreme Court Justice", "Melissa Hart", "Colorado"},
	localRepResponse{16, "CO Supreme Court Justice", "William W. Hood, III", "Colorado"},
	localRepResponse{17, "Denver City Clerk and Recorder", "Paul López", "Denver County"},
	localRepResponse{18, "Denver Mayor", "Michael Hancock", "Denver County"},
	localRepResponse{19, "Denver City Auditor", "Timothy M. O'Brien", "Denver County"},
	localRepResponse{20, "Denver City Council Member", "Deborah Ortega", "Denver County"},
	localRepResponse{21, "Denver City Council Member", "Robin Kniech", "Denver County"}}

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

// LocalRepsHandler looks up local representatives based on zipcode
func LocalRepsHandler(c *gin.Context) {
	//userGUID, _ := c.GetQuery("user_guid")
	userGUID := "1234"
	c.Header("Content-Type", "application/json")
	var tempUserRepList []localRepResponse
	for _, j := range userDB[userGUID] {
		fmt.Println("j", j)
		fmt.Println("Rep: ", localRepsDB[j].Name)
		tempUserRepList = append(tempUserRepList, localRepsDB[j])
	}
	msg := map[string]interface{}{"Status": "Ok", "user_guid": userGUID, "users_rep_list": tempUserRepList}
	c.JSON(http.StatusOK, msg)
}

// AddLocalRep adds a local rep to user's feed
func AddLocalRep(c *gin.Context) {
	userGUID, _ := c.GetQuery("user_guid")
	repGUID, _ := c.GetQuery("rep_guid")
	c.Header("Content-Type", "application/json")
	repGUIDInt, _ := strconv.Atoi(repGUID)
	userDB[userGUID] = append(userDB[userGUID], repGUIDInt)
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": "AddLocalRep handler not implemented yet",
	// })

	msg := map[string]interface{}{"Status": "Ok", "user_guid": userGUID, "users_rep_list": userDB[userGUID]}
	c.JSON(http.StatusOK, msg)
}

func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the client secret key
		err := jwtMiddleWare.CheckJWT(c.Writer, c.Request)
		if err != nil {
			// Token not found
			fmt.Println(err)
			c.Abort()
			c.Writer.WriteHeader(http.StatusUnauthorized)
			c.Writer.Write([]byte("Unauthorized"))
			return
		}
	}
}

func getPemCert(token *jwt.Token) (string, error) {
	cert := ""
	resp, err := http.Get(os.Getenv("AUTH0_DOMAIN") + ".well-known/jwks.json")
	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	x5c := jwks.Keys[0].X5c
	for k, v := range x5c {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + v + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		return cert, errors.New("unable to find appropriate key.")
	}

	return cert, nil
}

// JokeHandler returns a list of jokes available (in memory)
func JokeHandler(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	c.JSON(http.StatusOK, jokes)
}

// LikeJoke increments the likes of a particular joke Item
func LikeJoke(c *gin.Context) {
	// Check joke ID is valid
	if jokeid, err := strconv.Atoi(c.Param("jokeID")); err == nil {
		// find joke and increment likes
		for i := 0; i < len(jokes); i++ {
			if jokes[i].ID == jokeid {
				jokes[i].Likes = jokes[i].Likes + 1
			}
		}
		c.JSON(http.StatusOK, &jokes)
	} else {
		// the jokes ID is invalid
		c.AbortWithStatus(http.StatusNotFound)
	}
}

// Joke contains information about a single Joke
type Joke struct {
	ID    int    `json:"id" binding:"required"`
	Likes int    `json:"likes"`
	Joke  string `json:"joke" binding:"required"`
}

// We'll create a list of jokes
var jokes = []Joke{
	Joke{1, 0, "Did you hear about the restaurant on the moon? Great food, no atmosphere."},
	Joke{2, 0, "What do you call a fake noodle? An Impasta."},
	Joke{3, 0, "How many apples grow on a tree? All of them."},
	Joke{4, 0, "Want to hear a joke about paper? Nevermind it's tearable."},
	Joke{5, 0, "I just watched a program about beavers. It was the best dam program I've ever seen."},
	Joke{6, 0, "Why did the coffee file a police report? It got mugged."},
	Joke{7, 0, "How does a penguin build it's house? Igloos it together."},
}

var userDB = make(map[string][]int)
var jwtMiddleWare *jwtmiddleware.JWTMiddleware
var log = logrus.New()
var googleCivic civicResponse

func main() {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			aud := os.Getenv("AUTH0_API_AUDIENCE")
			checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
			if !checkAudience {
				return token, errors.New("invalid audience")
			}
			// verify iss claim
			iss := os.Getenv("AUTH0_DOMAIN")
			checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
			if !checkIss {
				return token, errors.New("invalid issuer")
			}

			cert, err := getPemCert(token)
			if err != nil {
				log.Fatalf("could not get cert: %+v", err)
			}

			result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
			return result, nil
		},
		SigningMethod: jwt.SigningMethodRS256,
	})

	// register our actual jwtMiddleware
	jwtMiddleWare = jwtMiddleware

	r := gin.Default()
	r.Use(static.Serve("/", static.LocalFile("./views", true)))

	userDB["1234"] = []int{2, 3, 4}
	userDB["MAGA"] = []int{0, 1}

	// config := cors.DefaultConfig()
	// config.AllowOrigins = []string{"*"}

	// r.Use(cors.New(config))
	// log.Info("Starting Application")
	// r.GET("/ready", getStatus)
	// r.GET("localReps", localRepresentatives)

	// Setup route group for the API
	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})
	}

	api.GET("/jokes", JokeHandler)
	api.POST("/jokes/like/:jokeID", LikeJoke)

	api.GET("/localreps", authMiddleware(), LocalRepsHandler)
	api.POST("/localreps/add", authMiddleware(), AddLocalRep)
	api.GET("/localreps/lookup", authMiddleware(), localRepresentatives)
	r.Run(":3000")
}
