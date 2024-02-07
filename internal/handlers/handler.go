package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sigit14ap/simple-dating-app-backend/internal/helpers"
	"github.com/sigit14ap/simple-dating-app-backend/internal/models"
	"github.com/sigit14ap/simple-dating-app-backend/internal/services"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService  services.UserServiceInterface
	matchService services.MatchServiceInterface
}

func NewHandler(
	userService services.UserServiceInterface,
	matchService services.MatchServiceInterface,
) *Handler {
	return &Handler{
		userService:  userService,
		matchService: matchService,
	}
}

func (h *Handler) Signup(c *gin.Context) {
	var signupRequest models.AuthRequest
	if err := c.ShouldBindJSON(&signupRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.Validate(signupRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, _ := h.userService.GetUserByUsername(signupRequest.Username)
	if user != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	err := h.userService.RegisterUser(signupRequest.Username, signupRequest.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully"})
}

func (h *Handler) Login(c *gin.Context) {
	var loginRequest models.AuthRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.Validate(loginRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.AuthenticateUser(loginRequest.Username, loginRequest.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := generateJWTToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate JWT token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *Handler) GetNextUser(c *gin.Context) {
	userID, err := helpers.GetUserIDFromToken(c.GetHeader("Authorization"))

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = h.userService.CheckEligibleSwipe(userID)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.matchService.GetNextUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get next user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": &user})
}

func (h *Handler) Swipe(c *gin.Context) {
	var swipeRequest models.SwipeRequest
	if err := c.ShouldBindJSON(&swipeRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	if err := helpers.Validate(swipeRequest); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userID, err := helpers.GetUserIDFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	_, err = h.matchService.FindMatch(userID, swipeRequest.TargetUserID, time.Now())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Match data not found"})
		return
	}

	_, err = h.matchService.UpdateMatch(userID, swipeRequest.TargetUserID, swipeRequest.Match)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update swipe data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Swipe Successful"})
}

func (h *Handler) BuyPackageUnlimitedSwipe(c *gin.Context) {
	userID, err := helpers.GetUserIDFromToken(c.GetHeader("Authorization"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	err = h.userService.UpdateUnlimitedSwipe(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to buy unlimited swipe package"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Buy Package Successful"})
}

func generateJWTToken(userID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":    userID,
		"expired_at": time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
