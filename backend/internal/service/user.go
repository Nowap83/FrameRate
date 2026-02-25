package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/Nowap83/FrameRate/backend/internal/dto"
	"github.com/Nowap83/FrameRate/backend/internal/model"
	"github.com/Nowap83/FrameRate/backend/internal/repository"
	"github.com/Nowap83/FrameRate/backend/internal/utils"
	"go.uber.org/zap"
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

// fetches all users (for admin)
func (s *UserService) GetAllUsers(page, limit int) (*dto.PaginatedUsersResponse, error) {
	users, total, err := s.userRepo.GetAllUsers(page, limit)
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.ToUserResponse(user))
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := &dto.PaginatedUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	return response, nil
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

	// fetch favorites
	favorites, _ := s.movieRepo.GetFavoriteMovies(userID, 4)
	response.Favorites = favorites

	return response, nil
}

// fetches paginated watched films with their ratings for a given user
func (s *UserService) GetMyFilms(userID uint, page, limit int) (*dto.PaginatedMoviesResponse, error) {
	moviesWithRatings, total, err := s.movieRepo.GetWatchedFilmsWithRatings(userID, page, limit)
	if err != nil {
		return nil, errors.New("failed to fetch user films")
	}

	var movieResponses []dto.MovieListResponse
	for _, m := range moviesWithRatings {
		movieResponses = append(movieResponses, dto.MovieListResponse{
			ID:                m.ID,
			TmdbID:            m.TmdbID,
			Title:             m.Title,
			ReleaseYear:       m.ReleaseYear,
			PosterURL:         m.PosterURL,
			AverageUserRating: 0.0,
			UserRating:        m.UserRating,
		})
	}

	totalPages := int((total + int64(limit) - 1) / int64(limit))

	response := &dto.PaginatedMoviesResponse{
		Movies:     movieResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

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
		utils.Log.Info("Updating Favorite Films", zap.Int("Count", len(input.FavoriteFilms)))
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

// updates the user's avatar
func (s *UserService) UpdateAvatar(userID uint, file *multipart.FileHeader) (string, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return "", ErrUserNotFound
	}

	// create uploads directory if not exists (double check)
	uploadDir := "./uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", errors.New("failed to create upload directory")
	}

	// generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("avatar_%d_%d%s", user.ID, time.Now().UnixNano(), ext)
	filepath := filepath.Join(uploadDir, filename)

	// save file
	src, err := file.Open()
	if err != nil {
		return "", errors.New("failed to open uploaded file")
	}
	defer src.Close()

	dst, err := os.Create(filepath)
	if err != nil {
		return "", errors.New("failed to create destination file")
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", errors.New("failed to save file")
	}

	// update DB
	fileURL := "/uploads/avatars/" + filename

	if err := s.userRepo.UpdateFields(userID, map[string]interface{}{"profile_picture_url": fileURL}); err != nil {
		return "", errors.New("failed to update profile picture in database")
	}

	return fileURL, nil
}

// checks if a username is available
func (s *UserService) CheckUsernameAvailability(username string) (bool, error) {
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return false, nil
	}
	return true, nil
}
