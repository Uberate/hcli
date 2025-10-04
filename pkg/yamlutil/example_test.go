package yamlutil

import (
	"fmt"
	"testing"
)

// Example struct from the documentation
type People struct {
	Name string `yaml:"Name" default:"abc" describe:"The name"`
	Age  int    `yaml:"Age" default:"456" describe:"The age of people"`
}

// Example with pointer fields
type PeopleWithPointer struct {
	Name *string `yaml:"Name" default:"abc" describe:"The name"`
	Age  *int    `yaml:"Age" default:"456" describe:"The age of people"`
}

// More complex example with arrays and nested structs
type Address struct {
	Street string `yaml:"street" default:"123 Main St" describe:"Street address"`
	City   string `yaml:"city" default:"Anytown" describe:"City name"`
	Zip    string `yaml:"zip" default:"12345" describe:"ZIP code"`
}

type Person struct {
	Name      string    `yaml:"name" default:"John Doe" describe:"Full name"`
	Age       int       `yaml:"age" default:"30" describe:"Age in years"`
	Active    bool      `yaml:"active" default:"true" describe:"Whether active"`
	Tags      []string  `yaml:"tags" describe:"List of tags"`
	Address   Address   `yaml:"address" describe:"Physical address"`
	Addresses []Address `yaml:"addresses" describe:"Multiple addresses"`
}

func TestGenerateYAMLFromStruct(t *testing.T) {
	// Test with the simple example from documentation
	person := People{
		Name: "test123",
		Age:  123,
	}

	yamlOutput, err := GenerateYAMLFromStruct(person)
	if err != nil {
		t.Fatalf("Failed to generate YAML: %v", err)
	}

	fmt.Println("=== Simple Example Output ===")
	fmt.Println(yamlOutput)

	// Test Render function (returns []byte)
	resBytes, err := Render(person)
	if err != nil {
		t.Fatalf("Failed to generate YAML with Render: %v", err)
	}

	fmt.Println("=== Render Output ===")
	fmt.Println(string(resBytes))

	// Test with pointer fields (nil pointers should use defaults)
	name := "pointer_name"
	pointerPerson := PeopleWithPointer{
		Name: &name,
		Age:  nil, // This should use default
	}

	pointerYAML, err := GenerateYAMLFromStruct(pointerPerson)
	if err != nil {
		t.Fatalf("Failed to generate pointer YAML: %v", err)
	}

	fmt.Println("=== Pointer Example Output ===")
	fmt.Println(pointerYAML)

	// Test with more complex example
	complexPerson := Person{
		Name:   "Bob Smith",
		Age:    35,
		Active: true,
		Tags:   []string{"admin", "user"},
		Address: Address{
			Street: "456 Oak Ave",
			City:   "Somewhere",
			Zip:    "67890",
		},
		Addresses: []Address{
			{
				Street: "789 Pine Rd",
				City:   "Elsewhere",
				Zip:    "11223",
			},
		},
	}

	complexYAML, err := GenerateYAMLFromStruct(complexPerson)
	if err != nil {
		t.Fatalf("Failed to generate complex YAML: %v", err)
	}

	fmt.Println("=== Complex Example Output ===")
	fmt.Println(complexYAML)

	// Test with zero values to see defaults
	emptyPerson := Person{}
	emptyYAML, err := GenerateYAMLFromStruct(emptyPerson)
	if err != nil {
		t.Fatalf("Failed to generate empty YAML: %v", err)
	}

	fmt.Println("=== Empty Example with Defaults ===")
	fmt.Println(emptyYAML)
}

func TestParseStruct(t *testing.T) {
	person := People{
		Name: "Charlie",
		Age:  40,
	}

	fields, err := ParseStruct(person)
	if err != nil {
		t.Fatalf("Failed to parse struct: %v", err)
	}

	fmt.Println("=== Parsed Fields ===")
	for _, field := range fields {
		fmt.Printf("Name: %s, Type: %s, Default: %s, Describe: %s, Value: %v\n",
			field.Name, field.Type, field.Default, field.Describe, field.Value)
	}
}
