#!/bin/sh

# if ! [ -f "simulate_nuclei.go" ]; then 
#    echo -e "simulate_nuclei.go not found"
#    exit
# else
#    echo -e "[*] Compiling simulate_nuclei.go\n" ;
#    go build ./snuclei -sizedwg .go
# fi

if ! [ -f "sysusage.go" ]; then 
   echo -e "sysusage.go not found"
   exit
else
   echo -e "[*] Compiling sysusage.go\n\n" ;
   go build ./sysusage.go
fi

# This depends mostly on network conditions and available RAM
# This should be in between 500 to 5000 for nuclei that too depending on host
CONCURRENCY=${1:2000}
URL=""

function checkalive(){
   curl -sw '%{http_code}' $1 -k -o /dev/null -m 5
   if  [ 0 -eq $? ]; then
   echo -e "URL Reachable\n"
   else
   echo -e "\n $1 Unreachable exiting\n" ;
   exit
   fi
}

if [ -z "$1" ]; then 
   echo -e "Concurrency not set. defaulting to 2000"
   CONCURRENCY="2000"
else
   CONCURRENCY=$1
   echo -e "[*] Target URL is $CONCURRENCY\n"
fi

if [ -z "$2" ]; then 
   echo -e "Target URL not given, defaulting to http://192.168.29.119:9000/ "
   checkalive "http://192.168.29.119:9000/"
else
   URL=$(printf " -u $2 ")
   echo -e "[*] Target URL is $2\n"
   checkalive $2
fi

# echo -e "[*] Benchmarking Nuclei With Unlimited Goroutines\n";

# echo -e "[+] Nuclei with 5k hosts\n"
# ./sysusage ./snuclei -sizedwg  -count 5000 $URL
# echo -e "\n"

# echo -e "[+] Nuclei with 10k hosts \n"
# ./sysusage ./snuclei -sizedwg  -count 10000 $URL
# echo -e "\n"

# echo -e "[+] Nuclei with 25k hosts \n"
# ./sysusage ./snuclei -sizedwg  -count 25000 $URL
# echo -e "\n"

# echo -e "[+] Nuclei with 50k hosts \n"
# ./sysusage ./snuclei -sizedwg  -count 50000 $URL
# echo -e "\n"

# echo -e "[+] Nuclei with 75k hosts \n"
# ./sysusage ./snuclei -sizedwg  -count 75000 $URL
# echo -e "\n"

# echo -e "[+] Nuclei with 100k hosts \n"
# ./sysusage ./snuclei -sizedwg  -count 100000 $URL
# echo -e "\n"

echo -e "[*] Benchmarking Nuclei With Fixed Goroutines ($CONCURRENCY) \n";


echo -e "[+] Nuclei with 5k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 5000 -fixed -c $CONCURRENCY 
echo -e "\n"

echo -e "[+] Nuclei with 10k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 10000 -fixed -c $CONCURRENCY 
echo -e "\n"

echo -e "[+] Nuclei with 25k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 25000 -fixed -c $CONCURRENCY
echo -e "\n"

echo -e "[+] Nuclei with 50k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 50000 -fixed -c $CONCURRENCY
echo -e "\n"

echo -e "[+] Nuclei with 75k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 75000 -fixed -c $CONCURRENCY
echo -e "\n"

echo -e "[+] Nuclei with 100k hosts with $CONCURRENCY Goroutines\n"
./sysusage ./snuclei -sizedwg  $URL -count 100000 -fixed -c $CONCURRENCY
echo -e "\n"


echo -e "[+] Done\n"