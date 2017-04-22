NAME="rest2command"
VERSION=$(git describe --abbrev=0 --tags)
PARSED_VERSION=$(echo $VERSION | sed "s/v//" | sed "s/\./_/g")
PACKAGE="${NAME}_${PARSED_VERSION}-1"
BINARY="rest2command-linux-amd64"

rm -rf target/

echo "Package: "$PACKAGE
mkdir -p target/${PACKAGE}/usr/bin
mkdir -p target/${PACKAGE}/etc/init.d/
mkdir -p target/${PACKAGE}/etc/rest2command

cp dist/${BINARY} target/${PACKAGE}/usr/bin/rest2command
cp dist/rest2command.sh target/${PACKAGE}/etc/init.d/rest2command
cp configuration.json target/${PACKAGE}/etc/rest2command/

mkdir -p target/${PACKAGE}/DEBIAN
cp control target/${PACKAGE}/DEBIAN/control

sed -i "" "s/_PACKAGE_NAME_/${NAME}/g" target/${PACKAGE}/DEBIAN/control
sed -i "" "s/_VERSION_/${VERSION}-1/g" target/${PACKAGE}/DEBIAN/control

echo "Permissions: "
chmod +x target/${PACKAGE}/usr/bin/rest2command
chmod +x target/${PACKAGE}/etc/init.d/rest2command

cd target && dpkg-deb --build target/${PACKAGE}
mv ${PACKAGE}_amd64.deb dist/