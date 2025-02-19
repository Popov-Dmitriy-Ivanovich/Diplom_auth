package models

type AccessRight uint

const (
	AR_READ = 1 << 0
	AR_CREATE_USER = 1 << 1
	AR_CREATE_ACTION = 1 << 2
	AR_VIEW_AND_RUN_ACTION = 1 << 3
	AR_ALL = AR_READ | AR_CREATE_ACTION | AR_CREATE_USER | AR_VIEW_AND_RUN_ACTION
)

type User struct {
	ID uint
	Login string
	Password []byte `json:"-"`
	AccessRights AccessRight
}

func (u User) HasAccessRight(ar AccessRight) bool {
	return (u.AccessRights & ar) != 0
}