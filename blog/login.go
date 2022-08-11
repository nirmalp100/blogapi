package blog

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

func UserRegistration(r *gin.Context) {

	var parsedBody map[string]string
	r.BindJSON(&parsedBody)

	password := parsedBody["password"]
	hashedpass := HashAndSalt(password)
	fmt.Println(len(hashedpass))
	fmt.Println(hashedpass)
	fmt.Println("")
	username := parsedBody["username"]
	email := parsedBody["email"]

	userid := GenerateRandomNumbers()
	_, err := db.Exec("insert into users(user_id,username,password,email) values($1,$2,$3,$4)", userid, username, hashedpass, email)
	CheckErr(err)
	r.JSON(http.StatusOK, "Registration successful")

}

func UserLogin(r *gin.Context) {
	var parsedBody map[string]string
	r.BindJSON(&parsedBody)
	newpassword := parsedBody["password"]
	email := parsedBody["email"]
	username := parsedBody["username"]

	//after parsing the login credentials the data is stored in database

	rows, err := db.Query("select username,password,email from users")
	CheckErr(err)

	var storedpassword, storedusername, storedemail string
	for rows.Next() {
		err := rows.Scan(&storedusername, &storedpassword, &storedemail)
		CheckErr(err)
	}

	//comparing registered username and email with the corresponding one's entered while login
	if storedusername == username && storedemail == email {

		//Comparing the password when the user registered
		//with the new login password
		//security purposes

		PasswordisCorrect := ComparePasswords(storedpassword, newpassword)
		if !PasswordisCorrect {
			r.JSON(http.StatusUnauthorized, "Wrong Password")
		} else {

			//After user login is authorized
			//accesstoken(JWT) is created from the login credentials
			//and the token is given as response

			accesstoken, err := GenerateJWT(email, username)
			CheckErr(err)
			r.JSON(http.StatusOK, accesstoken)

		}

	} else {
		r.JSON(http.StatusUnauthorized, "Wrong Credentials")
	}

}

//hashing the password
//adding salt before hashing to make the password unpredictable
func HashAndSalt(password string) string {
	pwd := []byte(password)
	hashed, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	return string(hashed)

}

func ComparePasswords(hashedPwd string, newpass string) bool {
	bytehash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(bytehash, []byte(newpass))
	if err != nil {
		log.Println(err)
		return false
	}

	return true
}
