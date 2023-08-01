#!/usr/bin/perl
use warnings;
use strict;

my $line;
my $prev;
my $label;

$line="";
$prev="";

while(<>){
	$line=$_;

	if ($prev =~ /^;;/) {
		$label=$prev;
		$label =~ s/^\s+|\s+$//g;
		
		if ($line =~ /^L([0-9A-F]+):/) {
			print "0x$1: \"$label\",\n"
		}
	}

	$prev=$line;
}
