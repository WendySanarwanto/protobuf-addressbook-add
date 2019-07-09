package main

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	pb "github.com/wendysanarwanto/protobuf-addressbook"
)

// Main reads the entire address book from a file, adds one person based on
// user input, then writes it back out to the same file.
func main() {
	// Expecting a CLI Argument when running this program, to be the filename to save/load
	if len(os.Args) != 2 {
		log.Fatalf("[ERROR] Usage: %s ADDRESS_BOOK_FILE\n", os.Args[0])
	}
	fileName := os.Args[1]

	// Read the existing address book
	binContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("[ERROR] '%s' file does not exist. \n", fileName)
		} else {
			log.Fatalln("[ERROR] Error reading file: ", err)
		}
		return
	}

	// [START marshal_proto]
	addressBook := &pb.AddressBook{}
	// [START_EXCLUDE]
	// 'Unmarshall', parse the binary content back into protobuf Message
	err = proto.Unmarshal(binContent, addressBook)
	if err != nil {
		log.Fatalln("[ERROR] Failed to parse address book:", err)
		return
	}

	// Add a new contact into the address book
	newContact, err := promptForAddress(os.Stdin)
	if err != nil {
		log.Fatalln("[ERROR] Error with entering a new Contact:", err)
	}
	addressBook.People = append(addressBook.People, newContact)
	// [END_EXCLUDE]

	// Write the modified addressBook to disk
	updatedBinContent, err := proto.Marshal(addressBook)
	if err != nil {
		log.Fatalln("[ERROR] Failed to encode address book:", err)
	}
	writtenFileAccessPermission := os.FileMode(0644)
	err = ioutil.WriteFile(fileName, updatedBinContent, writtenFileAccessPermission)
	if err != nil {
		log.Fatalln("[ERROR] Failed to write address book: ", err)
	}
	// [END marshal_proto]
}
