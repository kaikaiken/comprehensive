package user

import (
	"bangumiBackend/db"
	"bangumiBackend/model"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

const BaseURL = "/api/user"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e *echo.Echo) (err error) {
	e.GET(BaseURL+"/all", getAllUsers)
	e.GET(BaseURL+"/:id", getUser)
	e.POST(BaseURL+"/login", userLogin)
	e.PUT(BaseURL+"/username/:username", updateUser)
	e.GET(BaseURL+"/username/:username", getUserByUserName)
	e.POST(BaseURL, newUser)
	e.GET(BaseURL+"/username/:username/addBangumi/:bangumiId", addBangumi)
	e.GET(BaseURL+"/username/:username/removeBangumi/:bangumiId", removeBangumi)
	return nil
}

func getAllUsers(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var results []model.User
	err = collection.Find(nil).All(&results)
	if err != nil {
		return err
	}
	fmt.Println(results)
	return c.JSON(http.StatusOK, results)
}

func getUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	id := c.Param("id")
	var result model.User
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func updateUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var updatedUser model.User
	err = c.Bind(&updatedUser)
	if err != nil {
		return err
	}
	username := c.Param("username")
	err = collection.Update(bson.M{"username": username}, updatedUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func getUserByUserName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	username := c.Param("username")
	var result model.User
	err = collection.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newUser(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var newUser model.User
	var checkUser model.User
	err = c.Bind(&newUser)
	if err != nil {
		return err
	}
	collection.Find(bson.M{"username": newUser.UserName}).One(checkUser)
	if len(checkUser.UserName) > 0 {
		return err
	}

	newUser.Role = 2
	newUser.BangumiList = []string{}
	// newUser.BangumiList = ""
	newUser.Id = bson.NewObjectId()
	err = collection.Insert(&newUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newUser)
}

func addBangumi(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var updatedUser model.User
	username := c.Param("username")
	err = collection.Find(bson.M{"username": username}).One(&updatedUser)
	if err != nil {
		return err
	}
	bangumiId := c.Param("bangumiId")
	bangumiList := updatedUser.BangumiList
	bangumiList = append(bangumiList, bangumiId)
	// if len(bangumiList) == 0 {
	// 	bangumiList = bangumiId
	// } else {
	// 	bangumiList = bangumiList + "," + bangumiId
	// }
	updatedUser.BangumiList = bangumiList
	err = collection.Update(bson.M{"username": username}, updatedUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func removeBangumi(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var updatedUser model.User
	username := c.Param("username")
	err = collection.Find(bson.M{"username": username}).One(&updatedUser)
	if err != nil {
		return err
	}
	bangumiId := c.Param("bangumiId")
	bangumiList := updatedUser.BangumiList

	deleteInSlice(bangumiList, bangumiId)

	// bangumi_array := strings.Split(bangumiList, `,`)
	// for i, v := range bangumi_array {
	// 	if v == bangumiId {
	// 		bangumi_array = append(bangumi_array[:i], bangumi_array[i+1:]...)
	// 	}
	// }
	// fmt.Println(bangumi_array)
	// bangumiList = strings.Join(bangumi_array, ",")
	updatedUser.BangumiList = bangumiList
	err = collection.Update(bson.M{"username": username}, updatedUser)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedUser)
}

func deleteInSlice(sli []string, str string) []string {
	for k, v := range sli {
		if len(sli) == 1 {
			return []string{}
		}
		if v == str {
			sli = append(sli[:k], sli[k+1:]...)
		}
	}
	return sli
}

func userLogin(c echo.Context) (err error) {
	username := c.FormValue("username")
	password := c.FormValue("password")
	collection, closeConn := db.GlobalDatabase.User()
	defer closeConn()
	var result model.User
	err = collection.Find(bson.M{"username": username}).One(&result)
	fmt.Println(result)
	if err != nil {
		return c.String(http.StatusUnauthorized, "用户名不存在!")
	}
	if password != result.Password {
		return c.String(http.StatusUnauthorized, "密码不正确!")
	} else {
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["name"] = username
		claims["admin"] = true
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
		// Generate encoded token and send it as response.
		t, err := token.SignedString([]byte("secret"))
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}
