// Copyright (C) 2019 Storj Labs, Inc.
// See LICENSE for copying information.

package piecestore

import (
	"bufio"
	"sync"

	"github.com/zeebo/errs"

	"storj.io/storj/pkg/pb"
)

// BufferedUpload implements buffering for an Upload.
type BufferedUpload struct {
	buffer bufio.Writer
	upload *Upload
}

// NewBufferedUpload creates buffered upload with the specified size.
func NewBufferedUpload(upload *Upload, size int) Uploader {
	buffered := &BufferedUpload{}
	buffered.upload = upload
	buffered.buffer = *bufio.NewWriterSize(buffered.upload, size)
	return buffered
}

// Write writes content to the buffer and flushes it to the upload once enough data has been gathered.
func (upload *BufferedUpload) Write(data []byte) (int, error) {
	return upload.buffer.Write(data)
}

// Cancel aborts the upload.
func (upload *BufferedUpload) Cancel() error {
	return upload.upload.Cancel()
}

// Commit flushes any remaining content from buffer and commits the upload.
func (upload *BufferedUpload) Commit() (*pb.PieceHash, error) {
	flushErr := upload.buffer.Flush()
	piece, closeErr := upload.upload.Commit()
	return piece, errs.Combine(flushErr, closeErr)
}

// BufferedDownload implements buffering for download.
type BufferedDownload struct {
	buffer   bufio.Reader
	download *Download
}

// NewBufferedDownload creates a buffered download with the specified size.
func NewBufferedDownload(download *Download, size int) Downloader {
	buffered := &BufferedDownload{}
	buffered.download = download
	buffered.buffer = *bufio.NewReaderSize(buffered.download, size)
	return buffered
}

// Read reads from the buffer and downloading in batches once it's empty.
func (download *BufferedDownload) Read(p []byte) (int, error) {
	return download.buffer.Read(p)
}

// Close closes the buffered download.
func (download *BufferedDownload) Close() error {
	return download.download.Close()
}

// LockingUpload adds a lock around upload making it safe to use concurrently.
// TODO: this shouldn't be needed.
type LockingUpload struct {
	mu     sync.Mutex
	upload Uploader
}

// Write uploads data.
func (upload *LockingUpload) Write(p []byte) (int, error) {
	upload.mu.Lock()
	defer upload.mu.Unlock()
	return upload.upload.Write(p)
}

// Cancel aborts the upload.
func (upload *LockingUpload) Cancel() error {
	upload.mu.Lock()
	defer upload.mu.Unlock()
	return upload.upload.Cancel()
}

// Commit finishes the upload.
func (upload *LockingUpload) Commit() (*pb.PieceHash, error) {
	upload.mu.Lock()
	defer upload.mu.Unlock()
	return upload.upload.Commit()
}

// LockingDownload adds a lock around download making it safe to use concurrently.
// TODO: this shouldn't be needed.
type LockingDownload struct {
	mu       sync.Mutex
	download Downloader
}

// Read downloads content.
func (download *LockingDownload) Read(p []byte) (int, error) {
	download.mu.Lock()
	defer download.mu.Unlock()
	return download.download.Read(p)
}

// Close closes the deownload.
func (download *LockingDownload) Close() error {
	download.mu.Lock()
	defer download.mu.Unlock()
	return download.download.Close()
}
