package model

type User struct {
	ID        string `bson:"_id" json:"id"`
	FirstName string `bson:"firstName" json:"firstName"`
	LastName  string `bson:"lastName" json:"lastName"`
	Weight    int32  `bson:"weight" json:"weight"`
}

type UserList struct {
	Users []User `json:"users"`
}
