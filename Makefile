NAME := "rest2command"
VERSION := $(shell git describe --abbrev=0 --tags)
PARSED_VERSION := $(shell echo $(VERSION) | sed "s/v//" | sed "s/\./_/g")
PACKAGE := $(NAME)_$(PARSED_VERSION)-1
BINARY := "rest2command-linux-amd64"

deb:
	@echo "Package: "$(PACKAGE)
	$(shell mkdir -p target/$(PACKAGE)/usr/bin)
	$(shell mkdir -p target/$(PACKAGE)/etc/init.d/)
	$(shell mkdir -p target/$(PACKAGE)/etc/rest2command)

	$(shell cp dist/$(BINARY) target/$(PACKAGE)/usr/bin/rest2command)
	$(shell cp dist/rest2command.sh target/$(PACKAGE)/etc/init.d/rest2command)
	$(shell cp configuration.json target/$(PACKAGE)/etc/rest2command/)

	$(shell mkdir -p target/$(PACKAGE)/DEBIAN)
	$(shell cp control target/$(PACKAGE)/DEBIAN/control)

	$(shell sed -i 's/_PACKAGE_NAME_/$(NAME)/g' target/$(PACKAGE)/DEBIAN/control)
	$(shell sed -i 's/_VERSION_/$(VERSION)-1/g' target/$(PACKAGE)/DEBIAN/control)

	@echo "Permissions:"
	$(shell chmod +x target/$(PACKAGE)/usr/bin/rest2command)
	$(shell chmod +x target/$(PACKAGE)/etc/init.d/rest2command)

	$(shell cd target && dpkg-deb --build target/$(PACKAGE))

	$(shell mv $(PACKAGE)_amd64.deb dist/)

clean:
	$(shell rm -rf target/)
