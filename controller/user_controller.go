package controller

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"goshop/config"
	"goshop/entity"
	"goshop/helper"
	"goshop/model"
	"goshop/service"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const SESSION_ID = "id"

type userController struct {
	userService service.UserService
	authService config.AuthService
}

func NewUserController(userService service.UserService, authService config.AuthService) *userController {
	return &userController{userService, authService}
}

type UserFormatter struct {
	UserID int    `json:"id"`
	Email  string `json:"email"`
	Phone  string `json:"phone"`
	Token  string `json:"token"`
}

func FormatUser(user model.User, token string) UserFormatter { //Token akan didapatkan dari JWT
	formatter := UserFormatter{
		UserID: user.ID,
		Email:  user.Email,
		Phone:  user.Phone,
		Token:  token,
	}

	return formatter
}

func (h *userController) Login(c *gin.Context) {
	var input entity.LoginEmailInput
	err := c.ShouldBindJSON(&input)

	if err != nil {
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}

		responsError := helper.APIResponse("Login Failed #LOG001", http.StatusUnprocessableEntity, "fail", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, responsError)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"errors": err.Error()}

		responsError := helper.APIResponse("Login Failed #LOG002", http.StatusUnprocessableEntity, "fail", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, responsError)
		return
	}

	token, err := h.authService.GenerateTokenUser(loggedInUser.ID)
	if err != nil {
		responseError := helper.APIResponse("Login Failed", http.StatusBadGateway, "fail", "Unable to generate token")
		c.JSON(http.StatusBadGateway, responseError)
		return
	}

	response := helper.APIResponse("Login Success", http.StatusOK, "success", FormatUser(loggedInUser, token))

	c.JSON(http.StatusOK, response)
}

func (h *userController) UpdateProfile(c *gin.Context) {
	var input entity.DataUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// fmt.Println(err.Error())
		// return
		errorMessage := gin.H{"errors": helper.FormatValidationError(err)}

		responseError := helper.APIResponse("Create Account Failed", http.StatusUnprocessableEntity, "fail", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, responseError)
		return
	}

	input.ID = c.MustGet("currentUser").(model.User).ID
	updateUser, err := h.userService.UpdateProfile(input)
	if err != nil {
		responseError := helper.APIResponse("Create Account Failed", http.StatusBadRequest, "fail", nil)
		c.JSON(http.StatusBadRequest, responseError)
		return
	}

	response := helper.APIResponse("Account has been registered", http.StatusOK, "success", updateUser)
	c.JSON(http.StatusOK, response)
}

func (h *userController) LoginBE(c *gin.Context) {
	var input entity.LoginEmailInput
	err := c.ShouldBind(&input)
	if err != nil {
		fmt.Println(err)
		return
	}

	loggedInUser, err := h.userService.Login(input)

	fmt.Println(err)

	if err != nil {
		c.HTML(http.StatusOK, "pages/login", gin.H{
			"Error": "Credential Not Found",
			"Title": "Login",
		})
		return
	}

	token, _ := h.authService.GenerateTokenUser(loggedInUser.ID)

	session := sessions.Default(c)
	if session.Get("token") != token {
		session.Set("token", token)
		session.Options(sessions.Options{
			MaxAge: 3600 * 12, //Session Will Exp in 12 Hours
		})
		session.Save()
	}

	c.Redirect(http.StatusFound, "/register")
}

func (h *userController) RegisterStore(c *gin.Context) {

	var input entity.CreateUserInput

	err := c.ShouldBind(&input)

	if err != nil {
		fmt.Println(err)
	}

	file, err := c.FormFile("avatar")

	if err != nil {
		fmt.Println(err)
	}

	filePath := "storage/" + file.Filename
	c.SaveUploadedFile(file, filePath)

	input.Avatar = filePath
	_, err = h.userService.CreateUser(input)

	if err != nil {
		c.Redirect(http.StatusBadGateway, "/register3")
		return
	}

	rand.Seed(time.Now().UTC().UnixNano())       //Ambil random karakter
	otp := string(strconv.Itoa(rand.Int())[1:6]) //Convert to string agar bisa diambil 5 valuenya

	err = h.userService.SaveOtp(input.Email, otp)

	//Send Otp via email
	_, err = config.SendMail(input.Email, otp)

	if err != nil {
		fmt.Println("Unable to send email", err)
		return
	}

	session := sessions.Default(c)

	session.Set("email", input.Email)
	session.Options(sessions.Options{
		MaxAge: 3600 * 12, //Session Will Exp in 12 Hours
	})
	session.Save()

	c.HTML(http.StatusFound, "pages/verify", gin.H{})
}

func (h *userController) RegisterIndex(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/register", gin.H{
		"Title": "Register",
	})
}

func (h *userController) LoginIndex(c *gin.Context) {
	session := sessions.Default(c)
	if session.Get("token") != nil {
		c.Redirect(http.StatusFound, "/dashboard") //Force Redirect if authenticated
	}
	c.HTML(http.StatusOK, "pages/login", gin.H{
		"Title": "Login",
		"Year":  "2021",
	})
}

func (h *userController) Dashboard(c *gin.Context) {
	c.HTML(http.StatusOK, "pages/dashboard", gin.H{})
}

func (h *userController) DeleteSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Set("token", nil)
	session.Options(sessions.Options{
		MaxAge: -1,
	})
	session.Save()
	fmt.Println(session.Get("token"))
	//c.Redirect(301, "/login")
}

func (h *userController) Verify(c *gin.Context) {
	var input entity.OtpCodeInput

	session := sessions.Default(c)
	email := session.Get("email")
	input.Email = fmt.Sprintf("%v", email)
	err := c.ShouldBind(&input)

	if err != nil {
		c.HTML(http.StatusFound, "pages/verify", gin.H{
			"Error": err,
		})
	}

	_, err = h.userService.CheckOtp(input)

	if err != nil {
		c.HTML(http.StatusFound, "pages/verify", gin.H{
			"Error": err,
		})
	}

	//changeStatus
	err = h.userService.ChangeStatusUser(input.Email)

	if err != nil {
		c.HTML(http.StatusFound, "pages/verify", gin.H{
			"Error": err,
		})
	}

	c.Redirect(http.StatusFound, "/login")

}
