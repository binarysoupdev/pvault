module pvault

go 1.25.0

require (
	github.com/atotto/clipboard v0.1.4
	github.com/binarysoupdev/cryptool v1.0.0
	github.com/binarysoupdev/go-commando v1.4.0
	github.com/binarysoupdev/go-extensions v0.1.0
	github.com/binarysoupdev/got-style v1.1.0
	github.com/binarysoupdev/tinsel v1.0.0
	github.com/google/uuid v1.6.0
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/stretchr/testify v1.11.1
	golang.org/x/term v0.44.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/sys v0.46.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/binarysoupdev/go-commando => ../go-commando

replace github.com/binarysoupdev/cryptool => ../cryptool

replace github.com/binarysoupdev/go-extensions => ../go-extensions
