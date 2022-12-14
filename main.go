//
// SMBIOSKeygen
//
// Ported to Go from https://github.com/acidanthera/OpenCorePkg/tree/master/Utilities/macserial
//
// Original C version
// Copyright (c) 2018-2020 vit9696
// Copyright (c) 2020 Matis Schotte
//
// Go version
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
	crand "crypto/rand"
	"encoding/binary"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
	PROGRAM_VERSION = "2.1.8"
	SERIAL_WEEK_MIN = 1
	SERIAL_WEEK_MAX = 53
	SERIAL_YEAR_MIN = 2000
	SERIAL_YEAR_MAX = 2030

	SERIAL_YEAR_OLD_MIN  = 2003
	SERIAL_YEAR_OLD_MAX  = 2012
	SERIAL_YEAR_NEW_MIN  = 2010
	SERIAL_YEAR_NEW_MID  = 2020
	SERIAL_YEAR_NEW_MAX  = 2030
	SERIAL_COPY_MIN      = 1
	SERIAL_COPY_MAX      = 34
	SERIAL_LINE_MIN      = 0
	SERIAL_LINE_REPR_MAX = 1155
	SERIAL_LINE_MAX      = 3399 /* 68*33 + 33*34 + 33 */
	SERIAL_OLD_LEN       = 11
	SERIAL_NEW_LEN       = 12
	MODEL_CODE_OLD_LEN   = 3
	MODEL_CODE_NEW_LEN   = 4
	COUNTRY_OLD_LEN      = 2
	COUNTRY_NEW_LEN      = 3
	MLB_MAX_SIZE         = 32
)

const (
	MODE_SYSTEM_INFO = iota
	MODE_SERIAL_INFO
	MODE_MLB_INFO
	MODE_LIST_MODELS
	MODE_LIST_PRODUCTS
	MODE_GENERATE_MLB
	MODE_GENERATE_CURRENT
	MODE_GENERATE_ALL
	MODE_GENERATE_DERIVATIVES
)

type AppleModel uint32

type PlatformData struct {
	productName  string
	serialNumber string
}

type AppleModelDescription struct {
	code string
	name string
}

type Serial struct {
	// these are the items that compose the serial
	Country string  // 2 or 3 digits
	Year    [1]byte // 1 digit
	Week    [2]byte // 1 or 2 digits
	Line    [3]byte // 3 digits
	Model   string  // 3 or 4 digit model code

	// other available items
	CountryDesc string // the production location description
	ProductName string
	ModelDesc   string // complete model name string
	WeekStart   string
	WeekEnd     string
	// data
	DecodedYear int
	DecodedWeek int
	DecodedLine int
	DecodedCopy int
	Valid       bool
	Legacy      bool // true if legacy serial number
	// internal data
	index        int // the model index
	countryIndex int
}

// all the possible tunning parameters
// Index or ModelCode are mandatory but mutally exclusive
type Params struct {
	Index     int    // model index (can be found using -l option)
	Year      int    // production year
	Week      int    // production week
	Country   string // country code
	ModelCode string // the 3 or 4 digit model code (can be found using -l option)
	Line      int    // the production line
	Copy      int    //
}

var rnd *rand.Rand

// https://programming.guide/go/crypto-rand-int.html
type cryptoSource struct{}

// implement the math/rand.Source interface
func (s cryptoSource) Seed(seed int64) {}
func (s cryptoSource) Int63() int64    { return int64(s.Uint64() & ^uint64(1<<63)) }
func (s cryptoSource) Uint64() (v uint64) {
	err := binary.Read(crand.Reader, binary.BigEndian, &v)
	if err != nil {
		panic(err)
	}
	return v
}

// initialize the random generator with seed from the secure rng
func init() {
	var src cryptoSource
	rnd = rand.New(src)
}

func pseudoRandom() int {
	return rnd.Int()
}

// pseudoRandomBetween returns a random number between the half-open interval
// make sure that b > a otherwise it panics
func pseudoRandomBetween(a, b uint32) int {
	// no b<a check, it panics if n <= 0 ;-)
	return rnd.Intn(int(b-a)) + int(a)
}

// generateROM generates a MAC address based on Apple prefixes
func generateROM() string {
	prefix := AppleRomPrefix[pseudoRandom()%len(AppleRomPrefix)]
	mac := fmt.Sprintf("%s%02X%02X%02X", prefix, pseudoRandomBetween(0, 256), pseudoRandomBetween(0, 256), pseudoRandomBetween(0, 256))
	return mac
}

// Apple uses various conversion tables (e.g. AppleBase34) for value encoding.
func alphaToValue(c byte, conv []int, blacklist string) int {
	if c < 'A' || c > 'Z' {
		return -1
	}

	for i := 0; i < len(blacklist); i++ {
		if blacklist[i] == c {
			return -1
		}
	}

	return conv[c-'A']
}

// This is modified base34 used by Apple with I and O excluded.
func base34ToValue(c byte, mul int) int {
	if c >= '0' && c <= '9' {
		return int((c - '0')) * mul
	}
	if c >= 'A' && c <= 'Z' {
		tmp := alphaToValue(c, AppleTblBase34, AppleBase34Blacklist)
		if tmp >= 0 {
			return tmp * mul
		}
	}
	return -1
}

