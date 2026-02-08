package example

#Person: {
	name!: string & =~"^.{1,50}$"
	age?:  int & >=0 & <=150
	email?: string & =~"^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"
}

Person: #Person

// Example for testing
_example: Person & {
	name:  "John Doe"
	age:   30
	email: "john@example.com"
}
