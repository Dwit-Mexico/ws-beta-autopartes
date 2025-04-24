package repository

import (
	"time"

	"github.com/RomanshkVolkov/ws-beta-autopartes/internal/core/domain"
)

func (database *DSNSource) GetAllUsers() ([]domain.UserTableCRUD, error) {
	users := []domain.User{}
	database.DB.Preload("Profile").Model(&domain.User{}).Find(&users)

	crudUsers := []domain.UserTableCRUD{}
	for _, user := range users {
		crudUsers = append(crudUsers, domain.UserTableCRUD{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
		})
	}

	return crudUsers, nil
}

func (database *DSNSource) GetUserByID(id uint) (domain.EditableUser, error) {
	user := domain.User{}
	database.DB.Preload("Profile").Model(&domain.User{}).Where("id = ?", id).First(&user)

	if user.ID == 0 {
		return domain.EditableUser{}, nil
	}

	return domain.EditableUser{
		ID: user.ID,
		UserData: domain.UserData{
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
		},
		ShiftID: user.ShiftID,
	}, nil
}

func (database *DSNSource) CreateUser(data domain.CreateUserRequest) (domain.User, error) {
	user := domain.User{
		UserData: domain.UserData{
			Username: data.Username,
			Name:     CapitalizeAll(data.Name),
			Email:    data.Email,
			ShiftID:  data.ShiftID,
			IsActive: data.IsActive,
		},
	}

	hashedPassword, err := HashPassword(data.Password)
	if err != nil {
		return domain.User{}, err
	}
	user.Password = hashedPassword

	// save user
	if err := database.DB.Create(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (database *DSNSource) UpdateUser(data domain.EditableUser) (domain.User, error) {
	user := domain.User{}
	database.DB.Model(&domain.User{}).Where("id = ?", data.ID).First(&user)

	if user.ID == 0 {
		return domain.User{}, nil
	}

	user.Username = data.Username
	user.Name = data.Name
	user.Email = data.Email
	user.ShiftID = data.ShiftID
	user.IsActive = data.IsActive

	if err := database.DB.Save(&user).Error; err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (database *DSNSource) FindByUsername(username string) (domain.User, error) {
	user := domain.User{}
	// with profile
	database.DB.Preload("Profile").Model(&domain.User{}).Where("username = ?", username).First(&user)

	if user.ID == 0 {
		return domain.User{}, nil
	}

	return user, nil
}

func (database *DSNSource) FindByUsernameOrEmail(username, email string) (domain.User, error) {
	user := domain.User{}
	database.DB.Model(&domain.User{}).Where("username = ? OR email = ?", username, email).First(&user)

	if user.ID == 0 {
		return domain.User{}, nil
	}

	return user, nil
}

func (database *DSNSource) FindByID(id uint) (domain.User, error) {
	user := domain.User{}
	database.DB.Model(&domain.User{}).Where("id = ?", id).First(&user)

	if user.ID == 0 {
		return domain.User{}, nil
	}

	return user, nil
}

func (database *DSNSource) FindByUsernameAndOTP(username string) (domain.User, error) {
	user := domain.User{}
	if err := database.DB.Model(&domain.User{}).Where("username = ?", username).First(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (database *DSNSource) FindAndValidateOTP(username string, otp string) (domain.User, map[string][]string, error) {
	schemaError := map[string][]string{}
	user, err := database.FindByUsernameAndOTP(username)
	if err != nil || user.ID == 0 {
		schemaError["username"] = []string{"Usuario no encontrado"}
		return domain.User{}, schemaError, err
	}

	if user.OTP != otp {
		schemaError["code"] = []string{"Código OTP incorrecto"}
		return domain.User{}, schemaError, nil
	}

	if user.OTPExpirationDate.Before(time.Now().UTC()) {
		schemaError["otp"] = []string{"Tu código ha expirado"}
		return domain.User{}, schemaError, nil
	}

	return user, schemaError, nil
}

func (database *DSNSource) NewUser(request *domain.NewUser) (domain.UserData, error) {
	user := domain.User{
		UserData: domain.UserData{
			Username: request.Username,
			Name:     request.Name,
			Email:    request.Email,
		},
		Password: request.Password,
	}

	hashedPassword, err := HashPassword(user.Password)
	if err != nil {
		return domain.UserData{}, err
	}
	user.Password = hashedPassword

	if err := database.DB.Create(&user).Error; err != nil {
		return domain.UserData{}, err
	}

	user.Password = MaskString(user.Password)

	return domain.UserData{
		Username: user.Username,
		Name:     user.Name,
		Email:    user.Email,
	}, nil
}

func (database *DSNSource) SaveOTPCode(username string) (domain.User, error) {
	user, err := database.FindByUsername(username)
	if err != nil {
		return domain.User{}, err
	}

	if user.ID == 0 {
		return user, nil
	}

	otpCode := GenerateOTP(user.Username)
	user.OTP = otpCode
	user.OTPExpirationDate = time.Now().UTC().Add(time.Minute * 1)

	if err := database.DB.Save(&user).Error; err != nil {
		return user, err
	}

	return user, nil
}

func (database *DSNSource) UpdatePassword(userID uint, password string) error {
	user, err := database.FindByID(userID)
	if err != nil {
		return err
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	if err := database.DB.Save(&user).Error; err != nil {
		return err
	}

	return nil
}

func GenerateOTP(txt string) string {
	base := TxtToRandomNumbers(txt + "otp" + CurrentTime())
	return base[:6]
}

func RecordNotFound[T interface{}]() domain.APIResponse[T] {
	return domain.APIResponse[T]{
		Success: false,
		Message: domain.Message{
			En: "Record not found",
			Es: "Registro no encontrado",
		},
	}
}

func HandleDatabaseError[T interface{}](err error, message domain.Message) domain.APIResponse[T] {
	return domain.APIResponse[T]{
		Success: false,
		Message: message,
		Error:   err,
	}
}
