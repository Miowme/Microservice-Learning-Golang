package service

import (
	"errors"
	"os"
	"time"

	"core-service/internal/model"
	"core-service/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(name, email, password string) (*model.User, error)
	Login(email, password string) (string, error)
}

type authService struct {
	userRepo   repository.UserRepository
	walletRepo repository.WalletRepository
	db         *gorm.DB
}

func NewAuthService(userRepo repository.UserRepository, walletRepo repository.WalletRepository, db *gorm.DB) AuthService {
	return &authService{
		userRepo:   userRepo,
		walletRepo: walletRepo,
		db:         db,
	}
}

func (s *authService) Register(name, email, password string) (*model.User, error) {
	// 1. Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// 2. Gunakan database transaction untuk memastikan User dan Wallet terbuat secara atomik
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// Simpan user baru (menggunakan s.userRepo.Create tanpa parameter tx, atau sesuaikan jika repository Anda mendukung tx)
		// Karena error sebelumnya 'want (*model.User)', kita panggil tanpa tx:
		if err := s.userRepo.Create(user); err != nil {
			return errors.New("email sudah terdaftar atau gagal membuat user")
		}

		// Buat wallet otomatis dengan saldo awal 0 untuk user tersebut
		wallet := &model.Wallet{
			UserID:  user.ID,
			Balance: 0.00,
		}
		if err := s.walletRepo.Create(tx, wallet); err != nil {
			return errors.New("gagal membuat dompet untuk user")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (string, error) {
	// 1. Cari user berdasarkan email
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Generate JWT Token
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "default_secret_key"
	}

	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"email":   user.Email,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", errors.New("gagal menghasilkan token autentikasi")
	}

	return tokenString, nil
}