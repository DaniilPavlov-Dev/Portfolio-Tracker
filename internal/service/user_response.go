package service

import "PFnPTA/internal/model"

func ToUserResponse(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID: user.ID,
		Email: user.Email,
		CreatedAt: user.CreatedAt,
	}
}