.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/goal-create lambdas/goal/create/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/goal-edit lambdas/goal/edit/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/goal-get lambdas/goal/get/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/goal-list lambdas/goal/list/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/goal-delete lambdas/goal/delete/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/user-signup lambdas/user/signup/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/user-verify lambdas/user/verify/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/user-login lambdas/user/login/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/user-logout lambdas/user/logout/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/user-delete lambdas/user/delete/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
