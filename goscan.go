package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/color"
)

/* Command to generate top-x ports from nmap-services file:
   sort -r -k3 /usr/share/nmap/nmap-services | grep tcp | head -n 100 | cut -f 2 | cut -d '/' -f 1 | sed -n 'H;${x;s/\n/,/g;s/^,//;p;}'

   ┌──(root💀kali)-[/usr/share/nmap]
   └─# sort -r -k3 /usr/share/nmap/nmap-services | grep tcp | head -n 10 | cut -f 2 | cut -d '/' -f 1 | sed -n 'H;${x;s/\n/,/g;s/^,//;p;}'
   80,23,443,21,22,25,3389,110,445,139
*/

var top10ports = [10]int{80, 23, 443, 21, 22, 25, 3389, 110, 445, 139}
var top100ports = [100]int{80, 23, 443, 21, 22, 25, 3389, 110, 445, 139, 143, 53, 135, 3306, 8080, 1723, 111, 995, 993, 5900, 1025, 587, 8888, 199, 1720, 465, 548, 113, 81, 6001, 10000, 514, 5060, 179, 1026, 2000, 8443, 8000, 32768, 554, 26, 1433, 49152, 2001, 515, 8008, 49154, 1027, 5666, 646, 5000, 5631, 631, 49153, 8081, 2049, 88, 79, 5800, 106, 2121, 1110, 49155, 6000, 513, 990, 5357, 427, 49156, 543, 544, 5101, 144, 7, 389, 8009, 3128, 444, 9999, 5009, 7070, 5190, 3000, 5432, 1900, 3986, 13, 1029, 9, 5051, 6646, 49157, 1028, 873, 1755, 2717, 4899, 9100, 119, 37}
var top1kports = [1000]int{1, 3, 4, 6, 7, 9, 13, 17, 19, 20, 21, 22, 23, 24, 25, 26, 30, 32, 33, 37, 42, 43, 49, 53, 70, 79, 80, 81, 82, 83, 84, 85, 88, 89, 90, 99, 100, 106, 109, 110, 111, 113, 119, 125, 135, 139, 143, 144, 146, 161, 163, 179, 199, 211, 212, 222, 254, 255, 256, 259, 264, 280, 301, 306, 311, 340, 366, 389, 406, 407, 416, 417, 425, 427, 443, 444, 445, 458, 464, 465, 481, 497, 500, 512, 513, 514, 515, 524, 541, 543, 544, 545, 548, 554, 555, 563, 587, 593, 616, 617, 625, 631, 636, 646, 648, 666, 667, 668, 683, 687, 691, 700, 705, 711, 714, 720, 722, 726, 749, 765, 777, 783, 787, 800, 801, 808, 843, 873, 880, 888, 898, 900, 901, 902, 903, 911, 912, 981, 987, 990, 992, 993, 995, 999, 1000, 1001, 1002, 1007, 1009, 1010, 1011, 1021, 1022, 1023, 1024, 1025, 1026, 1027, 1028, 1029, 1030, 1031, 1032, 1033, 1034, 1035, 1036, 1037, 1038, 1039, 1040, 1041, 1042, 1043, 1044, 1045, 1046, 1047, 1048, 1049, 1050, 1051, 1052, 1053, 1054, 1055, 1056, 1057, 1058, 1059, 1060, 1061, 1062, 1063, 1064, 1065, 1066, 1067, 1068, 1069, 1070, 1071, 1072, 1073, 1074, 1075, 1076, 1077, 1078, 1079, 1080, 1081, 1082, 1083, 1084, 1085, 1086, 1087, 1088, 1089, 1090, 1091, 1092, 1093, 1094, 1095, 1096, 1097, 1098, 1099, 1100, 1102, 1104, 1105, 1106, 1107, 1108, 1110, 1111, 1112, 1113, 1114, 1117, 1119, 1121, 1122, 1123, 1124, 1126, 1130, 1131, 1132, 1137, 1138, 1141, 1145, 1147, 1148, 1149, 1151, 1152, 1154, 1163, 1164, 1165, 1166, 1169, 1174, 1175, 1183, 1185, 1186, 1187, 1192, 1198, 1199, 1201, 1213, 1216, 1217, 1218, 1233, 1234, 1236, 1244, 1247, 1248, 1259, 1271, 1272, 1277, 1287, 1296, 1300, 1301, 1309, 1310, 1311, 1322, 1328, 1334, 1352, 1417, 1433, 1434, 1443, 1455, 1461, 1494, 1500, 1501, 1503, 1521, 1524, 1533, 1556, 1580, 1583, 1594, 1600, 1641, 1658, 1666, 1687, 1688, 1700, 1717, 1718, 1719, 1720, 1721, 1723, 1755, 1761, 1782, 1783, 1801, 1805, 1812, 1839, 1840, 1862, 1863, 1864, 1875, 1900, 1914, 1935, 1947, 1971, 1972, 1974, 1984, 1998, 1999, 2000, 2001, 2002, 2003, 2004, 2005, 2006, 2007, 2008, 2009, 2010, 2013, 2020, 2021, 2022, 2030, 2033, 2034, 2035, 2038, 2040, 2041, 2042, 2043, 2045, 2046, 2047, 2048, 2049, 2065, 2068, 2099, 2100, 2103, 2105, 2106, 2107, 2111, 2119, 2121, 2126, 2135, 2144, 2160, 2161, 2170, 2179, 2190, 2191, 2196, 2200, 2222, 2251, 2260, 2288, 2301, 2323, 2366, 2381, 2382, 2383, 2393, 2394, 2399, 2401, 2492, 2500, 2522, 2525, 2557, 2601, 2602, 2604, 2605, 2607, 2608, 2638, 2701, 2702, 2710, 2717, 2718, 2725, 2800, 2809, 2811, 2869, 2875, 2909, 2910, 2920, 2967, 2968, 2998, 3000, 3001, 3003, 3005, 3006, 3007, 3011, 3013, 3017, 3030, 3031, 3052, 3071, 3077, 3128, 3168, 3211, 3221, 3260, 3261, 3268, 3269, 3283, 3300, 3301, 3306, 3322, 3323, 3324, 3325, 3333, 3351, 3367, 3369, 3370, 3371, 3372, 3389, 3390, 3404, 3476, 3493, 3517, 3527, 3546, 3551, 3580, 3659, 3689, 3690, 3703, 3737, 3766, 3784, 3800, 3801, 3809, 3814, 3826, 3827, 3828, 3851, 3869, 3871, 3878, 3880, 3889, 3905, 3914, 3918, 3920, 3945, 3971, 3986, 3995, 3998, 4000, 4001, 4002, 4003, 4004, 4005, 4006, 4045, 4111, 4125, 4126, 4129, 4224, 4242, 4279, 4321, 4343, 4443, 4444, 4445, 4446, 4449, 4550, 4567, 4662, 4848, 4899, 4900, 4998, 5000, 5001, 5002, 5003, 5004, 5009, 5030, 5033, 5050, 5051, 5054, 5060, 5061, 5080, 5087, 5100, 5101, 5102, 5120, 5190, 5200, 5214, 5221, 5222, 5225, 5226, 5269, 5280, 5298, 5357, 5405, 5414, 5431, 5432, 5440, 5500, 5510, 5544, 5550, 5555, 5560, 5566, 5631, 5633, 5666, 5678, 5679, 5718, 5730, 5800, 5801, 5802, 5810, 5811, 5815, 5822, 5825, 5850, 5859, 5862, 5877, 5900, 5901, 5902, 5903, 5904, 5906, 5907, 5910, 5911, 5915, 5922, 5925, 5950, 5952, 5959, 5960, 5961, 5962, 5963, 5987, 5988, 5989, 5998, 5999, 6000, 6001, 6002, 6003, 6004, 6005, 6006, 6007, 6009, 6025, 6059, 6100, 6101, 6106, 6112, 6123, 6129, 6156, 6346, 6389, 6502, 6510, 6543, 6547, 6565, 6566, 6567, 6580, 6646, 6666, 6667, 6668, 6669, 6689, 6692, 6699, 6779, 6788, 6789, 6792, 6839, 6881, 6901, 6969, 7000, 7001, 7002, 7004, 7007, 7019, 7025, 7070, 7100, 7103, 7106, 7200, 7201, 7402, 7435, 7443, 7496, 7512, 7625, 7627, 7676, 7741, 7777, 7778, 7800, 7911, 7920, 7921, 7937, 7938, 7999, 8000, 8001, 8002, 8007, 8008, 8009, 8010, 8011, 8021, 8022, 8031, 8042, 8045, 8080, 8081, 8082, 8083, 8084, 8085, 8086, 8087, 8088, 8089, 8090, 8093, 8099, 8100, 8180, 8181, 8192, 8193, 8194, 8200, 8222, 8254, 8290, 8291, 8292, 8300, 8333, 8383, 8400, 8402, 8443, 8500, 8600, 8649, 8651, 8652, 8654, 8701, 8800, 8873, 8888, 8899, 8994, 9000, 9001, 9002, 9003, 9009, 9010, 9011, 9040, 9050, 9071, 9080, 9081, 9090, 9091, 9099, 9100, 9101, 9102, 9103, 9110, 9111, 9200, 9207, 9220, 9290, 9415, 9418, 9485, 9500, 9502, 9503, 9535, 9575, 9593, 9594, 9595, 9618, 9666, 9876, 9877, 9878, 9898, 9900, 9917, 9929, 9943, 9944, 9968, 9998, 9999, 10000, 10001, 10002, 10003, 10004, 10009, 10010, 10012, 10024, 10025, 10082, 10180, 10215, 10243, 10566, 10616, 10617, 10621, 10626, 10628, 10629, 10778, 11110, 11111, 11967, 12000, 12174, 12265, 12345, 13456, 13722, 13782, 13783, 14000, 14238, 14441, 14442, 15000, 15002, 15003, 15004, 15660, 15742, 16000, 16001, 16012, 16016, 16018, 16080, 16113, 16992, 16993, 17877, 17988, 18040, 18101, 18988, 19101, 19283, 19315, 19350, 19780, 19801, 19842, 20000, 20005, 20031, 20221, 20222, 20828, 21571, 22939, 23502, 24444, 24800, 25734, 25735, 26214, 27000, 27352, 27353, 27355, 27356, 27715, 28201, 30000, 30718, 30951, 31038, 31337, 32768, 32769, 32770, 32771, 32772, 32773, 32774, 32775, 32776, 32777, 32778, 32779, 32780, 32781, 32782, 32783, 32784, 32785, 33354, 33899, 34571, 34572, 34573, 35500, 38292, 40193, 40911, 41511, 42510, 44176, 44442, 44443, 44501, 45100, 48080, 49152, 49153, 49154, 49155, 49156, 49157, 49158, 49159, 49160, 49161, 49163, 49165, 49167, 49175, 49176, 49400, 49999, 50000, 50001, 50002, 50003, 50006, 50300, 50389, 50500, 50636, 50800, 51103, 51493, 52673, 52822, 52848, 52869, 54045, 54328, 55055, 55056, 55555, 55600, 56737, 56738, 57294, 57797, 58080, 60020, 60443, 61532, 61900, 62078, 63331, 64623, 64680, 65000, 65129, 65389}

