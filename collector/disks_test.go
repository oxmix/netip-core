package collector

import (
	"log"
	"testing"
)

func TestParseSmart(t *testing.T) {
	t.Parallel()

	smartExamples := []string{`smartctl 7.3 2022-02-28 r5338 [x86_64-linux-6.1.0-13-amd64] (local build)
Copyright (C) 2002-22, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Model Family:     Seagate Enterprise Capacity 3.5 HDD
Device Model:     ST10000NM0016-1TT101
Serial Number:    ZA27LJ8H
LU WWN Device Id: 5 000c50 0b1a84875
Firmware Version: SND0
User Capacity:    10,000,831,348,736 bytes [10.0 TB]
Sector Sizes:     512 bytes logical, 4096 bytes physical
Rotation Rate:    7200 rpm
Form Factor:      3.5 inches
Device is:        In smartctl database 7.3/5319
ATA Version is:   ACS-3 T13/2161-D revision 5
SATA Version is:  SATA 3.1, 6.0 Gb/s (current: 6.0 Gb/s)
Local Time is:    Sat Feb 24 21:32:27 2024 UTC
SMART support is: Available - device has SMART capability.
SMART support is: Enabled

=== START OF READ SMART DATA SECTION ===
SMART overall-health self-assessment test result: PASSED

General SMART Values:
Offline data collection status:  (0x82)	Offline data collection activity
					was completed without error.
					Auto Offline Data Collection: Enabled.
Self-test execution status:      (   0)	The previous self-test routine completed
					without error or no self-test has ever 
					been run.
Total time to complete Offline 
data collection: 		(  567) seconds.
Offline data collection
capabilities: 			 (0x7b) SMART execute Offline immediate.
					Auto Offline data collection on/off support.
					Suspend Offline collection upon new
					command.
					Offline surface scan supported.
					Self-test supported.
					Conveyance Self-test supported.
					Selective Self-test supported.
SMART capabilities:            (0x0003)	Saves SMART data before entering
					power-saving mode.
					Supports SMART auto save timer.
Error logging capability:        (0x01)	Error logging supported.
					General Purpose Logging supported.
Short self-test routine 
recommended polling time: 	 (   1) minutes.
Extended self-test routine
recommended polling time: 	 ( 843) minutes.
Conveyance self-test routine
recommended polling time: 	 (   2) minutes.
SCT capabilities: 	       (0x50bd)	SCT Status supported.
					SCT Error Recovery Control supported.
					SCT Feature Control supported.
					SCT Data Table supported.

SMART Attributes Data Structure revision number: 10
Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  1 Raw_Read_Error_Rate     0x000f   081   064   044    Pre-fail  Always       -       116617481
  3 Spin_Up_Time            0x0003   093   091   000    Pre-fail  Always       -       0
  4 Start_Stop_Count        0x0032   100   100   020    Old_age   Always       -       72
  5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0
  7 Seek_Error_Rate         0x000f   095   060   045    Pre-fail  Always       -       3098627397
  9 Power_On_Hours          0x0032   060   060   000    Old_age   Always       -       35459 (77 134 0)
 10 Spin_Retry_Count        0x0013   100   100   097    Pre-fail  Always       -       0
 12 Power_Cycle_Count       0x0032   100   100   020    Old_age   Always       -       72
184 End-to-End_Error        0x0032   100   100   099    Old_age   Always       -       0
187 Reported_Uncorrect      0x0032   100   100   000    Old_age   Always       -       0
188 Command_Timeout         0x0032   100   100   000    Old_age   Always       -       0 0 0
189 High_Fly_Writes         0x003a   098   098   000    Old_age   Always       -       2
190 Airflow_Temperature_Cel 0x0022   056   045   040    Old_age   Always       -       44 (Min/Max 36/46)
191 G-Sense_Error_Rate      0x0032   099   099   000    Old_age   Always       -       3872
192 Power-Off_Retract_Count 0x0032   100   100   000    Old_age   Always       -       154
193 Load_Cycle_Count        0x0032   075   075   000    Old_age   Always       -       51831
194 Temperature_Celsius     0x0022   044   055   000    Old_age   Always       -       44 (0 12 0 0 0)
195 Hardware_ECC_Recovered  0x001a   030   001   000    Old_age   Always       -       116617481
197 Current_Pending_Sector  0x0012   100   100   000    Old_age   Always       -       0
198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0
240 Head_Flying_Hours       0x0000   100   253   000    Old_age   Offline      -       29306h+32m+33.331s
241 Total_LBAs_Written      0x0000   100   253   000    Old_age   Offline      -       126261969900
242 Total_LBAs_Read         0x0000   100   253   000    Old_age   Offline      -       805789028336

SMART Error Log Version: 1
No Errors Logged

SMART Self-test log structure revision number 1
No self-tests have been logged.  [To run self-tests, use: smartctl -t]

SMART Selective self-test log data structure revision number 1
 SPAN  MIN_LBA  MAX_LBA  CURRENT_TEST_STATUS
    1        0        0  Not_testing
    2        0        0  Not_testing
    3        0        0  Not_testing
    4        0        0  Not_testing
    5        0        0  Not_testing
Selective self-test flags (0x0):
  After scanning selected spans, do NOT read-scan remainder of disk.
If Selective self-test is pending on power-up, resume after 0 minute delay.
`,

		`smartctl 7.3 2022-02-28 r5338 [x86_64-linux-6.1.0-13-amd64] (local build)
Copyright (C) 2002-22, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Model Family:     Seagate Exos X16
Device Model:     ST14000NM001G-2KJ103
Serial Number:    ZL2ALN08
LU WWN Device Id: 5 000c50 0c7b680fe
Firmware Version: SN03
User Capacity:    14,000,519,643,136 bytes [14.0 TB]
Sector Sizes:     512 bytes logical, 4096 bytes physical
Rotation Rate:    7200 rpm
Form Factor:      3.5 inches
Device is:        In smartctl database 7.3/5319
ATA Version is:   ACS-4 (minor revision not indicated)
SATA Version is:  SATA 3.3, 6.0 Gb/s (current: 6.0 Gb/s)
Local Time is:    Sat Feb 24 21:32:27 2024 UTC
SMART support is: Available - device has SMART capability.
SMART support is: Enabled

=== START OF READ SMART DATA SECTION ===
SMART overall-health self-assessment test result: PASSED

General SMART Values:
Offline data collection status:  (0x82)	Offline data collection activity
					was completed without error.
					Auto Offline Data Collection: Enabled.
Self-test execution status:      (   0)	The previous self-test routine completed
					without error or no self-test has ever 
					been run.
Total time to complete Offline 
data collection: 		(  575) seconds.
Offline data collection
capabilities: 			 (0x7b) SMART execute Offline immediate.
					Auto Offline data collection on/off support.
					Suspend Offline collection upon new
					command.
					Offline surface scan supported.
					Self-test supported.
					Conveyance Self-test supported.
					Selective Self-test supported.
SMART capabilities:            (0x0003)	Saves SMART data before entering
					power-saving mode.
					Supports SMART auto save timer.
Error logging capability:        (0x01)	Error logging supported.
					General Purpose Logging supported.
Short self-test routine 
recommended polling time: 	 (   1) minutes.
Extended self-test routine
recommended polling time: 	 (1265) minutes.
Conveyance self-test routine
recommended polling time: 	 (   2) minutes.
SCT capabilities: 	       (0x70bd)	SCT Status supported.
					SCT Error Recovery Control supported.
					SCT Feature Control supported.
					SCT Data Table supported.

SMART Attributes Data Structure revision number: 10
Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  1 Raw_Read_Error_Rate     0x000f   083   064   044    Pre-fail  Always       -       207811977
  3 Spin_Up_Time            0x0003   092   090   000    Pre-fail  Always       -       0
  4 Start_Stop_Count        0x0032   100   100   020    Old_age   Always       -       44
  5 Reallocated_Sector_Ct   0x0033   100   100   010    Pre-fail  Always       -       0
  7 Seek_Error_Rate         0x000f   073   060   045    Pre-fail  Always       -       22037534
  9 Power_On_Hours          0x0032   088   088   000    Old_age   Always       -       11216
 10 Spin_Retry_Count        0x0013   100   100   097    Pre-fail  Always       -       0
 12 Power_Cycle_Count       0x0032   100   100   020    Old_age   Always       -       44
 18 Head_Health             0x000b   100   100   050    Pre-fail  Always       -       0
187 Reported_Uncorrect      0x0032   100   100   000    Old_age   Always       -       0
188 Command_Timeout         0x0032   100   100   000    Old_age   Always       -       0
190 Airflow_Temperature_Cel 0x0022   063   050   040    Old_age   Always       -       37 (Min/Max 33/37)
192 Power-Off_Retract_Count 0x0032   100   100   000    Old_age   Always       -       23
193 Load_Cycle_Count        0x0032   086   086   000    Old_age   Always       -       29794
194 Temperature_Celsius     0x0022   037   050   000    Old_age   Always       -       37 (0 23 0 0 0)
197 Current_Pending_Sector  0x0012   100   100   000    Old_age   Always       -       0
198 Offline_Uncorrectable   0x0010   100   100   000    Old_age   Offline      -       0
199 UDMA_CRC_Error_Count    0x003e   200   200   000    Old_age   Always       -       0
200 Pressure_Limit          0x0023   100   100   001    Pre-fail  Always       -       0
240 Head_Flying_Hours       0x0000   100   253   000    Old_age   Offline      -       955h+57m+40.755s
241 Total_LBAs_Written      0x0000   100   253   000    Old_age   Offline      -       45886802518
242 Total_LBAs_Read         0x0000   100   253   000    Old_age   Offline      -       488520849

SMART Error Log Version: 1
No Errors Logged

SMART Self-test log structure revision number 1
No self-tests have been logged.  [To run self-tests, use: smartctl -t]

SMART Selective self-test log data structure revision number 1
 SPAN  MIN_LBA  MAX_LBA  CURRENT_TEST_STATUS
    1        0        0  Not_testing
    2        0        0  Not_testing
    3        0        0  Not_testing
    4        0        0  Not_testing
    5        0        0  Not_testing
Selective self-test flags (0x0):
  After scanning selected spans, do NOT read-scan remainder of disk.
If Selective self-test is pending on power-up, resume after 0 minute delay.
`,
	}

	for k, smart := range smartExamples {
		s := (&Collector{}).parseSmart(smart)
		if s.Model == "" || s.Capacity == "" || s.Health == "" ||
			s.Working == 0 || s.Temperature == "" || s.Full == "" {
			log.Println(s.Model, "|", s.Capacity, "|", s.Health, "|", s.Working, "|", s.Temperature)
			log.Fatal("some filed is empty: arr key: ", k)
		}
	}
}

