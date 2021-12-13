package util

import (
    // "github.com/koyachi/go-nude"
    "github.com/nfnt/resize"
    "github.com/rwcarlsen/goexif/exif"
    "github.com/rwcarlsen/goexif/mknote"
    // qrcode "github.com/yeqown/go-qrcode"
    _ "image"
    _ "image/png"
// "mime/multipart"
    "image/jpeg"
    "errors"
    "os/exec"
    "strings"
    "context"
    "bytes"
    "time"
    // "fmt"
    "os"
    // "io"
)

func JpgValid(file *os.File) bool {
    // We only have read file header first 3 bytes
    head := make([]byte, 3)
    _, err := file.Read(head);
    // fmt.Println("JpgValid:Income BYTES:", head)
    if err != nil {
        return false
    }
    //.jpg: HEAD: [255 216 255]
    if head[0] != 255 && head[1] != 216 && head[2] != 255 {
        return false
    }
    return true
}

// Validations IsJPG, IsNude
func ImagefileValidations(fpath string) error {
    file, err := os.Open(fpath)
    if err != nil {
        return errors.New("Загруженная картинка не читается сайтом:)")
    }
    defer file.Close()
    // notNude, err := nude.IsNude(fpath)
    // if jpgOK := JpgValid(file); err != nil || jpgOK  || notNude {
    if jpgOK := JpgValid(file); err != nil || jpgOK {
        errors.New("Загруженная картинка не подходит для сайта")
    }
    return nil
}

// Remove unused file
func ImagefileRemove(fpath string) error {
    return os.Remove(fpath)
}

func ImagefileResizing(fpath, newfpath string, setwidth uint) error {
    file, err := os.Open(fpath)
    if err != nil {
// fmt.Println("RESIZE:err os.Open fpath")
        return err
    }    
    defer file.Close()
    img, err := jpeg.Decode(file)
    if err != nil {
// fmt.Println("RESIZE:err jpeg.Decode")
        return err
    }

//     m, _, err := image.DecodeConfig(file)
//     if err != nil {
// // fmt.Println("RESIZE:err DecodeConfig")
//         return err
//     }
// // fmt.Println("RESIZE:startSize:",m.Width)
//     if int(setwidth) >= m.Width {
        m := resize.Resize(setwidth, 0, img, resize.Lanczos3)
        out, err := os.Create(newfpath)
        if err != nil {
// fmt.Println("RESIZE:err create new path")
            return err
        }
        // Also deleting last file
        os.Remove(fpath)
        defer out.Close()
        jpeg.Encode(out, m, nil)
//     } else {
//         //  simple copy file
//             out, err := os.Create(newfpath)
//             if err != nil {
// // fmt.Println("RESIZE:err new path")
//                 return err
//             }
//             defer out.Close()
//             _, err = io.Copy(out, file)
//             if err != nil {
// // fmt.Println("RESIZE:err copy to new path")
//                 return err
//             }
//             return out.Close()
//     }
    return nil
}

// Parse exif data in image-bytes and get datetime string if exist.
// If not - get file datetime info. Output photo format is:
// IMG_20200301_175147.jpg
func GetDateTimeFromExif(imageBytes []byte) (string, error) {
    exif.RegisterParsers(mknote.All...)
    imgReader := bytes.NewReader(imageBytes)
    x, err := exif.Decode(imgReader)
    if err != nil {
        return "", err
    }
    dt, err := x.DateTime()
    if err != nil {
        return "", err
    }
    return dt.Format("20060102_150405"), nil
}

// Jpeg optimizing with jpegtran utilits then copy to -> userpic and maybe remove origfile
func ImagefileProgressiveOptimisation(ctx context.Context, fpath, targetpath string, removemeta bool) error {
    cmdTimeout := 3 * time.Second
    lsCmd := exec.Command("bash", "-c", "file " + fpath)
    lsOut, err := lsCmd.Output()
    if err != nil {
        return errors.New("Error exec: bash -c file" + fpath)
    }
    stripCommand := ""
    if removemeta {
        stripCommand = "--strip-all"
    } else {
        stripCommand = "--strip-none"
    }
    targetpathCommand := "--path=" + targetpath
    if targetpath == "" {
        targetpathCommand = ""
    }
    if !strings.Contains(string(lsOut), "progressive") {
        ctx, cancel := context.WithTimeout(ctx, cmdTimeout)
        defer cancel()
        if err := exec.CommandContext(ctx, "jpegoptim", stripCommand, "--all-progressive", "--quiet", "-ptm85", targetpathCommand, fpath).Run(); err != nil {
            return errors.New("Converting timeout for jpegoptim command occured")
        }
    }
    return nil
}

// func MakeQRfile(txtInfo, fpath string) error {
//     config := qrcode.Config{
//         EncMode: qrcode.EncModeByte,
//         EcLevel: qrcode.ErrorCorrectionQuart,
//     }
//     qrc, err := qrcode.NewWithConfig(txtInfo, &config, qrcode.WithQRWidth(9))
//     if err != nil {
//         return errors.New("could not generate QRCode", err)
//     }

//     // save file
//     if err := qrc.Save(fpath); err != nil {
//         return errors.New("could not save image", err)
//     }
//     return nil
// }