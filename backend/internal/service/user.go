package service

import (
	"errors"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUsernameTaken     = errors.New("username already taken")
	ErrPasswordIncorrect = errors.New("current password is incorrect")
)

type UserService struct {
	userRepo  repository.UserRepository
	movieRepo *repository.MovieRepository
}

func NewUserService(userRepo repository.UserRepository, movieRepo *repository.MovieRepository) *UserService {
	return &UserService{
		userRepo:  userRepo,
		movieRepo: movieRepo,
	}
}

// fetches a user by their ID
func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

// fetches user profile with statistics
func (s *UserService) GetProfile(userID uint) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	// fetch stats
	watchedCount, _ := s.movieRepo.CountWatched(userID)
	watchedYearCount, _ := s.movieRepo.CountWatchedThisYear(userID)
	reviewsCount, _ := s.movieRepo.CountReviews(userID)
	ratingDist, _ := s.movieRepo.GetRatingDistribution(userID)

	response := &dto.ProfileResponse{
		User: dto.ToUserResponse(user),
		Stats: &dto.UserStats{
			TotalFilms:         watchedCount,
			MoviesThisYear:     watchedYearCount,
			Reviews:            reviewsCount,
			Following:          0,
			Followers:          0,
			RatingDistribution: ratingDist,
		},
	}

	// fetch recent activity
	recentActivity, _ := s.movieRepo.GetRecentWatched(userID, 4)
	response.RecentActivity = recentActivity

	return response, nil
}

// updates user profile information
func (s *UserService) UpdateProfile(userID uint, input dto.UpdateProfileRequest) (*dto.ProfileResponse, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	if input.Username != nil && *input.Username != user.Username {
		existing, err := s.userRepo.GetByUsername(*input.Username)
		if err == nil && existing.ID != userID {
			return nil, ErrUsernameTaken
		}
	}

	updates := make(map[string]interface{})
	if input.Username != nil {
		updates["username"] = *input.Username
	}
	if input.Bio != nil {
		updates["bio"] = *input.Bio
	}
	if input.ProfilePicture != nil {
		updates["profile_picture_url"] = *input.ProfilePicture
	}
	if input.GivenName != nil {
		updates["given_name"] = *input.GivenName
	}
	if input.FamilyName != nil {
		updates["family_name"] = *input.FamilyName
	}
	if input.Location != nil {
		updates["location"] = *input.Location
	}
	if input.Website != nil {
		updates["website"] = *input.Website
	}

	if err := s.userRepo.UpdateFields(userID, updates); err != nil {
		return nil, errors.New("failed to update profile")
	}

	// updates favorite films if provided
	if input.FavoriteFilms != nil {
		if err := s.movieRepo.UpdateFavoriteFilms(userID, input.FavoriteFilms); err != nil {
			return nil, errors.New("failed to update favorites")
		}
	}

	return s.GetProfile(userID)
}

// changes the user's password
func (s *UserService) ChangePassword(userID uint, input dto.ChangePasswordRequest) error {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.CurrentPassword)); err != nil {
		return ErrPasswordIncorrect
	}

	hashedPassword, err := utils.HashPassword(input.NewPassword)
	if err != nil {
		return errors.New("failed to hash password")
	}

	if err := s.userRepo.UpdateFields(userID, map[string]interface{}{"password_hash": string(hashedPassword)}); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

// deletes the user's account
func (s *UserService) DeleteAccount(userID uint) error {
	if err := s.userRepo.Delete(userID); err != nil {
		return errors.New("failed to delete account")
	}
	return nil
}
