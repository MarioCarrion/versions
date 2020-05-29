module fixture.com/new_module_replace

go 1.14

require (
	github.com/MarioCarrion/nit v1.23.3
)

replace (
	github.com/MarioCarrion/nit => replaced/MarioCarrion/nit v9.0.0
)
