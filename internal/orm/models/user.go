package models

import (
	"github.com/ezavalishin/partygames/pkg/utils"
	"github.com/jinzhu/gorm"
	vkapi "github.com/leidruid/golang-vk-api"
)

type User struct {
	BaseModel
	VkUserId                int     `gorm:"NOT NULL" json:"vkUserId"`
	FirstName               *string `json:"firstName"`
	LastName                *string `json:"lastName"`
	Avatar                  *string `json:"avatar"`
	NotificationsAreEnabled bool    `json:"notificationsAreEnabled"`
	IsFollower              bool    `json:"isFollower"`
}

func (u *User) AfterCreate(scope *gorm.Scope) (err error) {

	fillUserFromVk(u, scope)

	return
}

func fillUserFromVk(u *User, scope *gorm.Scope) {
	client, err := vkapi.NewVKClientWithToken(utils.MustGet("VK_APP_SERVICE_KEY"), &vkapi.TokenOptions{
		ServiceToken:    true,
		TokenLanguage:   "ru",
		ValidateOnStart: true,
	})

	if err != nil {
		return
	}

	users, err := client.UsersGet([]int{u.VkUserId})

	if err != nil {
		return
	}

	vkUser := users[0]

	u.FirstName = &vkUser.FirstName
	u.LastName = &vkUser.LastName
	u.Avatar = &vkUser.PhotoMedium

	scope.DB().Save(&u)
}
