package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// this is a value from the fltk lib that is used for scrolling to the bottom of
// the helpview
const MAX_TOP_LINE = 1000000

// Wraps around log.Println() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Log(v ...any) {
	// activityText = fmt.Sprintf("%v<p>%v</p>", activityText, fmt.Sprintf("%v", v...))
	// activity.SetValue(activityText)
	log.Println(v...)
	// activity.SetTopLine(MAX_TOP_LINE)
	// activity.SetTopLine(activity.TopLine() - activity.H()) // scroll to the bottom
}

// Wraps around log.Printf() as well as adding activity to the
// activity text buffer. Always adds a newline to the activity buffer.
func Logf(format string, v ...any) {
	format = fmt.Sprintf("%v\n", format)
	// activityText = fmt.Sprintf("%v<p>%v</p>", activityText, fmt.Sprintf(format, v...))
	// activity.SetValue(activityText)
	log.Printf(format, v...)
	// activity.SetTopLine(MAX_TOP_LINE)
	// activity.SetTopLine(activity.TopLine() - activity.H()) // scroll to the bottom
}

// Runs a command with the provided command (such as `/bin/sh`) and args (such
// as ["-c","'echo hello'"]) and environment variables (such as 'DISPLAY=:0').
//
// `stdin`, `stdout`, and `stderr` can all be `nil`.
//
// Returns the exit code of the command when it finishes.
func RunCommand(command string, args []string, env []string, stdin io.Reader, stdout, stderr io.Writer) (int, error) {
	cmd := exec.Command(command, args...)
	cmd.Env = append(cmd.Env, env...)

	if stdin != nil {
		cmd.Stdin = stdin
	}

	if stdout != nil {
		cmd.Stdout = stdout
	}

	if stderr != nil {
		cmd.Stderr = stdout
	}

	err := cmd.Run()
	if err != nil {
		return cmd.ProcessState.ExitCode(), err
	}

	return cmd.ProcessState.ExitCode(), nil
}

// Returns a minimum value of 0 if the provided integer is less than 0.
// Otherwise, returns the int itself.
func floorz(i int) int {
	if i < 0 {
		return 0
	}
	return i
}

// Returns the lesser of the two numbers, with a floor of zero.
func minz(a, b int) int {
	if a < b {
		return floorz(a)
	}
	return floorz(b)
}

/*
func encr(s, key string) (string, error) {
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return "", err
	}

	// GCM mode requires a nonce (number used once)
	nonce := make([]byte, 12) // GCM standard nonce size is 12 bytes
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Encrypt the plaintext
	ciphertext := gcm.Seal(nonce, nonce, []byte(s), nil)

	// Return the base64 encoded ciphertext
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decr(s, key string) (string, error) {
	// Decode the base64 encoded ciphertext
	ciphertextBytes, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "", err
	}

	// Create a new AES cipher
	keyb := []byte(key)
	block, err := aes.NewCipher(keyb)
	if err != nil {
		return "", err
	}

	// GCM mode requires a nonce (number used once)
	nonce, ciphertextBytes := ciphertextBytes[:12], ciphertextBytes[12:]

	// Create a GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// Decrypt the ciphertext
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// Retrieves the encrypted keys (secrets) from the provided map `s`, decrypts
// each of them, and puts their unencrypted values into the resulting map.
func getDecryptedSecrets(s map[string]string, key string) (map[string]string, error) {
	r := make(map[string]string)
	if s == nil {
		return r, fmt.Errorf("received nil map")
	}

	for k, v := range s {
		decrypted, err := decr(k, key)
		if err != nil {
			return r, fmt.Errorf("failed to decrypt: %v", err.Error())
		}

		r[decrypted] = v
	}

	return r, nil
}
*/

func saveConfig(fileName string, c *AppConfig) error {
	if fileName == "" {
		return fmt.Errorf("received empty config filename")
	}

	if c == nil {
		return fmt.Errorf("config was nil")
	}

	b, err := json.Marshal(*c)
	if err != nil {
		return fmt.Errorf("failed to marshal app config to yaml: %v", err.Error())
	} else {
		dir, _ := filepath.Split(fileName)
		err := os.MkdirAll(dir, 0o755)
		if err != nil {
			return fmt.Errorf("failed to create app config parent dir %v: %v", dir, err.Error())
		}
		err = os.WriteFile(fileName, b, 0o644)
		if err != nil {
			return fmt.Errorf("failed to save app config to %v: %v", fileName, err.Error())
		}
	}

	return nil
}
