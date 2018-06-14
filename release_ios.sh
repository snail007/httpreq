#/bin/bash
VER="v1.1"
rm -rf httpreq-ios-*.tar.gz
rm -rf ios
mkdir ios

#ios  XCode required
gomobile bind -v -target=ios -ldflags="-s -w"
mv Httpreq.framework ios
cp README.md ios
tar zcfv httpreq-ios-${VER}.tar.gz ios
rm -rf ios

echo "done."
