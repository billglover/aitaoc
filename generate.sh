#!/bin/bash

declare -a arr=("scooter"
"visionsOfGrandeur"
"blueSkies"
"darkOcean"
"yoda"
"amin"
"harvey"
"flare"
"ultraViolet"
"sinCityRed"
"eveningNight"
"eXpresso"
"mono"
"white"
"coolSky"
"moonlitAsteroid"
"jShine")

mkdir img
for t in "${arr[@]}"
do
   echo "$t"
   ./aitaoc -height 2436 -width 1125 -numCols 8 -numRows 26 -theme $t -alpha 0.2 -strokeWidth 1 -exponential -sign=false --maxRotation 90
done
mv img iPhoneX

mkdir img
for t in "${arr[@]}"
do
   echo "$t"
   ./aitaoc -height 1334 -width 750 -numCols 8 -numRows 26 -theme $t -alpha 0.2 -strokeWidth 1 -exponential -sign=false --maxRotation 90
done
mv img iPhone8

mkdir img
for t in "${arr[@]}"
do
   echo "$t"
   ./aitaoc -height 1920 -width 1080 -numCols 8 -numRows 26 -theme $t -alpha 0.2 -strokeWidth 1 -exponential -sign=false --maxRotation 90
done
mv img iPhone8Plus

mkdir img
for t in "${arr[@]}"
do
   echo "$t"
   ./aitaoc -height 2224 -width 1668 -numCols 12 -numRows 22 -theme $t -alpha 0.2 -strokeWidth 1 -exponential -sign=false --maxRotation 90
done
mv img iPadPro-10.5

mkdir img
for t in "${arr[@]}"
do
   echo "$t"
   ./aitaoc -height 2732 -width 2048 -numCols 12 -numRows 22 -theme $t -alpha 0.2 -strokeWidth 1 -exponential -sign=false --maxRotation 90
done
mv img iPadPro-12.9
