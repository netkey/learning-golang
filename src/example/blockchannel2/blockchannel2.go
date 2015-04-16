//  文件块读入一个 channel,  再从 channel 写入另一文件

package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	//"runtime"
	"runtime/pprof"
	"time"
)

const (
	BlockSizeLimit    = 1024 * 1024 * 1
	ChannelBufferSize = 512
)

type dataBlock struct {
	BlockSequeue int    // 数据块的序号
	BlockSeek    int64  // 数据块在文件中的读取或写入位置点
	BlockSize    int    // 数据块的大小（真实数据大小）
	BlockData    []byte // 数据块的缓存
}

type DataBlockIO interface {

}
type BlockStreamIO interface {

}
type BlockStream struct {

}

func main() {

	// init all input var or parameter that used in program
	var (
		inFilename  = flag.String("src", "/Users/qinshen/vagrant/crossroads.mp4", "source fullfilename to be copy")
		outFilename = flag.String("dst", "/Users/qinshen/vagrant/cross_copy.mp4", "destion fullfilename be copy to")

		flagCpuprofile = flag.String("prof", "blockchannel2_porf", "runtime prof file")
		//delay   = flag.Duration("dst", 1*time.Second, "write delay")
	)
	flag.Parse()

	// 性能检查

	if *flagCpuprofile != "" {
		f, err := os.Create(*flagCpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	// 打开要读的文件，及要写入的文件
	inFile, outFile := os.Stdin, os.Stdout

	err := errors.New("")

	if len(*inFilename) > 0 && len(*outFilename) > 0 {
		//inFilename := "/Users/qinshen/vagrant/crossroads.mp4" //"/Users/qinshen/temp/GoogleChrome39.dmg" //
		if inFile, err = os.Open(*inFilename); err != nil {
			log.Fatal(err)
		}
		defer inFile.Close()

		//outFilename := "/Users/qinshen/vagrant/cross_copy.mp4" //"/Users/qinshen/temp/GoogleChrome39_copy.dmg" //
		if outFile, err = os.Create(*outFilename); err != nil {
			log.Fatal(err)
		}
		defer outFile.Close()
	} else {
		os.Exit(1)
	}

	inReader := bufio.NewReader(inFile)

	t0 := time.Now()

	// push data block to channel that the data block read from a file
	readBufferChan := blockReadStream(inReader)

	doneToWriteFile := blockWriteStream(readBufferChan, outFile)

	waitTofinish(doneToWriteFile)
	//defer close(doneToWriteFile)

	t1 := time.Now()
	fmt.Println()
	fmt.Println(" start at :", t0)
	fmt.Println(" end at   :", t1)
	fmt.Println(" lasted   :", t1.Sub(t0))
	fmt.Println()

}

// blockReadStream
// read data block and push to a channel
// return a channel
func blockReadStream(inReader io.Reader) <-chan dataBlock {
	// inFile *os.File
	// reader := bufio.NewReader(inFile)

	dataBlockChannel := make(chan dataBlock, ChannelBufferSize)

	go func() {

		var myBlockSequeue = 0
		var myBlockSeek int64 = 0

		for {
			//BlockSize := 0
			myBlockData := make([]byte, BlockSizeLimit)
			myBlockSize, err := inReader.Read(myBlockData)

			if err == io.EOF {
				err = nil // io.EOF isn't really an error

				break

			}
			/* else if err != nil {
				return //err //return nil, err // finish immediately for real errors
			}
			*/

			/*
			fmt.Println("block sequeue :", myBlockSequeue)
			fmt.Println("reader is read: ", myBlockSize)
			fmt.Println("file seek point :", myBlockSeek)
			*/
			fmt.Print(".")

			if myBlockSize > 0 {
				dataBlockChannel <- dataBlock{BlockSequeue: myBlockSequeue, BlockSeek: myBlockSeek, BlockSize: myBlockSize, BlockData: myBlockData}
			} else {
				break // file read to EOF
			}

			myBlockSequeue++
			myBlockSize64 := int64(myBlockSize)
			myBlockSeek += myBlockSize64

		}
		close(dataBlockChannel)
	}()

	return dataBlockChannel
}

// write data block to a file
// return a done channel
func blockWriteStream(dataBlockChannel <-chan dataBlock, outFile *os.File) <-chan bool {

	doneWrite := make(chan bool, 1)

	//var workers = runtime.NumCPU()
	//for i:=0; i < workers; i++ {
	go func() {
		for writeBlock := range dataBlockChannel {
			if _, err := outFile.WriteAt(writeBlock.BlockData[:writeBlock.BlockSize], writeBlock.BlockSeek); err != nil {

				break //return err
			}
			fmt.Print(">")

		}
		fmt.Println("\n file write finish ")

		outFile.Close()
		doneWrite <- true
	}()
	//}

	return doneWrite
}

// monitor a done channel
// return none
func waitTofinish(doneToWriteFile <-chan bool) {

	finishFlag := <-doneToWriteFile

	if finishFlag {
		fmt.Println("file copy finish.")
	}
	return
}
