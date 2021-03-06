package auth

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/ezavalishin/partygames/internal/orm"
	"github.com/ezavalishin/partygames/internal/orm/models"
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/ezavalishin/partygames/pkg/vk"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"net/http"
	"net/url"
	"strconv"
)

type wssContext struct {
	user models.User
}

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func validateAndGetUserID(vkParams string) (int, error) {

	decoded, err := base64.StdEncoding.DecodeString(vkParams)

	if err != nil {
		return 0, err
	}

	if ok, err := vk.ParamsVerify(string(decoded), utils.MustGet("VK_APP_SECRET")); err != nil || ok != true {
		return 0, errors.New("bad sign")
	}

	parsedUrl, err2 := url.ParseQuery(string(decoded))

	if err2 != nil {
		return 0, err2
	}

	vkUserId, _ := strconv.Atoi(parsedUrl["vk_user_id"][0])

	return vkUserId, nil
}

func getUserByID(orm *orm.ORM, vkUserId int) models.User {

	user := models.User{}

	orm.DB.FirstOrCreate(&user, models.User{VkUserId: vkUserId})

	return user
}

func WsValidateAndSetUser(orm *orm.ORM, s socketio.Conn, vkParams string) error {
	vkUserId, err := validateAndGetUserID(vkParams)

	if err != nil {
		return err
	}

	user := getUserByID(orm, vkUserId)

	s.SetContext(&wssContext{user: user})

	return nil
}

// Middleware decodes the share session cookie and packs the session into context
func Middleware(orm *orm.ORM) gin.HandlerFunc {
	return func(c *gin.Context) {

		var vkParams string

		vkParams = c.Request.Header.Get("Vk-Params")

		if vkParams == "" {
			vkParams = c.Query("vk-params")
		}

		if vkParams == "" {
			http.Error(c.Writer, "Required Vk-Params", http.StatusForbidden)
			c.Abort()
			return
		}

		vkUserId, err := validateAndGetUserID(vkParams)

		if err != nil {
			http.Error(c.Writer, "", http.StatusForbidden)
			c.Abort()
			return
		}

		// get the user from the database
		user := getUserByID(orm, vkUserId)

		// put it in context
		ctx := context.WithValue(c.Request.Context(), userCtxKey, &user)

		// and call the next with our new context
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *models.User {

	raw := ctx.Value(userCtxKey)
	if raw == nil {
		fmt.Println("could not retrieve gin.Context")
		return nil
	}

	gc, ok := raw.(*models.User)
	if !ok {
		fmt.Println("gin.Context has wrong type")
		return nil
	}
	return gc
}

func ForWssContext(s socketio.Conn) *models.User {

	wssC, ok := s.Context().(*wssContext)

	if !ok {
		return nil
	}

	if wssC == nil {
		fmt.Println("could not retrieve gin.Context")
		return nil
	}

	user := wssC.user

	return &user
}