func TestParseSmartNVME(t *testing.T) {
	t.Parallel()

	smartExamples := []string{
		`smartctl 7.3 2022-02-28 r5338 [x86_64-linux-6.1.0-13-amd64] (local build)
Copyright (C) 2002-22, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Model Number:                       Samsung SSD 960 EVO 500GB
Serial Number:                      S3X4NB0K402977T
Firmware Version:                   3B7QCXE7
PCI Vendor/Subsystem ID:            0x144d
IEEE OUI Identifier:                0x002538
Total NVM Capacity:                 500,107,862,016 [500 GB]
Unallocated NVM Capacity:           0
Controller ID:                      2
NVMe Version:                       1.2
Number of Namespaces:               1
Namespace 1 Size/Capacity:          500,107,862,016 [500 GB]
Namespace 1 Utilization:            484,289,253,376 [484 GB]
Namespace 1 Formatted LBA Size:     512
Namespace 1 IEEE EUI-64:            002538 5481b027b7
Local Time is:                      Sat Feb 24 21:32:28 2024 UTC
Firmware Updates (0x16):            3 Slots, no Reset required
Optional Admin Commands (0x0007):   Security Format Frmw_DL
Optional NVM Commands (0x001f):     Comp Wr_Unc DS_Mngmt Wr_Zero Sav/Sel_Feat
Log Page Attributes (0x03):         S/H_per_NS Cmd_Eff_Lg
Maximum Data Transfer Size:         512 Pages
Warning  Comp. Temp. Threshold:     83 Celsius
Critical Comp. Temp. Threshold:     85 Celsius

Supported Power States
St Op     Max   Active     Idle   RL RT WL WT  Ent_Lat  Ex_Lat
 0 +     6.04W       -        -    0  0  0  0        0       0
 1 +     5.09W       -        -    1  1  1  1        0       0
 2 +     4.08W       -        -    2  2  2  2        0       0
 3 -   0.0400W       -        -    3  3  3  3      210    1500
 4 -   0.0050W       -        -    4  4  4  4     2200    6000

Supported LBA Sizes (NSID 0x1)
Id Fmt  Data  Metadt  Rel_Perf
 0 +     512       0         0

=== START OF SMART DATA SECTION ===
SMART Health Status: PASSED

SMART/Health Information (NVMe Log 0x02)
Critical Warning:                   0x00
Temperature:                        53 Celsius
Available Spare:                    100%
Available Spare Threshold:          10%
Percentage Used:                    9%
Data Units Read:                    74,461,804 [38.1 TB]
Data Units Written:                 101,693,959 [52.0 TB]
Host Read Commands:                 1,174,678,878
Host Write Commands:                2,139,546,706
Controller Busy Time:               14,207
Power Cycles:                       1,751
Power On Hours:                     27,967
Unsafe Shutdowns:                   567
Media and Data Integrity Errors:    0
Error Information Log Entries:      3,676
Warning  Comp. Temperature Time:    0
Critical Comp. Temperature Time:    0
Temperature Sensor 1:               53 Celsius
Temperature Sensor 2:               73 Celsius

Error Information (NVMe Log 0x01, 16 of 64 entries)
Num   ErrCount  SQId   CmdId  Status  PELoc          LBA  NSID    VS
  0       3676     0  0x5010  0x4004      -            0     0     -
  1       3675     0  0x0010  0x4004      -            0     0     -
  2       3674     0  0x0000  0x4004      -            0     0     -
  3       3673     0  0x000c  0x4004      -            0     0     -
  4       3672     0  0x0008  0x4004      -            0     0     -
  5       3671     0  0x0014  0x4004      -            0     0     -
  6       3670     0  0x0014  0x4004      -            0     0     -
  7       3669     0  0x0018  0x4004      -            0     0     -
  8       3668     0  0x0010  0x4004      -            0     0     -
  9       3667     0  0x0014  0x4004      -            0     0     -
 10       3666     0  0x001c  0x4004      -            0     0     -
 11       3665     0  0x5018  0x4004      -            0     0     -
 12       3664     0  0x001c  0x4004      -            0     0     -
 13       3663     0  0x0014  0x4004      -            0     0     -
 14       3662     0  0x000c  0x4004      -            0     0     -
 15       3661     0  0x0061  0x4212  0x028            0     0     -
... (48 entries not read)
`,

		`smartctl 7.3 2022-02-28 r5338 [x86_64-linux-6.8.0-52-generic] (local build)
Copyright (C) 2002-22, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Model Number:                       XPG GAMMIX S11 Pro
Serial Number:                      2J3720127651
Firmware Version:                   S0118C
PCI Vendor/Subsystem ID:            0x1cc1
IEEE OUI Identifier:                0x000000
Controller ID:                      1
NVMe Version:                       1.3
Number of Namespaces:               1
Namespace 1 Size/Capacity:          512,110,190,592 [512 GB]
Namespace 1 Formatted LBA Size:     512
Local Time is:                      Fri Apr 25 13:33:42 2025 UTC
Firmware Updates (0x14):            2 Slots, no Reset required
Optional Admin Commands (0x0016):   Format Frmw_DL Self_Test
Optional NVM Commands (0x005f):     Comp Wr_Unc DS_Mngmt Wr_Zero Sav/Sel_Feat Timestmp
Log Page Attributes (0x0f):         S/H_per_NS Cmd_Eff_Lg Ext_Get_Lg Telmtry_Lg
Maximum Data Transfer Size:         64 Pages
Warning  Comp. Temp. Threshold:     75 Celsius
Critical Comp. Temp. Threshold:     80 Celsius

Supported Power States
St Op     Max   Active     Idle   RL RT WL WT  Ent_Lat  Ex_Lat
 0 +     9.00W       -        -    0  0  0  0        0       0
 1 +     4.60W       -        -    1  1  1  1        0       0
 2 +     3.80W       -        -    2  2  2  2        0       0
 3 -   0.0450W       -        -    3  3  3  3     2000    2000
 4 -   0.0040W       -        -    4  4  4  4     6000    8000

Supported LBA Sizes (NSID 0x1)
Id Fmt  Data  Metadt  Rel_Perf
 0 +     512       0         0

=== START OF SMART DATA SECTION ===
SMART overall-health self-assessment test result: PASSED

SMART/Health Information (NVMe Log 0x02)
Critical Warning:                   0x00
Temperature:                        48 Celsius
Available Spare:                    100%
Available Spare Threshold:          10%
Percentage Used:                    10%
Data Units Read:                    61,647,004 [31.5 TB]
Data Units Written:                 53,891,465 [27.5 TB]
Host Read Commands:                 543,039,882
Host Write Commands:                1,032,487,568
Controller Busy Time:               23,630
Power Cycles:                       424
Power On Hours:                     48,353
Unsafe Shutdowns:                   211
Media and Data Integrity Errors:    0
Error Information Log Entries:      0
Warning  Comp. Temperature Time:    0
Critical Comp. Temperature Time:    0
Thermal Temp. 1 Transition Count:   105
Thermal Temp. 1 Total Time:         559

Error Information (NVMe Log 0x01, 16 of 256 entries)
No Errors Logged`,
	}

	for k, smart := range smartExamples {
		s := (&Collector{}).parseSmart(smart)
		if s.Model == "" || s.Capacity == "" || s.Health == "" || s.Used == "" ||
			s.Working == 0 || s.Temperature == "" || s.Full == "" {
			log.Println(s.Model, "|", s.Capacity, "|", s.Health, "|", s.Working, "|", s.Used, "|", s.Temperature)
			log.Fatal("some filed is empty: arr key: ", k)
		}
	}
}

