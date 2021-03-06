/*
 * Simple GO interface test program
 * Copyright (c) 2017 struktur AG, Dirk Farin <farin@struktur.de>
 *
 * This file is part of heif, an example application using libheif.
 *
 * heif is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * heif is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with heif.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/strukturag/libheif/go/heif"
)

// ==================================================
//                      TEST
// ==================================================

func save_png(img image.Image, filename string) {
	var out bytes.Buffer
	if err := png.Encode(&out, img); err != nil {
		fmt.Printf("Could not encode image as PNG: %s\n", err)
	} else {
		if err := ioutil.WriteFile(filename, out.Bytes(), 0644); err != nil {
			fmt.Printf("Could not save PNG image as %s: %s\n", filename, err)
		} else {
			fmt.Printf("Written to %s\n", filename)
		}
	}
}

func test_heif_highlevel(filename string) {
	fmt.Printf("Performing highlevel conversion of %s\n", filename)
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not read file %s: %s\n", filename, err)
		return
	}
	defer file.Close()

	img, magic, err := image.Decode(file)
	if err != nil {
		fmt.Printf("Could not decode image: %s\n", err)
		return
	}

	fmt.Printf("Decoded image of type %s: %s\n", magic, img.Bounds())

	ext := filepath.Ext(filename)
	out_filename := filename[0:len(filename)-len(ext)] + "_highlevel.png"
	save_png(img, out_filename)
}

func test_heif_lowlevel(filename string) {
	fmt.Printf("Performing lowlevel conversion of %s\n", filename)
	c, err := heif.NewContext()
	if err != nil {
		fmt.Printf("Could not create context: %s\n", err)
		return
	}

	if err := c.ReadFromFile(filename); err != nil {
		fmt.Printf("Could not read file %s: %s\n", filename, err)
		return
	}

	var nImages = c.GetNumberofTopLevelImages()
	fmt.Printf("GetNumberofTopLevelImages: %v\n", nImages)

	var IDs = c.GetListOfTopLevelImageIDs()
	fmt.Printf("List of top level image IDs %s\n", IDs)

	if pID, err := c.GetPrimaryImageID(); err != nil {
		fmt.Printf("Could not get primary image id: %s\n", err)
	} else {
		fmt.Printf("Primary image: %v\n", pID)
	}

	handle, err := c.GetPrimaryImageHandle()
	if err != nil {
		fmt.Printf("Could not get primary image: %s\n", err)
		return
	}

	fmt.Printf("image size: %v %v\n", handle.GetWidth(), handle.GetHeight())

	if img, err := handle.DecodeImage(heif.ColorspaceUndefined,
		heif.ChromaUndefined,
		nil); err != nil {
		fmt.Printf("Could not decode image: %s\n", err)
	} else if i, err := img.GetImage(); err != nil {
		fmt.Printf("Could not get image: %s\n", err)
	} else {
		fmt.Printf("Rectangle: %v\n", i.Bounds())

		ext := filepath.Ext(filename)
		out_filename := filename[0:len(filename)-len(ext)] + "_lowlevel.png"
		save_png(i, out_filename)
	}
}

func main() {
	fmt.Printf("libheif version: %v\n", heif.GetVersion())
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s <filename>\n", os.Args[0])
		return
	}

	test_heif_lowlevel(os.Args[1])
	fmt.Println()
	test_heif_highlevel(os.Args[1])

	runtime.GC()

	time.Sleep(time.Second)
	fmt.Println("done.")
}
