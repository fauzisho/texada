tr -d '
mv temp1abcxyz $1
sed -i -e 's/[[:blank:]]*$//' -e 's/^[[:blank:]]*//' $1