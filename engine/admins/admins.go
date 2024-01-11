package admins

import (
	"aytp/engine"
	"aytp/graph/model"
	"aytp/models"
	"aytp/utils"
	"errors"
	

	uuid "github.com/satori/go.uuid"
)

func Login(input model.LoginInput) (*model.AuthPayload, error) {
	admin, err := engine.FetchAdmin(&models.Admin{Email: input.Email})
	if err != nil {
		return nil, engine.ErrWrongCreds
	}
	if !utils.CompareHashedString(admin.Password, input.Password) {
		return nil, engine.ErrWrongCreds
	}
	return engine.GenerateToken(admin.ID, admin.Name)
}
func ResetPassword(input model.ResetPasswordInput) (bool, error) {
	admin, err := engine.FetchAdminByToken(input.Token)
	if err != nil {
		return false, err
	}
	var bytes [16]byte
	copy(bytes[:], []byte{0})

	if admin.ID == bytes {
		return false, errors.New("invalid token")
	}

	error := engine.FetchUserByEmail(input.StudentsEmail)
	if !(error == nil) {
		return false, errors.New("email doesnt exists")
	}

	if !(input.NewPassword == input.ConfirmNewPassword) {
		return false, errors.New("passwords do not match")
	}

	newPassword, err := utils.HashString(input.NewPassword)
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	err = utils.DB.Model(&models.Student{}).Where("email = ?", input.StudentsEmail).Update("password", newPassword).Error
	if err != nil {
		return false, err
	}

	return true, nil
}


func AddRealAdmin(input model.NewSchoolAdmin) (bool, error) {
	if input.DevToken != "ngameagamesTokenHere" {
		return false, nil
	}
	err := engine.FetchUserByEmail(input.Email)
	if err == nil {
		return false, errors.New("email already registered")
	}
	password, err := utils.HashString(input.Password)
	if err != nil {
		return false, errors.New("internal encryption error")
	}

	newAdmin := models.Admin{
		SchoolID:   uuid.NewV4(),
		SchoolName: input.SchoolName,
		Name:       input.AdminName,
		Email:      input.Email,
		SchoolRef:  input.Ref,
		Password:   password,
	}

	err = utils.DB.Create(&newAdmin).Error
	if err != nil {
		return false, nil
	}
	return true, nil
}

func ChangeLogo(token, url string) (bool, error) {
	admin, err := engine.FetchAdminByToken(token)
	if err != nil {
		return false, err
	}
	err = utils.DB.Model(&models.Admin{}).Where("id = ?", admin.ID).Update("profile_photo", url).Error
	if err != nil {
		return false, err
	}
	return true, nil
}

func FetchLogo(ref string) (string, error) {
	var admin models.Admin
	err := utils.DB.Model(&models.Admin{}).Where("school_ref = ?", ref).First(&admin).Error
	if err != nil {
		return "", err
	}
	return admin.ProfilePhoto, nil
}

func FetchAllSchools() ([]*model.School, error) {
	adms, err := engine.FetchSchoolAdmins()
	if err != nil {
		return nil, err
	}
	gqlSchools := make([]*model.School, len(adms))
	for i, s := range adms {
		gqlSchools[i] = &model.School{
			ID:   s.SchoolID.String(),
			Name: s.SchoolName,
		}
	}
	return gqlSchools, nil
}
