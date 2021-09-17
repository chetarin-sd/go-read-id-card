package main

// Require pcscd, libpcsclite
import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/ebfe/scard"
	"golang.org/x/text/encoding/charmap"
)

var (
	cmdReq            = []byte{0x00, 0xc0, 0x00, 0x00}
	cmdSelectThaiCard = []byte{0x00, 0xA4, 0x04, 0x00, 0x08, 0xA0, 0x00, 0x00, 0x00, 0x54, 0x48, 0x00, 0x01}
	cmdCID            = []byte{0x80, 0xb0, 0x00, 0x04, 0x02, 0x00, 0x0d}
	cmdTHFullname     = []byte{0x80, 0xb0, 0x00, 0x11, 0x02, 0x00, 0x64}
	cmdENFullname     = []byte{0x80, 0xb0, 0x00, 0x75, 0x02, 0x00, 0x64}
	cmdBirth          = []byte{0x80, 0xb0, 0x00, 0xD9, 0x02, 0x00, 0x08}
	cmdGender         = []byte{0x80, 0xb0, 0x00, 0xE1, 0x02, 0x00, 0x01}
	cmdIssuer         = []byte{0x80, 0xb0, 0x00, 0xF6, 0x02, 0x00, 0x64}
	cmdIssueDate      = []byte{0x80, 0xb0, 0x01, 0x67, 0x02, 0x00, 0x08}
	cmdExpireDate     = []byte{0x80, 0xb0, 0x01, 0x6F, 0x02, 0x00, 0x08}
	cmdAddress        = []byte{0x80, 0xb0, 0x15, 0x79, 0x02, 0x00, 0x64}
	cmdPhoto          = [][]byte{
		{0x80, 0xb0, 0x01, 0x7B, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x02, 0x7A, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x03, 0x79, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x04, 0x78, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x05, 0x77, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x06, 0x76, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x07, 0x75, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x08, 0x74, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x09, 0x73, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0A, 0x72, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0B, 0x71, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0C, 0x70, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0D, 0x6F, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0E, 0x6E, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x0F, 0x6D, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x10, 0x6C, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x11, 0x6B, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x12, 0x6A, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x13, 0x69, 0x02, 0x00, 0xFF},
		{0x80, 0xb0, 0x14, 0x68, 0x02, 0x00, 0xFF},
	}
)