func TestParseSmartLife231(t *testing.T) {
	t.Parallel()

	smartExamples := `smartctl 7.3 2022-02-28 r5338 [x86_64-linux-6.1.0-32-amd64] (local build)
Copyright (C) 2002-22, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Device Model:     SSDSC2BW180A4
Serial Number:    BTDA339
LU WWN Device Id: 0 000120 000000000
Firmware Version: 5.8.2
User Capacity:    192,048,095,232 bytes [192 GB]
Sector Size:      512 bytes logical/physical
Rotation Rate:    Solid State Device
Form Factor:      2.5 inches
TRIM Command:     Available
Device is:        Not in smartctl database 7.3/5319
ATA Version is:   ATA8-ACS, ACS-2 T13/2015-D revision 3
SATA Version is:  SATA 3.0, 6.0 Gb/s (current: 6.0 Gb/s)
Local Time is:    Fri Apr 25 13:14:14 2025 UTC
SMART support is: Available - device has SMART capability.
SMART support is: Enabled

=== START OF READ SMART DATA SECTION ===
SMART overall-health self-assessment test result: PASSED

General SMART Values:
Offline data collection status:  (0x00)	Offline data collection activity
					was never started.
					Auto Offline Data Collection: Disabled.
Self-test execution status:      (   0)	The previous self-test routine completed
					without error or no self-test has ever 
					been run.
Total time to complete Offline 
data collection: 		(    0) seconds.
Offline data collection
capabilities: 			 (0x7d) SMART execute Offline immediate.
					No Auto Offline data collection support.
					Abort Offline collection upon new
					command.
					Offline surface scan supported.
					Self-test supported.
					Conveyance Self-test supported.
					Selective Self-test supported.
SMART capabilities:            (0x0003)	Saves SMART data before entering
					power-saving mode.
					Supports SMART auto save timer.
Error logging capability:        (0x01)	Error logging supported.
					General Purpose Logging supported.
Short self-test routine 
recommended polling time: 	 (   1) minutes.
Extended self-test routine
recommended polling time: 	 (  48) minutes.
Conveyance self-test routine
recommended polling time: 	 (   2) minutes.
SCT capabilities: 	       (0x0025)	SCT Status supported.
					SCT Data Table supported.

SMART Attributes Data Structure revision number: 10
Vendor Specific SMART Attributes with Thresholds:
ID# ATTRIBUTE_NAME          FLAG     VALUE WORST THRESH TYPE      UPDATED  WHEN_FAILED RAW_VALUE
  1 Raw_Read_Error_Rate     0x0032   095   095   050    Old_age   Always       -       4339657
  5 Reallocated_Sector_Ct   0x0033   099   099   003    Pre-fail  Always       -       0
  9 Power_On_Hours          0x0032   070   070   000    Old_age   Always       -       27105 (7 34 0)
 12 Power_Cycle_Count       0x0032   100   100   000    Old_age   Always       -       281
171 Unknown_Attribute       0x000a   100   100   000    Old_age   Always       -       0
172 Unknown_Attribute       0x0032   100   100   000    Old_age   Always       -       0
174 Unknown_Attribute       0x0030   000   000   000    Old_age   Offline      -       180
177 Wear_Leveling_Count     0x0000   000   000   000    Old_age   Offline      -       233
181 Program_Fail_Cnt_Total  0x000a   100   100   000    Old_age   Always       -       0
182 Erase_Fail_Count_Total  0x0032   100   100   000    Old_age   Always       -       0
187 Reported_Uncorrect      0x0012   100   100   000    Old_age   Always       -       0
194 Temperature_Celsius     0x0022   070   114   000    Old_age   Always       -       70 (Min/Max 45/114)
195 Hardware_ECC_Recovered  0x001c   120   120   000    Old_age   Offline      -       4339657
196 Reallocated_Event_Count 0x0033   099   099   003    Pre-fail  Always       -       0
201 Unknown_SSD_Attribute   0x001c   120   120   000    Old_age   Offline      -       4339657
204 Soft_ECC_Correction     0x001c   120   120   000    Old_age   Offline      -       4339657
230 Unknown_SSD_Attribute   0x0013   100   100   000    Pre-fail  Always       -       100
231 Unknown_SSD_Attribute   0x0013   095   095   010    Pre-fail  Always       -       51539607553
233 Media_Wearout_Indicator 0x0032   000   000   000    Old_age   Always       -       19572
234 Unknown_Attribute       0x0032   000   000   000    Old_age   Always       -       2603
241 Total_LBAs_Written      0x0032   000   000   000    Old_age   Always       -       2603
242 Total_LBAs_Read         0x0032   000   000   000    Old_age   Always       -       2255

SMART Error Log not supported

SMART Self-test log structure revision number 1
No self-tests have been logged.  [To run self-tests, use: smartctl -t]

SMART Selective self-test log data structure revision number 1
 SPAN  MIN_LBA  MAX_LBA  CURRENT_TEST_STATUS
    1        0        0  Not_testing
    2        0        0  Not_testing
    3        0        0  Not_testing
    4        0        0  Not_testing
    5        0        0  Not_testing
Selective self-test flags (0x0):
  After scanning selected spans, do NOT read-scan remainder of disk.
If Selective self-test is pending on power-up, resume after 0 minute delay.
`

	s := (&Collector{}).parseSmart(smartExamples)
	if s.Model == "" || s.Capacity == "" || s.Health == "" || s.Used != "5%" ||
		s.Working == 0 || s.Temperature == "" || s.Full == "" {
		log.Println(s.Model, "|", s.Capacity, "|", s.Health, "|", s.Working, "|", s.Used, "|", s.Temperature)
		log.Fatal("some filed is empty or incorrect used")
	}
}

