package service

// Person represents an RSS‐feed author for documentation.
// swagger:model Person
type Person struct {
	// example: Jane Doe
	Name string `json:"name"`
	// example: jane@example.com
	Email string `json:"email,omitempty"`
}

// Authors is a list of Person.
// swagger:model Authors
type Authors []Person
