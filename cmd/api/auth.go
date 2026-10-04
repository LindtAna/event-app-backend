package main

import (
	"event-app/internal/database"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type registerRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required,min=1"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type authResponse struct {
	AccessToken string        `json:"accessToken"`
	User        database.User `json:"user"`
}

type updateProfileRequest struct {
	Name      string `json:"name" binding:"required,min=2"`
	Bio       string `json:"bio"`
	AvatarUrl string `json:"avatarUrl"`
}

// Generierung des Access-Tokens (15 Minuten)
func (app *application) generateAccessToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(15 * time.Minute).Unix(),
	})
	return token.SignedString([]byte(app.jwtSecret))
}

// Generierung des Refresh-Tokens (7 Tage)
func (app *application) generateRefreshToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(7 * 24 * time.Hour).Unix(),
	})
	return token.SignedString([]byte(app.jwtSecret))
}

// Speichern des Refresh-Tokens in einem HttpOnly-Cookie
func (app *application) setRefreshTokenCookie(c *gin.Context, refreshToken string) {
	// MaxAge: 7 Tage = 7 * 24 * 3600 Sec
	// httpOnly: true (Schutz vor XSS), secure: false (für localhost; in der Produktion -> true)
	c.SetCookie("refreshToken", refreshToken, 7*24*3600, "/api/v1/auth", "localhost", false, true)
}

// RegisterUser registers a new user
// @Summary		Registers a new user
// @Description	Registers a new user
// @Tags auth
// @Accept json
// @Produce	json
// @Param user body	registerRequest	true "User"
// @Success 201	{object}	database.User
// @Router /api/v1/auth/register [post]
func (app *application) registerUser(c *gin.Context) {
	var register registerRequest

	if err := c.ShouldBindJSON(&register); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}

	// register.Password = string(hashedPassword)
	user := database.User{
		Email:    register.Email,
		Password: string(hashedPassword),
		Name:     register.Name,
	}

	err = app.models.Users.Insert(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
		return
	}

	c.JSON(http.StatusCreated, user)

}

func (app *application) updateProfile(c *gin.Context) {
	// Benutzer-ID aus der Middleware abrufen
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Daten in der Datenbank aktualisieren
	err := app.models.Users.UpdateProfile(userID.(int), req.Name, req.Bio, req.AvatarUrl)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile"})
		return
	}

	// Aktualisierte Benutzerdaten abrufen
	updatedUser, err := app.models.Users.Get(userID.(int))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// Login logs in a user
//
//	@Summary Logs in a user
//	@Description Logs in a user
//	@Tags auth
//	@Accept json
//	@Produce json
//	@Param user body loginRequest true "User"
//	@Success 200	{object} loginResponse
//	@Router /api/v1/auth/login [post]
func (app *application) login(c *gin.Context) {
	var auth loginRequest

	if err := c.ShouldBindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	existingUser, err := app.models.Users.GetByEmail(auth.Email)
	if err != nil {
		log.Println("ERROR GetByEmail:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Something went wrong"})
		return
	}
	if existingUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(auth.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, err := app.generateAccessToken(existingUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	refreshToken, err := app.generateRefreshToken(existingUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	app.setRefreshTokenCookie(c, refreshToken)

	c.JSON(http.StatusOK, authResponse{
		AccessToken: accessToken,
		User:        *existingUser,
	})
}

// Endpoint zur Aktualisierung des Access-Tokens mittels HttpOnly-Cookie
func (app *application) refreshToken(c *gin.Context) {
	cookie, err := c.Cookie("refreshToken")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token missing"})
		return
	}

	token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(app.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := int(claims["userId"].(float64))
	user, err := app.models.Users.Get(userID)
	if err != nil || user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}

	newAccessToken, err := app.generateAccessToken(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
		return
	}

	c.JSON(http.StatusOK, authResponse{
		AccessToken: newAccessToken,
		User:        *user,
	})
}

// Logout-Endpoint (löscht das Cookie)
func (app *application) logout(c *gin.Context) {
	c.SetCookie("refreshToken", "", -1, "/api/v1/auth", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Successfully logged out"})
}