func TestParseMdStat(t *testing.T) {
	t.Parallel()

	data := `Personalities : [raid6] [raid5] [raid4] [linear] [multipath] [raid0] [raid1] [raid10] 
md0 : active raid5 sdb[1] sdd[4] sdc[2] sda[0]
      29298914304 blocks super 1.2 level 5, 512k chunk, algorithm 2 [4/4] [UUUU]
      bitmap: 2/73 pages [8KB], 65536KB chunk

md1 : active raid1 sda3[2] sdb3[1]
      1927689152 blocks super 1.2 [2/1] [_U]
      [=========>...........]  recovery = 49.1% (948001664/1927689152) finish=431.9min speed=37796K/sec

md123 : active raid5 sdb[0] sde[4] sdc[2] sda[1]
      29298914304 blocks super 1.2 level 5, 512k chunk, algorithm 2 [4/4] [UUUU]
      [==============\u003e......]  resync = 72.1% (7047093632/9766304768) finish=341.1min speed=132823K/sec
      bitmap: 7/73 pages [28KB], 65536KB chunk

unused devices: <none>
`

	pms := (&Collector{}).parseMdStat(data)
	if len(pms) != 3 {
		t.Fatalf("expected 3 got %d", len(pms))
	}
	if len(pms["md1"].Disks) != 2 {
		t.Fatalf("expected 2 got %d", len(pms["md1"].Disks))
	}
	if len(pms["md123"].Disks) != 4 {
		t.Fatalf("expected 4 got %d", len(pms["md123"].Disks))
	}
	if pms["md123"].Proc.State != "resync" {
		t.Fatalf("expected resync got %q", pms["md123"].Proc.State)
	}
	if pms["md123"].Proc.Progress != 72.1 {
		t.Fatalf("expected 72.1 got %v", pms["md123"].Proc.Progress)
	}
	if pms["md123"].Proc.Left != 20466 {
		t.Fatalf("expected 20466 got %v", pms["md123"].Proc.Left)
	}
	if pms["md123"].Proc.Speed != 132823 {
		t.Fatalf("expected 132823 got %v", pms["md123"].Proc.Speed)
	}
}

