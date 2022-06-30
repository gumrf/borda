package middleware

import (
	"borda/internal/pkg/response"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
)

func AuthRequired(c *fiber.Ctx) error {
	token := c.Locals("token").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)

	// Get user id, scope from claims
	id := claims["iss"].(string)
	scope := claims["scope"].([]interface{})

	intId, _ := strconv.Atoi(id)

	// Store user id, scope in context for the following routes
	c.Locals("USER_ID", intId)
	c.Locals("SCOPE", scope[0])

	fmt.Println("User ID: "+id+", Scope: ", scope[0])

	return c.Next()
}

// func TeamRequired(c *fiber.Ctx, authService auth.AuthService) error {
// 	id := c.Locals("USER_ID").(int)

// 	teamId, ok := authService.VerifyUserTeam(id)
// 	if !ok {
// 		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse{
// 			Status: strconv.Itoa(fiber.StatusForbidden),
// 			Code:   response.MissingTeamIdCode,
// 			Title:  "Not a member of any group",
// 			Detail: "Join or create a team before requesting this route.",
// 		})
// 	}

// 	c.Locals("TEAM_ID", teamId)

// 	return c.Next()
// }

func AdminPermissionRequired(c *fiber.Ctx) error {
	scope := c.Locals("SCOPE")

	if scope != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(response.ErrorResponse{
			Status: strconv.Itoa(fiber.StatusForbidden),
			Code:   response.ForbiddenCode,
			Title:  "Admin permission required.",
			Detail: "Provide API key with right permission to access this route.",
		})
	}

	return c.Next()
}
