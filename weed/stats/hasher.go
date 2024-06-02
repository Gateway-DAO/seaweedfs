package stats

import (
	"hash"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/blake2b"
)

func Blake2b() (hash.Hash, error) { return blake2b.New256(nil) }

func hashFile(path string) ([]byte, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Create a new BLAKE2b hasher, here 256 bit for compatibility
	hasher, err := Blake2b()
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(hasher, file); err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func hashDirectory(directoryPath string) ([]byte, error) {
	// Create a new BLAKE2b hasher, here 256 bit for compatibility
	dirHasher, err := Blake2b()
	if err != nil {
		return nil, err
	}

	err = filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			fileHash, err := hashFile(path)
			if err != nil {
				return err
			}
			dirHasher.Write(fileHash)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirHasher.Sum(nil), nil
}
