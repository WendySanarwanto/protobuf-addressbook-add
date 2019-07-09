package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	pb "github.com/wendysanarwanto/protobuf-addressbook"
)

func promptForAddress(ioReader io.Reader) (*pb.Person, error) {
	// Create a Person message protobuf. It can be created like any struct.
	newContact := &pb.Person{}

	// Create a new buffer reader, to hold data entered by user
	bufferReader := bufio.NewReader(ioReader)
	// Show input ID prompt to user
	fmt.Print("Enter person ID number: ")

	// An int32 field in the .proto file is represented as an int32 field
	// in the generated Go struct.
	_, err := fmt.Fscanf(bufferReader, "%d\n", &newContact.Id)
	if err != nil {
		return newContact, err
	}

	// Show input Name prompt to user
	fmt.Print("Enter name: ")
	name, err := bufferReader.ReadString('\n')
	if err != nil {
		return newContact, err
	}

	// A string field in the .proto file results in a string field in Go.
	// We trim the whitespace because rd.ReadString includes the trailing
	// newline character in its output.
	newContact.Name = strings.TrimSpace(name)

	// Show input email address (blank for none)
	fmt.Print("Enter email address (blank for none): ")
	email, err := bufferReader.ReadString('\n')
	if err != nil {
		return newContact, err
	}
	newContact.Email = strings.TrimSpace(email)

	for {
		// Capture phone for current entry
		fmt.Print("Enter a phone number (or leave blank to finish): ")
		phone, err := bufferReader.ReadString('\n')
		if err != nil {
			return newContact, err
		}
		phone = strings.TrimSpace(phone)
		if phone == "" {
			break
		}

		// The PhoneNumber message type is nested within the Person
		// message in the .proto file.  This results in a Go struct
		// named using the name of the parent prefixed to the name of
		// the nested message.  Just as with pb.Person, it can be
		// created like any other struct.
		phoneNumber := &pb.Person_PhoneNumber{
			Number: phone,
		}

		fmt.Print("Is this a mobile, home, or work phone ?")
		phoneType, err := bufferReader.ReadString('\n')
		if err != nil {
			return newContact, err
		}
		phoneType = strings.TrimSpace(phoneType)

		switch phoneType {
		case "mobile":
			phoneNumber.Type = pb.Person_MOBILE
		case "home":
			phoneNumber.Type = pb.Person_HOME
		case "work":
			phoneNumber.Type = pb.Person_WORK
		default:
			fmt.Printf("Unknown phone type '%q' Using default.\n", phoneType)
		}

		// A repeated proto field maps to a slice field in Go.  We can
		// append to it like any other slice.
		newContact.Phones = append(newContact.Phones, phoneNumber)
	}

	return newContact, nil
}
