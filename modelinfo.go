// modelinfo.h converted to Go

package main

var AppleLegacyLocations = []string{
	"CK",
	"CY",
	"FC",
	"G8",
	"QP",
	"XA",
	"XB",
	"PT",
	"QT",
	"UV",
	"RN",
	"RM",
	"SG",
	"W8",
	"YM",
	// New
	"H0",
	"C0",
	"C3",
	"C7",
	"MB",
	"EE",
	"VM",
	"1C",
	"4H",
	"MQ",
	"WQ",
	"7J",
	"FK",
	"F1",
	"F2",
	"F7",
	"DL",
	"DM",
	"73",
}

const APPLE_LEGACY_LOCATION_OLDMAX = 15

var AppleLegacyLocationNames = []string{
	"Ireland (Cork)",
	"Korea",
	"USA (Fountain, Colorado)",
	"USA",
	"USA",
	"USA (ElkGrove/Sacramento, California)",
	"USA (ElkGrove/Sacramento, California)",
	"Korea",
	"Taiwan (Quanta Computer)",
	"Taiwan",
	"Mexico",
	"Refurbished Model",
	"Singapore",
	"China (Shanghai)",
	"China",
	// New
	"Unknown",
	"China (Quanta Computer, Tech-Com)",
	"China (Shenzhen, Foxconn)",
	"China (Shanghai, Pegatron)",
	"Malaysia",
	"Taiwan",
	"Czech Republic (Pardubice, Foxconn)",
	"China",
	"China",
	"China",
	"China",
	"China (Hon Hai/Foxconn)",
	"China (Zhengzhou, Foxconn)",
	"China (Zhengzhou, Foxconn)",
	"China (Zhengzhou, Foxconn)",
	"China",
	"China (Foxconn)",
	"China (Foxconn)",
	"Unknown",
}

var AppleLocations = []string{
	"C02",
	"C07",
	"C17",
	"C1M",
	"C2V",
	"CK2",
	"D25",
	"F5K",
	"W80",
	"W88",
	"W89",
	"CMV",
	"YM0",
	"DGK",
	"FVF",
}

var AppleLocationNames = []string{
	"China (Quanta Computer)",
	"China (Quanta Computer)",
	"China",
	"China",
	"China",
	"Ireland (Cork)",
	"Unknown",
	"USA (Flextronics)",
	"Unknown",
	"Unknown",
	"Unknown",
	"Unknown",
	"China (Hon Hai/Foxconn)",
	"Unknown",
	"Unknown",
}

var MLBBlock1 = []string{
	"200", "600", "403", "404", "405", "303", "108",
	"207", "609", "501", "306", "102", "701", "301",
	"501", "101", "300", "130", "100", "270", "310",
	"902", "104", "401", "902", "500", "700", "802",
}

var MLBBlock2 = []string{
	"GU", "4N", "J9", "QX", "OP", "CD", "GU",
}

var MLBBlock3 = []string{
	"1H", "1M", "AD", "1F", "A8", "UE", "JA", "JC", "8C", "CB", "FB",
}

// Used for line/copy:        A   B   C   D   E   F   G   H   I   J   K   L   M   N   O   P   Q   R   S   T   U   V   W   X   Y   Z
var AppleTblBase34 = []int{10, 11, 12, 13, 14, 15, 16, 17, 0, 18, 19, 20, 21, 22, 0, 23, 24, 25, 26, 27, 28, 29, 30, 31, 32, 33}
var AppleBase34Blacklist = "IO"
var AppleBase34Reverse = "0123456789ABCDEFGHJKLMNPQRSTUVWXYZ"

// Used for year:           A   B   C   D   E   F   G   H   I   J   K   L   M   N   O   P   Q   R   S   T   U   V   W   X   Y   Z
var AppleTblYear = []int{0, 0, 0, 0, 0, 1, 1, 2, 0, 2, 3, 3, 4, 4, 0, 5, 5, 6, 6, 7, 0, 7, 8, 8, 9, 9}
var AppleYearBlacklist = "ABEIOU"
var AppleYearReverse = "CDFGHJKLMNPQRSTVWXYZ"

