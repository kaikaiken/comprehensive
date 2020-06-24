package article

import (
	"bangumiBackend/db"
	"bangumiBackend/model"
	"net/http"

	"github.com/globalsign/mgo/bson"
	"github.com/labstack/echo"
)

const BaseURL = "/api/article"

type Controller struct{}

func NewController() *Controller {
	return &Controller{}
}

func (c *Controller) Initialize(e *echo.Echo) (err error) {
	e.GET(BaseURL+"/all", getAllArticles)
	e.GET(BaseURL+"/:id", getArticle)
	e.GET(BaseURL+"/title/:title", getArticleByTitle)
	e.POST(BaseURL, newArticle)
	e.DELETE(BaseURL+"/:id", deleteArticle)
	e.PUT(BaseURL+"/:id", updateArticle)

	e.GET(BaseURL+"/tag/all", getAllArticleTags)
	e.GET(BaseURL+"/tag/:tagName", getArticleTag)
	e.POST(BaseURL+"/tag", newArticleTag)
	e.DELETE(BaseURL+"/tag/tagName/:tagName", deleteArticleTagByName)

	e.GET(BaseURL+"/kind/all", getAllArticleKinds)
	e.GET(BaseURL+"/kind/:kindName", getArticleKind)
	e.POST(BaseURL+"/kind", newArticleKind)
	e.DELETE(BaseURL+"/kind/kindName/:kindName", deleteArticleKindByName)
	e.DELETE(BaseURL+"/kind/:id", deleteArticleKind)

	return nil
}

func getAllArticles(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var results []model.Article
	err = collection.Find(nil).Sort("-pub_time").All(&results)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	id := c.Param("id")
	var result model.Article
	err = collection.FindId(bson.ObjectIdHex(id)).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func getArticleByTitle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	title := c.Param("title")
	var result model.Article
	err = collection.Find(bson.M{"title": title}).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var newArticle model.Article
	err = c.Bind(&newArticle)
	if err != nil {
		return err
	}
	newArticle.Id = bson.NewObjectId()
	err = collection.Insert(&newArticle)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, newArticle)
}

func deleteArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	kindId := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(kindId))
	return c.NoContent(http.StatusNoContent)
}

func updateArticle(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.Article()
	defer closeConn()
	var updatedArticle model.Article
	err = c.Bind(&updatedArticle)
	if err != nil {
		return err
	}
	id := c.Param("id")
	err = collection.UpdateId(bson.ObjectIdHex(id), updatedArticle)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, updatedArticle)
}

func getAllArticleTags(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	var results []model.ArticleTag
	err = collection.Find(nil).All(&results)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticleTag(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	tagName := c.Param("tagName")
	var result model.ArticleTag
	err = collection.Find(bson.M{"articleTag": tagName}).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticleTag(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	var newArticleTag model.ArticleTag
	err = c.Bind(&newArticleTag)
	if err != nil {
		return err
	}
	newArticleTag.Id = bson.NewObjectId()
	err = collection.Insert(&newArticleTag)
	return c.JSON(http.StatusOK, newArticleTag)
}

func deleteArticleTagByName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleTag()
	defer closeConn()
	tagName := c.Param("tagName")
	err = collection.Remove(bson.M{"articleTag": tagName})
	return c.NoContent(http.StatusNoContent)
}

func getAllArticleKinds(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	var results []model.ArticleKind
	err = collection.Find(nil).All(&results)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, results)
}

func getArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindName := c.Param("kindName")
	var result model.ArticleKind
	err = collection.Find(bson.M{"articleKind": kindName}).One(&result)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, result)
}

func newArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	var newArticleKind model.ArticleKind
	err = c.Bind(&newArticleKind)
	if err != nil {
		return err
	}
	newArticleKind.Id = bson.NewObjectId()
	err = collection.Insert(&newArticleKind)
	return c.JSON(http.StatusOK, newArticleKind)
}

func deleteArticleKindByName(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindName := c.Param("kindName")
	err = collection.Remove(bson.M{"articleKind": kindName})
	return c.NoContent(http.StatusNoContent)
}

func deleteArticleKind(c echo.Context) (err error) {
	collection, closeConn := db.GlobalDatabase.ArticleKind()
	defer closeConn()
	kindId := c.Param("id")
	err = collection.RemoveId(bson.ObjectIdHex(kindId))
	return c.NoContent(http.StatusNoContent)
}
