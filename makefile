install:
	@go install .

generate-gifs:
	mkdir -p tmp

	# ROOT
	cd tmp && vhs ../doc/tapes/doctor.tape
	cd tmp && vhs ../doc/tapes/init.tape

	# DOMAIN
	cd tmp/my-project && vhs ../../doc/tapes/domain-new.tape
	cd tmp/my-project && vhs ../../doc/tapes/domain-ls.tape

	# PORT
	echo "package port\n\ntype UserController interface {}\ntype UserService interface {}\ntype UserRepository interface {}" > tmp/my-project/internal/port/user.go
	echo "package port\n\ntype CompanyController interface {}\ntype CompanyService interface {}\ntype CompanyRepository interface {}" > tmp/my-project/internal/port/company.go
	echo "package port\n\ntype Logger interface {}" > tmp/my-project/internal/port/logger.go
	cd tmp/my-project && vhs ../../doc/tapes/port-ls.tape

	# SERVICE
	cd tmp/my-project && vhs ../../doc/tapes/service-new.tape
	mkdir -p tmp/my-project/internal/domain/core/service/company
	echo "" > tmp/my-project/internal/domain/core/service/company/company.go
	cd tmp/my-project && vhs ../../doc/tapes/service-ls.tape

	# APPLICATION
	cd tmp/my-project && vhs ../../doc/tapes/app-new.tape
	mkdir -p tmp/my-project/internal/domain/core/application/companyrest
	echo "" > tmp/my-project/internal/domain/core/application/companyrest/companyrest.go
	cd tmp/my-project && vhs ../../doc/tapes/app-ls.tape

	# INFRASTRUCTURE
	cd tmp/my-project && vhs ../../doc/tapes/infra-new.tape
	mkdir -p tmp/my-project/internal/infrastructure/companyrepository
	echo "" > tmp/my-project/internal/infrastructure/companyrepository/companyrepository.go
	cd tmp/my-project && vhs ../../doc/tapes/infra-ls.tape

	# PACKAGE
	cd tmp/my-project && vhs ../../doc/tapes/pkg-new.tape
	mkdir -p tmp/my-project/internal/pkg/examplepkg1
	echo "" > tmp/my-project/internal/pkg/examplepkg1/examplepkg1.go
	mkdir -p tmp/my-project/internal/pkg/examplepkg2
	echo "" > tmp/my-project/internal/pkg/examplepkg2/examplepkg2.go
	mkdir -p tmp/my-project/pkg/external
	echo "" > tmp/my-project/pkg/external/external.go
	cd tmp/my-project && vhs ../../doc/tapes/pkg-ls.tape

	# ENTRYPOINT
	cd tmp/my-project && vhs ../../doc/tapes/cmd-new.tape
	mkdir -p tmp/my-project/cmd/cli
	echo "" > tmp/my-project/cmd/cli/main.go
	mkdir -p tmp/my-project/cmd/gui
	echo "" > tmp/my-project/cmd/gui/main.go
	cd tmp/my-project && vhs ../../doc/tapes/cmd-ls.tape

	# RUN
	cd tmp/my-project && vhs ../../doc/tapes/run.tape
	cd tmp/my-project && vhs ../../doc/tapes/run-with-env.tape

	# TREE
	cd tmp/my-project && vhs ../../doc/tapes/tree.tape
	cd tmp/my-project && tree .

	# CLEAN
	rm -rf tmp/my-project