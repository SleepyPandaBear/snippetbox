package main

import (
    "spbear/snippetbox/pkg/models"    
)

type templateData struct {
    Snippet *models.Snippet
    Snippets []*models.Snippet
}
