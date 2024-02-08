package _components

import "fmt"


type ViewBuilder struct {
	Title string
	ViewComponents []string
}

func NewViewBuilder(title string, viewComponents []string) *ViewBuilder {
	return &ViewBuilder{
		Title: title,
		ViewComponents: []string{},
	}
}

func (v *ViewBuilder) Build() string {
	viewString := "";
	for _, component := range v.ViewComponents {
		viewString += component
	}
	return fmt.Sprintf(`
		<!DOCTYPE html>
		<html>
		<head>
			<title>CFASuite</title>
		</head>
		<body>
			%s
		</body>
	`, viewString)
}

