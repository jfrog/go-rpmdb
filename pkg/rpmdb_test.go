package rpmdb

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	_ "github.com/glebarez/go-sqlite"
)

func TestPackageList(t *testing.T) {
	tests := []struct {
		name    string
		file    string // Test input file
		pkgList []*PackageInfo
	}{
		{
			name:    "CentOS5 plain",
			file:    "testdata/centos5-plain/Packages",
			pkgList: CentOS5Plain(),
		},
		{
			name:    "CentOS6 Plain",
			file:    "testdata/centos6-plain/Packages",
			pkgList: CentOS6Plain(),
		},
		{
			name:    "CentOS6 with Development tools",
			file:    "testdata/centos6-devtools/Packages",
			pkgList: CentOS6DevTools(),
		},
		{
			name:    "CentOS6 with many packages",
			file:    "testdata/centos6-many/Packages",
			pkgList: CentOS6Many(),
		},
		{
			name:    "CentOS7 Plain",
			file:    "testdata/centos7-plain/Packages",
			pkgList: CentOS7Plain(),
		},
		{
			name:    "CentOS7 with Development tools",
			file:    "testdata/centos7-devtools/Packages",
			pkgList: CentOS7DevTools(),
		},
		{
			name:    "CentOS7 with many packages",
			file:    "testdata/centos7-many/Packages",
			pkgList: CentOS7Many(),
		},
		{
			name:    "CentOS7 with Python 3.5",
			file:    "testdata/centos7-python35/Packages",
			pkgList: CentOS7Python35(),
		},
		{
			name:    "CentOS7 with httpd 2.4",
			file:    "testdata/centos7-httpd24/Packages",
			pkgList: CentOS7Httpd24(),
		},
		{
			name:    "CentOS8 with modules",
			file:    "testdata/centos8-modularitylabel/Packages",
			pkgList: CentOS8Modularitylabel(),
		},
		{
			name:    "RHEL UBI8 from s390x",
			file:    "testdata/ubi8-s390x/Packages",
			pkgList: UBI8s390x(),
		},
		{
			name:    "SLE15 with NDB style rpm database",
			file:    "testdata/sle15-bci/Packages.db",
			pkgList: SLE15WithNDB(),
		},
		{
			name:    "Fedora35 with SQLite3 style rpm database",
			file:    "testdata/fedora35/rpmdb.sqlite",
			pkgList: Fedora35WithSQLite3(),
		},
		{
			name:    "Fedora35 plus MongoDB with SQLite3 style rpm database",
			file:    "testdata/fedora35-plus-mongo/rpmdb.sqlite",
			pkgList: Fedora35PlusMongoDBWithSQLite3(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.ListPackages()
			require.NoError(t, err)

			// They are tested in another function.
			for _, g := range got {
				g.PGP = ""
				g.DigestAlgorithm = 0
				g.InstallTime = 0
				g.BaseNames = nil
				g.DirIndexes = nil
				g.DirNames = nil
				g.FileSizes = nil
				g.FileDigests = nil
				g.FileModes = nil
				g.FileFlags = nil
				g.UserNames = nil
				g.GroupNames = nil
				g.Provides = nil
				g.Requires = nil
			}

			for i, p := range tt.pkgList {
				assert.Equal(t, p, got[i])
			}
		})
	}
}