var golbalfullports bool = false
var globalTimeout int = 10
var globalTCPReaderTimeout int = 5
var globalThread int = 100
var globalAddr = "127.0.0.1"

var nmapServicesFilename = "nmap-services.csv"
var nmapTCPTable = []string{}

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", globalAddr, p)

		conn, err := net.DialTimeout("tcp", address, time.Duration(globalTimeout)*time.Second)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func portscan(scanports []int) {
	log.Println("[+] Start port scan...")
	fmt.Printf("[+] Threads: %d\n", globalThread)
	fmt.Printf("[+] total scan %d ports\n", len(scanports))
	fmt.Printf("%-10v%-10v%-20v%s\n", "ports", "status", "service", "extrainfo")
	ports := make(chan int, globalThread)
	results := make(chan int)

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for _, p := range scanports {
			ports <- p
		}
	}()

	for i := 0; i < len(scanports); i++ {
		port := <-results
		if port != 0 {
			extraInfo := strings.TrimSuffix(serviceprobe(port), "\n")
			nmapTCPDefault := portServiceMap(port)
			fmt.Printf("%-10v%-10v%-20v%s\n", port, "open", nmapTCPDefault, extraInfo)
		}
	}

	close(ports)
	close(results)
	log.Println("[+] Finished\n\n")
}

