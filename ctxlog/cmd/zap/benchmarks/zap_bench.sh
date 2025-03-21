#!/bin/bash

tests=("Zap" "ZapWithZax")

for test in "${tests[@]}"
do
    go test -bench="^Benchmark${test^}$" -benchmem -run="^$" > "${test}.log"
done
