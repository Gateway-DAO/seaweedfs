package stats

import (
	"encoding/base64"
	"fmt"
	"hash"
	"io"
	"os"
	"path/filepath"
	"regexp"

	"github.com/gateway-dao/seaweedfs/weed/glog"
	"golang.org/x/crypto/blake2b"
)

type Hash []byte

func (h Hash) ToString() string {
	return base64.RawStdEncoding.EncodeToString(h)
}

func HashFromString(encodedHash string) (Hash, error) {
	return base64.RawStdEncoding.DecodeString(encodedHash)
}

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

func hashDirectory(directoryPath string) (map[string]Hash, error) {
	return hashFilteredDirectory(directoryPath, "*")
}

func hashFilteredDirectory(dirPath, filter string) (map[string]Hash, error) {
	var dirHashes map[string]Hash = map[string]Hash{}

	rPOSIX, r_err := regexp.CompilePOSIX(filter)
	if r_err != nil {
		return nil, fmt.Errorf("filter %s is an invalid regex pattern", filter)
	}

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			if rPOSIX.MatchString(path) {
				glog.V(4).Infof("hashing file in %s", path)
				fileHash, err := hashFile(path)
				if err != nil {
					return err
				}

				dirHash, hash_err := Blake2b()
				if hash_err != nil {
					return fmt.Errorf("error constructing hasher for dirHash @ path %s: %s", path, hash_err)
				}
				dirHash.Write(fileHash)
				dirHashes[path] = dirHash.Sum(nil)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return dirHashes, nil
}