func loadNMAPServicesTable() {
	fi := "nmap-services_tcp.csv"

	fileBytes, err := ioutil.ReadFile(fi)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	nmapTCPTable = strings.Split(string(fileBytes), "\n")
}

func portServiceMap(port int) string {

	key := strconv.Itoa(port) + "/tcp"
	nmapService := ""

	for i := 0; i < len(nmapTCPTable); i++ {
		split := strings.Fields(nmapTCPTable[i])
		if key == split[1] {
			nmapService = split[0]
			break
		}
	}

	return nmapService
}

func serviceprobe(port int) string {
	CONNECT := globalAddr + ":" + strconv.Itoa(port)
	timeoutDuration := time.Duration(globalTCPReaderTimeout) * time.Second

	c, err := net.DialTimeout("tcp", CONNECT, time.Duration(globalTimeout)*time.Second)

	if err != nil {
		fmt.Println(err)
		return "error"
	}

	fmt.Fprintf(c, "\n")
	c.SetReadDeadline(time.Now().Add(timeoutDuration))
	message, _ := bufio.NewReader(c).ReadString('\n')

	return message
}

func printLogo() {
	const version = "0.0.2"
	const logo = `
    ________        _________                       ._. ._.
   /  _____/  ____ /   _____/ ____ _____    ____    | | | |
  /   \  ___ /  _ \\_____  \_/ ___\\__  \  /    \   | | | |
  \    \_\  (  <_> )        \  \___ / __ \|   |  \   \|  \|
   \______  /\____/_______  /\___  >____  /___|  /   __  __
          \/              \/     \/     \/     \/    \/  \/
  `
	cGreen := color.C256(72)
	cBlue := color.C256(38)

	cGreen.Println(logo)
	cBlue.Printf("%59s\n", "by  F4l13n5n0w")
	cBlue.Printf("%53s %s\n", "version:", version)
}

