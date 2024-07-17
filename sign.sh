#!/bin/sh
wails build clearn
#  Script.sh
#  sign
#path=$(cd "$(dirname "$0")";pwd)
#echo $path
#cd $path
#codesign -f -s "Wuly, Inc." -v --deep "build/bin/terminal.app"
#echo -n "检查签名？（y/n）"
#read is_sign
#if [[ $is_sign  = "y" ]]; then
#  spctl --verbose=4 --assess --type execute build/bin/terminal.app
#fi

cd build/bin
rm -rf *.app
rm *.dmg
create-dmg --volname terminal \
           --window-pos 200 120 \
           --window-size 800 400 \
           --icon-size 100 \
           --icon terminal.app 200 190 \
           --hide-extension terminal.app \
           --app-drop-link  600 185 \
           terminal_setup.dmg \
           terminal.app
