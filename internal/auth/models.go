package auth

type UserLoginDTO struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

type UserRegisterDTO struct {
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}
