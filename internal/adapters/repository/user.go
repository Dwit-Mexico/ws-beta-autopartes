package repository

import (
	"fmt"
	"time"

	"github.com/RomanshkVolkov/test-api/internal/core/domain"
)

func (database *DSNSource) GetAllUsers() ([]domain.UserTableCRUD, error) {
	users := []domain.User{}
	database.DB.Preload("Profile").Model(&domain.User{}).Find(&users)

	crudUsers := []domain.UserTableCRUD{}
	for _, user := range users {
		if user.Profile.Slug == "root" {
			continue
		}
		crudUsers = append(crudUsers, domain.UserTableCRUD{
			ID:       user.ID,
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
			Profile:  user.Profile.Name,
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

	userKitchens := []domain.UsersHasKitchens{}
	database.DB.Model(&domain.UsersHasKitchens{}).Where("user_id = ?", user.ID).Find(&userKitchens)

	kitchenIDs := []uint{}
	for _, userKitchen := range userKitchens {
		kitchenIDs = append(kitchenIDs, userKitchen.KitchenID)
	}

	return domain.EditableUser{
		ID: user.ID,
		UserData: domain.UserData{
			Username: user.Username,
			Name:     user.Name,
			Email:    user.Email,
			IsActive: user.IsActive,
		},
		ShiftID:    user.ShiftID,
		ProfileID:  user.ProfileID,
		KitchenIDs: kitchenIDs,
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
		ProfileID: data.ProfileID,
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

	userKitchens := data.KitchenIDs

	for _, kitchenID := range userKitchens {
		userKitchen := domain.UsersHasKitchens{
			UserID:    user.ID,
			KitchenID: kitchenID,
		}

		// save kitchens one by one
		if err := database.DB.Create(&userKitchen).Error; err != nil {
			return domain.User{}, err
		}
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
	user.ProfileID = data.ProfileID
	user.IsActive = data.IsActive

	if err := database.DB.Save(&user).Error; err != nil {
		return domain.User{}, err
	}

	// delete all kitchens
	database.DB.Where("user_id = ?", user.ID).Delete(&domain.UsersHasKitchens{})

	// save kitchens one by one
	for _, kitchenID := range data.KitchenIDs {
		userKitchen := domain.UsersHasKitchens{
			UserID:    user.ID,
			KitchenID: kitchenID,
		}

		if err := database.DB.Create(&userKitchen).Error; err != nil {
			return domain.User{}, err
		}
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

func (database *DSNSource) GetProfileByID(id uint) (domain.ProfileWithDetails, error) {
	profile := domain.UserProfiles{}
	// left join many to many profile with permissions
	database.DB.Model(&domain.UserProfiles{}).Where("id = ?", id).First(&profile)

	permissions := []domain.ProfilesHasPermissions{}

	database.DB.Model(&domain.ProfilesHasPermissions{}).Preload("Permission").Where("profile_id = ?", id).Find(&permissions)

	permisssionsMapped := []domain.PermissionsByProfile{}

	for _, permission := range permissions {
		permisssionsMapped = append(permisssionsMapped, domain.PermissionsByProfile{
			Permission: domain.Permission{
				ID:   permission.PermissionID,
				Name: permission.Permission.Name,
				Path: permission.Permission.Path,
			},
			Writing: permission.Writing,
		})
	}

	return domain.ProfileWithDetails{
		ID:          profile.ID,
		Name:        profile.Name,
		CreatedAt:   profile.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   profile.UpdatedAt.Format("2006-01-02 15:04:05"),
		Permissions: permisssionsMapped,
	}, nil
}

func (database *DSNSource) GetUsersProfiles() ([]domain.UserProfiles, error) {
	profiles := []domain.UserProfiles{}
	database.DB.Order("name ASC").
		Model(&domain.UserProfiles{}).
		Where("slug != ?", "root").
		Find(&profiles)

	return profiles, nil
}

func (database *DSNSource) CreateProfile(data domain.CreateProfile) (domain.UserProfiles, error) {
	profile := domain.UserProfiles{
		Name: CapitalizeAll(data.Name),
		Slug: Slugify(data.Name),
	}

	tx := database.DB.Begin()
	if err := tx.Create(&profile).Error; err != nil {
		tx.Rollback()
		return domain.UserProfiles{}, err
	}

	for _, permission := range data.Permissions {
		permissionsSerialized := domain.ProfilesHasPermissions{
			ProfileID:    profile.ID,
			PermissionID: permission.ID,
			Writing:      permission.Writing,
		}

		if err := tx.Create(&permissionsSerialized).Error; err != nil {
			tx.Rollback()
			return domain.UserProfiles{}, err
		}
	}

	tx.Commit()

	return profile, nil
}

func (database *DSNSource) UpdateProfile(data domain.EditableProfile) (domain.UserProfiles, error) {
	tx := database.DB.Begin()
	if err := tx.Model(&domain.UserProfiles{}).Where("id = ?", data.ID).Updates(domain.UserProfiles{
		Name: data.Name,
	}).Error; err != nil {
		tx.Rollback()
		return domain.UserProfiles{}, err
	}

	// delete all permissions
	if err := database.DB.Where("profile_id = ?", data.ID).Delete(&domain.ProfilesHasPermissions{}).Error; err != nil {
		tx.Rollback()
		return domain.UserProfiles{}, err
	}

	for _, permission := range data.Permissions {
		fmt.Println(permission)
		permissionsSerialized := domain.ProfilesHasPermissions{
			ProfileID:    data.ID,
			PermissionID: permission.ID,
			Writing:      permission.Writing,
		}

		if err := tx.Create(&permissionsSerialized).Error; err != nil {
			tx.Rollback()
			return domain.UserProfiles{}, err
		}
	}

	tx.Commit()

	profile := domain.UserProfiles{}

	database.DB.Model(&domain.UserProfiles{}).Where("id = ?", data.ID).First(&profile)

	return profile, nil
}

func (database *DSNSource) GetPermissions() ([]domain.Permission, error) {
	permissions := []domain.Permission{}
	database.DB.Model(&domain.Permission{}).Where("status = ?", true).Find(&permissions)

	return permissions, nil
}

func (database *DSNSource) GetKitchens() ([]domain.Kitchen, error) {
	kitchens := []domain.Kitchen{}
	database.DB.Model(&domain.Kitchen{}).Find(&kitchens)

	return kitchens, nil
}

func (database *DSNSource) CreateKitchen(kitchen domain.GenericCatalog) (domain.Kitchen, error) {
	createdItem := domain.Kitchen{
		Name: CapitalizeAll(kitchen.Name),
	}
	if err := database.DB.Create(&createdItem).Error; err != nil {
		return domain.Kitchen{}, err
	}

	return createdItem, nil
}

func (database *DSNSource) GetShifts() ([]domain.Shift, error) {
	shifts := []domain.Shift{}
	database.DB.Model(&domain.Shift{}).Find(&shifts)

	return shifts, nil
}

func (database *DSNSource) CreateShift(shift domain.GenericCatalog) (domain.Shift, error) {
	createdItem := domain.Shift{
		Name: CapitalizeAll(shift.Name),
	}
	if err := database.DB.Create(&createdItem).Error; err != nil {
		return domain.Shift{}, err
	}

	return createdItem, nil
}

// common functions

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
