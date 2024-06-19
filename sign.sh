#!/bin/sh
#  Script.sh
#  sign
path=$(cd "$(dirname "$0")";pwd)
echo $path
cd $path
codesign -f -s "Wuly, Inc." -v --deep "build/bin/terminal.app"
echo -n "检查签名？（y/n）"
read is_sign
if [[ $is_sign  = "y" ]]; then
  spctl --verbose=4 --assess --type execute build/bin/terminal.app
fi
