#!/bin/sh
6g experiment.go && 6l experiment.6 && ./6.out
rm -f experiment.6
rm -f 6.out

