package main

import (
	"encoding/csv"
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

// localRepResponse ... looks at
type localRepResponse struct {
	Index    int    `json:"index" binding:"required"`
	Office   string `json:"office"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	Division string `json:"division"`
}

// Representative ... object of representative
type Representative struct {
	GUID     string `json:"guid" binding:"required"`
	Office   string `json:"office"`
	Name     string `json:"name" binding:"required"`
	Location string `json:"location"`
	Division string `json:"division"`
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
	localRepResponse{0, "President of the United States", "Donald J. Trump", "United States", ""},
	localRepResponse{1, "Vice President of the United States", "Mike Pence", "United States", ""},
	localRepResponse{2, "U.S. Senator", "Cory Gardner", "Colorado", ""},
	localRepResponse{3, "U.S. Senator", "Michael F. Bennet", "Colorado", ""},
	localRepResponse{4, "U.S. Representative", "Diana DeGette", "Colorado's 1st congressional district", ""},
	localRepResponse{5, "Governor of Colorado", "Jared Polis", "Colorado", ""},
	localRepResponse{6, "Lieutenant Governor of Colorado", "Dianne Primavera", "Colorado", ""},
	localRepResponse{7, "CO Secretary of State", "Jena Griswold", "Colorado", ""},
	localRepResponse{8, "CO Attorney General", "Phil Weiser", "Colorado", ""},
	localRepResponse{9, "CO State Treasurer", "Dave Young", "Colorado", ""},
	localRepResponse{10, "CO Supreme Court Justice", "Carlos A. Samour, Jr.", "Colorado", ""},
	localRepResponse{11, "CO Supreme Court Justice", "Monica M. Márquez", "Colorado", ""},
	localRepResponse{12, "CO Supreme Court Justice", "Richard L. Gabriel", "Colorado", ""},
	localRepResponse{13, "CO Supreme Court Justice", "Brian D. Boatright", "Colorado", ""},
	localRepResponse{14, "CO Supreme Court Justice", "Nathan B. Coats", "Colorado", ""},
	localRepResponse{15, "CO Supreme Court Justice", "Melissa Hart", "Colorado", ""},
	localRepResponse{16, "CO Supreme Court Justice", "William W. Hood, III", "Colorado", ""},
	localRepResponse{17, "Denver City Clerk and Recorder", "Paul López", "Denver County", ""},
	localRepResponse{18, "Denver Mayor", "Michael Hancock", "Denver County", ""},
	localRepResponse{19, "Denver City Auditor", "Timothy M. O'Brien", "Denver County", ""},
	localRepResponse{20, "Denver City Council Member", "Deborah Ortega", "Denver County", ""},
	localRepResponse{21, "Denver City Council Member", "Robin Kniech", "Denver County", ""}}

func googleRepLookup(c *gin.Context) {
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
		tempResponse := localRepResponse{Index: i, Office: officeMap[i], Name: googleCivic.Officials[i].Name, Location: divisionMap[officeDivisionMap[officeMap[i]]], Division: officeDivisionMap[officeMap[i]]}
		response = append(response, tempResponse)
		// fmt.Println(strconv.Itoa(i) + " - " + officeMap[i] + " - " + googleCivic.Officials[i].Name + " - " + divisionMap[officeDivisionMap[officeMap[i]]])
	}

	msg := map[string]interface{}{"Status": "Ok", "address": address, "representatives": response}
	// fmt.Println(response)
	c.JSON(http.StatusOK, msg)
}

func localRepLookup(c *gin.Context) {
	address, _ := c.GetQuery("address")
	configuration := Configuration{}
	gonfig.GetConf("./data/config.json", &configuration)

	divisionList := zipDivisionMap[address]
	RepResponse := []Representative{}
	for _, division := range divisionList {
		fmt.Println(division)
		for _, rep := range divisionRepMap[division] {
			RepResponse = append(RepResponse, rep)
		}

	}

	msg := map[string]interface{}{"Status": "Ok", "address": address, "representatives": RepResponse}
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

// EditLocalRep adds or removes a local rep in a user's feed
func EditLocalRep(c *gin.Context) {
	userGUID, _ := c.GetQuery("user_guid")
	repGUID, _ := c.GetQuery("rep_guid")
	editTask, _ := c.GetQuery("editTask")
	c.Header("Content-Type", "application/json")
	targetRepIndex := -1
	repGUIDInt, _ := strconv.Atoi(repGUID)
	if editTask == "add" {
		userDB[userGUID] = append(userDB[userGUID], repGUIDInt)
	} else if editTask == "remove" {
		tempUserRepList := userDB[userGUID]
		for i, value := range tempUserRepList {
			if value == repGUIDInt {
				targetRepIndex = i
			}
		}
		if targetRepIndex != -1 {
			userDB[userGUID] = append(tempUserRepList[:targetRepIndex], tempUserRepList[targetRepIndex+1:]...)
		}
	} else {
		log.Info("edit Rep: provided invalid option")
	}
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

// downloadFromS3 is not yet in use
// func downloadFromS3(bucket string, key string) {
// 	awsS3Session, S3err := session.NewSessionWithOptions(session.Options{
// 		Config: aws.Config{Region: aws.String("us-west-2")},
// 	})
// 	if S3err != nil {
// 		fmt.Errorf("Not able to load S3 %q, %v", bucket, key)
// 	}

// 	downloader := s3manager.NewDownloader(awsS3Session)
// 	filePath := "data/" + bucket + key
// }

func loadRepDB(filePath string) map[string]Representative {
	f, err := os.Open(filePath)
	if err != nil {
		logrus.WithField("path", filePath).WithError(err).Error("Error while loading file")
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	// counter := 0
	// OfficialsMap := []([]string){}
	RepMap := map[string]Representative{}
	for _, line := range lines {
		tempGUID := line[4]
		tempOffice := line[0]
		tempName := line[1]
		tempLocation := line[2]
		tempDivision := line[3]
		RepMap[tempGUID] = Representative{tempGUID, tempOffice, tempName, tempLocation, tempDivision}
		// fmt.Println(OfficialsMap[line[6]])
	}

	return RepMap
}

func loadDivisionRepDB(filePath string) map[string][]Representative {
	f, err := os.Open(filePath)
	if err != nil {
		logrus.WithField("path", filePath).WithError(err).Error("Error while loading file")
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	// counter := 0
	// OfficialsMap := []([]string){}
	currentDivision := ""
	OfficialsMap := map[string][]Representative{}
	// mapTest := map[string]int{}
	tempRepList := []Representative{}
	for _, line := range lines {
		tempGUID := line[4]
		tempOffice := line[0]
		tempName := line[1]
		tempLocation := line[2]
		tempDivision := line[3]
		tempRep := Representative{tempGUID, tempOffice, tempName, tempLocation, tempDivision}

		if currentDivision == line[3] {
			tempRepList = append(tempRepList, tempRep)
		} else if currentDivision != line[3] {
			OfficialsMap[currentDivision] = tempRepList
			currentDivision = line[3]
			tempRepList = []Representative{}
			tempRepList = append(tempRepList, tempRep)
			fmt.Println(currentDivision)
			fmt.Println(tempRep)
		}
	}
	return OfficialsMap
}

func loadZipDivisionDB(filePath string) map[string][]string {
	f, err := os.Open(filePath)
	if err != nil {
		logrus.WithField("path", filePath).WithError(err).Error("Error while loading file")
	}
	defer f.Close()
	lines, err := csv.NewReader(f).ReadAll()
	currentZip := ""
	ZipMap := map[string][]string{}
	tempDivList := []string{}
	for _, line := range lines {
		tempDivision := line[1]
		if currentZip == line[0] {
			tempDivList = append(tempDivList, tempDivision)
		} else if currentZip != line[0] {
			ZipMap[currentZip] = tempDivList
			currentZip = line[0]
			tempDivList = []string{}
			tempDivList = append(tempDivList, tempDivision)
		}
	}
	return ZipMap
}

// func MapLocalRepsHandler(c *gin.Context) {
// 	zipcode, _ := c.GetQuery("zipcode")
// 	userGUID := "1234"
// 	c.Header("Content-Type", "application/json")
// 	var tempUserRepList []localRepResponse
// 	for _, j := range userDB[userGUID] {
// 		fmt.Println("j", j)
// 		fmt.Println("Rep: ", localRepsDB[j].Name)
// 		tempUserRepList = append(tempUserRepList, localRepsDB[j])
// 	}
// 	msg := map[string]interface{}{"Status": "Ok", "user_guid": userGUID, "users_rep_list": tempUserRepList}
// 	c.JSON(http.StatusOK, msg)
// }

var (
	userDB         = make(map[string][]int)
	jwtMiddleWare  *jwtmiddleware.JWTMiddleware
	log            = logrus.New()
	googleCivic    civicResponse
	repMap         map[string]Representative
	divisionRepMap map[string][]Representative
	zipDivisionMap map[string][]string
)

func init() {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
	// Load in-memory maps for reference
	repMap = loadRepDB("./data/officials.csv")
	divisionRepMap = loadDivisionRepDB("./data/officials.csv")
	zipDivisionMap = loadZipDivisionDB("./data/zip_divisions_db.csv")
}

func main() {

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

	// fmt.Println(divisionRepMap["ocd-division/country:us/state:co"])
	// fmt.Println(repMap["4a48e646caf9b30d8245038215efe4"])
	// fmt.Println("")
	// fmt.Println(zipDivisionMap["80204"])

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

	api.GET("/localreps", LocalRepsHandler)
	// api.POST("/localreps/add", authMiddleware(), AddLocalRep)
	api.POST("/localreps/edit", EditLocalRep)
	api.GET("/localreps/lookup", localRepLookup)
	api.GET("/localreps/google/lookup", googleRepLookup)
	r.Run(":3000")
}
