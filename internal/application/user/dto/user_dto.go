package dto

// UpdateProfileRequest contains fields the user can update.
type UpdateProfileRequest struct {
	UserID      string  `validate:"required,uuid4"`
	DisplayName *string `validate:"omitempty,min=2,max=100"`
	AvatarURL   *string `validate:"omitempty,url"`
	Bio         *string `validate:"omitempty,max=500"`
}

// UserProfileResponse is the safe user data returned to clients.
type UserProfileResponse struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Bio         string `json:"bio"`
	Status      string `json:"status"`
	Role        string `json:"role"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}
