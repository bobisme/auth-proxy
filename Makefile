run:
	cat <<EOF |
	-r '\.go$$' -s -- go run main.go
	-r '\.go$$' -s -- go run example/main.go
	EOF
	reflex --decoration=fancy -c -