func lineToRmin(line int) int {
	// info->line[0] is raw decoded copy, but it is not the real first produced unit.
	// To get the real copy we need to find the minimal allowed raw decoded copy,
	// which allows to obtain info->decodedLine.
	var rmin int
	if line > SERIAL_LINE_REPR_MAX {
		rmin = (line - SERIAL_LINE_REPR_MAX + 67) / 68
	}
	return rmin
}

// This one is modded to implement CCC algo for better generation.
// Changed base36 to base34, since that's what Apple uses.
// The algo is trash but is left for historical reasons.
func getAscii7(value uint32, size int) ([]byte, error) {
	// This is CCC conversion.
	if value < 1000000 {
		return []byte{}, fmt.Errorf("Invalid value argument")
	}

	for value > 10000000 {
		value /= 10
	}

	// log(2**64) / log(34) = 12.57 => max 13 char + '\0'
	ret := make([]byte, 14)
	offset := 13
	for {
		ret[offset] = AppleBase34Reverse[value%34]
		value /= 34
		if value == 0 {
			break
		}
		offset--
	}

	return ret[offset : offset+size], nil
}

func VerifyMLBChecksum(mlb string) bool {
	alphabet := "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"
	checksum := 0
	mlbLen := len(mlb)
	for i := 0; i < mlbLen; i++ {
		for j := 0; j < len(alphabet); j++ {
			if mlb[i] == alphabet[j] {
				// go doesn't let convert a bool to a int
				if (i & 1) == (mlbLen & 1) {
					checksum += 3 * j
				} else {
					checksum += 1 * j
				}
				break
			}
		}
	}
	return (checksum % len(alphabet)) == 0
}

func getProductionYear(model AppleModel, print bool) uint32 {
	var num uint32

	for num = 0; num < APPLE_MODEL_YEAR_MAX && AppleModelYear[model][num] > 0; num++ {
		if print {
			if num+1 != APPLE_MODEL_YEAR_MAX && AppleModelYear[model][num+1] > 0 {
				fmt.Printf("%d, ", AppleModelYear[model][num])
			} else {
				fmt.Printf("%d\n", AppleModelYear[model][num])
			}
		}
	}

	if ApplePreferredModelYear[model] > 0 {
		return ApplePreferredModelYear[model]
	}

	// XXX: improve this random? we just need a tiny number
	return AppleModelYear[model][uint32(pseudoRandom())%num]
}

func getModelCode(model AppleModel, print bool) string {
	if print {
		for i := 0; i < APPLE_MODEL_CODE_MAX && AppleModelCode[model][i] != ""; i++ {
			if i+1 != APPLE_MODEL_CODE_MAX && AppleModelCode[model][i+1] != "" {
				fmt.Printf("%s, ", AppleModelCode[model][i])
			} else {
				fmt.Printf("%s\n", AppleModelCode[model][i])
			}
		}
	}
	// Always choose the first model for stability by default.
	return AppleModelCode[model][0]
}

func getBoardCode(model AppleModel, print bool) string {
	if print {
		for i := 0; i < APPLE_BOARD_CODE_MAX && AppleBoardCode[model][i] != ""; i++ {
			if i+1 != APPLE_BOARD_CODE_MAX && AppleBoardCode[model][i+1] != "" {
				fmt.Printf("%s, ", AppleBoardCode[model][i])
			} else {
				fmt.Printf("%s\n", AppleBoardCode[model][i])
			}
		}
	}
	// Always choose the first model for stability by default.
	return AppleBoardCode[model][0]
}

