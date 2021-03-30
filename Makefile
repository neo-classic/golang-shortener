run-redis:
	go run main.go -port=8000 -db-provider=redis -redis-url=redis://localhost:6379
