// Package yamlutil provide an easy method to render the yaml file.
//
// The type of jsonschema:
// 1. Based type:
//  1. string
//  2. boolean
//  3. integer(number)
//  4. float(number)
//
// 2. Map or Struct：
//  1. Key must a string.
//  2. Values can be any type.
//
// 3. Array
//
// About the docs format
//
//	The output yaml files was follow these format:
//	1. Fields was define topped on fields.
//	2. Fields type was define inline at behind of fields. And follow these format: fields: "" # [string, default=""]
//	3. Field descriptions are displayed below the field as comments
//	4. Only based fields has type define
//
// Example:
//
// ```golang
//
//			type People struct {
//			    Name string `yaml:"Name" default="abc" describe:"The name"`
//			    Age int `yaml:"Age" default="456" describe:"The age of people"`
//	         Addresses []Address `yaml:"Addresses" describe:"The addresses"`
//			}
//
//	     type Address struct {
//	         Name string `yaml:"Name" describe:"The name of address"
//	     }
//
//		 func main() {
//		     people := People{
//		         Name: "test123",
//		         Age:  123,
//		     }
//
//		     resBytes, err := yamlutils.Render(people)
//		     fmt.Println(string(resBytes))
//		     fmt.Println(err) // should be
//		 }
//
// ```
//
// This file will output like:
// ```yaml
//
//	# The name
//
// Name: test123 # [string, default="abc"]
//
// # The age of people
// Age: 123 # [integer, default=456]
//
// # The addresses
// Addresses:
//   - # The name of address
//     Name: "" # [string]
//
// ```
// For output value, first used input struct value, if was nil(if ptr), used the struct default tag as the value.
//
// yamlutil used golang reflect and go tag to parse document.
package yamlutil