// parseSerial retrieves the information about a serial number
func parseSerial(serial string) (Serial, error) {
	info := Serial{}
	// Verify length.
	serial_len := len(serial)
	if serial_len != SERIAL_OLD_LEN && serial_len != SERIAL_NEW_LEN {
		if serial_len == 17 {
			return info, fmt.Errorf("Invalid serial, you probably inserted a MLB")
		} else {
			return info, fmt.Errorf("Invalid serial length, must be %d or %d", SERIAL_NEW_LEN, SERIAL_OLD_LEN)
		}
	}

	// Assume every serial valid by default.
	info.Valid = true

	// Verify alphabet (base34 with I and O exclued).
	for i := 0; i < serial_len; i++ {
		if !((serial[i] >= 'A' && serial[i] <= 'Z' && serial[i] != 'O' && serial[i] != 'I') ||
			(serial[i] >= '0' && serial[i] <= '9')) {
			fmt.Printf("WARN: Invalid symbol '%c' in serial!\n", serial[i])
			info.Valid = false
		}
	}

	model_len := 0

	var serialModel string
	switch serial_len {
	case SERIAL_NEW_LEN:
		serialModel = serial[serial_len-MODEL_CODE_NEW_LEN:]
	case SERIAL_OLD_LEN:
		serialModel = serial[serial_len-MODEL_CODE_OLD_LEN:]
	}
	// Start with looking up the model.
	info.index = -1
	for i := 0; i < len(AppleModelCode); i++ {
		for j := 0; j < APPLE_MODEL_CODE_MAX; j++ {
			code := AppleModelCode[i][j]
			if code == "" {
				break
			}
			if code == serialModel {
				info.Model = code
				info.index = i
				break
			}
		}
	}

	// Also lookup apple model.
	for i := 0; i < len(AppleModelDesc); i++ {
		code := AppleModelDesc[i].code
		if code == serialModel {
			info.ModelDesc = AppleModelDesc[i].name
			break
		}
	}

	// Fallback to possibly valid values if model is unknown.
	// XXX: does this makes sense???
	if info.index == -1 {
		if serial_len == SERIAL_NEW_LEN {
			model_len = MODEL_CODE_NEW_LEN
		} else {
			model_len = MODEL_CODE_OLD_LEN
		}
		// XXX: test this
		info.Model = serial[serial_len-model_len : serial_len]
	}

	// Lookup production location
	info.countryIndex = -1

	if serial_len == SERIAL_NEW_LEN {
		info.Country = serial[:COUNTRY_NEW_LEN]
		// serial += COUNTRY_NEW_LEN;
		for i := 0; i < len(AppleLocations); i++ {
			if info.Country == AppleLocations[i] {
				info.countryIndex = i
				info.CountryDesc = AppleLocationNames[i]
				break
			}
		}
	} else {
		info.Legacy = true
		info.Country = serial[:COUNTRY_OLD_LEN]
		// serial += COUNTRY_OLD_LEN;
		for i := 0; i < len(AppleLegacyLocations); i++ {
			if info.Country == AppleLegacyLocations[i] {
				info.countryIndex = i
				info.CountryDesc = AppleLegacyLocationNames[i]
				break
			}
		}
	}

	// Decode production year and week
	if serial_len == SERIAL_NEW_LEN {
		// These are not exactly year and week, lower year bit is used for week encoding.
		info.Year[0] = serial[COUNTRY_NEW_LEN]
		info.Week[0] = serial[COUNTRY_NEW_LEN+1]
		// New encoding started in 2010.
		info.DecodedYear = alphaToValue(info.Year[0], AppleTblYear, AppleYearBlacklist)
		// Since year can be encoded ambiguously, check the model code for 2010/2020 difference.
		// Old check relies on first letter of model to be greater than or equal to H, which breaks compatibility with iMac20,2 (=0).
		// Added logic checks provided model years `AppleModelYear` first year greater than or equal to 2020.
		if (info.index >= 0 && AppleModelYear[info.index][0] >= 2017 && info.DecodedYear < 7) ||
			(info.DecodedYear == 0 && info.Model[0] >= 'H') {
			info.DecodedYear += 2020
		} else if info.DecodedYear >= 0 {
			info.DecodedYear += 2010
		} else {
			fmt.Printf("WARN: Invalid year symbol '%c'!\n", info.Year[0])
			info.Valid = false
		}

		if info.Week[0] > '0' && info.Week[0] <= '9' {
			info.DecodedWeek = int(info.Week[0] - '0')
		} else {
			info.DecodedWeek = alphaToValue(info.Week[0], AppleTblWeek, AppleWeekBlacklist)
		}

		if info.DecodedWeek > 0 {
			if info.DecodedYear > 0 {
				info.DecodedWeek += alphaToValue(info.Year[0], AppleTblWeekAdd, "")
			}
		} else {
			fmt.Printf("WARN: Invalid week symbol '%c'!\n", info.Week[0])
			info.Valid = false
		}
	} else {
		info.Year[0] = serial[COUNTRY_OLD_LEN]
		info.Week[0] = serial[COUNTRY_OLD_LEN+1]
		info.Week[1] = serial[COUNTRY_OLD_LEN+2]

		// This is proven by MacPro5,1 valid serials from 2011 and 2012.
		if info.Year[0] >= '0' && info.Year[0] <= '2' {
			info.DecodedYear = 2010 + int(info.Year[0]-'0')
		} else if info.Year[0] >= '3' && info.Year[0] <= '9' {
			info.DecodedYear = 2000 + int(info.Year[0]-'0')
		} else {
			info.DecodedYear = -1
			fmt.Printf("WARN: Invalid year symbol '%c'!\n", info.Year[0])
			info.Valid = false
		}

		for i := 0; i < 2; i++ {
			if info.Week[i] >= '0' && info.Week[i] <= '9' {
				if i == 0 {
					info.DecodedWeek += 10 * int(info.Week[i]-'0')
				} else {
					info.DecodedWeek += 1 * int(info.Week[i]-'0')
				}
			} else {
				info.DecodedWeek = -1
				fmt.Printf("WARN: Invalid week symbol '%c'!\n", info.Week[i])
				info.Valid = false
				break
			}
		}
	}

	if info.DecodedWeek < SERIAL_WEEK_MIN || info.DecodedWeek > SERIAL_WEEK_MAX {
		fmt.Printf("WARN: Decoded week %d is out of valid range [%d, %d]!\n", info.DecodedWeek, SERIAL_WEEK_MIN, SERIAL_WEEK_MAX)
		info.DecodedWeek = -1
	}

	if info.DecodedYear > 0 && info.index >= 0 {
		found := false
		for i := 0; !found && i < APPLE_MODEL_YEAR_MAX && AppleModelYear[info.index][i] > 0; i++ {
			if int(AppleModelYear[info.index][i]) == info.DecodedYear {
				found = true
			}
		}
		if !found {
			fmt.Printf("WARN: Invalid year %d for model %s\n", info.DecodedYear, ApplePlatformData[info.index].productName)
			info.Valid = false
		}
	}

	if info.DecodedYear > 0 && info.DecodedWeek > 0 {
		day := 1 + 7*int((info.DecodedWeek-1))
		// the month must be set to 1 for this to match original macserial
		// we also don't need to add anything to the month
		t := time.Date(int(info.DecodedYear), 1, day, 0, 0, 0, 0, time.UTC)
		info.WeekStart = fmt.Sprintf("%02d.%02d.%04d", t.Day(), t.Month(), t.Year())

		if info.DecodedWeek == 53 && t.Day() != 31 {
			info.WeekEnd = fmt.Sprintf("31.12.%04d", t.Year())
		} else if info.DecodedWeek < 53 {
			t := time.Date(int(info.DecodedYear), 1, day+6, 0, 0, 0, 0, time.UTC)
			info.WeekEnd = fmt.Sprintf("%02d.%02d.%04d", t.Day(), t.Month(), t.Year())
		}
	}

	// Decode production line and copy
	mul := []int{68, 34, 1}
	serialPos := 0
	if serial_len == SERIAL_NEW_LEN {
		serialPos = COUNTRY_NEW_LEN + 2
	} else {
		serialPos = COUNTRY_OLD_LEN + 3
	}
	for i := 0; i < len(mul); i++ {
		info.Line[i] = serial[serialPos]
		tmp := base34ToValue(info.Line[i], mul[i])
		if tmp >= 0 {
			info.DecodedLine += tmp
		} else {
			fmt.Printf("WARN: Invalid line symbol '%c'!\n", info.Line[i])
			info.Valid = false
			break
		}
		serialPos++
	}

	if info.DecodedLine >= 0 {
		info.DecodedCopy = base34ToValue(info.Line[0], 1) - lineToRmin(info.DecodedLine)
	}

	info.ProductName = ApplePlatformData[info.index].productName

	return info, nil
}

