package loadbalancer

var instances = []string{
	"http://localhost:8080",
	"http://localhost:8081",
}

var current = 0

func GetNextInstance() string {
	instance := instances[current]
	current = (current + 1) % len(instances)
	return instance
}
