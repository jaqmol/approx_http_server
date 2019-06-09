#!/bin/sh

# generic section

cd ..
bin=${PWD##*/}
go build -o ./test/${bin}
cd test

# specific section
# TODO: ADJUST FOR HTTP SERVER

mkfifo input.pipe
mkfifo output_0.pipe
mkfifo output_1.pipe

echo
echo 'Please open 2 shells @'${PWD}
echo '  1$ cat < output_0.pipe'
echo '  2$ cat < output_1.pipe'
echo
echo 'Open a 3rd shell and paste:'
echo '  3$ ls > input.pipe'
echo
echo 'Check if ls output are being printed in 1$ and 2$'

DISTRIBUTE=copy \
IN_COUNT=1 \
IN_0=input.pipe \
OUT_COUNT=2 \
OUT_0=output_0.pipe \
OUT_1=output_1.pipe \
./${bin}

rm ${bin}
rm input.pipe
rm output_0.pipe
rm output_1.pipe