func main() {

	printLogo()

	targetFile := flag.String("iL", "", "Input the file path of target IP list")
	var singleIP string
	flag.StringVar(&singleIP, "ip", "", "target IP address")
	var portsRange string
	flag.StringVar(&portsRange, "p", "", "ports range, split by ',' or use '-p -' to scan full 65535 ports, support keywords 'top10' and 'top100'. (default is top1k)")
	var nTimeout = flag.Int("st", 5, "TCP scan connection timeout (seconds)")
	var nTCPReaderTimeout = flag.Int("rt", 5, "TCP message reader timeout (seconds), this is for service detection")
	var nThreads = flag.Int("thread", 100, "set thread number, make sure not too high")

	flag.Parse()

	globalTimeout = *nTimeout
	globalTCPReaderTimeout = *nTCPReaderTimeout
	globalThread = *nThreads

	scanports := []int{}

	loadNMAPServicesTable()

	if portsRange != "" {
		switch portsRange {
		case "-":
			for i := 1; i <= 65535; i++ {
				scanports = append(scanports, i)
			}
		case "top10":
			scanports = top10ports[0:]
		case "top100":
			scanports = top100ports[0:]
		default:
			tmp := strings.Split(portsRange, ",")
			for _, p := range tmp {
				port, _ := strconv.Atoi(p)
				scanports = append(scanports, port)
			}
		}
	} else {
		scanports = top1kports[0:]
	}

	if singleIP != "" {
		globalAddr = singleIP
		fmt.Printf("\n\n[+] scan %s\n", singleIP)
		portscan(scanports)
		return
	}

	if *targetFile != "" {
		ipRange := []string{""}

		file, err := os.Open(*targetFile)
		if err != nil {
			fmt.Println(err)
		} else {
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()
				ipRange = append(ipRange, line)
			}
		}
		for _, ip := range ipRange[1:] {
			globalAddr = ip
			fmt.Printf("\n\n[+] scan %s\n", ip)
			portscan(scanports)
		}
		return
	}
}
