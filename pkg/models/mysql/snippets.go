package mysql

import (
    "database/sql"
    "spbear/snippetbox/pkg/models"
)

type SnippetModel struct {
    DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
    // `?` is a placeholder for actual data when we call DB.Exec
    stmt := `INSERT INTO snippets (title, content, created, expires) VALUES(?,
    ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

    // Execute our statement with our data
    result, err := m.DB.Exec(stmt, title, content, expires)
    if err != nil {
        return 0, err
    }

    // Get the id of the last inserted record
    id, err := result.LastInsertId()
    if err != nil {
        return 0, err
    }

    return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
    // Our SQL statement for selecting specific non expired snippet
    stmt := `SELECT id, title, content, created, expires FROM snippets WHERE
    expires > UTC_TIMESTAMP() AND id = ?`


    // Execute statement, QueryRow returns only one result or none!
    // We can also chain functions: QueryRow("SELECT id, ...", id).Scan(&s.ID, ...)
    row := m.DB.QueryRow(stmt, id)

    // Empty struct that will be filled with data
    s := &models.Snippet{}

    // "Scan" resulted row into empty struct `s`. Our driver automatically
    // converts values to the right types (that we use for Snippet struct), for
    // converting date specific structures to time.Time we need to include
    // parseTime=true flag when connecting to the database.
    err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

    if err == sql.ErrNoRows {
        return nil, models.ErrNoRecord
    } else if err != nil {
        return nil, err
    }

    return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
    stmt := `SELECT id, title, content, created, expires FROM snippets WHERE
    expires > UTC_TIMESTAMP() ORDER BY created DESC LIMIT 10`

    rows, err := m.DB.Query(stmt)
    if err != nil {
        return nil, err
    }

    // Vital to call this, otherwise it will just use one connection of
    // database pool (can happen that all connections will be used)
    defer rows.Close()

    snippets := []*models.Snippet{}

    for rows.Next() {
        s := &models.Snippet{}
        err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
        if err != nil {
            return nil, err
        }

        snippets = append(snippets, s)
    }

    // We can't assume that successful iteration was completed over whole
    // resultset
    if err = rows.Err(); err != nil {
        return nil, err
    }

    return snippets, nil
}
