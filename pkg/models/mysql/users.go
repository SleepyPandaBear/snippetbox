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
    var id int
    var hashedPassword []byte
    row := m.DB.QueryRow("SELECT id, hashed_password FROM users WHERE email = ?", email)
    err := row.Scan(&id, &hashedPassword)
    if err == sql.ErrNoRows {
        return 0, models.ErrInvalidCredentials
    } else if err != nil {
        return 0, err
    }

    err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(password))
    if err == bcrypt.ErrMismatchedHashAndPassword {
        return 0, models.ErrInvalidCredentials
    } else if err != nil {
        return 0, err
    }

    return id, nil
}

func (m *UserModel) Get(id int) (*models.User, error) {
    s := &models.User{}

    stmt := `SELECT id, name, email, created FROM users WHERE id = ?`

    err := m.DB.QueryRow(stmt, id).Scan(&s.ID, &s.Name, &s.Email, &s.Created)
    if err == sql.ErrNoRows {
        return nil, models.ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    return s, nil
}

func (m *UserModel) ChangePassword(id int, currentPassword, newPassword string) error {
    var currentHashedPassword []byte
    row := m.DB.QueryRow("SELECT hashed_password FROM users WHERE id = ?", id)
    err := row.Scan(&currentHashedPassword)
    if err != nil {
        return err
    }

    err = bcrypt.CompareHashAndPassword(currentHashedPassword, []byte(currentPassword))
    if err == bcrypt.ErrMismatchedHashAndPassword {
        return models.ErrInvalidCredentials
    } else if err != nil {
        return err
    }

    newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
    if err != nil {
        return err
    }

    stmt := "UPDATE users SET hashed_password = ? WHERE id = ?"
    _, err = m.DB.Exec(stmt, string(newHashedPassword), id)

    return err
}