func main() {
	// Create context with pcscd
	context, err := scard.EstablishContext()
	if err != nil {
		fmt.Println("Error EstablishContext:", err)
		return
	}

	// Release context after exit
	defer context.Release()

	// Get all reader
	readers, err := context.ListReaders()
	if err != nil {
		fmt.Println("Error ListReaders:", err)
		return
	}
	if len(readers) == 0 {
		fmt.Println("Error card readers not found.")
		return
	}

	fmt.Println("Readers: ")
	for rIdx, rItem := range readers {
		fmt.Println("Card reader ID: ", rIdx, "Item: ", rItem)
	}

	// Select reader
	readerIDx := 0
	if len(readers) > 1 {
		buf := bufio.NewReader(os.Stdin)
		fmt.Print("Select card reader ID[0]: ")
		readerIDxInput, err := buf.ReadString('\n')
		if err != nil {
			fmt.Println("Error select reader:", err)
			return
		}
		fmt.Println("Selected: ", readerIDxInput)
		readerIDx, err := strconv.Atoi(readerIDxInput)
		if err != nil {
			readerIDx = 0
		}
		if readerIDx < 0 || readerIDx > len(readers)-1 {
			fmt.Println("Error select reader: index out of bound")
			return
		}
	}
	reader := readers[readerIDx]

	// Connect to card
	card, err := context.Connect(reader, scard.ShareShared, scard.ProtocolAny)
	if err != nil {
		fmt.Println("Error Connect:", err)
		return
	}

	// Disconnect from card after exit
	defer card.Disconnect(scard.ResetCard)

	fmt.Println("Card status:")
	status, err := card.Status()
	if err != nil {
		fmt.Println("Error status:", err)
		return
	}
	fmt.Printf("\treader: %s\n\tstate: %x\n\tactive protocol: %x\n\tatr: % x\n",
		status.Reader, status.State, status.ActiveProtocol, status.Atr)

	atr, err := card.GetAttrib(scard.AttrAtrString)
	if err != nil {
		fmt.Println("Error GetAttrib:", err)
		return
	}

	// Get card attribute
	fmt.Println("Card ATR: ", string(atr))
	if atr[0] == 0x3B && atr[1] == 0x67 {
		cmdReq = []byte{0x00, 0xc0, 0x00, 0x01}
	} else {
		cmdReq = []byte{0x00, 0xc0, 0x00, 0x00}
	}
	fmt.Println("Req: ", cmdReq)

	// Select thai national ID card
	card.Transmit(cmdSelectThaiCard)

	dataTofile := ""

	cid, _ := getString(card, cmdCID, cmdReq)
	fmt.Println("cid: ", cid)
	dataTofile = "cid: " + cid + "\n"

	thFullname, _ := getString(card, cmdTHFullname, cmdReq)
	fmt.Println("thFullname: ", thFullname)
	dataTofile += "thFullName: " + thFullname + "\n"

	enFullname, _ := getString(card, cmdENFullname, cmdReq)
	fmt.Println("enFullname: ", enFullname)
	dataTofile += "enFullName: " + enFullname + "\n"

	dateOfBirth, _ := getString(card, cmdBirth, cmdReq)
	fmt.Println("dateOfBirth: ", dateOfBirth)
	dataTofile += "dateOfBirth: " + dateOfBirth + "\n"

	gender, _ := getString(card, cmdGender, cmdReq)
	fmt.Println("gender: ", gender)
	dataTofile += "gender: " + gender + "\n"

	issuer, _ := getString(card, cmdIssuer, cmdReq)
	fmt.Println("issuer: ", issuer)
	dataTofile += "issuer: " + issuer + "\n"

	issueDate, _ := getString(card, cmdIssueDate, cmdReq)
	fmt.Println("issueDate: ", issueDate)
	dataTofile += "issueDate: " + issueDate + "\n"

	expireDate, _ := getString(card, cmdExpireDate, cmdReq)
	fmt.Println("expireDate: ", expireDate)
	dataTofile += "expireDate: " + expireDate + "\n"

	address, _ := getString(card, cmdAddress, cmdReq)
	fmt.Println("address: ", address)
	dataTofile += "address: " + address + "\n"

	photo, _ := getPhoto(card, cmdReq)

	os.Mkdir("image", 0755)
	err = ioutil.WriteFile("./image/"+cid+".jpg", photo, 0664)
	if err != nil {
		fmt.Printf("Error write photo: %+v", err)
		return
	}

	writeFile(dataTofile)
}

func getString(card *scard.Card, cmd, req []byte) (resp string, err error) {
	rawResp, err := getData(card, cmd, cmdReq)
	if err != nil {
		fmt.Printf("getString: %+v", err)
		return resp, err
	}
	thResp, err := thaiToUnicode(rawResp)
	if err != nil {
		fmt.Printf("getString: %+v", err)
		return resp, err
	}

	// Remove unused bytes
	thResp = bytes.Trim(thResp, " ")

	return string(thResp), err
}

func thaiToUnicode(data []byte) (out []byte, err error) {
	decoder := charmap.Windows874.NewDecoder()
	out, err = decoder.Bytes(data)
	return out, err
}

func getData(card *scard.Card, cmd, req []byte) (resp []byte, err error) {
	// Send cmd
	_, err = card.Transmit(cmd)
	if err != nil {
		fmt.Printf("getData: %+v", err)
		return resp, err
	}
	// Send select cmd
	req = append(req, cmd[len(cmd)-1])
	resp, err = card.Transmit(req)
	if err != nil {
		fmt.Printf("getData: %+v", err)
		return resp, err
	}
	// Remove unused bytes
	resp = resp[:len(resp)-2]
	return resp, err
}

func getPhoto(card *scard.Card, req []byte) (resp []byte, err error) {
	for _, itemCmd := range cmdPhoto {
		tmpArray, err := getData(card, itemCmd, req)
		if err != nil {
			fmt.Printf("getPhoto: %+v", err)
		}
		resp = append(resp, tmpArray...)
	}
	return resp, err
}

func writeFile(text string) {
	val := text
	data := []byte(val)

	err := ioutil.WriteFile("data.txt", data, 0644)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Created file.")

}
