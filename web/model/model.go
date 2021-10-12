package model

func GetUsername(mobile string) string {
	var user User
	Gdb.Where("mobile=?",mobile).Select("name").Find(&user)
	return user.Name
}

func UpdataAvatar(username,avatar string ) error {
	return  Gdb.Model(new(User)).Where("name=?",username).
		Update("avatar_url",avatar).Error
}