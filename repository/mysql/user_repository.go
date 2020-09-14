package mysql

import (
	"errors"
	"sync"

	"github.com/micro-community/auth/db"
	"github.com/micro-community/auth/models"
	"golang.org/x/crypto/bcrypt"
)

//userRepository data
type UserRepository struct {
	mu    *sync.Mutex
	users []*models.User
}

func (UserRepository) TableName() string {
	return "user"
}

// 获取用户数据
func (u UserRepository) Get() (UserRepository UserRepository, err error) {
	table := db.DB().Table(u.TableName()).Select([]string{"user.*", "role.role_name"})
	table = table.Joins("left join role on user.role_id=role.role_id")
	if u.UserId != 0 {
		table = table.Where("user_id = ?", u.UserId)
	}

	if u.Username != "" {
		table = table.Where("username = ?", u.Username)
	}

	if u.Password != "" {
		table = table.Where("password = ?", u.Password)
	}

	if u.RoleId != 0 {
		table = table.Where("role_id = ?", u.RoleId)
	}

	if u.DeptId != 0 {
		table = table.Where("dept_id = ?", u.DeptId)
	}

	if u.PostId != 0 {
		table = table.Where("post_id = ?", u.PostId)
	}

	if err = table.First(&UserRepository).Error; err != nil {
		return
	}

	UserRepository.Password = ""
	return
}

//加密
func (u *UserRepository) Encrypt() (err error) {
	if u.Password == "" {
		return
	}

	var hash []byte
	if hash, err = bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost); err != nil {
		return
	} else {
		u.Password = string(hash)
		return
	}
}

//添加
func (u UserRepository) Insert() (id int64, err error) {
	if err = u.Encrypt(); err != nil {
		return
	}

	// check 用户名
	var count int64
	db.DB().Table(u.TableName()).Where("username = ?", u.Username).Count(&count)
	if count > 0 {
		err = errors.New("账户已存在！")
		return
	}

	//添加数据
	if err = db.DB().Table(u.TableName()).Create(&u).Error; err != nil {
		return
	}
	id = u.UserId
	return
}

//修改
func (u *UserRepository) Update(id int64) (update UserRepository, err error) {
	if u.Password != "" {
		if err = u.Encrypt(); err != nil {
			return
		}
	}
	if err = db.DB().Table(u.TableName()).First(&update, id).Error; err != nil {
		return
	}
	if u.RoleId == 0 {
		u.RoleId = update.RoleId
	}

	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = db.DB().Table(u.TableName()).Model(&update).Updates(&u).Error; err != nil {
		return
	}
	return
}
func (u *UserRepository) BatchDelete(id []int) (Result bool, err error) {
	if err = db.DB().Table(u.TableName()).Where("user_id in (?)", id).Delete(&UserRepository{}).Error; err != nil {
		return
	}
	Result = true
	return
}

func (u *UserRepository) ToView() *user.UserInfo {
	var v user.UserInfo
	v.Name = u.NickName
	//.....

	return &v
}
