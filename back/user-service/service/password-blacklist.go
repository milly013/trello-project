package service

import (
    "bufio"
    "os"
    "log"
    "strings"
)

// Crna lista lozinki
var blacklistedPasswords = make(map[string]struct{})

// LoadBlacklistedPasswords učitava crne liste lozinki iz fajla.
func LoadBlacklistedPasswords(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        log.Fatalf("Greška prilikom otvaranja fajla za crnu listu lozinki: %v", err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        password := strings.TrimSpace(scanner.Text())
        if password != "" {
            blacklistedPasswords[password] = struct{}{}
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("Greška prilikom čitanja fajla za crnu listu lozinki: %v", err)
    }

    log.Println("Crna lista lozinki uspešno učitana.")
}

// IsPasswordBlacklisted proverava da li je lozinka na crnoj listi.
func IsPasswordBlacklisted(password string) bool {
    _, exists := blacklistedPasswords[password]
    return exists
}