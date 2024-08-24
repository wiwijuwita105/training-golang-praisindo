package entity

type User struct {
	ID   int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"`
	Name string `json:"name"`
}