func TestParseMdAdm(t *testing.T) {
	t.Parallel()

	data := `/dev/md1:
           Version : 1.2
     Creation Time : Wed Jan 11 19:10:46 2023
        Raid Level : raid5
        Array Size : 29298914304 (27.29 TiB 30.00 TB)
     Used Dev Size : 9766304768 (9.10 TiB 10.00 TB)
      Raid Devices : 4
     Total Devices : 4
       Persistence : Superblock is persistent

     Intent Bitmap : Internal

       Update Time : Sat Feb 24 19:16:58 2024
             State : active 
    Active Devices : 4
   Working Devices : 4
    Failed Devices : 0
     Spare Devices : 0

            Layout : left-symmetric
        Chunk Size : 512K

Consistency Policy : bitmap

              Name : krl-2:1
              UUID : 4de0b7bd:e314de52:7e68e937:0b94e9b3
            Events : 123186

    Number   Major   Minor   RaidDevice State
       0       8        0        0      active sync   /dev/sda
       1       8       16        1      active sync   /dev/sdb
       2       8       32        2      active sync   /dev/sdc
       4       8       48        3      active sync   /dev/sdd
`

	p := (&Collector{}).parseMdAdm(data)
	if p.Name != "krl-2:1" {
		t.Fatal("expected krl-2:1 got", p.Name)
	}
	if p.State != "active" {
		t.Fatal("expected active got", p.State)
	}
	if p.Level != "RAID5" {
		t.Fatal("expected RAID5 got", p.Level)
	}
	if p.Capacity != "30.00 TB" {
		t.Fatal("expected 30.00 TB got", p.Capacity)
	}
	if p.CreatedAt != "Wed Jan 11 19:10:46 2023" {
		t.Fatal("expected Wed Jan 11 19:10:46 2023 got", p.CreatedAt)
	}
	if p.Active != "4" {
		t.Fatal("expected 4 got", p.Active)
	}
	if p.Working != "4" {
		t.Fatal("expected 4 got", p.Working)
	}
	if p.Failed != "0" {
		t.Fatal("expected 0 got", p.Failed)
	}
	if p.Spare != "0" {
		t.Fatal("expected 0 got", p.Spare)
	}
}

