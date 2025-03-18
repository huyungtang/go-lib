package reader

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/huyungtang/go-lib/file"
	"golang.org/x/text/transform"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(filename string, opts ...file.Option) (r Reader, err error) {
	if isExist := file.IsExist(filename); isExist != file.IsFile {
		return nil, errors.New("file not found")
	}

	var f *os.File
	if f, err = os.OpenFile(filename, os.O_RDONLY, os.ModePerm); err != nil {
		return
	}

	return &reader{
		File:    f,
		Context: new(file.Context).ApplyOptions(opts),
	}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// reader *********************************************************************************************************************************
type reader struct {
	*os.File
	*file.Context
}

// reader *********************************************************************************************************************************
func (o *reader) reader() io.Reader {
	o.File.Seek(0, 0)

	if o.Context.Encoding == nil {
		return o.File
	}

	return transform.NewReader(o.File, o.Context.Encoding.NewDecoder().Transformer)
}

// Reader
// ****************************************************************************************************************************************
type Reader interface {
	ReadAll() ([]byte, error)
	Readln() *readerEvent
	Close() error
}

// ReadAll
// ****************************************************************************************************************************************
func (o *reader) ReadAll() (bs []byte, err error) {
	return io.ReadAll(o.reader())
}

// Readln
// ****************************************************************************************************************************************
func (o *reader) Readln() (evt *readerEvent) {
	evt = &readerEvent{
		Read:  make(chan *readerContext),
		EOF:   make(chan bool),
		Error: make(chan error),
	}

	go func(r io.Reader, c *readerEvent) {
		s := bufio.NewScanner(r)
		ln := uint64(0)

	LOOP:
		for {
			if s.Scan() {
				ln++
				c.Read <- &readerContext{LineNo: ln, Content: s.Text()}
			} else {
				if err := s.Err(); err != nil {
					c.Error <- err
				} else {
					c.EOF <- true
				}
				break LOOP
			}
		}
	}(o.reader(), evt)

	return evt
}

// Close
// ****************************************************************************************************************************************
func (o *reader) Close() (err error) {
	if o.File != nil {
		err = o.File.Close()
	}

	return
}

// readerEvent ****************************************************************************************************************************
type readerEvent struct {
	Read  chan *readerContext
	EOF   chan bool
	Error chan error
}

// readerContext **************************************************************************************************************************
type readerContext struct {
	LineNo  uint64
	Content string
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
