#!/bin/bash

bearLocation="~/Library/Group\ Containers/9K33E3U3T4.net.shinyfrog.bear/Application\ Data/database.sqlite"
newLocation="${PWD}/data/database.sqllite"

echo "Copying from: $bearLocation"
echo "Copying to: $newLocation"

if [[ -e $newLocation ]]; then
    echo "Old file detected, removing"
    rm $newLocation
fi
echo "Copying File"
eval cp "$bearLocation" "$newLocation"