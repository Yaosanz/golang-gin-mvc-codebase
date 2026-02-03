package helpers

import "github.com/gin-gonic/gin"

type Role struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type AuthUser struct {
	ID       string  `json:"user_id"`
	Email    string  `json:"email"`
	Username string  `json:"username"`
	Name     string  `json:"name"`
	Exp      float64 `json:"exp"`
	Roles    []Role  `json:"roles"`
}

// Auth extracts user claims from the Gin context
func Auth(ctx *gin.Context) *AuthUser {
	// Retrieve "user" from the context
	user, ok := ctx.Get("user")
	if !ok {
		return nil
	}

	// Ensure the user is a map[string]interface{}
	mapClaim, ok := user.(map[string]interface{})
	if !ok {
		return nil
	}

	// Safely extract values
	var data AuthUser

	if id, ok := mapClaim["id"].(string); ok {
		data.ID = id
	}
	if email, ok := mapClaim["email"].(string); ok {
		data.Email = email
	}
	if username, ok := mapClaim["username"].(string); ok {
		data.Username = username
	}
	if name, ok := mapClaim["name"].(string); ok {
		data.Name = name
	}
	if exp, ok := mapClaim["exp"].(float64); ok {
		data.Exp = exp
	}

	// Handle roles as a nested object
	if roles, ok := mapClaim["roles"].([]interface{}); ok {
		for _, role := range roles {
			if roleMap, ok := role.(map[string]interface{}); ok {
				roleData := Role{}
				if id, ok := roleMap["id"].(string); ok {
					roleData.ID = id
				}
				if name, ok := roleMap["name"].(string); ok {
					roleData.Name = name
				}
				if slug, ok := roleMap["slug"].(string); ok {
					roleData.Slug = slug
				}
				data.Roles = append(data.Roles, roleData)
			}
		}
	}

	return &data
}

// // HasRole check if user has role: specific role
// func HasRole(authUser *AuthUser, role string, selectedRole string) bool {
// 	// check if auth user roles is empty
// 	if len(authUser.Roles) == 0 {
// 		return false
// 	}

// 	// iterate through the roles to check for the given role value in slug or name
// 	isFound := false
// 	for _, authRole := range authUser.Roles {
// 		if role == authRole.Slug || authRole.Name == role {
// 			isFound = true
// 		}
// 	}
// 	if role != selectedRole {
// 		isFound = false
// 	}
// 	return isFound
// }

// // HasRoles check if user has roles: multiple roles
// func HasRoles(authUser *AuthUser, roles []string, selectedRole string) bool {
// 	// check if auth user roles is empty
// 	if len(authUser.Roles) == 0 {
// 		return false
// 	}

// 	// Convert roles to a map for faster lookup
// 	roleSet := make(map[string]struct{})
// 	for _, role := range roles {
// 		roleSet[role] = struct{}{}
// 	}

// 	// iterate through the roles to check for the given role value in slug or name
// 	isFound := false
// 	for _, authRole := range authUser.Roles {
// 		if _, found := roleSet[authRole.Slug]; found {
// 			isFound = true
// 		}
// 		if _, found := roleSet[authRole.Name]; found {
// 			isFound = true
// 		}
// 	}
// 	if _, found := roleSet[selectedRole]; !found {
// 		isFound = false
// 	}

// 	return isFound
// }

func CanAccessRoles(authUser *AuthUser, allowedRoles []string, selectedRole string) bool {
	// Ensure the selected role is part of the user's assigned roles
	roleSet := make(map[string]bool)
	for _, role := range authUser.Roles {
		roleSet[role.Name] = true
	}
	if !roleSet[selectedRole] {
		return false
	}

	// Check if the selected role is allowed to access the function
	for _, role := range allowedRoles {
		if selectedRole == role {
			return true
		}
	}
	return false
}
