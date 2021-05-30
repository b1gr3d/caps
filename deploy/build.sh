BUILD_DESTDIR=/tmp/build/${CI_JOB_ID}

# go build
collapse_section_start bootstrap_build_env "${TXT_GREEN}#### BOOTSTRAPPING BUILD ENVIRONMENT ####${TXT_CLEAR}" true
apt install -y jq devscripts dh-make golang
cd ${CI_PROJECT_DIR}
collapse_section_end bootstrap_build_env

collapse_section_start build_go "${TXT_GREEN}#### BUILDING GO BINARY ####${TXT_CLEAR}" true
go build .
check_success "go build failed"
collapse_section_end build_go

# Copy compiled binary into BUILD_DESTDIR for package.sh to pick up
mkdir -p ${BUILD_DESTDIR}/src
cp ./caps ${BUILD_DESTDIR}/src
cp -fpr ./deploy ${BUILD_DESTDIR}