func TestParseMdAdm_Resync(t *testing.T) {
	t.Parallel()

	data := `/dev/md1:
           Version : 1.2
     Creation Time : Wed Jan 11 19:10:46 2023
        Raid Level : raid5
        Array Size : 29298914304 (27.29 TiB 30.00 TB)
     Used Dev Size : 9766304768 (9.10 TiB 10.00 TB)
      Raid Devices : 4
     Total Devices : 4
       Persistence : Superblock is persistent
       
     Intent Bitmap : Internal
     
       Update Time : Thu Apr  4 13:45:08 2024
             State : clean, resyncing 
    Active Devices : 4
   Working Devices : 4
    Failed Devices : 0
     Spare Devices : 0
     
            Layout : left-symmetric
        Chunk Size : 512K
        Consistency Policy : bitmap
        
     Resync Status : 72% complete
     
              Name : krl-2:1  (local to host krl-2)
              UUID : 4de0b7bd:e314de52:7e68e937:0b94e9b3
            Events : 129394
            
    Number   Major   Minor   RaidDevice State
       0       8       16        0      active sync   /dev/sdb
       1       8        0        1      active sync   /dev/sda
       2       8       32        2      active sync   /dev/sdc
       4       8       64        3      active sync   /dev/sde
`

	p := (&Collector{}).parseMdAdm(data)
	if p.Name != "krl-2:1" {
		t.Fatal("expected krl-2:1 got", p.Name)
	}
	if p.State != "clean, resyncing" {
		t.Fatal("expected clean, resyncing got", p.State)
	}
	if p.Level != "RAID5" {
		t.Fatal("expected RAID5 got", p.Level)
	}
	if p.Capacity != "30.00 TB" {
		t.Fatal("expected 30.00 TB got", p.Capacity)
	}
	if p.CreatedAt != "Wed Jan 11 19:10:46 2023" {
		t.Fatal("expected Wed Jan 11 19:10:46 2023 got", p.CreatedAt)
	}
	if p.Active != "4" {
		t.Fatal("expected 4 got", p.Active)
	}
	if p.Working != "4" {
		t.Fatal("expected 4 got", p.Working)
	}
	if p.Failed != "0" {
		t.Fatal("expected 0 got", p.Failed)
	}
	if p.Spare != "0" {
		t.Fatal("expected 0 got", p.Spare)
	}
}
