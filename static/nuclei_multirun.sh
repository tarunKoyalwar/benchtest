#!/bin/bash

for i in {1..3}
do
  printf "\n[*] Starting run $i\n"
  filename=$(printf "nuclei_run_$i.txt")
  nuclei -u https://scanme.sh -c 50 -o $filename ;
done

wc -l nuclei_run_*