// generateSerial generates a new serial
func generateSerial(param Params) (Serial, error) {

	if param.Index < 0 && param.ModelCode == "" {
		// fmt.Printf("ERROR: Unable to determine model!\n")
		return Serial{}, fmt.Errorf("Unable to determine Mac model")
	}

	var model string
	if param.ModelCode == "" {
		model = getModelCode(AppleModel(param.Index), false)
	} else {
		// XXX: validate the 3 digit model code
		model = param.ModelCode
		// XXX: set the model index
	}

	country := param.Country
	country_len := len(param.Country)
	if country_len == 0 {
		// Random country choice strongly decreases key verification probability.
		if len(model) == MODEL_CODE_NEW_LEN {
			country_len = COUNTRY_NEW_LEN
		} else {
			country_len = COUNTRY_OLD_LEN
		}
		if param.Index < 0 {
			if country_len == COUNTRY_OLD_LEN {
				country = AppleLegacyLocations[0]
			} else {
				country = AppleLocations[0]
			}
		} else {
			// extract it from a legit serial from internal database
			country = ApplePlatformData[param.Index].serialNumber[:country_len]
		}
	}

	year := param.Year
	if param.Year < 0 {
		// XXX: this enters in conflict with ModelCode
		if param.Index < 0 {
			if country_len == COUNTRY_OLD_LEN {
				year = SERIAL_YEAR_OLD_MAX
			} else {
				year = SERIAL_YEAR_NEW_MID
			}
		} else {
			year = int(getProductionYear(AppleModel(param.Index), false))
		}
	}

	week := param.Week
	// Last week is too rare to care
	if param.Week < 0 {
		week = pseudoRandomBetween(SERIAL_WEEK_MIN, SERIAL_WEEK_MAX-1)
	}

	var yearData [1]byte
	var weekData [2]byte
	var weekString string
	var yearString string
	if country_len == COUNTRY_OLD_LEN {
		if year < SERIAL_YEAR_OLD_MIN || year > SERIAL_YEAR_OLD_MAX {
			return Serial{}, fmt.Errorf("Year %d is out of valid legacy range [%d, %d]", year, SERIAL_YEAR_OLD_MIN, SERIAL_YEAR_OLD_MAX)
		}
		yearData[0] = '0' + byte((year-2000)%10)
		weekData[0] = '0' + byte(week/10)
		weekData[1] = '0' + byte(week%10)
		weekString = fmt.Sprintf("%s", string(weekData[:]))
		yearString = fmt.Sprintf("%s", string(yearData[:]))
	} else {
		if year < SERIAL_YEAR_NEW_MIN || year > SERIAL_YEAR_NEW_MAX {
			return Serial{}, fmt.Errorf("Year %d is out of valid modern range [%d, %d]", year, SERIAL_YEAR_NEW_MIN, SERIAL_YEAR_NEW_MAX)
		}

		base_new_year := 2010
		if year >= SERIAL_YEAR_NEW_MID {
			base_new_year = 2020
		}

		if week >= 27 {
			yearData[0] = AppleYearReverse[(year-base_new_year)*2+1]
		} else {
			yearData[0] = AppleYearReverse[(year-base_new_year)*2]
		}
		weekData[0] = AppleWeekReverse[week]
		// if we print directly the bytes it will encode the nul values and they will count towards the length
		weekString = fmt.Sprintf("%s", string(weekData[:1]))
		yearString = fmt.Sprintf("%s", string(yearData[:]))
	}

	line := param.Line
	if param.Line < 0 {
		line = pseudoRandomBetween(SERIAL_LINE_MIN, SERIAL_LINE_MAX)
	}

	rmin := lineToRmin(line)

	// Verify and apply user supplied copy if any
	if param.Copy >= 0 {
		rmin += param.Copy - 1
		if rmin*68 > line {
			return Serial{}, fmt.Errorf("Copy %d cannot represent line %d", param.Copy, line)
		}
	}
	var lineData [3]byte
	lineData[0] = AppleBase34Reverse[rmin]
	lineData[1] = AppleBase34Reverse[(line-rmin*68)/34]
	lineData[2] = AppleBase34Reverse[(line-rmin*68)%34]

	// print a serial to send to the parser
	serial := fmt.Sprintf("%s%s%s%s%s", country, yearString, weekString, lineData, model)
	// parse it and return the structure
	s, err := parseSerial(serial)
	if err != nil {
		return s, err
	}
	return s, nil
}

