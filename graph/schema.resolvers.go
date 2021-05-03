package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"time"

	"gitlab.com/sirinibin/go-mysql-graphql/config"
	"gitlab.com/sirinibin/go-mysql-graphql/graph/generated"
	"gitlab.com/sirinibin/go-mysql-graphql/graph/model"
)

func (r *mutationResolver) Register(ctx context.Context, input model.RegisterRequest) (*model.RegisterResponse, error) {
	currentTime := time.Now()

	user := &model.User{
		Name:      input.Name,
		Username:  input.Username,
		Email:     input.Email,
		Password:  model.HashPassword(input.Password),
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}

	// Validate data
	err := user.Validate()
	if err != nil {
		return nil, err
	}

	err = user.Insert()
	if err != nil {
		return nil, err
	}
	r.users = append(r.users, user)

	response := &model.RegisterResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
	return response, nil
}

func (r *mutationResolver) Authorize(ctx context.Context, input model.AuthorizeRequest) (*model.AuthCodeResponse, error) {
	auth := &model.AuthorizeRequest{
		Username: input.Username,
		Password: input.Password,
	}

	err := auth.Authenticate()
	if err != nil {
		return nil, err
	}

	authCode, err := auth.GenerateAuthCode()
	if err != nil {
		return nil, err
	}
	return &authCode, nil
}

func (r *mutationResolver) Accesstoken(ctx context.Context, input model.AccesstokenRequest) (*model.AccessTokenResponse, error) {
	accessTokenRequest := &model.AccesstokenRequest{
		AuthCode: input.AuthCode,
	}

	tokenClaims, err := model.AuthenticateByAuthCode(accessTokenRequest.AuthCode)
	if err != nil {
		return nil, err
	}

	accessToken, err := model.GenerateAccesstoken(tokenClaims.Username)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (r *mutationResolver) Refreshtoken(ctx context.Context, input model.RefreshtokenRequest) (*model.AccessTokenResponse, error) {
	refreshTokenRequest := &model.RefreshtokenRequest{
		RefreshToken: input.RefreshToken,
	}

	tokenClaims, err := model.AuthenticateByRefreshToken(refreshTokenRequest.RefreshToken)
	if err != nil {
		return nil, err
	}

	accessToken, err := model.GenerateAccesstoken(tokenClaims.Username)
	if err != nil {
		return nil, err
	}

	return &accessToken, nil
}

func (r *mutationResolver) Logout(ctx context.Context) (string, error) {
	AccessUUID := ctx.Value("AccessUUID")

	deleted, err := config.RedisClient.Del(AccessUUID.(string)).Result()
	if err != nil || deleted == 0 {
		return "", err

	}
	return "Successfully logged out", nil
}

func (r *mutationResolver) CreateEmployee(ctx context.Context, input model.CreateEmployeeRequest) (*model.Employee, error) {
	UserID := ctx.Value("UserID").(string)

	currentTime := time.Now().Local()

	employee := &model.Employee{
		Name:      input.Name,
		Email:     input.Email,
		CreatedBy: UserID,
		UpdatedBy: UserID,
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
	// Validate data
	err := employee.Validate("create")
	if err != nil {
		return nil, err
	}

	err = employee.Insert()
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *mutationResolver) UpdateEmployee(ctx context.Context, input model.UpdateEmployeeRequest) (*model.Employee, error) {
	UserID := ctx.Value("UserID").(string)

	currentTime := time.Now().Local()

	employee := &model.Employee{
		ID:        input.ID,
		Name:      input.Name,
		Email:     input.Email,
		UpdatedBy: UserID,
		UpdatedAt: currentTime,
	}
	// Validate data
	err := employee.Validate("update")
	if err != nil {
		return nil, err
	}

	employee, err = employee.Update()
	if err != nil {
		return nil, err
	}

	return employee, nil
}

func (r *mutationResolver) DeleteEmployee(ctx context.Context, id string) (string, error) {
	res, err := model.DeleteEmployee(id)
	if err != nil || res == 0 {
		return "", errors.New("Unable to delete")
	}

	return "Deleted successfully", nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return r.users, nil
}

func (r *queryResolver) Me(ctx context.Context) (*model.MeResponse, error) {
	UserID := ctx.Value("UserID")
	user, err := model.FindUserByID(UserID.(string))
	if err != nil {
		return nil, err
	}
	user.Password = ""
	response := &model.MeResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	return response, nil
}

func (r *queryResolver) Employees(ctx context.Context, page *model.PageCriterias, filter *model.FilterCriterias) ([]*model.Employee, error) {
	employees, err := model.FindEmployees(page, filter)
	if err != nil {
		return nil, err
	}
	return employees, nil
}

func (r *queryResolver) ViewEmployee(ctx context.Context, id string) (*model.Employee, error) {
	employee, err := model.FindEmployeeByID(id)
	if err != nil {
		return nil, err
	}

	return employee, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
