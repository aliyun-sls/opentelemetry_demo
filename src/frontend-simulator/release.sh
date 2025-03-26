#!/bin/bash

sed -i "s/MALL_HOST/$MALL_HOST/g" $(grep MALL_HOST -rl .)