func TestRpmDB_Package(t *testing.T) {
	tests := []struct {
		name                   string
		pkgName                string
		file                   string // Test input file
		want                   *PackageInfo
		wantInstalledFiles     []FileInfo
		wantInstalledFileNames []string
		wantErr                string
	}{
		{
			name:    "centos5 python",
			pkgName: "python",
			file:    "testdata/centos5-plain/Packages",
			want: &PackageInfo{
				Name:        "python",
				Version:     "2.4.3",
				Release:     "56.el5",
				Arch:        "x86_64",
				Size:        74377,
				SourceRpm:   "python-2.4.3-56.el5.src.rpm",
				License:     "PSF - see LICENSE",
				Vendor:      "CentOS",
				Summary:     "An interpreted, interactive, object-oriented programming language.",
				SigMD5:      "ebfb56be33b146ef39180a090e581258",
				PGP:         "",
				InstallTime: 1459411575,
				Provides: []string{
					"Distutils",
					"python(abi)",
					"python-abi",
					"python-x86_64",
					"python2",
					"python",
				},
				Requires: []string{
					"/usr/bin/env",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libdl.so.2()(64bit)",
					"libm.so.6()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpython2.4.so.1.0()(64bit)",
					"libutil.so.1()(64bit)",
					"python-libs-x86_64",
					"rpmlib(CompressedFileNames)",
					"rpmlib(PartialHardlinkSets)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(VersionedDependencies)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     CentOS5PythonInstalledFiles,
			wantInstalledFileNames: CentOS5PythonInstalledFileNames,
		},
		{
			name:    "centos6 glibc",
			pkgName: "glibc",
			file:    "testdata/centos6-plain/Packages",
			want: &PackageInfo{
				Name:            "glibc",
				Version:         "2.12",
				Release:         "1.212.el6",
				Arch:            "x86_64",
				Size:            13117447,
				SourceRpm:       "glibc-2.12-1.212.el6.src.rpm",
				License:         "LGPLv2+ and LGPLv2+ with exceptions and GPLv2+",
				Vendor:          "CentOS",
				Summary:         "The GNU libc libraries",
				SigMD5:          "89e843d7979a50a26e2ea1924ef3e213",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "RSA/SHA1, Wed Jun 20 11:36:27 2018, Key ID 0946fca2c105b9de",
				InstallTime:     1538857091,
				Provides: []string{
					"ANSI_X3.110.so()(64bit)",
					"ARMSCII-8.so()(64bit)",
					"ASMO_449.so()(64bit)",
					"BIG5.so()(64bit)",
					"BIG5HKSCS.so()(64bit)",
					"BRF.so()(64bit)",
					"CP10007.so()(64bit)",
					"CP1125.so()(64bit)",
					"CP1250.so()(64bit)",
					"CP1251.so()(64bit)",
					"CP1252.so()(64bit)",
					"CP1253.so()(64bit)",
					"CP1254.so()(64bit)",
					"CP1255.so()(64bit)",
					"CP1256.so()(64bit)",
					"CP1257.so()(64bit)",
					"CP1258.so()(64bit)",
					"CP737.so()(64bit)",
					"CP775.so()(64bit)",
					"CP932.so()(64bit)",
					"CSN_369103.so()(64bit)",
					"CWI.so()(64bit)",
					"DEC-MCS.so()(64bit)",
					"EBCDIC-AT-DE-A.so()(64bit)",
					"EBCDIC-AT-DE.so()(64bit)",
					"EBCDIC-CA-FR.so()(64bit)",
					"EBCDIC-DK-NO-A.so()(64bit)",
					"EBCDIC-DK-NO.so()(64bit)",
					"EBCDIC-ES-A.so()(64bit)",
					"EBCDIC-ES-S.so()(64bit)",
					"EBCDIC-ES.so()(64bit)",
					"EBCDIC-FI-SE-A.so()(64bit)",
					"EBCDIC-FI-SE.so()(64bit)",
					"EBCDIC-FR.so()(64bit)",
					"EBCDIC-IS-FRISS.so()(64bit)",
					"EBCDIC-IT.so()(64bit)",
					"EBCDIC-PT.so()(64bit)",
					"EBCDIC-UK.so()(64bit)",
					"EBCDIC-US.so()(64bit)",
					"ECMA-CYRILLIC.so()(64bit)",
					"EUC-CN.so()(64bit)",
					"EUC-JISX0213.so()(64bit)",
					"EUC-JP-MS.so()(64bit)",
					"EUC-JP.so()(64bit)",
					"EUC-KR.so()(64bit)",
					"EUC-TW.so()(64bit)",
					"GB18030.so()(64bit)",
					"GBBIG5.so()(64bit)",
					"GBGBK.so()(64bit)",
					"GBK.so()(64bit)",
					"GEORGIAN-ACADEMY.so()(64bit)",
					"GEORGIAN-PS.so()(64bit)",
					"GOST_19768-74.so()(64bit)",
					"GREEK-CCITT.so()(64bit)",
					"GREEK7-OLD.so()(64bit)",
					"GREEK7.so()(64bit)",
					"HP-GREEK8.so()(64bit)",
					"HP-ROMAN8.so()(64bit)",
					"HP-ROMAN9.so()(64bit)",
					"HP-THAI8.so()(64bit)",
					"HP-TURKISH8.so()(64bit)",
					"IBM037.so()(64bit)",
					"IBM038.so()(64bit)",
					"IBM1004.so()(64bit)",
					"IBM1008.so()(64bit)",
					"IBM1008_420.so()(64bit)",
					"IBM1025.so()(64bit)",
					"IBM1026.so()(64bit)",
					"IBM1046.so()(64bit)",
					"IBM1047.so()(64bit)",
					"IBM1097.so()(64bit)",
					"IBM1112.so()(64bit)",
					"IBM1122.so()(64bit)",
					"IBM1123.so()(64bit)",
					"IBM1124.so()(64bit)",
					"IBM1129.so()(64bit)",
					"IBM1130.so()(64bit)",
					"IBM1132.so()(64bit)",
					"IBM1133.so()(64bit)",
					"IBM1137.so()(64bit)",
					"IBM1140.so()(64bit)",
					"IBM1141.so()(64bit)",
					"IBM1142.so()(64bit)",
					"IBM1143.so()(64bit)",
					"IBM1144.so()(64bit)",
					"IBM1145.so()(64bit)",
					"IBM1146.so()(64bit)",
					"IBM1147.so()(64bit)",
					"IBM1148.so()(64bit)",
					"IBM1149.so()(64bit)",
					"IBM1153.so()(64bit)",
					"IBM1154.so()(64bit)",
					"IBM1155.so()(64bit)",
					"IBM1156.so()(64bit)",
					"IBM1157.so()(64bit)",
					"IBM1158.so()(64bit)",
					"IBM1160.so()(64bit)",
					"IBM1161.so()(64bit)",
					"IBM1162.so()(64bit)",
					"IBM1163.so()(64bit)",
					"IBM1164.so()(64bit)",
					"IBM1166.so()(64bit)",
					"IBM1167.so()(64bit)",
					"IBM12712.so()(64bit)",
					"IBM1364.so()(64bit)",
					"IBM1371.so()(64bit)",
					"IBM1388.so()(64bit)",
					"IBM1390.so()(64bit)",
					"IBM1399.so()(64bit)",
					"IBM16804.so()(64bit)",
					"IBM256.so()(64bit)",
					"IBM273.so()(64bit)",
					"IBM274.so()(64bit)",
					"IBM275.so()(64bit)",
					"IBM277.so()(64bit)",
					"IBM278.so()(64bit)",
					"IBM280.so()(64bit)",
					"IBM281.so()(64bit)",
					"IBM284.so()(64bit)",
					"IBM285.so()(64bit)",
					"IBM290.so()(64bit)",
					"IBM297.so()(64bit)",
					"IBM420.so()(64bit)",
					"IBM423.so()(64bit)",
					"IBM424.so()(64bit)",
					"IBM437.so()(64bit)",
					"IBM4517.so()(64bit)",
					"IBM4899.so()(64bit)",
					"IBM4909.so()(64bit)",
					"IBM4971.so()(64bit)",
					"IBM500.so()(64bit)",
					"IBM5347.so()(64bit)",
					"IBM803.so()(64bit)",
					"IBM850.so()(64bit)",
					"IBM851.so()(64bit)",
					"IBM852.so()(64bit)",
					"IBM855.so()(64bit)",
					"IBM856.so()(64bit)",
					"IBM857.so()(64bit)",
					"IBM860.so()(64bit)",
					"IBM861.so()(64bit)",
					"IBM862.so()(64bit)",
					"IBM863.so()(64bit)",
					"IBM864.so()(64bit)",
					"IBM865.so()(64bit)",
					"IBM866.so()(64bit)",
					"IBM866NAV.so()(64bit)",
					"IBM868.so()(64bit)",
					"IBM869.so()(64bit)",
					"IBM870.so()(64bit)",
					"IBM871.so()(64bit)",
					"IBM874.so()(64bit)",
					"IBM875.so()(64bit)",
					"IBM880.so()(64bit)",
					"IBM891.so()(64bit)",
					"IBM901.so()(64bit)",
					"IBM902.so()(64bit)",
					"IBM903.so()(64bit)",
					"IBM9030.so()(64bit)",
					"IBM904.so()(64bit)",
					"IBM905.so()(64bit)",
					"IBM9066.so()(64bit)",
					"IBM918.so()(64bit)",
					"IBM921.so()(64bit)",
					"IBM922.so()(64bit)",
					"IBM930.so()(64bit)",
					"IBM932.so()(64bit)",
					"IBM933.so()(64bit)",
					"IBM935.so()(64bit)",
					"IBM937.so()(64bit)",
					"IBM939.so()(64bit)",
					"IBM943.so()(64bit)",
					"IBM9448.so()(64bit)",
					"IEC_P27-1.so()(64bit)",
					"INIS-8.so()(64bit)",
					"INIS-CYRILLIC.so()(64bit)",
					"INIS.so()(64bit)",
					"ISIRI-3342.so()(64bit)",
					"ISO-2022-CN-EXT.so()(64bit)",
					"ISO-2022-CN.so()(64bit)",
					"ISO-2022-JP-3.so()(64bit)",
					"ISO-2022-JP.so()(64bit)",
					"ISO-2022-KR.so()(64bit)",
					"ISO-IR-197.so()(64bit)",
					"ISO-IR-209.so()(64bit)",
					"ISO646.so()(64bit)",
					"ISO8859-1.so()(64bit)",
					"ISO8859-10.so()(64bit)",
					"ISO8859-11.so()(64bit)",
					"ISO8859-13.so()(64bit)",
					"ISO8859-14.so()(64bit)",
					"ISO8859-15.so()(64bit)",
					"ISO8859-16.so()(64bit)",
					"ISO8859-2.so()(64bit)",
					"ISO8859-3.so()(64bit)",
					"ISO8859-4.so()(64bit)",
					"ISO8859-5.so()(64bit)",
					"ISO8859-6.so()(64bit)",
					"ISO8859-7.so()(64bit)",
					"ISO8859-8.so()(64bit)",
					"ISO8859-9.so()(64bit)",
					"ISO8859-9E.so()(64bit)",
					"ISO_10367-BOX.so()(64bit)",
					"ISO_11548-1.so()(64bit)",
					"ISO_2033.so()(64bit)",
					"ISO_5427-EXT.so()(64bit)",
					"ISO_5427.so()(64bit)",
					"ISO_5428.so()(64bit)",
					"ISO_6937-2.so()(64bit)",
					"ISO_6937.so()(64bit)",
					"JOHAB.so()(64bit)",
					"KOI-8.so()(64bit)",
					"KOI8-R.so()(64bit)",
					"KOI8-RU.so()(64bit)",
					"KOI8-T.so()(64bit)",
					"KOI8-U.so()(64bit)",
					"LATIN-GREEK-1.so()(64bit)",
					"LATIN-GREEK.so()(64bit)",
					"MAC-CENTRALEUROPE.so()(64bit)",
					"MAC-IS.so()(64bit)",
					"MAC-SAMI.so()(64bit)",
					"MAC-UK.so()(64bit)",
					"MACINTOSH.so()(64bit)",
					"MIK.so()(64bit)",
					"NATS-DANO.so()(64bit)",
					"NATS-SEFI.so()(64bit)",
					"PT154.so()(64bit)",
					"RK1048.so()(64bit)",
					"SAMI-WS2.so()(64bit)",
					"SHIFT_JISX0213.so()(64bit)",
					"SJIS.so()(64bit)",
					"T.61.so()(64bit)",
					"TCVN5712-1.so()(64bit)",
					"TIS-620.so()(64bit)",
					"TSCII.so()(64bit)",
					"UHC.so()(64bit)",
					"UNICODE.so()(64bit)",
					"UTF-16.so()(64bit)",
					"UTF-32.so()(64bit)",
					"UTF-7.so()(64bit)",
					"VISCII.so()(64bit)",
					"config(glibc)",
					"ld-linux-x86-64.so.2()(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.2.5)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.3)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.4)(64bit)",
					"ldconfig",
					"libBrokenLocale.so.1()(64bit)",
					"libBrokenLocale.so.1(GLIBC_2.2.5)(64bit)",
					"libCNS.so()(64bit)",
					"libGB.so()(64bit)",
					"libISOIR165.so()(64bit)",
					"libJIS.so()(64bit)",
					"libJISX0213.so()(64bit)",
					"libKSC.so()(64bit)",
					"libSegFault.so()(64bit)",
					"libanl.so.1()(64bit)",
					"libanl.so.1(GLIBC_2.2.5)(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.10)(64bit)",
					"libc.so.6(GLIBC_2.11)(64bit)",
					"libc.so.6(GLIBC_2.12)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.2.6)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.3)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.5)(64bit)",
					"libc.so.6(GLIBC_2.6)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libc.so.6(GLIBC_2.8)(64bit)",
					"libc.so.6(GLIBC_2.9)(64bit)",
					"libcidn.so.1()(64bit)",
					"libcrypt.so.1()(64bit)",
					"libcrypt.so.1(GLIBC_2.2.5)(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libdl.so.2(GLIBC_2.3.3)(64bit)",
					"libdl.so.2(GLIBC_2.3.4)(64bit)",
					"libm.so.6()(64bit)",
					"libm.so.6(GLIBC_2.2.5)(64bit)",
					"libm.so.6(GLIBC_2.4)(64bit)",
					"libmemusage.so()(64bit)",
					"libnsl.so.1()(64bit)",
					"libnsl.so.1(GLIBC_2.2.5)(64bit)",
					"libnss_compat.so.2()(64bit)",
					"libnss_dns.so.2()(64bit)",
					"libnss_files.so.2()(64bit)",
					"libnss_hesiod.so.2()(64bit)",
					"libnss_nis.so.2()(64bit)",
					"libnss_nisplus.so.2()(64bit)",
					"libpcprofile.so()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.11)(64bit)",
					"libpthread.so.0(GLIBC_2.12)(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0(GLIBC_2.2.6)(64bit)",
					"libpthread.so.0(GLIBC_2.3.2)(64bit)",
					"libpthread.so.0(GLIBC_2.3.3)(64bit)",
					"libpthread.so.0(GLIBC_2.3.4)(64bit)",
					"libpthread.so.0(GLIBC_2.4)(64bit)",
					"libresolv.so.2()(64bit)",
					"libresolv.so.2(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2(GLIBC_2.3.2)(64bit)",
					"libresolv.so.2(GLIBC_2.9)(64bit)",
					"librt.so.1()(64bit)",
					"librt.so.1(GLIBC_2.2.5)(64bit)",
					"librt.so.1(GLIBC_2.3.3)(64bit)",
					"librt.so.1(GLIBC_2.3.4)(64bit)",
					"librt.so.1(GLIBC_2.4)(64bit)",
					"librt.so.1(GLIBC_2.7)(64bit)",
					"libthread_db.so.1()(64bit)",
					"libthread_db.so.1(GLIBC_2.2.5)(64bit)",
					"libthread_db.so.1(GLIBC_2.3)(64bit)",
					"libthread_db.so.1(GLIBC_2.3.3)(64bit)",
					"libutil.so.1()(64bit)",
					"libutil.so.1(GLIBC_2.2.5)(64bit)",
					"rtld(GNU_HASH)",
					"glibc",
					"glibc(x86-64)",
				},
				Requires: []string{
					"/sbin/ldconfig",
					"/usr/sbin/glibc_post_upgrade.x86_64",
					"basesystem",
					"config(glibc)",
					"glibc-common",
					"ld-linux-x86-64.so.2()(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.2.5)(64bit)",
					"ld-linux-x86-64.so.2(GLIBC_2.3)(64bit)",
					"libBrokenLocale.so.1()(64bit)",
					"libCNS.so()(64bit)",
					"libGB.so()(64bit)",
					"libISOIR165.so()(64bit)",
					"libJIS.so()(64bit)",
					"libJISX0213.so()(64bit)",
					"libKSC.so()(64bit)",
					"libanl.so.1()(64bit)",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.3)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libcidn.so.1()(64bit)",
					"libcrypt.so.1()(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libfreebl3.so()(64bit)",
					"libfreebl3.so(NSSRAWHASH_3.12.3)(64bit)",
					"libgcc",
					"libm.so.6()(64bit)",
					"libnsl.so.1()(64bit)",
					"libnsl.so.1(GLIBC_2.2.5)(64bit)",
					"libnss_compat.so.2()(64bit)",
					"libnss_dns.so.2()(64bit)",
					"libnss_files.so.2()(64bit)",
					"libnss_hesiod.so.2()(64bit)",
					"libnss_nis.so.2()(64bit)",
					"libnss_nisplus.so.2()(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2()(64bit)",
					"libresolv.so.2(GLIBC_2.2.5)(64bit)",
					"libresolv.so.2(GLIBC_2.9)(64bit)",
					"librt.so.1()(64bit)",
					"libthread_db.so.1()(64bit)",
					"libutil.so.1()(64bit)",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PartialHardlinkSets)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(VersionedDependencies)",
					"rpmlib(PayloadIsXz)",
				},
			},
			wantInstalledFiles:     CentOS6GlibcInstalledFiles,
			wantInstalledFileNames: CentOS6GlibcInstalledFileNames,
		},
		{
			name:    "centos8 nodejs",
			pkgName: "nodejs",
			file:    "testdata/centos8-modularitylabel/Packages",
			want: &PackageInfo{
				Epoch:           intRef(1),
				Name:            "nodejs",
				Version:         "10.21.0",
				Release:         "3.module_el8.2.0+391+8da3adc6",
				Arch:            "x86_64",
				Size:            31483781,
				SourceRpm:       "nodejs-10.21.0-3.module_el8.2.0+391+8da3adc6.src.rpm",
				License:         "MIT and ASL 2.0 and ISC and BSD",
				Vendor:          "CentOS",
				Modularitylabel: "nodejs:10:8020020200707141642:6a468ee4",
				Summary:         "JavaScript runtime",
				SigMD5:          "bac7919c2369f944f9da510bbd01370b",
				PGP:             "RSA/SHA256, Tue Jul  7 16:08:24 2020, Key ID 05b555b38483c65d",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				InstallTime:     1606911097,
				Provides: []string{
					"bundled(brotli)",
					"bundled(c-ares)",
					"bundled(http-parser)",
					"bundled(icu)",
					"bundled(libuv)",
					"bundled(nghttp2)",
					"bundled(v8)",
					"nodejs",
					"nodejs(abi)",
					"nodejs(abi10)",
					"nodejs(engine)",
					"nodejs(v8-abi)",
					"nodejs(v8-abi6)",
					"nodejs(x86-64)",
					"nodejs-punycode",
					"npm(punycode)",
				},
				Requires: []string{
					"/bin/sh",
					"ca-certificates",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.14)(64bit)",
					"libc.so.6(GLIBC_2.15)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.28)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.2)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.6)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libc.so.6(GLIBC_2.9)(64bit)",
					"libcrypto.so.1.1()(64bit)",
					"libcrypto.so.1.1(OPENSSL_1_1_0)(64bit)",
					"libcrypto.so.1.1(OPENSSL_1_1_1)(64bit)",
					"libdl.so.2()(64bit)",
					"libdl.so.2(GLIBC_2.2.5)(64bit)",
					"libgcc_s.so.1()(64bit)",
					"libgcc_s.so.1(GCC_3.0)(64bit)",
					"libgcc_s.so.1(GCC_3.4)(64bit)",
					"libm.so.6()(64bit)",
					"libm.so.6(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0()(64bit)",
					"libpthread.so.0(GLIBC_2.2.5)(64bit)",
					"libpthread.so.0(GLIBC_2.3.2)(64bit)",
					"libpthread.so.0(GLIBC_2.3.3)(64bit)",
					"librt.so.1()(64bit)",
					"librt.so.1(GLIBC_2.2.5)(64bit)",
					"libssl.so.1.1()(64bit)",
					"libssl.so.1.1(OPENSSL_1_1_0)(64bit)",
					"libssl.so.1.1(OPENSSL_1_1_1)(64bit)",
					"libstdc++.so.6()(64bit)",
					"libstdc++.so.6(CXXABI_1.3)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.2)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.5)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.8)(64bit)",
					"libstdc++.so.6(CXXABI_1.3.9)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.11)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.14)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.15)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.18)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.20)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.21)(64bit)",
					"libstdc++.so.6(GLIBCXX_3.4.9)(64bit)",
					"libz.so.1()(64bit)",
					"npm",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
					"rpmlib(PayloadIsXz)",
					"rtld(GNU_HASH)",
				},
			},
			wantInstalledFiles:     CentOS8NodejsInstalledFiles,
			wantInstalledFileNames: CentOS8NodejsInstalledFileNames,
		},
		{
			name:    "CBL-Mariner 2.0 curl",
			pkgName: "curl",
			file:    "testdata/cbl-mariner-2.0/rpmdb.sqlite",
			want: &PackageInfo{
				Name:            "curl",
				Version:         "7.76.0",
				Release:         "6.cm2",
				Arch:            "x86_64",
				Size:            326023,
				SourceRpm:       "curl-7.76.0-6.cm2.src.rpm",
				License:         "MIT",
				Vendor:          "Microsoft Corporation",
				Summary:         "An URL retrieval utility and library",
				SigMD5:          "b5f5369ae91df3672fa3338669ec5ca2",
				DigestAlgorithm: PGPHASHALGO_SHA256,
				PGP:             "RSA/SHA256, Thu Jan 27 09:02:11 2022, Key ID 0cd9fed33135ce90",
				InstallTime:     1643279454,
				Provides: []string{
					"curl",
					"curl(x86-64)",
				},
				Requires: []string{
					"/bin/sh",
					"/sbin/ldconfig",
					"/sbin/ldconfig",
					"curl-libs",
					"krb5",
					"libc.so.6()(64bit)",
					"libc.so.6(GLIBC_2.14)(64bit)",
					"libc.so.6(GLIBC_2.2.5)(64bit)",
					"libc.so.6(GLIBC_2.3)(64bit)",
					"libc.so.6(GLIBC_2.3.4)(64bit)",
					"libc.so.6(GLIBC_2.33)(64bit)",
					"libc.so.6(GLIBC_2.34)(64bit)",
					"libc.so.6(GLIBC_2.4)(64bit)",
					"libc.so.6(GLIBC_2.7)(64bit)",
					"libcurl.so.4()(64bit)",
					"libssh2",
					"libz.so.1()(64bit)",
					"openssl",
					"rpmlib(CompressedFileNames)",
					"rpmlib(FileDigests)",
					"rpmlib(PayloadFilesHavePrefix)",
				},
			},
			wantInstalledFiles:     Mariner2CurlInstalledFiles,
			wantInstalledFileNames: Mariner2CurlInstalledFileNames,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(tt.file)
			require.NoError(t, err)

			got, err := db.Package(tt.pkgName)
			if tt.wantErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
				return
			}

			assert.NoError(t, err)

			gotInstalledFiles, err := got.InstalledFiles()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFiles, gotInstalledFiles)

			gotInstalledFileNames, err := got.InstalledFileNames()
			assert.NoError(t, err)
			assert.Equal(t, tt.wantInstalledFileNames, gotInstalledFileNames)

			// These fields are tested through InstalledFiles() above
			got.BaseNames = nil
			got.DirIndexes = nil
			got.DirNames = nil
			got.FileSizes = nil
			got.FileDigests = nil
			got.FileModes = nil
			got.FileFlags = nil
			got.UserNames = nil
			got.GroupNames = nil

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCorruptedPackage(t *testing.T) {
	db, err := Open("testdata/corrupted/Packages")
	require.NoError(t, err)
	_, err = db.ListPackages()
	assert.Error(t, err, "failed to parse")
}

func TestCorruptedPackageWithTimeout(t *testing.T) {
	db, err := Open("testdata/corrupted/Packages")
	require.NoError(t, err)
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Microsecond)
	defer cancel()
	_, err = db.ListPackagesWithContext(ctxWithTimeout)
	assert.Equal(t, "timeout for parse page", err.Error())
}

func TestCorruptedPackageWithContext(t *testing.T) {
	db, err := Open("testdata/corrupted/Packages")
	require.NoError(t, err)
	ctxWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err = db.ListPackagesWithContext(ctxWithTimeout)
	assert.Error(t, err, "failed to parse")
}
