#!/bin/sh

if [ "$#" -ne 2 ]; then
  echo "Usage: $0 <FILTER> <COUNT>"
  echo "Example 1: $0 . 10"
  echo "Example 2: $0 Base16Uint16 10"
  exit 1
fi

FILTER=$1
COUNT=$2

# Replace any '/' with '_'.
# Replace '/' with '_'
FILTER_AS_FILENAME="${FILTER//\//_}"

OUT_STRCONV=".strconv_$FILTER_AS_FILENAME.txt"
OUT_PARSEINT=".parseint_$FILTER_AS_FILENAME.txt"
FUNC_STRCONV="strconv"
FUNC_PARSEINT="parseint"

echo "Benchmarking function $FUNC_STRCONV to $OUT_STRCONV"
go test \
  -test.timeout=0 \
  -benchmem \
  -bench $FILTER \
  -benchfunc $FUNC_STRCONV \
  -count $COUNT \
  | tee $OUT_STRCONV
clear


echo "Benchmarking function $FUNC_PARSEINT to $OUT_PARSEINT"
go test \
  -test.timeout=0 \
  -benchmem \
  -bench $FILTER \
  -benchfunc $FUNC_PARSEINT \
  -count $COUNT \
  | tee $OUT_PARSEINT
clear

go run golang.org/x/perf/cmd/benchstat $OUT_STRCONV $OUT_PARSEINT