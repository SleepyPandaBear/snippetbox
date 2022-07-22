package mysql

import (
    "strings"
    "database/sql"
    "spbear/snippetbox/pkg/models"
    "github.com/go-sql-driver/mysql"
    "golang.org/x/crypto/bcrypt"
)

type UserModel struct {
    DB *sql.DB
}

func (m *UserModel) Insert(name, email, password string) error {
    // Hash passoword with 2^12 bcrypt iterations (function also adds random
    // salt)
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    if err != nil {
        return err
    }

    stmt := `INSERT INTO users (name, email, hashed_password, created)
    VALUES(?, ?, ?, UTC_TIMESTAMP())`

    // We could also have a function to check if an email is taken, but that
    // would insert a race condition. If two users would signup with the same 
    // email, both emails could be verified but when inserting into database
    // would be successful for only one of them (other would return an error)
    _, err = m.DB.Exec(stmt, name, email, string(hashedPassword))
    if err != nil {
        if mysqlErr, ok := err.(*mysql.MySQLError); ok {
            if mysqlErr.Number == 1062 && strings.Contains(mysqlErr.Message,
                "users.users_uc_email") {
                return models.ErrDuplicateEmail
            }
        }
    }
    return err
}

func (m *UserModel) Authenticate(email, password string) (int, error) {
    return 0, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
    return nil, nil
}