func (s *Serial) String() string {
	var serial string
	if s.Legacy {
		serial = fmt.Sprintf("%s%s%s%s%s", s.Country, s.Year, s.Week, s.Line, s.Model)
	} else {
		serial = fmt.Sprintf("%s%s%s%s%s", s.Country, s.Year, s.Week[:1], s.Line, s.Model)
	}
	return serial
}

func (s *Serial) Print() {
	fmt.Printf("%14s: %4s - %s\n", "Country", s.Country, s.CountryDesc)
	fmt.Printf("%14s: %4s - %d\n", "Year", s.Year, s.DecodedYear)
	fmt.Printf("%14s: %4s - %d", "Week", s.Week, s.DecodedWeek)
	fmt.Printf(" (%s-%s)\n", s.WeekStart, s.WeekEnd)

	if s.DecodedCopy >= 0 {
		fmt.Printf("%14s: %4s - %d (copy %d)\n", "Line", s.Line, s.DecodedLine, s.DecodedCopy+1)
	} else {
		fmt.Printf("%14s: %4s - %d (copy %d)\n", "Line", s.Line, s.DecodedLine, -1)
	}
	if s.index >= 0 {
		fmt.Printf("%14s: %4s - %s\n", "Model", s.Model, ApplePlatformData[s.index].productName)
	} else {
		fmt.Printf("%14s: %4s - %s\n", "Model", s.Model, "Unknown")
	}
	if s.ModelDesc != "" {
		fmt.Printf("%14s: %s\n", "SystemModel", s.ModelDesc)
	} else {
		fmt.Printf("%14s: %s\n", "SystemModel", "Unknown, please report!")
	}
	if s.Valid {
		fmt.Printf("%14s: %s\n", "Valid", "Possibly")
	} else {
		fmt.Printf("%14s: %s\n", "Valid", "Unlikely")
	}
}

// generates a MLB from the serial number
func (s *Serial) MLB() string {
	// This is a direct reverse from CCC, rework it later...
	if s.index < 0 {
		fmt.Printf("WARN: Unknown model, assuming default!\n")
		s.index = APPLE_MODEL_MAX - 1
	}
	for {
		year := uint32(0)
		week := uint32(0)

		legacy := false
		if len(s.Country) == COUNTRY_OLD_LEN {
			legacy = true
		}

		if legacy {
			year = uint32(s.Year[0] - '0')
			week = uint32(s.Week[0]-'0')*10 + uint32(s.Week[1]-'0')
		} else {
			syear := s.Year[0]
			sweek := s.Week[0]

			srcyear := "CDFGHJKLMNPQRSTVWXYZ"
			dstyear := "00112233445566778899"
			for i := 0; i < len(srcyear); i++ {
				if syear == srcyear[i] {
					year = uint32(dstyear[i] - '0')
					break
				}
			}

			overrides := "DGJLNQSVXZ"
			for i := 0; i < len(overrides); i++ {
				if syear == overrides[i] {
					week = 27
					break
				}
			}

			srcweek := "123456789CDFGHJKLMNPQRSTVWXYZ"
			for i := 0; i < len(srcweek); i++ {
				if sweek == srcweek[i] {
					week += uint32(i) + 1
					break
				}
			}
			// This is silently not handled, and it should not be needed for normal serials.
			// Bugged MacBookPro6,2 and MacBookPro7,1 will gladly hit it.
			if week < SERIAL_WEEK_MIN {
				return fmt.Sprintf("FAIL-ZERO-%c", sweek)
			}
		}

		week--

		if week <= 9 {
			if week == 0 {
				week = SERIAL_WEEK_MAX
				if year == 0 {
					year = 9
				} else {
					year--
				}
			}
		}

		var serial string
		if legacy {
			// The loop is not present in CCC, but it throws an exception here,
			// and effectively generates nothing. The logic is crazy :/.
			// Also, it was likely meant to be written as pseudoRandom() % 0x8000.
			var code []byte
			var err error
			for {
				code, err = getAscii7(uint32(pseudoRandomBetween(0, 0x7FFE))*0x73BA1C, 3)
				if err == nil {
					break
				}
			}
			board := getBoardCode(AppleModel(s.index), false)
			suffix := AppleBase34Reverse[pseudoRandom()%34]
			// For old MLB, this is a variant of base 34 value. First item character is always 0.
			serial = fmt.Sprintf("%s%d%02d0%s%s%c", s.Country, year, week, string(code), board, suffix)
		} else {
			part1 := MLBBlock1[pseudoRandom()%len(MLBBlock1)]
			part2 := MLBBlock2[pseudoRandom()%len(MLBBlock2)]
			board := getBoardCode(AppleModel(s.index), false)
			part3 := MLBBlock3[pseudoRandom()%len(MLBBlock3)]
			serial = fmt.Sprintf("%s%d%02d%s%s%s%s", s.Country, year, week, part1, part2, board, part3)
		}
		// there is no other exit other than a valid serial, usually this function is called
		// after serial has been validated
		if VerifyMLBChecksum(serial) {
			return serial
		}
	}
}

