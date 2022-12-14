//
// SMBIOSKeygen
//
// Copyright (c) 2022 Pedro Vilaça
//
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// 1. Redistributions of source code must retain the above copyright notice, this
// list of conditions and the following disclaimer.
//
// 2. Redistributions in binary form must reproduce the above copyright notice,
// this list of conditions and the following disclaimer in the documentation and/
// or other materials provided with the distribution.
//
// 3. Neither the name of the copyright holder nor the names of its contributors
// may be used to endorse or promote products derived from this software without
// specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
// ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
// WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
//

package main

import (
	"testing"
)

// serial numbers
/*
https://giuliomac.wordpress.com/2014/03/01/5-real-mac-serial-numbers-for-your-hackintosh/

1. iMac 21.5-inch, Late 2013
CPU: Intel Core i5 – 2,70 GHz
RAM: 8GB 1600 MHz DDR3
GPU: Intel Iris Pro 1024 MB
Serial number: C02L13ECF8J2 [Verification]

2. iMac 27-inch, Late 2013
CPU: Intel Core i5 – 3,50 GHz
RAM: 8GB 1600 MHz DDR3
GPU: NVIDIA GeForce GTX 780M 4096 MB
Serial number: C02LC1T5FLHH [Verification]

3. MacBook Pro Retina 13-inch, Late 2013
CPU: Intel Core i5 – 2,40 GHz
RAM: 8GB 1600 MHz DDR3
GPU: Intel Iris Pro 1024 MB
Serial number: C02LJ41LFH00 [Verification]

4. MacBook Pro Retina 15-inch, Late 2013
CPU: Intel Core i7 – 2,0 GHz
RAM: 8GB 1600 MHz DDR3
GPU: Intel Iris Pro 1024 MB
Serial number: C02LJ6QSFD56 [Verification]

5. MacPro, Late 2013
CPU: Intel Xenon E5 6-Core – 3,50 GHz
RAM: 16GB 1867 MHz DDR3
GPU: AMD FirePro D500 3072 MB
Serial number: F5KLV0H8F693 [Verification]

W88401231AX
C02CG123DC79

*/

type serialTests struct {
	serial string
	model  string
	year   int
}

var serials = []serialTests{
	{serial: "W88401231AX", model: "MacBook5,1", year: 2008},
	{serial: "C02CG123DC79", model: "MacBookPro6,1", year: 2010},
	{serial: "C02L13ECF8J2", model: "iMac14,1", year: 2013},
	{serial: "C02LC1T5FLHH", model: "iMac14,2", year: 2013},
	{serial: "C02LJ41LFH00", model: "MacBookPro11,1", year: 2013},
	{serial: "C02LJ6QSFD56", model: "MacBookPro11,2", year: 2013},
	{serial: "F5KLV0H8F693", model: "MacPro6,1", year: 2013},
}

func TestAlphaToValue(t *testing.T) {
	// new serial: C02VCWY4HH27 - V character
	decodedYear := alphaToValue('V', AppleTblYear, AppleYearBlacklist)
	if decodedYear != 7 {
		t.Fatal("Invalid decoded year")
	}
}

func TestBase34ToValue(t *testing.T) {
	// new serial: C02VCWY4HH27 - WY4
	// input values from serial WY4
	line := []byte{87, 89, 52}
	// the expected values
	exp := []int{2040, 1088, 4}

	mul := []int{68, 34, 1}
	for i := 0; i < len(mul); i++ {
		tmp := base34ToValue(line[i], mul[i])
		if tmp < 0 {
			t.Fatal("Invalid line symbol")
		}
		if tmp != exp[i] {
			t.Fatalf("Value not expected: %d vs %d", exp[i], tmp)
		}
	}
}

func TestGetAscii7(t *testing.T) {
	code, err := getAscii7(0x73BA1C*10, 3)
	if err != nil {
		t.Fatal("Failed to get ascii7")
	}
	if string(code) != "5NY" {
		t.Fatal("Failed to get correct ascii7")
	}
}

func TestLineToRmin(t *testing.T) {

}

func TestVerifyMLBChecksum(t *testing.T) {
	// retrieve somewhere from the internet
	if VerifyMLBChecksum("C02443500KZG2QDA7") == false {
		t.Fatal("C02443500KZG2QDA7 should be valid checksum")
	}
	// small permutation that should be invalid
	if VerifyMLBChecksum("C02443500KZG2QDA8") == true {
		t.Fatal("C02443500KZG2QDA8 should be invalid checksum")
	}
}

func TestGetProductionYear(t *testing.T) {
	// only one year available for this model
	year := getProductionYear(0, false)
	if year != 2006 {
		t.Fatal("Bad production year for index 0")
	}
	// valid are 2017 2018 2019
	year = getProductionYear(1, false)
	if year < 2017 || year > 2019 {
		t.Fatal("Bad production year for index 1")
	}
	// valid are 2006 2007
	year = getProductionYear(2, false)
	if year < 2006 || year > 2007 {
		t.Fatal("Bad production year for index 2")
	}
}

func TestGetModelCode(t *testing.T) {
	// code is always retrieving first code from the table for each model
	model := getModelCode(0, false)
	if model != "U9B" {
		t.Fatal("Bad model for index 0")
	}
	model = getModelCode(1, false)
	if model != "HH27" {
		t.Fatal("Bad model for index 1")
	}
}

func TestGetBoardCode(t *testing.T) {
	// code is always retrieving first code from the table for each model
	board := getBoardCode(0, false)
	if board != "V3G" {
		t.Fatal("Bad board for index 0")
	}
	board = getBoardCode(1, false)
	if board != "HJ9L" {
		t.Fatal("Bad board for index 1")
	}
}

func TestParseSerial(t *testing.T) {
	for _, v := range serials {
		s, err := parseSerial(v.serial)
		if err != nil {
			t.Fatalf("%s", err)
		}
		if !s.Valid {
			t.Fatal("Serial is not valid")
		}
		if s.ProductName != v.model {
			t.Fatal("Model name is wrong")
		}
		if s.DecodedYear != v.year {
			t.Fatal("Year is wrong")
		}
	}
}

func TestGenerateSerial(t *testing.T) {
	args := Params{
		Index: 0,
		Year:  -1,
		Week:  -1,
		Copy:  -1,
		Line:  -1,
	}
	// generate a serial
	s, err := generateSerial(args)
	if err != nil {
		t.Fatal(err)
	}
	// try to parse it
	s, err = parseSerial(s.String())
	if err != nil {
		t.Fatal(err)
	}
}

func TestMLB(t *testing.T) {
	s, err := parseSerial("W88401231AX")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if !s.Valid {
		t.Fatal("Serial is not valid")
	}
	// this function never fails for now
	_ = s.MLB()
}
