package lesson3

import "errors"

type databaseClient interface {
	Read(id int) (string, error)
}

type fileSystemClient interface {
	Write(record string) error
}

type BatchJobRunner struct {
	databaseClient   databaseClient
	fileSystemClient fileSystemClient
}

func (b *BatchJobRunner) Run() error {
	return errors.New("the job failed ahhhhhh")
}

type brokenDatabaseClient struct{}

func (db *brokenDatabaseClient) Read(id int) (string, error) {
	return "", errors.New("ahhh something broke because nothing in life is fair")
}

type brokenFilesystemClient struct{}

func (fs *brokenFilesystemClient) Write(record string) error {
	return errors.New("wow, you just wrote everyone's SSN to a file")
}

func main() {
	// Instantiate our dependencies...
	databaseClient := &brokenDatabaseClient{}
	filesystemClient := &brokenFilesystemClient{}

	jobRunner := &BatchJobRunner{databaseClient: databaseClient, fileSystemClient: filesystemClient}

	// Run our job...
	jobRunner.Run()
}
