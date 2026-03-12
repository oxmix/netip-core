package collector

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPackage_parseDebianLikeDpkg(t *testing.T) {
	content := []byte(`Package: wget
Status: install ok installed
Priority: standard
Section: web
Installed-Size: 3605
Maintainer: No  l K  the <noel@debian.org>
Architecture: amd64
Multi-Arch: foreign
Source: wget (1.21.3-1)
Version: 1.21.3-1+b2
Depends: libc6 (>= 2.33), libgnutls30 (>= 3.7.2), libidn2-0 (>= 0.6), libnettle8, libpcre2-8-0 (>= 10.22), libpsl5 (>= 0.16.0), libuuid1 (>= 2.16), zlib1>
Recommends: ca-certificates
Conflicts: wget-ssl
Conffiles:
 /etc/wgetrc c43064699caf6109f4b3da0405c06ebb
Description: retrieves files from the web
 Wget is a network utility to retrieve files from the web
 using HTTP(S) and FTP, the two most widely used internet
 protocols. It works non-interactively, so it will work in
 the background, after having logged off. The program supports
 recursive retrieval of web-authoring pages as well as FTP
 sites -- you can use Wget to make mirrors of archives and
 home pages or to travel the web like a WWW robot.
 .
 Wget works particularly well with slow or unstable connections
 by continuing to retrieve a document until the document is fully
 downloaded. Re-getting files from where it left off works on
 servers (both HTTP and FTP) that support it. Both HTTP and FTP
 retrievals can be time stamped, so Wget can see if the remote
 file has changed since the last retrieval and automatically
 retrieve the new version if it has.
 .
 Wget supports proxy servers; this can lighten the network load,
 speed up retrieval, and provide access behind firewalls.
Homepage: https://www.gnu.org/software/wget/

Package: zip
Status: install ok installed
Priority: optional
Section: utils
Installed-Size: 616
Maintainer: Santiago Vila <sanvila@debian.org>
Architecture: amd64
Multi-Arch: foreign
Version: 3.0-13
Depends: libbz2-1.0, libc6 (>= 2.34)
Recommends: unzip
Description: Archiver for .zip files
 This is InfoZIP's zip program. It produces files that are fully
 compatible with the popular PKZIP program; however, the command line
 options are not identical. In other words, the end result is the same,
 but the methods differ. :-)
 .
 This version supports encryption.
Homepage: https://infozip.sourceforge.net/Zip.html

`)

	file := filepath.Join(t.TempDir(), "status")
	err := os.WriteFile(file, content, 0600)
	if err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	ipk, err := parseDebianLikeDpkg(file)
	if err != nil {
		t.Fatalf("failed to parse debian dpkg file: %v", err)
	}
	if ipk[0].Name != "wget" && ipk[0].Version != "1.21.3-1+b2" {
		t.Fatal("expected wget version 1.21.3-1+b2, got ", ipk[0].Version)
	}
	if ipk[1].Name != "zip" && ipk[1].Version != "3.0-13" {
		t.Fatal("expected zip version 3.0-13, got ", ipk[1].Version)
	}
}
