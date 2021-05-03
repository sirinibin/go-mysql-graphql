package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
	"gitlab.com/sirinibin/go-mysql-graphql/config"
)

type SearchCriterias struct {
	Page     uint32                 `bson:"page,omitempty" json:"page,omitempty"`
	Size     uint32                 `bson:"size,omitempty" json:"size,omitempty"`
	SearchBy map[string]interface{} `bson:"search_by,omitempty" json:"search_by,omitempty"`
	SortBy   string                 `bson:"sort_by,omitempty" json:"sort_by,omitempty"`
	ID       string                 `bson:"id,omitempty" json:"id,omitempty"`
	Name     string                 `bson:"name,omitempty" json:"name,omitempty"`
	Email    string                 `bson:"email,omitempty" json:"email,omitempty"`
}

var SortFields = map[string]bool{"id": true, "name": true, "email": true, "created_by": true, "updated_by": true}

func ValidateSortString(str string) error {

	split := strings.Split(str, ",")
	for _, v := range split {
		split2 := strings.Split(v, " ")
		if len(split2) > 2 {
			return errors.New("Invalid sort value")
		} else if len(split2) == 2 {
			if split2[1] != "asc" && split2[1] != "desc" {
				return errors.New("Invalid sort order value ( Note:use asc or desc, default:asc )")
			}
		}
		fieldName := split2[0]
		if !SortFields[fieldName] {
			return errors.New("Invalid field " + fieldName)
		}

	}
	return nil
}

func FindEmployees(pageCriterias *PageCriterias, filter *FilterCriterias) ([]*Employee, error) {

	var employees []*Employee

	page := 1
	size := 10

	if pageCriterias != nil {

		if pageCriterias.Page != nil {
			page = *pageCriterias.Page
		}

		if pageCriterias.Size != nil {
			size = *pageCriterias.Size
		}
	}

	offset := (page - 1) * (size)

	searchString := ""

	args := []interface{}{}

	if filter != nil {

		if filter.ID != nil {
			searchString += "id = ? "
			args = append(args, *filter.ID)

		}

		if filter.Name != nil {
			if len(args) > 0 {
				searchString += " AND "
			}
			searchString += "name like ? "
			args = append(args, *filter.Name)
		}

		if filter.Email != nil {
			if len(args) > 0 {
				searchString += " AND "
			}
			searchString += "email = ? "
			args = append(args, *filter.Email)
		}

		if len(args) > 0 {
			searchString = " WHERE " + searchString
		}

	}

	sortString := ""
	if pageCriterias != nil && pageCriterias.Sort != nil && *pageCriterias.Sort != "" {
		err := ValidateSortString(*pageCriterias.Sort)
		if err != nil {
			return nil, err
		}
		sortString = fmt.Sprintf("ORDER BY %s", *pageCriterias.Sort)
	}
	args = append(args, offset)
	args = append(args, size)

	query := "SELECT id,name,email,created_by,updated_by,created_At,updated_at FROM employee " + searchString + sortString + " limit ?,?"
	res, err := config.DB.Query(query, args...)
	defer res.Close()

	if err != nil {
		return nil, err
	}

	for res.Next() {
		var employee Employee

		var createdAt string
		var updatedAt string
		err := res.Scan(&employee.ID, &employee.Name, &employee.Email, &employee.CreatedBy, &employee.UpdatedBy, &createdAt, &updatedAt)
		if err != nil {
			return nil, err
		}

		layout := "2006-01-02 15:04:05"

		employee.CreatedAt, err = time.Parse(layout, createdAt)
		if err != nil {
			return nil, err
		}

		employee.UpdatedAt, err = time.Parse(layout, updatedAt)
		if err != nil {
			return nil, err
		}

		employees = append(employees, &employee)

	}
	return employees, nil

}

func DeleteEmployee(employeeID string) (int64, error) {

	res, err := config.DB.Exec("DELETE from employee where id=?", employeeID)
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()

}

func IsEmployeeExists(employeeID string) (exists bool, err error) {

	var id uint64

	err = config.DB.QueryRow("SELECT id from employee where id=?", employeeID).Scan(&id)

	return id != 0, err
}

func FindEmployeeByID(id string) (*Employee, error) {

	var createdAt string
	var updatedAt string
	var employee Employee

	err := config.DB.QueryRow("SELECT id,created_by,updated_by,name,email,created_at,updated_at from employee where id=?", id).Scan(&employee.ID, &employee.CreatedBy, &employee.UpdatedBy, &employee.Name, &employee.Email, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	layout := "2006-01-02 15:04:05"

	employee.CreatedAt, err = time.Parse(layout, createdAt)
	if err != nil {
		return &employee, err
	}

	employee.UpdatedAt, err = time.Parse(layout, updatedAt)
	if err != nil {
		return &employee, err
	}

	return &employee, err
}

func (employee *Employee) IsEmailExists() (exists bool, err error) {

	var id uint64

	//i, err := strconv.Atoi("-42")

	if employee.ID != "" {
		//Old Record
		err = config.DB.QueryRow("SELECT id from employee where email=? and id!=?", employee.Email, employee.ID).Scan(&id)
	} else {
		//New Record
		err = config.DB.QueryRow("SELECT id from employee where email=?", employee.Email).Scan(&id)
	}
	return id != 0, err
}

func (employee *Employee) Validate(scenario string) error {

	if scenario == "update" {
		if employee.ID == "" {
			return errors.New("ID is required")
		}
		exists, err := IsEmployeeExists(employee.ID)
		if err != nil || !exists {
			return err
		}

	}

	if govalidator.IsNull(employee.Name) {
		return errors.New("Name is required")
	}

	if govalidator.IsNull(employee.Email) {
		return errors.New("E-mail is required")
	}

	emailExists, err := employee.IsEmailExists()
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if emailExists {
		return errors.New("E-mail is Already in use")
	}

	return nil
}

func (employee *Employee) Insert() error {

	res, err := config.DB.Exec("insert into employee (name,created_by,updated_by, email,created_at,updated_at) VALUES (?, ?, ?, ?, ?, ?)", employee.Name, employee.CreatedBy, employee.UpdatedBy, employee.Email, employee.CreatedAt, employee.UpdatedAt)
	if err != nil {
		return err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Error %s when finding last insert Id", err)
		return err
	}
	employee.ID = strconv.FormatInt(id, 10)
	log.Print("user.ID:")
	log.Print(employee.ID)
	log.Printf("%d employee created ", rows)

	return nil
}

func (employee *Employee) Update() (*Employee, error) {

	res, err := config.DB.Exec("UPDATE employee SET name=?, updated_by=? ,email=?, updated_at=? WHERE id=?", employee.Name, employee.UpdatedBy, employee.Email, employee.UpdatedAt, employee.ID)
	if err != nil {
		return nil, err
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when finding rows affected", err)
		return nil, err
	}

	employee, err = FindEmployeeByID(employee.ID)
	if err != nil {
		return nil, err
	}

	log.Print("user.ID:")
	log.Print(employee.ID)
	log.Printf("%d employee updated ", rows)

	return employee, nil
}
