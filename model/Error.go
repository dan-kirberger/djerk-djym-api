package model

type ErrorResponse struct {
	Status  int    `bson:"_id" json:"status"`
	Message string `bson:"firstName" json:"message"`
}