func usage(app string) {
	fmt.Printf(
		"  ___ __  __ ___ ___ ___  ___ _  __                       \n"+
			" / __|  \\/  | _ )_ _/ _ \\/ __| |/ /___ _  _ __ _ ___ _ _  \n"+
			" \\__ \\ |\\/| | _ \\| | (_) \\__ \\ ' </ -_) || / _` / -_)    \\ \n"+
			" |___/_|  |_|___/___\\___/|___/_|\\_\\___|\\_, \\__, \\___|_||_|\n"+
			"                                       |__/|___/          \n"+
			"Usage:\n"+
			"%s command [options]\n\n"+
			"Commands:\n"+
			" --help           (-h)  show this help\n"+
			" --version        (-v)  show program version\n"+
			" --keygen         (-k)  generate necessary OpenCore serials\n"+
			" --deriv <serial> (-d)  generate all derivative serials\n"+
			" --generate       (-g)  generate serial (requires at least model option)\n"+
			" --generate-all   (-a)  generate serial for all models\n"+
			" --info <serial>  (-i)  decode serial information\n"+
			" --verify <mlb>         verify MLB checksum\n"+
			" --list           (-l)  list known mac models\n"+
			" --list-products  (-lp) list known product codes\n"+
			" --mlb <serial>         generate MLB based on serial\n"+
			" --sys            (-s)  get system info\n"+
			" --uuid           (-u)  generate UUID\n\n"+
			"Options:\n"+
			" --model <model>  (-m)  mac model (index or string) used for generation\n"+
			" --num <num>      (-n)  number of generated pairs\n"+
			" --year <year>    (-y)  year used for generation\n"+
			" --week <week>    (-w)  week used for generation\n"+
			" --country <loc>  (-c)  country location used for generation\n"+
			" --copy <copy>    (-o)  production copy index\n"+
			" --line <line>    (-e)  production line\n"+
			" --platform <ppp> (-p)  3 or 4 digit string model code used for generation\n\n", app)
}