// Used for week to year add:  A   B   C   D   E   F   G   H   I   J   K   L   M   N   O   P   Q   R   S   T   U   V   W   X   Y   Z
var AppleTblWeekAdd = []int{0, 0, 0, 26, 0, 0, 26, 0, 0, 26, 0, 26, 0, 26, 0, 0, 26, 0, 26, 0, 0, 26, 0, 26, 0, 26}

// Used for week            A   B   C   D   E   F   G   H   I   J   K   L   M   N   O   P   Q   R   S   T   U   V   W   X   Y   Z
var AppleTblWeek = []int{0, 0, 10, 11, 0, 12, 13, 14, 0, 15, 16, 17, 18, 19, 0, 20, 21, 22, 0, 23, 0, 24, 25, 26, 27, 0}
var AppleWeekBlacklist = "ABEIOSUZ"
var AppleWeekReverse = "0123456789CDFGHJKLMNPQRTVWX123456789CDFGHJKLMNPQRTVWXY"

// Apple assigned MAC address prefixes from https://gitlab.com/wireshark/wireshark/-/raw/master/manuf
var AppleRomPrefix = []string{
	"000393",
	"000A27",
	"000A95",
	"000D93",
	"0010FA",
	"001124",
	"001451",
	"0016CB",
	"0017F2",
	"0019E3",
	"001B63",
	"001CB3",
	"001D4F",
	"001E52",
	"001EC2",
	"001F5B",
	"001FF3",
	"0021E9",
	"002241",
	"002312",
	"002332",
	"00236C",
	"0023DF",
	"002436",
	"002500",
	"00254B",
	"0025BC",
	"002608",
	"00264A",
	"0026B0",
	"0026BB",
	"003065",
	"003EE1",
	"0050E4",
	"0056CD",
	"006171",
	"006D52",
	"008865",
	"00B362",
	"00C610",
	"00CDFE",
	"00F4B9",
	"00F76F",
	"040CCE",
	"041552",
	"041E64",
	"042665",
	"04489A",
	"044BED",
	"0452F3",
	"045453",
	"0469F8",
	"04D3CF",
	"04DB56",
	"04E536",
	"04F13E",
	"04F7E4",
	"086698",
	"086D41",
	"087045",
	"087402",
	"0C1539",
	"0C3021",
	"0C3E9F",
	"0C4DE9",
	"0C5101",
	"0C74C2",
	"0C771A",
	"0CBC9F",
	"0CD746",
	"101C0C",
	"1040F3",
	"10417F",
	"1093E9",
	"109ADD",
	"10DDB1",
	"14109F",
	"145A05",
	"148FC6",
	"1499E2",
	"14BD61",
	"182032",
	"183451",
	"186590",
	"189EFC",
	"18AF61",
	"18AF8F",
	"18E7F4",
	"18EE69",
	"18F643",
	"1C1AC0",
	"1C5CF2",
	"1C9148",
	"1C9E46",
	"1CABA7",
	"1CE62B",
	"203CAE",
	"20768F",
	"2078F0",
	"207D74",
	"209BCD",
	"20A2E4",
	"20AB37",
	"20C9D0",
	"241EEB",
	"24240E",
	"245BA7",
	"24A074",
	"24A2E1",
	"24AB81",
	"24E314",
	"24F094",
	"280B5C",
	"283737",
	"285AEB",
	"286AB8",
	"286ABA",
	"28A02B",
	"28CFDA",
	"28CFE9",
	"28E02C",
	"28E14C",
	"28E7CF",
	"28ED6A",
	"28F076",
	"2C1F23",
	"2C200B",
	"2C3361",
	"2CB43A",
	"2CBE08",
	"2CF0A2",
	"2CF0EE",
	"3010E4",
	"30636B",
	"3090AB",
	"30F7C5",
	"341298",
	"34159E",
	"34363B",
	"3451C9",
	"34A395",
	"34AB37",
	"34C059",
	"34E2FD",
	"380F4A",
	"38484C",
	"3871DE",
	"38B54D",
	"38C986",
	"38CADA",
	"3C0754",
	"3C15C2",
	"3CAB8E",
	"3CD0F8",
	"3CE072",
	"403004",
	"40331A",
	"403CFC",
	"404D7F",
	"406C8F",
	"40A6D9",
	"40B395",
	"40D32D",
	"440010",
	"442A60",
	"444C0C",
	"44D884",
	"44FB42",
	"483B38",
	"48437C",
	"484BAA",
	"4860BC",
	"48746E",
	"48A195",
	"48BF6B",
	"48D705",
	"48E9F1",
	"4C3275",
	"4C57CA",
	"4C74BF",
	"4C7C5F",
	"4C8D79",
	"4CB199",
	"503237",
	"507A55",
	"5082D5",
	"50EAD6",
	"542696",
	"544E90",
	"54724F",
	"549F13",
	"54AE27",
	"54E43A",
	"54EAA8",
	"581FAA",
	"58404E",
	"5855CA",
	"587F57",
	"58B035",
	"5C5948",
	"5C8D4E",
	"5C95AE",
	"5C969D",
	"5C97F3",
	"5CADCF",
	"5CF5DA",
	"5CF7E6",
	"5CF938",
	"600308",
	"60334B",
	"606944",
	"609217",
	"609AC1",
	"60A37D",
	"60C547",
	"60D9C7",
	"60F445",
	"60F81D",
	"60FACD",
	"60FB42",
	"60FEC5",
	"64200C",
	"6476BA",
	"649ABE",
	"64A3CB",
	"64A5C3",
	"64B0A6",
	"64B9E8",
	"64E682",
	"680927",
	"685B35",
	"68644B",
	"68967B",
	"689C70",
	"68A86D",
	"68AE20",
	"68D93C",
	"68DBCA",
	"68FB7E",
	"6C19C0",
	"6C3E6D",
	"6C4008",
	"6C709F",
	"6C72E7",
	"6C8DC1",
	"6C94F8",
	"6CAB31",
	"6CC26B",
	"701124",
	"7014A6",
	"703EAC",
	"70480F",
	"705681",
	"70700D",
	"7073CB",
	"7081EB",
	"70A2B3",
	"70CD60",
	"70DEE2",
	"70E72C",
	"70ECE4",
	"70F087",
	"741BB2",
	"748114",
	"748D08",
	"74E1B6",
	"74E2F5",
	"7831C1",
	"783A84",
	"784F43",
	"786C1C",
	"787E61",
	"789F70",
	"78A3E4",
	"78CA39",
	"78D75F",
	"78FD94",
	"7C0191",
	"7C04D0",
	"7C11BE",
	"7C5049",
	"7C6D62",
	"7C6DF8",
	"7CC3A1",
	"7CC537",
	"7CD1C3",
	"7CF05F",
	"7CFADF",
	"80006E",
	"804971",
	"80929F",
	"80BE05",
	"80D605",
	"80E650",
	"80EA96",
	"80ED2C",
	"842999",
	"843835",
	"84788B",
	"848506",
	"8489AD",
	"848E0C",
	"84A134",
	"84B153",
	"84FCAC",
	"84FCFE",
	"881FA1",
	"885395",
	"8863DF",
	"8866A5",
	"886B6E",
	"88C663",
	"88CB87",
	"88E87F",
	"8C006D",
	"8C2937",
	"8C2DAA",
	"8C5877",
	"8C7B9D",
	"8C7C92",
	"8C8EF2",
	"8C8FE9",
	"8CFABA",
	"9027E4",
	"903C92",
	"9060F1",
	"907240",
	"90840D",
	"908D6C",
	"90B0ED",
	"90B21F",
	"90B931",
	"90C1C6",
	"90FD61",
	"949426",
	"94E96A",
	"94F6A3",
	"9801A7",
	"9803D8",
	"9810E8",
	"985AEB",
	"989E63",
	"98B8E3",
	"98D6BB",
	"98E0D9",
	"98F0AB",
	"98FE94",
	"9C04EB",
	"9C207B",
	"9C293F",
	"9C35EB",
	"9C4FDA",
	"9C84BF",
	"9C8BA0",
	"9CF387",
	"9CF48E",
	"9CFC01",
	"A01828",
	"A03BE3",
	"A0999B",
	"A0D795",
	"A0EDCD",
	"A43135",
	"A45E60",
	"A46706",
	"A4B197",
	"A4B805",
	"A4C361",
	"A4D18C",
	"A4D1D2",
	"A4F1E8",
	"A82066",
	"A85B78",
	"A860B6",
	"A8667F",
	"A886DD",
	"A88808",
	"A88E24",
	"A8968A",
	"A8BBCF",
	"A8FAD8",
	"AC293A",
	"AC3C0B",
	"AC61EA",
	"AC7F3E",
	"AC87A3",
	"ACBC32",
	"ACCF5C",
	"ACFDEC",
	"B03495",
	"B0481A",
	"B065BD",
	"B0702D",
	"B09FBA",
	"B418D1",
	"B44BD2",
	"B48B19",
	"B49CDF",
	"B4F0AB",
	"B8098A",
	"B817C2",
	"B844D9",
	"B853AC",
	"B8782E",
	"B88D12",
	"B8C75D",
	"B8E856",
	"B8F6B1",
	"B8FF61",
	"BC3BAF",
	"BC4CC4",
	"BC52B7",
	"BC5436",
	"BC6778",
	"BC6C21",
	"BC926B",
	"BC9FEF",
	"BCA920",
	"BCEC5D",
	"C01ADA",
	"C06394",
	"C0847A",
	"C09F42",
	"C0CCF8",
	"C0CECD",
	"C0D012",
	"C0F2FB",
	"C42C03",
	"C4B301",
	"C81EE7",
	"C82A14",
	"C8334B",
	"C869CD",
	"C86F1D",
	"C88550",
	"C8B5B7",
	"C8BCC8",
	"C8E0EB",
	"C8F650",
	"CC088D",
	"CC08E0",
	"CC20E8",
	"CC25EF",
	"CC29F5",
	"CC4463",
	"CC785F",
	"CCC760",
	"D0034B",
	"D023DB",
	"D02598",
	"D03311",
	"D04F7E",
	"D0A637",
	"D0C5F3",
	"D0E140",
	"D4619D",
	"D49A20",
	"D4DCCD",
	"D4F46F",
	"D8004D",
	"D81D72",
	"D83062",
	"D89695",
	"D89E3F",
	"D8A25E",
	"D8BB2C",
	"D8CF9C",
	"D8D1CB",
	"DC0C5C",
	"DC2B2A",
	"DC2B61",
	"DC3714",
	"DC415F",
	"DC86D8",
	"DC9B9C",
	"DCA4CA",
	"DCA904",
	"E05F45",
	"E06678",
	"E0ACCB",
	"E0B52D",
	"E0B9BA",
	"E0C767",
	"E0C97A",
	"E0F5C6",
	"E0F847",
	"E425E7",
	"E48B7F",
	"E498D6",
	"E49A79",
	"E4C63D",
	"E4CE8F",
	"E4E4AB",
	"E8040B",
	"E80688",
	"E8802E",
	"E88D28",
	"E8B2AC",
	"EC3586",
	"EC852F",
	"ECADB8",
	"F02475",
	"F07960",
	"F099BF",
	"F0B0E7",
	"F0B479",
	"F0C1F1",
	"F0CBA1",
	"F0D1A9",
	"F0DBE2",
	"F0DBF8",
	"F0DCE2",
	"F0F61C",
	"F40F24",
	"F41BA1",
	"F431C3",
	"F437B7",
	"F45C89",
	"F4F15A",
	"F4F951",
	"F80377",
	"F81EDF",
	"F82793",
	"F86214",
	"FC253F",
	"FCD848",
	"FCE998",
	"FCFC48",
}