func main() {
	var cmdHelp bool
	var cmdVersion bool
	var cmdGenerate bool
	var cmdGenerateAll bool
	var cmdInfo string
	var cmdVerify string
	var cmdList bool
	var cmdListProds bool
	var cmdMLB string
	var cmdDeriv string
	var cmdSys bool
	var cmdUuid bool
	var cmdKeygen bool
	var optModel string
	var optNum int
	var optYear int
	var optWeek int
	var optCountry string
	var optModelCode string
	var optCopy int
	var optLine int
	// https://www.antoniojgutierrez.com/posts/2021-05-14-short-and-long-options-in-go-flags-pkg/
	flag.BoolVar(&cmdHelp, "h", false, "show this help")
	flag.BoolVar(&cmdHelp, "help", false, "show this help")
	flag.BoolVar(&cmdVersion, "v", false, "")
	flag.BoolVar(&cmdVersion, "version", false, "")
	flag.BoolVar(&cmdGenerate, "g", false, "")
	flag.BoolVar(&cmdGenerate, "generate", false, "")
	flag.BoolVar(&cmdGenerateAll, "a", false, "")
	flag.BoolVar(&cmdGenerateAll, "generate-all", false, "")
	flag.StringVar(&cmdInfo, "i", "", "")
	flag.StringVar(&cmdInfo, "info", "", "")
	flag.StringVar(&cmdVerify, "verify", "", "")
	flag.BoolVar(&cmdList, "l", false, "")
	flag.BoolVar(&cmdList, "list", false, "")
	flag.BoolVar(&cmdListProds, "lp", false, "")
	flag.BoolVar(&cmdListProds, "list-products", false, "")
	flag.StringVar(&cmdMLB, "mlb", "", "")
	flag.StringVar(&cmdDeriv, "d", "", "")
	flag.StringVar(&cmdDeriv, "deriv", "", "")
	flag.BoolVar(&cmdSys, "s", false, "")
	flag.BoolVar(&cmdSys, "sys", false, "")
	flag.BoolVar(&cmdUuid, "u", false, "")
	flag.BoolVar(&cmdUuid, "uuid", false, "")
	flag.BoolVar(&cmdKeygen, "k", false, "")
	flag.BoolVar(&cmdKeygen, "keygen", false, "")
	flag.StringVar(&optModel, "m", "", "")
	flag.StringVar(&optModel, "model", "", "")
	flag.IntVar(&optNum, "n", 5, "")
	flag.IntVar(&optNum, "num", 5, "")
	// XXX: default value
	flag.IntVar(&optYear, "y", -1, "")
	flag.IntVar(&optYear, "year", -1, "")
	flag.IntVar(&optWeek, "w", -1, "")
	flag.IntVar(&optWeek, "week", -1, "")
	flag.StringVar(&optCountry, "c", "", "")
	flag.StringVar(&optCountry, "country", "", "")
	flag.StringVar(&optModelCode, "p", "", "")
	flag.StringVar(&optModelCode, "platform", "", "")
	flag.IntVar(&optCopy, "o", -1, "")
	flag.IntVar(&optCopy, "copy", -1, "")
	flag.IntVar(&optLine, "e", -1, "")
	flag.IntVar(&optLine, "line", -1, "")
	// set the usage because of duplicate commands
	flag.Usage = func() { usage(os.Args[0]) }
	flag.Parse()

	// commands that don't depend on options
	if cmdHelp {
		flag.Usage()
		os.Exit(0)
	} else if cmdVersion {
		fmt.Printf("SMBIOSKeygen v%s\n", PROGRAM_VERSION)
		os.Exit(0)
	} else if cmdSys {
		GetSystemInfo()
		os.Exit(0)
	} else if cmdUuid {
		uuid := uuid.New()
		fmt.Println(strings.ToUpper(uuid.String()))
		os.Exit(0)
	}

	// this is the most used model
	defaultIndex := 0
	for i := 0; i < APPLE_MODEL_MAX; i++ {
		if "iMacPro1,1" == ApplePlatformData[i].productName {
			defaultIndex = i
			break
		}
	}

	args := Params{
		Index: defaultIndex,
		Year:  -1,
		Week:  -1,
		Copy:  -1,
		Line:  -1,
	}

	// parse options

	// this can be a number (index) or model description
	if optModel != "" {
		value, err := strconv.Atoi(optModel)
		// error means the user inserted the model string instead
		if err != nil {
			for i := 0; i < APPLE_MODEL_MAX; i++ {
				if optModel == ApplePlatformData[i].productName {
					args.Index = i
					break
				}
			}
		} else {
			args.Index = value
		}
	}

	if optYear != -1 {
		if optYear < SERIAL_YEAR_MIN || optYear > SERIAL_YEAR_MAX {
			fmt.Printf("ERROR: Year %d is out of valid range [%d, %d]!\n", optYear, SERIAL_YEAR_MIN, SERIAL_YEAR_MAX)
			os.Exit(1)
		}
		args.Year = optYear
	}

	// seems buggy with week 2 for example
	if optWeek != -1 {
		if optWeek < SERIAL_WEEK_MIN || optWeek > SERIAL_WEEK_MAX {
			fmt.Printf("ERROR: Week %d is out of valid range [%d, %d]!\n", optWeek, SERIAL_WEEK_MIN, SERIAL_WEEK_MAX)
			os.Exit(1)
		}
		args.Week = optWeek
	}

	// XXX: it will accept any country - shouldn't we have a lookup table?
	if optCountry != "" {
		len := len(optCountry)
		if len != COUNTRY_OLD_LEN && len != COUNTRY_NEW_LEN {
			fmt.Printf("ERROR: Country location %s is neither %d nor %d symbols long!\n", optCountry, COUNTRY_OLD_LEN, COUNTRY_NEW_LEN)
			os.Exit(1)
		}
		args.Country = optCountry
	}

	// XXX: test this
	if optModelCode != "" {
		len := len(optModelCode)
		if len != MODEL_CODE_OLD_LEN && len != MODEL_CODE_NEW_LEN {
			fmt.Printf("ERROR: Platform code %s is neither %d nor %d symbols long!\n", optModelCode, MODEL_CODE_OLD_LEN, MODEL_CODE_NEW_LEN)
			os.Exit(1)
		}
		args.ModelCode = optModelCode
	}

	if optCopy != -1 {
		if optCopy < SERIAL_COPY_MIN || optCopy > SERIAL_COPY_MAX {
			fmt.Printf("ERROR: Copy %d is out of valid range [%d, %d]!\n", optCopy, SERIAL_COPY_MIN, SERIAL_COPY_MAX)
			os.Exit(1)
		}
		args.Copy = optCopy
	}

	if optLine != -1 {
		if optLine < SERIAL_LINE_MIN || optLine > SERIAL_LINE_MAX {
			fmt.Printf("ERROR: Line %d is out of valid range [%d, %d]!\n", optLine, SERIAL_LINE_MIN, SERIAL_LINE_MAX)
			os.Exit(1)
		}
		args.Line = optLine
	}

	if args.Index >= 0 && args.ModelCode != "" {
		fmt.Printf("ERROR: --model and --platform options are mutually exclusive. Please set only one.\n")
		os.Exit(1)
	}

	// and now execute the commands
	// -l  || --list
	if cmdList {
		fmt.Printf("Available models:\n")
		for j := 0; j < APPLE_MODEL_MAX; j++ {
			fmt.Printf("%14s: %s\n", "Model", ApplePlatformData[j].productName)
			fmt.Printf("%14s: %d\n", "Model Index", j)
			fmt.Printf("%14s: ", "Prod years")
			getProductionYear(AppleModel(j), true)
			fmt.Printf("%14s: %s\n", "Base Serial", ApplePlatformData[j].serialNumber)
			fmt.Printf("%14s: ", "Model codes")
			getModelCode(AppleModel(j), true)
			fmt.Printf("%14s: ", "Board codes")
			getBoardCode(AppleModel(j), true)
			fmt.Println("")
		}
		fmt.Printf("Available legacy location codes:\n")
		for j := 0; j < len(AppleLegacyLocations); j++ {
			fmt.Printf(" - %s, %s\n", AppleLegacyLocations[j], AppleLegacyLocationNames[j])
		}
		fmt.Printf("\nAvailable new location codes:\n")
		for j := 0; j < len(AppleLocations); j++ {
			fmt.Printf(" - %s, %s\n", AppleLocations[j], AppleLocationNames[j])
		}
		os.Exit(0)
	}
	// -lp || --list-products
	if cmdListProds {
		for j := 0; j < len(AppleModelDesc); j++ {
			fmt.Printf("%4s - %s\n", AppleModelDesc[j].code, AppleModelDesc[j].name)
		}
		os.Exit(0)
	}
	// -i || --info
	if cmdInfo != "" {
		s, err := parseSerial(cmdInfo)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		s.Print()
		os.Exit(0)
	}
	// --verify
	if cmdVerify != "" {
		slen := len(cmdVerify)
		switch slen {
		case 13:
			fmt.Printf("Valid MLB length: legacy\n")
		case 17:
			fmt.Printf("Valid MLB length: modern\n")
		default:
			fmt.Printf("ERROR: Invalid MLB length: %d\n", slen)
			os.Exit(1)
		}
		if VerifyMLBChecksum(cmdVerify) {
			fmt.Printf("Valid MLB checksum\n")
		} else {
			fmt.Printf("WARNING: Invalid MLB checksum\n")
		}
		os.Exit(0)
	}
	// -g || --generate
	if cmdGenerate {
		if args.Index == -1 && args.ModelCode == "" {
			fmt.Printf("ERROR: Please set at least a model or platform option\n")
			flag.Usage()
			os.Exit(1)
		}
		for i := 0; i < optNum; i++ {
			s, err := generateSerial(args)
			if err != nil {
				fmt.Printf("ERROR: %s\n", err)
				continue
			}
			mlb := s.MLB()
			fmt.Printf("%s | Serial: %s | MLB: %s\n", s.ProductName, s.String(), mlb)
		}
		os.Exit(0)
	}
	// -a || --generate-all
	if cmdGenerateAll {
		for i := 0; i < APPLE_MODEL_MAX; i++ {
			args.Index = i
			for j := 0; j < optNum; j++ {
				s, err := generateSerial(args)
				if err != nil {
					fmt.Printf("ERROR: %s\n", err)
					continue
				}
				mlb := s.MLB()
				fmt.Printf("%14s | %s | %s\n", ApplePlatformData[i].productName, s.String(), mlb)
			}
		}
		os.Exit(0)
	}
	// --mlb
	if cmdMLB != "" {
		s, err := parseSerial(cmdMLB)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		if !s.Valid {
			fmt.Printf("ERROR: Serial is not valid\n")
			os.Exit(1)
		}
		mlb := s.MLB()
		fmt.Printf("%s\n", mlb)
		os.Exit(0)
	}
	// -d || --deriv
	if cmdDeriv != "" {
		s, err := parseSerial(cmdDeriv)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}
		rmin := lineToRmin(s.DecodedLine)
		for k := 0; k < 34; k++ {
			start := k * 68
			if s.DecodedLine > start && s.DecodedLine-start <= SERIAL_LINE_REPR_MAX {
				rem := s.DecodedLine - start
				fmt.Printf("%s%s%s%c%c%c%s - copy %d\n", s.Country, s.Year, s.Week, AppleBase34Reverse[k],
					AppleBase34Reverse[rem/34], AppleBase34Reverse[rem%34], s.Model, k-rmin+1)
			}
		}
		os.Exit(0)
	}

	if cmdKeygen {
		if args.Index == -1 && args.ModelCode == "" {
			fmt.Printf("ERROR: Please set at least a model or platform option\n")
			flag.Usage()
			os.Exit(1)
		}
		s, err := generateSerial(args)
		if err != nil {
			fmt.Printf("ERROR: %s\n", err)
			os.Exit(1)
		}

		mlb := s.MLB()
		uuid := uuid.New()
		rom := generateROM()

		fmt.Printf("Type:         %s\n", s.ProductName)
		fmt.Printf("Serial:       %s\n", s.String())
		fmt.Printf("Board Serial: %s\n", mlb)
		fmt.Printf("UUID:         %s\n", strings.ToUpper(uuid.String()))
		fmt.Printf("ROM:          %s\n", rom)
		// fmt.Printf("\nYou can verify serial validity at https://checkcoverage.apple.com/\n")
		// fmt.Printf("You should be looking for a \"We're sorry, we're unable to check coverage for this serial number.\" error message.\n")
		os.Exit(0)
	}

	// "unreachable" aka no commands issued
	flag.Usage()
	os.Exit(0)
}
